// Package patch allows softpatching ROMs based on the presence of a patch file
// next to the ROM. This is useful to apply fan translations without altering
// No-Intro ROMs. Softpatching only works for cores where NeedFullPath is false.
package patch

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Try to apply different patches located next to the game
// Currently only .ups is supported
func Try(gamePath string, bytes []byte) (*[]byte, error) {
	upsFile := strings.TrimSuffix(gamePath, filepath.Ext(gamePath)) + ".ups"
	if _, err := os.Stat(upsFile); !os.IsNotExist(err) {
		pbytes, err := ioutil.ReadFile(upsFile)
		if err != nil {
			return nil, err
		}

		patched, err := applyUPS(pbytes, bytes)
		if err != nil {
			return nil, err
		}

		return patched, nil
	}

	ipsFile := strings.TrimSuffix(gamePath, filepath.Ext(gamePath)) + ".ips"
	if _, err := os.Stat(ipsFile); !os.IsNotExist(err) {
		pbytes, err := ioutil.ReadFile(ipsFile)
		if err != nil {
			return nil, err
		}

		patched, err := applyIPS(pbytes, bytes)
		if err != nil {
			return nil, err
		}

		return patched, nil
	}

	return nil, nil
}
