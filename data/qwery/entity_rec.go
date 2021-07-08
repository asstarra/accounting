package qwery

import (
	s "accounting/data/db"
	"fmt"
)

// Выборка из таблицы "EntityRec" с возможным известным Id родителя,
// где для дочерней детали помимо Id указывается название и тип.
// Порядок: IdParent, IdChild, TypeTitleChild, TitleChild, count.
func SelectEntityRecChild(idP *int64) string {
	rTable := s.Tab["EntityRec"]
	rIdP := rTable.Columns["Parent"].Name
	rIdC := rTable.Columns["Child"].Name
	rCount := rTable.Columns["Count"].Name

	eTable := s.Tab["Entity"]
	eId := eTable.Columns["Id"].Name
	eTitle := eTable.Columns["Title"].Name
	eType := eTable.Columns["Type"].Name

	tTable := s.Tab["EntityType"]
	tId := tTable.Columns["Id"].Name
	tTitle := tTable.Columns["Title"].Name
	strArr := make([]string, 1, 2)
	strArr[0] = fmt.Sprintf("SELECT r.%s AS id_p, r.%s AS id_с, j.type, j.title, r.%s AS count FROM %s AS r JOIN "+
		"(SELECT e.%s AS id, t.%s AS type, e.%s AS title FROM %s AS e JOIN %s AS t ON e.%s = t.%s) AS j "+
		"ON r.%s = j.id", rIdP, rIdC, rCount, rTable.Name,
		eId, tTitle, eTitle, eTable.Name, tTable.Name, eType, tId,
		rIdC)
	if idP != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", rIdP, *idP))
	}
	return Merger(strArr)
}

func InsertEntityRec(idP, idC int64, count int32) string {
	table := s.Tab["EntityRec"]
	rIdP := table.Columns["Parent"].Name
	rIdC := table.Columns["Child"].Name
	rCount := table.Columns["Count"].Name
	return fmt.Sprintf("INSERT %s (%s, %s, %s) VALUES (%d, %d, %d)",
		table.Name, rIdP, rIdC, rCount, idP, idC, count)
}

func UpdateEntityRec(idP, idC int64, count int32) string {
	table := s.Tab["EntityRec"]
	rIdP := table.Columns["Parent"].Name
	rIdC := table.Columns["Child"].Name
	rCount := table.Columns["Count"].Name
	return fmt.Sprintf("UPDATE %s SET %s = %d WHERE %s = %d AND %s = %d",
		table.Name, rCount, count, rIdP, idP, rIdC, idC)
}

func DeleteEntityRec(idP, idC int64) string {
	table := s.Tab["EntityRec"]
	rIdP := table.Columns["Parent"].Name
	rIdC := table.Columns["Child"].Name
	return fmt.Sprintf("DELETE FROM %s WHERE %s = %d AND %s = %d",
		table.Name, rIdP, idP, rIdC, idC)
}
