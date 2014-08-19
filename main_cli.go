package main

import (
	"./artnet"
	"./cfread"
	//	"./cmdinterface"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func logger(queue chan string) {
	for {
		str := <-queue
		fmt.Fprintf(os.Stdout, "\r")
		const layout = "Jan 2, 2006 at 3:04:05pm (MST)"
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

var configName string

func runApp() {

	flag.StringVar(&configName, "conf", "", "Config.json")
	flag.StringVar(&configName, "c", "", "Config.json (short version)")
	flag.Parse()
	conf := cfread.ParseConfig(configName)
	conf.Status()

	go logger(conf.ChLog.Queue)
	go cli()
	an := artnet.Artnet{}
	err := an.Setup(&conf)
	if err == nil {
		an.Connect()
	} else {
		conf.ChLog.Log(fmt.Sprintf("Error setup ArtNet: %s", err.Error()))
		conf.ChLog.Log("Pause for 1 second...")
		time.Sleep(1 * time.Second)
	}
}

func main() {
	runApp()
}
