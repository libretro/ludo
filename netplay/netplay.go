package netplay

import (
	"net"

	ntf "github.com/libretro/ludo/notifications"
)

// Listen is used by the netplay host, listening address and port
var Listen string

// Join is used by the netplay guest, address of the host
var Join string

// Conn is the connection between two players
var Conn net.Conn

// Init initialises a netplay session between two players
func Init() {
	if Listen != "" { // Host mode
		ln, err := net.Listen("tcp", Listen)
		if err != nil {
			ntf.DisplayAndLog(ntf.Error, "Netplay", err.Error())
			return
		}

		Conn, err = ln.Accept()
		if err != nil {
			ntf.DisplayAndLog(ntf.Error, "Netplay", err.Error())
			return
		}
		ntf.DisplayAndLog(ntf.Success, "Netplay", "Player #2 is connected.")
	} else if Join != "" { // Guest mode
		var err error
		Conn, err = net.Dial("tcp", Join)
		if err != nil {
			ntf.DisplayAndLog(ntf.Error, "Netplay", err.Error())
			return
		}
		ntf.DisplayAndLog(ntf.Success, "Netplay", "Connected.")
	}
}
