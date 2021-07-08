package window

import (
	"accounting/data"
	e "accounting/data/errors"
	l "accounting/data/log"
	"accounting/data/text"
	. "accounting/window/data"
	"database/sql"
	"log"

	"fmt"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

// Структура, содержащая описание и переменные окна.
type windowsFormMarkedDetail struct {
	*walk.Dialog
	buttonChooseWidget     *walk.PushButton
	Map3                   *Map3
	orderW, detailW, lineW *walk.ComboBox
	orderM, detailM, lineM []*Id64Title
}

// Инициализация модели окна.
func newWindowsFormMarkedDetail(db *sql.DB, map3 *Map3, detail *MarkedDetail) (*windowsFormMarkedDetail, error) {
	// var err error
	if db == nil || map3 == nil || detail == nil { // GO-TO в других файлах проверить корректность указателей.
		return nil, errors.New(e.Err.ErrorNil)
	}
	wf := new(windowsFormMarkedDetail)
	wf.Map3 = map3
	if detail.Id == 0 {
		wf.orderM = map3.ModelOrders(0, 0, true)
		wf.detailM = map3.ModelDetails(0, 0, true)
		wf.lineM = map3.ModelMarkingLines(0, 0, false)
	} else {
		wf.orderM = map3.ModelOrders(0, detail.Marking, false)
		wf.detailM = map3.ModelDetails(0, detail.Marking, false)
		wf.lineM = map3.ModelMarkingLines(wf.orderM[0].Id, wf.detailM[0].Id, false)
	}
	return wf, nil
}

// Описание и запуск диалогового окна.
func MarkedDetailRunDialog(owner walk.Form, db *sql.DB, map3 *Map3, isChange bool, detail *MarkedDetail) (int, error) {
	log.Printf(l.BeginWindow, l.MarkedDetail)
	var databind *walk.DataBinder
	wf, err := newWindowsFormMarkedDetail(db, map3, detail)
	if err != nil {
		return 0, errors.Wrap(err, e.Err.ErrorInit)
	}
	log.Printf(l.InitWindow, l.MarkedDetail)
	if err := (dec.Dialog{
		AssignTo: &wf.Dialog,
		Title:    text.T.HeadingMarkedDetail,
		DataBinder: dec.DataBinder{
			AssignTo:       &databind,
			Name:           "detail",
			DataSource:     detail,
			ErrorPresenter: dec.ToolTipErrorPresenter{},
		},
		Layout: dec.VBox{},
		Children: []dec.Widget{
			dec.Composite{
				Layout: dec.HBox{MarginsZero: true},
				Children: []dec.Widget{
					dec.Label{
						Text: "Заказ:",
					},
					dec.ComboBox{
						AssignTo:      &wf.orderW,
						BindingMember: "Id",
						DisplayMember: "Title",
						MinSize:       dec.Size{120, 0},
						Model:         wf.orderM,
						OnCurrentIndexChanged: func() {
							wf.setLineCmbx()
						},
					},
					dec.HSpacer{Size: 20},

					dec.Label{
						Text: "Деталь:",
					},
					dec.ComboBox{
						AssignTo:      &wf.detailW,
						BindingMember: "Id",
						DisplayMember: "Title",
						MinSize:       dec.Size{120, 0},
						Model:         wf.detailM,
						OnCurrentIndexChanged: func() {
							wf.setLineCmbx()
						},
					},
				},
			},
			dec.Composite{
				Layout: dec.Grid{Columns: 2, MarginsZero: true},
				Children: []dec.Widget{
					dec.Label{
						Text: "Иерархия:",
					},
					dec.ComboBox{
						AssignTo:      &wf.lineW,
						Value:         dec.Bind("Marking", dec.SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Title",
						MinSize:       dec.Size{150, 0},
						Model:         wf.lineM,
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
						Text: "Родитель:", // GO-TO переименовать
					},
					dec.PushButton{
						AssignTo: &wf.buttonChooseWidget,
						MinSize:  dec.Size{150, 10},
						Text:     wf.Map3.MarkedDetailMinToString(detail.Parent),
						OnClicked: func() {
							log.Println(l.Info, l.LogChoose)
							if err := (func() error {
								var parent MarkedDetailMin = detail.Parent
								if parent.Id == 0 {
									parent.Marking = detail.Marking
								}
								fmt.Println(parent)                                        //GO-TO
								cmd, err := MarkedDetailsRunDialog(wf, db, false, &parent) // GO-TO выбор только из возможных родителей.
								log.Printf(l.EndWindow, l.MarkedDetails, cmd)
								if err != nil {
									return errors.Wrapf(err, l.In.InMarkedDetailsRunDialog, false, parent)
								}
								if cmd == walk.DlgCmdOK {
									detail.Parent = parent
									wf.buttonChooseWidget.SetText(wf.Map3.MarkedDetailMinToString(parent))
								}
								return nil
							}()); err != nil {
								err = errors.Wrap(err, e.Err.ErrorChoose)
								log.Println(l.Error, err)
								walk.MsgBox(wf, text.T.MsgBoxError, MsgError(err), data.Icon.Error)
							}
						},
					},
				},
			},
			dec.Composite{
				Layout: dec.HBox{MarginsZero: true},
				Children: []dec.Widget{
					dec.PushButton{
						Text: text.T.ButtonOK,
						OnClicked: func() {
							log.Println(l.Info, l.LogOk) // GO-TO проверить на корректность mark
							if err := databind.Submit(); err != nil {
								err = errors.Wrap(err, e.Err.ErrorSubmit)
								log.Println(l.Error, err)
								walk.MsgBox(wf, text.T.MsgBoxError, MsgError(err), data.Icon.Error)
								return
							}
							wf.Accept()
						},
					},
					dec.PushButton{
						Text: text.T.ButtonCansel,
						OnClicked: func() {
							log.Println(l.Info, l.LogCansel)
							wf.Cancel()
						},
					},
				},
			},
		},
	}.Create(owner)); err != nil {
		err = errors.Wrap(err, e.Err.ErrorCreateWindow)
		return 0, err
	}
	log.Printf(l.CreateWindow, l.MarkedDetail)

	log.Printf(l.RunWindow, l.MarkedDetail)
	return wf.Run(), nil
}

func (wf windowsFormMarkedDetail) setLineCmbx() {
	oi := wf.orderM[MaxInt(wf.orderW.CurrentIndex(), 0)].Id
	di := wf.detailM[MaxInt(wf.detailW.CurrentIndex(), 0)].Id

	wf.lineM = wf.Map3.ModelMarkingLines(oi, di, true)
	wf.lineW.SetModel(wf.lineM)
}
