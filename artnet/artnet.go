package artnet

import (
	"../cfread"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Status1 struct {
	UBEA      bool // true if Present
	RDM       bool // true if RDM
	BootROM   bool // true if booted from ROM
	PortAddr  byte //  Port Address "unknown" 00 - unknown, "front" 01 - set from front panel, "net" 10 - programmed from network, "unused" 11 - not used
	Indicator byte // "unknown" 00- no indicators, "locate" 01- locatemode, "mute" 10- mute mode, "normal" 11- normal mode
}

func CreateStatus1() Status1 {
	var s Status1
	s.UBEA = false
	s.RDM = false
	s.BootROM = false
	s.PortAddr = 3
	s.Indicator = 3
	return s
}

func (s *Status1) SetPortAddr(v string) {
	switch v {
	case "unknown":
		s.PortAddr = 0
	case "front":
		s.PortAddr = 1
	case "net":
		s.PortAddr = 2
	case "unused":
		s.PortAddr = 3
	}
}

func (s *Status1) SetIndicator(v string) {
	switch v {
	case "unknown":
		s.Indicator = 0
	case "locate":
		s.Indicator = 1
	case "mute":
		s.Indicator = 2
	case "normal":
		s.Indicator = 3
	}
}

func (s *Status1) GetStatus1() byte {
	var rv byte
	if s.UBEA {
		rv |= 1
	}
	if s.RDM {
		rv |= 2
	}
	if s.BootROM {
		rv |= 4
	}
	rv |= s.PortAddr << 4
	rv |= s.Indicator << 6
	return rv
}

type Status2 struct {
	hasWeb      bool
	ipDhcp      bool
	dhcpCapable bool
}

func CreateStatus2() Status2 {
	var s Status2
	s.hasWeb = true
	s.ipDhcp = false
	s.dhcpCapable = true
	return s
}

func (s *Status2) GetStatus2() byte {
	var status2 byte
	if s.hasWeb {
		status2 |= 1
	}
	if s.ipDhcp {
		status2 |= 2
	}
	if s.dhcpCapable {
		status2 |= 4
	}
	status2 |= 8 // ArtNet3
	return status2
}

func (s *Status2) DHCP(flag bool) {
	s.ipDhcp = flag
}

type Artnet struct {
	universe    uint16
	broadcast   net.IP
	localIp     net.IP
	macAddr     net.HardwareAddr
	OEM         uint16
	ESTA        uint16
	Stat1       Status1
	Stat2       Status2
	conf        *cfread.CFReader
	conn        *net.UDPConn
	PoolCounter uint16
	StatusCode  uint16
	StatusMsg   string
}

func (a *Artnet) SetStatus(stat string) {
	switch stat {
	case "RcDebug":
		a.StatusCode = 0x0000
	case "RcPowerOk":
		a.StatusCode = 0x0001
	case "RcPowerFail":
		a.StatusCode = 0x0002
	case "RcSocketWr1":
		a.StatusCode = 0x0003
	case "RcParseFail":
		a.StatusCode = 0x0004
	case "RcUdpFail":
		a.StatusCode = 0x0005
	case "RcShNameOk":
		a.StatusCode = 0x0006
	case "RcLoNameOk":
		a.StatusCode = 0x0007
	case "RcDmxError":
		a.StatusCode = 0x0008
	case "RcDmxUdpFull":
		a.StatusCode = 0x0009
	case "RcDmxRxFull":
		a.StatusCode = 0x000a
	case "RcSwitchErr":
		a.StatusCode = 0x000b
	case "RcConfigErr":
		a.StatusCode = 0x000c
	case "RcDmxShort":
		a.StatusCode = 0x000d
	case "RcFirmwareFail":
		a.StatusCode = 0x000e
	case "RcUserFail":
		a.StatusCode = 0x000f
	}
}

func (a *Artnet) GetNodeReport() string {
	a.PoolCounter++
	return fmt.Sprintf("#%04x [%05d] %s", a.StatusCode, a.PoolCounter, a.StatusMsg)
}

func (a *Artnet) Logf(format string, i ...interface{}) {
	a.conf.ChLog.Log(fmt.Sprintf(format, i))
}

func (a *Artnet) Log(i ...interface{}) {
	a.conf.ChLog.Log(fmt.Sprint(i))
}

func (a *Artnet) Setup(conf *cfread.CFReader) error {

	a.universe = 0x1
	a.conf = conf

	iface, err := net.InterfaceByName(a.conf.Interface)
	if err != nil {
		return err
	}
	//a.Logf("Got Interface: %s", iface)

	// Got MAC Address
	a.macAddr = iface.HardwareAddr
	a.Logf("Got Interface with MAC: %s", a.macAddr.String())

	iAddrs, err := iface.Addrs()
	if err != nil {
		return err
	}

	for _, addr := range iAddrs {
		ip, ipNet, err := net.ParseCIDR(addr.String())
		if err == nil {
			if ip.To4() != nil {
				mask := ipNet.Mask
				network := ip.Mask(mask)
				a.localIp = ip
				a.broadcast = net.IPv4(
					network[0]|^mask[0],
					network[1]|^mask[1],
					network[2]|^mask[2],
					network[3]|^mask[3],
				)
				a.Logf("Addr :     %s", ip.String())
				a.Logf("Network :  %s", network.String())
				a.Logf("Mask:      %s", mask.String())
				a.Logf("Broadcast: %s", a.broadcast.String())

				oem, err := strconv.ParseUint(a.conf.OEMString, 0, 16)
				if err != nil {
					return err
				}
				a.OEM = uint16(oem)
				a.Logf("OEM Number: %d", a.OEM)

				esta, err := strconv.ParseUint(a.conf.ESTAString, 0, 16)
				if err != nil {
					return err
				}
				a.ESTA = uint16(esta)
				a.Logf("ESTA Number: %d", a.ESTA)

				a.Stat1 = CreateStatus1()
				a.Stat1.UBEA = a.conf.UBEA
				a.Stat1.RDM = a.conf.RDM
				a.Stat1.BootROM = a.conf.BootROM
				a.Stat1.SetPortAddr(a.conf.PortAddr)
				a.Stat1.SetIndicator(a.conf.Indicator)

				a.Stat2 = CreateStatus2()
				a.Stat2.DHCP(a.conf.DHCP)

				a.SetStatus("RcPowerOk")
				a.StatusMsg = "Everything is OK"

				return nil
			}
		}
	}

	return errors.New("Can`t find IPv4 Address")
}

func (a *Artnet) Connect() {
	a.Logf("The broadcast address is: %s", a.broadcast.String())
	service := ":6454"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		a.Logf("Error Resovle : %s", err.Error())
		return
	}

	a.conn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		a.Logf("Error Listen on : %s", err.Error())
		return
	}

	if a.conf.Controller {
		go a.Pooler()
	}
	a.ListenArtnet()

}

func (a *Artnet) Pooler() {
	ticker := time.NewTicker(time.Second * 8)
	go func() {
		for _ = range ticker.C {
			a.SendArtPoll()
		}
	}()
}

func (a *Artnet) ParsePacket(buf [1024]byte, addr *net.UDPAddr, n int) {
	//a.Logf("Read ArtNet: %d bytes", n))

	id := string(buf[0:8])
	if (id[0] != 'A') || (id[1] != 'r') {
		a.Logf("Id: %s", id)
		return
	}
	var OpCode uint
	OpCode = uint(buf[8]) + uint(buf[9])*256
	//a.Logf("Opcode: 0x%04x", OpCode))
	switch OpCode {
	case 0x2000:
		a.Log("ArtPoolRequest: IN Progress")
		protVer := uint(buf[10])*256 + uint(buf[11])
		if protVer < 14 {
			a.Logf("ProtVer: %d is lower than", protVer)
			return
		}
		a.SendArtPollReply(addr)
	case 0x2100:
		a.Log("ArtPoolReply: NOT REALIZED")
	case 0x2300:
		a.Log("OpDiagData: NOT REALIZED")
	case 0x2400:
		a.Log("OpCommand: NOT REALIZED")
	case 0x5000:
		a.Log("OpOutput/OpDmx: NOT REALIZED")
	case 0x5100:
		a.Log("OpNzs: NOT REALIZED")
	case 0x6000:
		a.Log("OpAddress: NOT REALIZED")
	case 0x7000:
		a.Log("OpInput: NOT REALIZED")
	case 0x8000:
		a.Log("OpTodRequest: NOT REALIZED")
	case 0x8100:
		a.Log("OpTodData: NOT REALIZED")
	case 0x8200:
		a.Log("OpTodControl: NOT REALIZED")
	case 0x8300:
		a.Log("OpRdm: NOT REALIZED")
	case 0x8400:
		a.Log("OpRdmSub: NOT REALIZED")
	case 0xa010:
		a.Log("OpVideoSetup: NOT REALIZED")
	case 0xa020:
		a.Log("OpVideoPalette: NOT REALIZED")
	case 0xa040:
		a.Log("OpVideoData: NOT REALIZED")
	case 0xf000:
		a.Log("OpMacMaster: NOT REALIZED")
	case 0xf100:
		a.Log("OpMacSlave: NOT REALIZED")
	case 0xf200:
		a.Log("OpFirmwareMaster: NOT REALIZED")
	case 0xf300:
		a.Log("OpFirmwareReply: NOT REALIZED")
	case 0xf400:
		a.Log("OpFileTnMaster: NOT REALIZED")
	case 0xf500:
		a.Log("OpFileFnMaster: NOT REALIZED")
	case 0xf600:
		a.Log("OpFileFnReply: NOT REALIZED")
	case 0xf800:
		a.Log("OpIpProg: NOT REALIZED")
	case 0xf900:
		a.Log("OpIpProgReply: NOT REALIZED")
	case 0x9000:
		a.Log("OpMedia: NOT REALIZED")
	case 0x9100:
		a.Log("OpMediaPatch: NOT REALIZED")
	case 0x9200:
		a.Log("OpMediaControl: NOT REALIZED")
	case 0x9300:
		a.Log("OpMediaControlReply: NOT REALIZED")
	case 0x9700:
		a.Log("OpTimeCode: NOT REALIZED")
	case 0x9800:
		a.Log("OpTimeSync: NOT REALIZED")
	case 0x9900:
		a.Log("OpTrigger: NOT REALIZED")
	case 0x9a00:
		a.Log("OpDirectory: NOT REALIZED")
	case 0x9b00:
		a.Log("OpDirectoryReply: NOT REALIZED")
	}

}

func (a *Artnet) ListenArtnet() {

	for {
		var buf [1024]byte
		n, addr, err := a.conn.ReadFromUDP(buf[0:])
		if err != nil {
			a.Logf("Error Read From UDP: %s", err.Error())
			return
		}
		//a.Logf("Read ArtNet from %s : %d bytes", addr.String(), n))
		a.ParsePacket(buf, addr, n)
	}
}

func (a *Artnet) SendPacket(buf [512]byte, Len int) {

	service := a.broadcast.String() + ":6454"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		a.Logf("Error in resolve broadcast: %s", err.Error())
		return
	}

	_, err = a.conn.WriteToUDP(buf[0:Len], udpAddr)
	if err != nil {
		a.Logf("Error Write to UDP: %s", err.Error())
		return
	}
	/*
		var retbuf [512]byte
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, err := conn.Read(retbuf[0:])
		if err != nil {
			fmt.Printf("Error read from UDP: %s", err.Error())
			return
		}

		fmt.Printf("Received: %d bytes\n", n)
	*/
}

func (a *Artnet) SendArtPoll() {
	var buf [512]byte

	OpCode := 0x2000
	ProtVerHi := 00
	ProtVerLo := 14

	TalkToMe := 0
	Priority := 0

	buf[0] = byte('A')           // A
	buf[1] = byte('r')           // r
	buf[2] = byte('t')           // t
	buf[3] = byte('-')           // -
	buf[4] = byte('N')           // N
	buf[5] = byte('e')           // e
	buf[6] = byte('t')           // t
	buf[7] = 0                   // 0x00
	buf[8] = byte(OpCode & 0xff) // OpCode[0]
	buf[9] = byte(OpCode >> 8)   // OpCode[1]
	buf[10] = byte(ProtVerHi)    // ProtVerHi
	buf[11] = byte(ProtVerLo)    // ProtVerLo
	buf[12] = byte(TalkToMe)     // TalkToMe
	buf[13] = byte(Priority)     // Priority
	bufLen := 14

	a.SendPacket(buf, bufLen)

}

func getLow(word uint16) byte {
	return byte(word & 255)
}

func getHi(word uint16) byte {
	return byte((word >> 8) & 255)
}

func (a *Artnet) SendArtPollReply(addr *net.UDPAddr) {

	//a.Logf("Fake SendArtPollReply: %s", addr.String()))
	var buf [512]byte

	OpCode := uint16(0x2100)

	// idAddress := a.localIP
	// port

	buf[0] = byte('A')                // A
	buf[1] = byte('r')                // r
	buf[2] = byte('t')                // t
	buf[3] = byte('-')                // -
	buf[4] = byte('N')                // N
	buf[5] = byte('e')                // e
	buf[6] = byte('t')                // t
	buf[7] = 0                        // 0x00
	buf[8] = getLow(OpCode)           // OpCode[0]
	buf[9] = getHi(OpCode)            // OpCode[1]
	buf[10] = a.localIp[0]            // IPV4 [0]
	buf[11] = a.localIp[1]            // IPV4 [1]
	buf[12] = a.localIp[2]            // IPV4 [2]
	buf[13] = a.localIp[3]            // IPV4 [3]
	buf[14] = 0x36                    // IP Port Low
	buf[15] = 0x19                    // IP Port Hi
	buf[16] = getHi(a.conf.ProgVers)  // High byte of Version
	buf[17] = getLow(a.conf.ProgVers) // Low byte of Version
	buf[18] = 0x00                    // NetSwitch
	buf[19] = 0x00                    // Net Sub Switch
	buf[20] = getHi(a.OEM)            // OEMHi
	buf[21] = getLow(a.OEM)           // OEMLow
	buf[22] = a.conf.UBEAVer          // Ubea Version
	buf[23] = a.Stat1.GetStatus1()    // Status1
	buf[24] = getLow(a.ESTA)          // ESTA LO
	buf[25] = getHi(a.ESTA)           // ESTA HI

	for i, c := range a.conf.ShortName {
		if i < 16 {
			buf[26+i] = byte(c)
		} else {
			break
		}
	}
	for i, c := range a.conf.LongName {
		if i < 63 {
			buf[44+i] = byte(c)
		} else {
			break
		}
	}
	NodeReport := a.GetNodeReport()
	for i, s := range NodeReport {
		if i < 63 {
			buf[108+i] = byte(s)
		}
	}
	// 108 + 64
	//buf[108 : 108+64] = NodeReport
	buf[172] = 0 // NumPorts Hi
	buf[173] = 0 // NumPorts Lo
	buf[174] = 0 // Port 0 Type
	buf[175] = 0 // Port 1 Type
	buf[176] = 0 // Port 2 Type
	buf[177] = 0 // Port 3 Type
	buf[178] = 0 // GoodInput 0
	buf[179] = 0 // GoodInput 1
	buf[180] = 0 // GoodInput 2
	buf[181] = 0 // GoodInput 3
	buf[182] = 0 // GoodOutput 0
	buf[183] = 0 // GoodOutput 1
	buf[184] = 0 // GoodOutput 2
	buf[185] = 0 // GoodOutput 3
	buf[186] = 0 // SwIn 0
	buf[187] = 0 // SwIn 1
	buf[188] = 0 // SwIn 2
	buf[189] = 0 // SwIn 3
	buf[190] = 0 // SwOut 0
	buf[191] = 0 // SwOut 1
	buf[192] = 0 // SwOut 2
	buf[193] = 0 // SwOut 3
	buf[194] = 0 // SwVideo
	buf[195] = 0 // SwMacro
	buf[196] = 0 // SwRemote
	buf[197] = 0 // Spare
	buf[198] = 0 // Spare
	buf[199] = 0 // Spare
	buf[200] = 0 // Style
	// MAC ADDRESS
	buf[201] = a.macAddr[0] // MAC HI
	buf[202] = a.macAddr[1] // MAC
	buf[203] = a.macAddr[2] // MAC
	buf[204] = a.macAddr[3] // MAC
	buf[205] = a.macAddr[4] // MAC
	buf[206] = a.macAddr[5] // MAC LO

	buf[207] = a.localIp[0] // BIND IP 0
	buf[208] = a.localIp[1] // BIND IP 1
	buf[209] = a.localIp[2] // BIND IP 2
	buf[210] = a.localIp[3] // BIND IP 3
	buf[211] = 0            // BInd Index

	buf[212] = a.Stat2.GetStatus2() // Status2
	// 212 + 26 = 238
	bufLen := 238

	a.SendPacket(buf, bufLen)

}
