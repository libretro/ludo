// Package options deals with configuration at the libretro core level. Each
// core exports a list of variables that can take different values. This package
// can list them, load them, and save them.
package options

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
	"github.com/pelletier/go-toml"
)

// Variable represents one core option. A variable can take a limited number of
// values. The possibilities are stored in v.Choices. The current value
// can be accessed with v.Choices[v.Choice]
type Variable struct {
	Key     string   // unique id of the variable
	Desc    string   // human readable name of the variable
	Choices []string // available values
	Choice  int      // index of the current value
	Default string
}

// Options is a container type for core options internals
type Options struct {
	Vars    []*Variable // the variables exposed by the core
	Updated bool        // notify the core that values have been updated

	sync.Mutex
}

// VariableInterface is used as a compatibility layer for old v0 options and v1 options
type VariableInterface interface {
	Key() string
	Desc() string
	Choices() []string
	DefaultValue() string
}

// New instantiate a core options manager
func New(vars []VariableInterface) (*Options, error) {
	o := &Options{}

	// Cache core options
	for _, v := range vars {
		v := v
		o.Vars = append(o.Vars, &Variable{
			Key:     v.Key(),
			Desc:    v.Desc(),
			Choices: v.Choices(),
			Default: v.DefaultValue(),
			Choice:  utils.IndexOfString(v.DefaultValue(), v.Choices()),
		})
	}
	o.Updated = true
	err := o.load()
	return o, err
}

// Save core options to a file
func (o *Options) Save() error {
	o.Lock()
	defer o.Unlock()

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	m := make(map[string]string)
	for _, v := range o.Vars {
		m[strings.Replace(v.Key, ".", "___", 1)] = v.Choices[v.Choice]
	}
	b, err := toml.Marshal(m)
	if err != nil {
		return err
	}

	name := utils.FileName(state.CorePath)
	fd, err := os.Create(filepath.Join(home, ".ludo", name+".toml"))
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(fd, bytes.NewReader(b))
	if err != nil {
		return err
	}

	return fd.Sync()
}

// Load core options from a file
func (o *Options) load() error {
	o.Lock()
	defer o.Unlock()

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	name := utils.FileName(state.CorePath)
	b, err := ioutil.ReadFile(filepath.Join(home, ".ludo", name+".toml"))
	if err != nil {
		return err
	}

	var opts map[string]string
	err = toml.Unmarshal(b, &opts)
	if err != nil {
		return err
	}

	for optk, optv := range opts {
		for _, variable := range o.Vars {
			if variable.Key == strings.Replace(optk, "___", ".", 1) {
				for j, c := range variable.Choices {
					if c == optv {
						variable.Choice = j
					}
				}
			}
		}
	}

	return nil
}
