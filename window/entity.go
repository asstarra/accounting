package window

import (
	"database/sql"
	"fmt"
	"log"

	// "strconv"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

type Entity struct {
	Id             int
	Title          string
	Type           int
	Specification  string
	ProductionLine bool
	Note           string
}

type modelEntityComponent struct {
	walk.TableModelBase
	items []*EntityRecChild
}

type windowsFormEntity struct {
	*walk.Dialog

	textTitle         *walk.LineEdit
	textSpecification *walk.LineEdit
	textNote          *walk.TextEdit
}

func newModelEntityComponent() *modelEntityComponent {
	m := new(modelEntityComponent)
	m.items = make([]*EntityRecChild, 0)
	m.items = append(m.items, &EntityRecChild{IdTitle: IdTitle{Id: 2, Title: "asd"}, Count: 4})
	m.items = append(m.items, &EntityRecChild{IdTitle: IdTitle{Id: 3, Title: "qwe"}, Count: 8})
	m.items = append(m.items, &EntityRecChild{IdTitle: IdTitle{Id: 4, Title: "zxc"}, Count: 0})
	return m
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
	log.Println("PANIC!", "unexpected col")
	panic("unexpected col")
}

func EntityRunDialog(owner walk.Form, entity *Entity, db *sql.DB) (int, error) {
	log.Println("INFO -", "BEGIN window - ENTITY")
	// sHeading := "Внимание"
	// msgBoxIcon := walk.MsgBoxIconWarning
	// var dlg *walk.Dialog
	var databind *walk.DataBinder
	modelType, err := SelectIdTitle("EntityType", db)
	if err != nil {
		err = errors.Wrap(err, "Не удалось узнать список типов.")
		log.Println("ERROR!", err)
		return 0, err
	}
	model := newModelEntityComponent()
	var tv *walk.TableView

	wf := &windowsFormEntity{}

	if err := (dec.Dialog{
		AssignTo: &wf.Dialog,
		Title:    "Сущность",
		DataBinder: dec.DataBinder{
			AssignTo:       &databind,
			Name:           "entity",
			DataSource:     entity,
			ErrorPresenter: dec.ToolTipErrorPresenter{},
		},
		// Size:     dec.Size{100, 100},
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
								Model:         modelType,
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
								AssignTo: &tv,
								Columns: []dec.TableViewColumn{
									{Title: "Название"},
									{Title: "Количество"},
								},
								ColumnSpan: 2,

								Model: model,
								OnCurrentIndexChanged: func() {
									fmt.Printf("CurrentIndexes: %v\n", tv.CurrentIndex())
								},
								OnSelectedIndexesChanged: func() {
									fmt.Printf("SelectedIndexes: %v\n", tv.SelectedIndexes())
								},
							},

							dec.PushButton{
								ColumnSpan: 2,
								Text:       "Добавить компонент",
								OnClicked: func() {
									// wf.Cancel()
								},
							},
							dec.PushButton{
								ColumnSpan: 2,
								Text:       "Изменить компонент",
								OnClicked: func() {
									// wf.Cancel()
								},
							},
							dec.PushButton{
								ColumnSpan: 2,
								Text:       "Удалить компонент",
								OnClicked: func() {
									// wf.Cancel()
								},
							},

							dec.PushButton{
								Text: "OK",
								OnClicked: func() {
									if err := databind.Submit(); err != nil {
										log.Println(err)
										return
									}
									wf.Accept()
								},
							},
							dec.PushButton{
								Text:      "Отмена",
								OnClicked: func() { wf.Cancel() },
							},
						},
					},
				},
			},
		},
	}.Create(owner)); err != nil {
		err = errors.Wrap(err, "Could not create Window Form Entity")
		log.Println("ERROR!", err)
		return 0, err
	}
	log.Println("INFO -", "CREATE window - ENTITY")

	log.Println("INFO -", "RUN window - ENTITY")
	return wf.Run(), nil
}
