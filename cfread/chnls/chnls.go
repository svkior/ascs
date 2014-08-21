package chnls

type Channel struct {
	Name string // Name of interface
	// "DMX-512", "MIDI", "Avab", "Colortran CMX", "ADB 62.5", "Art-Net"
	Enabled     bool // Enabled of interface
	IsArtnetIn  bool
	IsArtnetOut bool
}

func (c *Channel) GetPortType() byte {
	var out byte

	if c.Enabled {

		if c.IsArtnetIn {
			out |= 128
		}

		if c.IsArtnetOut {
			out |= 64
		}

		switch c.Name {
		case "DMX-512":
			out |= 0
		case "MIDI":
			out |= 1
		case "Avab":
			out |= 2
		case "Colortran CMX":
			out |= 3
		case "ADB 62.5":
			out |= 4
		case "Art-Net":
			out |= 5
		}

	}
	return out
}

func (c *Channel) SetPortType(b byte) {
	if b == 0 {
		c.Enabled = false
	} else {
		c.IsArtnetIn = (b&128 > 0)
		c.IsArtnetOut = (b&64 > 0)
		switch b & 63 {
		case 0:
			c.Name = "DMX-512"
		case 1:
			c.Name = "MIDI"
		case 2:
			c.Name = "Avab"
		case 3:
			c.Name = "Colortran CMX"
		case 4:
			c.Name = "ADB 62.5"
		case 5:
			c.Name = "Art-Net"
		}
	}
}

func (c *Channel) String() string {
	var s string
	if c.Enabled {
		s = c.Name
		if c.IsArtnetIn {
			s += " (ArtnetIn) "
		}
		if c.IsArtnetOut {
			s += " (ArtnetOut)"
		}
	} else {
		s = "Disabled"
	}
	return s
}
