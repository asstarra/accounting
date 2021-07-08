package window

import (
	e "accounting/data/errors"
	l "accounting/data/log"
	"accounting/data/qwery"
	"accounting/data/text"
	. "accounting/window/data"
	"database/sql"

	// "fmt"
	"log"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

// Выборка идентификатора и названия из таблиц типа Type.
func SelectId16Title(db *sql.DB, tableName string, id *int16,
	title *string) ([]*Id16Title, map[int16]string, error) {
	arr := make([]*Id16Title, 0, 20)
	m := make(map[int16]string, 20)
	if err := (func() error {
		QwStr := qwery.SelectType16(tableName, id, title)
		if err := db.Ping(); err != nil { // Пинг БД.
			return errors.Wrap(err, e.Err.ErrorPingDB)
		}
		rows, err := db.Query(QwStr) // Запрос к БД.
		if err != nil {
			return errors.Wrapf(err, e.Err.ErrorQueryDB, QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			var row Id16Title
			err := rows.Scan(&row.Id, &row.Title)
			if err != nil {
				return errors.Wrap(err, e.Err.ErrorDecryptRow)
			}
			arr = append(arr, &row)
			m[row.Id] = row.Title
		}
		return nil
	}()); err != nil {
		err = errors.Wrapf(err, l.In.InSelectId16Title, tableName, qwery.ToStr(id), qwery.ToStr(title))
	}
	return arr, m, nil
}

// Сруктура, содержащая модель таблицы.
type modelTypeComponent struct {
	walk.TableModelBase
	items []*Id16Title // Содержит названия.
	item  Id16Title    // Элемент для сравнения при изменении.
}

func (m *modelTypeComponent) RowCount() int {
	return len(m.items)
}

func (m *modelTypeComponent) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return item.Title
	}
	log.Println(l.Panic, e.Err.ErrorUnexpectedColumn) // Лог.
	panic(e.Err.ErrorUnexpectedColumn)
}

func (m *modelTypeComponent) Equal(row, col int) bool {
	item := m.items[row]
	switch col {
	case 0:
		return m.item.Title == item.Title
	}
	log.Println(l.Panic, e.Err.ErrorUnexpectedColumn) // Лог.
	panic(e.Err.ErrorUnexpectedColumn)
}

// Структура, содержащая описание и переменные окна.
type windowsFormType struct {
	*walk.Dialog
	modelTable *modelTypeComponent // Модель таблицы, в которой содержатся названия.
	tv         *walk.TableView     // Виджет таблицы, в которой содержатся названия.
	textW      *walk.LineEdit      // Виджит текстовой строки.
}

// Инициализация модели окна.
func newWindowsFormType(db *sql.DB, tableName string) (*windowsFormType, error) {
	var err error
	wf := new(windowsFormType)
	wf.modelTable = new(modelTypeComponent)
	wf.modelTable.items, _, err = SelectId16Title(db, tableName, nil, nil)
	if err != nil {
		err = errors.Wrap(err, e.Err.ErrorTableInit)
		return nil, err
	}
	return wf, nil
}

// Описание и запуск диалогового окна для задания "типов".
func TypeRunDialog(owner walk.Form, db *sql.DB, tableName string) (int, error) {
	log.Printf(l.BeginWindow, l.Type)            // Лог.
	wf, err := newWindowsFormType(db, tableName) // Инициализация.
	if err != nil {
		return 0, errors.Wrap(err, e.Err.ErrorInit)
	}
	log.Printf(l.InitWindow, l.Type) // Лог.
	if err := (dec.Dialog{           // Описание окна.
		AssignTo: &wf.Dialog,                  // Привязка окна.
		Title:    text.HeadingType[tableName], // Название.
		Layout:   dec.VBox{MarginsZero: true},
		Children: []dec.Widget{ // Элементы окна.
			dec.HSplitter{
				Children: []dec.Widget{
					dec.Composite{ // Левая половина.
						Layout: dec.VBox{},
						Children: []dec.Widget{
							dec.TableView{ // Таблица.
								AssignTo: &wf.tv, // Привязка к виджету.
								Columns: []dec.TableViewColumn{
									{Title: text.TitleType[tableName]},
								},
								MinSize: dec.Size{120, 0},
								Model:   wf.modelTable, // Привязка к модели.
							},
						},
					}, // Правая половина.
					dec.Composite{
						Layout: dec.VBox{},
						Children: []dec.Widget{
							dec.LineEdit{ // Текстовая строка.
								AssignTo:  &wf.textW, // Привязка к виджету.
								MaxLength: 255,       // Ограничение на количество букв.
								MinSize:   dec.Size{120, 0},
							},
							dec.PushButton{ // Кнопка добавить.
								Text: text.T.ButtonAdd,
								OnClicked: func() {
									log.Println(l.Info, l.LogAdd)                 // Лог.
									if err := wf.add(db, tableName); err != nil { // Обработка ошибки.
										MsgBoxAdd(wf, err)
									}
								},
							},
							dec.PushButton{ // Кнопка изменить.
								Text: text.T.ButtonChange,
								OnClicked: func() {
									log.Println(l.Info, l.LogChange)                 // Лог.
									if err := wf.change(db, tableName); err != nil { // Обработка ошибки.
										MsgBoxChange(wf, err)
									}
								},
							},
							dec.PushButton{ // Кнопка удалить.
								Text: text.T.ButtonDelete,
								OnClicked: func() {
									log.Println(l.Info, l.LogDelete)                 // Лог.
									if err := wf.delete(db, tableName); err != nil { // Обработка ошибки.
										MsgBoxDelete(wf, err)
									}
								},
							},
						},
					},
				},
			},
		},
	}.Create(owner)); err != nil {
		err = errors.Wrap(err, e.Err.ErrorCreateWindow) // Обработка ошибки создания окна.
		return 0, err
	}
	log.Printf(l.CreateWindow, l.Type) // Лог.

	log.Printf(l.RunWindow, l.Type) // Лог.
	return wf.Run(), nil            // Запуск окна.
}

// Функция, для добавления строки в таблицу.
func (wf *windowsFormType) add(db *sql.DB, tableName string) error {
	// Обработка ограничинения на пустую строку и повторяемость значений.
	wf.modelTable.item = Id16Title{Id: 0, Title: wf.textW.Text()}
	if IsStringEmpty(wf, wf.modelTable.item.Title) || IsRepeat(wf, wf.modelTable, []int{0}) {
		return nil
	}
	QwStr := qwery.InsertType16(tableName, wf.modelTable.item.Title)
	if err := db.Ping(); err != nil { // Пинг БД.
		return errors.Wrap(err, e.Err.ErrorPingDB)
	}
	result, err := db.Exec(QwStr) // Запрос к БД.
	if err != nil {
		return errors.Wrapf(err, e.Err.ErrorAddDB, QwStr)
	}
	if id, err := result.LastInsertId(); err != nil { // Выбор Id.
		MsgBoxNotInsertedId(wf, err)
	} else {
		wf.modelTable.item.Id = int16(id)
	}
	wf.textW.SetText("") // Обнуление текстовой строки.
	// Обновление таблицы.
	trackLatest := wf.tv.ItemVisible(len(wf.modelTable.items) - 1)
	wf.modelTable.items = append(wf.modelTable.items, &wf.modelTable.item)
	index := len(wf.modelTable.items) - 1
	wf.modelTable.PublishRowsReset()
	wf.modelTable.PublishRowsInserted(index, index)
	if trackLatest {
		wf.tv.EnsureItemVisible(len(wf.modelTable.items) - 1)
	}
	return nil
}

// Функция, для изменения строки в таблице.
func (wf *windowsFormType) change(db *sql.DB, tableName string) error {
	// Обработка ограничинения на пустую строку и повторяемость значений.
	wf.modelTable.item = Id16Title{Id: 0, Title: wf.textW.Text()}
	if IsStringEmpty(wf, wf.modelTable.item.Title) || IsRepeat(wf, wf.modelTable, []int{0}) {
		return nil
	}
	// Проверка на выделение изменяемой строчки и возможность ее изменения.
	if IsCorrectIndex(wf, wf.modelTable, wf.tv) || wf.isConstraint(tableName) {
		return nil
	}
	index := wf.tv.CurrentIndex()
	QwStr := qwery.UpdateType16(tableName, wf.modelTable.items[index].Id, wf.textW.Text())
	if err := db.Ping(); err != nil { // Пинг БД.
		return errors.Wrap(err, e.Err.ErrorPingDB)
	}
	if _, err := db.Exec(QwStr); err != nil { // Запрос к БД.
		return errors.Wrapf(err, e.Err.ErrorChangeDB, QwStr)
	}
	// Обновление таблицы и текстовой строки.
	wf.textW.SetText("")
	wf.modelTable.items[index].Title = wf.modelTable.item.Title
	wf.modelTable.PublishRowsReset()
	wf.modelTable.PublishRowChanged(index)
	return nil
}

// Функция, для удаления строки из таблицы.
func (wf *windowsFormType) delete(db *sql.DB, tableName string) error {
	// Проверка на выделение изменяемой строчки и возможность ее изменения.
	if IsCorrectIndex(wf, wf.modelTable, wf.tv) || wf.isConstraint(tableName) {
		return nil
	}
	index := wf.tv.CurrentIndex()
	QwStr := qwery.DeleteType16(tableName, wf.modelTable.items[index].Id)
	if err := db.Ping(); err != nil { // Пинг БД.
		return errors.Wrap(err, e.Err.ErrorPingDB)
	}
	if _, err := db.Exec(QwStr); err != nil { // Запрос к БД.
		return errors.Wrapf(err, e.Err.ErrorDeleteDB, QwStr)
	}
	// Обновление таблицы.
	trackLatest := wf.tv.ItemVisible(len(wf.modelTable.items) - 1)
	wf.modelTable.items = wf.modelTable.items[:index+copy(wf.modelTable.items[index:], wf.modelTable.items[index+1:])]
	wf.modelTable.PublishRowsReset()
	wf.modelTable.PublishRowsRemoved(index, index)
	if trackLatest {
		wf.tv.EnsureItemVisible(len(wf.modelTable.items) - 1)
	}
	wf.modelTable.PublishRowsChanged(index, wf.modelTable.RowCount()-1)
	return nil
}

// Функция, которая проверяет ограничения при изменении и удалении,
// которые зависят от имени таблицы. Строчка должна быть выделена.
func (wf *windowsFormType) isConstraint(tableName string) bool {
	index := wf.tv.CurrentIndex()
	id := wf.modelTable.items[index].Id
	switch tableName { // TO-DO Для других типов.
	case "EntityType":
		if id >= 1 || id <= 5 {
			MsgBoxNotChange(wf)
			return true
		}
	}
	return false
}
