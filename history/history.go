package history

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/laverya/yaml.v3"
)

// Game represents a game in the history file
type Game struct {
	Path   string // Absolute path of the game on the filesystem
	Name   string // Human readable name of the game, comes from the RDB
	System string
	Core   string
}

// History is a list of games
type History []Game

// List is the list of recently played games
var List History

// Push pushes a game onto the stack
func Push(g Game) {
	List = append([]Game{g}, List...)

	// Deduplicate
	l := History{}
	exist := map[string]bool{}
	for _, g := range List {
		if !exist[g.Path] {
			l = append(l, g)
			exist[g.Path] = true
		}
	}
	List = l

	err := Save()
	if err != nil {
		log.Println(err)
	}
}

// Load loads history.yml in memory
func Load() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(filepath.Join(homeDir, ".ludo", "history.yml"))
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, &List)
}

// Save persists the history as a yml file
func Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	file, err := os.Create(homeDir + "/.ludo/history.yml")
	if err != nil {
		return err
	}
	defer file.Close()

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	defer enc.Close()
	enc.SetLineLength(-1)
	if err := enc.Encode(List); err != nil {
		return err
	}

	_, err = buf.WriteTo(file)
	return err
}
