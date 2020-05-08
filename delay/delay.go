package delay

import (
	"log"

	"github.com/libretro/ludo/netplay"
)

var InputQueue chan [20]bool

var Count uint64

func init() {
	InputQueue = make(chan [20]bool, 60)
}

func ReceiveInputs() {
	for {
		log.Println("receive inputs")

		netinput := [20]byte{}
		if _, err := netplay.Conn.Read(netinput[:]); err != nil {
			log.Fatalln(err)
		}

		Count++
		log.Println("incr", Count)

		log.Println(netinput)

		playerInput := [20]bool{}
		for i, b := range netinput {
			if b == 1 {
				playerInput[i] = true
			}
		}

		InputQueue <- playerInput
	}
}
