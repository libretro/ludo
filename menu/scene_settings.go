package menu

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/structs"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/libretro/ludo/audio"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
	"github.com/libretro/ludo/video"
)

type sceneSettings struct {
	entry
}

func buildSettings() Scene {
	var list sceneSettings
	list.label = "Settings"

	fields := structs.Fields(&settings.Current)
	for _, f := range fields {
		f := f
		// Don't expose settings without label
		if f.Tag("label") == "" {
			continue
		}

		if f.Tag("widget") == "dir" {
			// Directory settings
			list.children = append(list.children, entry{
				label: f.Tag("label"),
				icon:  "folder",
				value: f.Value,
				stringValue: func() string {
					return "[" + utils.FileName(f.Value().(string)) + "]"
				},
				widget: widgets[f.Tag("widget")],
				callbackOK: func() {
					list.segueNext()
					menu.stack = append(menu.stack, buildExplorer(
						f.Value().(string),
						nil,
						func(path string) { dirExplorerCb(path, f) },
						&entry{
							label: "<Select this directory>",
							icon:  "scan",
						}),
					)
				},
			})
		} else {
			// Regular settings
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
	}

	list.segueMount()

	return &list
}

// triggered when selecting a directory in the settings file explorer
func dirExplorerCb(path string, f *structs.Field) {
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		ntf.DisplayAndLog(ntf.Error, "Settings", err.Error())
		return
	}
	f.Set(path)
	ntf.DisplayAndLog(ntf.Success, "Settings", "%s set to %s", f.Tag("label"), f.Value().(string))
	err = settings.Save()
	if err != nil {
		ntf.DisplayAndLog(ntf.Error, "Settings", err.Error())
		return
	}
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
		vid.Reconfigure(settings.Current.VideoFullscreen)
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
		vid.Reconfigure(settings.Current.VideoFullscreen)
		menu.ContextReset()
		settings.Save()
	},
	"VideoFilter": func(f *structs.Field, direction int) {
		filters := []string{"nearest", "linear", "sharp-bilinear"}
		v := f.Value().(string)
		i := utils.IndexOfString(v, filters)
		i += direction
		if i < 0 {
			i = len(filters) - 1
		}
		if i > len(filters)-1 {
			i = 0
		}
		f.Set(filters[i])
		vid.UpdateFilter(filters[i])
		settings.Save()
	},
	"AudioVolume": func(f *structs.Field, direction int) {
		v := f.Value().(float32)
		v += 0.1 * float32(direction)
		if v < 0 {
			v = 0
		}
		if v > 1 {
			v = 1
		}
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

func (s *sceneSettings) Entry() *entry {
	return &s.entry
}

func (s *sceneSettings) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneSettings) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneSettings) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneSettings) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *sceneSettings) render() {
	genericRender(&s.entry)
}

func (s *sceneSettings) drawHintBar() {
	w, h := vid.Window.GetFramebufferSize()
	menu.ratio = float32(w) / 1920
	vid.DrawRect(0.0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 1.0, video.Color{R: 0.75, G: 0.75, B: 0.75, A: 1})

	var stack float32
	if state.Global.CoreRunning {
		stackHint(&stack, "key-p", "RESUME", h)
	}
	stackHint(&stack, "key-up-down", "NAVIGATE", h)
	stackHint(&stack, "key-z", "BACK", h)

	list := menu.stack[len(menu.stack)-1].Entry()
	if list.children[list.ptr].callbackOK != nil {
		stackHint(&stack, "key-x", "SET", h)
	} else {
		stackHint(&stack, "key-left-right", "SET", h)
	}
}
