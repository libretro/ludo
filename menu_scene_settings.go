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
			incr: func(direction int) {
				incrCallbacks[f.Name()](f, direction)
			},
			value: f.Value,
			stringValue: func() string {
				return fmt.Sprintf(f.Tag("fmt"), f.Value())
			},
			widget: widgets[f.Tag("widget")],
		})
	}

	list.segueMount()

	return &list
}

var widgets = map[string]func(*entry){
	"switch": func(e *entry) {
		icon := "on"
		if e.value().(bool) {
			icon = "off"
		}
		w, h := window.GetFramebufferSize()
		drawImage(menu.icons[icon],
			float32(w)-650*menu.ratio-128*menu.ratio,
			float32(h)*e.yp-64*1.25*menu.ratio,
			128*menu.ratio, 128*menu.ratio,
			1.25, color{1, 1, 1, e.iconAlpha})
	},
}

// Generic stuff

func (s *screenSettings) Entry() *entry {
	return &s.entry
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
