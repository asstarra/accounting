package window

import (
	"accounting/data"
	"database/sql"
	"fmt"
	"log"
	"reflect"

	"github.com/pkg/errors"
)

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

func SelectMarkingLineNew(db *sql.DB) ([]*MarkingLine, map[int64]*Entity, error) {
	var MapIdEntity map[int64]*Entity
	var MarkingLines []*MarkingLine
	if err := (func() error {
		entities, err := SelectEntities(db, "", 0, true)
		if err != nil {
			return err
		}
		MapIdEntity = make(map[int64]*Entity, len(entities))
		MarkingLines = make([]*MarkingLine, 0, len(entities))
		for _, val := range entities {
			MapIdEntity[val.Id] = val
		}
		entityRec, err := SelectEntityRec(db)
		if err != nil {
			return err
		}
		for _, val := range entityRec {
			MapIdEntity[val.IdP].Children = append(MapIdEntity[val.IdP].Children, &val.EntityRecChild)
		}
		for _, val := range MapIdEntity {
			if val.Type == 1 {
				createLineRec([]int64{val.Id}, val.Children, &MarkingLines, &MapIdEntity)
			}
		}
		return nil
	}()); err != nil {
		return MarkingLines, MapIdEntity, errors.Wrapf(err, data.S.InSelectMarkingLineNew)
	}
	return MarkingLines, MapIdEntity, nil
}

func createLineRec(hierarchy []int64, children []*EntityRecChild, MarkingLines *[]*MarkingLine, MapIdEntity *map[int64]*Entity) {
	for _, val := range children {
		entityChild := (*MapIdEntity)[val.Id]
		hierarchy2 := make([]int64, 0, 10)
		hierarchy2 = append(hierarchy2, hierarchy...)
		hierarchy2 = append(hierarchy2, entityChild.Id)
		if entityChild.Marking != MarkingNo {
			*MarkingLines = append(*MarkingLines, &MarkingLine{Hierarchy: hierarchy2})
		}
		createLineRec(hierarchy2, entityChild.Children, MarkingLines, MapIdEntity)
	}
}

func SelectMarkingLineEntity(db *sql.DB, id int64) (map[int8]*Entity, error) {
	mapNumberEntity := make(map[int8]*Entity)
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
			mapNumberEntity[number] = &row
		}
		return nil
	}()); err != nil {
		return mapNumberEntity, errors.Wrapf(err, data.S.InSelectMarkingLineEntity, id)
	}
	return mapNumberEntity, nil
}

func SelectMarkingLineOld(db *sql.DB) ([]*MarkingLine, map[int64]*Entity, error) {
	MarkingLines := make([]*MarkingLine, 0)
	MapIdEntity := make(map[int64]*Entity)
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
			row := MarkingLine{Hierarchy: make([]int64, 0, 10)}
			err := rows.Scan(&row.Id)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorDecryptRow)
			}
			MarkingLines = append(MarkingLines, &row)
		}

		for index, val := range MarkingLines {
			mapNumberEntity, err := SelectMarkingLineEntity(db, val.Id)
			if err != nil {
				return err
			}
			for number := 1; number <= len(mapNumberEntity); number++ {
				entity := mapNumberEntity[int8(number)]
				MarkingLines[index].Hierarchy = append(MarkingLines[index].Hierarchy, entity.Id)
				MapIdEntity[entity.Id] = entity
			}
		}
		return nil
	}()); err != nil {
		return MarkingLines, MapIdEntity, errors.Wrapf(err, data.S.InSelectMarkingLineOld)
	}
	return MarkingLines, MapIdEntity, nil
}

func UpdateMarkingLine(db *sql.DB) {
	if err := (func() error {
		now, MapIdEntity, err := SelectMarkingLineNew(db)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorRead)
		}
		old, _, err := SelectMarkingLineOld(db)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorRead)
		}
		update := func() error {
			for _, valO := range old {
				for j, valN := range now {
					if reflect.DeepEqual(valO.Hierarchy, valN.Hierarchy) {
						now[j].Id = valO.Id
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
					for number, entityId := range valN.Hierarchy {
						QwStr2 := data.InsertMarkingLine(now[j].Id, entityId, int8(number+1))
						if err := db.Ping(); err != nil {
							return errors.Wrap(err, data.S.ErrorPingDB)
						}
						if _, err := db.Exec(QwStr2); err != nil {
							return errors.Wrap(err, data.S.ErrorAddDB+QwStr)
						}
					}
					old = append(old, now[j])
				}
			}
			return nil
		}
		err = update()
		fmt.Println(now)
		fmt.Println(old)
		fmt.Println(MapIdEntity)
		return errors.Wrap(err, data.S.ErrorUpdate)
	}()); err != nil {
		err = errors.Wrap(err, data.S.ErrorUpdateMarkingLine)
		log.Println(data.S.Error, err)
		ErrorRunWindow(err.Error())
	}
}
