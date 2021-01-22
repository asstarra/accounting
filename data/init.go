package data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

type DataBase struct {
	Host     string `json:"Host"`
	Password string `json:"Password"`
	Database string `json:"Database"`
}

var DB DataBase

func DataSourseTcp() string {
	return fmt.Sprint(DB.Host, ":", DB.Password, "@tcp/", DB.Database)
}

type Column struct {
	Name string `json:"Name"`
}

type Table struct {
	Name    string            `json:"Name"`
	Columns map[string]Column `json:"Columns"`
}

type Tables map[string]Table

var Tab Tables

type Register struct {
	IsSaveDialog *walk.MutableCondition
}

func (r Register) String() string {
	s := "{"
	s += "IsSaveDialog:" + fmt.Sprint(r.IsSaveDialog)
	s += "}"
	return s
}

var Reg Register

func initRegister() {
	Reg.IsSaveDialog = walk.NewMutableCondition()
	dec.MustRegisterCondition("IsSaveDialog", Reg.IsSaveDialog)
}

func initFromFile(filename string, data interface{}) error {
	configFile, err := os.Open(filename)
	if err != nil {
		err = errors.Wrap(err, "Не удалось открыть файл "+filename)
		return err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(data)
	if err != nil {
		err = errors.Wrap(err, "Ошибка чтения данных в файле "+filename)
		return err
	}
	return nil
}

func Init() error {
	err := initFromFile("config/database.json", &DB)
	if err != nil {
		return err
	}
	log.Println("DEBUG -", "InitDatabase", DB)
	err = initFromFile("config/table.json", &Tab)
	if err != nil {
		return err
	}
	log.Println("DEBUG -", "InitTables", Tab)
	initRegister()
	log.Println("DEBUG -", "InitRegister", Reg)
	return nil
}
