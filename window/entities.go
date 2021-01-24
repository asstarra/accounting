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

func SelectEntities(db *sql.DB, title string, entityType int64) ([]*Entity, error) {
	arr := make([]*Entity, 0)
	if err := (func() error {
		if err := db.Ping(); err != nil {
			err = errors.Wrap(err, data.S.ErrorPingDB)
			return err
		}
		table := data.Tab["Entity"]
		sId := table.Columns["Id"].Name
		sTitle := table.Columns["Title"].Name
		sType := table.Columns["Type"].Name
		sSpec := table.Columns["Specification"].Name
		sProdLine := table.Columns["ProductionLine"].Name
		sNote := table.Columns["Note"].Name
		QwStr := fmt.Sprintf("SELECT %s AS id, %s AS title, %s AS type, %s AS spec, %s AS pline, %s AS note FROM %s",
			sId, sTitle, sType, sSpec, sProdLine, sNote, table.Name)
		//GO-TO выборка строчек по условию
		rows, err := db.Query(QwStr)
		if err != nil {
			return errors.Wrap(err, data.S.ErrorQuery+QwStr)
		}
		defer rows.Close()
		for rows.Next() {
			row := Entity{}
			err := rows.Scan(&row.Id, &row.Title, &row.Type, &row.Specification, &row.ProductionLine, &row.Note)
			if err != nil {
				err = errors.Wrap(err, data.S.ErrorDecryptRow)
				return err
			}
			arr = append(arr, &row)
		}
		return nil
	}()); err != nil {
		return arr, errors.Wrap(err, fmt.Sprintf("In SelectEntities(title = \"%s\", entityType = %d)", title, entityType))
	}
	return arr, nil
}

type modelEntitiesComponent struct {
	walk.TableModelBase
	items      []*Entity
	mapIdTitle map[int64]string
}

type windowsFormEntities struct {
	*walk.Dialog
	modelType  []*IdTitle
	modelTable *modelEntitiesComponent
	tv         *walk.TableView
}

func newWindowsFormEntities(db *sql.DB) (*windowsFormEntities, error) {
	var err error
	wf := new(windowsFormEntities)
	wf.modelTable, err = newModelEntitiesComponent(db)
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorTableInit)
		return nil, err
	}
	wf.modelType, wf.modelTable.mapIdTitle, err = SelectIdTitle(db, "EntityType")
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorTypeInit)
		return nil, err
	}
	return wf, nil
}

func newModelEntitiesComponent(db *sql.DB) (*modelEntitiesComponent, error) {
	var err error
	m := new(modelEntitiesComponent)
	m.items, err = SelectEntities(db, "", 0)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *modelEntitiesComponent) RowCount() int {
	return len(m.items)
}

func (m *modelEntitiesComponent) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return m.mapIdTitle[item.Type]
	case 1:
		return item.Title
	case 2:
		return item.Specification
	case 3:
		return item.Note
	}
	log.Println(data.S.Panic, data.S.ErrorUnexpectedColumn)
	panic(data.S.ErrorUnexpectedColumn)
}

func EntitiesRunDialog(owner walk.Form, db *sql.DB, isChange bool, idTitle *IdTitle) (int, error) {
	log.Printf(data.S.BeginWindow, data.S.Entities)
	var err error
	wf, err := newWindowsFormEntities(db)
	if err != nil {
		return 0, errors.Wrap(err, "In newWindowsFormEntities()")
	}

	if err = (dec.Dialog{
		AssignTo: &wf.Dialog,
		Title:    data.S.HeadingEntities,
		// Size:     dec.Size{100, 100},
		Layout:  dec.VBox{},
		MinSize: dec.Size{450, 0},
		Children: []dec.Widget{
			// GO-TO
			dec.TableView{
				AssignTo: &wf.tv,
				Columns: []dec.TableViewColumn{
					{Title: "Тип"},
					{Title: "Название"},
					{Title: "Спецификация"},
					{Title: "Примичание"},
				},
				MinSize: dec.Size{0, 200},
				Model:   wf.modelTable,
				OnCurrentIndexChanged: func() { //GO-TO
					fmt.Printf("CurrentIndexes: %v\n", wf.tv.CurrentIndex())
				},
				OnSelectedIndexesChanged: func() {
					fmt.Printf("SelectedIndexes: %v\n", wf.tv.SelectedIndexes())
				},
			},
			dec.Composite{
				Layout:  dec.HBox{},
				Visible: isChange,
				Children: []dec.Widget{
					dec.PushButton{
						Text: data.S.ButtonAdd,
						OnClicked: func() {
							log.Println(data.S.Info, data.S.LogAdd)
							if err := wf.add(db); err != nil {
								err = errors.Wrap(err, data.S.ErrorAdd)
								log.Println(data.S.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, err.Error(), data.Icon.Error)
							}
						},
					},
					dec.PushButton{
						Text: data.S.ButtonChange,
						OnClicked: func() {
							log.Println(data.S.Info, data.S.LogChange)
							if err := wf.change(db); err != nil {
								err = errors.Wrap(err, data.S.ErrorChange)
								log.Println(data.S.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, err.Error(), data.Icon.Error)
							}
						},
					},
					dec.PushButton{
						Text: data.S.ButtonDelete,
						OnClicked: func() {
							log.Println(data.S.Info, data.S.LogDelete)
							if err := wf.delete(db); err != nil {
								err = errors.Wrap(err, data.S.ErrorDelete)
								log.Println(data.S.Error, err)
								walk.MsgBox(wf, data.S.MsgBoxError, err.Error(), data.Icon.Error)
							}
						},
					},
				},
			},
			dec.Composite{
				Layout:  dec.HBox{},
				Visible: !isChange,
				Children: []dec.Widget{
					dec.PushButton{
						Text: data.S.ButtonOK,
						OnClicked: func() {
							log.Println(data.S.Info, data.S.LogOk)
							if wf.modelTable.RowCount() > 0 && wf.tv.CurrentIndex() != -1 {
								fmt.Println(wf.tv.CurrentIndex())
								index := wf.tv.CurrentIndex()
								idTitle.Id = wf.modelTable.items[index].Id
								idType := wf.modelTable.items[index].Type
								sType := wf.modelTable.mapIdTitle[idType]
								idTitle.Title = sType + " " + wf.modelTable.items[index].Title
								wf.Accept()
							} else {
								walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChoseRow, data.Icon.Info)
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
	log.Printf(data.S.CreateWindow, data.S.Entities)

	log.Printf(data.S.RunWindow, data.S.Entities)
	return wf.Run(), nil
}

func (wf windowsFormEntities) add(db *sql.DB) error {
	var entity Entity
	children := make([]*EntityRecChild, 0)
	entity.Children = &children
	cmd, err := EntityRunDialog(wf, db, &entity)
	log.Printf(data.S.EndWindow, data.S.Entity, cmd)
	if err != nil {
		return errors.Wrapf(err, "In EntityRunDialog(entity = %v)", entity) //GO-TO
	}
	if cmd != walk.DlgCmdOK {
		return nil
	}
	if err = db.Ping(); err != nil {
		return errors.Wrap(err, data.S.ErrorPingDB)
	}
	table := data.Tab["Entity"]
	sTitle := table.Columns["Title"].Name
	sType := table.Columns["Type"].Name
	sSpec := table.Columns["Specification"].Name
	sProdLine := table.Columns["ProductionLine"].Name
	sNote := table.Columns["Note"].Name
	QwStr := fmt.Sprintf("INSERT %s (%s, %s, %s, %s, %s) VALUES (%s, %d, %s, %t, %s)",
		table.Name, sTitle, sType, sSpec, sProdLine, sNote,
		Empty(entity.Title), entity.Type, Empty(entity.Specification), entity.ProductionLine, Empty(entity.Note))
	fmt.Println(QwStr)
	result, err := db.Exec(QwStr)
	if err != nil {
		return errors.Wrap(err, data.S.ErrorAddDB+QwStr)
	}

	trackLatest := wf.tv.ItemVisible(len(wf.modelTable.items) - 1) //&& len(wf.tv.SelectedIndexes()) <= 1
	wf.modelTable.items = append(wf.modelTable.items, &entity)
	index := len(wf.modelTable.items) - 1
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
	table = data.Tab["EntityRec"]
	sIdC := table.Columns["Child"].Name
	sIdP := table.Columns["Parent"].Name
	sCount := table.Columns["Count"].Name
	for _, v := range *entity.Children {
		QwStrChild := fmt.Sprintf("INSERT %s (%s, %s, %s) VALUES (%d, %d, %d)",
			table.Name, sIdP, sIdC, sCount, id, v.Id, v.Count)
		_, err := db.Exec(QwStrChild)
		if err != nil {
			err = errors.Wrap(err, data.S.ErrorAddDB+QwStrChild)
			log.Println(data.S.Error, err)
			walk.MsgBox(wf, data.S.MsgBoxError, err.Error(), data.Icon.Critical)
		}
	}
	return nil
}

func (wf windowsFormEntities) change(db *sql.DB) error {
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 {
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChoseRow, data.Icon.Info)
		return nil
	}
	index := wf.tv.CurrentIndex()
	entity := *wf.modelTable.items[index]
	//GO-TO
	children := make([]*EntityRecChild, 0)
	entity.Children = &children
	cmd, err := EntityRunDialog(wf, db, &entity)
	log.Printf(data.S.EndWindow, data.S.Entity, cmd)
	if err != nil {
		return errors.Wrapf(err, "In EntityRunDialog(entity = %v)", entity) //GO-TO
	}
	if cmd != walk.DlgCmdOK {
		return nil
	}
	if err = db.Ping(); err != nil {
		return errors.Wrap(err, data.S.ErrorPingDB)
	}
	table := data.Tab["Entity"]
	sId := table.Columns["Id"].Name
	sTitle := table.Columns["Title"].Name
	sType := table.Columns["Type"].Name
	sSpec := table.Columns["Specification"].Name
	sProdLine := table.Columns["ProductionLine"].Name
	sNote := table.Columns["Note"].Name
	QwStr := fmt.Sprintf("UPDATE %s SET %s = %s, %s = %d, %s = %s, %s = %t, %s = %s WHERE %s = %d",
		table.Name, sTitle, Empty(entity.Title), sType, entity.Type, sSpec, Empty(entity.Specification),
		sProdLine, entity.ProductionLine, sNote, Empty(entity.Note), sId, entity.Id)
	fmt.Println(QwStr)
	_, err = db.Exec(QwStr)
	if err != nil {
		return errors.Wrap(err, data.S.ErrorChangeDB+QwStr)
	}

	wf.modelTable.items[index] = &entity
	wf.modelTable.PublishRowChanged(index)
	return nil
}

func (wf windowsFormEntities) delete(db *sql.DB) error {
	if wf.modelTable.RowCount() <= 0 || wf.tv.CurrentIndex() == -1 {
		walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChoseRow, data.Icon.Info)
		return nil
	}
	// if wf.modelTable.RowCount() > 0 && wf.tv.CurrentIndex() != -1 {
	index := wf.tv.CurrentIndex()
	log.Println(index, *wf.modelTable.items[index])
	id := wf.modelTable.items[index].Id
	table := data.Tab["Entity"]
	sId := table.Columns["Id"].Name
	QwStr := fmt.Sprintf("DELETE FROM %s WHERE %s = %d", table.Name, sId, id)
	if err := db.Ping(); err != nil {
		return errors.Wrap(err, data.S.ErrorPingDB)
	}
	_, err := db.Exec(QwStr)
	if err != nil {
		return errors.Wrap(err, data.S.ErrorDeleteDB+QwStr)
	}

	trackLatest := wf.tv.ItemVisible(len(wf.modelTable.items) - 1) //&& len(wf.tv.SelectedIndexes()) <= 1
	wf.modelTable.items = wf.modelTable.items[:index+copy(wf.modelTable.items[index:], wf.modelTable.items[index+1:])]
	wf.modelTable.PublishRowsRemoved(index, index)
	if trackLatest {
		wf.tv.EnsureItemVisible(len(wf.modelTable.items) - 1)
	}
	// if l := len(wf.modelTable.items); l <= index {
	// 	index = l - 1
	// }
	// if index >= 0 {
	// 	wf.tv.SetCurrentIndex(index)
	// }
	// } else {
	// 	walk.MsgBox(wf, data.S.MsgBoxInfo, data.S.MsgChoseRow, data.Icon.Info)
	// }
	return nil
}
