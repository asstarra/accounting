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
	Count int
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

type MarkingLine struct {
	Id       int64
	Entities []*Entity
}

func (m MarkingLine) String() string {
	var c string
	for _, val := range m.Entities {
		c += fmt.Sprintf("%s, ", val.Title)
	}
	return fmt.Sprintf("\n{Id: %d, Entities: [%v]}", m.Id, c)
}

func MsgError(err error) string {
	return strings.Replace(err.Error(), ": ", ":\n", -1)
}
