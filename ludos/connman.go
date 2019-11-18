package ludos

import (
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CurrentNetwork is the network we're connected to
var CurrentNetwork Network
var counter int

const connmanPath = "/storage/.cache/connman/"

// Network is a network as detected by connman
type Network struct {
	SSID string
	Path string
}

var cache map[string]string

// ScanNetworks enables connman and returns the list of available SSIDs
func ScanNetworks() ([]Network, error) {
	cache = map[string]string{} // clear the status cache
	networks := []Network{}

	err := exec.Command("/usr/bin/connmanctl", "enable", "wifi").Run()
	if err != nil {
		return networks, err
	}
	err = exec.Command("/usr/bin/connmanctl", "scan", "wifi").Run()
	if err != nil {
		return networks, err
	}
	out, err := exec.Command("/usr/bin/connmanctl", "services").Output()
	if err != nil {
		return networks, err
	}

	for _, line := range strings.Split(string(out), "\n") {
		if len(line) == 0 {
			continue
		}
		network := Network{
			SSID: strings.TrimSpace(line[4:24]),
			Path: line[25:],
		}
		networks = append(networks, network)
	}

	return networks, nil
}

// NetworkStatus returns the status of a network
func NetworkStatus(network Network) string {
	if counter++; counter%120 != 0 {
		return cache[network.Path]
	}

	out, err := exec.Command(
		"/usr/bin/bash", "-c",
		"connmanctl services "+network.Path+" | grep State",
	).Output()
	if err != nil {
		return cache[network.Path]
	}

	if strings.Contains(string(out), "online") {
		cache[network.Path] = "Online"
		CurrentNetwork = network
	} else if strings.Contains(string(out), "ready") {
		cache[network.Path] = "Ready"
		CurrentNetwork = network
	} else if strings.Contains(string(out), "association") {
		cache[network.Path] = "Association"
	} else {
		cache[network.Path] = ""
	}

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
IPv4.method=dhcp
`, network.Path, network.SSID, hexSSID, passphrase)

	netPath := filepath.Join(connmanPath, network.Path)

	if err := os.MkdirAll(netPath, os.ModePerm); err != nil {
		return err
	}

	fd, err := os.Create(filepath.Join(netPath, "settings"))
	if err != nil {
		return err
	}

	if _, err := fd.WriteString(config); err != nil {
		return err
	}

	// We want the sync to happen before connmanctl connect, don't defer
	if err := fd.Sync(); err != nil {
		return err
	}

	if err := fd.Close(); err != nil {
		return err
	}

	if err := exec.Command("/usr/bin/connmanctl", "connect", network.Path).
		Run(); err != nil {
		return err
	}

	cache = map[string]string{} // clear the status cache
	cache[network.Path] = "Connecting"

	return nil
}
