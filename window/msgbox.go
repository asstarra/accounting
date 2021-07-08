package window

import (
	e "accounting/data/errors"
	l "accounting/data/log"
	"accounting/data/text"
	"log"
	"strings"

	"github.com/lxn/walk"
	"github.com/pkg/errors"
)

// Функция конвертирующая ошибки для показа пользователю.
func MsgError(err error) string {
	return strings.Replace(err.Error(), ": ", ":\n", -1)
}

// Проверка на пустую строку.
func IsStringEmpty(owner walk.Form, str string) bool {
	if str == "" {
		walk.MsgBox(owner, text.T.MsgBoxInfo, text.T.MsgEmptyTitle, walk.MsgBoxIconInformation)
		return true
	}
	return false
}

type modelTable interface {
	RowCount() int
	Value(row, col int) interface{}
	Equal(row, col int) bool
}

// Проверка на совпадение значений. len(cols) должно быть > 0.
func IsRepeat(owner walk.Form, tab modelTable, cols []int) bool {
	for i := 0; i < tab.RowCount(); i++ {
		var flag bool = true
		for _, colNum := range cols {
			flag = flag && tab.Equal(i, colNum)
		}
		if flag {
			walk.MsgBox(owner, text.T.MsgBoxInfo, text.T.MsgRepeat, walk.MsgBoxIconInformation)
			return true
		}
	}
	return false
}

// Проверка на выделение изменяемой строчки.
func IsCorrectIndex(owner walk.Form, tab modelTable, tabView *walk.TableView) bool {
	if tab.RowCount() <= 0 || tabView.CurrentIndex() < 0 {
		walk.MsgBox(owner, text.T.MsgBoxInfo, text.T.MsgChooseRow, walk.MsgBoxIconInformation)
		return true
	}
	return false
}

// Сообщение о том, что строчку нельзя изменить.
func MsgBoxNotChange(owner walk.Form) {
	walk.MsgBox(owner, text.T.MsgBoxInfo, text.T.MsgNotChange, walk.MsgBoxIconInformation)
}

// Сообщение о том, что не получен индекс вставляемой строки.
func MsgBoxNotInsertedId(owner walk.Form, err error) {
	log.Println(l.Error, errors.Wrap(err, l.LogNotInsertedId)) // Лог.
	err = errors.Wrap(err, text.T.MsgInsertedIndex)
	walk.MsgBox(owner, text.T.MsgBoxError, MsgError(err), walk.MsgBoxIconError)
}

func MsgBoxAdd(owner walk.Form, err error) {
	err = errors.Wrap(err, e.Err.ErrorAddRow)
	log.Println(l.Error, err) // Лог.
	walk.MsgBox(owner, text.T.MsgBoxError, MsgError(err), walk.MsgBoxIconError)
}

func MsgBoxChange(owner walk.Form, err error) {
	err = errors.Wrap(err, e.Err.ErrorChangeRow)
	log.Println(l.Error, err) // Лог.
	walk.MsgBox(owner, text.T.MsgBoxError, MsgError(err), walk.MsgBoxIconError)
}

func MsgBoxDelete(owner walk.Form, err error) {
	err = errors.Wrap(err, e.Err.ErrorDeleteRow)
	log.Println(l.Error, err) // Лог.
	walk.MsgBox(owner, text.T.MsgBoxError, MsgError(err), walk.MsgBoxIconError)
}
