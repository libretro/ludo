package main

import (
	"fmt"

	"github.com/fatih/structs"
)

type screenSettings struct {
	entry
}

func buildSettings() screen {
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

	list.present()

	return &list
}

func (s *screenSettings) present() {
	initEntries(&s.entry)
}

func (s *screenSettings) makeRoomForChildren() {
	genericMakeRoomForChildren(&s.entry)
}

func (s *screenSettings) getFocusBack() {
	animateEntries(&s.entry)
}

func (s *screenSettings) update() {
	verticalInput(&s.entry)
}

func (s *screenSettings) render() {
	verticalRender(&s.entry)
}
