package optimization

import (
	"accounting/data"
	"database/sql"

	// "fmt"
	// "sort"
	"time"

	//"log"
	"github.com/pkg/errors"
)

// Структура для выборки данных из БД.
type PersonTime struct {
	Person int16     // Id человека.
	Start  time.Time // Время начала временного промежутка.
	Finish time.Time // Время окончания. День совпадает со временем начала,
	// за исключением случая, когда заканчивается в 00:00:00 следующего дня.
	Detail int64 // Id детали.
	Number int32 // Номер операции.
}

// Непрерывный период.
type Period struct {
	Start  time.Time // Начало.
	Finish time.Time // Конец.
}

// Продолжительность периода.
func (p Period) Duration() int32 {
	return int32(p.Finish.Sub(p.Start).Seconds())
}

// Состояние
type ValidPersonTime int8

const (
	Wait   = 0 // Считали из БД. Ждет.
	InDB   = 1 // Считали из БД. Работает.
	GoWait = 2 // Есть изменения по сравнению с БД. Ждет.
	GoWork = 3 // Есть изменения по сравнению с БД. Работает.
)

// Период, над одной операцией.
type PeriodLisp struct {
	Duration int32           // Продолжительность. +
	Person   int16           // Ид человека +
	Detail   int64           // Ид Детали +
	Number   int32           // Номер операции. +
	Valid    ValidPersonTime // Состояние +
	Lisp     []Period        // Список непрерывных временных отрезков. +
}

// Разбиваем период на части и выделяем интервал,
// когда человек работает над данной деталью.
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
			Valid:    GoWait,
			Lisp:     begin,
		})
	}
	pl2 = append(pl2, &PeriodLisp{
		Duration: getDur(medium),
		Person:   pl.Person,
		Detail:   det.Id,
		Number:   det.Way[det.NowStage].OperationNumber,
		Valid:    GoWork,
		Lisp:     medium,
	})
	if len(end) != 0 {
		pl2 = append(pl2, &PeriodLisp{
			Duration: getDur(end),
			Person:   pl.Person,
			Detail:   0,
			Number:   0,
			Valid:    GoWait,
			Lisp:     end,
		})
	}
	return pl2
}

// Возращаем время окончания.
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

// Возращаем время начала.
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

// Возращаем продолжительность.
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

// Считываем информацию о расписании одного человека из БД.
func SelectPersonTime(db *sql.DB, Person *int16, Start, Finish *time.Time,
	Detail, Entity *int64, Number *int32) ([]PersonTime, error) {
	arr := make([]PersonTime, 0, 10)
	if err := (func() error {
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		QwStr := data.SelectPersonTime(Person, Start, Finish, Detail, Entity, Number)
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorQueryDB+QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			pt := PersonTime{}
			var detail, entity sql.NullInt64
			var number sql.NullInt32
			var start, finish []uint8
			err := rows.Scan(&pt.Person, &start, &finish, &detail, &entity, &number)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorDecryptRow)
			}
			if pt.Start, err = time.Parse(data.C.TimeLayoutMySql, string(start)); err != nil {
				return err // GO-TO error
			}
			if pt.Finish, err = time.Parse(data.C.TimeLayoutMySql, string(finish)); err != nil {
				return err // GO-TO error
			}
			pt.Detail, pt.Number = detail.Int64, number.Int32
			arr = append(arr, pt)
		}
		return nil
	}()); err != nil {
		return arr, errors.Wrapf(err, "In SelectPersonTime") //GO-TO строка
	}
	return arr, nil
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
					if pt[i].Detail != pt[j].Detail || pt[i].Number != pt[j].Number {
						break
					}
				}
				pl := PeriodLisp{
					Person: idPerson,
					Detail: pt[i].Detail,
					Number: pt[i].Number,
					Lisp:   make([]Period, 0, j-i),
					Valid:  Wait,
				}
				if pl.Detail != 0 || pl.Number != 0 {
					pl.Valid = InDB
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
		return make(map[int16]*PersonDB, 0), errors.Wrapf(err, "In SelectPerson") //GO-TO строка
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

// Выбираем допустимые периоды для конкретной операции над определенной деталью.
func (p *PersonDB) GetPeriod(det *DetailDB) []PeriodChoose {
	arr := make([]PeriodChoose, 0, len(p.TimeTable))
	for i, val := range p.TimeTable {
		if val.Valid == Wait {
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
