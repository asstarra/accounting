package qwery

import (
	s "accounting/data"
	"fmt"
)

func SelectMarking() string {
	table := s.Tab["Marking"]
	sId := table.Columns["Id"].Name
	return fmt.Sprintf("SELECT %s AS id FROM %s", sId, table.Name)
}

func InsertMarking() string {
	table := s.Tab["Marking"]
	return fmt.Sprintf("INSERT %s () VALUES ()", table.Name)
}

// Выборка из таблиц "MarkingLine" и "Entity" значений входящих в 1 линию.
// На вход получает Id линии. Порядок вывода: Number in Line, значения "Entity":
// Id, Title, Type, Enum, Mark, Specification, Note.
func SelectMarkingLineEntity(markId int64) string {
	mTable := s.Tab["MarkingLine"]
	mIdM := mTable.Columns["Marking"].Name
	mIdE := mTable.Columns["Entity"].Name
	mNumber := mTable.Columns["Number"].Name

	eTable := s.Tab["Entity"]
	eId := eTable.Columns["Id"].Name
	eTitle := eTable.Columns["Title"].Name
	eType := eTable.Columns["Type"].Name
	eEnum := eTable.Columns["Enumerable"].Name
	eMark := eTable.Columns["Marking"].Name
	eSpec := eTable.Columns["Specification"].Name
	eNote := eTable.Columns["Note"].Name
	return fmt.Sprintf("SELECT m.%s AS number, e.%s AS id, e.%s AS title, e.%s AS type, "+
		"%s AS enum, e.%s AS mark, e.%s AS spec, e.%s AS note "+
		"FROM %s AS e JOIN %s AS m ON e.%s = m.%s WHERE m.%s = %d",
		mNumber, eId, eTitle, eType, eEnum, eMark, eSpec, eNote,
		eTable.Name, mTable.Name, eId, mIdE, mIdM, markId)
}

func InsertMarkingLine(idM, idE int64, number int8) string {
	table := s.Tab["MarkingLine"]
	sIdM := table.Columns["Marking"].Name
	sIdE := table.Columns["Entity"].Name
	sNumber := table.Columns["Number"].Name
	return fmt.Sprintf("INSERT %s (%s, %s, %s) VALUES (%d, %d, %d)",
		table.Name, sIdM, sIdE, sNumber, idM, idE, number)
}
