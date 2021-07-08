package qwery

import (
	s "accounting/data/db"
	"fmt"
)

// Выборка из таблицы "Entity" по заданным параметрам.
// Порядок: Id, Title, Type, Enum, Mark, Specification, Note.
func SelectEntity(vId *int64, vTitle *string, vType *int16, vEnum *bool,
	vMark *int8, vSpec, vNote *string, isChange bool) string {
	eTable := s.Tab["Entity"]
	eId := eTable.Columns["Id"].Name
	eTitle := eTable.Columns["Title"].Name
	eType := eTable.Columns["Type"].Name
	eEnum := eTable.Columns["Enumerable"].Name
	eMark := eTable.Columns["Marking"].Name
	eSpec := eTable.Columns["Specification"].Name
	eNote := eTable.Columns["Note"].Name

	strArr := make([]string, 1, 9)
	strArr[0] = fmt.Sprintf("SELECT %s AS id, %s AS title, %s AS type, %s AS enum,"+
		" %s AS mark, %s AS spec, %s AS note FROM %s",
		eId, eTitle, eType, eEnum, eMark, eSpec, eNote, eTable.Name)
	if vId != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", eId, *vId))
	}
	if vTitle != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s LIKE '%%%s%%'", eTitle, *vTitle))
	}
	if vType != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", eType, *vType))
	}
	if vEnum != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %s", eEnum, ToStr(vEnum)))
	}
	if vSpec != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s LIKE '%%%s%%'", eSpec, *vSpec))
	}
	if vMark != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", eMark, *vMark))
	}
	if vNote != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s LIKE '%%%s%%'", eNote, *vNote))
	}
	if !isChange {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s != 1", eType)) //GO-TO
	}
	return Merger(strArr)
}

func InsertEntity(vTitle string, vType int16, vEnum bool, vMark int8, vSpec, vNote string) string {
	eTable := s.Tab["Entity"]
	eTitle := eTable.Columns["Title"].Name
	eType := eTable.Columns["Type"].Name
	eEnum := eTable.Columns["Enumerable"].Name
	eMark := eTable.Columns["Marking"].Name
	eSpec := eTable.Columns["Specification"].Name
	eNote := eTable.Columns["Note"].Name
	return fmt.Sprintf("INSERT %s (%s, %s, %s, %s, %s, %s) VALUES ('%s', %d, %s, %d, '%s', '%s')",
		eTable.Name, eTitle, eType, eEnum, eMark, eSpec, eNote,
		vTitle, vType, ToStr(vEnum), vMark, vSpec, vNote)
}

func UpdateEntity(vId int64, vTitle string, vType int16, vEnum bool, vMark int8, vSpec, vNote string) string {
	eTable := s.Tab["Entity"]
	eId := eTable.Columns["Id"].Name
	eTitle := eTable.Columns["Title"].Name
	eType := eTable.Columns["Type"].Name
	eEnum := eTable.Columns["Enumerable"].Name
	eMark := eTable.Columns["Marking"].Name
	eSpec := eTable.Columns["Specification"].Name
	eNote := eTable.Columns["Note"].Name
	return fmt.Sprintf("UPDATE %s SET %s = '%s', %s = %d, %s = %s, %s = %d, %s = '%s', %s = '%s' WHERE %s = %d",
		eTable.Name, eTitle, vTitle, eType, vType, eEnum, ToStr(vEnum), eMark, vMark,
		eSpec, vSpec, eNote, vNote, eId, vId)
}

func DeleteEntity(id int64) string {
	eTable := s.Tab["Entity"]
	eId := eTable.Columns["Id"].Name
	return fmt.Sprintf("DELETE FROM %s WHERE %s = %d", eTable.Name, eId, id)
}
