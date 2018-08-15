package main

import (
	"io/ioutil"
	"os/user"
	"path/filepath"

	"github.com/libretro/go-playthemall/state"
)

func savestateName() string {
	name := filepath.Base(state.Global.GamePath)
	ext := filepath.Ext(name)
	name = name[0 : len(name)-len(ext)]
	return name + ".state"
}

func saveState() error {
	usr, _ := user.Current()
	s := state.Global.Core.SerializeSize()
	bytes, err := state.Global.Core.Serialize(s)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(usr.HomeDir+"/.playthemall/savestates/"+savestateName(), bytes, 0644)
}

func loadState() error {
	usr, _ := user.Current()
	s := state.Global.Core.SerializeSize()
	bytes, err := ioutil.ReadFile(usr.HomeDir + "/.playthemall/savestates/" + savestateName())
	if err != nil {
		return err
	}
	err = state.Global.Core.Unserialize(bytes, s)
	return err
}
