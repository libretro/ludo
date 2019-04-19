package ludos

import (
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
	out, _ := exec.Command("/usr/bin/connmanctl", "services").Output()

	networks := []Network{}
	for _, line := range strings.Split(string(out), "\n") {
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
func ConnectNetwork(network Network, passphrase string) error {
	hexSSID := hex.EncodeToString([]byte(network.SSID))

	config := fmt.Sprintf(`[%s]
Name=%s
SSID=%s
Favorite=true
AutoConnect=true
Passphrase=%s
IPv4.method=dhcp`, network.ID, network.SSID, hexSSID, passphrase)

	err := os.MkdirAll("/var/lib/connman/"+network.ID, os.ModePerm)
	if err != nil {
		return err
	}

	fd, err := os.Create("/var/lib/connman/" + network.ID + "/service")
	if err != nil {
		return err
	}
	defer fd.Close()
	fd.WriteString(config)

	err = exec.Command("/usr/bin/connmanctl", "connect", network.ID).Run()
	if err != nil {
		return err
	}

	return nil
}
