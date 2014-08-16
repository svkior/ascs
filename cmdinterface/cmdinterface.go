package cmdinterface

type CmdIface struct {
	Queue chan string
}

func (c *CmdIface) Init() {
	c.Queue = make(chan string)
}

func (c *CmdIface) Log(log string) {
	c.Queue <- log
}
