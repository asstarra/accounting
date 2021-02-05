package window

import (
	"accounting/data"
	"database/sql"
	"log"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

func SelectEntityRecChild(db *sql.DB, parent int64) ([]*EntityRecChild, error) {
	arr := make([]*EntityRecChild, 0)
	if err := (func() error {
		if err := db.Ping(); err != nil {
			err = errors.Wrap(err, data.S.ErrorPingDB)
			return err
		}

		QwStr := data.SelectEntityRecChild(parent)
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorQuery+QwStr)
		}
		defer rows.Close()
		var e_type, title string
		for rows.Next() {
			row := EntityRecChild{}
			err := rows.Scan(&row.Id, &e_type, &title, &row.Count)
			if err != nil {
				err = errors.Wrap(err, data.S.ErrorDecryptRow)
				return err
			}
			row.Title = e_type + " " + title
			arr = append(arr, &row)
		}
		return nil
	}()); err != nil {
		return arr, errors.Wrapf(err, data.S.InSelectEntityRecChild, parent)
	}
	return arr, nil
}

type modelEntityComponent struct {
	walk.TableModelBase
	items []*EntityRecChild
}

type windowsFormEntity struct {
	*walk.Dialog
	modelType  []*IdTitle
	modelTable *modelEntityComponent
	tv         *walk.TableView
}

func newWindowsFormEntity(db *sql.DB, entity *Entity) (*windowsFormEntity, error) {
	var err error
	wf := new(windowsFormEntity)
	wf.modelTable = new(modelEntityComponent)
	wf.modelTable.items = *entity.Children
	wf.modelType, _, err = SelectIdTitle(db, "EntityType")
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorTypeInit)
		return nil, err
	}
	return wf, nil
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
	log.Println(data.S.Panic, data.S.ErrorUnexpectedColumn)
	panic(data.S.ErrorUnexpectedColumn)
}

func EntityRunDialog(owner walk.Form, db *sql.DB, entity *Entity) (int, error) {
	log.Printf(data.S.BeginWindow, data.S.Entity)
	sButtonAdd := " компонент" // GO-TO возможно нужно? вынести строку
	var databind *walk.DataBinder
	wf, err := newWindowsFormEntity(db, entity)
	if err != nil {
		return 0, errors.Wrap(err, data.S.ErrorInit)
	}
	log.Printf(data.S.InitWindow, data.S.Entity)
	if err := (dec.Dialog{
		AssignTo: &wf.Dialog,
		Title:    data.S.HeadingEntity,
		DataBinder: dec.DataBinder{
			AssignTo:       &databind,
			Name:           "entity",
			DataSource:     entity,
			ErrorPresenter: dec.ToolTipErrorPresenter{},
		},
		Layout: dec.VBox{MarginsZero: true},
		Children: []dec.Widget{
			dec.HSplitter{
				Children: []dec.Widget{
					dec.Composite{
						Layout: dec.Grid{Columns: 2},
						Children: []dec.Widget{
							dec.Label{
								Text: "Название:",
							},
							dec.LineEdit{
								MinSize: dec.Size{170, 0},
								Text:    dec.Bind("Title"),
							},

							dec.Label{
								Text: "Тип:",
							},
							dec.ComboBox{
								Value:         dec.Bind("Type", dec.SelRequired{}),
								BindingMember: "Id",
								DisplayMember: "Title",
								Model:         wf.modelType,
							},

							dec.Label{
								Text: "Спецификация:",
							},
							dec.LineEdit{
								Text: dec.Bind("Specification"),
							},

							dec.RadioButtonGroupBox{
								ColumnSpan: 2,
								Title:      "Маркировка:",
								Layout:     dec.HBox{},
								DataMember: "Marking",
								Buttons: []dec.RadioButton{
									{Text: MarkingTitle[MarkingNo], Value: MarkingNo},
									{Text: MarkingTitle[MarkingAll], Value: MarkingAll},
									{Text: MarkingTitle[MarkingYear], Value: MarkingYear},
								},
							},

							dec.Label{
								ColumnSpan: 2,
								Text:       "Примечание:",
							},
							dec.TextEdit{
								ColumnSpan: 2,
								MinSize:    dec.Size{0, 100},
								Text:       dec.Bind("Note"),
							},
						},
					},
					dec.Composite{
						Layout:  dec.Grid{Columns: 2},
						MinSize: dec.Size{230, 300},
						Children: []dec.Widget{
							dec.Label{
								ColumnSpan: 2,
								Text:       "Компоненты:",
							},
							dec.TableView{
								AssignTo: &wf.tv,
								Columns: []dec.TableViewColumn{
									{Title: "Название"},
									{Title: "Количество"},
								},
								ColumnSpan: 2,
								Model:      wf.modelTable,
							},

							dec.PushButton{
								ColumnSpan: 2,
								Text:       data.S.ButtonAdd + sButtonAdd,
								OnClicked: func() {
									log.Println(data.S.Info, data.S.LogAdd)
									if err := wf.add(db, entity); err != nil {
										err = errors.Wrap(err, data.S.ErrorAdd)
										log.Println(data.S.Error, err)
										walk.MsgBox(wf, data.S.MsgBoxError, err.Error(), data.Icon.Error)
									}
								},
							},
							dec.PushButton{
								ColumnSpan: 2,
								Text:       data.S.ButtonChange + sButtonAdd,
								OnClicked: func() {
									log.Println(data.S.Info, data.S.LogChange)
									if err := wf.change(db, entity); err != nil {
										err = errors.Wrap(err, data.S.ErrorChange)
										log.Println(data.S.Error, err)
										walk.MsgBox(wf, data.S.MsgBoxError, err.Error(), data.Icon.Error)
									}
								},
							},
							dec.PushButton{
								ColumnSpan: 2,
								Text:       data.S.ButtonDelete + sButtonAdd,
								OnClicked: func() {
									log.Println(data.S.Info, data.S.LogDelete)
									if err := wf.delete(db, entity); err != nil {
										err = errors.Wrap(err, data.S.ErrorDelete)
										log.Println(data.S.Error, err)
										walk.MsgBox(wf, data.S.MsgBoxError, err.Error(), data.Icon.Error)
									}
								},
							},

							dec.PushButton{
								Text: data.S.ButtonOK,
								OnClicked: func() {
									log.Println(data.S.Info, data.S.LogOk)
									if err := databind.Submit(); err != nil {
										err = errors.Wrap(err, data.S.ErrorSubmit)
										log.Println(data.S.Error, err)
										walk.MsgBox(wf, data.S.MsgBoxError, err.Error(), data.Icon.Error)
										return
									}
									entity.Children = &wf.modelTable.items
									wf.Accept()
								},
							},
							dec.PushButton{
								Text: data.S.ButtonCansel,
								OnClicked: func() {
									log.Println(data.S.Info, data.S.LogCansel)
									wf.Cancel()
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
	log.Printf(data.S.CreateWindow, data.S.Entity)

	log.Printf(data.S.RunWindow, data.S.Entity)
	return wf.Run(), nil
}

func (wf windowsFormEntity) add(db *sql.DB, entity *Entity) error {
	child := EntityRecChild{}
	cmd, err := EntityRecRunDialog(wf, db, false, &child)
	log.Printf(data.S.EndWindow, data.S.EntityRec, cmd)
	if err != nil {
		return errors.Wrapf(err, data.S.InEntityRecRunDialog, child)
	}
	if cmd != walk.DlgCmdOK {
		return nil
	}
	if entity.Id != 0 {
		QwStr := data.InsertEntityRec(entity.Id, child.Id, child.Count)
		if err = db.Ping(); err != nil {
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		result, err := db.Exec(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorAddDB+QwStr)
		}
		if id, err := result.LastInsertId(); err != nil {
			log.Println(data.S.Error, errors.Wrap(err, data.S.ErrorInsertIndexLog))
			walk.MsgBox(wf, data.S.MsgBoxError, data.S.ErrorInsertIndex, data.Icon.Critical)
		} else {
			child.Id = id
		}
	}
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

func (wf windowsFormEntity) change(db *sql.DB, entity *Entity) error {
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 {
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChooseRow, data.Icon.Info)
		return nil
	}
	index := wf.tv.CurrentIndex()
	child := wf.modelTable.items[index]
	cmd, err := EntityRecRunDialog(wf, db, true, child)
	log.Printf(data.S.EndWindow, data.S.EntityRec, cmd)
	if err != nil {
		return errors.Wrapf(err, data.S.InEntityRecRunDialog, child)
	}
	if cmd != walk.DlgCmdOK {
		return nil
	}
	if entity.Id != 0 {
		QwStr := data.UpdateEntityRec(entity.Id, child.Id, child.Count)
		if err = db.Ping(); err != nil {
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		if _, err := db.Exec(QwStr); err != nil {
			return errors.Wrap(err, data.S.ErrorChangeDB+QwStr)
		}
	}
	wf.modelTable.PublishRowsReset()
	wf.modelTable.PublishRowChanged(index)
	return nil
}

func (wf windowsFormEntity) delete(db *sql.DB, entity *Entity) error {
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 {
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChooseRow, data.Icon.Info)
		return nil
	}
	index := wf.tv.CurrentIndex()
	id := wf.modelTable.items[index].Id
	if entity.Id != 0 {
		QwStr := data.DeleteEntityRec(entity.Id, id)
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		if _, err := db.Exec(QwStr); err != nil {
			return errors.Wrap(err, data.S.ErrorDeleteDB+QwStr)
		}
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
