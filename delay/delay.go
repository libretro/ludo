package delay

import (
	"log"

	"github.com/libretro/ludo/netplay"
)

var LocalQueue chan [20]bool
var RemoteQueue chan [20]bool

var Count uint64

func init() {
	LocalQueue = make(chan [20]bool, 60)
	RemoteQueue = make(chan [20]bool, 60)
}

func ReceiveInputs() {
	for {
		netinput := [20]byte{}
		if _, err := netplay.Conn.Read(netinput[:]); err != nil {
			log.Fatalln(err)
		}

		Count++

		playerInput := [20]bool{}
		for i, b := range netinput {
			if b == 1 {
				playerInput[i] = true
			}
		}

		RemoteQueue <- playerInput
	}
}
