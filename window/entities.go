package window

import (
	"accounting/data"
	e "accounting/data/errors"
	l "accounting/data/log"
	"accounting/data/qwery"
	. "accounting/data/table"
	"accounting/data/text"
	. "accounting/window/data"
	"database/sql"

	//"fmt"
	"log"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

// Сруктура, содержащая модель таблицы.
type modelEntitiesComponent struct {
	walk.TableModelBase
	items        []*Entity        // Список сущностей.
	mapIdToTitle map[int16]string // Отображения Id типа в его название.
}

func (m *modelEntitiesComponent) RowCount() int {
	return len(m.items)
}

func (m *modelEntitiesComponent) Value(row, col int) interface{} { // TO-DO
	item := m.items[row]
	switch col {
	case 0:
		return m.mapIdToTitle[item.Type]
	case 1:
		return item.Title
	case 2:
		return item.Specification
	case 3:
		return item.Marking.Title()
	case 4:
		return item.Note
	}
	log.Println(l.Panic, e.UnexpectedColumn)
	panic(e.UnexpectedColumn)
}

// Инициализация модели таблицы.
func newModelEntitiesComponent(db *sql.DB, isChange bool) (*modelEntitiesComponent, error) { // TO-DO
	var err error
	m := new(modelEntitiesComponent)
	m.items, err = SelectEntity(db, nil, nil, nil, nil, nil, nil, nil, isChange) // TO-DO
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Структура, содержащая описание и переменные окна.
type windowsFormEntities struct {
	*walk.Dialog
	modelType  []*Id16Title            // Модель выпадающего списка, содержащая типы сущности.
	modelTable *modelEntitiesComponent // Модель таблицы, содержащей компоненты.
	tv         *walk.TableView         // Виджет таблицы, содержащей компоненты.
}

// Инициализация модели окна.
func newWindowsFormEntities(db *sql.DB, isChange bool) (*windowsFormEntities, error) {
	var err error
	wf := new(windowsFormEntities)
	wf.modelTable, err = newModelEntitiesComponent(db, isChange)
	if err != nil {
		err = errors.Wrap(err, e.Err.ErrorTableInit)
		return nil, err
	}
	wf.modelType, wf.modelTable.mapIdToTitle, err = SelectId16Title(db, TableEntityType, nil, nil)
	if err != nil {
		err = errors.Wrap(err, e.Err.ErrorTypeInit)
		return nil, err
	}
	wf.modelType = append([]*Id16Title{new(Id16Title)}, wf.modelType...) // TO-DO оптмизировать по памяти.
	return wf, nil
}

// Описание и запуск диалогового окна.
func EntitiesRunDialog(owner walk.Form, db *sql.DB, isChange bool, idTitle *Id64Title) (int, error) {
	log.Printf(l.BeginWindow, l.Entities) // Лог.
	var err error
	var databind *walk.DataBinder
	search := new(Id16Title)
	wf, err := newWindowsFormEntities(db, isChange) // Инициализация.
	if err != nil {
		return 0, errors.Wrap(err, e.Err.ErrorInit)
	}
	log.Printf(l.InitWindow, l.Entities) // Лог.
	if err = (dec.Dialog{                // Описание окна.
		AssignTo: &wf.Dialog,             // Привязка окна.
		Title:    text.T.HeadingEntities, // Название.
		DataBinder: dec.DataBinder{ // Привязка к структуре
			AssignTo:       &databind,
			Name:           "search",
			DataSource:     search,
			ErrorPresenter: dec.ToolTipErrorPresenter{},
		},
		Layout:  dec.VBox{},
		MinSize: dec.Size{550, 0},
		Children: []dec.Widget{
			dec.Composite{
				Layout: dec.HBox{MarginsZero: true},
				Children: []dec.Widget{
					dec.Label{ // Лэйбэл название.
						Text: "Название:", // TO-DO
					},
					dec.LineEdit{ // Текстовая строка для ввода названия.
						Text: dec.Bind("Title"),
					},
					dec.HSpacer{Size: 20},

					dec.Label{ // Лэйбэл тип.
						Text: "Тип:", // TO-DO
					},
					dec.ComboBox{ // Выпадающий список для выбора типа.
						Value:         dec.Bind("Id", dec.SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Title",
						MinSize:       dec.Size{120, 0},
						Model:         wf.modelType,
					},
					dec.HSpacer{Size: 20},

					dec.PushButton{ // Кнопка поиск.
						Text: text.T.ButtonSearch,
						OnClicked: func() {
							log.Println(l.Info, l.LogSearch)          // Лог.
							if err := databind.Submit(); err != nil { // Обновление данных.
								err = errors.Wrap(err, e.Err.ErrorSubmit) // Обработка ошибок.
								log.Println(l.Error, err)                 // Лог.
								walk.MsgBox(wf, text.T.MsgBoxError, MsgError(err), data.Icon.Error)
								return
							}
							lastLen := wf.modelTable.RowCount()
							if items, err := SelectEntity(db, nil, &search.Title, &search.Id,
								nil, nil, nil, nil, isChange); err != nil { // Выборка из БД.
								err = errors.Wrap(err, e.Err.ErrorSubquery) // Обработка ошибок.
								log.Println(l.Error, err)                   // Лог.
								walk.MsgBox(wf, text.T.MsgBoxError, MsgError(err), data.Icon.Error)
								return
							} else {
								wf.modelTable.items = items
							}
							nowLen := wf.modelTable.RowCount() // Обновление таблицы.
							wf.modelTable.PublishRowsReset()
							wf.modelTable.PublishRowsRemoved(0, lastLen)
							wf.modelTable.PublishRowsInserted(0, nowLen)
							wf.modelTable.PublishRowsChanged(0, nowLen)
						},
					},
				},
			},
			dec.TableView{ // Таблица с компонентами.
				AssignTo: &wf.tv, // Привязка к виджету.
				Columns: []dec.TableViewColumn{ // TO-DO
					{Title: "Тип"},
					{Title: "Название"},
					{Title: "Спецификация"},
					{Title: "Маркировка"},
					{Title: "Примечание"},
				},
				MinSize: dec.Size{0, 200},
				Model:   wf.modelTable, // Привязка к модели.
			},
			dec.Composite{
				Layout:  dec.HBox{MarginsZero: true},
				Visible: isChange, // Видимость.
				Children: []dec.Widget{
					dec.PushButton{ // Кнопка добавить.
						Text: text.T.ButtonAdd,
						OnClicked: func() {
							log.Println(l.Info, l.LogAdd) // Лог.
							if err := wf.add(db); err != nil {
								err = errors.Wrap(err, e.Err.ErrorAddRow) // Обработка ошибок.
								log.Println(l.Error, err)                 // Лог.
								walk.MsgBox(wf, text.T.MsgBoxError, MsgError(err), data.Icon.Error)
							}
						},
					},
					dec.PushButton{ // Кнопка изменить.
						Text: text.T.ButtonChange,
						OnClicked: func() {
							log.Println(l.Info, l.LogChange) // Лог.
							if err := wf.change(db); err != nil {
								err = errors.Wrap(err, e.Err.ErrorChangeRow) // Обработка ошибок.
								log.Println(l.Error, err)                    // Лог.
								walk.MsgBox(wf, text.T.MsgBoxError, MsgError(err), data.Icon.Error)
							}
						},
					},
					dec.PushButton{ // Кнопка удалить.
						Text: text.T.ButtonDelete,
						OnClicked: func() {
							log.Println(l.Info, l.LogDelete) // Лог.
							if err := wf.delete(db); err != nil {
								err = errors.Wrap(err, e.Err.ErrorDeleteRow) // Обработка ошибок.
								log.Println(l.Error, err)                    // Лог
								walk.MsgBox(wf, text.T.MsgBoxError, MsgError(err), data.Icon.Error)
							}
						},
					},
				},
			},
			dec.Composite{
				Layout:  dec.HBox{MarginsZero: true},
				Visible: !isChange, // Видимость.
				Children: []dec.Widget{
					dec.PushButton{ // Кнопка Ок.
						Text: text.T.ButtonOK,
						OnClicked: func() {
							log.Println(l.Info, l.LogOk) // Лог.
							if wf.modelTable.RowCount() > 0 && wf.tv.CurrentIndex() != -1 {
								index := wf.tv.CurrentIndex()
								idTitle.Id = wf.modelTable.items[index].Id
								idType := wf.modelTable.items[index].Type
								sType := wf.modelTable.mapIdToTitle[idType]
								idTitle.Title = sType + " " + wf.modelTable.items[index].Title
								wf.Accept()
							} else {
								walk.MsgBox(wf, text.T.MsgBoxInfo, text.T.MsgChooseRow, data.Icon.Info)
							}
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
	}.Create(owner)); err != nil {
		err = errors.Wrap(err, e.Err.ErrorCreateWindow) // Обработка ошибок создания окна.
		return 0, err
	}
	log.Printf(l.CreateWindow, l.Entities) // Лог.

	log.Printf(l.RunWindow, l.Entities) // Лог.
	return wf.Run(), nil                // Запуск окна.
}

// Функция, для добавления строки в таблицу.
func (wf *windowsFormEntities) add(db *sql.DB) error {
	entity := NewEntity()
	cmd, err := EntityRunDialog(wf, db, &entity)
	log.Printf(l.EndWindow, l.Entity, cmd) // Лог.
	if err != nil {
		return errors.Wrapf(err, l.In.InEntityRunDialog, entity)
	}
	if cmd != walk.DlgCmdOK {
		return nil
	}

	QwStr := qwery.InsertEntity(entity.Title, entity.Type, entity.Enumerable, int8(entity.Marking), entity.Specification, entity.Note)
	if err = db.Ping(); err != nil { // Пинг БД.
		return errors.Wrap(err, e.Err.ErrorPingDB)
	}
	result, err := db.Exec(QwStr) // Запрос к БД.
	if err != nil {
		return errors.Wrapf(err, e.Err.ErrorAddDB, QwStr)
	}

	trackLatest := wf.tv.ItemVisible(len(wf.modelTable.items) - 1) //&& len(wf.tv.SelectedIndexes()) <= 1
	wf.modelTable.items = append(wf.modelTable.items, &entity)
	index := len(wf.modelTable.items) - 1
	wf.modelTable.PublishRowsReset()
	wf.modelTable.PublishRowsInserted(index, index)
	if trackLatest {
		wf.tv.EnsureItemVisible(len(wf.modelTable.items) - 1)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println(l.Error, e.Err.ErrorInsertIndexLog)                                 // Лог.
		walk.MsgBox(wf, text.T.MsgBoxError, e.Err.ErrorInsertIndex, data.Icon.Critical) // TO-DO
		return nil
	}
	wf.modelTable.items[index].Id = id
	for _, val := range entity.Children {
		QwStrChild := qwery.InsertEntityRec(id, val.Id, val.Count)
		if _, err := db.Exec(QwStrChild); err != nil { // Запрос к БД.
			err = errors.Wrap(err, e.Err.ErrorAddDB+QwStrChild) // Обработка ошибок.
			log.Println(l.Error, err)                           // Лог.
			walk.MsgBox(wf, text.T.MsgBoxError, MsgError(err), data.Icon.Error)
		}
	}
	return nil
}

// Функция, для изменения строки в таблице.
func (wf *windowsFormEntities) change(db *sql.DB) error {
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 {
		walk.MsgBox(wf, text.T.MsgBoxInfo, text.T.MsgChooseRow, data.Icon.Info)
		return nil
	}
	var err error
	index := wf.tv.CurrentIndex()
	entity := wf.modelTable.items[index]
	_, children, err := SelectEntityRecChild(db, &entity.Id)
	if err != nil {
		return errors.Wrap(err, e.Err.ErrorSubquery)
	}
	entity.Children = children
	cmd, err := EntityRunDialog(wf, db, entity)
	log.Printf(l.EndWindow, l.Entity, cmd) // Лог.
	if err != nil {
		return errors.Wrapf(err, l.In.InEntityRunDialog, *entity)
	}
	if cmd != walk.DlgCmdOK {
		return nil
	}

	QwStr := qwery.UpdateEntity(entity.Id, entity.Title, entity.Type, entity.Enumerable,
		int8(entity.Marking), entity.Specification, entity.Note)
	if err = db.Ping(); err != nil { // Пинг БД.
		return errors.Wrap(err, e.Err.ErrorPingDB)
	}
	_, err = db.Exec(QwStr) // Запрос к БД.
	if err != nil {
		return errors.Wrapf(err, e.Err.ErrorChangeDB, QwStr)
	}

	wf.modelTable.items[index] = entity // Обновление таблицы.
	wf.modelTable.PublishRowsReset()
	wf.modelTable.PublishRowChanged(index)
	return nil
}

// Функция, для удаления строки из таблицы.
func (wf *windowsFormEntities) delete(db *sql.DB) error {
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 {
		walk.MsgBox(wf, text.T.MsgBoxInfo, text.T.MsgChooseRow, data.Icon.Info)
		return nil
	}
	index := wf.tv.CurrentIndex()
	id := wf.modelTable.items[index].Id

	QwStr := qwery.DeleteEntity(id)
	if err := db.Ping(); err != nil { // Пинг БД.
		return errors.Wrap(err, e.Err.ErrorPingDB)
	}
	_, err := db.Exec(QwStr) // Запрос к БД.
	if err != nil {
		return errors.Wrapf(err, e.Err.ErrorDeleteDB, QwStr)
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
