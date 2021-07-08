package window

import (
	"strings"
)

// Функция конвертирующая ошибки для показа пользователю.
func MsgError(err error) string {
	return strings.Replace(err.Error(), ": ", ":\n", -1)
}

func MaxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}
