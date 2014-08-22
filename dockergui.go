package main

import (
	"code.google.com/p/goncurses"
	//	"fmt"
	"log"
)

const (
	WEIGHT = 10
	WIDTH  = 30
)

func main() {
	stdscr, err := goncurses.Init()
	if err != nil {
		log.Fatal("init:", err)
	}
	defer goncurses.End()

	goncurses.Raw(true)
	goncurses.Echo(false)
	goncurses.Cursor(0)
	stdscr.Clear()
	stdscr.Keypad(true)

	menu_items := []string{
		"Add controller [ ]",
		"Add Switcher [ ]",
	}

	items := make([]*goncurses.MenuItem, len(menu_items))
	for i, val := range menu_items {
		items[i], _ = goncurses.NewItem(val, "")
		defer items[i].Free()
	}

	menu, err := goncurses.NewMenu(items)
	if err != nil {
		stdscr.Print(err)
		return
	}

	defer menu.Free()

	menu.Option(goncurses.O_ONEVALUE, false)

	y, _ := stdscr.MaxYX()
	stdscr.MovePrint(y-3, 0, "UpDown to move, spacebar to toggle, enter: log. q to exit")

	stdscr.Refresh()
	menu.Post()
	defer menu.UnPost()

	for {
		goncurses.Update()
		ch := stdscr.GetChar()

		switch ch {
		case 'q':
			return
		case ' ':
			menu.Driver(goncurses.REQ_TOGGLE)
		case goncurses.KEY_RETURN, goncurses.KEY_ENTER:
			var list string
			for _, item := range menu.Items() {
				if item.Value() {
					list += "\"" + item.Name() + "\""
				}
			}
			stdscr.Move(20, 0)
			stdscr.ClearToEOL()
			stdscr.MovePrint(20, 0, list)
			stdscr.Refresh()
		default:
			menu.Driver(goncurses.DriverActions[ch])
		}
	}
}
