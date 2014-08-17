package main

import (
	"./artnet"
	"./cmdinterface"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func logger(queue chan string) {
	for {
		str := <-queue
		fmt.Fprintf(os.Stdout, "\r")
		const layout = "Jan 2, 2006 at 3:04pm (MST)"
		strSplits := strings.Split(str, "\n")
		var firstTime bool = true
		for _, s := range strSplits {
			if len(s) > 0 {
				if firstTime {
					firstTime = false
					s = time.Now().Format(layout) + " : " + s
				}
				fmt.Fprintf(os.Stdout, "%s\n", s)
			}
		}
	}
}

func cli() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')
		if len(command) > 1 {
			switch command {
			case "quit\n":
				println("Quit.")
				os.Exit(0)
			default:
				fmt.Print("Unknown command: " + command)
			}
		}
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
