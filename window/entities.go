package window

import (
	"database/sql"
	"log"

	// "strconv"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

type windowsFormEntities struct {
	*walk.Dialog
}

func EntitiesRunDialog(owner walk.Form, entity *EntityRecChild, db *sql.DB) (int, error) {
	log.Println("INFO -", "BEGIN window - ENTITIES")
	// sHeading := "Внимание"
	// msgBoxIcon := walk.MsgBoxIconWarning
	// var dlg *walk.Dialog
	var databind *walk.DataBinder
	// var wf *walk.Dialog

	wf := &windowsFormEntities{}

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
		Layout: dec.VBox{},
		Children: []dec.Widget{

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
		err = errors.Wrap(err, "Could not create Window Form Entities")
		log.Println("ERROR!", err)
		return 0, err
	}
	log.Println("INFO -", "CREATE window - ENTITIES")

	log.Println("INFO -", "RUN window - ENTITIES")
	return wf.Run(), nil
}
