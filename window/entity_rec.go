package window

import (
	"accounting/data"
	"accounting/data/qwery"

	// "accounting/window"
	"database/sql"
	"fmt"
	"log"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

type Id64Title struct {
	Id    int64  // Ид.
	Title string // Описание. Название типа + название сущности.
}

func (it Id64Title) String() string {
	return fmt.Sprintf("{Id = %d, Title = '%s'}", it.Id, it.Title)
}

type EntityRecChild struct { // GO-TO rename.
	Id64Title       // Ид + текстовое описание.
	Count     int32 // Количество дочерних компонентов.
}

type EntityRec struct {
	IdP            int64 // Ид родителя.
	EntityRecChild       // Ид + описание + количество.
}

func (it EntityRecChild) String() string {
	return fmt.Sprintf("'%s'", it.Title)
	// return fmt.Sprintf("{Id = %d, Title = '%s', Count = %d}", it.Id, it.Title, it.Count) //GO-TO
}

// Поиск дочерних сущностей для заданной сущности (id родительской таблицы = parent).
// Выборка из таблицы EntityRec id дочерней таблицы и количества.
// Поле title = тип сущности + название сущности.
func SelectEntityRecChild(db *sql.DB, parent *int64) ([]*EntityRec, []*EntityRecChild, error) {
	arr := make([]*EntityRec, 0, 1000)
	if err := (func() error {
		QwStr := qwery.SelectEntityRecChild(parent)
		if err := db.Ping(); err != nil { // Пинг БД.
			return errors.Wrap(err, data.S.ErrorPingDB)
		}
		rows, err := db.Query(QwStr) // Запрос к БД.
		if err != nil {
			return errors.Wrapf(err, data.S.ErrorQueryDB, QwStr)
		}
		defer rows.Close()
		var e_type, title string
		for rows.Next() {
			row := EntityRec{}
			err := rows.Scan(&row.IdP, &row.Id, &e_type, &title, &row.Count)
			if err != nil {
				return errors.Wrap(err, data.S.ErrorDecryptRow)
			}
			row.Title = e_type + " " + title
			arr = append(arr, &row)
		}
		return nil
	}()); err != nil {
		return arr, []*EntityRecChild{}, errors.Wrapf(err, data.Log.InSelectEntityRecChild, parent)
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
	log.Printf(data.Log.BeginWindow, data.Log.EntityRec) // Лог.
	var databind *walk.DataBinder                        // Иициализация.
	wf := windowsFormEntityRec{}
	log.Printf(data.Log.InitWindow, data.Log.EntityRec) // Лог.
	if err := (dec.Dialog{                              // Описание окна.
		AssignTo: &wf.Dialog,              // Привязка окна.
		Title:    data.S.HeadingEntityRec, // Название.
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
							log.Println(data.Log.Info, data.Log.LogChoose)
							if err := (func() error {
								log.Println(data.Log.Info, "Выбор") // Лог. GO-TO
								idTitle := child.Id64Title
								// var it = window.IdTitle{Id: idTitle.Id, Title: idTitle.Title} //GO-TO
								cmd, err := EntitiesRunDialog(wf, db, false, &idTitle)
								log.Printf(data.Log.EndWindow, data.Log.Entities, cmd) // Лог.
								if err != nil {
									return errors.Wrapf(err, data.Log.InEntitiesRunDialog, false, idTitle)
								}
								if cmd == walk.DlgCmdOK {
									child.Id64Title = idTitle
									wf.entitiesBW.SetText(child.Title)
								}
								return nil
							}()); err != nil {
								err = errors.Wrap(err, data.S.ErrorChoose) // Обработка ошибки.
								log.Println(data.Log.Error, err)           // Лог.
								walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
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
						Text: data.S.ButtonOK,
						OnClicked: func() {
							log.Println(data.Log.Info, data.Log.LogOk) // Лог.
							if err := databind.Submit(); err != nil {
								err = errors.Wrap(err, data.S.ErrorSubmit) // Обработка ошибки.
								log.Println(data.Log.Error, err)           // Лог.
								walk.MsgBox(wf, data.S.MsgBoxError, MsgError(err), data.Icon.Error)
								return
							}
							wf.Accept()
						},
					},
					dec.PushButton{ // Кнопка Отмена.
						Text: data.S.ButtonCansel,
						OnClicked: func() {
							log.Println(data.Log.Info, data.Log.LogCansel) // Лог.
							wf.Cancel()
						},
					},
				},
			},
		},
	}.Create(owner)); err != nil {
		err = errors.Wrap(err, data.S.ErrorCreateWindow) // Обработка ошибки создания окна.
		return 0, err
	}
	log.Printf(data.Log.CreateWindow, data.Log.EntityRec) // Лог.

	log.Printf(data.Log.RunWindow, data.Log.EntityRec) // Лог.
	return wf.Run(), nil                               // Запуск окна.
}
