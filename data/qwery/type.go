package qwery

import (
	s "accounting/data"
	"fmt"
)

// Выборка из таблицы tableName Id и названия, в которой Id занимает 2 байта.
// Опционально, выборка названия по Id или где в названиях есть строка.
// Порядок: Id, Title.
func SelectType16(tableName string, vId *int16, vTitle *string) string {
	table := s.Tab[tableName]
	sId := table.Columns["Id"].Name
	sTitle := table.Columns["Title"].Name
	strArr := make([]string, 1, 3)
	strArr[0] = fmt.Sprintf("SELECT %s AS id, %s AS title FROM %s", sId, sTitle, table.Name)
	if vId != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sId, *vId))
	}
	if vTitle != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s LIKE '%%%s%%'", sTitle, *vTitle))
	}
	return Merger(strArr)
}

func InsertType16(tableName, vTitle string) string {
	table := s.Tab[tableName]
	sTitle := table.Columns["Title"].Name
	return fmt.Sprintf("INSERT %s (%s) VALUES ('%s')", table.Name, sTitle, vTitle)
}

func UpdateType16(tableName string, vId int16, vTitle string) string {
	table := s.Tab[tableName]
	sId := table.Columns["Id"].Name
	sTitle := table.Columns["Title"].Name
	return fmt.Sprintf("UPDATE %s SET %s = '%s' WHERE %s = %d", table.Name, sTitle, vTitle, sId, vId)
}

func DeleteType16(tableName string, vId int16) string {
	table := s.Tab[tableName]
	sId := table.Columns["Id"].Name
	return fmt.Sprintf("DELETE FROM %s WHERE %s = %d", table.Name, sId, vId)
}
