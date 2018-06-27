package main

import (
	"encoding/json"
	"os/user"

	"github.com/kivutar/go-playthemall/libretro"
)

var options struct {
	Updated bool
	Vars    []libretro.Variable
	Choices []int
}

// // Options wraps a []libretro.Variable to provide marshaling and unmarshaling
// type Options []libretro.Variable

// // MarshalJSON does JSON marshaling for a []libretro.Variable
// func (w Options) MarshalJSON() ([]byte, error) {
// 	m := make(map[string]string)
// 	for _, v := range w {
// 		m[v.Key()] = v.Choices()[0]
// 	}
// 	return json.Marshal(m)
// }

// func saveOptions() error {
// 	lock.Lock()
// 	defer lock.Unlock()

// 	usr, _ := user.Current()

// 	b, _ := json.MarshalIndent(Options(options.Vars), "", "  ")

// 	name := filename(g.corePath)
// 	f, err := os.Create(usr.HomeDir + "/.playthemall/" + name + ".json")
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	_, err = io.Copy(f, bytes.NewReader(b))
// 	return err
// }

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
