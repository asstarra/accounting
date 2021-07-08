package window

import (
	"accounting/data"
	e "accounting/data/errors"
	l "accounting/data/log"
	"accounting/data/qwery"
	"accounting/data/text"
	. "accounting/window/data"

	// "accounting/window"
	"database/sql"
	// "fmt"
	"log"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

// Поиск дочерних сущностей для заданной сущности (id родительской таблицы = parent).
// Выборка из таблицы EntityRec id дочерней таблицы и количества.
// Поле title = тип сущности + название сущности.
func SelectEntityRecChild(db *sql.DB, parent *int64) ([]*EntityRec, []*EntityRecChild, error) {
	arr := make([]*EntityRec, 0, 1000)
	if err := (func() error {
		QwStr := qwery.SelectEntityRecChild(parent)
		if err := db.Ping(); err != nil { // Пинг БД.
			return errors.Wrap(err, e.Err.ErrorPingDB)
		}
		rows, err := db.Query(QwStr) // Запрос к БД.
		if err != nil {
			return errors.Wrapf(err, e.Err.ErrorQueryDB, QwStr)
		}
		defer rows.Close()
		var e_type, title string
		for rows.Next() {
			row := EntityRec{}
			err := rows.Scan(&row.IdP, &row.Id, &e_type, &title, &row.Count)
			if err != nil {
				return errors.Wrap(err, e.Err.ErrorDecryptRow)
			}
			row.Title = e_type + " " + title
			arr = append(arr, &row)
		}
		return nil
	}()); err != nil {
		return arr, []*EntityRecChild{}, errors.Wrapf(err, l.In.InSelectEntityRecChild, parent)
	}
	childs := make([]*EntityRecChild, 0, len(arr))
	for _, val := range arr {
		childs = append(childs, &(val.EntityRecChild)) //GO-TO проверить совпадения адресов.
	}
	return arr, childs, nil
}

// Структура, содержащая описание и переменные окна.
type windowsFormEntityRec struct {
	*walk.Dialog
	entitiesBW *walk.PushButton // Виджет кнопки для выбора Entity.
}

// Описание и запуск диалогового окна.
func EntityRecRunDialog(owner walk.Form, db *sql.DB, isChange bool, child *EntityRecChild) (int, error) {
	log.Printf(l.BeginWindow, l.EntityRec) // Лог.
	var databind *walk.DataBinder          // Иициализация.
	wf := windowsFormEntityRec{}
	log.Printf(l.InitWindow, l.EntityRec) // Лог.
	if err := (dec.Dialog{                // Описание окна.
		AssignTo: &wf.Dialog,              // Привязка окна.
		Title:    text.T.HeadingEntityRec, // Название.
		DataBinder: dec.DataBinder{ // Привязка к структуре.
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
					dec.Label{ // Лэйбел название.
						Text: "Название:", //GO-TO
					},
					dec.PushButton{ // Кнопка для выбора дочерней "Entity"
						AssignTo: &wf.entitiesBW, // Привязка к виджету.
						Enabled:  !isChange,      // Доступ.
						MinSize:  dec.Size{150, 10},
						Text:     child.Title, // Текст.
						OnClicked: func() {
							log.Println(l.Info, l.LogChoose)
							if err := (func() error {
								log.Println(l.Info, "Выбор") // Лог. GO-TO
								idTitle := child.Id64Title
								// var it = window.IdTitle{Id: idTitle.Id, Title: idTitle.Title} //GO-TO
								cmd, err := EntitiesRunDialog(wf, db, false, &idTitle)
								log.Printf(l.EndWindow, l.Entities, cmd) // Лог.
								if err != nil {
									return errors.Wrapf(err, l.In.InEntitiesRunDialog, false, idTitle)
								}
								if cmd == walk.DlgCmdOK {
									child.Id64Title = idTitle
									wf.entitiesBW.SetText(child.Title)
								}
								return nil
							}()); err != nil {
								err = errors.Wrap(err, e.Err.ErrorChoose) // Обработка ошибки.
								log.Println(l.Error, err)                 // Лог.
								walk.MsgBox(wf, text.T.MsgBoxError, MsgError(err), data.Icon.Error)
							}
						},
					},

					dec.Label{ // Лэйбел количество.
						Text: "Количество:", //GO-TO
					},
					dec.NumberEdit{ // Числовое поле. Количество штук.
						Value:    dec.Bind("Count", dec.Range{0, 1000}),
						Suffix:   " шт",
						Decimals: 0,
					},
				},
			},
			dec.Composite{
				Layout: dec.HBox{},
				Children: []dec.Widget{
					dec.PushButton{ // Кнопка Ок.
						Text: text.T.ButtonOK,
						OnClicked: func() {
							log.Println(l.Info, l.LogOk) // Лог.
							if err := databind.Submit(); err != nil {
								err = errors.Wrap(err, e.Err.ErrorSubmit) // Обработка ошибки.
								log.Println(l.Error, err)                 // Лог.
								walk.MsgBox(wf, text.T.MsgBoxError, MsgError(err), data.Icon.Error)
								return
							}
							wf.Accept()
						},
					},
					dec.PushButton{ // Кнопка Отмена.
						Text: text.T.ButtonCansel,
						OnClicked: func() {
							log.Println(l.Info, l.LogCansel) // Лог.
							wf.Cancel()
						},
					},
				},
			},
		},
	}.Create(owner)); err != nil {
		err = errors.Wrap(err, e.Err.ErrorCreateWindow) // Обработка ошибки создания окна.
		return 0, err
	}
	log.Printf(l.CreateWindow, l.EntityRec) // Лог.

	log.Printf(l.RunWindow, l.EntityRec) // Лог.
	return wf.Run(), nil                 // Запуск окна.
}
