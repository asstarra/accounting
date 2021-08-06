package window

import (
	e "accounting/data/errors"
	l "accounting/data/log"
	"accounting/data/qwery"
	. "accounting/data/table"
	"accounting/data/text"
	. "accounting/window/data"
	"database/sql"

	// "fmt"
	"log"
	"sort"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

// Выборка из таблицы "Entity" по заданным параметрам.
// Порядок: Id, Title, Type, Enum, Mark, Specification, Note.
// isChange определяет, разрешено ли выбирать строчки, где тип сущности это заказ.
// Информация о дочерних сущностях не выбирается.
func SelectEntity(db *sql.DB, id *int64, title *string, eType *int16, enum *bool,
	mark *int8, spec, note *string, isChange bool) ([]*Entity, error) {
	arr := make([]*Entity, 0)
	if err := (func() error {
		QwStr := qwery.SelectEntity(id, title, eType, enum, mark, spec, note, isChange)
		if err := db.Ping(); err != nil { // Пинг БД.
			return errors.Wrap(err, e.Err.ErrorPingDB)
		}
		rows, err := db.Query(QwStr) // Запрос к БД.
		if err != nil {
			return errors.Wrapf(err, e.Err.ErrorQueryDB, QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			var row Entity = NewEntity()
			err := rows.Scan(&row.Id, &row.Title, &row.Type, &row.Enumerable, &row.Marking, &row.Specification, &row.Note)
			if err != nil {
				return errors.Wrap(err, e.Err.ErrorDecryptRow)
			}
			arr = append(arr, &row)
		}
		return nil
	}()); err != nil {
		return arr, qwery.Wrapf(err, l.In.InSelectEntity, id, title, eType, enum, mark, isChange)
	}
	return arr, nil
}

// Сруктура, содержащая модель таблицы.
type modelEntityComponent struct {
	walk.TableModelBase
	walk.SorterBase
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
	log.Println(l.Panic, e.UnexpectedColumn)
	panic(e.UnexpectedColumn)
}

func (m *modelEntityComponent) Sort(col int, order walk.SortOrder) error {
	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]
		c := func(ls bool) bool {
			if order == walk.SortAscending {
				return ls
			}
			return !ls
		}
		switch col {
		case 0:
			return c(a.Title < b.Title)
		case 1:
			return c(a.Count < b.Count)
		}
		log.Println(l.Panic, e.UnexpectedColumn) // Лог.
		panic(e.UnexpectedColumn)
	})
	return m.SorterBase.Sort(col, order)
}

func (m *modelEntityComponent) Equal(row, col int, itemPtr interface{}) bool {
	val, ok := itemPtr.(*EntityRecChild)
	if !ok || val == nil {
		log.Println(l.Panic, e.WrongType) // Лог.
		panic(e.WrongType)
	}
	item := m.items[row]
	switch col {
	case 0:
		return item.Title == val.Title
	case 1:
		return item.Count == val.Count
	}
	log.Println(l.Panic, e.UnexpectedColumn) // Лог.
	panic(e.UnexpectedColumn)
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
	wf.modelType, wf.mapIdToTitle, err = SelectId16Title(db, TableEntityType, nil, nil)
	if err != nil {
		err = errors.Wrap(err, e.Err.ErrorTypeInit)
		return nil, err
	}
	return wf, nil
}

// Описание и запуск диалогового окна.
func EntityRunDialog(owner walk.Form, db *sql.DB, entity *Entity) (int, error) {
	log.Printf(l.BeginWindow, l.Entity) // Лог.
	var databind *walk.DataBinder       // Инициализация
	wf, err := newWindowsFormEntity(db, entity)
	if err != nil {
		return 0, errors.Wrap(err, e.Err.ErrorInit)
	}
	log.Printf(l.InitWindow, l.Entity) // Лог.
	if err := (dec.Dialog{             // Описание окна.
		AssignTo: &wf.Dialog,           // Привязка окна.
		Title:    text.T.HeadingEntity, // Название.
		DataBinder: dec.DataBinder{ // Привязка к структуре.
			AssignTo:       &databind,
			Name:           "entity",
			DataSource:     entity,
			ErrorPresenter: dec.ToolTipErrorPresenter{},
		},
		Layout: dec.VBox{MarginsZero: true},
		Children: []dec.Widget{
			dec.HSplitter{
				Children: []dec.Widget{
					dec.Composite{ // Левая половина.
						Layout: dec.Grid{Columns: 2},
						Children: []dec.Widget{
							dec.Label{ // Лэйбэл название.
								Text: text.T.LabelTitle,
							},
							dec.LineEdit{ // Текстова строка для названия.
								MaxLength: 255,
								MinSize:   dec.Size{170, 0},
								Text:      dec.Bind("Title"),
							},

							dec.Label{ // Лэйбэл тип.
								Text: text.T.LabelType,
							},
							dec.ComboBox{ // Выпадающий список для выбора типа.
								Value:         dec.Bind("Type", dec.SelRequired{}),
								BindingMember: "Id",
								DisplayMember: "Title",
								Model:         wf.modelType,
							},

							dec.CheckBox{
								ColumnSpan:     2,
								Text:           text.T.LabelEnumerable,
								TextOnLeftSide: true,
								Checked:        dec.Bind("Enumerable"),
							},

							dec.RadioButtonGroupBox{ // Выбор способа маркировки
								ColumnSpan: 2,
								Title:      text.T.LabelMarking,
								Layout:     dec.HBox{},
								DataMember: "Marking",
								Buttons: []dec.RadioButton{
									{Text: MarkingNo.Title(), Value: MarkingNo},
									{Text: MarkingAll.Title(), Value: MarkingAll},
									{Text: MarkingYear.Title(), Value: MarkingYear},
								},
							},

							dec.Label{ // Лэйбэл спецификация.
								Text: text.T.LabelSpecification,
							},
							dec.LineEdit{ // Текстовая строка.
								MaxLength: 255,
								Text:      dec.Bind("Specification"),
							},

							dec.Label{ // Лэйбэл примичание.
								ColumnSpan: 2,
								Text:       text.T.LabelNote,
							},
							dec.TextEdit{ // Текстовое поле для примечания.
								ColumnSpan: 2,
								MaxLength:  1023,
								MinSize:    dec.Size{0, 100},
								Text:       dec.Bind("Note"),
							},
						},
					},
					dec.Composite{ // Правая половина.
						Layout:  dec.Grid{Columns: 2},
						MinSize: dec.Size{230, 300},
						Children: []dec.Widget{
							dec.Label{ // Лэйбэл компоненты.
								ColumnSpan: 2,
								Text:       text.T.LabelComponents,
							},
							dec.TableView{ // Таблица с дочерними компонентами.
								AssignTo: &wf.tv, // Привязка к виджету.
								Columns: []dec.TableViewColumn{
									{Title: text.T.ColumnTitle, Width: 120},
									{Title: text.T.ColumnCount, Width: 80},
								},
								ColumnSpan: 2,
								Model:      wf.modelTable, // Привязка к модели.
							},

							dec.PushButton{ // Кнопка добавить дочерний компонент.
								ColumnSpan: 2,
								Text:       text.T.ButtonAdd + text.T.SuffixComponent,
								OnClicked: func() {
									log.Println(l.Info, l.LogAdd)              // Лог.
									if err := wf.add(db, entity); err != nil { // Обработка ошибок.
										MsgBoxError(wf, err, e.Err.ErrorAddRow)
									}
								},
							},
							dec.PushButton{ // Кнопка изменения дочернего компонента.
								ColumnSpan: 2,
								Text:       text.T.ButtonChange + text.T.SuffixComponent,
								OnClicked: func() {
									log.Println(l.Info, l.LogChange)              // Лог.
									if err := wf.change(db, entity); err != nil { // Обработка ошибок.
										MsgBoxError(wf, err, e.Err.ErrorChangeRow)
									}
								},
							},
							dec.PushButton{ // Кнопка удаления дочернего компонента.
								ColumnSpan: 2,
								Text:       text.T.ButtonDelete + text.T.SuffixComponent,
								OnClicked: func() {
									log.Println(l.Info, l.LogDelete)              // Лог.
									if err := wf.delete(db, entity); err != nil { // Обработка ошибок.
										MsgBoxError(wf, err, e.Err.ErrorDeleteRow)
									}
								},
							},

							dec.PushButton{ // Кнопка Ок.
								Text: text.T.ButtonOK,
								OnClicked: func() {
									log.Println(l.Info, l.LogOk)              // Лог.
									if err := databind.Submit(); err != nil { // Обработка ошибок.
										MsgBoxError(wf, err, e.Err.ErrorSubmit)
										return
									}
									if IsStringEmpty(wf, entity.Title) {
										return
									}
									entity.Children = wf.modelTable.items
									wf.Accept()
								},
							},
							dec.PushButton{ // Кнопка отмена.
								Text: text.T.ButtonCansel,
								OnClicked: func() {
									log.Println(l.Info, l.LogCansel) // Лог.
									wf.Cancel()
								},
							},
						},
					},
				},
			},
		},
	}.Create(owner)); err != nil {
		err = errors.Wrap(err, e.Err.ErrorCreateWindow) // Обработка ошибок создания окна.
		return 0, err
	}
	log.Printf(l.CreateWindow, l.Entity) // Лог.

	log.Printf(l.RunWindow, l.Entity) // Лог.
	return wf.Run(), nil              // Запуск окна.
}

// Функция, для добавления строки в таблицу.
func (wf *windowsFormEntity) add(db *sql.DB, entity *Entity) error {
	child := EntityRecChild{}
	cmd, err := EntityRecRunDialog(wf, db, false, &child)
	log.Printf(l.EndWindow, l.EntityRec, cmd) // Лог.
	if err != nil {
		return errors.Wrapf(err, l.In.InEntityRecRunDialog, child)
	}
	if cmd != walk.DlgCmdOK {
		return nil
	}
	if IsRepeat(wf, wf.modelTable, []int{1}, &child) {
		return nil
	}
	// Изменение БД при добавлении дочерней детали в составе(entity_rec), когда
	// родительская деталь не внесена БД происходит в entities.go при добавлении
	// родительской детали в БД.
	if entity.Id != 0 { // Здесь добавляем дочернюю сущность при известной родительской.
		// Проверяем на зацикленость состав.
		s, err := checkEntityRec(db, entity.Id, []*EntityRecChild{&child})
		if err != nil {
			return err
		}
		if s != "" {
			s = wf.mapIdToTitle[entity.Type] + " " + entity.Title + " -> " + s
			return errors.Wrap(errors.New(s), e.Err.ErrorGraphCircle)
		}
		QwStr := qwery.InsertEntityRec(entity.Id, child.Id, child.Count)
		if err = db.Ping(); err != nil { // Пинг БД.
			return errors.Wrap(err, e.Err.ErrorPingDB)
		}
		result, err := db.Exec(QwStr) // Запрос к БД.
		if err != nil {
			return errors.Wrapf(err, e.Err.ErrorAddDB, QwStr)
		}
		if id, err := result.LastInsertId(); err != nil {
			MsgBoxNotInsertedId(wf, err)
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
	if IsCorrectIndex(wf, wf.modelTable, wf.tv) {
		return nil
	}
	index := wf.tv.CurrentIndex()
	child := wf.modelTable.items[index]
	cmd, err := EntityRecRunDialog(wf, db, true, child)
	log.Printf(l.EndWindow, l.EntityRec, cmd) // Лог.
	if err != nil {
		return errors.Wrapf(err, l.In.InEntityRecRunDialog, child)
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
			return errors.Wrap(err, e.Err.ErrorPingDB)
		}
		if _, err := db.Exec(QwStr); err != nil { // Запрос к БД.
			return errors.Wrapf(err, e.Err.ErrorChangeDB, QwStr)
		}
	}
	wf.modelTable.PublishRowsReset() // Обновление таблицы.
	wf.modelTable.PublishRowChanged(index)
	return nil
}

// Функция, для удаления строки из таблицы.
func (wf *windowsFormEntity) delete(db *sql.DB, entity *Entity) error {
	if IsCorrectIndex(wf, wf.modelTable, wf.tv) {
		return nil
	}
	index := wf.tv.CurrentIndex()
	id := wf.modelTable.items[index].Id
	// Изменение БД при удалении дочерней детали в составе(entity_rec), когда
	// родительская деталь не внесена БД не происходит, т.к. этих данных нет в БД.
	if entity.Id != 0 {
		QwStr := qwery.DeleteEntityRec(entity.Id, id)
		if err := db.Ping(); err != nil { // Пинг БД.
			return errors.Wrap(err, e.Err.ErrorPingDB)
		}
		if _, err := db.Exec(QwStr); err != nil { // Запрос к БД.
			return errors.Wrapf(err, e.Err.ErrorDeleteDB, QwStr)
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
			return "", errors.Wrap(err, e.Err.ErrorSubquery)
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

func deleteMarkingLineEnd(db *sql.DB, end int64) (bool, error) {

	return false, nil
}
