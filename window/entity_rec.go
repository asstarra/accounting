package window

import (
	"accounting/data"
	"database/sql"
	"log"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

// Структура, содержащая описание и переменные окна.
type windowsFormEntityRec struct {
	*walk.Dialog
	buttonEntitiesWidget *walk.PushButton
}

// Описание и запуск диалогового окна.
func EntityRecRunDialog(owner walk.Form, db *sql.DB, isChange bool, child *EntityRecChild) (int, error) {
	log.Printf(data.Log.BeginWindow, data.Log.EntityRec)
	var databind *walk.DataBinder
	wf := windowsFormEntityRec{}
	log.Printf(data.Log.InitWindow, data.Log.EntityRec)
	if err := (dec.Dialog{
		AssignTo: &wf.Dialog,
		Title:    data.S.HeadingEntityRec,
		DataBinder: dec.DataBinder{
			AssignTo:       &databind,
			Name:           "child",
			DataSource:     child,
			ErrorPresenter: dec.ToolTipErrorPresenter{},
		},
		Layout: dec.VBox{MarginsZero: true},
		Children: []dec.Widget{
			dec.Composite{
				Layout: dec.Grid{Columns: 2},
				Children: []dec.Widget{
					dec.Label{
						Text: "Название:",
					},
					dec.PushButton{
						AssignTo: &wf.buttonEntitiesWidget,
						Enabled:  !isChange,
						MinSize:  dec.Size{150, 10},
						Text:     child.Title,
						OnClicked: func() {
							log.Println(data.Log.Info, data.Log.LogChoose)
							if err := (func() error {
								idTitle := child.IdTitle
								cmd, err := EntitiesRunDialog(wf, db, false, &idTitle)
								log.Printf(data.Log.EndWindow, data.Log.Entities, cmd)
								if err != nil {
									return errors.Wrapf(err, data.Log.InEntitiesRunDialog, false, idTitle)
								}
								if cmd == walk.DlgCmdOK {
									child.IdTitle = idTitle
									wf.buttonEntitiesWidget.SetText(child.Title)
								}
								return nil
							}()); err != nil {
								err = errors.Wrap(err, data.S.ErrorChoose)
								log.Println(data.Log.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
							}
						},
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
						Text: data.S.ButtonOK,
						OnClicked: func() {
							log.Println(data.Log.Info, data.Log.LogOk)
							if err := databind.Submit(); err != nil {
								err = errors.Wrap(err, data.S.ErrorSubmit)
								log.Println(data.Log.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
								return
							}
							wf.Accept()
						},
					},
					dec.PushButton{
						Text: data.S.ButtonCansel,
						OnClicked: func() {
							log.Println(data.Log.Info, data.Log.LogCansel)
							wf.Cancel()
						},
					},
				},
			},
		},
	}.Create(owner)); err != nil {
		err = errors.Wrap(err, data.S.ErrorCreateWindow)
		return 0, err
	}
	log.Printf(data.Log.CreateWindow, data.Log.EntityRec)

	log.Printf(data.Log.RunWindow, data.Log.EntityRec)
	return wf.Run(), nil
}
