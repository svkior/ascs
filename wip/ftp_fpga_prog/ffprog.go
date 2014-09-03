package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jlaffaye/ftp"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

var version string

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

const (
	IP_DEFAULT     = "192.168.0.136"
	BIT_DEFAULT    = "./*.bit"
	NEED_DEFAULT   = false
	REPEAT_DEFAULT = false
	DEST_DEFAULT   = "/root/firmware1.bit"
	REPROG_DEFAULT = 0
)

type Config struct {
	DestIp     string // IP адрес устройства
	BitFile    string // паттерн, который нужно зашивать
	NeedProg   bool   // Нужно ли выполнять команду TODO: Сделать параметр с названием команды, которую нужно выполнить
	NeedRepeat bool   // Нужно ли сканировать директорию на предмет новых файлов. TODO: нужно задавать время
	DestFile   string // файл, который нужно сохранять
	Reprog     int    // Время в секундах для автоматической перезаливки программы
}

var config Config

func configuringProgam() {

	printVersion := flag.Bool("v", false, "Print version and exit")
	writeYamlPtr := flag.String("wyaml", "", "Write config in YAML format and exit")
	writeJsonPtr := flag.String("wjson", "", "Write config in JSON format and exit")
	debugPtr := flag.String("debug", "", "Working use cases are:\n     -debug=config  : Prints current config and exit normally")
	configPtr := flag.String("json", "", "Read config from JSON config file")
	yamlPtr := flag.String("yaml", "", "Read config from YAML config file")
	ipPtr := flag.String("ip", IP_DEFAULT, "DevBoard ip address")
	bitPtr := flag.String("bit", BIT_DEFAULT, "full path of bit file")
	needPtr := flag.Bool("prog", NEED_DEFAULT, "Download and program firmware")
	repeatPtr := flag.Bool("repeat", REPEAT_DEFAULT, "Wait for 1 second, rescan and flash if changed")
	destPtr := flag.String("dest", DEST_DEFAULT, "Specify destination for file name")
	reprogPtr := flag.Int("reprog", REPROG_DEFAULT, "time in seconds to auto reprog")

	flag.Parse()

	if len(version) < 1 {
		version = "DEVELOPMENT"
	}

	fmt.Println("AT91SAM9+FPGA DevBoard, Fpga programmer Version: ", version)

	// Если попросили версию, то после этого выходим
	if *printVersion {
		os.Exit(0)
	}

	config.DestIp = IP_DEFAULT
	config.BitFile = BIT_DEFAULT
	config.NeedProg = NEED_DEFAULT
	config.NeedRepeat = REPEAT_DEFAULT
	config.DestFile = DEST_DEFAULT
	config.Reprog = REPROG_DEFAULT

	if (len(*configPtr) > 0) && (len(*yamlPtr) > 0) {
		fmt.Println("Макс, не дури. Используй только один тип конфигов")
		os.Exit(99)
	}

	if len(*configPtr) > 0 {

		//log.Println("Reading config from file ", *configPtr)

		jdata, err := ioutil.ReadFile(*configPtr)
		if err != nil {
			fmt.Printf("Error opening config file : %s", err.Error())
			os.Exit(1)
		}

		err = json.Unmarshal(jdata, &config)

		if err != nil {
			fmt.Println("error decoding config file: ", err.Error())
			os.Exit(2)
		}
	}

	if len(*yamlPtr) > 0 {
		yfile, err := ioutil.ReadFile(*yamlPtr)
		if err != nil {
			fmt.Println("error opening YAML config file: ", err.Error())
			os.Exit(3)
		}

		err = yaml.Unmarshal(yfile, &config)
		if err != nil {
			fmt.Println("error decoding YAML config file:", err.Error())
		}

	}

	// Сверху накладываем дополнительные ключи в командной строке

	if *ipPtr != IP_DEFAULT {
		config.DestIp = *ipPtr
	}
	if *bitPtr != BIT_DEFAULT {
		config.BitFile = *bitPtr
	}
	if *needPtr != NEED_DEFAULT {
		config.NeedProg = *needPtr
	}
	if *repeatPtr != REPEAT_DEFAULT {
		config.NeedRepeat = *repeatPtr
	}
	if *destPtr != DEST_DEFAULT {
		config.DestFile = *destPtr
	}
	if *reprogPtr != REPROG_DEFAULT {
		config.Reprog = *reprogPtr
	}

	// Записываем конфиги обратно на диск
	var needExit bool

	if len(*writeYamlPtr) > 0 {
		log.Println("Writing YAML config to: ", *writeYamlPtr)
		encodedY, err := yaml.Marshal(&config)
		if err != nil {
			fmt.Println("Error encoding YAML: ", err.Error())
			os.Exit(6)
		}

		err = ioutil.WriteFile(*writeYamlPtr, encodedY, 0644)
		if err != nil {
			fmt.Println("Error writing YAML: ", err.Error())
			os.Exit(7)
		}

		needExit = true
	}

	if len(*writeJsonPtr) > 0 {
		log.Println("Writing JSPN config to: ", *writeJsonPtr)
		b, err := json.Marshal(config)
		if err != nil {
			fmt.Println("Error encoding JSON: ", err.Error())
			os.Exit(8)
		}

		err = ioutil.WriteFile(*writeJsonPtr, b, 0644)
		if err != nil {
			fmt.Println("Error writing JSON: ", err.Error())
		}

		needExit = true

	}

	if len(*debugPtr) > 0 {
		switch *debugPtr {
		case "config":
			fmt.Println("Program config:")
			fmt.Println("\tDestIp     = ", config.DestIp)
			fmt.Println("\tBitFile    = ", config.BitFile)
			fmt.Println("\tNeedProg   = ", config.NeedProg)
			fmt.Println("\tNeedRepeat = ", config.NeedRepeat)
			fmt.Println("\tDestFile   = ", config.DestFile)
			fmt.Println("\tReprog     = ", config.Reprog)
			os.Exit(0)
		default:
			fmt.Println("Unknown option for -debug flag: ", *debugPtr)
			os.Exit(4)
		}
	}

	if needExit {
		os.Exit(0)
	}

}

func main() {

	var counter int
	var modTime time.Time
	var lastName string
	var telNet net.Conn
	var telBuf *bufio.Reader

	var reprogCounter int

	configuringProgam()

	ftpAddr := config.DestIp + ":21"
	telnetAddr := config.DestIp + ":23"

	firstTime := true
	modified := false

	for {
		files, err := filepath.Glob(config.BitFile)
		if err != nil {
			log.Fatal(err.Error())
		}

		if len(files) < 1 {
			log.Fatal("Could not find files matching: ", config.BitFile)
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
			modified = false
			log.Println("Selecting file: ", lastName)

			uploadFile(&ftpAddr, &lastName, &config.DestFile)

			if config.NeedProg {

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
				telNet.Write([]byte("fpga_loader " + config.DestFile + "\n"))
				log.Println("Programming done")
				TelnetWaitCommand(telBuf)
			}

		}
		if config.NeedRepeat == false {
			log.Println("Exiting...")
			break
		} else {
			counter += 1
			//log.Printf("<%04d> Sleeping for 5 seconds...\n", counter)
			time.Sleep(1 * time.Second)
			//log.Println("Wake Up!")
			if config.Reprog > 0 {
				if reprogCounter == 0 {
					reprogCounter = config.Reprog
					modified = true
				} else {
					reprogCounter -= 1
				}
			}
		}

	}
}
