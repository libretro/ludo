package menu

import (
	"fmt"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/libretro/go-playthemall/audio"
	"github.com/libretro/go-playthemall/settings"
	"github.com/libretro/go-playthemall/video"

	"github.com/fatih/structs"
)

type screenSettings struct {
	entry
}

func buildSettings() Scene {
	var list screenSettings
	list.label = "Settings"

	fields := structs.Fields(&settings.Settings)
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
		icon := "off"
		if e.value().(bool) {
			icon = "on"
		}
		w, h := vid.Window.GetFramebufferSize()
		vid.DrawImage(menu.icons[icon],
			float32(w)-650*menu.ratio-128*menu.ratio,
			float32(h)*e.yp-64*1.25*menu.ratio,
			128*menu.ratio, 128*menu.ratio,
			1.25, video.Color{R: 1, G: 1, B: 1, A: e.iconAlpha})
	},

	// Range widget for audio volume and similat float settings
	"range": func(e *entry) {
		fbw, fbh := vid.Window.GetFramebufferSize()
		x := float32(fbw) - 650*menu.ratio - 256*e.scale*menu.ratio
		y := float32(fbh)*e.yp - 3*menu.ratio
		w := 256 * e.scale * menu.ratio
		h := 6 * menu.ratio
		x1, y1, x2, y2, x3, y3, x4, y4 := video.XYWHTo4points(x, y, w, h, float32(fbh))
		vid.DrawQuad(x1, y1, x2, y2, x3, y3, x4, y4, video.Color{R: 1, G: 1, B: 1, A: e.iconAlpha / 4})

		w = 256 * e.scale * menu.ratio * e.value().(float32)
		x1, y1, x2, y2, x3, y3, x4, y4 = video.XYWHTo4points(x, y, w, h, float32(fbh))
		vid.DrawQuad(x1, y1, x2, y2, x3, y3, x4, y4, video.Color{R: 1, G: 1, B: 1, A: e.iconAlpha})
	},
}

type callbackIncrement func(*structs.Field, int)

// incrCallbacks is a map of callbacks called when a setting value is changed.
var incrCallbacks = map[string]callbackIncrement{
	"VideoFullscreen": func(f *structs.Field, direction int) {
		v := f.Value().(bool)
		v = !v
		f.Set(v)
		vid = video.Init(settings.Settings.VideoFullscreen)
		settings.Save()
	},
	"VideoMonitorIndex": func(f *structs.Field, direction int) {
		v := f.Value().(int)
		v += direction
		if v < 0 {
			v = 0
		}
		if v > len(glfw.GetMonitors())-1 {
			v = len(glfw.GetMonitors()) - 1
		}
		f.Set(v)
		vid = video.Init(settings.Settings.VideoFullscreen)
		settings.Save()
	},
	"AudioVolume": func(f *structs.Field, direction int) {
		v := f.Value().(float32)
		v += 0.1 * float32(direction)
		f.Set(v)
		audio.SetVolume(v)
		settings.Save()
	},
	"ShowHiddenFiles": func(f *structs.Field, direction int) {
		v := f.Value().(bool)
		v = !v
		f.Set(v)
		settings.Save()
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
