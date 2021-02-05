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

var S = struct {
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
}{
	Panic:   "PANIC!",
	Error:   "ERROR!",
	Warning: "WARNING!",
	Info:    "INFO:",
	Debug:   "DEBUG:",

	MsgBoxError:   "Ошибка!",
	MsgBoxWarning: "Внимание!",
	MsgBoxInfo:    "Информация",

	BeginWindow:  "INFO: BEGIN window %s",
	InitWindow:   "INFO: INIT window %s",
	CreateWindow: "INFO: CREATE window %s",
	RunWindow:    "INFO: RUN window %s",
	EndWindow:    "INFO: END window %s, cmd %v",

	Entity:    "ENTITY",
	Entities:  "ENTITIES",
	EntityRec: "ENTITY_REC",
	Type:      "TYPE",

	HeadingEntity:    "Учет - Сущность",
	HeadingEntities:  "Учет - Сущности",
	HeadingEntityRec: "Учет - Дочерний компонент",
	HeadingType:      "Учет - Типы",

	LogOk:     "Ok",
	LogCansel: "Cansel",
	LogAdd:    "Add",
	LogChange: "Change",
	LogDelete: "Delete",
	LogChoose: "Choose",
	LogSearch: "Search",

	ButtonOK:     "OK",
	ButtonCansel: "Отмена",
	ButtonAdd:    "Добавить",
	ButtonChange: "Изменить",
	ButtonDelete: "Удалить",
	ButtonSearch: "Поиск",

	InEntitiesRunDialog:    "In EntitiesRunDialog(isChage = %t, IdTitle = %v)",
	InEntityRunDialog:      "In EntityRunDialog(entity = %v)",
	InEntityRecRunDialog:   "In EntityRecRunDialog(child = %v)",
	InTypeRunDialog:        "In TypeRunDialog(tableName = %s)",
	InSelectEntities:       "In SelectEntities(title = \"%s\", entityType = %d)",
	InSelectEntityRecChild: "In SelectEntityRecChild(parent = %d)",
	InSelectIdTitle:        "In SelectIdTitle(tableName = %s)",

	MsgChooseRow:  "Выберите строчку",
	MsgEmptyTitle: "Название не может состоять из пустой строки",

	ErrorTableInit:        "При заполнении таблицы произошла ошибка",
	ErrorTypeInit:         "Не удалось узнать список типов",
	ErrorCreateWindow:     "Could not create Window Form",
	ErrorUnexpectedColumn: "Unexpected column",
	ErrorOpenFile:         "Не удалось открыть файл ",
	ErrorReadFile:         "Ошибка чтения данных в файле ",
	ErrorInit:             "Ошибка инициализации",
	ErrorOpedDB:           "Не удалось открыть соединение к базе данных",
	ErrorPingDB:           "Не удалось подключится к базе данных",
	ErrorQuery:            "Ошибка запроса к базе данных. Строка запроса = ",
	ErrorDecryptRow:       "Не удалось расшифровать строку",
	ErrorAdd:              "Не удалось добавить строку",
	ErrorChange:           "Не удалось изменить строку",
	ErrorDelete:           "Не удалось удалить строку",
	ErrorAddDB:            "Не удалось добавить строку в базу данных. Строка запроса = ",
	ErrorChangeDB:         "Не удалось изменить строку в базе данных. Строка запроса = ",
	ErrorDeleteDB:         "Не удалось удалить строку из базы данных. Строка запроса = ",
	ErrorInsertIndexLog:   "При вставке новой строки в базу данных не удалось узнать индекс вставляемой строки",
	ErrorInsertIndex: "Это сообщение не должно показываться.\n" +
		"При вставке новой строки в базу данных не удалось узнать индекс  вставляемой строки.\n" +
		"Следует перезапустить программу и проверить корректность данных в последней вставленной строке.",
	ErrorSubmit:   "Не удалось сохранить данные",
	ErrorChoose:   "Не удалось выбрать данные",
	ErrorSubquery: "Не удалось сделать подзапрос",
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
