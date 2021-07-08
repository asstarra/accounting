package qwery

import (
	s "accounting/data"
	"fmt"
)

// Выборка из таблицы "Entity" по заданным параметрам.
// Порядок: Id, Title, Type, Enum, Mark, Specification, Note.
func SelectEntity(vId *int64, vTitle *string, vType *int16, vEnum *bool,
	vMark *int8, vSpec, vNote *string, isChange bool) string {
	table := s.Tab["Entity"]
	sId := table.Columns["Id"].Name
	sTitle := table.Columns["Title"].Name
	sType := table.Columns["Type"].Name
	sEnum := table.Columns["Enumerable"].Name
	sMark := table.Columns["Marking"].Name
	sSpec := table.Columns["Specification"].Name
	sNote := table.Columns["Note"].Name
	strArr := make([]string, 1, 9)
	strArr[0] = fmt.Sprintf("SELECT %s AS id, %s AS title, %s AS type, %s AS enum,"+
		" %s AS mark, %s AS spec, %s AS note FROM %s",
		sId, sTitle, sType, sEnum, sMark, sSpec, sNote, table.Name)
	if vId != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sId, *vId))
	}
	if vTitle != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s LIKE '%%%s%%'", sTitle, *vTitle))
	}
	if vType != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sType, *vType))
	}
	if vEnum != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %s", sEnum, ToStr(vEnum)))
	}
	if vSpec != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s LIKE '%%%s%%'", sSpec, *vSpec))
	}
	if vMark != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sMark, *vMark))
	}
	if vNote != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s LIKE '%%%s%%'", sNote, *vNote))
	}
	if !isChange {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s != 1", sType)) //GO-TO
	}
	return Merger(strArr)
}

func InsertEntity(vTitle string, vType int16, vEnum bool, vMark int8, vSpec, vNote string) string {
	table := s.Tab["Entity"]
	sTitle := table.Columns["Title"].Name
	sType := table.Columns["Type"].Name
	sEnum := table.Columns["Enumerable"].Name
	sMark := table.Columns["Marking"].Name
	sSpec := table.Columns["Specification"].Name
	sNote := table.Columns["Note"].Name
	return fmt.Sprintf("INSERT %s (%s, %s, %s, %s, %s, %s) VALUES ('%s', %d, %s, %d, '%s', '%s')",
		table.Name, sTitle, sType, sEnum, sMark, sSpec, sNote,
		vTitle, vType, ToStr(vEnum), vMark, vSpec, vNote)
}

func UpdateEntity(vId int64, vTitle string, vType int16, vEnum bool, vMark int8, vSpec, vNote string) string {
	table := s.Tab["Entity"]
	sId := table.Columns["Id"].Name
	sTitle := table.Columns["Title"].Name
	sType := table.Columns["Type"].Name
	sEnum := table.Columns["Enumerable"].Name
	sMark := table.Columns["Marking"].Name
	sSpec := table.Columns["Specification"].Name
	sNote := table.Columns["Note"].Name
	return fmt.Sprintf("UPDATE %s SET %s = '%s', %s = %d, %s = %s, %s = %d, %s = '%s', %s = '%s' WHERE %s = %d",
		table.Name, sTitle, vTitle, sType, vType, sEnum, ToStr(vEnum), sMark, vMark,
		sSpec, vSpec, sNote, vNote, sId, vId)
}

func DeleteEntity(id int64) string {
	table := s.Tab["Entity"]
	sId := table.Columns["Id"].Name
	return fmt.Sprintf("DELETE FROM %s WHERE %s = %d", table.Name, sId, id)
}
