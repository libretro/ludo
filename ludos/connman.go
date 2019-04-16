package ludos

import (
	"time"

	"github.com/libretro/ludo/state"
)

var currentNetwork string
var connectingTo string

// ScanNetworks enables connman and returns the list of available SSIDs
func ScanNetworks() []string {
	time.Sleep(time.Second * 2)

	return []string{
		"Fake Network 1",
		"Fake Network 2",
		"Fake Network 3",
		"Fake Network 4",
		"Fake Network 5",
		"Fake Network 6",
		"Fake Network 7",
	}
}

// NetworkStatus returns the status of a network
func NetworkStatus(network string) string {
	if network == currentNetwork {
		return "Connected"
	}
	if network == connectingTo {
		return "Connecting"
	}
	return ""
}

// ConnectNetwork attempt to establish a connection to the given network
func ConnectNetwork(network string) {
	connectingTo = network
	time.Sleep(time.Second * 3)
	currentNetwork = network
	state.Global.Network = network
}
