package main

import (
	"database/sql"
	"fmt"
	"time"

	// "io"
	"log"
	"os"

	"accounting/data"
	"accounting/data/db"
	e "accounting/data/errors"
	l "accounting/data/log"
	"accounting/data/qwery"
	. "accounting/data/table"
	"accounting/optimization"
	"accounting/window"
	// "accounting/window2"
)

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

func main() {
	sErr := "Error: "
	f, err := os.OpenFile("logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		err = errors.Wrap(err, e.Err.ErrorOpenFile+"logfile.txt")
		window.ErrorRunWindow(sErr + err.Error())
		log.Println(l.Error, err)
		return
	}
	defer f.Close()
	log.SetOutput(f)

	log.Println("-------------------------------------------------------------")

	err = data.Init()
	if err != nil {
		err = errors.Wrap(err, e.Err.ErrorInit)
		window.ErrorRunWindow(sErr + err.Error())
		log.Println(l.Error, err)
		return
	}

	db, err := sql.Open("mysql", db.DataSourseTcp())
	if err != nil {
		err = errors.Wrap(err, e.Err.ErrorOpedDB)
		window.ErrorRunWindow(sErr + err.Error())
		log.Println(l.Error, err)
		return
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		err = errors.Wrap(err, e.Err.ErrorPingDB)
		window.ErrorRunWindow(sErr + err.Error())
		log.Println(l.Error, err)
		return
	}

	// optimization.A()
	// if opt, err := optimization.NewQualificationTable(db); err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(opt)
	// }
	// if opt, err := optimization.SelectPerson(db, nil, nil); err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(opt)
	// }
	// if opt, err := optimization.SelectDetail(db, nil, nil); err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(opt)
	// }
	// if opt, err := optimization.SelectDays(db, nil, nil); err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(opt)
	// }
	// var o optimization.Optimization
	// o.Init(db, nil, nil)
	// o.Run()
	var a = time.Now()
	fmt.Println(qwery.ToStr(a), qwery.ToStr(&a))
	var b int32 = 7
	var c int8 = 8
	fmt.Println(qwery.ToStr(&b), qwery.ToStr(b), qwery.ToStr(&c))
	fmt.Println(qwery.Sprintf("%s, %s, %s, %s", &a, &b, 874.3, c))

	var mw *walk.MainWindow
	if err := (dec.MainWindow{
		AssignTo: &mw,
		Title:    "Учет",
		Size:     dec.Size{300, 80},
		Layout:   dec.VBox{},
		Children: []dec.Widget{
			// dec.PushButton{
			// 	Text: "Тип компонента",
			// 	OnClicked: func() {
			// 		cmd, err := window2.TypeRunDialog(mw, db, "EntityType")
			// 		if err != nil {
			// 			log.Println("ERROR!", err)
			// 		}
			// 		log.Printf(l.EndWindow, l.Type, cmd)
			// 	},
			// },
			// dec.PushButton{
			// 	Text: "Компоненты",
			// 	OnClicked: func() {
			// 		cmd, err := window2.EntitiesRunDialog(mw, db, true, nil)
			// 		if err != nil {
			// 			log.Println("ERROR!", err)
			// 		}
			// 		log.Printf(l.EndWindow, l.Entities, cmd)
			// 	},
			// },

			dec.PushButton{
				Text: "EntityType",
				OnClicked: func() {
					cmd, err := window.TypeRunDialog(mw, db, TableEntityType)
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Printf(l.EndWindow, l.Type, cmd)
				},
			},
			dec.PushButton{
				Text: "StatusType",
				OnClicked: func() {
					cmd, err := window.TypeRunDialog(mw, db, TableStatusType)
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Printf(l.EndWindow, l.Type, cmd)
				},
			},
			dec.PushButton{
				Text: "Person",
				OnClicked: func() {
					cmd, err := window.TypeRunDialog(mw, db, TablePerson)
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Printf(l.EndWindow, l.Type, cmd)
				},
			},
			dec.PushButton{
				Text: "Operation",
				OnClicked: func() {
					cmd, err := window.TypeRunDialog(mw, db, TableOperation)
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Printf(l.EndWindow, l.Type, cmd)
				},
			},
			dec.PushButton{
				Text: "Entities",
				OnClicked: func() {
					cmd, err := window.EntitiesRunDialog(mw, db, true, nil)
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Printf(l.EndWindow, l.Entities, cmd)
				},
			},
			dec.PushButton{
				Text: "MarkedDetail",
				OnClicked: func() {
					cmd, err := window.MarkedDetailsRunDialog(mw, db, true, nil)
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Printf(l.EndWindow, l.MarkedDetails, cmd)
				},
			},
			dec.PushButton{
				Text: "Маршрутка",
				OnClicked: func() {
					var o optimization.Optimization
					o.Init(db, nil, nil)
					o.Run()
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Println(err)
				},
			},
		},
	}.Create()); err != nil {
		err = errors.Wrap(err, e.Err.ErrorCreateWindow)
		window.ErrorRunWindow(sErr + err.Error())
		log.Println(l.Error, err)
		return
	}
	mw.Run()
}
