package window

import (
	e "accounting/data/errors"
	l "accounting/data/log"
	"accounting/data/text"
	"log"

	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

// Описание и запуск диалогового окна.
func ErrorRunWindow(s string) {
	if _, err := (dec.MainWindow{
		Title:  text.T.MsgBoxError,
		Size:   dec.Size{300, 80},
		Layout: dec.VBox{},
		Children: []dec.Widget{
			dec.Label{
				Text: s,
			},
		},
	}.Run()); err != nil {
		log.Println(l.Error, errors.Wrap(err, e.Err.ErrorCreateWindowErr+s)) // Лог.
	}
}
