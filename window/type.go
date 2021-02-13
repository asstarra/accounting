package window

import (
	"accounting/data"
	"database/sql"
	"log"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

// Выборка идентификатора и названия из таблиц типа Type.
func SelectIdTitle(db *sql.DB, tableName string) ([]*IdTitle, map[int64]string, error) {
	arr := make([]*IdTitle, 0)
	m := make(map[int64]string)
	if err := (func() error {
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		QwStr := data.SelectType(tableName)
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorQueryDB+QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			row := IdTitle{}
			err := rows.Scan(&row.Id, &row.Title)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorDecryptRow)
			}
			arr = append(arr, &row)
			m[row.Id] = row.Title
		}
		return nil
	}()); err != nil {
		err = errors.Wrapf(err, data.S.InSelectIdTitle, tableName)
	}
	return arr, m, nil
}

// Сруктура, содержащая модель таблицы.
type modelTypeComponent struct {
	walk.TableModelBase
	items []*IdTitle
}

// Структура, содержащая описание и переменные окна.
type windowsFormType struct {
	*walk.Dialog
	modelTable *modelTypeComponent
	tv         *walk.TableView
	textW      *walk.LineEdit
}

// Инициализация модели окна.
func newWindowsFormType(db *sql.DB, tableName string) (*windowsFormType, error) {
	var err error
	wf := new(windowsFormType)
	wf.modelTable = new(modelTypeComponent)
	wf.modelTable.items, _, err = SelectIdTitle(db, tableName)
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorTableInit)
		return nil, err
	}
	return wf, nil
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
	log.Println(data.S.Panic, data.S.ErrorUnexpectedColumn)
	panic(data.S.ErrorUnexpectedColumn)
}

// Описание и запуск диалогового окна.
func TypeRunDialog(owner walk.Form, db *sql.DB, tableName string) (int, error) {
	log.Printf(data.S.BeginWindow, data.S.Type)
	wf, err := newWindowsFormType(db, tableName)
	if err != nil {
		return 0, errors.Wrap(err, data.S.ErrorInit)
	}
	log.Printf(data.S.InitWindow, data.S.Type)
	if err := (dec.Dialog{
		AssignTo: &wf.Dialog,
		Title:    data.S.HeadingType,
		Layout:   dec.VBox{MarginsZero: true},
		Children: []dec.Widget{
			dec.HSplitter{
				Children: []dec.Widget{
					dec.Composite{
						Layout: dec.VBox{},
						Children: []dec.Widget{
							dec.TableView{
								AssignTo: &wf.tv,
								Columns: []dec.TableViewColumn{
									{Title: "Название"},
								},
								MinSize: dec.Size{120, 0},
								Model:   wf.modelTable,
							},
						},
					},
					dec.Composite{
						Layout: dec.VBox{},
						Children: []dec.Widget{
							dec.LineEdit{
								AssignTo:  &wf.textW,
								MaxLength: 255,
								MinSize:   dec.Size{120, 0},
							},
							dec.PushButton{
								Text: data.S.ButtonAdd,
								OnClicked: func() {
									log.Println(data.S.Info, data.S.LogAdd)
									if err := wf.add(db, tableName); err != nil {
										err = errors.Wrap(err, data.S.ErrorAddRow)
										log.Println(data.S.Error, err)
										walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
									}
								},
							},
							dec.PushButton{
								Text: data.S.ButtonChange,
								OnClicked: func() {
									log.Println(data.S.Info, data.S.LogChange)
									if err := wf.change(db, tableName); err != nil {
										err = errors.Wrap(err, data.S.ErrorChangeRow)
										log.Println(data.S.Error, err)
										walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
									}
								},
							},
							dec.PushButton{
								Text: data.S.ButtonDelete,
								OnClicked: func() {
									log.Println(data.S.Info, data.S.LogDelete)
									if err := wf.delete(db, tableName); err != nil {
										err = errors.Wrap(err, data.S.ErrorDeleteRow)
										log.Println(data.S.Error, err)
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
		err = errors.Wrap(err, data.S.ErrorCreateWindow)
		return 0, err
	}
	log.Printf(data.S.CreateWindow, data.S.Type)

	log.Printf(data.S.RunWindow, data.S.Type)
	return wf.Run(), nil
}

// Функция, для добавления строки в таблицу.
func (wf windowsFormType) add(db *sql.DB, tableName string) error {
	if wf.textW.Text() == "" {
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgEmptyTitle, data.Icon.Info)
		return nil
	}
	var row IdTitle
	row.Title = wf.textW.Text()
	QwStr := data.InsertType(tableName, row.Title)
	if err := db.Ping(); err != nil {
		return errors.Wrap(err, data.S.ErrorPingDB)
	}
	result, err := db.Exec(QwStr)
	if err != nil {
		return errors.Wrap(err, data.S.ErrorAddDB+QwStr)
	}
	if row.Id, err = result.LastInsertId(); err != nil {
		log.Println(data.S.Error, errors.Wrap(err, data.S.ErrorInsertIndexLog))
		walk.MsgBox(wf, data.S.MsgBoxError, data.S.ErrorInsertIndex, data.Icon.Critical)
		row.Id = 0
	}
	wf.textW.SetText("")
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
func (wf windowsFormType) change(db *sql.DB, tableName string) error {
	if wf.textW.Text() == "" {
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgEmptyTitle, data.Icon.Info)
		return nil
	}
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 {
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChooseRow, data.Icon.Info)
		return nil
	}
	index := wf.tv.CurrentIndex()
	QwStr := data.UpdateType(tableName, wf.textW.Text(), wf.modelTable.items[index].Id)
	if err := db.Ping(); err != nil {
		return errors.Wrap(err, data.S.ErrorPingDB)
	}
	if _, err := db.Exec(QwStr); err != nil {
		return errors.Wrap(err, data.S.ErrorChangeDB+QwStr)
	}
	wf.modelTable.items[index].Title = wf.textW.Text()
	wf.textW.SetText("")
	wf.modelTable.PublishRowsReset()
	wf.modelTable.PublishRowChanged(index)
	return nil
}

// Функция, для удаления строки из таблицы.
func (wf windowsFormType) delete(db *sql.DB, tableName string) error {
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 {
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChooseRow, data.Icon.Info)
		return nil
	}
	index := wf.tv.CurrentIndex()
	id := wf.modelTable.items[index].Id
	QwStr := data.DeleteType(tableName, id)
	if err := db.Ping(); err != nil {
		return errors.Wrap(err, data.S.ErrorPingDB)
	}
	if _, err := db.Exec(QwStr); err != nil {
		return errors.Wrap(err, data.S.ErrorDeleteDB+QwStr)
	}
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
