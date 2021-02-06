package window

import (
	"accounting/data"
	"database/sql"
	"fmt"
	"log"

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

func SelectMarkingLine(db *sql.DB) ([]*MarkingLine, error) {
	arr := make([]*MarkingLine, 0)
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
			row := MarkingLine{Entities: make([]*Entity, 0, 10)}
			err := rows.Scan(&row.Id)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorDecryptRow)
			}
			arr = append(arr, &row)
		}

		for index, val := range arr {
			mapNumberEntity, err := SelectMarkingLineEntity(db, val.Id)
			if err != nil {
				return err
			}
			for number := 1; number <= len(mapNumberEntity); number++ {
				arr[index].Entities = append(arr[index].Entities, mapNumberEntity[int8(number)])
			}
		}
		return nil
	}()); err != nil {
		return arr, errors.Wrapf(err, data.S.InSelectMarkingLine)
	}
	return arr, nil
}

type markingLineGraph struct {
	entities   []*Entity
	mapIdIndex map[int64]int
	linesN     []*MarkingLine
	linesO     []*MarkingLine
}

func (m markingLineGraph) String() string {
	return fmt.Sprintf("{entities: %v,\n\nmap: %v,\n\nlinesN: %v,\n\nlinesO: %v}", m.entities, m.mapIdIndex, m.linesN, m.linesO)
}

func newMarkingLineGraph(db *sql.DB) (markingLineGraph, error) {
	var err error
	g := markingLineGraph{}
	g.entities, err = SelectEntities(db, "", 0, true)
	if err != nil {
		return g, err
	}
	g.mapIdIndex = make(map[int64]int, len(g.entities))
	for index, val := range g.entities {
		g.mapIdIndex[val.Id] = index
	}
	entityRec, err := SelectEntityRec(db)
	if err != nil {
		return g, err
	}
	for _, val := range entityRec {
		indexP := g.mapIdIndex[val.IdP]
		g.entities[indexP].Children = append(g.entities[indexP].Children, &val.EntityRecChild)
	}
	g.linesN = make([]*MarkingLine, 0)
	for _, val := range g.entities {
		if val.Type == 1 {
			g.createLineRec([]*Entity{val}, val.Children)
		}
	}
	if g.linesO, err = SelectMarkingLine(db); err != nil {
		return g, err
	}
	return g, nil
}

func (g *markingLineGraph) createLineRec(entities []*Entity, children []*EntityRecChild) {
	for _, val := range children {
		entity := g.entities[g.mapIdIndex[val.Id]]
		entities2 := make([]*Entity, 0, 10)
		entities2 = append(entities2, entities...)
		entities2 = append(entities2, entity)
		if entity.Marking != MarkingNo {
			g.linesN = append(g.linesN, &MarkingLine{Entities: entities2})
		}
		g.createLineRec(entities2, entity.Children)
	}
}

func (g *markingLineGraph) update(db *sql.DB) error {
	Equal := func(a, b []*Entity) bool {
		if len(a) != len(b) {
			return false
		}
		for index, val := range a {
			if val.Id != b[index].Id {
				return false
			}
		}
		return true
	}
	for _, valO := range g.linesO {
		for j, valN := range g.linesN {
			if Equal(valO.Entities, valN.Entities) {
				g.linesN[j].Id = valO.Id
				break
			}
		}
	}
	for _, valN := range g.linesN {
		if valN.Id == 0 {
			QwStr := data.InsertMarking()
			if err := db.Ping(); err != nil {
				return errors.Wrap(err, data.S.ErrorPingDB)
			}
			result, err := db.Exec(QwStr)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorAddDB+QwStr)
			}
			idM, err := result.LastInsertId()
			if err != nil {
				return errors.Wrap(err, data.S.ErrorInsertIndexLog)
			}

			for number, entity := range valN.Entities {
				QwStr2 := data.InsertMarkingLine(idM, entity.Id, int8(number+1))
				if err := db.Ping(); err != nil {
					return errors.Wrap(err, data.S.ErrorPingDB)
				}
				if _, err := db.Exec(QwStr2); err != nil {
					return errors.Wrap(err, data.S.ErrorAddDB+QwStr)
				}
			}
		}
	}
	return nil
}

func UpdateMarkingLine(db *sql.DB) {
	if err := (func() error {
		g, err := newMarkingLineGraph(db)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorRead)
		}
		err = g.update(db)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorUpdate)
		}
		fmt.Println(g)
		return nil
	}()); err != nil {
		err = errors.Wrap(err, data.S.ErrorUpdateMarkingLine)
		log.Println(data.S.Error, err)
		ErrorRunWindow(err.Error())
	}
}
