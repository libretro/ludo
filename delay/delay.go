package delay

import (
	"log"

	"github.com/libretro/ludo/netplay"
)

var InputQueue chan [10][20]bool

var Count uint64

var playerInput [10][20]bool

func init() {
	InputQueue = make(chan [10][20]bool, 60)
	playerInput = [10][20]bool{}
}

func ReceiveInputs() {
	for {
		log.Println("receive inputs")

		netinput := [20]byte{}
		InputQueue <- playerInput
		if _, err := netplay.Conn.Read(netinput[:]); err != nil {
			log.Fatalln(err)
		}

		Count++
		log.Println("incr", Count)

		log.Println(netinput)

		for i, b := range netinput {
			if b == 1 {
				playerInput[Count%10][i] = true
			}
		}

		InputQueue <- playerInput
	}
}
