package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/user"

	"github.com/kivutar/go-playthemall/libretro"
)

var options struct {
	Updated bool
	Vars    Options
}

// Options wraps a []libretro.Variable to provide marshaling and unmarshaling
type Options []libretro.Variable

// MarshalJSON does JSON marshaling for a []libretro.Variable
func (w Options) MarshalJSON() ([]byte, error) {
	m := make(map[string]string)
	for _, v := range w {
		m[v.Key()] = v.Choices()[0]
	}
	return json.Marshal(m)
}

func saveOptions() error {
	lock.Lock()
	defer lock.Unlock()

	usr, _ := user.Current()

	b, _ := json.MarshalIndent(Options(options.Vars), "", "  ")
	name := filename(g.corePath)
	f, err := os.Create(usr.HomeDir + "/.playthemall/" + name + ".json")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, bytes.NewReader(b))
	return err
}
