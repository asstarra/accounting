package window

import (
	"database/sql"
	"log"

	// "strconv"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

type EntityRecChild struct {
	IdTitle
	Count int
}

func EntityRecRunDialog(owner walk.Form, entity *EntityRecChild, db *sql.DB) (int, error) {
	log.Println("INFO -", "BEGIN window - ENTITY_REC")
	// sHeading := "Внимание"
	// msgBoxIcon := walk.MsgBoxIconWarning
	// var dlg *walk.Dialog
	var databind *walk.DataBinder
	// var wf *walk.Dialog
	var wf struct {
		*walk.Dialog
	}
	//wf := &windowsFormEntity{}

	if err := (dec.Dialog{
		AssignTo: &wf.Dialog,
		Title:    "Компонент",
		DataBinder: dec.DataBinder{
			AssignTo:       &databind,
			Name:           "entity",
			DataSource:     entity,
			ErrorPresenter: dec.ToolTipErrorPresenter{},
		},
		// Size:     dec.Size{100, 100},
		Layout: dec.VBox{MarginsZero: true},
		Children: []dec.Widget{
			dec.Composite{
				Layout: dec.Grid{Columns: 2},
				Children: []dec.Widget{
					dec.Label{
						Text: "Название:",
					},
					dec.PushButton{
						MinSize: dec.Size{150, 10},
						Text:    dec.Bind("Title"),
					},

					dec.Label{
						Text: "Количество:",
					},
					dec.NumberEdit{
						Value:    dec.Bind("Count", dec.Range{0, 1000}),
						Suffix:   " шт",
						Decimals: 0,
					},
				},
			},
			dec.Composite{
				Layout: dec.HBox{},
				Children: []dec.Widget{
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
	}.Create(owner)); err != nil {
		err = errors.Wrap(err, "Could not create Window Form Entity_rec")
		log.Println("ERROR!", err)
		return 0, err
	}
	log.Println("INFO -", "CREATE window - ENTITY_REC")

	log.Println("INFO -", "RUN window - ENTITY_REC")
	return wf.Run(), nil
}
