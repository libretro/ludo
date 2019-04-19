package ludos

import (
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CurrentNetwork is the network we're connected to
var CurrentNetwork Network
var counter int

// Network is a network as detected by connman
type Network struct {
	SSID string
	Path string
}

var cache map[string]string

// ScanNetworks enables connman and returns the list of available SSIDs
func ScanNetworks() []Network {
	cache = map[string]string{}

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
			Path: line[25:len(line)],
		}
		networks = append(networks, network)
	}

	return networks
}

// NetworkStatus returns the status of a network
func NetworkStatus(network Network) string {
	_, ok := cache[network.Path]
	if !ok && counter%120 == 0 {
		out, _ := exec.Command(
			"/usr/bin/bash",
			"-c",
			"connmanctl services "+network.Path+" | grep State",
		).Output()
		if strings.Contains(string(out), "online") {
			cache[network.Path] = "Online"
			CurrentNetwork = network
		} else {
			cache[network.Path] = ""
		}
	}
	counter++
	return cache[network.Path]
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
IPv4.method=dhcp`, network.Path, network.SSID, hexSSID, passphrase)

	err := os.MkdirAll("/var/lib/connman/"+network.Path, os.ModePerm)
	if err != nil {
		return err
	}

	fd, err := os.Create("/var/lib/connman/" + network.Path + "/service")
	if err != nil {
		return err
	}
	defer fd.Close()
	fd.WriteString(config)

	err = exec.Command("/usr/bin/connmanctl", "connect", network.Path).Run()
	if err != nil {
		return err
	}

	cache = map[string]string{}

	return nil
}
