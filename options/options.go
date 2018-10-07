// Package options deals with configuration at the libretro core level. Each
// core exports a list of variables that can take different values. This package
// can list them, load them, and save them.
package options

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/user"
	"sync"

	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
)

var lock sync.Mutex

// Options is a container type for core options internals
type Options struct {
	Vars    []libretro.Variable
	Choices []int
	Updated bool
}

// New instanciate a core options manager
func New(vars []libretro.Variable) *Options {
	o := &Options{}
	o.Vars = vars
	o.Choices = make([]int, len(o.Vars))
	o.Updated = true
	o.load()
	return o
}

// NumChoices returns the number of choices for a given variable
func (o *Options) NumChoices(choiceIndex int) int {
	return len(o.Vars[choiceIndex].Choices())
}

// Save core options to a file
func (o *Options) Save() error {
	lock.Lock()
	defer lock.Unlock()

	usr, _ := user.Current()

	m := make(map[string]string)
	for i, v := range o.Vars {
		m[v.Key()] = v.Choices()[o.Choices[i]]
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	name := utils.Filename(state.Global.CorePath)
	f, err := os.Create(usr.HomeDir + "/.ludo/" + name + ".json")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, bytes.NewReader(b))
	return err
}

// Load core options from a file
func (o *Options) load() error {
	lock.Lock()
	defer lock.Unlock()

	usr, _ := user.Current()

	name := utils.Filename(state.Global.CorePath)
	b, err := utils.Slurp(usr.HomeDir + "/.ludo/" + name + ".json")
	if err != nil {
		return err
	}

	var opts map[string]string
	err = json.Unmarshal(b, &opts)

	for optk, optv := range opts {
		for i, variable := range o.Vars {
			if variable.Key() == optk {
				for j, c := range variable.Choices() {
					if c == optv {
						o.Choices[i] = j
					}
				}
			}
		}
	}

	return err
}
