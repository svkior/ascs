package artnet

import (
	"../cfread"
	"errors"
	"fmt"
	"net"
	"time"
)

type Artnet struct {
	universe  uint16
	broadcast net.IP
	localIp   net.IP
	macAddr   net.HardwareAddr
	conf      *cfread.CFReader
	conn      *net.UDPConn
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
		a.Logf("Error get interface by name: %s", err.Error())
		return err
	}
	//a.Logf("Got Interface: %s", iface)

	// Got MAC Address
	a.macAddr = iface.HardwareAddr
	a.Logf("Got Interface with MAC: %s", a.macAddr.String())

	iAddrs, err := iface.Addrs()
	if err != nil {
		a.Logf("Error get interface addresses: %s", err.Error())
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
	ticker := time.NewTicker(time.Second * 5)
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
	OpCode := uint(buf[8]) + uint(buf[9])*256
	//a.Logf("Opcode: 0x%04x", OpCode))
	switch OpCode {
	case 0x2000:
		a.Log("ArpPoolRequest: IN Progress")
		protVer := uint(buf[10])*256 + uint(buf[11])
		if protVer < 14 {
			a.Logf("ProtVer: %d is lower than", protVer)
			return
		}
		a.SendArtPollReply(addr)
	case 0x2100:
		a.Log("ArpPoolReply: NOT REALIZED")
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

func (a *Artnet) SendArtPollReply(addr *net.UDPAddr) {

	//a.Logf("Fake SendArtPollReply: %s", addr.String()))
	var buf [512]byte

	OpCode := 0x2100

	// idAddress := a.localIP
	// port

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
	buf[10] = a.localIp[0]       // IPV4 [0]
	buf[11] = a.localIp[1]       // IPV4 [1]
	buf[12] = a.localIp[2]       // IPV4 [2]
	buf[13] = a.localIp[3]       // IPV4 [3]
	buf[14] = 0x36               // IP Port Low
	buf[15] = 0x19               // IP Port Hi
	buf[16] = 0x00               // High byte of Version
	buf[17] = 0x00               // Low byte of Version
	buf[18] = 0x00               // NetSwitch
	buf[19] = 0x00               // Net Sub Switch
	buf[20] = 0x04               // OEMHi
	buf[21] = 0x30               // OEMLow
	buf[22] = 0x00               // Ubea Version
	buf[23] = 0x00               // Status1
	buf[24] = byte('p')          // ESTA LO
	buf[25] = byte('z')          // ESTA HI

	for i, c := range a.conf.ShortName {
		if i < 16 {
			buf[26+i] = byte(c)
		} else {
			break
		}
	}
	/*
		buf[26] = byte('A') // Short Name 0
		buf[27] = byte('A') // Short Name 1
		buf[28] = byte('A') // Short Name 2
		buf[29] = byte('A') // Short Name 3
		buf[30] = byte('A') // Short Name 4
		buf[31] = byte('A') // Short Name 5
		buf[32] = byte('A') // Short Name 6
		buf[33] = byte('A') // Short Name 7
		buf[34] = byte('A') // Short Name 8
		buf[35] = byte('A') // Short Name 9
		buf[36] = byte('A') // Short Name 10
		buf[37] = byte('A') // Short Name 11
		buf[38] = byte('A') // Short Name 12
		buf[39] = byte('A') // Short Name 13
		buf[40] = byte('A') // Short Name 14
		buf[41] = byte('A') // Short Name 15
		buf[42] = byte('A') // Short Name 16
	*/
	buf[43] = 0 // Short Name END MUST BE 0
	longName := []byte(a.conf.LongName)
	for i, s := range longName {
		if i < 63 {
			buf[44+i] = s
		} else {
			break
		}
	}
	//buf[44 : 44+64] = LongName // Long Name 0
	// 44 + 64 =
	NodeReport := []byte("Node Report")
	for i, s := range NodeReport {
		buf[108+i] = s
	}
	// 108 + 64
	//buf[108 : 108+64] = NodeReport
	buf[172] = 0    // NumPorts Hi
	buf[173] = 0    // NumPorts Lo
	buf[174] = 0    // Port 0 Type
	buf[175] = 0    // Port 1 Type
	buf[176] = 0    // Port 2 Type
	buf[177] = 0    // Port 3 Type
	buf[178] = 0    // GoodInput 0
	buf[179] = 0    // GoodInput 1
	buf[180] = 0    // GoodInput 2
	buf[181] = 0    // GoodInput 3
	buf[182] = 0    // GoodOutput 0
	buf[183] = 0    // GoodOutput 1
	buf[184] = 0    // GoodOutput 2
	buf[185] = 0    // GoodOutput 3
	buf[186] = 0    // SwIn 0
	buf[187] = 0    // SwIn 1
	buf[188] = 0    // SwIn 2
	buf[189] = 0    // SwIn 3
	buf[190] = 0    // SwOut 0
	buf[191] = 0    // SwOut 1
	buf[192] = 0    // SwOut 2
	buf[193] = 0    // SwOut 3
	buf[194] = 0    // SwVideo
	buf[195] = 0    // SwMacro
	buf[196] = 0    // SwRemote
	buf[197] = 0    // Spare
	buf[198] = 0    // Spare
	buf[199] = 0    // Spare
	buf[200] = 0    // Style
	buf[201] = 0xff // MAC HI
	buf[202] = 0xff // MAC
	buf[203] = 0xff // MAC
	buf[204] = 0xff // MAC
	buf[205] = 0xff // MAC
	buf[206] = 0xff // MAC LO
	buf[207] = 0x0  // BIND IP 0
	buf[208] = 0x0  // BIND IP 1
	buf[209] = 0x0  // BIND IP 2
	buf[210] = 0x0  // BIND IP 3
	buf[211] = 0    // BInd Index
	buf[212] = 0    // Status2
	// 212 + 26 = 238
	bufLen := 238

	a.SendPacket(buf, bufLen)

}
