package main

import (
	"./artnet"
	"./ascsgui"
	"./cmdinterface"
	//"fmt"
)

func main() {

	chLog := cmdinterface.CmdIface{}
	chLog.Init()

	gui := ascsgui.Gui{}
	gui.RunGui(&chLog)
	an := artnet.Artnet{}
	an.Setup(&chLog)
	go an.Connect("192.168.97.100")
	
	gui.Wait()
}
