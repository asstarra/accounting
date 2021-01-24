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

func SelectIdTitle(db *sql.DB, tableName string) ([]*IdTitle, map[int64]string, error) {
	arr := make([]*IdTitle, 0)
	m := make(map[int64]string)
	if err := (func() error {
		if err := db.Ping(); err != nil {
			err = errors.Wrap(err, data.S.ErrorPingDB)
			return err
		}
		table := data.Tab[tableName]
		sId := table.Columns["Id"].Name
		sTitle := table.Columns["Title"].Name
		QwStr := fmt.Sprintf("SELECT %s AS id, %s AS title FROM %s", sId, sTitle, table.Name)
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorQuery+QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			row := IdTitle{}
			err := rows.Scan(&row.Id, &row.Title)
			if err != nil {
				err = errors.Wrap(err, data.S.ErrorDecryptRow)
				return err
			}
			arr = append(arr, &row)
			m[row.Id] = row.Title
		}
		return nil
	}()); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("In SelectIdTitle(tableName = %s)", tableName))
	}
	return arr, m, nil
}

type itemType struct {
	id   int64
	last string
	now  string
}

type modelType struct {
	walk.ListModelBase
	items []itemType
}

type windowsFormType struct {
	*walk.Dialog
	lbModel         *modelType
	lbWidget        *walk.ListBox
	textToAddWidget *walk.LineEdit
}

func itemAdd_EntityType(s string) itemType {
	return itemType{
		id:   0,
		last: "",
		now:  s,
	}
}

func (item *itemType) update(now string) {
	item.now = now
}

func newModelType(tableName string, db *sql.DB) (*modelType, error) {
	arr, _, err := SelectIdTitle(db, tableName)
	if err != nil {
		err = errors.Wrap(err, "In newModelType произошла ошибка запроса данных")
		log.Println("ERROR!", err)
		return nil, err
	}
	m := &modelType{items: make([]itemType, len(arr))}
	for i, v := range arr {
		m.items[i].id = v.Id
		m.items[i].last = v.Title
		m.items[i].now = v.Title
	}
	return m, nil
}

func (m *modelType) ItemCount() int {
	return len(m.items)
}

func (m *modelType) Value(index int) interface{} {
	return m.items[index].now
}

func (wf *windowsFormType) save(tableName string, db *sql.DB) error {
	if err := db.Ping(); err != nil {
		err = errors.Wrap(err, "Не удалось подключиться к базе данных")
		log.Println("ERROR!", err)
		return err
	}
	table := data.Tab[tableName]
	title := table.Columns["Title"].Name
	id := table.Columns["Id"].Name
	QwStrIns := "INSERT INTO " + table.Name + " (" + title + ") VALUES (?)"
	QwStrUpd := "UPDATE " + table.Name + " SET " + title + " = ? WHERE title = ?"
	QwStrDel := "DELETE FROM " + table.Name + " WHERE " + id + " NOT IN ("
	pointFlag := false
	for _, v := range wf.lbModel.items {
		if v.id != 0 {
			if pointFlag {
				QwStrDel += ", "
			}
			QwStrDel += strconv.FormatInt(v.id, 64) //strconv.Itoa(v.id)
			pointFlag = true
		}
	}
	QwStrDel += ")"

	tx, err := db.Begin()
	if err != nil {
		err = errors.Wrap(err, "Не удалось начать транзакцию")
		log.Println("ERROR!", err)
		return err
	}
	_, err = tx.Exec(QwStrDel)
	if err != nil {
		err = errors.Wrap(err, "Не удалось удалить строчки")
		log.Println("ERROR!", err)
		if errBack := tx.Rollback(); errBack != nil {
			errBack = errors.Wrap(errBack, "Не удалось откатить транзакцию")
			err := errors.Wrap(errBack, err.Error()+"\n")
			log.Println("ERROR!", err)
		}
		return err
	}
	for i, v := range wf.lbModel.items {
		if v.last != v.now {
			if v.last == "" {
				result, err := tx.Exec(QwStrIns, v.now)
				if err != nil {
					err = errors.Wrap(err, "Не удалось вставить строчку со значением "+v.now)
					log.Println("ERROR!", err)
					if errBack := tx.Rollback(); errBack != nil {
						errBack = errors.Wrap(errBack, "Не удалось откатить транзакцию")
						err := errors.Wrap(errBack, err.Error()+"\n")
						log.Println("ERROR!", err)
					}
					return err
				} else {
					wf.lbModel.items[i].last = v.now
					if id, err := result.LastInsertId(); err != nil {
						err = errors.Wrap(err, "Не удалось узнать индекс вставленной строчки")
						err = errors.Wrap(err, "При успешном сохранении, рекомендуется закрыть вкладку и открыть ее заново, иначе работа программы не гарантируется\n")
						log.Println("ERROR!", err)
						walk.MsgBox(wf, "Внимание, возможна некорректная работа программы", err.Error(), walk.MsgBoxIconWarning)
					} else {
						wf.lbModel.items[i].id = id
					}
				}
			} else {
				_, err := tx.Exec(QwStrUpd, v.now, v.last)
				if err != nil {
					err = errors.Wrap(err, "Не удалось обновить значение "+v.last+" на значение "+v.now)
					log.Println("ERROR!", err)
					if errBack := tx.Rollback(); errBack != nil {
						errBack = errors.Wrap(errBack, "Не удалось откатить транзакцию")
						err := errors.Wrap(errBack, err.Error()+"\n")
						log.Println("ERROR!", err)
					}
					return err
				} else {
					wf.lbModel.items[i].last = v.now
				}
			}
		}
	}
	if err = tx.Commit(); err != nil {
		err = errors.Wrap(err, "Не удалось завершить транзакцию")
		log.Println("ERROR!", err)
		return err
	}
	return nil
}

func TypeDialogRun(owner walk.Form, sTitle, tableName string, db *sql.DB) (int, error) {
	log.Println("INFO -", "BEGIN window - TYPE, tableName -", tableName)

	sHeading := "Внимание"
	msgBoxIcon := walk.MsgBoxIconWarning
	// var acceptPB, cancelPB *walk.PushButton
	lbModel, err := newModelType(tableName, db)
	if err != nil {
		err = errors.Wrap(err, "При заполнении списка произошла ошибка")
		log.Println("ERROR!", err)
		return 0, err
	}
	wf := &windowsFormType{
		lbModel: lbModel,
	}
	data.Reg.IsSaveDialog.SetSatisfied(false)
	checkTitle := func(text string) bool {
		if text == "" {
			walk.MsgBox(wf, sHeading, "Название не должно состоять из пустой строки.", msgBoxIcon)
			return false
		}
		if len(text) > 255 {
			walk.MsgBox(wf, sHeading, "Длинна названия должна быть меньше 255.", msgBoxIcon)
			return false
		}
		for i := range wf.lbModel.items {
			if wf.lbModel.Value(i) == text {
				walk.MsgBox(wf, sHeading, "Такое название уже существует.", msgBoxIcon)
				return false
			}
		}
		return true
	}
	checkIndex := func() bool {
		if wf.lbWidget.CurrentIndex() < 0 {
			walk.MsgBox(wf, sHeading, "Выберите строчку, которую хотите изменить.", msgBoxIcon)
			return false
		}
		return true
	}

	if err := (dec.Dialog{
		AssignTo: &wf.Dialog,
		Title:    sTitle,
		// DefaultButton: &acceptPB,
		// CancelButton:  &cancelPB,
		// MinSize: dec.Size{400, 250},
		Layout: dec.VBox{MarginsZero: true},
		Children: []dec.Widget{
			dec.HSplitter{
				Children: []dec.Widget{
					dec.Composite{
						MinSize: dec.Size{150, 150},
						Layout:  dec.VBox{},
						Children: []dec.Widget{
							dec.ListBox{
								AssignTo: &wf.lbWidget,
								Model:    wf.lbModel,
							},
						},
					},
					dec.Composite{
						// MinSize: dec.Size{200, 50},
						Layout: dec.VBox{},
						Children: []dec.Widget{
							dec.LineEdit{
								AssignTo:  &wf.textToAddWidget,
								MaxLength: 63,
							},
							dec.PushButton{
								Text: "Добавить",
								OnClicked: func() {
									log.Println("INFO -", "Add")
									if text := wf.textToAddWidget.Text(); checkTitle(text) {
										trackLatest := wf.lbWidget.ItemVisible(len(wf.lbModel.items)-1) && len(wf.lbWidget.SelectedIndexes()) <= 1
										wf.lbModel.items = append(wf.lbModel.items, itemAdd_EntityType(wf.textToAddWidget.Text()))
										index := len(wf.lbModel.items) - 1
										wf.lbModel.PublishItemsInserted(index, index)
										if trackLatest {
											wf.lbWidget.EnsureItemVisible(len(wf.lbModel.items) - 1)
										}
										wf.textToAddWidget.SetText("")
										data.Reg.IsSaveDialog.SetSatisfied(true)
									}
									log.Println("INFO -", wf.lbModel.items)
								},
							},
							dec.PushButton{
								Text: "Переименовать",
								OnClicked: func() {
									log.Println("INFO -", "Update")
									if !checkIndex() {
										return
									}
									if text := wf.textToAddWidget.Text(); checkTitle(text) {
										index := wf.lbWidget.CurrentIndex()
										wf.lbModel.items[index].now = text
										wf.lbModel.PublishItemChanged(index)
										wf.textToAddWidget.SetText("")
										data.Reg.IsSaveDialog.SetSatisfied(true)
									}
									log.Println("INFO -", wf.lbModel.items)
								},
							},
							dec.PushButton{
								Text: "Удалить",
								OnClicked: func() {
									log.Println("INFO -", "Delete")
									if !checkIndex() {
										return
									}
									trackLatest := wf.lbWidget.ItemVisible(len(wf.lbModel.items)-1) && len(wf.lbWidget.SelectedIndexes()) <= 1
									index := wf.lbWidget.CurrentIndex()
									wf.lbModel.items = wf.lbModel.items[:index+copy(wf.lbModel.items[index:], wf.lbModel.items[index+1:])]
									wf.lbModel.PublishItemsRemoved(index, index)
									if l := len(wf.lbModel.items); l <= index {
										index = l - 1
									}
									if index >= 0 {
										wf.lbWidget.SetCurrentIndex(index)
									}
									if trackLatest {
										wf.lbWidget.EnsureItemVisible(len(wf.lbModel.items) - 1)
									}
									data.Reg.IsSaveDialog.SetSatisfied(true)
									log.Println("INFO -", wf.lbModel.items)
								},
							},
							dec.Composite{
								Layout: dec.HBox{},
								Children: []dec.Widget{
									dec.PushButton{
										Text: "OK",
										// AssignTo: &acceptPB,
										Enabled: dec.Bind("!IsSaveDialog"),
										OnClicked: func() {
											log.Println("INFO -", "Ok")
											wf.Accept()
										},
									},
									dec.PushButton{
										Text:    "Сохранить",
										Enabled: dec.Bind("IsSaveDialog"),
										OnClicked: func() {
											log.Println("INFO -", "Save")
											if err := wf.save(tableName, db); err != nil {
												err = errors.Wrap(err, "In "+tableName+" не удалось сохранить данные.")
												log.Println("ERROR!", err)
												walk.MsgBox(wf, sHeading, err.Error(), walk.MsgBoxIconExclamation)
											} else {
												data.Reg.IsSaveDialog.SetSatisfied(false)
											}
											log.Println("INFO -", wf.lbModel.items)
										},
									},
									dec.PushButton{
										Text: "Отмена",
										// AssignTo: &cancelPB,
										Enabled: dec.Bind("IsSaveDialog"),
										OnClicked: func() {
											log.Println("INFO -", "Cancel")
											wf.Cancel()
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}.Create(owner)); err != nil {
		err = errors.Wrap(err, "Could not create Window Form Type")
		log.Println("ERROR!", err)
		return 0, err
	}
	log.Println("INFO -", "CREATE window - TYPE, tableName -", tableName)
	log.Println("INFO -", wf.lbModel.items)
	log.Println("INFO -", "RUN window - TYPE, tableName -", tableName)
	return wf.Run(), nil
}
