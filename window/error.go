package window

import (
	"accounting/data"
	"log"

	dec "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
)

// Описание и запуск диалогового окна.
func ErrorRunWindow(s string) {
	if _, err := (dec.MainWindow{
		Title:  data.S.MsgBoxError,
		Size:   dec.Size{300, 80},
		Layout: dec.VBox{},
		Children: []dec.Widget{
			dec.Label{
				Text: s,
			},
		},
	}.Run()); err != nil {
		log.Println(data.S.Error, errors.Wrap(err, data.S.ErrorCreateWindowErr+s))
	}
}
