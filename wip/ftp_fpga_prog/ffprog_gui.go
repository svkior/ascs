package main

import (
	"github.com/andlabs/ui"
	//"image"
	//"reflect"
)

var w ui.Window

func initGUI() {
	b := ui.NewButton("Button")
	stack := ui.NewVerticalStack(b)
	w = ui.NewWindow("Window", 400, 500, stack)
	w.OnClosing(func() bool {
		ui.Stop()
		return true
	})
	w.Show()

}

func main() {
	go ui.Do(initGUI)
	err := ui.Go()
	if err != nil {
		panic(err)
	}
}
