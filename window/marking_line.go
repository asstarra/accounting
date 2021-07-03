package window

import (
	"accounting/data"
	"database/sql"
	"log"

	"github.com/pkg/errors"
)

// Выборка всех значений из таблицы описывающей состав сущности (EntityRec):
// Родительская сущность, дочерняя сущность и количество.
// Родительская сущность описывается только идентификатором, для дочерней
// также выбирается поле "тип сущности + название сущности".
func SelectEntityRec(db *sql.DB) ([]*EntityRec, error) {
	arr := make([]*EntityRec, 0)
	if err := (func() error {
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		QwStr := data.SelectEntityRec()
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorQueryDB+QwStr)
		}
		defer rows.Close()
		var e_type, title string
		for rows.Next() {
			row := EntityRec{}
			err := rows.Scan(&row.IdP, &row.Id, &e_type, &title, &row.Count)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorDecryptRow)
			}
			row.Title = e_type + " " + title
			arr = append(arr, &row)
		}
		return nil
	}()); err != nil {
		return arr, errors.Wrapf(err, data.S.InSelectEntityRec)
	}
	return arr, nil
}

// Выборка иерархических линий по составу сущностей.
func SelectMarkingLineNow(db *sql.DB) ([]*MarkingLine, map[int64]*Entity, error) {
	var mapIdToEntity map[int64]*Entity
	var MarkingLines []*MarkingLine
	if err := (func() error {
		entities, err := SelectEntities(db, "", 0, true)
		if err != nil {
			return err
		}
		mapIdToEntity = make(map[int64]*Entity, len(entities))
		MarkingLines = make([]*MarkingLine, 0, len(entities))
		for _, val := range entities {
			mapIdToEntity[val.Id] = val
		}
		entityRec, err := SelectEntityRec(db)
		if err != nil {
			return err
		}
		for _, val := range entityRec {
			mapIdToEntity[val.IdP].Children = append(mapIdToEntity[val.IdP].Children, &val.EntityRecChild)
		}
		for _, val := range mapIdToEntity {
			if val.Type == 1 {
				createLineRec([]IdCount{IdCount{Id: val.Id, Count: 1}}, val.Children, &MarkingLines, &mapIdToEntity)
				appendOrderDetail(&MarkingLines)
			}
		}
		return nil
	}()); err != nil {
		return MarkingLines, mapIdToEntity, errors.Wrapf(err, data.S.InSelectMarkingLineNew)
	}
	return MarkingLines, mapIdToEntity, nil
}

// Добавление в список иерархических линий
// линий, состоящих только из первого и последнего элемента.
func appendOrderDetail(MarkingLines *[]*MarkingLine) {
	appendOnlyOne := func(order, detail IdCount) {
		contains := false
		for _, val := range *MarkingLines {
			if len(val.Hierarchy) == 2 && val.Hierarchy[0].Id == order.Id && val.Hierarchy[len(val.Hierarchy)-1].Id == detail.Id {
				contains = true
			}
		}
		if !contains {
			*MarkingLines = append(*MarkingLines, &MarkingLine{Hierarchy: []IdCount{order, detail}})
		}
	}
	type Ids struct {
		order, detail int64
	}
	mapIdsToCount := make(map[Ids]int32)
	for _, val := range *MarkingLines {
		var count int32 = 1
		for _, val := range val.Hierarchy {
			count = count * val.Count
		}
		ids := Ids{
			order:  val.Hierarchy[0].Id,
			detail: val.Hierarchy[len(val.Hierarchy)-1].Id,
		}
		mapIdsToCount[ids] += count
	}
	for key, val := range mapIdsToCount {
		appendOnlyOne(IdCount{Id: key.order}, IdCount{Id: key.detail, Count: val})
	}
}

// Рекурсивная функция для построения иерархических линий.
func createLineRec(hierarchy []IdCount, children []*EntityRecChild, MarkingLines *[]*MarkingLine, mapIdToEntity *map[int64]*Entity) {
	for _, val := range children {
		entityChild := (*mapIdToEntity)[val.Id]
		hierarchy2 := make([]IdCount, 0, 10)
		hierarchy2 = append(hierarchy2, hierarchy...)
		hierarchy2 = append(hierarchy2, IdCount{Id: val.Id, Count: val.Count})
		if entityChild.Marking != MarkingNo {
			*MarkingLines = append(*MarkingLines, &MarkingLine{Hierarchy: hierarchy2})
		}
		createLineRec(hierarchy2, entityChild.Children, MarkingLines, mapIdToEntity)
	}
}

// Выборка сущностей входящих в иерархическую линию.
func SelectMarkingLineEntity(db *sql.DB, id int64) (map[int8]*Entity, error) {
	mapNumberToEntity := make(map[int8]*Entity)
	if err := (func() error {
		QwStr := data.SelectMarkingLineEntity(id)
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorQueryDB+QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			var number int8
			row := NewEntity()
			err := rows.Scan(&number, &row.Id, &row.Title, &row.Type, &row.Specification, &row.Marking, &row.Note)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorDecryptRow)
			}
			mapNumberToEntity[number] = &row
		}
		return nil
	}()); err != nil {
		return mapNumberToEntity, errors.Wrapf(err, data.S.InSelectMarkingLineEntity, id)
	}
	return mapNumberToEntity, nil
}

// Выборка уже определенных иерархических линий.
func SelectMarkingLineOld(db *sql.DB) ([]*MarkingLine, map[int64]*Entity, error) {
	MarkingLines := make([]*MarkingLine, 0)
	mapIdToEntity := make(map[int64]*Entity)
	if err := (func() error {
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		QwStr := data.SelectMarking()
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorQueryDB+QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			row := MarkingLine{Hierarchy: make([]IdCount, 0, 10)}
			err := rows.Scan(&row.Id)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorDecryptRow)
			}
			MarkingLines = append(MarkingLines, &row)
		}

		for index, val := range MarkingLines {
			mapNumberToEntity, err := SelectMarkingLineEntity(db, val.Id)
			if err != nil {
				return err
			}
			for number := 1; number <= len(mapNumberToEntity); number++ {
				entity := mapNumberToEntity[int8(number)]
				MarkingLines[index].Hierarchy = append(MarkingLines[index].Hierarchy, IdCount{Id: entity.Id})
				mapIdToEntity[entity.Id] = entity
			}
		}
		return nil
	}()); err != nil {
		return MarkingLines, mapIdToEntity, errors.Wrapf(err, data.S.InSelectMarkingLineOld)
	}
	return MarkingLines, mapIdToEntity, nil
}

// Обновление таблиц (Marking, MarkingLine) отвечающих за иерархию маркируемых деталей
// после изменения таблиц (Entity, EntityRec) описывающих сущности и их состав.
// Возращает отображение из id в иерархию маркируемых деталей (MarkingLine) и
// отображение из id в сущность.
func UpdateMarkingLine(db *sql.DB, isAllEntities bool) (map[int64]*MarkingLine, map[int64]*Entity, error) {
	Equal := func(a, b []IdCount) bool {
		if len(a) != len(b) {
			return false
		}
		for i, val := range a {
			if val.Id != b[i].Id {
				return false
			}
		}
		return true
	}
	var now, old []*MarkingLine
	var mapIdToEntityNow, mapIdToEntityOld map[int64]*Entity
	var mapIdToMarkingLine map[int64]*MarkingLine
	var err error
	if err = (func() error {
		now, mapIdToEntityNow, err = SelectMarkingLineNow(db)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorRead)
		}
		old, mapIdToEntityOld, err = SelectMarkingLineOld(db)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorRead)
		}
		update := func() error {
			for i, valO := range old {
				for j, valN := range now {
					if Equal(valO.Hierarchy, valN.Hierarchy) {
						now[j].Id = old[i].Id
						for k, _ := range valN.Hierarchy {
							old[i].Hierarchy[k].Count = now[j].Hierarchy[k].Count
						}
						break
					}
				}
			}
			for j, valN := range now {
				if valN.Id == 0 {
					QwStr := data.InsertMarking()
					if err := db.Ping(); err != nil {
						return errors.Wrap(err, data.S.ErrorPingDB)
					}
					result, err := db.Exec(QwStr)
					if err != nil {
						return errors.Wrap(err, data.S.ErrorAddDB+QwStr)
					}
					now[j].Id, err = result.LastInsertId()
					if err != nil {
						return errors.Wrap(err, data.S.ErrorInsertIndexLog)
					}
					for number, entityIdCount := range valN.Hierarchy {
						QwStr2 := data.InsertMarkingLine(now[j].Id, entityIdCount.Id, int8(number+1))
						if err := db.Ping(); err != nil {
							return errors.Wrap(err, data.S.ErrorPingDB)
						}
						if _, err := db.Exec(QwStr2); err != nil {
							return errors.Wrap(err, data.S.ErrorAddDB+QwStr)
						}
						mapIdToEntityOld[entityIdCount.Id] = mapIdToEntityNow[entityIdCount.Id]
					}
					old = append(old, now[j])
				}
			}
			return nil
		}
		err = update()
		return errors.Wrap(err, data.S.ErrorUpdate)
	}()); err != nil {
		err = errors.Wrap(err, data.S.ErrorUpdateMarkingLine)
		log.Println(data.S.Error, err)
		ErrorRunWindow(err.Error()) // GO-TO
	}
	mapIdToMarkingLine = make(map[int64]*MarkingLine, len(now))
	for _, val := range now {
		mapIdToMarkingLine[val.Id] = val
	}
	// fmt.Println(now)
	// fmt.Println(old)
	// fmt.Println(mapIdToEntityNow)
	// fmt.Println(mapIdToMarkingLine)
	if isAllEntities {
		return mapIdToMarkingLine, mapIdToEntityNow, err
	}
	return mapIdToMarkingLine, mapIdToEntityOld, err
}
