package cfread

import (
	"../cmdinterface"
	"./chnls"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type CFReader struct {
	NoConsole  bool   // if daemon mode
	Interface  string // Interface for ArtNet
	LongName   string // Long name of Device
	ShortName  string // Short name of Device
	OEMString  string // OEM Number
	ESTAString string // Manufacturer code
	// Status 2
	DHCP bool // true if DHCP capable

	Channels []chnls.Channel

	// Status 1
	UBEA      bool   // true if Present
	RDM       bool   // true if RDM
	BootROM   bool   // true if booted from ROM
	PortAddr  string //  Port Address "unknown" - unknown, "front" - set from front panel, "net" - programmed from network, "unused" - not used
	Indicator string // "unknown" - no indicators, "locate" - locatemode, "mute" - mute mode, "normal" - normal mode

	UBEAVer byte // Версия UBEA

	Controller bool
	ChLog      cmdinterface.CmdIface
	ProgVers   uint16

	Users  []string // STUB
	Groups []string // STUB
}

func (c *CFReader) Status() {
	fmt.Printf("Interface from Config: %s\n", c.Interface)
}

func ParseConfig(fname string) CFReader {
	cfr := CFReader{}
	cfr.ProgVers = 0x0002  // 0.002
	cfr.Interface = "eth0" // Default config for linux
	cfr.LongName = "Default TTS Device"
	cfr.ChLog.Init()

	if len(fname) < 1 {
		println("Creating default config")
	} else {
		println("Reading config from file : " + fname)
		file, err := os.Open(fname)
		if err != nil {
			fmt.Printf("Error opening config file : %s", err.Error())
			return cfr
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&cfr)
		if err != nil {
			fmt.Println("error:", err.Error())
			return cfr
		}
	}

	if cfr.Interface == "auto" {
		ints, err := net.Interfaces()
		if err != nil {
			fmt.Printf("Error get interfaces: %s", err.Error())
			return cfr
		}
		for n, intf := range ints {
			fmt.Printf("Int: %s", intf)
			if (intf.Flags & net.FlagUp) != 0 {
				fmt.Printf("Configured auto interface : %s", n, intf.Name)
				cfr.Interface = intf.Name
				fmt.Printf("[123]")
				break
			}
		}
	}

	return cfr
}
