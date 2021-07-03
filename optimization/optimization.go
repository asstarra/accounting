package optimization

import (
	// "accounting/data"
	"container/heap"
	"database/sql"
	"fmt"
	"sort"
	"time"

	//"log"
	"github.com/pkg/errors"
)

// Item - это то, чем мы управляем в приоритетной очереди.
type Item struct {
	value    *DetailDB // Значение элемента; произвольное.
	priority int32     // Приоритет элемента в очереди.
	// Индекс необходим для обновления
	// и поддерживается методами heap.Interface.
	index int // Индекс элемента в куче.
}

func (item *Item) String() string {
	return fmt.Sprintf("ID = %d, priority  = %d, state = %v, Way = %v\n",
		item.value.Id, item.priority, item.value.State, item.value.Way)
}

// func (i *Item) GetPriority() int32 {
// 	return i.value.Priority
// }

// PriorityQueue реализует heap.Interface и содержит Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	// Мы хотим, чтобы Pop давал нам самый высокий,
	// а не самый низкий приоритет,
	// поэтому здесь мы используем оператор больше.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq) // pq.Len()
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // избежать утечки памяти
	item.index = -1 // для безопасности
	*pq = old[0 : n-1]
	return item
}

// update изменяет приоритет и значение Item в очереди.
func (pq *PriorityQueue) update(item *Item, value *DetailDB, priority int32) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

type Optimization struct {
	start, finish time.Time
	opt           OperationPersonTable
	persons       map[int16]*PersonDB
	detail        map[int64]*DetailDB
	tt            Timetable
	pqDet         PriorityQueue
}

func (o *Optimization) getDurationParentRec(idDetail int64) (int32, time.Time) {
	det := o.detail[idDetail]
	var sum int32 = 0
	var finish time.Time = o.detail[idDetail].Finish
	if det.Parent.Valid {
		sum, finish = o.getDurationParentRec(o.detail[det.Parent.Int64].Id)
	}
	return o.detail[idDetail].WorkDuration() + sum, finish
}

func (o *Optimization) SetPriority(idDetail int64) bool {
	var childrenReadyFlag bool = true
	var det *DetailDB = o.detail[idDetail]
	for _, val := range det.Children {
		childrenReadyFlag = childrenReadyFlag && o.detail[val].State == Ready
	}
	if !childrenReadyFlag {
		return false
	} else {
		det.State = InWork
	}
	if det.WorkDuration() == 0 {
		det.State = Ready
		return false // GO-TO ?
	}
	dur, finish := o.getDurationParentRec(idDetail)
	free := o.tt.GetDuration(det.NowTime(), finish)
	det.Priority = dur - free
	return true
}

func (o *Optimization) Init(db *sql.DB, startPtr, finishPtr *time.Time) error {
	if err := (func() error {
		var err error
		if o.opt, err = NewQualificationTable(db); err != nil {
			return err
		}
		if o.persons, err = SelectPerson(db, startPtr, finishPtr); err != nil {
			return err
		}
		if o.detail, err = SelectDetail(db, startPtr, finishPtr); err != nil {
			return err
		}
		if o.tt, err = SelectDays(db, startPtr, finishPtr); err != nil {
			return err
		}
		o.pqDet = make(PriorityQueue, 0, len(o.detail))
		var i = 0
		for ind, det := range o.detail {
			if o.SetPriority(ind) {
				o.pqDet = append(o.pqDet, &Item{
					value:    det,
					priority: det.Priority,
					index:    i,
				})
			}
		}
		return nil
	}()); err != nil {
		return errors.Wrapf(err, "In OptimizationInit") //GO-TO строка
	}
	// for o.pqDet.Len() > 0 {
	// 	item := heap.Pop(&o.pqDet).(*Item)
	// 	fmt.Printf("%.2d:%v - %d\n", item.priority, item.value.Id, item.value.WorkDuration())
	// }
	o.Print()
	return nil
}

type PeriodsPriority struct {
	list     map[int16]PeriodChoose
	priority int32
}

func (o *Optimization) selectPeriod(det *DetailDB) []PeriodChoose { // GO-TO упрощенная версия, доделать.
	arr := make([]PeriodChoose, 0, 20)
	persons := o.opt.GetPersons(det.Way[det.NowStage].OperationId)
	fmt.Println("persons", persons)
	for _, person := range persons {
		pc := o.persons[person].GetPeriod(det)
		fmt.Println("pc", pc)
		if len(pc) != 0 {
			arr = append(arr, pc[0])
		}
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Finish.Before(arr[j].Finish)
	})
	fmt.Println("periods arr", arr, arr[:det.Way[det.NowStage].PersonCount])
	if len(arr) < int(det.Way[det.NowStage].PersonCount) {
		return arr // GO-TO деталь выбыла, из-за отсутствия периодов.
	}
	return arr[:det.Way[det.NowStage].PersonCount]
}

func (o *Optimization) update() {
	det := heap.Pop(&(o.pqDet)).(*Item).value
	periods := o.selectPeriod(det)
	fmt.Println("periods", periods)
	if len(periods) < int(det.Way[det.NowStage].PersonCount) {
		return // GO-TO деталь выбыла, из-за отсутствия периодов.
	}
	for _, period := range periods {
		fmt.Println(period)
		fmt.Println(period.PL.Person, o.persons)
		fmt.Println(o.persons[period.PL.Person])
		o.persons[period.PL.Person].Change(period.Index, det.NowTime(), period.Finish, det)
	}
	det.Change(periods)
	if det.State == Ready {
		if det.Parent.Valid {
			parent := o.detail[det.Parent.Int64]
			var childrenReadyFlag bool = true
			for _, val := range parent.Children {
				childrenReadyFlag = childrenReadyFlag && o.detail[val].State == Ready
			}

			if childrenReadyFlag {
				o.SetPriority(parent.Id)
				parent.State = InWork
				heap.Push(&(o.pqDet), &Item{
					value:    parent,
					priority: parent.Priority,
				})
				// GO-TO удалить стадии где 0 человек.
			}
		}
		return
	}
	o.SetPriority(det.Id)
	heap.Push(&(o.pqDet), &Item{
		value:    det,
		priority: det.Priority,
	})
}

func (o *Optimization) Print() {
	fmt.Println(o.detail)
	fmt.Println(o.persons)
	fmt.Println(o.pqDet)
}

func (o *Optimization) Run() {
	var i = 0
	for o.pqDet.Len() > 0 {
		o.update()
		i++
		fmt.Println("------", i)
	}
	o.Print()
}
