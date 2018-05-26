package main

import (
	"io/ioutil"
	"os/user"
	"path/filepath"
)

func savestateName() string {
	name := filepath.Base(g.gamePath)
	ext := filepath.Ext(name)
	name = name[0 : len(name)-len(ext)]
	return name + ".state"
}

func saveState() error {
	usr, _ := user.Current()
	s := g.core.SerializeSize()
	bytes, err := g.core.Serialize(s)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(usr.HomeDir+"/.playthemall/savestates/"+savestateName(), bytes, 0644)
}

func loadState() error {
	usr, _ := user.Current()
	s := g.core.SerializeSize()
	bytes, err := ioutil.ReadFile(usr.HomeDir + "/.playthemall/savestates/" + savestateName())
	if err != nil {
		return err
	}
	err = g.core.Unserialize(bytes, s)
	return err
}
