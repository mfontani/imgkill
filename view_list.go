package main

import (
	"fmt"

	"github.com/rivo/tview"
)

func displayList(app *tview.Application, images []imageJSONLine) {
	list := tview.NewList()
	for _, v := range images {
		addFuncItem(list, fmt.Sprintf("%s:%s", v.Repository, v.Tag), "", 0)
	}
	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})
	list.SetDoneFunc(func() {
		app.Stop()
	})
	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}

func addFuncItem(l *tview.List, mainText, secondaryText string, shortcut rune) {
	l.AddItem(mainText, secondaryText, shortcut, func() {
		if selectedImages[mainText] {
			delete(selectedImages, mainText)
		} else {
			selectedImages[mainText] = true
		}
		// refresh "view"
		for i := 0; i < l.GetItemCount()-1; i++ {
			main, secondary := l.GetItemText(i)
			if selectedImages[main] {
				secondary = "SELECTED"
			} else {
				secondary = ""
			}
			l.SetItemText(i, main, secondary)
		}
	})
}
