package cfread

import (
	"../cmdinterface"
	"encoding/json"
	"fmt"
	"os"
)

type CFReader struct {
	Interface  string   // Interface for ArtNet
	LongName   string   // Long name of Device
	ShortName  string   // Short name of Device
	OEMString  string   // OEM Number (Manufacturer)
	Users      []string // STUB
	Groups     []string // STUB
	Controller bool
	ChLog      cmdinterface.CmdIface
	ProgVers   uint16
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
		}
	}
	return cfr
}
