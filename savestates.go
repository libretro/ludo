package main

import (
	"io/ioutil"
	"path/filepath"
)

func savestateName() string {
	name := filepath.Base(g.gamePath)
	ext := filepath.Ext(name)
	name = name[0 : len(name)-len(ext)]
	return name + ".state"
}

func saveState() error {
	s := g.core.SerializeSize()
	bytes, err := g.core.Serialize(s)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(savestateName(), bytes, 0644)
}

func loadState() error {
	s := g.core.SerializeSize()
	bytes, err := ioutil.ReadFile(savestateName())
	if err != nil {
		return err
	}
	err = g.core.Unserialize(bytes, s)
	if err != nil {
		return err
	}
	return nil
}
