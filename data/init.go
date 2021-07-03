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

// Переменные для подключения к БД.
var DB struct {
	Host     string `json:"Host"`
	Password string `json:"Password"`
	Database string `json:"Database"`
}

// Получение строки подключения.
func DataSourseTcp() string {
	return fmt.Sprint(DB.Host, ":", DB.Password, "@tcp/", DB.Database)
}

// Название полей и таблиц.
var Tab map[string]Table

type Table struct {
	Name    string            `json:"Name"`
	Columns map[string]Column `json:"Columns"`
}

type Column struct {
	Name string `json:"Name"`
}

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

var C = struct {
	TimeLayoutMySql string
	TimeLayoutDay   string
}{
	TimeLayoutMySql: "2006-01-02 15:04:05",
	TimeLayoutDay:   "2006-01-02",
}

// Строковые константы.
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

	Entity        string
	Entities      string
	EntityRec     string
	Type          string
	MarkedDetail  string
	MarkedDetails string

	HeadingEntity        string
	HeadingEntities      string
	HeadingEntityRec     string
	HeadingType          string
	HeadingMarkedDetail  string
	HeadingMarkedDetails string

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

	InEntitiesRunDialog       string
	InEntityRunDialog         string
	InEntityRecRunDialog      string
	InTypeRunDialog           string
	InMarkedDetailRunDialog   string
	InMarkedDetailsRunDialog  string
	InSelectEntities          string
	InSelectEntityRecChild    string
	InSelectEntityRec         string
	InSelectIdTitle           string
	InSelectMarkingLineNew    string
	InSelectMarkingLineOld    string
	InSelectMarkingLineEntity string
	InSelectMarkedDetail      string
	InSelectMarkedDetails     string

	MsgChooseRow  string
	MsgEmptyTitle string

	ErrorTableInit         string
	ErrorTypeInit          string
	ErrorCreateWindow      string
	ErrorCreateWindowErr   string
	ErrorUnexpectedColumn  string
	ErrorOpenFile          string
	ErrorReadFile          string
	ErrorInit              string
	ErrorOpedDB            string
	ErrorPingDB            string
	ErrorQueryDB           string
	ErrorAddDB             string
	ErrorChangeDB          string
	ErrorDeleteDB          string
	ErrorDecryptRow        string
	ErrorAddRow            string
	ErrorChangeRow         string
	ErrorDeleteRow         string
	ErrorInsertIndexLog    string
	ErrorInsertIndex       string
	ErrorSubmit            string
	ErrorChoose            string
	ErrorRead              string
	ErrorUpdate            string
	ErrorSubquery          string
	ErrorGraphCircle       string
	ErrorUpdateMarkingLine string
	ErrorNil               string
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

	Entity:        "ENTITY",
	Entities:      "ENTITIES",
	EntityRec:     "ENTITY_REC",
	Type:          "TYPE",
	MarkedDetail:  "MARKED_DETAIL",
	MarkedDetails: "MARKED_DETAILS",

	HeadingEntity:        "Учет - Сущность",
	HeadingEntities:      "Учет - Сущности",
	HeadingEntityRec:     "Учет - Дочерний компонент",
	HeadingType:          "Учет - Типы",
	HeadingMarkedDetail:  "Учет - Маркировка детали",
	HeadingMarkedDetails: "Учет - Список деталей",

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

	InEntitiesRunDialog:       "In EntitiesRunDialog(isChage = %t, IdTitle = %v)",
	InEntityRunDialog:         "In EntityRunDialog(entity = %v)",
	InEntityRecRunDialog:      "In EntityRecRunDialog(child = %v)",
	InTypeRunDialog:           "In TypeRunDialog(tableName = %s)",
	InMarkedDetailRunDialog:   "In MarkedDetailRunDialog(detail = %v)",
	InMarkedDetailsRunDialog:  "In MarkedDetailsRunDialog(isChage = %t, parent detail = %v)",
	InSelectEntities:          "In SelectEntities(title = \"%s\", entityType = %d)",
	InSelectEntityRecChild:    "In SelectEntityRecChild(parent = %d)",
	InSelectEntityRec:         "In SelectEntityRec()",
	InSelectIdTitle:           "In SelectIdTitle(tableName = %s)",
	InSelectMarkingLineNew:    "In SelectMarkingLineNew()",
	InSelectMarkingLineOld:    "In SelectMarkingLineOld()",
	InSelectMarkingLineEntity: "In SelectMarkingLineEntity(id = %d)",
	InSelectMarkedDetail:      "",
	InSelectMarkedDetails:     "In SelectMarkedDetails(marking = %v)",

	MsgChooseRow:  "Выберите строчку",
	MsgEmptyTitle: "Название не может состоять из пустой строки",

	ErrorTableInit:        "При заполнении таблицы произошла ошибка",
	ErrorTypeInit:         "Не удалось узнать список типов",
	ErrorCreateWindow:     "Не удалось создать окно",
	ErrorCreateWindowErr:  "Не удалось создать окно для ошибки. Текст ошибки = ",
	ErrorUnexpectedColumn: "Обращение к неизвестному столбцу",
	ErrorOpenFile:         "Не удалось открыть файл ",
	ErrorReadFile:         "Не корректные данные в файле ",
	ErrorInit:             "Ошибка инициализации",
	ErrorOpedDB:           "Не удалось открыть соединение к базе данных",
	ErrorPingDB:           "Не удалось подключится к базе данных",
	ErrorQueryDB:          "Ошибка запроса к базе данных.\nСтрока запроса = ",
	ErrorAddDB:            "Не удалось добавить строку в базу данных.\nСтрока запроса = ",
	ErrorChangeDB:         "Не удалось изменить строку в базе данных.\nСтрока запроса = ",
	ErrorDeleteDB:         "Не удалось удалить строку из базы данных.\nСтрока запроса = ",
	ErrorDecryptRow:       "Не удалось расшифровать строку",
	ErrorAddRow:           "Не удалось добавить строку",
	ErrorChangeRow:        "Не удалось изменить строку",
	ErrorDeleteRow:        "Не удалось удалить строку",
	ErrorInsertIndexLog:   "При вставке новой строки в базу данных не удалось узнать индекс вставляемой строки",
	ErrorInsertIndex: "Это сообщение не должно показываться.\n" +
		"При вставке новой строки в базу данных не удалось узнать индекс  вставляемой строки.\n" +
		"Следует перезапустить программу и проверить корректность данных в последней вставленной строке.",
	ErrorSubmit:            "Не удалось разместить данные",
	ErrorChoose:            "Не удалось выбрать данные",
	ErrorRead:              "Не удалось считать данные",
	ErrorUpdate:            "Не удалось обновить данные",
	ErrorSubquery:          "Не удалось сделать подзапрос",
	ErrorGraphCircle:       "Иерархия не может быть циклической",
	ErrorUpdateMarkingLine: "При обновлении иерархии производственных линий произошла ошибка",
	ErrorNil:               "Некорректные данные (nil)",
}

// Чтение json из файла.
func initFromFile(filename string, data interface{}) error {
	configFile, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, S.ErrorOpenFile+filename)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(data)
	if err != nil {
		return errors.Wrap(err, S.ErrorReadFile+filename)
	}
	return nil
}

// Инициализация глобальных переменных.
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
