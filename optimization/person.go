package optimization

import (
	"accounting/data"
	. "accounting/data/constants"
	e "accounting/data/errors"
	l "accounting/data/log"
	"accounting/data/qwery"
	"database/sql"

	// "fmt"
	// "sort"
	"time"

	//"log"
	"github.com/pkg/errors"
)

// Состояние, описывающее периоды.
type StatePersonTime int8

const (
	PersonTimeWaitDB   = 1 // Считали из БД. Ждет.
	PersonTimeWorkDB   = 2 // Считали из БД. Работает.
	PersonTimeWaitNoDB = 3 // Есть изменения по сравнению с БД. Ждет.
	PersonTimeWorkNoDB = 4 // Есть изменения по сравнению с БД. Работает.
)

// Ожидает ли рабочий в данный временной период.
func (pt StatePersonTime) IsWait() bool {
	return pt == PersonTimeWaitDB || pt == PersonTimeWaitNoDB
}

// Структура для выборки непрерывных временных интервалов из БД.
// Во время непрерывного интервала, человек не может отойти, а потом опять
// вернуться к выполнению операции. У одной операции над деталью может быть
// больше одного непрерывного временного интервала.
type PersonTime struct {
	Person int16           // Id рабочего.
	Start  time.Time       // Время начала временного интервала.
	Finish time.Time       // Время окончания временного интервала.
	Detail int64           // Id детали.
	Number int32           // Номер операции.
	State  StatePersonTime // Статус.
}

// Считываем информацию о непрерывных интервалах рабочих из БД.
func SelectPersonTime(db *sql.DB, Person *int16, Start, Finish *time.Time,
	Detail, Entity *int64, Number *int32) ([]PersonTime, error) {
	arr := make([]PersonTime, 0, 10)
	if err := (func() error {
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, e.Err.ErrorPingDB)
		}
		QwStr := data.SelectPersonTime(Person, Start, Finish, Detail, Entity, Number)
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrapf(err, e.Err.ErrorQueryDB, QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			pt := PersonTime{}
			var detail, entity sql.NullInt64
			var number sql.NullInt32
			var start, finish []uint8
			err := rows.Scan(&pt.Person, &start, &finish, &detail, &entity, &number)
			if err != nil {
				return errors.Wrap(err, e.Err.ErrorDecryptRow)
			}
			if pt.Start, err = time.Parse(TimeLayoutMySql, string(start)); err != nil {
				return errors.Wrapf(err, e.Err.ErrorDecryptTime, string(start))
			}
			if pt.Finish, err = time.Parse(TimeLayoutMySql, string(finish)); err != nil {
				return errors.Wrapf(err, e.Err.ErrorDecryptTime, string(finish))
			}
			if detail.Valid && entity.Valid && number.Valid {
				pt.Detail, pt.Number, pt.State = detail.Int64, number.Int32, PersonTimeWorkDB
			} else {
				pt.State = PersonTimeWaitDB
			}
			arr = append(arr, pt)
		}
		return nil
	}()); err != nil {
		return arr, qwery.Wrapf(err, l.In.InSelectPersonTime,
			Person, Start, Finish, Detail, Entity, Number)
	}
	return arr, nil
}

// Непрерывный временной интервал.
type Period struct {
	Start  time.Time // Начало.
	Finish time.Time // Конец.
}

// Продолжительность непрерывного временного интервала.
func (p Period) Duration() int32 {
	return int32(p.Finish.Sub(p.Start).Seconds())
}

// Период, над одной операцией.
type PeriodLisp struct {
	Duration int32           // Продолжительность. +
	Person   int16           // Ид человека +
	Detail   int64           // Ид Детали +
	Number   int32           // Номер операции. +
	Valid    StatePersonTime // Состояние +
	Lisp     []Period        // Список непрерывных временных отрезков. +
}

// Разбиваем период на части и выделяем интервал,
// когда человек работает над данной деталью. // GO-TO возможно заменить окончание на продолжительность.
func (pl *PeriodLisp) Change(start, finish time.Time, det *DetailDB) []*PeriodLisp {
	begin := make([]Period, 0, 10)
	medium := make([]Period, 0, 10)
	end := make([]Period, 0, 10)
	for _, val := range pl.Lisp {
		if val.Finish.Before(start) || val.Finish.Equal(start) {
			begin = append(begin, val)
		} else if val.Start.After(finish) || val.Start.Equal(finish) {
			end = append(end, val)
		} else if start.Before(val.Start) || start.Equal(val.Start) {
			if finish.Before(val.Finish) {
				medium = append(medium, Period{
					Start:  val.Start,
					Finish: finish,
				})
				end = append(end, Period{
					Start:  finish,
					Finish: val.Finish,
				})
			} else {
				medium = append(medium, val)
			}
		} else {
			if finish.Before(val.Finish) {
				begin = append(begin, Period{
					Start:  val.Start,
					Finish: start,
				})
				medium = append(medium, Period{
					Start:  start,
					Finish: finish,
				})
				end = append(end, Period{
					Start:  finish,
					Finish: val.Finish,
				})
			} else {
				begin = append(begin, Period{
					Start:  val.Start,
					Finish: start,
				})
				medium = append(medium, Period{
					Start:  start,
					Finish: finish,
				})
			}
		}
	}
	getDur := func(ps []Period) int32 {
		var dur int32 = 0
		for _, val := range ps {
			dur += val.Duration()
		}
		return dur
	}
	pl2 := make([]*PeriodLisp, 0, 3)
	if len(begin) != 0 {
		pl2 = append(pl2, &PeriodLisp{
			Duration: getDur(begin),
			Person:   pl.Person,
			Detail:   0,
			Number:   0,
			Valid:    PersonTimeWaitNoDB,
			Lisp:     begin,
		})
	}
	pl2 = append(pl2, &PeriodLisp{
		Duration: getDur(medium),
		Person:   pl.Person,
		Detail:   det.Id,
		Number:   det.Way[det.NowStage].OperationNumber,
		Valid:    PersonTimeWorkNoDB,
		Lisp:     medium,
	})
	if len(end) != 0 {
		pl2 = append(pl2, &PeriodLisp{
			Duration: getDur(end),
			Person:   pl.Person,
			Detail:   0,
			Number:   0,
			Valid:    PersonTimeWaitNoDB,
			Lisp:     end,
		})
	}
	return pl2
}

// По известному началу и продолжительности возращаем время окончания.
func (pl *PeriodLisp) GetFinish(start time.Time, duration int32) (time.Time, bool) {
	for i := 0; i < len(pl.Lisp); i++ {
		pi := pl.Lisp[i]
		if start.Before(pi.Start) {
			if dur := pi.Duration(); duration > dur {
				duration = duration - dur
			} else {
				return AddDur(pi.Start, duration), true
				// return pi.Start.Add(duration * time.Second), true
			}
		} else if start.Before(pi.Finish) {
			if dur := int32(pi.Finish.Sub(start).Seconds()); duration > dur {
				duration = duration - dur
			} else {
				return AddDur(start, duration), true
				// return start.Add(duration * time.Second), true
			}
		}
	}
	return start, false
}

// По известному окончанию и продолжительности возращаем время начала.
func (pl *PeriodLisp) GetStart(finish time.Time, duration int32) (time.Time, bool) {
	for i := len(pl.Lisp) - 1; i >= 0; i-- {
		pi := pl.Lisp[i]
		if pi.Finish.Before(finish) {
			if dur := pi.Duration(); duration > dur {
				duration = duration - dur
			} else {
				return AddDur(pi.Finish, -duration), true
				// return pi.Finish.Add(-duration * time.Second), true
			}
		} else if pi.Start.Before(finish) {
			if dur := int32(finish.Sub(pi.Start).Seconds()); duration > dur {
				duration = duration - dur
			} else {
				return AddDur(finish, -duration), true
				// return finish.Add(-duration * time.Second), true
			}
		}
	}
	return finish, false
}

// По известному началу и окончанию возращаем продолжительность.
func (pl *PeriodLisp) GetDuration(start, finish time.Time) int32 {
	var duration int32 = 0
	for i := 0; i < len(pl.Lisp); i++ {
		pi := pl.Lisp[i]
		if start.Before(pi.Start) {
			if pi.Finish.Before(finish) {
				duration = duration + pi.Duration()
			} else if pi.Start.Before(finish) {
				return duration + int32(finish.Sub(pi.Start).Seconds())
			}
		} else if start.Before(pi.Finish) {
			if pi.Finish.Before(finish) {
				duration = duration + int32(pi.Finish.Sub(start).Seconds())
			} else {
				return duration + int32(finish.Sub(start).Seconds())
			}
		}
	}
	return duration
}

// Храним информацию о человеке.
type PersonDB struct {
	Id        int16         // Ид человека.
	TimeTable []*PeriodLisp // Расписание.
}

// Считываем расписание всех людей из БД.
func SelectPerson(db *sql.DB, startPtr, finishPtr *time.Time) (map[int16]*PersonDB, error) {
	arr := make([]PersonDB, 0, 0)
	if err := (func() error {
		// Считать идентификаторы людей.
		persons, err := SelectId16(db, "Person")
		if err != nil {
			return err //GO-TO добавить обработку ошибок
		}
		if len(persons) == 0 {
			return errors.New("len(persons) == 0") //GO-TO добавить обработку ошибок
		}
		// arr = make([]PersonDB, 0, len(persons))

		// Для каждого человека считать расписание.
		for _, idPerson := range persons {
			pt, err := SelectPersonTime(db, &idPerson, startPtr, finishPtr, nil, nil, nil)
			if err != nil {
				return err //GO-TO добавить обработку ошибок
			}
			// sort.Slice(pt, func(i, j int) bool { //GO-TO сортировка в запросе БД.
			// 	return pt[i].Start.Before(pt[j].Start)
			// })
			periodLisps := make([]*PeriodLisp, 0, len(pt)/2)
			// Объединяем интервалы времени с одинаковой работой в один период.
			for i := 0; i < len(pt); {
				var j int
				for j = i + 1; j < len(pt); j++ {
					if pt[i].State != pt[j].State || pt[i].Detail != pt[j].Detail ||
						pt[i].Number != pt[j].Number {
						break
					}
				}
				pl := PeriodLisp{
					Person: idPerson,
					Detail: pt[i].Detail,
					Number: pt[i].Number,
					Lisp:   make([]Period, 0, j-i),
					Valid:  pt[i].State,
				}
				duration := 0.0
				for k := i; k < j; k++ {
					duration += pt[k].Finish.Sub(pt[k].Start).Seconds()
					pl.Lisp = append(pl.Lisp, Period{Start: pt[k].Start, Finish: pt[k].Finish})
				}
				pl.Duration = int32(duration)
				periodLisps = append(periodLisps, &pl)
				i = j
			}
			arr = append(arr, PersonDB{Id: idPerson, TimeTable: periodLisps})
		}
		return nil
	}()); err != nil {
		return make(map[int16]*PersonDB, 0), qwery.Wrapf(err, l.In.InSelectPerson, startPtr, finishPtr)
	}
	mp := make(map[int16]*PersonDB, len(arr))
	for i, person := range arr {
		mp[person.Id] = &arr[i]
	}
	return mp, nil
}

// Изменяем информацию о расписании человека.
func (p *PersonDB) Change(i int, start, finish time.Time, det *DetailDB) {
	pl2 := p.TimeTable[i].Change(start, finish, det) // GO-TO оптимизировать выделение памяти для массива.
	pl := p.TimeTable[0:i]
	pl = append(pl, pl2...)
	pl = append(pl, p.TimeTable[i+1:]...)
	p.TimeTable = pl
}

// Выбранный период для конкретной операции над определенной деталью.
type PeriodChoose struct {
	PL     *PeriodLisp // Периода в PersonDB.TimeTable.
	Index  int         // Индекс периода в PersonDB.TimeTable.
	Finish time.Time   // Окончание периода для данной детали.
}

// Выбираем допустимые периоды у данного рабочего для текущей операции над определенной деталью.
func (p *PersonDB) GetPeriod(det *DetailDB) []PeriodChoose {
	arr := make([]PeriodChoose, 0, len(p.TimeTable))
	for i, val := range p.TimeTable {
		if val.Valid.IsWait() {
			if finish, ok := val.GetFinish(det.NowTime(), det.NowDuration()); ok {
				arr = append(arr, PeriodChoose{
					PL:     val,
					Index:  i,
					Finish: finish,
				})
			}
		}
	}
	return arr
}
