package ludos

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CurrentNetwork is the name of the Wi-Fi network we're connected to
var CurrentNetwork string

// ConnectingTo is the name of the Wi-Fi network we're connecting to
var ConnectingTo string

// Network is a network as detected by connman
type Network struct {
	SSID string
	ID   string
}

// ScanNetworks enables connman and returns the list of available SSIDs
func ScanNetworks() []Network {
	exec.Command("/usr/bin/connmanctl", "enable", "wifi").Run()
	exec.Command("/usr/bin/connmanctl", "scan", "wifi").Run()

	var stdout bytes.Buffer
	cmd := exec.Command("/usr/bin/connmanctl", "services")
	cmd.Stdout = &stdout
	cmd.Run()

	networks := []Network{}
	for _, line := range strings.Split(string(stdout.Bytes()), "\n") {
		if len(line) == 0 {
			continue
		}
		network := Network{
			SSID: strings.TrimSpace(line[4:24]),
			ID:   line[25:len(line)],
		}
		networks = append(networks, network)
	}

	return networks
}

// NetworkStatus returns the status of a network
func NetworkStatus(network string) string {
	if network == CurrentNetwork {
		return "Connected"
	}
	if network == ConnectingTo {
		return "Connecting"
	}
	return ""
}

// ConnectNetwork attempt to establish a connection to the given network
func ConnectNetwork(network Network, pass string) {
	var hexSSID []byte
	hex.Encode(hexSSID, []byte(network.SSID))

	config := fmt.Sprintf(`[%s]
Name=%s
SSID=%s
Favorite=true
AutoConnect=true
Passphrase=%s
IPv4.method=dhcp`, network.ID, network.SSID, hexSSID, pass)

	err := os.MkdirAll("/var/lib/connman/"+network.ID, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	fd, _ := os.Create("/var/lib/connman/" + network.ID + "/service")
	defer fd.Close()
	fd.Write([]byte(config))

	exec.Command("/usr/bin/connmanctl", "connect", network.ID).Run()
}
