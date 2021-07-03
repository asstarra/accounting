package optimization

import (
	"accounting/data"
	"database/sql"
	"fmt"
	"time"

	//"log"
	"github.com/pkg/errors"
)

type Waypoint struct {
	OperationId     int16 // Идентификатор типа операции (stage_number). +
	OperationNumber int32 // Номер операции в маршрутке (stage_name). +
	Duration        int32 // Продолжительность в секундах. +
	PersonCount     int8  // Количество человек. +

	Persons []int16   // Список идентификаторов людей. +
	Start   time.Time // Дата и время начала. +
	Finish  time.Time // Дата и время окончания. +

	Check bool // Выполнена ли операция. +
}

func (w Waypoint) Copy() Waypoint {
	persons := make([]int16, len(w.Persons), w.PersonCount)
	for i, val := range w.Persons {
		persons[i] = val
	}
	return Waypoint{
		OperationId:     w.OperationId,
		OperationNumber: w.OperationNumber,
		PersonCount:     w.PersonCount,
		Duration:        w.Duration,
		Persons:         persons,
		Start:           w.Start,
		Finish:          w.Finish,
	}
}

type StateDetail int8

const (
	No           = 0
	WaitChildren = 1
	InWork       = 2
	Ready        = 3
	NoWork       = 4
)

type DetailDB struct {
	Id     int64         // Идентификатор детали. +
	Entity int64         // Сущность. +
	State  StateDetail   // Состояние (пересчитывается каждый раз заново, не БД) +
	Start  time.Time     // Дата и время начала. +
	Finish time.Time     // Дата и время окончания. +
	Parent sql.NullInt64 // Родитель. +
	Way    []Waypoint    // Маршрутка. +

	NowStage int     // Текущая стадия. +
	Children []int64 // Индексы дочерних деталей. +

	Priority int32 // Приоритет.

	// Now time.Time // Текущее дата и время. + func
	// DetailCount  int     // Количество деталей в партии. ??? --
	// Priority     float64 // Приоритет выполнения.
	// WorkDuration float64 // Оставшаяся продолжительность. ??? +
	// Freetime     float64 // Оставшийся запас времени, которое деталь может простоять в очереди.
	// Durtime      float64 // Общая продолжительность на коэфицент. Для приритета. ???
}

func (d DetailDB) NowTime() time.Time {
	if d.NowStage == 0 {
		return d.Start
	}
	return d.Way[d.NowStage-1].Finish
}

func (d DetailDB) MyDuration() int32 {
	var dur int32 = 0
	for _, val := range d.Way {
		dur += val.Duration
	}
	return dur
}

func (d DetailDB) WorkDuration() int32 {
	var dur int32 = 0
	for i := d.NowStage; i < len(d.Way); i++ {
		dur += d.Way[i].Duration
	}
	return dur
}

func (d DetailDB) NowDuration() int32 {
	if d.NowStage < len(d.Way) {
		return d.Way[d.NowStage].Duration
	}
	return 0
}

func ChoosePersonIdsAndTime(db *sql.DB, detail, entity *int64,
	number *int32, personCount int8) (time.Time, time.Time, []int16, error) {
	arr := make([]int16, 0, personCount)
	var start, finish time.Time
	if err := (func() error {
		pt, err := SelectPersonTime(db, nil, nil, nil, detail, entity, number)
		if err != nil {
			return err // GO-TO error
		}
		if len(pt) == 0 {
			return nil
		}
		mp := make(map[int16]bool, personCount)
		start, finish = pt[0].Start, pt[0].Finish
		for _, val := range pt {
			if val.Start.Before(start) {
				start = val.Start
			}
			if finish.Before(val.Finish) {
				finish = val.Finish
			}
			mp[val.Person] = true
		}
		for key, _ := range mp {
			arr = append(arr, key)
		}
		return nil
	}()); err != nil {
		return start, finish, arr, errors.Wrapf(err, "In ChoosePersonIdsAndTime") //GO-TO строка
	}
	return start, finish, arr, nil
}

func SelectRouteSheetForDetail(db *sql.DB, idDetail, idEntity int64, start time.Time) ([]Waypoint, error) {
	arr := make([]Waypoint, 0)
	if err := (func() error {
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		QwStr := data.SelectRouteSheet(&idEntity, nil, nil, nil, nil)
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorQueryDB+QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			wp := Waypoint{Persons: make([]int16, 0, 2)}
			var entity int64
			err := rows.Scan(&entity, &wp.OperationNumber, &wp.Duration, &wp.OperationId, &wp.PersonCount)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorDecryptRow)
			}
			if wp.PersonCount != 0 {
				if wp.Start, wp.Finish, wp.Persons, err = ChoosePersonIdsAndTime(db,
					&idDetail, &idEntity, &wp.OperationNumber, wp.PersonCount); err != nil {
					return err // GO-TO error
				}
				wp.Check = len(wp.Persons) != 0
			} else {
				// Проверка выполнения операции при равном 0 количестве человек находится в функции detailRec.
				// if len(arr) == 0 {
				// 	wp.Check, wp.Start, wp.Finish = true, start, start.Add(SecToDur(wp.Duration))
				// } else if last := arr[len(arr)-1]; last.Check {
				// 	wp.Check, wp.Start, wp.Finish = true, last.Finish, last.Finish.Add(SecToDur(wp.Duration))
				// }
			}
			arr = append(arr, wp)
		}
		return nil
	}()); err != nil {
		return arr, errors.Wrapf(err, "In SelectRouteSheetForDetail") //GO-TO строка
	}
	return arr, nil
}

func SelectDetail(db *sql.DB, startPtr, finishPtr *time.Time) (map[int64]*DetailDB, error) {
	arr := make([]DetailDB, 0, 20)
	if err := (func() error {
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		QwStr := data.SelectDetail(nil, nil, nil, startPtr, finishPtr, nil)
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorQueryDB+QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			det := DetailDB{}
			var start, finish []uint8
			var state int8
			err := rows.Scan(&det.Id, &det.Entity, &state, &start, &finish, &det.Parent)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorDecryptRow)
			}
			if det.Start, err = time.Parse(data.C.TimeLayoutMySql, string(start)); err != nil {
				return err // GO-TO error
			}
			if det.Finish, err = time.Parse(data.C.TimeLayoutMySql, string(finish)); err != nil {
				return err // GO-TO error
			}
			if det.Way, err = SelectRouteSheetForDetail(db, det.Id, det.Entity, det.Start); err != nil {
				return err // GO-TO error
			}
			arr = append(arr, det)
		}
		return nil
	}()); err != nil {
		return make(map[int64]*DetailDB, 0), errors.Wrapf(err, "In SelectDetail") //GO-TO строка
	}
	mp := make(map[int64]*DetailDB, len(arr))
	for i, det := range arr {
		mp[det.Id] = &arr[i]
	}
	for _, det := range arr {
		if det.Parent.Valid {
			mp[det.Parent.Int64].Children = append(mp[det.Parent.Int64].Children, det.Id)
		}
	}
	for _, det := range arr {
		detailRec(&mp, det.Id)
	}
	return mp, nil
}

func detailRec(mp *map[int64]*DetailDB, idDetail int64) (bool, time.Time) {
	// Условие, чтобы не обрабатывать одну и ту же детль дважды.
	det := (*mp)[idDetail]
	if det.State == Ready {
		return true, det.NowTime()
	} else if det.State != No {
		return false, det.NowTime()
	}
	// Устанавливаем деталей текущую стадию и состояние. Старт, финиш в операции.
	way := func() bool { // функция запускается только для тех деталей, чьи дочернии ужи обработаны.
		for i, val := range det.Way {
			// Проверка, выполнена ли операция для тех где количество человек равно 0.
			if wp := &((*mp)[idDetail].Way[i]); val.PersonCount == 0 {
				if i == 0 {
					wp.Check, wp.Start, wp.Finish = true, (*mp)[idDetail].Start,
						(*mp)[idDetail].Start.Add(SecToDur(wp.Duration))
				} else if last := det.Way[i-1]; last.Check {
					wp.Check, wp.Start, wp.Finish = true, last.Finish,
						last.Finish.Add(SecToDur(wp.Duration))
				}
			}
			if !(*mp)[idDetail].Way[i].Check { // Провекрка на выход из цикла.
				(*mp)[idDetail].NowStage = i
				(*mp)[idDetail].State = InWork
				return false
			}
		}
		(*mp)[idDetail].NowStage = len(det.Way)
		(*mp)[idDetail].State = Ready
		return true
	}
	// if len(det.Children) == 0 { // GO-TO Можно упростить конструкцию условия.
	// 	return way(), det.NowTime()
	// } else {
	// Меняем старт, чек детали, вызов вэй
	var childrenReadyFlag bool = true
	var start time.Time = det.Start
	for _, val := range det.Children {
		check, childFinishTime := detailRec(mp, val)
		childrenReadyFlag = childrenReadyFlag && check
		if check && start.Before(childFinishTime) {
			start = childFinishTime
		}
	}
	(*mp)[idDetail].Start = start
	if childrenReadyFlag {
		return way(), (*mp)[idDetail].NowTime()
	} else {
		(*mp)[idDetail].State = WaitChildren
		return false, (*mp)[idDetail].Start
	}
	// }
}

func (det *DetailDB) Change(periods []PeriodChoose) {
	w := &(det.Way[det.NowStage])
	finish := periods[0].Finish
	for _, p := range periods {
		w.Persons = append(w.Persons, p.PL.Person)
		if finish.Before(p.Finish) {
			finish = p.Finish
		}
	}
	w.Finish = finish
	// GO-TO start
	w.Check = true
	det.NowStage++
	for i := det.NowStage; i < len(det.Way); i++ {
		if wi, wl := &(det.Way[i]), &(det.Way[i-1]); wi.PersonCount == 0 {
			wi.Start, wi.Finish, wi.Check = wl.Finish, wl.Finish.Add(SecToDur(wi.Duration)), true
		} else {
			break
		}
	}
	det.State = Ready
}

func (det *DetailDB) String() string {
	return fmt.Sprintf("ID = %d, state = %v, finish = %v, parent = %d, way = %v\n", det.Id, det.State, det.Finish,
		det.Parent.Int64, det.Way)
}
