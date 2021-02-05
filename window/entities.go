package window

import (
	"accounting/data"
	"database/sql"
	"log"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

func SelectEntities(db *sql.DB, title string, entityType int64) ([]*Entity, error) {
	arr := make([]*Entity, 0)
	if err := (func() error {
		if err := db.Ping(); err != nil {
			err = errors.Wrap(err, data.S.ErrorPingDB)
			return err
		}
		QwStr := data.SelectEntity(title, entityType)
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorQuery+QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			row := Entity{}
			err := rows.Scan(&row.Id, &row.Title, &row.Type, &row.Specification, &row.Marking, &row.Note)
			if err != nil {
				err = errors.Wrap(err, data.S.ErrorDecryptRow)
				return err
			}
			arr = append(arr, &row)
		}
		return nil
	}()); err != nil {
		return arr, errors.Wrapf(err, data.S.InSelectEntities, title, entityType)
	}
	return arr, nil
}

type modelEntitiesComponent struct {
	walk.TableModelBase
	items      []*Entity
	mapIdTitle map[int64]string
}

type windowsFormEntities struct {
	*walk.Dialog
	modelType  []*IdTitle
	modelTable *modelEntitiesComponent
	tv         *walk.TableView
}

func newWindowsFormEntities(db *sql.DB) (*windowsFormEntities, error) {
	var err error
	wf := new(windowsFormEntities)
	wf.modelTable, err = newModelEntitiesComponent(db)
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorTableInit)
		return nil, err
	}
	wf.modelType, wf.modelTable.mapIdTitle, err = SelectIdTitle(db, "EntityType")
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorTypeInit)
		return nil, err
	}
	wf.modelType = append([]*IdTitle{new(IdTitle)}, wf.modelType...)
	return wf, nil
}

func newModelEntitiesComponent(db *sql.DB) (*modelEntitiesComponent, error) {
	var err error
	m := new(modelEntitiesComponent)
	m.items, err = SelectEntities(db, "", 0)
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
		return m.mapIdTitle[item.Type]
	case 1:
		return item.Title
	case 2:
		return item.Specification
	case 3:
		return item.Marking //GO-TO
	case 4:
		return item.Note
	}
	log.Println(data.S.Panic, data.S.ErrorUnexpectedColumn)
	panic(data.S.ErrorUnexpectedColumn)
}

func EntitiesRunDialog(owner walk.Form, db *sql.DB, isChange bool, idTitle *IdTitle) (int, error) {
	log.Printf(data.S.BeginWindow, data.S.Entities)
	var err error
	var databind *walk.DataBinder
	search := new(IdTitle)
	wf, err := newWindowsFormEntities(db)
	if err != nil {
		return 0, errors.Wrap(err, data.S.ErrorInit)
	}
	log.Printf(data.S.InitWindow, data.S.Entities)
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
		MinSize: dec.Size{450, 0},
		Children: []dec.Widget{
			dec.Composite{
				Layout: dec.HBox{},
				Children: []dec.Widget{
					dec.Label{
						Text: "Название:",
					},
					dec.LineEdit{
						Text: dec.Bind("Title"),
					},

					dec.Label{
						Text: "Тип:",
					},
					dec.ComboBox{
						Value:         dec.Bind("Id", dec.SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Title",
						Model:         wf.modelType,
					},

					dec.PushButton{
						Text: data.S.ButtonSearch,
						OnClicked: func() {
							log.Println(data.S.Info, data.S.LogSearch)
							if err := databind.Submit(); err != nil {
								err = errors.Wrap(err, data.S.ErrorSubmit)
								log.Println(data.S.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, err.Error(), data.Icon.Error)
								return
							}
							lastLen := wf.modelTable.RowCount()
							if wf.modelTable.items, err = SelectEntities(db, search.Title, search.Id); err != nil {
								err = errors.Wrap(err, data.S.ErrorSubquery)
								log.Println(data.S.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, err.Error(), data.Icon.Error)
								return
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
				Layout:  dec.HBox{},
				Visible: isChange,
				Children: []dec.Widget{
					dec.PushButton{
						Text: data.S.ButtonAdd,
						OnClicked: func() {
							log.Println(data.S.Info, data.S.LogAdd)
							if err := wf.add(db); err != nil {
								err = errors.Wrap(err, data.S.ErrorAdd)
								log.Println(data.S.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, err.Error(), data.Icon.Error)
							}
						},
					},
					dec.PushButton{
						Text: data.S.ButtonChange,
						OnClicked: func() {
							log.Println(data.S.Info, data.S.LogChange)
							if err := wf.change(db); err != nil {
								err = errors.Wrap(err, data.S.ErrorChange)
								log.Println(data.S.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, err.Error(), data.Icon.Error)
							}
						},
					},
					dec.PushButton{
						Text: data.S.ButtonDelete,
						OnClicked: func() {
							log.Println(data.S.Info, data.S.LogDelete)
							if err := wf.delete(db); err != nil {
								err = errors.Wrap(err, data.S.ErrorDelete)
								log.Println(data.S.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, err.Error(), data.Icon.Error)
							}
						},
					},
				},
			},
			dec.Composite{
				Layout:  dec.HBox{},
				Visible: !isChange,
				Children: []dec.Widget{
					dec.PushButton{
						Text: data.S.ButtonOK,
						OnClicked: func() {
							log.Println(data.S.Info, data.S.LogOk)
							if wf.modelTable.RowCount() > 0 && wf.tv.CurrentIndex() != -1 {
								index := wf.tv.CurrentIndex()
								idTitle.Id = wf.modelTable.items[index].Id
								idType := wf.modelTable.items[index].Type
								sType := wf.modelTable.mapIdTitle[idType]
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
	log.Printf(data.S.CreateWindow, data.S.Entities)

	log.Printf(data.S.RunWindow, data.S.Entities)
	return wf.Run(), nil
}

func (wf windowsFormEntities) add(db *sql.DB) error {
	var entity Entity
	children := make([]*EntityRecChild, 0)
	entity.Children = &children
	cmd, err := EntityRunDialog(wf, db, &entity)
	log.Printf(data.S.EndWindow, data.S.Entity, cmd)
	if err != nil {
		return errors.Wrapf(err, data.S.InEntityRunDialog, entity)
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
		return errors.Wrap(err, data.S.ErrorAddDB+QwStr)
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
		log.Println(data.S.Error, data.S.ErrorInsertIndexLog)
		walk.MsgBox(wf, data.S.MsgBoxError, data.S.ErrorInsertIndex, data.Icon.Critical)
		return nil
	}
	wf.modelTable.items[index].Id = id
	for _, v := range *entity.Children {
		QwStrChild := data.InsertEntityRec(id, v.Id, v.Count)
		if _, err := db.Exec(QwStrChild); err != nil {
			err = errors.Wrap(err, data.S.ErrorAddDB+QwStrChild)
			log.Println(data.S.Error, err)
			walk.MsgBox(wf, data.S.MsgBoxError, err.Error(), data.Icon.Error)
		}
	}
	return nil
}

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
	entity.Children = &children
	cmd, err := EntityRunDialog(wf, db, entity)
	log.Printf(data.S.EndWindow, data.S.Entity, cmd)

	if err != nil {
		return errors.Wrapf(err, data.S.InEntityRunDialog, *entity)
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
		return errors.Wrap(err, data.S.ErrorChangeDB+QwStr)
	}
	wf.modelTable.items[index] = entity
	wf.modelTable.PublishRowsReset()
	wf.modelTable.PublishRowChanged(index)
	return nil
}

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
		return errors.Wrap(err, data.S.ErrorDeleteDB+QwStr)
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
