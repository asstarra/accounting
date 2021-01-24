package main

import (
	"database/sql"
	// "fmt"

	// "io"
	"log"
	"os"

	"accounting/data"
	"accounting/window"
)

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

func errorWindow(s string) {
	if _, err := (dec.MainWindow{
		Title:  "КРИТИЧЕСКАЯ ОШИБКА!",
		Size:   dec.Size{300, 80},
		Layout: dec.VBox{},
		Children: []dec.Widget{
			dec.Label{
				Text: s,
			},
		},
	}.Run()); err != nil {
		log.Printf("ERROR! In errorWindow, err = %v", err)
	}
}

func main() {
	sErr := "Error: "
	f, err := os.OpenFile("logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorOpenFile+"logfile.txt")
		errorWindow(sErr + err.Error())
		log.Println("ERROR!", err)
		return
	}
	defer f.Close()
	log.SetOutput(f)

	// err = syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
	// if err != nil {
	// 	log.Println("ERROR!", err)
	// 	return
	// }

	log.Println("-------------------------------------------------------------")

	err = data.Init()
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorInit)
		errorWindow(sErr + err.Error())
		log.Println(data.S.Error, err)
		return
	}

	db, err := sql.Open("mysql", data.DataSourseTcp())
	if err != nil {
		err = errors.Wrap(err, data.S.ErrorOpedDB)
		errorWindow(sErr + err.Error())
		log.Println(data.S.Error, err)
		return
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		err = errors.Wrap(err, data.S.ErrorPingDB)
		errorWindow(sErr + err.Error())
		log.Println(data.S.Error, err)
		return
	}

	entity := window.Entity{}
	entityRec := window.EntityRecChild{}
	idTitle := window.IdTitle{}

	var mw *walk.MainWindow
	if err := (dec.MainWindow{
		AssignTo: &mw,
		Title:    "Учет",
		Size:     dec.Size{300, 80},
		Layout:   dec.VBox{},
		Children: []dec.Widget{
			dec.PushButton{
				Text: "Type",
				OnClicked: func() {
					cmd, err := window.TypeDialogRun(mw, "Учет - Тип сущности", "EntityType", db)
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Printf(data.S.EndWindow, data.S.Type, cmd)
				},
			},
			dec.PushButton{
				Text: "Entity",
				OnClicked: func() {
					cmd, err := window.EntityRunDialog(mw, db, &entity)
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Printf(data.S.EndWindow, data.S.Entity, cmd)
				},
			},
			dec.PushButton{
				Text: "Entity_rec",
				OnClicked: func() {
					cmd, err := window.EntityRecRunDialog(mw, db, false, &entityRec)
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Printf(data.S.EndWindow, data.S.EntityRec, cmd)
				},
			},
			dec.PushButton{
				Text: "Entities",
				OnClicked: func() {
					cmd, err := window.EntitiesRunDialog(mw, db, true, &idTitle)
					if err != nil {
						log.Println("ERROR!", err)
					}
					log.Printf(data.S.EndWindow, data.S.Entities, cmd)
				},
			},
		},
	}.Create()); err != nil {
		err = errors.Wrap(err, data.S.ErrorCreateWindow)
		errorWindow(sErr + err.Error())
		log.Println(data.S.Error, err)
		return
	}
	mw.Run()
}
