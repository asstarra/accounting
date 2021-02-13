package window

import (
	"fmt"
	"strings"
)

type IdTitle struct {
	Id    int64
	Title string
}

func (a IdTitle) String() string {
	return fmt.Sprintf("{Id = %d, Title = '%s'}", a.Id, a.Title)
}

type EntityRecChild struct {
	IdTitle
	Count int32
}

func (a EntityRecChild) String() string {
	return fmt.Sprintf("'%s'", a.Title)
	// return fmt.Sprintf("{Id = %d, Title = '%s', Count = %d}", a.Id, a.Title, a.Count)
}

type Marking int8

const (
	MarkingNo Marking = 1 + iota
	MarkingAll
	MarkingYear
)

var MarkingTitle = map[Marking]string{
	MarkingNo:   "Нет",
	MarkingAll:  "Сквозная",
	MarkingYear: "По годам",
}

type Entity struct {
	Id            int64
	Title         string
	Type          int64
	Specification string
	Marking       Marking
	Note          string
	Children      []*EntityRecChild
}

func NewEntity() Entity {
	return Entity{Children: make([]*EntityRecChild, 0)}
}

func (a Entity) String() string {
	return fmt.Sprintf("{Id = %d, Title = '%s', Type = %d, Specification = '%s', Marking = %s, Note = '%s', Children = %v}\n",
		a.Id, a.Title, a.Type, a.Specification, MarkingTitle[a.Marking], a.Note, a.Children)
}

type EntityRec struct {
	IdP int64
	EntityRecChild
}

type IdCount struct {
	Id    int64
	Count int32
}

type MarkingLine struct {
	Id        int64
	Hierarchy []IdCount
}

func (m MarkingLine) String() string {
	return fmt.Sprintf("{%d -> %v}\n", m.Id, m.Hierarchy)
}

// type MarkingLineGraph struct {
// 	MapIdEntity  map[int64]*Entity
// 	MarkingLines []*MarkingLine
// }

// Функция конвертирующая ошибки для показа пользователю.
func MsgError(err error) string {
	return strings.Replace(err.Error(), ": ", ":\n", -1)
}
