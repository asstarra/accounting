package window

import (
	"accounting/data"
	e "accounting/data/errors"
	l "accounting/data/log"
	"accounting/data/qwery"
	. "accounting/window/data"
	"database/sql"

	// "fmt"
	"log"

	"github.com/pkg/errors"
)

// Выборка иерархических линий по составу сущностей.
// Выбираются все сущности и записываются в карту, затем выбирается состав,
// из состава дочерние сущности добавляются к родительским.
// Возращает массив линий и карту сущностей.
func SelectMarkingLineNow(db *sql.DB) ([]*MarkingLine, map[int64]*Entity, error) {
	var mapIdToEntity map[int64]*Entity
	var MarkingLines []*MarkingLine
	if err := (func() error {
		// Выбираем все сущности.
		entities, err := SelectEntity(db, nil, nil, nil, nil, nil, nil, nil, true)
		if err != nil {
			return err
		}
		mapIdToEntity = make(map[int64]*Entity, len(entities))
		for _, val := range entities {
			mapIdToEntity[val.Id] = val
		}
		// Выбираем состав.
		entityRec, _, err := SelectEntityRecChild(db, nil)
		if err != nil {
			return err
		}
		for _, val := range entityRec { // Вносим информацию о составе в карту сущностей.
			mapIdToEntity[val.IdP].Children = append(mapIdToEntity[val.IdP].Children, &val.EntityRecChild)
		}
		MarkingLines = make([]*MarkingLine, 0, len(entities))
		for _, val := range mapIdToEntity {
			if val.Type == 1 { // По данным строим линии.
				createLineRec([]IdCount{IdCount{Id: val.Id, Count: 1}}, val.Children, &MarkingLines, &mapIdToEntity)
				appendOrderDetail(&MarkingLines)
			}
		}
		return nil
	}()); err != nil { // Обработка ошибок.
		return MarkingLines, mapIdToEntity, errors.Wrapf(err, l.In.InSelectMarkingLineNew)
	}
	return MarkingLines, mapIdToEntity, nil
}

// Добавление в список иерархических линий, таких линий, которые состоят только
// из первого и последнего элемента, с указанием количества.
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

// Выборка сущностей входящих в одну иерархическую линию.
// Возращает отображение номера в линии в сущность.
func SelectMarkingLineEntity(db *sql.DB, id int64) (map[int8]*Entity, error) {
	mapNumberToEntity := make(map[int8]*Entity)
	if err := (func() error {
		QwStr := qwery.SelectMarkingLineEntity(id)
		if err := db.Ping(); err != nil { // Пинг БД.
			return errors.Wrap(err, e.Err.ErrorPingDB)
		}
		rows, err := db.Query(QwStr) // Запрос к БД.
		if err != nil {
			return errors.Wrapf(err, e.Err.ErrorQueryDB, QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			var number int8
			row := NewEntity()
			err := rows.Scan(&number, &row.Id, &row.Title, &row.Type, &row.Enumerable, &row.Marking, &row.Specification, &row.Note)
			if err != nil {
				return errors.Wrap(err, e.Err.ErrorDecryptRow)
			}
			mapNumberToEntity[number] = &row
		}
		return nil
	}()); err != nil { // Обработка ошибок.
		return mapNumberToEntity, errors.Wrapf(err, l.In.InSelectMarkingLineEntity, id)
	}
	return mapNumberToEntity, nil
}

// Выборка уже определенных иерархических линий.
func SelectMarkingLineOld(db *sql.DB) ([]*MarkingLine, map[int64]*Entity, error) {
	MarkingLines := make([]*MarkingLine, 0)
	mapIdToEntity := make(map[int64]*Entity)
	if err := (func() error {
		QwStr := qwery.SelectMarking()    // Выбор всех Ид линий.
		if err := db.Ping(); err != nil { // Пинг БД.
			return errors.Wrap(err, e.Err.ErrorPingDB)
		}
		rows, err := db.Query(QwStr) // Запрос к БД.
		if err != nil {
			return errors.Wrapf(err, e.Err.ErrorQueryDB, QwStr)
		}
		defer rows.Close()
		for rows.Next() { // Создание линий с пустой иерархией.
			row := MarkingLine{Hierarchy: make([]IdCount, 0, 10)}
			err := rows.Scan(&row.Id)
			if err != nil {
				return errors.Wrap(err, e.Err.ErrorDecryptRow)
			}
			MarkingLines = append(MarkingLines, &row)
		}

		for index, val := range MarkingLines {
			// Выбираем Entity входящие в данную линию.
			// Отображение номера линии в Entity.
			mapNumberToEntity, err := SelectMarkingLineEntity(db, val.Id)
			if err != nil {
				return err
			}
			for number := 1; number <= len(mapNumberToEntity); number++ {
				entity := mapNumberToEntity[int8(number)]
				// Добавление в линию.
				MarkingLines[index].Hierarchy = append(MarkingLines[index].Hierarchy, IdCount{Id: entity.Id})
				mapIdToEntity[entity.Id] = entity // Внесение сущности в общее отображение из Ид в сущность.
			}
		}
		return nil
	}()); err != nil { // Обработка ошибок.
		return MarkingLines, mapIdToEntity, errors.Wrapf(err, l.In.InSelectMarkingLineOld)
	}
	return MarkingLines, mapIdToEntity, nil
}

// Обновление таблиц (Marking, MarkingLine) отвечающих за иерархию маркируемых деталей
// после изменения таблиц (Entity, EntityRec) описывающих сущности и их состав.
// Возращает отображение из id в иерархию маркируемых деталей (MarkingLine) и
// отображение из id в сущность.
// Считает, что нет ошибок в БД.
func UpdateMarkingLine(db *sql.DB, isAllEntities bool) (map[int64]*MarkingLine, map[int64]*Entity, error) {
	// Функция, которая проверяет эквивалентность элементов линии.
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
			return errors.Wrap(err, e.Err.ErrorRead)
		}
		old, mapIdToEntityOld, err = SelectMarkingLineOld(db)
		if err != nil {
			return errors.Wrap(err, e.Err.ErrorRead)
		}
		update := func() error {
			for i, valO := range old {
				for j, valN := range now { // У совпадающих добавляем недостоющие данные.
					if Equal(valO.Hierarchy, valN.Hierarchy) {
						now[j].Id = old[i].Id
						for k, _ := range valN.Hierarchy {
							old[i].Hierarchy[k].Count = now[j].Hierarchy[k].Count
						}
						break
					}
				} // Удаление старых ненужных линий происходит в при изменении состава.
			}
			for j, valN := range now { // Вставка новых линий в иерархию.
				if valN.Id == 0 {
					QwStr := qwery.InsertMarking()
					if err := db.Ping(); err != nil {
						return errors.Wrap(err, e.Err.ErrorPingDB)
					}
					result, err := db.Exec(QwStr)
					if err != nil {
						return errors.Wrapf(err, e.Err.ErrorAddDB, QwStr)
					}
					now[j].Id, err = result.LastInsertId()
					if err != nil {
						return errors.Wrap(err, e.Err.ErrorInsertIndexLog)
					}
					for number, entityIdCount := range valN.Hierarchy {
						QwStr2 := data.InsertMarkingLine(now[j].Id, entityIdCount.Id, int8(number+1))
						if err := db.Ping(); err != nil {
							return errors.Wrap(err, e.Err.ErrorPingDB)
						}
						if _, err := db.Exec(QwStr2); err != nil {
							return errors.Wrapf(err, e.Err.ErrorAddDB, QwStr)
						}
						mapIdToEntityOld[entityIdCount.Id] = mapIdToEntityNow[entityIdCount.Id]
					}
					old = append(old, now[j])
				}
			}
			return nil
		}
		err = update()
		return errors.Wrap(err, e.Err.ErrorUpdate)
	}()); err != nil { // Обработка ошибок.
		err = errors.Wrap(err, e.Err.ErrorUpdateMarkingLine)
		log.Println(l.Error, err)     // Лог.
		ErrorRunWindow(MsgError(err)) // GO-TO
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
