package netplay

import (
	"log"
	"net"

	ntf "github.com/libretro/ludo/notifications"
)

// Listen is used by the netplay host, listening address and port
var Listen int

// Join is used by the netplay guest, address of the host
var Join string

// Conn is the connection between two players
var Conn net.Conn

// Init initialises a netplay session between two players
func Init() {
	if Listen > 0 { // Host mode
		Conn, err := net.ListenUDP("udp", &net.UDPAddr{
			Port: Listen,
		})
		if err != nil {
			log.Println("Netplay", err.Error())
			return
		}

		Conn.SetReadBuffer(1048576)

		msg := [2]byte{}
		Conn.Read(msg[:])
		log.Println(msg)
		log.Println(ntf.Success, "Netplay", "Player #2 is connected.")
	} else if Join != "" { // Guest mode
		var err error
		Conn, err = net.Dial("udp", Join)
		if err != nil {
			log.Println("Netplay", err.Error())
			return
		}
		log.Println(ntf.Success, "Netplay", "Connected.")
		Conn.Write([]byte("hi"))
	}
}
