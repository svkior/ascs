package artnet

import (
	"fmt"
	"net"
	//	"time"
)

type Artnet struct {
	universe  uint16
	broadcast net.IP
	localIp   net.IP
}

func (a *Artnet) Setup(universe uint16) {
	a.universe = universe
}

func (a *Artnet) Connect(ipAddr string) {
	addr := net.ParseIP(ipAddr)
	if addr == nil {
		fmt.Println("Invalid Address: %s", ipAddr)
		return
	}
	fmt.Printf("The address is: %s\n", addr.String())
	a.localIp = addr

	mask := addr.DefaultMask()
	network := addr.Mask(mask)

	a.broadcast = net.IPv4(
		network[0]|^mask[0],
		network[1]|^mask[1],
		network[2]|^mask[2],
		network[3]|^mask[3],
	)

	fmt.Printf("The broadcast address is: %s\n", a.broadcast.String())

}

func (a *Artnet) SendPacket(buf [512]byte, Len int) {

	service := a.broadcast.String() + ":6454"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		fmt.Printf("Error in resolve broadcast: %s", err.Error())
		return
	}

	serviceRet := a.localIp.String() + ":6454"
	srvAddr, err := net.ResolveUDPAddr("udp4", serviceRet)
	if err != nil {
		fmt.Printf("Error in resolve server: %s", err.Error())
		return
	}

	conn, err := net.DialUDP("udp", srvAddr, udpAddr)
	if err != nil {
		fmt.Printf("Error Dial UDP: %s", err.Error())
		return
	}

	_, err = conn.Write(buf[0:Len])
	if err != nil {
		fmt.Printf("Error Write to UDP: %s", err.Error())
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
