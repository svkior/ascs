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
	"time"
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

func uploadFile(ftpAddr *string, fileName *string, destName *string) {

	var conn *ftp.ServerConn
	var err error
	//		log.Println("CONNECTION: ", conn)
	log.Printf("Connecting to ftp://%s\n", *ftpAddr)
	conn, err = ftp.Connect(*ftpAddr)
	if err != nil {
		log.Fatal("Error connecting to ", *ftpAddr, " : ", err.Error())
	}
	defer func() {
		log.Println("FTP Disconnecting from remote host")
		conn.Quit()
	}()

	log.Println("Logging as ftp:ftp")
	err = conn.Login("root", "1")
	if err != nil {
		log.Fatal("Error connecting to ", *ftpAddr, " : ", err.Error())
	}

	log.Println("Remove old firmware")
	err = conn.Delete(*destName)
	if err != nil {
		log.Println("Error delete file: ", err.Error())
	}

	// open bit file
	fi, err := os.Open(*fileName)
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
	log.Printf("Putting firmware %s\n", *fileName)
	err = conn.Stor(*destName, r)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	fmt.Println("\n<== AT91SAM9+FPGA DevBoard, Fpga programmer V:0.001a ==>\n")
	ipPtr := flag.String("ip", "192.168.0.136", "DevBoard ip address")
	bitPtr := flag.String("bit", "./*.bit", "full path of bit file")
	needPtr := flag.Bool("prog", false, "Download and program firmware")
	repeatPtr := flag.Bool("repeat", false, "Wait for 5 second, rescan and flash if changed")
	destPtr := flag.String("dest", "/root/firmware1.bit", "Specify destination for file name")

	flag.Parse()

	var counter int
	var modTime time.Time
	var lastName string
	var telNet net.Conn
	var telBuf *bufio.Reader

	ftpAddr := *ipPtr + ":21"
	telnetAddr := *ipPtr + ":23"

	firstTime := true
	modified := false

	for {
		modified = false
		files, err := filepath.Glob(*bitPtr)
		if err != nil {
			log.Fatal(err.Error())
		}

		if len(files) < 1 {
			log.Fatal("Could not find files matching: ", *bitPtr)
		}

		for _, file := range files {
			//			log.Println("==:", file)
			info, err := os.Stat(file)
			if err != nil {
				log.Fatal(err.Error())
			}
			if firstTime {
				modTime = info.ModTime()
				lastName = file
				firstTime = false
				modified = true

			} else {
				if info.ModTime().After(modTime) {
					modTime = info.ModTime()
					lastName = file
					modified = true
				}
			}
			//log.Println("File: ", file, " Mod time: ", info.ModTime())
		}

		// Если нашли новый файл
		if modified {
			log.Println("Selecting file: ", lastName)

			uploadFile(&ftpAddr, &lastName, destPtr)

			if *needPtr {

				log.Println(telNet)
				if true {
					log.Println("Connecting to telnet ", telnetAddr)
					telNet, err = net.Dial("tcp", telnetAddr)
					if err != nil {
						log.Fatal(err.Error())
					}
					defer func() {
						log.Println("Telnet Disconnecting from remote host")
						telNet.Close()
					}()
					telBuf = bufio.NewReader(telNet)
				} else {
					telNet.Write([]byte("\n"))
				}
				TelnetWaitCommand(telBuf)
				telNet.Write([]byte("sync\n"))
				TelnetWaitCommand(telBuf)
				telNet.Write([]byte("fpga_loader " + *destPtr + "\n"))
				log.Println("Programming done")
				TelnetWaitCommand(telBuf)
			}

		}
		if *repeatPtr == false {
			log.Println("Exiting...")
			break
		} else {
			counter += 1
			//log.Printf("<%04d> Sleeping for 5 seconds...\n", counter)
			time.Sleep(5 * time.Second)
			//log.Println("Wake Up!")
		}

	}
}
