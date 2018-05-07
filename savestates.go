package main

import (
	"io/ioutil"
)

func saveState() error {
	s := g.core.SerializeSize()
	bytes, err := g.core.Serialize(s)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("savestate1", bytes, 0644)
}

func loadState() error {
	s := g.core.SerializeSize()
	bytes, err := ioutil.ReadFile("savestate1")
	if err != nil {
		return err
	}
	err = g.core.Unserialize(bytes, s)
	if err != nil {
		return err
	}
	return nil
}
