package data

import (
	"fmt"
)

func SelectEntity(title string, entityType int64) string {
	table := Tab["Entity"]
	sId := table.Columns["Id"].Name
	sTitle := table.Columns["Title"].Name
	sType := table.Columns["Type"].Name
	sSpec := table.Columns["Specification"].Name
	sProdLine := table.Columns["ProductionLine"].Name
	sNote := table.Columns["Note"].Name
	s := fmt.Sprintf("SELECT %s AS id, %s AS title, %s AS type, %s AS spec, %s AS pline, %s AS note FROM %s",
		sId, sTitle, sType, sSpec, sProdLine, sNote, table.Name)
	// isWhere := false
	// var sWhere string
	// if entityType != 0
	//GO-TO выборка строчек по условию
	return s
}

func InsertEntity(vTitle, vSpec, vNote string, vProdL bool, vType int64) string {
	table := Tab["Entity"]
	sTitle := table.Columns["Title"].Name
	sType := table.Columns["Type"].Name
	sSpec := table.Columns["Specification"].Name
	sProdLine := table.Columns["ProductionLine"].Name
	sNote := table.Columns["Note"].Name
	return fmt.Sprintf("INSERT %s (%s, %s, %s, %s, %s) VALUES ('%s', %d, '%s', %t, '%s')",
		table.Name, sTitle, sType, sSpec, sProdLine, sNote,
		vTitle, vType, vSpec, vProdL, vNote)
}

func UpdateEntity(vTitle, vSpec, vNote string, vProdL bool, vType, vId int64) string {
	table := Tab["Entity"]
	sId := table.Columns["Id"].Name
	sTitle := table.Columns["Title"].Name
	sType := table.Columns["Type"].Name
	sSpec := table.Columns["Specification"].Name
	sProdLine := table.Columns["ProductionLine"].Name
	sNote := table.Columns["Note"].Name
	return fmt.Sprintf("UPDATE %s SET %s = '%s', %s = %d, %s = '%s', %s = %t, %s = '%s' WHERE %s = %d",
		table.Name, sTitle, vTitle, sType, vType, sSpec, vSpec,
		sProdLine, vProdL, sNote, vNote, sId, vId)
}

func DeleteEntity(id int64) string {
	table := Tab["Entity"]
	sId := table.Columns["Id"].Name
	return fmt.Sprintf("DELETE FROM %s WHERE %s = %d", table.Name, sId, id)
}

func SelectEntityRecChild(idP int64) string {
	rTable := Tab["EntityRec"]
	sIdC := rTable.Columns["Child"].Name
	sIdP := rTable.Columns["Parent"].Name
	sCount := rTable.Columns["Count"].Name
	eTable := Tab["Entity"]
	eId := eTable.Columns["Id"].Name
	eTitle := eTable.Columns["Title"].Name
	eType := eTable.Columns["Type"].Name
	tTable := Tab["EntityType"]
	tId := tTable.Columns["Id"].Name
	tTitle := tTable.Columns["Title"].Name
	return fmt.Sprintf("SELECT r.%s AS id, j.type, j.title, r.%s AS count FROM %s AS r JOIN "+
		"(SELECT e.%s AS id, t.%s AS type, e.%s AS title FROM %s AS e JOIN %s AS t ON e.%s = t.%s) AS j "+
		"ON r.%s = j.id WHERE r.%s = %d", sIdC, sCount, rTable.Name,
		eId, tTitle, eTitle, eTable.Name, tTable.Name, eType, tId,
		sIdC, sIdP, idP)
	// return fmt.Sprintf("SELECT %s AS idc, %s AS count FROM %s WHERE %s = %d",
	// 	sIdC, sCount, rTable.Name, sIdP, idP)
}

func InsertEntityRec(idP, idC int64, count int) string {
	table := Tab["EntityRec"]
	sIdP := table.Columns["Parent"].Name
	sIdC := table.Columns["Child"].Name
	sCount := table.Columns["Count"].Name
	return fmt.Sprintf("INSERT %s (%s, %s, %s) VALUES (%d, %d, %d)",
		table.Name, sIdP, sIdC, sCount, idP, idC, count)
}

func UpdateEntityRec(idP, idC int64, count int) string {
	table := Tab["EntityRec"]
	sIdP := table.Columns["Parent"].Name
	sIdC := table.Columns["Child"].Name
	sCount := table.Columns["Count"].Name
	return fmt.Sprintf("UPDATE %s SET %s = %d WHERE %s = %d AND %s = %d",
		table.Name, sCount, count, sIdP, idP, sIdC, idC)
}

func DeleteEntityRec(idP, idC int64) string {
	table := Tab["EntityRec"]
	sIdP := table.Columns["Parent"].Name
	sIdC := table.Columns["Child"].Name
	return fmt.Sprintf("DELETE FROM %s  WHERE %s = %d AND %s = %d",
		table.Name, sIdP, idP, sIdC, idC)
}
