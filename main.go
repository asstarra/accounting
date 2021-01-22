package main

import (
	"database/sql"
	// "fmt"
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

func main() {
	f, err := os.OpenFile("logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("-----------------------------------------------------------------------------------------------------------------------------------")

	err = data.Init()
	if err != nil {
		err = errors.Wrap(err, "Ошибка инициализации")
		log.Fatalln("ERROR!", err)
	}

	db, err := sql.Open("mysql", data.DataSourseTcp())
	if err != nil {
		err = errors.Wrap(err, "Не удалось открыть соединение к БД")
		log.Fatalln("ERROR!", err)
	}
	defer db.Close()

	entity := window.Entity{}
	entityRec := window.EntityRecChild{}

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
					if cmd, err := window.TypeDialogRun(mw, "Учет - Тип сущности", "EntityType", db); err != nil {
						log.Println("ERROR!", err)
					} else {
						log.Println("INFO -", "END window - TYPE, tableName - EntityType, cmd -", cmd)
					}
				},
			},
			dec.PushButton{
				Text: "Entity",
				OnClicked: func() {
					if cmd, err := window.EntityRunDialog(mw, &entity, db); err != nil {
						log.Println("ERROR!", err)
					} else {
						log.Println("INFO -", "END window - ENTITY, cmd -", cmd)
					}
				},
			},
			dec.PushButton{
				Text: "Entity_rec",
				OnClicked: func() {
					if cmd, err := window.EntityRecRunDialog(mw, &entityRec, db); err != nil {
						log.Println("ERROR!", err)
					} else {
						log.Println("INFO -", "END window - ENTITY, cmd -", cmd)
					}
				},
			},
		},
	}.Create()); err != nil {
		err = errors.Wrap(err, "Could not create Window Form")
		log.Fatalln("ERROR!", err)
	}
	mw.Run()
}
