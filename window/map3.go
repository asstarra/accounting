package window

import (
	"accounting/data"
	"database/sql"
	"fmt"
	"log"
	"sort"

	"github.com/pkg/errors"
)

type Map3 struct {
	mapIdToMarkingLine map[int64]*MarkingLine
	mapIdToEntity      map[int64]*Entity
	mapIdToEntityType  map[int64]string
	ToMarkingIds       func(order, detail, line int64) []int64
}

func NewMap3(db *sql.DB, isAllEntities bool) (Map3, error) {
	var m Map3
	var err error
	m.mapIdToMarkingLine, m.mapIdToEntity, err = UpdateMarkingLine(db, isAllEntities)
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorSubquery)
		return m, err
	}
	_, m.mapIdToEntityType, err = SelectIdTitle(db, "EntityType")
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorTypeInit)
		return m, err
	}

	// Определение функции ToMarkingIds, сохраняющей результаты предыдущих вычислений.
	type Id3 struct {
		order, detail, line int64
	}
	mapId3ToMarkingLine := make(map[Id3][]int64, 0)
	m.ToMarkingIds = func(order, detail, line int64) []int64 {
		// fmt.Println(mapId3ToMarkingLine)
		id3 := Id3{order: order, detail: detail, line: line}
		if arr, ok := mapId3ToMarkingLine[id3]; ok {
			return arr
		}
		mapId3ToMarkingLine[id3] = m.toMarkingIds(order, detail, line)
		return mapId3ToMarkingLine[id3]
	}
	return m, nil
}

// По ид заказа, детали и линии собирает все допустимые индификаторы линий
func (m *Map3) toMarkingIds(order, detail, line int64) []int64 {
	markings := make([]int64, 0, 20)
	if line != 0 {
		markings = append(markings, line)
	} else {
		for _, val := range m.mapIdToMarkingLine {
			if val.Hierarchy[0].Id != order && order != 0 {
				continue
			}
			if val.Hierarchy[len(val.Hierarchy)-1].Id != detail && detail != 0 {
				continue
			}
			markings = append(markings, val.Id)
		}
	}
	sort.Slice(markings, func(i, j int) bool {
		si := m.MarkingToString(markings[i])
		sj := m.MarkingToString(markings[j])
		return si < sj
	})
	return markings
}

// Переводит индификатор маркировочной линии в список указателей на структуру IdTitle
func (m *Map3) MarkingsToIdTitles(ids []int64) []*IdTitle {
	arr := make([]*IdTitle, 0, len(ids))
	for _, val := range ids {
		arr = append(arr, &IdTitle{Id: val, Title: m.MarkingToString(val)})
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Title < arr[j].Title
	})
	return arr
}

// Переводит индификатор маркировочной линии в строку с информацией о ней
func (m *Map3) MarkingToString(id int64) string {
	var s string
	mline, ok := m.mapIdToMarkingLine[id]
	if !ok {
		log.Println(data.Log.Warning, "В карте mapIdToMarkingLine не найдено значение ", id)
		return ""
	}
	for _, val := range mline.Hierarchy {
		s += fmt.Sprintf(" -> (%s x%d)", m.EntityToString(val.Id), val.Count)
	}
	s = s[4:]
	// s = fmt.Sprintf("%d {%s}", id, s)
	return s
}

// Переводит индификатор сущности в строку.
func (m *Map3) EntityToString(id int64) string {
	e, ok := m.mapIdToEntity[id]
	if !ok {
		log.Println(data.Log.Warning, "В карте mapIdToEntity не найдено значение ", id)
		return ""
	}
	eType, ok := m.mapIdToEntityType[e.Type]
	if !ok {
		log.Println(data.Log.Warning, "В карте mapIdToEntityType не найдено значение ", e.Type)
		return ""
	}
	return fmt.Sprintf("%s %s", eType, e.Title)
}

// Переводит структуру MarkedDetailMin в строку.
func (m *Map3) MarkedDetailMinToString(md MarkedDetailMin) string {
	if md.Id == 0 {
		return "Нет"
	}
	line := m.mapIdToMarkingLine[md.Marking]
	if line == nil {
		log.Println(data.Log.Warning, "В карте mapIdToMarkingLine не найдено значение ", md.Marking)
		return "ERROR"
	}
	eId := line.Hierarchy[len(line.Hierarchy)-1].Id
	// e := m.mapIdToEntity[eId]
	// if e == nil {
	// 	log.Println(data.Log.Warning, "В карте mapIdToEntity не найдено значение ", eId)
	// 	return "ERROR"
	// }
	// return fmt.Sprintf("%s %s %s", m.mapIdToEntityType[e.Type], e.Title, md.Mark)
	return fmt.Sprintf("%s %s", m.EntityToString(eId), md.Mark)
}

// Переводит список индификаторов маркировочных линий в список ид сущностей
// с названием заказов.
func (m *Map3) IdsToModelOrders(ids []int64, isAll bool) []*IdTitle {
	mapOrder := make(map[int64]bool, len(ids)/5)
	for _, val := range ids {
		line := m.mapIdToMarkingLine[val]
		mapOrder[line.Hierarchy[0].Id] = true
	}
	arr := make([]*IdTitle, 0, len(mapOrder)+1)
	if isAll {
		arr = append(arr, &IdTitle{})
	}
	for key, _ := range mapOrder {
		arr = append(arr, &IdTitle{Id: key, Title: m.EntityToString(key)})
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Title < arr[j].Title
	})
	return arr
}

// По ид заказа, детали и линии собирает список индификаторов с названием заказов.
func (m *Map3) ModelOrders(detail, line int64, isAll bool) []*IdTitle {
	return m.IdsToModelOrders(m.ToMarkingIds(0, detail, line), isAll)
}

// Переводит список индификаторов маркировочных линий в список ид сущностей
// с названием деталей.
func (m *Map3) IdsToModelDetails(ids []int64, isAll bool) []*IdTitle {
	mapDetails := make(map[int64]bool, len(ids)/2)
	for _, val := range ids {
		line := m.mapIdToMarkingLine[val]
		mapDetails[line.Hierarchy[len(line.Hierarchy)-1].Id] = true
	}
	arr := make([]*IdTitle, 0, len(mapDetails)+1)
	if isAll {
		arr = append(arr, &IdTitle{})
	}
	for key, _ := range mapDetails {
		arr = append(arr, &IdTitle{Id: key, Title: m.EntityToString(key)})
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Title < arr[j].Title
	})
	return arr
}

// По ид заказа, детали и линии собирает список индификаторов с названием деталей.
func (m *Map3) ModelDetails(order, line int64, isAll bool) []*IdTitle {
	return m.IdsToModelDetails(m.ToMarkingIds(order, 0, line), isAll)
}

// По ид заказа, детали и линии собирает список индификаторов с названием маркировочных линий.
func (m *Map3) ModelMarkingLines(order, detail int64, isAll bool) []*IdTitle {
	ids := m.ToMarkingIds(order, detail, 0)
	if isAll {
		ids = append(ids, 0)
	}
	return m.MarkingsToIdTitles(ids)
}

// Переводит IdTitle в []int64
func (m *Map3) ModelToIds(model []*IdTitle) []int64 {
	arr := make([]int64, 0, len(model))
	for _, val := range model {
		arr = append(arr, val.Id)
	}
	return arr
}

// Возращает ид заказа по ид маркировочной линии.
func (m *Map3) MarkingLineToOrder(line int64) int64 {
	return m.mapIdToMarkingLine[line].Hierarchy[0].Id
}

// Возращает ид детали по ид маркировочной линии.
func (m *Map3) MarkingLineToDetail(line int64) int64 {
	hierarchy := m.mapIdToMarkingLine[line].Hierarchy
	return hierarchy[len(hierarchy)-1].Id
}

// func (m *Map3) ParentMarkingLineIds(childLine int64) []int64 {
// 	child := m.mapIdToMarkingLine[childLine]
// 	if child == nil {
// 		return []int64{}
// 	}
// 	fmt.Println(m.MarkingToString(child.Id))
// 	for _, parent := range m.mapIdToMarkingLine {
// 		if len(child.Hierarchy)-1 == len(parent.Hierarchy) {
// 			var equal bool = true
// 			for i, val := range parent.Hierarchy {
// 				if child.Hierarchy[i] != val {
// 					equal = false
// 				}
// 			}
// 			if equal {
// 				return []int64{parent.Id}
// 			}
// 		}
// 	}
// 	arr := make([]int64, 0, 20)
// 	if len(child.Hierarchy) == 2 && child.Hierarchy[0].Count == 0 {
// 		ids := m.ToMarkingIds(child.Hierarchy[0].Id, child.Hierarchy[len(child.Hierarchy)-1].Id, 0)
// 		for _, val := range ids {
// 			arr = append(arr, m.ParentMarkingLineIds(val)...)
// 		}
// 	}
// 	fmt.Println(arr)
// 	return arr
// }
