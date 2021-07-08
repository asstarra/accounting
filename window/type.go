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

type Id16Title struct {
	Id    int16  // Id.
	Title string // Название.
}

func (a Id16Title) String() string {
	return fmt.Sprintf("{Id = %d, Title = '%s'}", a.Id, a.Title)
}

// Выборка идентификатора и названия из таблиц типа Type.
func SelectId16Title(db *sql.DB, tableName string, id *int16,
	title *string) ([]*Id16Title, map[int16]string, error) {
	arr := make([]*Id16Title, 0, 20)
	m := make(map[int16]string, 20)
	if err := (func() error {
		QwStr := qwery.SelectType16(tableName, id, title)
		if err := db.Ping(); err != nil { // Пинг БД.
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		rows, err := db.Query(QwStr) // Запрос к БД.
		if err != nil {
			return errors.Wrapf(err, data.S.ErrorQueryDB, QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			var row Id16Title
			err := rows.Scan(&row.Id, &row.Title)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorDecryptRow)
			}
			arr = append(arr, &row)
			m[row.Id] = row.Title
		}
		return nil
	}()); err != nil {
		err = errors.Wrapf(err, data.Log.InSelectIdTitle, tableName) //GO-TO
	}
	return arr, m, nil
}

// Сруктура, содержащая модель таблицы.
type modelTypeComponent struct {
	walk.TableModelBase
	items []*Id16Title // Содержит названия.
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
	log.Println(data.Log.Panic, data.S.ErrorUnexpectedColumn)
	panic(data.S.ErrorUnexpectedColumn)
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
		err = errors.Wrap(err, data.S.ErrorTableInit)
		return nil, err
	}
	return wf, nil
}

// Описание и запуск диалогового окна для задания "типов".
func TypeRunDialog(owner walk.Form, db *sql.DB, tableName string) (int, error) {
	log.Printf(data.Log.BeginWindow, data.Log.Type) // Лог.
	wf, err := newWindowsFormType(db, tableName)    // Инициализация.
	if err != nil {
		return 0, errors.Wrap(err, data.S.ErrorInit)
	}
	log.Printf(data.Log.InitWindow, data.Log.Type) // Лог.
	if err := (dec.Dialog{                         // Описание окна.
		AssignTo: &wf.Dialog,         // Привязка окна.
		Title:    data.S.HeadingType, // Название.
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
									{Title: "Название"}, //GO-TO
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
								Text: data.S.ButtonAdd,
								OnClicked: func() {
									log.Println(data.Log.Info, data.Log.LogAdd) // Лог.
									if err := wf.add(db, tableName); err != nil {
										err = errors.Wrap(err, data.S.ErrorAddRow) // Обработка ошибки.
										log.Println(data.Log.Error, err)           // Лог.
										walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
									}
								},
							},
							dec.PushButton{ // Кнопка изменить.
								Text: data.S.ButtonChange,
								OnClicked: func() {
									log.Println(data.Log.Info, data.Log.LogChange) // Лог.
									if err := wf.change(db, tableName); err != nil {
										err = errors.Wrap(err, data.S.ErrorChangeRow) // Обработка ошибки.
										log.Println(data.Log.Error, err)              // Лог.
										walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
									}
								},
							},
							dec.PushButton{ // Кнопка удалить.
								Text: data.S.ButtonDelete,
								OnClicked: func() {
									log.Println(data.Log.Info, data.Log.LogDelete)
									if err := wf.delete(db, tableName); err != nil { // Лог.
										err = errors.Wrap(err, data.S.ErrorDeleteRow) // Обработка ошибки.
										log.Println(data.Log.Error, err)              // Лог.
										walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
									}
								},
							},
						},
					},
				},
			},
		},
	}.Create(owner)); err != nil {
		err = errors.Wrap(err, data.S.ErrorCreateWindow) // Обработка ошибки создания окна.
		return 0, err
	}
	log.Printf(data.Log.CreateWindow, data.Log.Type) // Лог.

	log.Printf(data.Log.RunWindow, data.Log.Type) // Лог.
	return wf.Run(), nil                          // Запуск окна.
}

// Функция, для добавления строки в таблицу.
func (wf *windowsFormType) add(db *sql.DB, tableName string) error {
	var row = Id16Title{Id: 0, Title: wf.textW.Text()}
	if row.Title == "" { // Обработка ограничинения на пустую строку. GO-TO вынести. проверить на повторяемость.
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgEmptyTitle, data.Icon.Info)
		return nil
	}
	QwStr := qwery.InsertType16(tableName, row.Title)
	if err := db.Ping(); err != nil { // Пинг БД.
		return errors.Wrap(err, data.S.ErrorPingDB)
	}
	result, err := db.Exec(QwStr) // Запрос к БД.
	if err != nil {
		return errors.Wrapf(err, data.S.ErrorAddDB, QwStr)
	}
	if id, err := result.LastInsertId(); err != nil { // Выбор Id.
		log.Println(data.Log.Error, errors.Wrap(err, data.S.ErrorInsertIndexLog))        // Лог.
		walk.MsgBox(wf, data.S.MsgBoxError, data.S.ErrorInsertIndex, data.Icon.Critical) // GO-TO? вынести.
	} else {
		row.Id = int16(id)
	}
	wf.textW.SetText("") // Обнуление текстовой строки.
	// Обновление таблицы.
	trackLatest := wf.tv.ItemVisible(len(wf.modelTable.items) - 1)
	wf.modelTable.items = append(wf.modelTable.items, &row)
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
	var title = wf.textW.Text()
	if title == "" { // Обработка ограничинения на пустую строку. GO-TO вынести. проверить на повторяемость.
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgEmptyTitle, data.Icon.Info)
		return nil
	}
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 { // Проверка на выделение изменяемой строчки.
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChooseRow, data.Icon.Info) // GO-TO вынести.
		return nil
	}
	if wf.isConstraint(tableName) {
		return nil
	}
	index := wf.tv.CurrentIndex()
	QwStr := qwery.UpdateType16(tableName, wf.modelTable.items[index].Id, wf.textW.Text())
	if err := db.Ping(); err != nil { // Пинг БД.
		return errors.Wrap(err, data.S.ErrorPingDB)
	}
	if _, err := db.Exec(QwStr); err != nil { // Запрос к БД.
		return errors.Wrapf(err, data.S.ErrorChangeDB, QwStr)
	}
	// Обновление таблицы и текстовой строки.
	wf.textW.SetText("")
	wf.modelTable.items[index].Title = title
	wf.modelTable.PublishRowsReset()
	wf.modelTable.PublishRowChanged(index)
	return nil
}

// Функция, для удаления строки из таблицы.
func (wf *windowsFormType) delete(db *sql.DB, tableName string) error {
	index := wf.tv.CurrentIndex()
	if wf.modelTable.RowCount() <= 0 || index == -1 { // Проверка на выделение изменяемой строчки.
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChooseRow, data.Icon.Info) // GO-TO вынести.
		return nil
	}
	if wf.isConstraint(tableName) {
		return nil
	}
	QwStr := qwery.DeleteType16(tableName, wf.modelTable.items[index].Id)
	if err := db.Ping(); err != nil { // Пинг БД.
		return errors.Wrap(err, data.S.ErrorPingDB)
	}
	if _, err := db.Exec(QwStr); err != nil { // Запрос к БД.
		return errors.Wrapf(err, data.S.ErrorDeleteDB, QwStr)
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

// Функция, которая проверяет ограничения при изменении и удалении,
// которые зависят от имени таблицы.
func (wf *windowsFormType) isConstraint(tableName string) bool {
	index := wf.tv.CurrentIndex()
	id := wf.modelTable.items[index].Id // GO-TO
	switch tableName {                  // GO-TO
	case "EntityType":
		if id >= 1 || id <= 5 {
			walk.MsgBox(wf, data.S.MsgBoxInfo, "Данную строчку нельзя изменить.", data.Icon.Info) // GO-TO вынести.
			return true
		}
	}
	return false
}
