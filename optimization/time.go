package optimization

import (
	"accounting/data"
	"database/sql"
	"fmt"
	"sort"
	"time"

	//"log"
	"github.com/pkg/errors"
)

// Первод int32 (количество секунд) в тип time.Duration.
func SecToDur(second int32) time.Duration {
	return time.Duration(second) * time.Second
}

func AddDur(start time.Time, second int32) time.Time {
	return start.Add(time.Duration(second) * time.Second)
}

// Округляем дату и время до того же дня и время = 00:00:00.
func ClearClock(t time.Time) time.Time {
	t2, _ := time.Parse(data.C.TimeLayoutDay, t.Format(data.C.TimeLayoutDay))
	return t2
}

// Разница в днях.
func GetDifDay(start, day time.Time) int32 {
	return int32(ClearClock(day).Sub(ClearClock(start)).Hours() / 24)
}

// Хранит статистическую информацию об одном дне.
type Day struct {
	StartMin     time.Time // +
	FinishMax    time.Time // +
	StartMean    time.Time // +
	FinishMean   time.Time // +
	Duration     int32     // +
	FreeDuration int32     // +
	CountPerson  int       // +
}

func (d Day) String() string {
	return fmt.Sprintf("S=%s, F=%s, Sn=%s, Fn=%s, %d, %d, %d\n",
		d.StartMin.Format(data.C.TimeLayoutMySql), d.FinishMax.Format(data.C.TimeLayoutMySql),
		d.StartMean.Format(data.C.TimeLayoutMySql), d.FinishMean.Format(data.C.TimeLayoutMySql),
		d.Duration, d.FreeDuration, d.CountPerson)
}

func (d Day) GetDay() time.Time {
	return ClearClock(d.StartMean)
}

func (d Day) GetDuration() int32 {
	return int32(d.Finish().Sub(d.Start()).Seconds())
}

func (d Day) Start() time.Time {
	return d.StartMin
}

func (d Day) Finish() time.Time {
	return d.FinishMax
}

func SelectDays(db *sql.DB, startPtr, finishPtr *time.Time) (Timetable, error) {
	var arr []Day = make([]Day, 0, 100)
	if startPtr != nil && finishPtr != nil {
		arr = make([]Day, 0, int(finishPtr.Sub(*startPtr).Hours()/24)+1)
	}
	if err := (func() error {
		pt, err := SelectPersonTime(db, nil, startPtr, finishPtr, nil, nil, nil)
		if err != nil {
			return err //GO-TO добавить обработку ошибок
		}
		if len(pt) == 0 {
			return nil //GO-TO ??? error
		}
		sort.Slice(pt, func(i, j int) bool {
			return pt[i].Start.Before(pt[j].Start)
		})
		if startPtr == nil {
			start := (ClearClock(pt[0].Start))
			startPtr = &start
		}
		fmt.Println("1-1")
		for i := 0; i < len(pt); {
			var j, k int
			d := ClearClock(pt[i].Start)
			for j = i + 1; j < len(pt); j++ {
				if !d.Equal(ClearClock(pt[j].Start)) {
					break
				}
			}
			fmt.Println("i-j-k", i, j, k)
			day := Day{
				StartMin:  pt[i].Start,
				FinishMax: pt[i].Finish,
			}
			var start, finish int
			persons := make(map[int16]bool, 10)
			for k = i; k < j; k++ {
				day.Duration += int32(pt[k].Finish.Sub(pt[k].Start).Seconds())
				if pt[k].Detail == 0 && pt[k].Number == 0 {
					day.FreeDuration += int32(pt[k].Finish.Sub(pt[k].Start).Seconds())
				}
				if day.FinishMax.Before(pt[k].Finish) {
					day.FinishMax = pt[k].Finish
				}
				start += int(pt[k].Start.Sub(*startPtr).Seconds())
				finish += int(pt[k].Finish.Sub(*startPtr).Seconds())
				persons[pt[k].Person] = true
			}
			day.StartMean = (*startPtr).Add(SecToDur(int32(start / (j - i))))
			day.FinishMean = (*startPtr).Add(SecToDur(int32(finish / (j - i))))
			day.CountPerson = len(persons)
			arr = append(arr, day)
			i = j
		}
		return nil
	}()); err != nil {
		return Timetable{}, errors.Wrapf(err, "In SelectDays") //GO-TO строка
	}
	var dur int32 = 0
	for _, val := range arr {
		dur += val.GetDuration()
	}
	dur = dur / int32(len(arr))
	tt := Timetable{
		days:     arr,
		timeMean: dur,
	}
	return tt, nil
}

type Timetable struct {
	days     []Day
	timeMean int32
}

func (t Timetable) String() string {
	str := fmt.Sprintf("TT: timeMean=%d, days=[", t.timeMean)
	for _, val := range t.days {
		str += val.String()
	}
	return str + "]\n"
}

func (t Timetable) GetDuration(start, finish time.Time) int32 {
	var dur int32
	var si, fi int
	if i := GetDifDay(t.days[0].GetDay(), start); i < 0 {
		dur += -i * t.timeMean
		si = 0
	} else if t.days[i].Start().Before(start) {
		dur += int32(t.days[i].Finish().Sub(start).Seconds())
		si = int(i) + 1
	} else {
		si = int(i)
	}
	if i := GetDifDay(t.days[len(t.days)-1].GetDay(), finish); i > 0 {
		dur += i * t.timeMean
		fi = len(t.days)
	} else if finish.Before(t.days[len(t.days)-1-int(i)].Finish()) {
		dur += int32(finish.Sub(t.days[len(t.days)-1-int(i)].Start()))
		fi = len(t.days) - 1 - int(i)
	} else {
		fi = len(t.days) - int(i)
	}
	for ; si < fi; si++ {
		dur += t.days[si].GetDuration()
	}
	return dur
}

func check(a int) bool {
	fmt.Println("check", a)
	return a < 5
}

func A() {
	// var a, b time.Time
	a, _ := time.Parse("2006-01-02 15:04:05", "2021-06-18 08:00:00")
	b := a.Add(160379 * time.Second)
	ac := time.Duration(160379)
	fmt.Println("1-1", b.Before(b))
	fmt.Println(a, b)
	fmt.Println(a, a.Add(ac*time.Second))
	fmt.Println(b.Sub(a).Hours(), b.Sub(a).Minutes(), b.Sub(a).Seconds())
	fmt.Println(true || check(4))
	fmt.Println(true || check(6))
}
