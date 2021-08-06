package text

// import (
// 	. "accounting/data/table"
// )

var T = struct {
	MsgBoxError   string
	MsgBoxWarning string
	MsgBoxInfo    string

	MsgChooseRow  string
	MsgEmptyTitle string
	MsgRepeat     string
	MsgNotChange  string

	MsgInsertedIndex string

	HeadingEntities      string
	HeadingEntity        string
	HeadingEntityRec     string
	HeadingEntityType    string
	HeadingMarkedDetail  string
	HeadingMarkedDetails string

	ColumnCount string
	ColumnTitle string

	LabelEnumerable    string
	LabelCount         string
	LabelComponents    string
	LabelMarking       string
	LabelNote          string
	LabelSpecification string
	LabelTitle         string
	LabelType          string

	SuffixPieces    string
	SuffixComponent string

	ButtonOK     string
	ButtonCansel string
	ButtonAdd    string
	ButtonChange string
	ButtonDelete string
	ButtonSearch string
	ButtonChoose string

	MarkingNo   string
	MarkingAll  string
	MarkingYear string
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
		"Следует проверить корректность данных в последней вставленной строке и перезапустить программу.",

	HeadingEntities:      "Учет - Список компонентов",
	HeadingEntity:        "Учет - Компонент",
	HeadingEntityRec:     "Учет - Дочерний компонент",
	HeadingEntityType:    "Учет - Типы компонентов",
	HeadingMarkedDetail:  "Учет - Маркировка детали",
	HeadingMarkedDetails: "Учет - Список деталей",

	ColumnCount: "Количество",
	ColumnTitle: "Название",

	LabelEnumerable:    "Можно сосчитать:  ",
	LabelCount:         "Количество:",
	LabelComponents:    "Компоненты:",
	LabelMarking:       "Маркировка:",
	LabelSpecification: "Спецификация:",
	LabelNote:          "Примечание:",
	LabelTitle:         "Название:",
	LabelType:          "Тип:",

	SuffixPieces:    " шт",
	SuffixComponent: " компонент",

	ButtonOK:     "OK",
	ButtonCansel: "Отмена",
	ButtonAdd:    "Добавить",
	ButtonChange: "Изменить",
	ButtonDelete: "Удалить",
	ButtonSearch: "Поиск",
	ButtonChoose: "Выбрать",

	MarkingNo:   "Нет",
	MarkingAll:  "Сквозная",
	MarkingYear: "По годам",
}
