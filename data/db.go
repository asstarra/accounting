package data

import (
	"fmt"
)

// Строки запроса к БД.

func SelectType(tableName string) string {
	table := Tab[tableName]
	sId := table.Columns["Id"].Name
	sTitle := table.Columns["Title"].Name
	return fmt.Sprintf("SELECT %s AS id, %s AS title FROM %s", sId, sTitle, table.Name)
}

func InsertType(tableName, vTitle string) string {
	table := Tab[tableName]
	sTitle := table.Columns["Title"].Name
	return fmt.Sprintf("INSERT %s (%s) VALUES ('%s')", table.Name, sTitle, vTitle)
}

func UpdateType(tableName, vTitle string, vId int64) string {
	table := Tab[tableName]
	sId := table.Columns["Id"].Name
	sTitle := table.Columns["Title"].Name
	return fmt.Sprintf("UPDATE %s SET %s = '%s' WHERE %s = %d", table.Name, sTitle, vTitle, sId, vId)
}

func DeleteType(tableName string, vId int64) string {
	table := Tab[tableName]
	sId := table.Columns["Id"].Name
	return fmt.Sprintf("DELETE FROM %s WHERE %s = %d", table.Name, sId, vId)
}

func SelectEntity(title string, entityType int64, isChange bool) string {
	table := Tab["Entity"]
	sId := table.Columns["Id"].Name
	sTitle := table.Columns["Title"].Name
	sType := table.Columns["Type"].Name
	sSpec := table.Columns["Specification"].Name
	sMarking := table.Columns["Marking"].Name
	sNote := table.Columns["Note"].Name
	s := fmt.Sprintf("SELECT %s AS id, %s AS title, %s AS type, %s AS spec, %s AS mark, %s AS note FROM %s WHERE %s LIKE '%%%s%%'",
		sId, sTitle, sType, sSpec, sMarking, sNote, table.Name, sTitle, title)
	if entityType != 0 {
		s += fmt.Sprintf(" AND %s = %d", sType, entityType)
	}
	if !isChange {
		s += fmt.Sprintf(" AND %s != 1", sType)
	}
	return s
}

func InsertEntity(vTitle, vSpec, vNote string, vMarking int8, vType int64) string {
	table := Tab["Entity"]
	sTitle := table.Columns["Title"].Name
	sType := table.Columns["Type"].Name
	sSpec := table.Columns["Specification"].Name
	sMarking := table.Columns["Marking"].Name
	sNote := table.Columns["Note"].Name
	return fmt.Sprintf("INSERT %s (%s, %s, %s, %s, %s) VALUES ('%s', %d, '%s', %d, '%s')",
		table.Name, sTitle, sType, sSpec, sMarking, sNote,
		vTitle, vType, vSpec, vMarking, vNote)
}

func UpdateEntity(vTitle, vSpec, vNote string, vMarking int8, vType, vId int64) string {
	table := Tab["Entity"]
	sId := table.Columns["Id"].Name
	sTitle := table.Columns["Title"].Name
	sType := table.Columns["Type"].Name
	sSpec := table.Columns["Specification"].Name
	sMarking := table.Columns["Marking"].Name
	sNote := table.Columns["Note"].Name
	return fmt.Sprintf("UPDATE %s SET %s = '%s', %s = %d, %s = '%s', %s = %d, %s = '%s' WHERE %s = %d",
		table.Name, sTitle, vTitle, sType, vType, sSpec, vSpec,
		sMarking, vMarking, sNote, vNote, sId, vId)
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
}

func SelectEntityRec() string {
	// rTable := Tab["EntityRec"]
	// sIdC := rTable.Columns["Child"].Name
	// sIdP := rTable.Columns["Parent"].Name
	// return fmt.Sprintf("SELECT %s AS id_p, %s AS id_c FROM %s",
	// 	sIdP, sIdC, rTable.Name)
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
	return fmt.Sprintf("SELECT r.%s AS id_p, r.%s AS id_c, j.type, j.title, r.%s AS count FROM %s AS r JOIN "+
		"(SELECT e.%s AS id, t.%s AS type, e.%s AS title FROM %s AS e JOIN %s AS t ON e.%s = t.%s) AS j "+
		"ON r.%s = j.id", sIdP, sIdC, sCount, rTable.Name,
		eId, tTitle, eTitle, eTable.Name, tTable.Name, eType, tId,
		sIdC)
}

func InsertEntityRec(idP, idC int64, count int32) string {
	table := Tab["EntityRec"]
	sIdP := table.Columns["Parent"].Name
	sIdC := table.Columns["Child"].Name
	sCount := table.Columns["Count"].Name
	return fmt.Sprintf("INSERT %s (%s, %s, %s) VALUES (%d, %d, %d)",
		table.Name, sIdP, sIdC, sCount, idP, idC, count)
}

func UpdateEntityRec(idP, idC int64, count int32) string {
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
	return fmt.Sprintf("DELETE FROM %s WHERE %s = %d AND %s = %d",
		table.Name, sIdP, idP, sIdC, idC)
}

func SelectMarking() string {
	table := Tab["Marking"]
	sId := table.Columns["Id"].Name
	return fmt.Sprintf("SELECT %s AS id FROM %s", sId, table.Name)
}

func InsertMarking() string {
	table := Tab["Marking"]
	sId := table.Columns["Id"].Name
	return fmt.Sprintf("INSERT %s (%s) VALUES (0)", table.Name, sId)
}

func SelectMarkingLineEntity(markId int64) string {
	mTable := Tab["MarkingLine"]
	mIdM := mTable.Columns["Marking"].Name
	mIdE := mTable.Columns["Entity"].Name
	mNumber := mTable.Columns["Number"].Name
	table := Tab["Entity"]
	sId := table.Columns["Id"].Name
	sTitle := table.Columns["Title"].Name
	sType := table.Columns["Type"].Name
	sSpec := table.Columns["Specification"].Name
	sMarking := table.Columns["Marking"].Name
	sNote := table.Columns["Note"].Name
	return fmt.Sprintf("SELECT m.%s AS number, e.%s AS id, e.%s AS title, e.%s AS type, e.%s AS spec, e.%s AS mark, e.%s AS note "+
		"FROM %s AS e JOIN %s AS m ON e.%s = m.%s WHERE m.%s = %d",
		mNumber, sId, sTitle, sType, sSpec, sMarking, sNote,
		table.Name, mTable.Name, sId, mIdE, mIdM, markId)
}

func InsertMarkingLine(idM, idE int64, number int8) string {
	table := Tab["MarkingLine"]
	sIdM := table.Columns["Marking"].Name
	sIdE := table.Columns["Entity"].Name
	sNumber := table.Columns["Number"].Name
	return fmt.Sprintf("INSERT %s (%s, %s, %s) VALUES (%d, %d, %d)",
		table.Name, sIdM, sIdE, sNumber, idM, idE, number)
}

func SelectMarkedDetail(markings []int64) string {
	dTable := Tab["MarkedDetail"]
	dId := dTable.Columns["Id"].Name
	dMarking := dTable.Columns["Marking"].Name
	dMark := dTable.Columns["Mark"].Name
	dParent := dTable.Columns["Parent"].Name
	var marking string
	for _, val := range markings {
		marking += fmt.Sprintf(", %d", val)
	}
	if len(marking) > 2 {
		marking = marking[2:]
	}
	s := fmt.Sprintf("SELECT a.%s AS id, a.%s AS marking, a.%s AS mark, b.%s AS id, b.%s AS marking, b.%s AS mark"+
		" FROM %s AS a LEFT JOIN %s AS b ON a.%s <=> b.%s;",
		dId, dMarking, dMark, dId, dMarking, dMark,
		dTable.Name, dTable.Name, dParent, dId)
	// s := fmt.Sprintf("SELECT %s AS id, %s AS marking, %s AS mark, %s AS perent FROM %s",
	// 	dId, dMarking, dMark, dParent, dTable.Name)
	if len(markings) != 0 {
		s += fmt.Sprintf(" WHERE %s IN (%s)", dMarking, marking)
	}
	return s
}

func InsertMarkedDetail(marking, parent int64, mark string) string {
	dTable := Tab["MarkedDetail"]
	dMarking := dTable.Columns["Marking"].Name
	dMark := dTable.Columns["Mark"].Name
	dParent := dTable.Columns["Parent"].Name
	vParent := fmt.Sprintf("%d", parent)
	if parent == 0 {
		vParent = "NULL"
	}
	return fmt.Sprintf("INSERT %s (%s, %s, %s) VALUES (%d, '%s', %s)",
		dTable.Name, dMarking, dMark, dParent, marking, mark, vParent)
}

func UpdateMarkedDetail(id, marking, parent int64, mark string) string {
	dTable := Tab["MarkedDetail"]
	dId := dTable.Columns["Id"].Name
	dMarking := dTable.Columns["Marking"].Name
	dMark := dTable.Columns["Mark"].Name
	dParent := dTable.Columns["Parent"].Name
	vParent := fmt.Sprintf("%d", parent)
	if parent == 0 {
		vParent = "NULL"
	}
	return fmt.Sprintf("UPDATE %s SET %s = %d, %s = '%s', %s = %s WHERE %s = %d",
		dTable.Name, dMarking, marking, dMark, mark, dParent, vParent, dId, id)
}

func DeleteMarkedDetail(id int64) string {
	dTable := Tab["MarkedDetail"]
	dId := dTable.Columns["Id"].Name
	return fmt.Sprintf("DELETE FROM %s WHERE %s = %d",
		dTable.Name, dId, id)
}
