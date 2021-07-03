package main

import (
	"database/sql"
	// "fmt"

	// "io"
	"log"
	"os"

	"accounting/data"
	"accounting/optimization"
	"accounting/window"
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
		err = errors.Wrap(err, data.S.ErrorOpenFile+"logfile.txt")
		window.ErrorRunWindow(sErr + err.Error())
		log.Println(data.S.Error, err)
		return
	}
	defer f.Close()
	log.SetOutput(f)

	log.Println("-------------------------------------------------------------")

	err = data.Init()
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorInit)
		window.ErrorRunWindow(sErr + err.Error())
		log.Println(data.S.Error, err)
		return
	}

	db, err := sql.Open("mysql", data.DataSourseTcp())
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorOpedDB)
		window.ErrorRunWindow(sErr + err.Error())
		log.Println(data.S.Error, err)
		return
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		err = errors.Wrap(err, data.S.ErrorPingDB)
		window.ErrorRunWindow(sErr + err.Error())
		log.Println(data.S.Error, err)
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

	var mw *walk.MainWindow
	if err := (dec.MainWindow{
		AssignTo: &mw,
		Title:    "Учет",
		Size:     dec.Size{300, 80},
		Layout:   dec.VBox{},
		Children: []dec.Widget{
			dec.PushButton{
				Text: "EntityType",
				OnClicked: func() {
					cmd, err := window.TypeRunDialog(mw, db, "EntityType")
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Printf(data.S.EndWindow, data.S.Type, cmd)
				},
			},
			dec.PushButton{
				Text: "StatusType",
				OnClicked: func() {
					cmd, err := window.TypeRunDialog(mw, db, "StatusType")
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Printf(data.S.EndWindow, data.S.Type, cmd)
				},
			},
			dec.PushButton{
				Text: "Person",
				OnClicked: func() {
					cmd, err := window.TypeRunDialog(mw, db, "Person")
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Printf(data.S.EndWindow, data.S.Type, cmd)
				},
			},
			dec.PushButton{
				Text: "Operation",
				OnClicked: func() {
					cmd, err := window.TypeRunDialog(mw, db, "Operation")
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Printf(data.S.EndWindow, data.S.Type, cmd)
				},
			},
			dec.PushButton{
				Text: "Entities",
				OnClicked: func() {
					cmd, err := window.EntitiesRunDialog(mw, db, true, nil)
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Printf(data.S.EndWindow, data.S.Entities, cmd)
				},
			},
			dec.PushButton{
				Text: "MarkedDetail",
				OnClicked: func() {
					cmd, err := window.MarkedDetailsRunDialog(mw, db, true, nil)
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Printf(data.S.EndWindow, data.S.MarkedDetails, cmd)
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
		err = errors.Wrap(err, data.S.ErrorCreateWindow)
		window.ErrorRunWindow(sErr + err.Error())
		log.Println(data.S.Error, err)
		return
	}
	mw.Run()
}
