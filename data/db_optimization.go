package data

import (
	"database/sql"
	"fmt"
	"time"
)

func Prefix(strArr []string) string {
	if len(strArr) == 1 {
		return " WHERE "
	} else {
		return " AND "
	}
}

func Merger(strArr []string) string {
	str := ""
	for _, val := range strArr {
		str += val
	}
	return str
}

// Строки запроса к БД.

func SelectId(tableName string) string {
	table := Tab[tableName]
	sId := table.Columns["Id"].Name
	return fmt.Sprintf("SELECT %s AS id FROM %s", sId, table.Name)
}

func SelectQualification(vPerson, vOperation *int16, vLevel *int8) string {
	table := Tab["Qualification"]
	sPerson := table.Columns["Person"].Name
	sOperation := table.Columns["Operation"].Name
	sLevel := table.Columns["Level"].Name
	strArr := make([]string, 1, 4)
	strArr[0] = fmt.Sprintf("SELECT %s AS id_person, %s AS id_operation, %s AS level FROM %s",
		sPerson, sOperation, sLevel, table.Name)
	if vPerson != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sPerson, *vPerson))
	}
	if vOperation != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sOperation, *vOperation))
	}
	if vLevel != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sLevel, *vLevel))
	}
	return Merger(strArr)
}

func SelectRouteSheet(vEntity *int64, vNumber, vDuration *int32, vOperation *int16, vPersonCount *int8) string {
	table := Tab["RouteSheet"]
	sEntity := table.Columns["Entity"].Name
	sNumber := table.Columns["Number"].Name
	sDuration := table.Columns["Duration"].Name
	sOperation := table.Columns["Operation"].Name
	sPersonCount := table.Columns["PersonCount"].Name
	strArr := make([]string, 1, 7)
	strArr[0] = fmt.Sprintf("SELECT %s AS entity, %s AS number, "+
		"%s AS duration, %s AS operation, %s AS personCount FROM %s",
		sEntity, sNumber, sDuration, sOperation, sPersonCount, table.Name)
	if vEntity != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sEntity, *vEntity))
	}
	if vNumber != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sNumber, *vNumber))
	}
	if vDuration != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sDuration, *vDuration))
	}
	if vOperation != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sOperation, *vOperation))
	}
	if vPersonCount != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sPersonCount, *vPersonCount))
	}
	strArr = append(strArr, fmt.Sprintf(" ORDER BY %s, %s ASC", sEntity, sNumber))
	return Merger(strArr)
}

func SelectDetail(vId, vEntity *int64, vState *int8, vStart, vFinish *time.Time, vParent *sql.NullInt64) string { // GO-TO start/finish
	table := Tab["Detail"]
	sId := table.Columns["Id"].Name
	sEntity := table.Columns["Entity"].Name
	sState := table.Columns["State"].Name
	sStart := table.Columns["Start"].Name
	sFinish := table.Columns["Finish"].Name
	sParent := table.Columns["Parent"].Name
	strArr := make([]string, 1, 7)
	strArr[0] = fmt.Sprintf("SELECT %s AS id, %s AS entity, %s AS state, %s AS start, %s AS finish, %s AS parent FROM %s",
		sId, sEntity, sState, sStart, sFinish, sParent, table.Name)
	if vId != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sId, *vId))
	}
	if vEntity != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sEntity, *vEntity))
	}
	if vState != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sState, *vState))
	}
	if vStart != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s > %d", sFinish, *vStart))
	}
	if vFinish != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s < %d", sStart, *vFinish))
	}
	if vParent != nil {
		strParent := "NULL"
		if vParent.Valid {
			strParent = fmt.Sprint(vParent.Int64)
		}
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %s", sParent, strParent))
	}
	return Merger(strArr)
}

func SelectPersonTime(vPerson *int16, vStart, vFinish *time.Time, vDetail, vEntity *int64, vNumber *int32) string {
	table := Tab["PersonTime"]
	sPerson := table.Columns["Person"].Name
	sStart := table.Columns["Start"].Name
	sFinish := table.Columns["Finish"].Name
	sDetail := table.Columns["Detail"].Name
	sEntity := table.Columns["Entity"].Name
	sNumber := table.Columns["Number"].Name
	strArr := make([]string, 1, 8)
	strArr[0] = fmt.Sprintf("SELECT %s AS person, %s AS start, %s AS finish,"+
		" %s AS detail, %s AS entity, %s AS number FROM %s",
		sPerson, sStart, sFinish, sDetail, sEntity, sNumber, table.Name)
	if vPerson != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sPerson, *vPerson))
	}
	if vStart != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s > %d", sFinish, *vStart))
	}
	if vFinish != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s < %d", sStart, *vFinish))
	}
	if vDetail != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sDetail, *vDetail))
	}
	if vEntity != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sEntity, *vEntity))
	}
	if vNumber != nil {
		strArr = append(strArr, Prefix(strArr)+fmt.Sprintf("%s = %d", sNumber, *vNumber))
	}
	strArr = append(strArr, fmt.Sprintf(" ORDER BY %s, %s ASC", sPerson, sStart))
	return Merger(strArr)
}
