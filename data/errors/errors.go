package errors

var Err = struct {
	ErrorTableInit         string
	ErrorTypeInit          string
	ErrorCreateWindow      string
	ErrorCreateWindowErr   string
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
	ErrorDecryptTime       string
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
	ErrorTableInit:       "При заполнении таблицы произошла ошибка",
	ErrorTypeInit:        "Не удалось узнать список типов",
	ErrorCreateWindow:    "Не удалось создать окно",
	ErrorCreateWindowErr: "Не удалось создать окно для ошибки. Текст ошибки = ",
	ErrorOpenFile:        "Не удалось открыть файл ",
	ErrorReadFile:        "Не корректные данные в файле ",
	ErrorInit:            "Ошибка инициализации",
	ErrorOpedDB:          "Не удалось открыть соединение к базе данных",
	ErrorPingDB:          "Не удалось подключится к базе данных",
	ErrorQueryDB:         "Ошибка запроса к базе данных.\nСтрока запроса = \"%s\"",
	ErrorAddDB:           "Не удалось добавить строку в базу данных.\nСтрока запроса = \"%s\"",
	ErrorChangeDB:        "Не удалось изменить строку в базе данных.\nСтрока запроса = \"%s\"",
	ErrorDeleteDB:        "Не удалось удалить строку из базы данных.\nСтрока запроса = \"%s\"",
	ErrorDecryptRow:      "Не удалось расшифровать строку",
	ErrorDecryptTime:     "Не удалось расшифровать время. Строка = \"%s\"",
	ErrorAddRow:          "Не удалось добавить строку",
	ErrorChangeRow:       "Не удалось изменить строку",
	ErrorDeleteRow:       "Не удалось удалить строку",
	ErrorInsertIndexLog:  "При вставке новой строки в базу данных не удалось узнать индекс вставляемой строки",
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

var (
	NilPointer       = "invalid memory address or nil pointer dereference"
	WrongType        = "nil pointer or incorrect type in interface"
	UnexpectedColumn = "unexpected column"
)
