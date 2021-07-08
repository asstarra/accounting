package text

var T = struct {
	MsgBoxError   string
	MsgBoxWarning string
	MsgBoxInfo    string

	MsgChooseRow  string
	MsgEmptyTitle string
	MsgRepeat     string
	MsgNotChange  string

	MsgInsertedIndex string

	HeadingEntity        string
	HeadingEntities      string
	HeadingEntityRec     string
	HeadingEntityType    string
	HeadingMarkedDetail  string
	HeadingMarkedDetails string

	TextEntityTypeTitle string

	ButtonOK     string
	ButtonCansel string
	ButtonAdd    string
	ButtonChange string
	ButtonDelete string
	ButtonSearch string
}{
	MsgBoxError:   "Ошибка!",
	MsgBoxWarning: "Внимание!",
	MsgBoxInfo:    "Информация",

	MsgChooseRow:  "Выберите строчку, которую хотите изменить.",
	MsgEmptyTitle: "Название не может состоять из пустой строки.",
	MsgRepeat:     "Такой элемент в таблице уже существует.",
	MsgNotChange:  "Данную строчку нельзя изменить.",

	MsgInsertedIndex: "Это сообщение не должно показываться.\n" +
		"При вставке новой строки в базу данных не удалось узнать индекс  вставляемой строки.\n" +
		"Следует перезапустить программу и проверить корректность данных в последней вставленной строке.",

	HeadingEntity:        "Учет - Сущность",
	HeadingEntities:      "Учет - Сущности",
	HeadingEntityRec:     "Учет - Дочерний компонент",
	HeadingEntityType:    "Учет - Типы компонентов",
	HeadingMarkedDetail:  "Учет - Маркировка детали",
	HeadingMarkedDetails: "Учет - Список деталей",

	TextEntityTypeTitle: "Название",

	ButtonOK:     "OK",
	ButtonCansel: "Отмена",
	ButtonAdd:    "Добавить",
	ButtonChange: "Изменить",
	ButtonDelete: "Удалить",
	ButtonSearch: "Поиск",
}

var HeadingType = map[string]string{ // TO-DO
	"EntityType": T.HeadingEntityType,
}

var TitleType = map[string]string{ // TO-DO
	"EntityType": T.TextEntityTypeTitle,
}
