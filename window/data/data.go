package data

import (
	"fmt"
)

type Id16Title struct {
	Id    int16  // Id.
	Title string // Название.
}

func (a Id16Title) String() string {
	return fmt.Sprintf("{Id = %d, Title = '%s'}", a.Id, a.Title)
}

type Id64Title struct {
	Id    int64  // Ид.
	Title string // Описание. Название типа + название сущности.
}

func (it Id64Title) String() string {
	return fmt.Sprintf("{Id = %d, Title = '%s'}", it.Id, it.Title)
}

type EntityRecChild struct { // GO-TO rename.
	Id64Title       // Ид + текстовое описание.
	Count     int32 // Количество дочерних компонентов.
}

type EntityRec struct {
	IdP            int64 // Ид родителя.
	EntityRecChild       // Ид + описание + количество.
}

func (it EntityRecChild) String() string {
	return fmt.Sprintf("'%s'", it.Title)
	// return fmt.Sprintf("{Id = %d, Title = '%s', Count = %d}", it.Id, it.Title, it.Count) //GO-TO
}

// Тип маркировки компонента
type Marking int8

const (
	MarkingNo   Marking = 1 + iota // Не маркируется.
	MarkingAll                     // Маркировка сквозная.
	MarkingYear                    // Маркировка по годам.
)

var MapMarkingToTitle = map[Marking]string{ // TO-DO
	MarkingNo:   "Нет",
	MarkingAll:  "Сквозная",
	MarkingYear: "По годам",
}

type Entity struct {
	Id            int64             // Ид.
	Title         string            // Название.
	Type          int16             // Ид типа.
	Enumerable    bool              // Можно посчитать?.
	Marking       Marking           // Способ маркировки.
	Specification string            // Спицификация.
	Note          string            // Примечание.
	Children      []*EntityRecChild // Дочерние детали: ид с описанием.
}

func NewEntity() Entity {
	return Entity{Children: make([]*EntityRecChild, 0, 0)}
}

func (e Entity) String() string {
	return fmt.Sprintf("{Id = %d, Title = '%s', Type = %d, Enum = %v, Mark = %s, Spec = '%s', Note = '%s', Children = %v}\n",
		e.Id, e.Title, e.Type, e.Enumerable, MapMarkingToTitle[e.Marking], e.Specification, e.Note, e.Children)
}

type IdCount struct {
	Id    int64 // Ид.
	Count int32 // Количество.
}

type MarkingLine struct {
	Id        int64     // Ид линии.
	Hierarchy []IdCount // Линия.
}

func (m MarkingLine) String() string {
	return fmt.Sprintf("{%d -> %v}\n", m.Id, m.Hierarchy)
}

type MarkedDetailMin struct {
	Id      int64  // Ид детали.
	Marking int64  // Ид линии.
	Mark    string // Маркировка.
}

type MarkedDetail struct {
	MarkedDetailMin                 // Дочерняя.
	Parent          MarkedDetailMin // Родительская.
}
