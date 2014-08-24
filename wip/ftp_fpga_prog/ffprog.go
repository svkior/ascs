package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/jlaffaye/ftp"
	"log"
	"net"
	"os"
	"path/filepath"
)

func TelnetWaitCommand(telBuf *bufio.Reader) {
	for {
		symbol, err := telBuf.ReadByte()
		if err != nil {
			log.Fatal(err.Error())
		}
		if symbol == byte('#') {
			break
		}
	}
}

const (
	fwName = "firmware1.bit"
)

func main() {
	fmt.Println("AT91SAM9+FPGA DevBoard, Fpga programmer")
	ipPtr := flag.String("ip", "192.168.0.136", "DevBoard ip address")
	bitPtr := flag.String("bit", "./top_arm.bit", "full path of bit file")
	needPtr := flag.Bool("prog", false, "Download and program firmware")

	flag.Parse()

	files, err := filepath.Glob("*.bit")
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, file := range files {
		log.Println("==:", file)
	}

	log.Fatal("Help")

	log.Printf("Connecting to ftp://%s\n", *ipPtr)
	conn, err := ftp.Connect(*ipPtr + ":21")
	if err != nil {
		log.Fatal("Error connecting to ", *ipPtr, " : ", err.Error())
	}
	defer func() {
		log.Println("FTP Disconnecting from remote host")
		conn.Quit()
	}()

	log.Println("Logging as ftp:ftp")
	err = conn.Login("root", "1")
	if err != nil {
		log.Fatal("Error connecting to ", *ipPtr, " : ", err.Error())
	}

	log.Println("Remove old firmware")
	err = conn.Delete(fwName)
	if err != nil {
		log.Println("Error delete file: ", err.Error())
	}

	// open bit file
	fi, err := os.Open(*bitPtr)
	if err != nil {
		log.Fatal(err.Error())
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	r := bufio.NewReader(fi)
	log.Printf("Putting firmware %s\n", *bitPtr)
	err = conn.Stor(fwName, r)
	if err != nil {
		log.Fatal(err.Error())
	}

	if *needPtr {
		telNet, err := net.Dial("tcp", *ipPtr+":23")
		if err != nil {
			log.Fatal(err.Error())
		}
		defer func() {
			log.Println("Telnet Disconnecting from remote host")
			telNet.Close()
		}()
		telBuf := bufio.NewReader(telNet)

		TelnetWaitCommand(telBuf)
		telNet.Write([]byte("sync\n"))
		TelnetWaitCommand(telBuf)
		telNet.Write([]byte("fpga_loader /root/firmware.bit\n"))
		log.Println("Programming done")
		TelnetWaitCommand(telBuf)
	}

}
