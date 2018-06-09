package main

import (
	"fmt"

	"github.com/fatih/structs"
)

type screenSettings struct {
	entry
}

func buildSettings() scene {
	var list screenSettings
	list.label = "Settings"

	fields := structs.Fields(&settings)
	for _, f := range fields {
		f := f
		list.children = append(list.children, entry{
			label: f.Tag("label"),
			icon:  "subsetting",
			callbackIncr: func(direction int) {
				incrCallbacks[f.Name()](f, direction)
			},
			callbackValue: func() string {
				return fmt.Sprintf(f.Tag("fmt"), f.Value())
			},
		})
	}

	list.segueMount()

	return &list
}

func (s *screenSettings) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *screenSettings) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *screenSettings) segueBack() {
	genericAnimate(&s.entry)
}

func (s *screenSettings) update() {
	genericInput(&s.entry)
}

func (s *screenSettings) render() {
	genericRender(&s.entry)
}
