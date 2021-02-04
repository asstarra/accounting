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

var DB struct {
	Host     string `json:"Host"`
	Password string `json:"Password"`
	Database string `json:"Database"`
}

func DataSourseTcp() string {
	return fmt.Sprint(DB.Host, ":", DB.Password, "@tcp/", DB.Database)
}

var Tab map[string]Table

type Table struct {
	Name    string            `json:"Name"`
	Columns map[string]Column `json:"Columns"`
}

type Column struct {
	Name string `json:"Name"`
}

var Reg Register

type Register struct {
	IsSaveDialog *walk.MutableCondition
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

var Icon struct {
	Critical walk.MsgBoxStyle
	Error    walk.MsgBoxStyle
	Warning  walk.MsgBoxStyle
	Info     walk.MsgBoxStyle
}

func initIcon() {
	Icon.Critical = walk.MsgBoxIconError
	Icon.Error = walk.MsgBoxIconError
	Icon.Warning = walk.MsgBoxIconWarning
	Icon.Info = walk.MsgBoxIconInformation
}

var S struct {
	Panic   string
	Error   string
	Warning string
	Info    string
	Debug   string

	MsgBoxError   string
	MsgBoxWarning string
	MsgBoxInfo    string

	BeginWindow  string
	InitWindow   string
	CreateWindow string
	RunWindow    string
	EndWindow    string

	Entity    string
	Entities  string
	EntityRec string
	Type      string

	HeadingEntity    string
	HeadingEntities  string
	HeadingEntityRec string
	HeadingType      string

	LogOk     string
	LogCansel string
	LogAdd    string
	LogChange string
	LogDelete string
	LogChoose string
	LogSearch string

	ButtonOK     string
	ButtonCansel string
	ButtonAdd    string
	ButtonChange string
	ButtonDelete string
	ButtonSearch string

	InEntitiesRunDialog    string
	InEntityRunDialog      string
	InEntityRecRunDialog   string
	InTypeRunDialog        string
	InSelectEntities       string
	InSelectEntityRecChild string
	InSelectIdTitle        string

	MsgChooseRow  string
	MsgEmptyTitle string

	ErrorTableInit        string
	ErrorTypeInit         string
	ErrorCreateWindow     string
	ErrorUnexpectedColumn string
	ErrorOpenFile         string
	ErrorReadFile         string
	ErrorInit             string
	ErrorOpedDB           string
	ErrorPingDB           string
	ErrorQuery            string
	ErrorDecryptRow       string
	ErrorAdd              string
	ErrorChange           string
	ErrorDelete           string
	ErrorAddDB            string
	ErrorChangeDB         string
	ErrorDeleteDB         string
	ErrorInsertIndexLog   string
	ErrorInsertIndex      string
	ErrorSubmit           string
	ErrorChoose           string
	ErrorSubquery         string
}

func initString() {
	S.Panic = "PANIC!"
	S.Error = "ERROR!"
	S.Warning = "WARNING!"
	S.Info = "INFO:"
	S.Debug = "DEBUG:"

	S.MsgBoxError = "Ошибка!"
	S.MsgBoxWarning = "Внимание!"
	S.MsgBoxInfo = "Информация"

	S.BeginWindow = "INFO: BEGIN window %s"
	S.InitWindow = "INFO: INIT window %s"
	S.CreateWindow = "INFO: CREATE window %s"
	S.RunWindow = "INFO: RUN window %s"
	S.EndWindow = "INFO: END window %s, cmd %v"

	S.Entity = "ENTITY"
	S.Entities = "ENTITIES"
	S.EntityRec = "ENTITY_REC"
	S.Type = "TYPE"

	S.HeadingEntity = "Учет - Сущность"
	S.HeadingEntities = "Учет - Сущности"
	S.HeadingEntityRec = "Учет - Дочерний компонент"
	S.HeadingType = "Учет - Типы"

	S.LogOk = "Ok"
	S.LogCansel = "Cansel"
	S.LogAdd = "Add"
	S.LogChange = "Change"
	S.LogDelete = "Delete"
	S.LogChoose = "Choose"
	S.LogSearch = "Search"

	S.ButtonOK = "OK"
	S.ButtonCansel = "Отмена"
	S.ButtonAdd = "Добавить"
	S.ButtonChange = "Изменить"
	S.ButtonDelete = "Удалить"
	S.ButtonSearch = "Поиск"

	S.InEntitiesRunDialog = "In EntitiesRunDialog(isChage = %t, IdTitle = %v)"
	S.InEntityRunDialog = "In EntityRunDialog(entity = %v)"
	S.InEntityRecRunDialog = "In EntityRecRunDialog(child = %v)"
	S.InTypeRunDialog = "In TypeRunDialog(tableName = %s)"
	S.InSelectEntities = "In SelectEntities(title = \"%s\", entityType = %d)"
	S.InSelectEntityRecChild = "In SelectEntityRecChild(parent = %d)"
	S.InSelectIdTitle = "In SelectIdTitle(tableName = %s)"

	S.MsgChooseRow = "Выберите строчку"
	S.MsgEmptyTitle = "Название не может состоять из пустой строки"

	S.ErrorTableInit = "При заполнении таблицы произошла ошибка"
	S.ErrorTypeInit = "Не удалось узнать список типов"
	S.ErrorCreateWindow = "Could not create Window Form"
	S.ErrorUnexpectedColumn = "Unexpected column"
	S.ErrorOpenFile = "Не удалось открыть файл "
	S.ErrorReadFile = "Ошибка чтения данных в файле "
	S.ErrorInit = "Ошибка инициализации"
	S.ErrorOpedDB = "Не удалось открыть соединение к базе данных"
	S.ErrorPingDB = "Не удалось подключится к базе данных"
	S.ErrorQuery = "Ошибка запроса к базе данных. Строка запроса = "
	S.ErrorDecryptRow = "Не удалось расшифровать строку"
	S.ErrorAdd = "Не удалось добавить строку"
	S.ErrorChange = "Не удалось изменить строку"
	S.ErrorDelete = "Не удалось удалить строку"
	S.ErrorAddDB = "Не удалось добавить строку в базу данных. Строка запроса = "
	S.ErrorChangeDB = "Не удалось изменить строку в базе данных. Строка запроса = "
	S.ErrorDeleteDB = "Не удалось удалить строку из базы данных. Строка запроса = "
	S.ErrorInsertIndexLog = "При вставке новой строки в базу данных не удалось узнать индекс вставляемой строки"
	S.ErrorInsertIndex = "Это сообщение не должно показываться.\n" +
		"При вставке новой строки в базу данных не удалось узнать индекс  вставляемой строки.\n" +
		"Следует перезапустить программу и проверить корректность данных в последней вставленной строке."
	S.ErrorSubmit = "Не удалось сохранить данные"
	S.ErrorChoose = "Не удалось выбрать данные"
	S.ErrorSubquery = "Не удалось сделать подзапрос"
}

func initFromFile(filename string, data interface{}) error {
	configFile, err := os.Open(filename)
	if err != nil {
		err = errors.Wrap(err, S.ErrorOpenFile+filename)
		return err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(data)
	if err != nil {
		err = errors.Wrap(err, S.ErrorReadFile+filename)
		return err
	}
	return nil
}

func Init() error {
	initString()
	initIcon()
	log.Println(S.Debug, "InitString")
	err := initFromFile("config/database.json", &DB)
	if err != nil {
		return err
	}
	log.Println(S.Debug, "InitDatabase")
	err = initFromFile("config/table.json", &Tab)
	if err != nil {
		return err
	}
	log.Println(S.Debug, "InitTables")
	initRegister()
	log.Println(S.Debug, "InitRegister")
	return nil
}
