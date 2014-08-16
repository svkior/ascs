package main

import (
	//	"./artnet"
	"./ascsgui"
	"./cmdinterface"
	//"fmt"
)

func main() {
	/*
		an := artnet.Artnet{}
		an.Setup(0x1)
		an.Connect("192.168.97.100")
		an.SendArtPoll()
	*/
	chLog := cmdinterface.CmdIface{}
	chLog.Init()

	gui := ascsgui.Gui{}
	gui.RunGui(&chLog)

	chLog.Log("Hello, World!")

	gui.Wait()
}
