package main

import (
	"fmt"

	"github.com/fatih/structs"
)

func buildSettings() entry {
	var list entry
	list.label = "Settings"
	list.input = verticalInput
	list.render = verticalRender

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

	initEntries(list)

	return list
}
