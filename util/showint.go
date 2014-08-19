package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	ints, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Error get interfaces: %s", err.Error())
		os.Exit(1)
	}

	for n, intf := range ints {
		if (intf.Flags & net.FlagUp) != 0 {
			fmt.Printf("Int(%d) : %s", n, intf.Name)

			iface, err := net.InterfaceByName(intf.Name)
			if err != nil {
				fmt.Printf("Error get interfaces: %s", err.Error())
				os.Exit(1)
			}

			fmt.Printf("Interface Name: %s", iface)
		}
	}
}
