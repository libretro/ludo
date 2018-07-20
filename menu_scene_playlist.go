package main

import (
	"bufio"
	"os"
)

type screenPlaylist struct {
	entry
}

func buildPlaylist(path string) scene {
	var list screenPlaylist
	list.label = filename(path)

	file, err := os.Open(path)
	if err != nil {
		notifyAndLog("Menu", err.Error())
		list.segueMount()
		return &list
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for {
		more := scanner.Scan()
		if !more {
			break
		}
		scanner.Scan() // path
		name := scanner.Text()
		scanner.Scan() // unused
		scanner.Scan() // unused
		scanner.Scan() // CRC
		scanner.Scan() // lpl

		list.children = append(list.children, entry{
			label: name,
			icon:  "subsetting",
		})
	}

	list.segueMount()

	return &list
}

// Generic stuff

func (s *screenPlaylist) Entry() *entry {
	return &s.entry
}

func (s *screenPlaylist) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *screenPlaylist) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *screenPlaylist) segueBack() {
	genericAnimate(&s.entry)
}

func (s *screenPlaylist) update() {
	genericInput(&s.entry)
}

func (s *screenPlaylist) render() {
	genericRender(&s.entry)
}
