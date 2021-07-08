package data

import (
	// "database/sql"
	"accounting/data/db"
	e "accounting/data/errors"
	l "accounting/data/log"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

// Храним глобальные регистры.
var Reg Register

type Register struct {
	IsSaveDialog *walk.MutableCondition // GO-TO
}

func (r Register) String() string {
	s := "{"
	s += "IsSaveDialog:" + fmt.Sprint(r.IsSaveDialog)
	s += "}"
	return s
}

func initRegister() {
	Reg.IsSaveDialog = walk.NewMutableCondition()
	dec.MustRegisterCondition("IsSaveDialog", Reg.IsSaveDialog)
}

// Храним иконки для разных сообщений.
var Icon = struct {
	Critical walk.MsgBoxStyle
	Error    walk.MsgBoxStyle
	Warning  walk.MsgBoxStyle
	Info     walk.MsgBoxStyle
}{
	Critical: walk.MsgBoxIconError,
	Error:    walk.MsgBoxIconError,
	Warning:  walk.MsgBoxIconWarning,
	Info:     walk.MsgBoxIconInformation,
}

// Чтение json из файла.
func initFromFile(filename string, data interface{}) error {
	configFile, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, e.Err.ErrorOpenFile+filename)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(data)
	if err != nil {
		return errors.Wrap(err, e.Err.ErrorReadFile+filename)
	}
	return nil
}

// Инициализация глобальных переменных.
func Init() error {
	err := initFromFile("config/database.json", &db.DB)
	if err != nil {
		return err
	}
	log.Println(l.Debug, "InitDatabase")
	err = initFromFile("config/table.json", &db.Tab)
	if err != nil {
		return err
	}
	log.Println(l.Debug, "InitTables")
	initRegister()
	log.Println(l.Debug, "InitRegister")
	return nil
}
