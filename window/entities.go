package window

import (
	"accounting/data"
	"database/sql"

	// "fmt"
	"log"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

// Выборка из таблицы Entity всех ее полей удовлетворяющих условию,
// где в значения поля Title входит title,
// значение поля Type равно entityType (при 0, разрешен любой тип),
// isChange определяет, разрешено ли выбирать строчки, где тип сущности это заказ.
// Информация о дочерних сущностях не выбирается.
func SelectEntities(db *sql.DB, title string, entityType int64, isChange bool) ([]*Entity, error) {
	arr := make([]*Entity, 0)
	if err := (func() error {
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		QwStr := data.SelectEntity(title, entityType, isChange)
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrapf(err, data.S.ErrorQueryDB, QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			row := NewEntity()
			err := rows.Scan(&row.Id, &row.Title, &row.Type, &row.Specification, &row.Marking, &row.Note)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorDecryptRow)
			}
			arr = append(arr, &row)
		}
		return nil
	}()); err != nil {
		return arr, errors.Wrapf(err, data.Log.InSelectEntities, title, entityType)
	}
	return arr, nil
}

// Сруктура, содержащая модель таблицы.
type modelEntitiesComponent struct {
	walk.TableModelBase
	items        []*Entity
	mapIdToTitle map[int64]string
}

// Структура, содержащая описание и переменные окна.
type windowsFormEntities struct {
	*walk.Dialog
	modelType  []*IdTitle
	modelTable *modelEntitiesComponent
	tv         *walk.TableView
}

// Инициализация модели окна.
func newWindowsFormEntities(db *sql.DB, isChange bool) (*windowsFormEntities, error) {
	var err error
	wf := new(windowsFormEntities)
	wf.modelTable, err = newModelEntitiesComponent(db, isChange)
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorTableInit)
		return nil, err
	}
	wf.modelType, wf.modelTable.mapIdToTitle, err = SelectIdTitle(db, "EntityType")
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorTypeInit)
		return nil, err
	}
	wf.modelType = append([]*IdTitle{new(IdTitle)}, wf.modelType...)
	return wf, nil
}

// Инициализация модели таблицы.
func newModelEntitiesComponent(db *sql.DB, isChange bool) (*modelEntitiesComponent, error) {
	var err error
	m := new(modelEntitiesComponent)
	m.items, err = SelectEntities(db, "", 0, isChange)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *modelEntitiesComponent) RowCount() int {
	return len(m.items)
}

func (m *modelEntitiesComponent) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return m.mapIdToTitle[item.Type]
	case 1:
		return item.Title
	case 2:
		return item.Specification
	case 3:
		return MapMarkingToTitle[item.Marking]
	case 4:
		return item.Note
	}
	log.Println(data.Log.Panic, data.S.ErrorUnexpectedColumn)
	panic(data.S.ErrorUnexpectedColumn)
}

// Описание и запуск диалогового окна.
func EntitiesRunDialog(owner walk.Form, db *sql.DB, isChange bool, idTitle *IdTitle) (int, error) {
	log.Printf(data.Log.BeginWindow, data.Log.Entities)
	var err error
	var databind *walk.DataBinder
	search := new(IdTitle)
	wf, err := newWindowsFormEntities(db, isChange)
	if err != nil {
		return 0, errors.Wrap(err, data.S.ErrorInit)
	}
	log.Printf(data.Log.InitWindow, data.Log.Entities)
	if err = (dec.Dialog{
		AssignTo: &wf.Dialog,
		Title:    data.S.HeadingEntities,
		DataBinder: dec.DataBinder{
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
					dec.Label{
						Text: "Название:",
					},
					dec.LineEdit{
						Text: dec.Bind("Title"),
					},
					dec.HSpacer{Size: 20},

					dec.Label{
						Text: "Тип:",
					},
					dec.ComboBox{
						Value:         dec.Bind("Id", dec.SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Title",
						MinSize:       dec.Size{120, 0},
						Model:         wf.modelType,
					},
					dec.HSpacer{Size: 20},

					dec.PushButton{
						Text: data.S.ButtonSearch,
						OnClicked: func() {
							log.Println(data.Log.Info, data.Log.LogSearch)
							if err := databind.Submit(); err != nil {
								err = errors.Wrap(err, data.S.ErrorSubmit)
								log.Println(data.Log.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
								return
							}
							lastLen := wf.modelTable.RowCount()
							if items, err := SelectEntities(db, search.Title, search.Id, isChange); err != nil {
								err = errors.Wrap(err, data.S.ErrorSubquery)
								log.Println(data.Log.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
								return
							} else {
								wf.modelTable.items = items
							}
							nowLen := wf.modelTable.RowCount()
							wf.modelTable.PublishRowsReset()
							wf.modelTable.PublishRowsRemoved(0, lastLen)
							wf.modelTable.PublishRowsInserted(0, nowLen)
							wf.modelTable.PublishRowsChanged(0, nowLen)
						},
					},
				},
			},
			dec.TableView{
				AssignTo: &wf.tv,
				Columns: []dec.TableViewColumn{
					{Title: "Тип"},
					{Title: "Название"},
					{Title: "Спецификация"},
					{Title: "Маркировка"},
					{Title: "Примичание"},
				},
				MinSize: dec.Size{0, 200},
				Model:   wf.modelTable,
			},
			dec.Composite{
				Layout:  dec.HBox{MarginsZero: true},
				Visible: isChange,
				Children: []dec.Widget{
					dec.PushButton{
						Text: data.S.ButtonAdd,
						OnClicked: func() {
							log.Println(data.Log.Info, data.Log.LogAdd)
							if err := wf.add(db); err != nil {
								err = errors.Wrap(err, data.S.ErrorAddRow)
								log.Println(data.Log.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
							}
						},
					},
					dec.PushButton{
						Text: data.S.ButtonChange,
						OnClicked: func() {
							log.Println(data.Log.Info, data.Log.LogChange)
							if err := wf.change(db); err != nil {
								err = errors.Wrap(err, data.S.ErrorChangeRow)
								log.Println(data.Log.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
							}
						},
					},
					dec.PushButton{
						Text: data.S.ButtonDelete,
						OnClicked: func() {
							log.Println(data.Log.Info, data.Log.LogDelete)
							if err := wf.delete(db); err != nil {
								err = errors.Wrap(err, data.S.ErrorDeleteRow)
								log.Println(data.Log.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
							}
						},
					},
				},
			},
			dec.Composite{
				Layout:  dec.HBox{MarginsZero: true},
				Visible: !isChange,
				Children: []dec.Widget{
					dec.PushButton{
						Text: data.S.ButtonOK,
						OnClicked: func() {
							log.Println(data.Log.Info, data.Log.LogOk)
							if wf.modelTable.RowCount() > 0 && wf.tv.CurrentIndex() != -1 {
								index := wf.tv.CurrentIndex()
								idTitle.Id = wf.modelTable.items[index].Id
								idType := wf.modelTable.items[index].Type
								sType := wf.modelTable.mapIdToTitle[idType]
								idTitle.Title = sType + " " + wf.modelTable.items[index].Title
								wf.Accept()
							} else {
								walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChooseRow, data.Icon.Info)
							}
						},
					},
					dec.PushButton{
						Text:      data.S.ButtonCansel,
						OnClicked: func() { wf.Cancel() },
					},
				},
			},
		},
	}.Create(owner)); err != nil {
		err = errors.Wrap(err, data.S.ErrorCreateWindow)
		return 0, err
	}
	log.Printf(data.Log.CreateWindow, data.Log.Entities)

	log.Printf(data.Log.RunWindow, data.Log.Entities)
	return wf.Run(), nil
}

// Функция, для добавления строки в таблицу.
func (wf windowsFormEntities) add(db *sql.DB) error {
	entity := NewEntity()
	cmd, err := EntityRunDialog(wf, db, &entity)
	log.Printf(data.Log.EndWindow, data.Log.Entity, cmd)
	if err != nil {
		return errors.Wrapf(err, data.Log.InEntityRunDialog, entity)
	}
	if cmd != walk.DlgCmdOK {
		return nil
	}
	QwStr := data.InsertEntity(entity.Title, entity.Specification, entity.Note, int8(entity.Marking), entity.Type)
	if err = db.Ping(); err != nil {
		return errors.Wrap(err, data.S.ErrorPingDB)
	}

	result, err := db.Exec(QwStr)
	if err != nil {
		return errors.Wrapf(err, data.S.ErrorAddDB, QwStr)
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
		log.Println(data.Log.Error, data.S.ErrorInsertIndexLog)
		walk.MsgBox(wf, data.S.MsgBoxError, data.S.ErrorInsertIndex, data.Icon.Critical)
		return nil
	}
	wf.modelTable.items[index].Id = id
	for _, val := range entity.Children {
		QwStrChild := data.InsertEntityRec(id, val.Id, val.Count)
		if _, err := db.Exec(QwStrChild); err != nil {
			err = errors.Wrap(err, data.S.ErrorAddDB+QwStrChild)
			log.Println(data.Log.Error, err)
			walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
		}
	}
	return nil
}

// Функция, для изменения строки в таблице.
func (wf windowsFormEntities) change(db *sql.DB) error {
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 {
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChooseRow, data.Icon.Info)
		return nil
	}
	var err error
	index := wf.tv.CurrentIndex()
	entity := wf.modelTable.items[index]
	children, err := SelectEntityRecChild(db, entity.Id)
	if err != nil {
		return errors.Wrap(err, data.S.ErrorSubquery)
	}
	entity.Children = children
	cmd, err := EntityRunDialog(wf, db, entity)
	log.Printf(data.Log.EndWindow, data.Log.Entity, cmd)

	if err != nil {
		return errors.Wrapf(err, data.Log.InEntityRunDialog, *entity)
	}
	if cmd != walk.DlgCmdOK {
		return nil
	}
	QwStr := data.UpdateEntity(entity.Title, entity.Specification, entity.Note, int8(entity.Marking), entity.Type, entity.Id)
	if err = db.Ping(); err != nil {
		return errors.Wrap(err, data.S.ErrorPingDB)
	}
	_, err = db.Exec(QwStr)
	if err != nil {
		return errors.Wrapf(err, data.S.ErrorChangeDB, QwStr)
	}
	wf.modelTable.items[index] = entity
	wf.modelTable.PublishRowsReset()
	wf.modelTable.PublishRowChanged(index)
	return nil
}

// Функция, для удаления строки из таблицы.
func (wf windowsFormEntities) delete(db *sql.DB) error {
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 {
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChooseRow, data.Icon.Info)
		return nil
	}
	index := wf.tv.CurrentIndex()
	id := wf.modelTable.items[index].Id
	QwStr := data.DeleteEntity(id)
	if err := db.Ping(); err != nil {
		return errors.Wrap(err, data.S.ErrorPingDB)
	}
	_, err := db.Exec(QwStr)
	if err != nil {
		return errors.Wrapf(err, data.S.ErrorDeleteDB, QwStr)
	}

	trackLatest := wf.tv.ItemVisible(len(wf.modelTable.items) - 1) //&& len(wf.tv.SelectedIndexes()) <= 1
	wf.modelTable.items = wf.modelTable.items[:index+copy(wf.modelTable.items[index:], wf.modelTable.items[index+1:])]
	wf.modelTable.PublishRowsReset()
	wf.modelTable.PublishRowsRemoved(index, index)
	if trackLatest {
		wf.tv.EnsureItemVisible(len(wf.modelTable.items) - 1)
	}
	// if l := len(wf.modelTable.items); l <= index {
	// 	index = l - 1
	// }
	// if index >= 0 {
	// 	wf.tv.SetCurrentIndex(index)
	// }
	wf.modelTable.PublishRowsChanged(index, wf.modelTable.RowCount()-1)
	return nil
}
