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

func SelectMarkedDetails(db *sql.DB, markings []int64) ([]*MarkedDetail, error) { //GO-TO
	arr := make([]*MarkedDetail, 0)
	if err := (func() error {
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		QwStr := data.SelectMarkedDetail(markings) //GO-TO
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorQueryDB+QwStr)
		}
		defer rows.Close()
		var parentId, parentMarking sql.NullInt64
		var parentMark sql.NullString
		for rows.Next() {
			row := MarkedDetail{}
			err := rows.Scan(&row.Id, &row.Marking, &row.Mark, &parentId, &parentMarking, &parentMark)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorDecryptRow)
			}
			row.Parent.Id = parentId.Int64
			row.Parent.Marking = parentMarking.Int64
			row.Parent.Mark = parentMark.String
			arr = append(arr, &row)
		}
		return nil
	}()); err != nil {
		return arr, errors.Wrapf(err, data.S.InSelectMarkedDetails, markings)
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
	modelTable *modelMarkedDetailsComponent
	tv         *walk.TableView
}

// Инициализация модели окна.
func newWindowsFormMarkedDetails(db *sql.DB) (*windowsFormMarkedDetails, error) {
	var err error
	wf := new(windowsFormMarkedDetails)
	wf.modelTable = new(modelMarkedDetailsComponent)
	wf.modelTable.MapIdToMarkingLine, wf.modelTable.MapIdToEntity, err = UpdateMarkingLine(db, false)
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorSubquery)
		return nil, err
	}
	wf.modelTable.items, err = SelectMarkedDetails(db, []int64{})
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorTableInit)
		return nil, err
	}
	_, wf.modelTable.MapIdToEntityType, err = SelectIdTitle(db, "EntityType")
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorTypeInit)
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
		line := m.Map3.MapIdToMarkingLine[item.Parent.Marking]
		if line == nil {
			return "Нет"
		}
		e := m.Map3.MapIdToEntity[line.Hierarchy[len(line.Hierarchy)-1].Id]
		if e == nil {
			log.Println(data.S.Warning, "Обращение к отсутствующему значению")
			return "ERROR"
		}
		return fmt.Sprintf("%s %s %s", m.MapIdToEntityType[e.Type], e.Title, item.Mark)
	}
	log.Println(data.S.Panic, data.S.ErrorUnexpectedColumn)
	panic(data.S.ErrorUnexpectedColumn)
}

// Описание и запуск диалогового окна.
func MarkedDetailsRunDialog(owner walk.Form, db *sql.DB, isChange bool, parent *int64) (int, error) {
	log.Printf(data.S.BeginWindow, data.S.MarkedDetails)
	var err error
	var databind *walk.DataBinder
	var search struct {
		Order, Detail, Line int64
	}
	wf, err := newWindowsFormMarkedDetails(db)
	if err != nil {
		return 0, errors.Wrap(err, data.S.ErrorInit)
	}
	log.Printf(data.S.InitWindow, data.S.MarkedDetails)
	if err = (dec.Dialog{
		AssignTo: &wf.Dialog,
		Title:    data.S.HeadingMarkedDetails,
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
						Value:         dec.Bind("Order", dec.SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Title",
						MinSize:       dec.Size{120, 0},
						Model:         wf.modelTable.Map3.Orders(),
					},
					dec.HSpacer{Size: 20},

					dec.Label{
						Text: "Деталь:",
					},
					dec.ComboBox{
						Value:         dec.Bind("Detail", dec.SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Title",
						MinSize:       dec.Size{120, 0},
						Model:         wf.modelTable.Map3.MarkedDetails(),
					},
					dec.HSpacer{Size: 20},

					dec.PushButton{
						Text: data.S.ButtonSearch,
						OnClicked: func() {
							log.Println(data.S.Info, data.S.LogSearch)
							if err := databind.Submit(); err != nil {
								err = errors.Wrap(err, data.S.ErrorSubmit)
								log.Println(data.S.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
								return
							}
							lastLen := wf.modelTable.RowCount()
							markings := make([]int64, 0, 20)
							if search.Line != 0 {
								markings = append(markings, search.Line)
							} else {
								for _, val := range wf.modelTable.MapIdToMarkingLine {
									if val.Hierarchy[0].Id != search.Order && search.Order != 0 {
										continue
									}
									if val.Hierarchy[len(val.Hierarchy)-1].Id != search.Order && search.Detail != 0 {
										continue
									}
									markings = append(markings, val.Id)
								}
							}
							if wf.modelTable.items, err = SelectMarkedDetails(db, markings); err != nil {
								err = errors.Wrap(err, data.S.ErrorSubquery)
								log.Println(data.S.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
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
			dec.Composite{
				Layout: dec.HBox{MarginsZero: true},
				Children: []dec.Widget{
					dec.Label{
						Text: "Иерархия:",
					},
					dec.ComboBox{
						Value:         dec.Bind("Line", dec.SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Title",
						MinSize:       dec.Size{150, 0},
						Model:         wf.modelTable.Map3.MarkingLines(),
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
				Visible: isChange,
				Children: []dec.Widget{
					dec.PushButton{
						Text: data.S.ButtonAdd,
						OnClicked: func() {
							log.Println(data.S.Info, data.S.LogAdd)
							if err := wf.add(db); err != nil {
								err = errors.Wrap(err, data.S.ErrorAddRow)
								log.Println(data.S.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
							}
						},
					},
					dec.PushButton{
						Text: data.S.ButtonChange,
						OnClicked: func() {
							log.Println(data.S.Info, data.S.LogChange)
							if err := wf.change(db); err != nil {
								err = errors.Wrap(err, data.S.ErrorChangeRow)
								log.Println(data.S.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
							}
						},
					},
					dec.PushButton{
						Text: data.S.ButtonDelete,
						OnClicked: func() {
							log.Println(data.S.Info, data.S.LogDelete)
							if err := wf.delete(db); err != nil {
								err = errors.Wrap(err, data.S.ErrorDeleteRow)
								log.Println(data.S.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
							}
						},
					},
				},
			},
			dec.Composite{
				Layout:  dec.HBox{MarginsZero: true},
				Visible: !isChange,
				Children: []dec.Widget{
					dec.PushButton{
						Text: data.S.ButtonOK,
						OnClicked: func() {
							log.Println(data.S.Info, data.S.LogOk)
							if wf.modelTable.RowCount() > 0 && wf.tv.CurrentIndex() != -1 {
								index := wf.tv.CurrentIndex()
								*parent = wf.modelTable.items[index].Id
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
	log.Printf(data.S.CreateWindow, data.S.MarkedDetails)

	log.Printf(data.S.RunWindow, data.S.MarkedDetails)
	return wf.Run(), nil
}

// Функция, для добавления строки в таблицу.
func (wf windowsFormMarkedDetails) add(db *sql.DB) error {
	var detail MarkedDetail
	cmd, err := MarkedDetailRunDialog(wf, db, &wf.modelTable.Map3, false, &detail)
	log.Printf(data.S.EndWindow, data.S.Entity, cmd)
	if err != nil {
		return errors.Wrapf(err, data.S.InMarkedDetailRunDialog, false, detail)
	}
	if cmd != walk.DlgCmdOK {
		return nil
	}
	QwStr := data.InsertMarkedDetail(detail.Marking, detail.Parent.Id, detail.Mark)
	if err = db.Ping(); err != nil {
		return errors.Wrap(err, data.S.ErrorPingDB)
	}

	result, err := db.Exec(QwStr)
	if err != nil {
		return errors.Wrap(err, data.S.ErrorAddDB+QwStr)
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
		log.Println(data.S.Error, data.S.ErrorInsertIndexLog)
		walk.MsgBox(wf, data.S.MsgBoxError, data.S.ErrorInsertIndex, data.Icon.Critical)
		return nil
	}
	wf.modelTable.items[index].Id = id
	return nil
}

// Функция, для изменения строки в таблице.
func (wf windowsFormMarkedDetails) change(db *sql.DB) error {
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 {
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChooseRow, data.Icon.Info)
		return nil
	}
	var err error
	index := wf.tv.CurrentIndex()
	detail := wf.modelTable.items[index]
	cmd, err := MarkedDetailRunDialog(wf, db, &wf.modelTable.Map3, true, detail)
	log.Printf(data.S.EndWindow, data.S.MarkedDetail, cmd)

	if err != nil {
		return errors.Wrapf(err, data.S.InMarkedDetailRunDialog, true, *detail)
	}
	if cmd != walk.DlgCmdOK {
		return nil
	}
	QwStr := data.UpdateMarkedDetail(detail.Id, detail.Marking, detail.Parent.Id, detail.Mark)
	if err = db.Ping(); err != nil {
		return errors.Wrap(err, data.S.ErrorPingDB)
	}
	_, err = db.Exec(QwStr)
	if err != nil {
		return errors.Wrap(err, data.S.ErrorChangeDB+QwStr)
	}
	wf.modelTable.items[index] = detail
	wf.modelTable.PublishRowsReset()
	wf.modelTable.PublishRowChanged(index)
	return nil
}

// Функция, для удаления строки из таблицы.
func (wf windowsFormMarkedDetails) delete(db *sql.DB) error {
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 {
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChooseRow, data.Icon.Info)
		return nil
	}
	index := wf.tv.CurrentIndex()
	id := wf.modelTable.items[index].Id
	QwStr := data.DeleteMarkedDetail(id)
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
	wf.modelTable.PublishRowsChanged(index, wf.modelTable.RowCount()-1)
	return nil
}