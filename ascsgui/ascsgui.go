package ascsgui

import (
	"../cmdinterface"
	br "./bridge"
	"fmt"
	"gopkg.in/qml.v0"
)

type Gui struct {
	win       *qml.Window
	eng       *qml.Engine
	component qml.Object
	context   *qml.Context
}

func (g *Gui) RunGui(logger *cmdinterface.CmdIface) {
	qml.Init(nil)
	g.eng = qml.NewEngine()

	var err error
	g.component, err = g.eng.LoadFile("./ascsgui/main.qml")
	if err != nil {
		fmt.Printf("Error load file: %s", err.Error())
		return
	}

	logs := br.Logs{}

	bridge := br.Bridge{}
	bridge.BLogs = &logs

	g.context = g.eng.Context()
	g.context.SetVar("bridge", &bridge)
	g.context.SetVar("logs", &logs)

	g.win = g.component.CreateWindow(nil)

	bridge.SetupVars(g.win.Root())

	g.win.Show()
	go bridge.UDPDayTimeServer()
	go bridge.StartLogger(logger)
}

func (g *Gui) Wait() {
	g.win.Wait()

}
