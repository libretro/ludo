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
	Vars    []libretro.Variable
	Choices []int
}

func saveOptions() error {
	lock.Lock()
	defer lock.Unlock()

	usr, _ := user.Current()

	m := make(map[string]string)
	for i, v := range options.Vars {
		m[v.Key()] = v.Choices()[options.Choices[i]]
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	name := filename(g.corePath)
	f, err := os.Create(usr.HomeDir + "/.playthemall/" + name + ".json")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, bytes.NewReader(b))
	return err
}

func loadOptions() error {
	lock.Lock()
	defer lock.Unlock()

	usr, _ := user.Current()

	name := filename(g.corePath)
	b, err := slurp(usr.HomeDir + "/.playthemall/" + name + ".json")
	if err != nil {
		return err
	}

	var opts map[string]string
	err = json.Unmarshal(b, &opts)

	for optk, optv := range opts {
		for i, variable := range options.Vars {
			if variable.Key() == optk {
				for j, c := range variable.Choices() {
					if c == optv {
						options.Choices[i] = j
					}
				}
			}
		}
	}

	return err
}
