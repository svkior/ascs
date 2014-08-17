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
		log.Println(str)
	}
}

func runApp() {

	chLog := cmdinterface.CmdIface{}
	chLog.Init()
	an := artnet.Artnet{}
	an.Setup(&chLog)
	an.Connect("192.168.97.104")
	an.SendArtPoll()

	go logger(chLog.Queue)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')
		fmt.Println(command)
	}

}

func main() {
	runApp()
}
