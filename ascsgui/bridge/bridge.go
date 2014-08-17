package bridge

import (
	"../../cmdinterface"
	"fmt"
	"gopkg.in/qml.v0"
	"net"
	"strings"
	"time"
)

type LogRecord struct {
	record string
	rType  string
	color  string
}

type Logs struct {
	list []LogRecord
	Len  int
}

func (l *Logs) Add(r LogRecord) {
	l.list = append(l.list, r)
	l.Len = len(l.list)
	qml.Changed(l, &l.Len)
	if l.Len > 100 {
		l.list = l.list[1:]
		l.Len = len(l.list)
		qml.Changed(l, &l.Len)
	}
}

func (l *Logs) Type(index int) string {
	return l.list[index].rType
}

func (l *Logs) Record(index int) string {
	return l.list[index].record
}

func (l *Logs) Color(index int) string {
	return l.list[index].color
}

type Bridge struct {
	started   bool
	ipName    string
	BLogs     *Logs
	listView  qml.Object
	inputText qml.Object
}

func (b *Bridge) SetupVars(root qml.Object) {
	b.listView = root.ObjectByName("logview")
	b.inputText = root.ObjectByName("inputtext1")
}

func (b *Bridge) ColorLog(log string, color string) {
	const layout = "Jan 2, 2006 at 3:04pm (MST)"
	strSplits := strings.Split(log, "\n")
	var firstTime bool = true
	for _, s := range strSplits {
		if len(s) > 0 {
			if firstTime {
				firstTime = false
				s = time.Now().Format(layout) + " : " + s
			}
			b.BLogs.Add(LogRecord{record: s, rType: "info", color: color})
		}
	}
}

func (b *Bridge) StartLogger(logger *cmdinterface.CmdIface) {
	var logLine string
	for {
		logLine = <-logger.Queue
		b.Log(logLine)
	}
}

func (b *Bridge) Log(log string) {
	b.ColorLog(log, "gold")
}

func (b *Bridge) ClientLog(log string) {
	b.ColorLog(log, "green")
}

func (b *Bridge) ErrorLog(log string) {
	b.BLogs.Add(LogRecord{record: log, rType: "error"})
}

func (b *Bridge) HandleClick() {
	b.ClientLog("Clicked")
	go b.UDPClient()
}

func (b *Bridge) UDPClient() {
	service := "127.0.0.1:1200"

	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		b.ClientLog(fmt.Sprintf("Error resolveUDPAddr: %s", err.Error()))
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		b.ClientLog(fmt.Sprintf("Error DialUDP: %s", err.Error()))
		return
	}

	_, err = conn.Write([]byte("anything"))
	if err != nil {
		b.ClientLog(fmt.Sprintf("Error Write to UDP: %s", err.Error()))
		return
	}

	var buf [512]byte
	n, err := conn.Read(buf[0:])
	if err != nil {
		b.ClientLog(fmt.Sprintf("Error Read from UDP: %s", err.Error()))
		return
	}

	b.ClientLog(fmt.Sprintf("Got: %s", string(buf[0:n])))
}

func (b *Bridge) UDPDayTimeServer() {
	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		b.Log(fmt.Sprintf("Error resolveUDPAddr: %s", err.Error()))
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		b.Log(fmt.Sprintf("Error ListenUDP: %s", err.Error()))
		return
	}

	for {
		b.handleClient(conn)
	}

}

func (b *Bridge) handleClient(conn *net.UDPConn) {

	var buf [512]byte
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	b.Log(fmt.Sprintf("Request from %s", addr.String()))
	daytime := time.Now().String()
	b.Log(fmt.Sprintf("Send time: %s", daytime))
	conn.WriteToUDP([]byte(daytime), addr)
}
