package main

import (
	"./artnet"
	"./cmdinterface"
	"bufio"
	"fmt"
	//	"github.com/codegangsta/cli"
	"log"
	"os"
)

func logger(queue chan string) {
	for {
		str := <-queue
		println(str)
	}
}

func cli() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')
		fmt.Println(command)
	}
}

func runApp() {

	chLog := cmdinterface.CmdIface{}
	chLog.Init()

	go logger(chLog.Queue)
	go cli()
	an := artnet.Artnet{}
	an.Setup(&chLog)
	an.Connect("192.168.97.102")
}

func main() {
	runApp()
}
