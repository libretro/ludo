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
	patchFile := strings.TrimSuffix(gamePath, filepath.Ext(gamePath)) + ".ups"
	if _, err := os.Stat(patchFile); !os.IsNotExist(err) {
		pbytes, err := ioutil.ReadFile(patchFile)
		if err != nil {
			return nil, err
		}

		patched, err := applyUPS(pbytes, bytes)
		if err != nil {
			return nil, err
		}

		return patched, nil
	}
	return nil, nil
}
