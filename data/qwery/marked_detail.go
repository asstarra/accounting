package qwery

import (
	s "accounting/data/db"
	"fmt"
)

func SelectMarkedDetail(markings []int64) string {
	dTable := s.Tab["MarkedDetail"]
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
		" FROM %s AS a LEFT JOIN %s AS b ON a.%s <=> b.%s",
		dId, dMarking, dMark, dId, dMarking, dMark,
		dTable.Name, dTable.Name, dParent, dId)
	// s := fmt.Sprintf("SELECT %s AS id, %s AS marking, %s AS mark, %s AS perent FROM %s",
	// 	dId, dMarking, dMark, dParent, dTable.Name)
	if len(markings) != 0 {
		s += fmt.Sprintf(" WHERE a.%s IN (%s)", dMarking, marking)
	}
	// fmt.Println(s)
	return s
}

func InsertMarkedDetail(marking, parent int64, mark string) string {
	dTable := s.Tab["MarkedDetail"]
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
	dTable := s.Tab["MarkedDetail"]
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
	dTable := s.Tab["MarkedDetail"]
	dId := dTable.Columns["Id"].Name
	return fmt.Sprintf("DELETE FROM %s WHERE %s = %d",
		dTable.Name, dId, id)
}
