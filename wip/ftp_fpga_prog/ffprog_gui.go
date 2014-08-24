package main

import (
	"fmt"
	"gopkg.in/qml.v0"
)

func main() {
	qml.Init(nil)
	eng := qml.NewEngine()
	comp, err := eng.LoadFile("./main.qml")
	if err != nil {
		fmt.Printf("Error load file: %s", err.Error())
		return
	}

	win := comp.CreateWindow(nil)
	win.Show()
	win.Wait()
}
