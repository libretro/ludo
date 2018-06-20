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

// Widgets to display settings values
var widgets = map[string]func(*entry){

	// On/Off switch for boolean settings
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

	// Range widget for audio volume and similat float settings
	"range": func(e *entry) {
		fbw, fbh := window.GetFramebufferSize()
		x := float32(fbw) - 650*menu.ratio - 128*menu.ratio
		y := float32(fbh)*e.yp - 3*menu.ratio
		w := 128 * menu.ratio
		h := 6 * menu.ratio
		x1, y1, x2, y2, x3, y3, x4, y4 := xywhTo4points(x, y, w, h, float32(fbh))
		drawQuad(x1, y1, x2, y2, x3, y3, x4, y4, color{1, 1, 1, e.iconAlpha / 4})

		w = 128 * menu.ratio * e.value().(float32)
		x1, y1, x2, y2, x3, y3, x4, y4 = xywhTo4points(x, y, w, h, float32(fbh))
		drawQuad(x1, y1, x2, y2, x3, y3, x4, y4, color{1, 1, 1, e.iconAlpha})
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
