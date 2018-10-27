package menu

import (
	"fmt"

	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"

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
		// Don't expose settings without label
		if f.Tag("label") == "" {
			continue
		}
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
		color := video.Color{R: 0, G: 0, B: 0, A: e.iconAlpha}
		if state.Global.CoreRunning {
			color = video.Color{R: 1, G: 1, B: 1, A: e.iconAlpha}
		}
		vid.DrawImage(menu.icons[icon],
			float32(w)-128*menu.ratio-128*menu.ratio,
			float32(h)*e.yp-64*1.25*menu.ratio,
			128*menu.ratio, 128*menu.ratio,
			1.25, color)
	},

	// Range widget for audio volume and similat float settings
	"range": func(e *entry) {
		fbw, fbh := vid.Window.GetFramebufferSize()
		x := float32(fbw) - 128*menu.ratio - 175*menu.ratio
		y := float32(fbh)*e.yp - 4*menu.ratio
		w := 175 * menu.ratio
		h := 8 * menu.ratio
		c := video.Color{R: 0, G: 0, B: 0, A: e.iconAlpha}
		if state.Global.CoreRunning {
			c = video.Color{R: 1, G: 1, B: 1, A: e.iconAlpha}
		}
		vid.DrawRoundedRect(x, y, w, h, 0.9, video.Color{R: c.R, G: c.G, B: c.B, A: e.iconAlpha / 4})
		w = 175 * menu.ratio * e.value().(float32)
		vid.DrawRoundedRect(x, y, w, h, 0.9, c)
		vid.DrawCircle(x+w, y+4*menu.ratio, 38*menu.ratio, c)
	},
}

type callbackIncrement func(*structs.Field, int)

// incrCallbacks is a map of callbacks called when a setting value is changed.
var incrCallbacks = map[string]callbackIncrement{
	"VideoFullscreen": func(f *structs.Field, direction int) {
		v := f.Value().(bool)
		v = !v
		f.Set(v)
		vid.Reconfigure(settings.Settings.VideoFullscreen)
		menu.ContextReset()
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
		vid.Reconfigure(settings.Settings.VideoFullscreen)
		menu.ContextReset()
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
