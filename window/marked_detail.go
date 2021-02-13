package window

import (
	"accounting/data"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

// Структура, содержащая описание и переменные окна.
type windowsFormMarkedDetail struct {
	*walk.Dialog
	buttonChooseWidget *walk.PushButton
}

// Описание и запуск диалогового окна.
func MarkedDetailRunDialog(owner walk.Form, db *sql.DB, map3 *Map3, isChange bool, detail *MarkedDetail) (int, error) {
	log.Printf(data.S.BeginWindow, data.S.MarkedDetail)
	var databind *walk.DataBinder
	wf := windowsFormMarkedDetail{}
	log.Printf(data.S.InitWindow, data.S.MarkedDetail)
	if err := (dec.Dialog{
		AssignTo: &wf.Dialog,
		Title:    data.S.HeadingMarkedDetail,
		DataBinder: dec.DataBinder{
			AssignTo:       &databind,
			Name:           "detail",
			DataSource:     detail,
			ErrorPresenter: dec.ToolTipErrorPresenter{},
		},
		Layout: dec.VBox{MarginsZero: true},
		Children: []dec.Widget{
			dec.Composite{
				Layout: dec.Grid{Columns: 2},
				Children: []dec.Widget{
					dec.Label{
						Text: "Иерархия:",
					},
					dec.ComboBox{
						Value:         dec.Bind("Marking", dec.SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Title",
						MinSize:       dec.Size{150, 0},
						Model:         map3.MarkingLines(),
					},

					dec.Label{
						Text: "Маркировка:",
					},
					dec.LineEdit{
						MaxLength: 15,
						MinSize:   dec.Size{50, 0},
						Text:      dec.Bind("Mark"),
					},

					dec.Label{
						Text: "Родитель:", // GO-TO
					},
					dec.PushButton{
						AssignTo: &wf.buttonChooseWidget,
						MinSize:  dec.Size{150, 10},
						Text:     strconv.FormatInt(detail.Parent.Id, 10), // GO-TO
						OnClicked: func() {
							log.Println(data.S.Info, data.S.LogChoose)
							if err := (func() error {
								var parent int64
								cmd, err := MarkedDetailsRunDialog(wf, db, false, &parent)
								log.Printf(data.S.EndWindow, data.S.MarkedDetails, cmd)
								if err != nil {
									return errors.Wrapf(err, data.S.InMarkedDetailsRunDialog, false, parent)
								}
								if cmd == walk.DlgCmdOK {
									fmt.Println(parent)
									detail.Parent.Id = parent
									wf.buttonChooseWidget.SetText(strconv.FormatInt(parent, 10)) // GO-TO
								}
								return nil
							}()); err != nil {
								err = errors.Wrap(err, data.S.ErrorChoose)
								log.Println(data.S.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
							}
						},
					},
				},
			},
			dec.Composite{
				Layout: dec.HBox{},
				Children: []dec.Widget{
					dec.PushButton{
						Text: data.S.ButtonOK,
						OnClicked: func() {
							log.Println(data.S.Info, data.S.LogOk) // GO-TO проверить на корректность mark
							if err := databind.Submit(); err != nil {
								err = errors.Wrap(err, data.S.ErrorSubmit)
								log.Println(data.S.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
								return
							}
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
	}.Create(owner)); err != nil {
		err = errors.Wrap(err, data.S.ErrorCreateWindow)
		return 0, err
	}
	log.Printf(data.S.CreateWindow, data.S.MarkedDetail)

	log.Printf(data.S.RunWindow, data.S.MarkedDetail)
	return wf.Run(), nil
}
