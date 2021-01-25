package window

import (
	"fmt"
)

type IdTitle struct {
	Id    int64
	Title string
}

func (a IdTitle) String() string {
	return fmt.Sprintf("{Id = %d, Title = %s}", a.Id, a.Title)
}

type EntityRecChild struct {
	IdTitle
	Count int
}

func (a EntityRecChild) String() string {
	return fmt.Sprintf("{Id = %d, Title = %s, Count = %d}", a.Id, a.Title, a.Count)
}

type Entity struct {
	Id             int64
	Title          string
	Type           int64
	Specification  string
	ProductionLine bool
	Note           string
	Children       *[]*EntityRecChild
}

func (a Entity) String() string {
	var c string
	if a.Children == nil {
		c = "nil"
	} else {
		c = fmt.Sprint(*a.Children)
	}
	return fmt.Sprintf("{Id = %d, Title = %s, Type = %d, Specification = %s, ProductionLine = %t, Note = %s, Children = %s}",
		a.Id, a.Title, a.Type, a.Specification, a.ProductionLine, a.Note, c)
}
