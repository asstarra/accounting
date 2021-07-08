package window

import (
	"accounting/data"
	e "accounting/data/errors"
	l "accounting/data/log"
	"accounting/data/qwery"
	"accounting/data/text"
	. "accounting/window/data"
	"database/sql"
	"log"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

func SelectMarkedDetails(db *sql.DB, markings []int64) ([]*MarkedDetail, error) { //GO-TO ?
	arr := make([]*MarkedDetail, 0)
	if err := (func() error {
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, e.Err.ErrorPingDB)
		}
		QwStr := qwery.SelectMarkedDetail(markings) //GO-TO ?
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrapf(err, e.Err.ErrorQueryDB, QwStr)
		}
		defer rows.Close()
		var parentId, parentMarking sql.NullInt64
		var parentMark sql.NullString
		for rows.Next() {
			row := MarkedDetail{}
			err := rows.Scan(&row.Id, &row.Marking, &row.Mark, &parentId, &parentMarking, &parentMark)
			if err != nil {
				return errors.Wrap(err, e.Err.ErrorDecryptRow)
			}
			row.Parent.Id = parentId.Int64
			row.Parent.Marking = parentMarking.Int64
			row.Parent.Mark = parentMark.String
			arr = append(arr, &row)
		}
		return nil
	}()); err != nil {
		return arr, errors.Wrapf(err, l.In.InSelectMarkedDetails, markings)
	}
	return arr, nil
}

// Сруктура, содержащая модель таблицы.
type modelMarkedDetailsComponent struct {
	walk.TableModelBase
	items []*MarkedDetail
	Map3
}

// Структура, содержащая описание и переменные окна.
type windowsFormMarkedDetails struct {
	*walk.Dialog
	modelTable             *modelMarkedDetailsComponent
	tv                     *walk.TableView
	orderW, detailW, lineW *walk.ComboBox
	orderM, detailM, lineM []*Id64Title
}

// Инициализация модели окна.
func newWindowsFormMarkedDetails(db *sql.DB, parent *MarkedDetailMin) (*windowsFormMarkedDetails, error) {
	if db == nil {
		return nil, errors.New(e.Err.ErrorNil)
	}
	var err error
	wf := new(windowsFormMarkedDetails)
	wf.modelTable = new(modelMarkedDetailsComponent)
	wf.modelTable.Map3, err = NewMap3(db, false)
	if err != nil {
		return nil, err
	}
	wf.orderM = wf.modelTable.Map3.ModelOrders(0, 0, true)
	wf.detailM = wf.modelTable.Map3.ModelDetails(0, 0, true)
	wf.lineM = wf.modelTable.Map3.ModelMarkingLines(0, 0, true)
	wf.modelTable.items, err = SelectMarkedDetails(db, []int64{})
	if err != nil {
		err = errors.Wrap(err, e.Err.ErrorTableInit)
		return nil, err
	}
	return wf, nil
}

func (m *modelMarkedDetailsComponent) RowCount() int {
	return len(m.items)
}

func (m *modelMarkedDetailsComponent) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return m.Map3.MarkingToString(item.Marking)
	case 1:
		return item.Mark
	case 2:
		return m.Map3.MarkedDetailMinToString(item.Parent)
	}
	log.Println(l.Panic, e.Err.ErrorUnexpectedColumn)
	panic(e.Err.ErrorUnexpectedColumn)
}

// Описание и запуск диалогового окна.
func MarkedDetailsRunDialog(owner walk.Form, db *sql.DB, isChange bool, parent *MarkedDetailMin) (int, error) {
	log.Printf(l.BeginWindow, l.MarkedDetails)
	var err error
	var databind *walk.DataBinder
	var search = new(struct {
		Order, Detail, Line int64
	})
	wf, err := newWindowsFormMarkedDetails(db, parent)
	if err != nil {
		return 0, errors.Wrap(err, e.Err.ErrorInit)
	}
	log.Printf(l.InitWindow, l.MarkedDetails)
	if err = (dec.Dialog{
		AssignTo: &wf.Dialog,
		Title:    text.T.HeadingMarkedDetails,
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
						Text: "Заказ:",
					},
					dec.ComboBox{
						AssignTo:      &wf.orderW,
						Value:         dec.Bind("Order", dec.SelRequired{}),
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
						Value:         dec.Bind("Detail", dec.SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Title",
						MinSize:       dec.Size{120, 0},
						Model:         wf.detailM,
						OnCurrentIndexChanged: func() {
							wf.setLineCmbx()
						},
					},
					dec.HSpacer{Size: 20},

					dec.PushButton{
						Text: text.T.ButtonSearch,
						OnClicked: func() {
							log.Println(l.Info, l.LogSearch)
							err := databind.Submit()
							if err != nil {
								err = errors.Wrap(err, e.Err.ErrorSubmit)
								log.Println(l.Error, err)
								walk.MsgBox(wf, text.T.MsgBoxError, MsgError(err), data.Icon.Error)
								return
							}
							markings := wf.modelTable.Map3.ToMarkingIds(search.Order, search.Detail, search.Line)
							lastLen := wf.modelTable.RowCount()
							if items, err := SelectMarkedDetails(db, markings); err != nil {
								err = errors.Wrap(err, e.Err.ErrorSubquery)
								log.Println(l.Error, err)
								walk.MsgBox(wf, text.T.MsgBoxError, MsgError(err), data.Icon.Error)
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
			dec.Composite{
				Layout: dec.HBox{MarginsZero: true},
				Children: []dec.Widget{
					dec.Label{
						Text: "Иерархия:",
					},
					dec.ComboBox{
						AssignTo:      &wf.lineW,
						Value:         dec.Bind("Line", dec.SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Title",
						MinSize:       dec.Size{150, 0},
						Model:         wf.lineM,
					},
				},
			},
			dec.TableView{
				AssignTo: &wf.tv,
				Columns: []dec.TableViewColumn{
					{Title: "Линия"},
					{Title: "Маркировка"},
					{Title: "Родитель"}, // GO-TO
				},
				MinSize: dec.Size{0, 200},
				Model:   wf.modelTable,
			},
			dec.Composite{
				Layout:  dec.HBox{MarginsZero: true},
				Visible: parent == nil,
				Children: []dec.Widget{
					dec.PushButton{
						Text: text.T.ButtonAdd,
						OnClicked: func() {
							log.Println(l.Info, l.LogAdd)
							if err := wf.add(db); err != nil {
								err = errors.Wrap(err, e.Err.ErrorAddRow)
								log.Println(l.Error, err)
								walk.MsgBox(wf, text.T.MsgBoxError, MsgError(err), data.Icon.Error)
							}
						},
					},
					dec.PushButton{
						Text: text.T.ButtonChange,
						OnClicked: func() {
							log.Println(l.Info, l.LogChange)
							if err := wf.change(db); err != nil {
								err = errors.Wrap(err, e.Err.ErrorChangeRow)
								log.Println(l.Error, err)
								walk.MsgBox(wf, text.T.MsgBoxError, MsgError(err), data.Icon.Error)
							}
						},
					},
					dec.PushButton{
						Text: text.T.ButtonDelete,
						OnClicked: func() {
							log.Println(l.Info, l.LogDelete)
							if err := wf.delete(db); err != nil {
								err = errors.Wrap(err, e.Err.ErrorDeleteRow)
								log.Println(l.Error, err)
								walk.MsgBox(wf, text.T.MsgBoxError, MsgError(err), data.Icon.Error)
							}
						},
					},
				},
			},
			dec.Composite{
				Layout:  dec.HBox{MarginsZero: true},
				Visible: parent != nil,
				Children: []dec.Widget{
					dec.PushButton{
						Text: text.T.ButtonOK,
						OnClicked: func() {
							log.Println(l.Info, l.LogOk)
							if wf.modelTable.RowCount() > 0 && wf.tv.CurrentIndex() != -1 {
								index := wf.tv.CurrentIndex()
								*parent = wf.modelTable.items[index].MarkedDetailMin
								wf.Accept()
							} else {
								walk.MsgBox(wf, text.T.MsgBoxInfo, text.T.MsgChooseRow, data.Icon.Info)
							}
						},
					},
					dec.PushButton{
						Text:      text.T.ButtonCansel,
						OnClicked: func() { wf.Cancel() },
					},
				},
			},
		},
	}.Create(owner)); err != nil {
		err = errors.Wrap(err, e.Err.ErrorCreateWindow)
		return 0, err
	}
	log.Printf(l.CreateWindow, l.MarkedDetails)

	log.Printf(l.RunWindow, l.MarkedDetails)
	return wf.Run(), nil
}

// Функция, для добавления строки в таблицу.
func (wf windowsFormMarkedDetails) add(db *sql.DB) error {
	var detail MarkedDetail
	cmd, err := MarkedDetailRunDialog(wf, db, &wf.modelTable.Map3, false, &detail)
	log.Printf(l.EndWindow, l.Entity, cmd)
	if err != nil {
		return errors.Wrapf(err, l.In.InMarkedDetailRunDialog, false, detail)
	}
	if cmd != walk.DlgCmdOK {
		return nil
	}
	QwStr := qwery.InsertMarkedDetail(detail.Marking, detail.Parent.Id, detail.Mark)
	if err = db.Ping(); err != nil {
		return errors.Wrap(err, e.Err.ErrorPingDB)
	}

	result, err := db.Exec(QwStr)
	if err != nil {
		return errors.Wrapf(err, e.Err.ErrorAddDB, QwStr)
	}

	trackLatest := wf.tv.ItemVisible(len(wf.modelTable.items) - 1)
	wf.modelTable.items = append(wf.modelTable.items, &detail)
	index := len(wf.modelTable.items) - 1
	wf.modelTable.PublishRowsReset()
	wf.modelTable.PublishRowsInserted(index, index)
	if trackLatest {
		wf.tv.EnsureItemVisible(len(wf.modelTable.items) - 1)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println(l.Error, e.Err.ErrorInsertIndexLog)
		walk.MsgBox(wf, text.T.MsgBoxError, e.Err.ErrorInsertIndex, data.Icon.Critical)
		return nil
	}
	wf.modelTable.items[index].Id = id
	return nil
}

// Функция, для изменения строки в таблице.
func (wf windowsFormMarkedDetails) change(db *sql.DB) error {
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 {
		walk.MsgBox(wf, text.T.MsgBoxInfo, text.T.MsgChooseRow, data.Icon.Info)
		return nil
	}
	var err error
	index := wf.tv.CurrentIndex()
	detail := wf.modelTable.items[index]
	cmd, err := MarkedDetailRunDialog(wf, db, &wf.modelTable.Map3, true, detail)
	log.Printf(l.EndWindow, l.MarkedDetail, cmd)

	if err != nil {
		return errors.Wrapf(err, l.In.InMarkedDetailRunDialog, true, *detail)
	}
	if cmd != walk.DlgCmdOK {
		return nil
	}
	QwStr := qwery.UpdateMarkedDetail(detail.Id, detail.Marking, detail.Parent.Id, detail.Mark)
	if err = db.Ping(); err != nil {
		return errors.Wrap(err, e.Err.ErrorPingDB)
	}
	_, err = db.Exec(QwStr)
	if err != nil {
		return errors.Wrapf(err, e.Err.ErrorChangeDB, QwStr)
	}
	wf.modelTable.items[index] = detail
	wf.modelTable.PublishRowsReset()
	wf.modelTable.PublishRowChanged(index)
	return nil
}

// Функция, для удаления строки из таблицы.
func (wf windowsFormMarkedDetails) delete(db *sql.DB) error {
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 {
		walk.MsgBox(wf, text.T.MsgBoxInfo, text.T.MsgChooseRow, data.Icon.Info)
		return nil
	}
	index := wf.tv.CurrentIndex()
	id := wf.modelTable.items[index].Id
	QwStr := qwery.DeleteMarkedDetail(id)
	if err := db.Ping(); err != nil {
		return errors.Wrap(err, e.Err.ErrorPingDB)
	}
	_, err := db.Exec(QwStr)
	if err != nil {
		return errors.Wrapf(err, e.Err.ErrorDeleteDB, QwStr)
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

func (wf windowsFormMarkedDetails) setLineCmbx() {
	oi := wf.orderM[MaxInt(wf.orderW.CurrentIndex(), 0)].Id
	di := wf.detailM[MaxInt(wf.detailW.CurrentIndex(), 0)].Id

	wf.lineM = wf.modelTable.Map3.ModelMarkingLines(oi, di, true)
	wf.lineW.SetModel(wf.lineM)
}
