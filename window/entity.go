package window

import (
	"accounting/data"
	"database/sql"
	"fmt"
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
		table := data.Tab["EntityRec"]
		sIdC := table.Columns["Child"].Name
		sIdP := table.Columns["Parent"].Name
		sCount := table.Columns["Count"].Name
		// GO-TO выборка названия и типа из таблицы entity
		QwStr := fmt.Sprintf("SELECT %s AS idc, %s AS count FROM %s WHERE %s = %d",
			sIdC, sCount, table.Name, sIdP, parent)
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorQuery+QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			row := EntityRecChild{}
			err := rows.Scan(&row.Id, &row.Count)
			if err != nil {
				err = errors.Wrap(err, data.S.ErrorDecryptRow)
				return err
			}
			arr = append(arr, &row)
		}
		return nil
	}()); err != nil {
		return arr, errors.Wrap(err, fmt.Sprintf("In SelectEntityRecChild(parent = %d)", parent))
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

func newModelEntityComponent() (*modelEntityComponent, error) { //GO-TO
	m := new(modelEntityComponent)
	m.items = make([]*EntityRecChild, 0)
	m.items = append(m.items, &EntityRecChild{IdTitle: IdTitle{Id: 2, Title: "asd"}, Count: 4})
	m.items = append(m.items, &EntityRecChild{IdTitle: IdTitle{Id: 3, Title: "qwe"}, Count: 8})
	m.items = append(m.items, &EntityRecChild{IdTitle: IdTitle{Id: 4, Title: "zxc"}, Count: 0})
	return m, nil
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
	// msgBoxIcon := walk.MsgBoxIconWarning
	sButtonAdd := " компонент"
	var err error
	var databind *walk.DataBinder
	wf := &windowsFormEntity{}
	wf.modelType, _, err = SelectIdTitle(db, "EntityType")
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorTypeInit)
		return 0, err
	}
	wf.modelTable, err = newModelEntityComponent()
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorTableInit)
		return 0, err
	}

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
								MinSize: dec.Size{150, 0},
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

							dec.Label{
								Text: "Маркировка:",
							},
							dec.CheckBox{
								Checked: dec.Bind("ProductionLine"),
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

								Model: wf.modelTable,
								OnCurrentIndexChanged: func() { //GO-TO
									fmt.Printf("CurrentIndexes: %v\n", wf.tv.CurrentIndex())
								},
								OnSelectedIndexesChanged: func() {
									fmt.Printf("SelectedIndexes: %v\n", wf.tv.SelectedIndexes())
								},
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
									// GO-TO
								},
							},
							dec.PushButton{
								ColumnSpan: 2,
								Text:       data.S.ButtonDelete + sButtonAdd,
								OnClicked: func() {
									// GO-TO
								},
							},

							dec.PushButton{
								Text: data.S.ButtonOK,
								OnClicked: func() {
									if err := databind.Submit(); err != nil {
										log.Println(err)
										// GO-TO
										return
									}
									wf.Accept()
								},
							},
							dec.PushButton{
								Text:      data.S.ButtonCansel,
								OnClicked: func() { wf.Cancel() },
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
	// cmd, err := EntityRunDialog(wf, db, false, &entity)
	log.Printf(data.S.EndWindow, data.S.EntityRec, cmd)
	if err != nil {
		return errors.Wrapf(err, "In EntityRecRunDialog(entity = %v)", entity)
	}
	if cmd != walk.DlgCmdOK {
		return nil
	}
	if entity.Id != 0 {
		if err = db.Ping(); err != nil {
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		table := data.Tab["EntityRec"]
		sIdC := table.Columns["Child"].Name
		sIdP := table.Columns["Parent"].Name
		sCount := table.Columns["Count"].Name
		QwStr := fmt.Sprintf("INSERT %s (%s, %s, %s) VALUES (%d, %d, %d)",
			table.Name, sIdP, entity.Id, sIdC, child.Id, sCount, child.Count)
		fmt.Println(QwStr)
		_, err := db.Exec(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorAddDB+QwStr)
		}
	}
	trackLatest := wf.tv.ItemVisible(len(wf.modelTable.items) - 1) //&& len(wf.tv.SelectedIndexes()) <= 1
	wf.modelTable.items = append(wf.modelTable.items, &child)
	index := len(wf.modelTable.items) - 1
	wf.modelTable.PublishRowsInserted(index, index)
	if trackLatest {
		wf.tv.EnsureItemVisible(len(wf.modelTable.items) - 1)
	}
	//GO-TO
	return nil
}
