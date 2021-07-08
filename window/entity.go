package window

import (
	"accounting/data"
	"accounting/data/qwery"
	"database/sql"

	"fmt"
	"log"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

// Тип маркировки компонента
type Marking int8

const (
	MarkingNo   Marking = 1 + iota // Не маркируется.
	MarkingAll                     // Маркировка сквозная.
	MarkingYear                    // Маркировка по годам.
)

var MapMarkingToTitle = map[Marking]string{ // TO-DO
	MarkingNo:   "Нет",
	MarkingAll:  "Сквозная",
	MarkingYear: "По годам",
}

type Entity struct {
	Id            int64             // Ид.
	Title         string            // Название.
	Type          int16             // Ид типа.
	Enumerable    bool              // Можно посчитать?.
	Marking       Marking           // Способ маркировки.
	Specification string            // Спицификация.
	Note          string            // Примечание.
	Children      []*EntityRecChild // Дочерние детали: ид с описанием.
}

func NewEntity() Entity {
	return Entity{Children: make([]*EntityRecChild, 0, 0)}
}

func (e Entity) String() string {
	return fmt.Sprintf("{Id = %d, Title = '%s', Type = %d, Enum = %v, Mark = %s, Spec = '%s', Note = '%s', Children = %v}\n",
		e.Id, e.Title, e.Type, e.Enumerable, MapMarkingToTitle[e.Marking], e.Specification, e.Note, e.Children)
}

// Выборка из таблицы Entity всех ее полей удовлетворяющих условию, GO-TO
// где в значения поля Title входит title,
// значение поля Type равно entityType (при 0, разрешен любой тип),
// isChange определяет, разрешено ли выбирать строчки, где тип сущности это заказ.
// Информация о дочерних сущностях не выбирается.
func SelectEntity(db *sql.DB, id *int64, title *string, eType *int16, enum *bool,
	mark *int8, spec, note *string, isChange bool) ([]*Entity, error) {
	arr := make([]*Entity, 0)
	if err := (func() error {
		QwStr := qwery.SelectEntity(id, title, eType, enum, mark, spec, note, isChange)
		if err := db.Ping(); err != nil { // Пинг БД.
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		rows, err := db.Query(QwStr) // Запрос к БД.
		if err != nil {
			return errors.Wrapf(err, data.S.ErrorQueryDB, QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			var row Entity = NewEntity()
			err := rows.Scan(&row.Id, &row.Title, &row.Type, &row.Enumerable, &row.Marking, &row.Specification, &row.Note)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorDecryptRow)
			}
			arr = append(arr, &row)
		}
		return nil
	}()); err != nil {
		return arr, errors.Wrapf(err, data.Log.InSelectEntities, title, eType) //GO-TO
	}
	return arr, nil
}

// Сруктура, содержащая модель таблицы.
type modelEntityComponent struct { //GO-TO rename.
	walk.TableModelBase
	items []*EntityRecChild // Список дочерних компонентов.
}

func (m *modelEntityComponent) RowCount() int {
	return len(m.items)
}

func (m *modelEntityComponent) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return item.Title
	case 1:
		return item.Count
	}
	log.Println(data.Log.Panic, data.S.ErrorUnexpectedColumn)
	panic(data.S.ErrorUnexpectedColumn)
}

// Структура, содержащая описание и переменные окна.
type windowsFormEntity struct {
	*walk.Dialog
	modelType    []*Id16Title          // Модель выпадающего списка, содержащая типы сущности.
	modelTable   *modelEntityComponent // Модель таблицы, содержащей дочерние компоненты.
	tv           *walk.TableView       // Виджет таблицы, содержащей дочерние компоненты.
	mapIdToTitle map[int16]string      // Отображение из Id в название типа.
}

// Инициализация модели окна.
func newWindowsFormEntity(db *sql.DB, entity *Entity) (*windowsFormEntity, error) {
	var err error
	wf := new(windowsFormEntity)
	wf.modelTable = new(modelEntityComponent)
	wf.modelTable.items = entity.Children
	wf.modelType, wf.mapIdToTitle, err = SelectId16Title(db, "EntityType", nil, nil)
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorTypeInit)
		return nil, err
	}
	return wf, nil
}

// Описание и запуск диалогового окна.
func EntityRunDialog(owner walk.Form, db *sql.DB, entity *Entity) (int, error) {
	log.Printf(data.Log.BeginWindow, data.Log.Entity) // Лог.
	sButtonAdd := " компонент"                        // GO-TO возможно нужно? вынести строку
	var databind *walk.DataBinder                     // Инициализация
	wf, err := newWindowsFormEntity(db, entity)
	if err != nil {
		return 0, errors.Wrap(err, data.S.ErrorInit)
	}
	log.Printf(data.Log.InitWindow, data.Log.Entity) // Лог.
	if err := (dec.Dialog{                           // Описание окна.
		AssignTo: &wf.Dialog,           // Привязка окна.
		Title:    data.S.HeadingEntity, // Название.
		DataBinder: dec.DataBinder{ // Привязка к структуре.
			AssignTo:       &databind,
			Name:           "entity",
			DataSource:     entity,
			ErrorPresenter: dec.ToolTipErrorPresenter{},
		},
		Layout: dec.VBox{MarginsZero: true},
		Children: []dec.Widget{
			dec.HSplitter{ // Левая половина.
				Children: []dec.Widget{
					dec.Composite{
						Layout: dec.Grid{Columns: 2},
						Children: []dec.Widget{
							dec.Label{ // Лэйбэл название.
								Text: "Название:", // TO-DO
							},
							dec.LineEdit{ // Текстова строка для названия.
								MaxLength: 255,
								MinSize:   dec.Size{170, 0},
								Text:      dec.Bind("Title"),
							},

							dec.Label{ // Лэйбэл тип.
								Text: "Тип:", // TO-DO
							},
							dec.ComboBox{ // Выпадающий список для выбора типа.
								Value:         dec.Bind("Type", dec.SelRequired{}),
								BindingMember: "Id",
								DisplayMember: "Title",
								Model:         wf.modelType,
							},

							dec.Label{ // Лэйбэл считаемость. // TO-DO
								ColumnSpan: 2,
								Text:       "Считаемость:", // TO-DO
							},

							dec.RadioButtonGroupBox{ // Выбор способа маркировки
								ColumnSpan: 2,
								Title:      "Маркировка:", // TO-DO
								Layout:     dec.HBox{},
								DataMember: "Marking",
								Buttons: []dec.RadioButton{
									{Text: MapMarkingToTitle[MarkingNo], Value: MarkingNo},
									{Text: MapMarkingToTitle[MarkingAll], Value: MarkingAll},
									{Text: MapMarkingToTitle[MarkingYear], Value: MarkingYear},
								},
							},

							dec.Label{ // Лэйбэл спецификация.
								Text: "Спецификация:", // TO-DO
							},
							dec.LineEdit{ // Текстовая строка.
								MaxLength: 255,
								Text:      dec.Bind("Specification"),
							},

							dec.Label{ // Лэйбэл примичание.
								ColumnSpan: 2,
								Text:       "Примечание:", // TO-DO
							},
							dec.TextEdit{ // Текстовое поле для примечания.
								ColumnSpan: 2,
								MaxLength:  1023,
								MinSize:    dec.Size{0, 100},
								Text:       dec.Bind("Note"),
							},
						},
					},
					dec.Composite{
						Layout:  dec.Grid{Columns: 2},
						MinSize: dec.Size{230, 300},
						Children: []dec.Widget{
							dec.Label{ // Лэйбэл компоненты.
								ColumnSpan: 2,
								Text:       "Компоненты:", // TO-DO
							},
							dec.TableView{ // Таблица с дочерними компонентами.
								AssignTo: &wf.tv, // Привязка к виджету.
								Columns: []dec.TableViewColumn{
									{Title: "Название"},   // TO-DO
									{Title: "Количество"}, // TO-DO
								},
								ColumnSpan: 2,
								Model:      wf.modelTable, // Привязка к модели.
							},

							dec.PushButton{ // Кнопка добавить дочерний компонент.
								ColumnSpan: 2,
								Text:       data.S.ButtonAdd + sButtonAdd,
								OnClicked: func() {
									log.Println(data.Log.Info, data.Log.LogAdd) // Лог.
									if err := wf.add(db, entity); err != nil {
										err = errors.Wrap(err, data.S.ErrorAddRow) // Обработка ошибок.
										log.Println(data.Log.Error, err)           // Лог.
										walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
									}
								},
							},
							dec.PushButton{ // Кнопка изменения дочернего компонента.
								ColumnSpan: 2,
								Text:       data.S.ButtonChange + sButtonAdd,
								OnClicked: func() {
									log.Println(data.Log.Info, data.Log.LogChange) // Лог.
									if err := wf.change(db, entity); err != nil {
										err = errors.Wrap(err, data.S.ErrorChangeRow) // Обработка ошибок.
										log.Println(data.Log.Error, err)              // Лог.
										walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
									}
								},
							},
							dec.PushButton{ // Кнопка удаления дочернего компонента.
								ColumnSpan: 2,
								Text:       data.S.ButtonDelete + sButtonAdd,
								OnClicked: func() {
									log.Println(data.Log.Info, data.Log.LogDelete) // Лог.
									if err := wf.delete(db, entity); err != nil {
										err = errors.Wrap(err, data.S.ErrorDeleteRow) // Обработка ошибок.
										log.Println(data.Log.Error, err)              // Лог.
										walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
									}
								},
							},

							dec.PushButton{ // Кнопка Ок.
								Text: data.S.ButtonOK,
								OnClicked: func() {
									log.Println(data.Log.Info, data.Log.LogOk) // Лог.
									if err := databind.Submit(); err != nil {
										err = errors.Wrap(err, data.S.ErrorSubmit) // Обработка ошибок.
										log.Println(data.Log.Error, err)           // Лог.
										walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
										return
									}
									entity.Children = wf.modelTable.items
									wf.Accept()
								},
							},
							dec.PushButton{ // Кнопка отмена.
								Text: data.S.ButtonCansel,
								OnClicked: func() {
									log.Println(data.Log.Info, data.Log.LogCansel) // Лог.
									wf.Cancel()
								},
							},
						},
					},
				},
			},
		},
	}.Create(owner)); err != nil {
		err = errors.Wrap(err, data.S.ErrorCreateWindow) // Обработка ошибок создания окна.
		return 0, err
	}
	log.Printf(data.Log.CreateWindow, data.Log.Entity) // Лог.

	log.Printf(data.Log.RunWindow, data.Log.Entity) // Лог.
	return wf.Run(), nil                            // Запуск окна.
}

// Функция, для добавления строки в таблицу.
func (wf *windowsFormEntity) add(db *sql.DB, entity *Entity) error {
	child := EntityRecChild{}
	cmd, err := EntityRecRunDialog(wf, db, false, &child)
	log.Printf(data.Log.EndWindow, data.Log.EntityRec, cmd) // Лог.
	if err != nil {
		return errors.Wrapf(err, data.Log.InEntityRecRunDialog, child)
	}
	if cmd != walk.DlgCmdOK {
		return nil
	}
	// Изменение БД при добавлении дочерней детали в составе(entity_rec), когда
	// родительская деталь не внесена БД происходит в entities.go при добавлении
	// родительской детали в БД.
	if entity.Id != 0 { // Здесь добавляем дочернюю сущность при известной родительской.
		s, err := checkEntityRec(db, entity.Id, []*EntityRecChild{&child})
		if err != nil {
			return err
		}
		if s != "" {
			s = wf.mapIdToTitle[entity.Type] + " " + entity.Title + " -> " + s
			return errors.Wrap(errors.New(s), data.S.ErrorGraphCircle)
		}
		QwStr := qwery.InsertEntityRec(entity.Id, child.Id, child.Count)
		if err = db.Ping(); err != nil { // Пинг БД.
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		result, err := db.Exec(QwStr) // Запрос к БД.
		if err != nil {
			return errors.Wrapf(err, data.S.ErrorAddDB, QwStr)
		}
		if id, err := result.LastInsertId(); err != nil {
			log.Println(data.Log.Error, errors.Wrap(err, data.S.ErrorInsertIndexLog))        // Лог.
			walk.MsgBox(wf, data.S.MsgBoxError, data.S.ErrorInsertIndex, data.Icon.Critical) // TO-DO
		} else {
			child.Id = id
		}
	}
	// Обновление таблицы.
	trackLatest := wf.tv.ItemVisible(len(wf.modelTable.items) - 1) //&& len(wf.tv.SelectedIndexes()) <= 1
	wf.modelTable.items = append(wf.modelTable.items, &child)
	index := len(wf.modelTable.items) - 1
	wf.modelTable.PublishRowsReset()
	wf.modelTable.PublishRowsInserted(index, index)
	if trackLatest {
		wf.tv.EnsureItemVisible(len(wf.modelTable.items) - 1)
	}
	return nil
}

// Функция, для изменения строки в таблице.
func (wf *windowsFormEntity) change(db *sql.DB, entity *Entity) error {
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 {
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChooseRow, data.Icon.Info)
		return nil
	}
	index := wf.tv.CurrentIndex()
	child := wf.modelTable.items[index]
	cmd, err := EntityRecRunDialog(wf, db, true, child)
	log.Printf(data.Log.EndWindow, data.Log.EntityRec, cmd) // Лог.
	if err != nil {
		return errors.Wrapf(err, data.Log.InEntityRecRunDialog, child)
	}
	if cmd != walk.DlgCmdOK {
		return nil
	}
	// Изменение БД при изменении дочерней детали в составе(entity_rec), когда
	// родительская деталь не внесена БД происходит в entities.go при добавлении
	// родительской детали в БД.
	if entity.Id != 0 { // Здесь изменяем дочернюю сущность при известной родительской.
		QwStr := qwery.UpdateEntityRec(entity.Id, child.Id, child.Count)
		if err = db.Ping(); err != nil { // Пинг БД.
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		if _, err := db.Exec(QwStr); err != nil { // Запрос к БД.
			return errors.Wrapf(err, data.S.ErrorChangeDB, QwStr)
		}
	}
	wf.modelTable.PublishRowsReset() // Обновление таблицы.
	wf.modelTable.PublishRowChanged(index)
	return nil
}

// Функция, для удаления строки из таблицы.
func (wf *windowsFormEntity) delete(db *sql.DB, entity *Entity) error {
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 {
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChooseRow, data.Icon.Info)
		return nil
	}
	index := wf.tv.CurrentIndex()
	id := wf.modelTable.items[index].Id
	// Изменение БД при удалении дочерней детали в составе(entity_rec), когда
	// родительская деталь не внесена БД не происходит, т.к. этих данных нет в БД.
	if entity.Id != 0 {
		QwStr := qwery.DeleteEntityRec(entity.Id, id)
		if err := db.Ping(); err != nil { // Пинг БД.
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		if _, err := db.Exec(QwStr); err != nil { // Запрос к БД.
			return errors.Wrapf(err, data.S.ErrorDeleteDB, QwStr)
		}
	}
	// Обновление таблицы.
	trackLatest := wf.tv.ItemVisible(len(wf.modelTable.items) - 1) //&& len(wf.tv.SelectedIndexes()) <= 1
	wf.modelTable.items = wf.modelTable.items[:index+copy(wf.modelTable.items[index:], wf.modelTable.items[index+1:])]
	wf.modelTable.PublishRowsReset()
	wf.modelTable.PublishRowsRemoved(index, index)
	if trackLatest {
		wf.tv.EnsureItemVisible(len(wf.modelTable.items) - 1)
	}
	wf.modelTable.PublishRowsChanged(index, wf.modelTable.RowCount()-1)
	return nil
}

// Рекурсивная функция, для проверки состава (таблица EntityRec) на наличие циклов.
// Возращает строку, которая опиывает первый найденный цикл.
// Если циклов нет, возращает пустую строку.
func checkEntityRec(db *sql.DB, parent int64, children []*EntityRecChild) (string, error) {
	for _, val := range children { // Ищем родителя в списке дочерних компонентов.
		if parent == val.Id {
			return val.Title, nil
		}
	}
	for _, val := range children { // Для каждой дочерней сущности, вызываем функцию рекурсивно.
		_, children2, err := SelectEntityRecChild(db, &val.Id)
		if err != nil {
			return "", errors.Wrap(err, data.S.ErrorSubquery)
		}
		s, err := checkEntityRec(db, parent, children2)
		if err != nil {
			return "", err
		}
		if s == "" {
			return s, nil
		} else {
			return val.Title + " -> " + s, nil
		}
	}
	return "", nil
}
