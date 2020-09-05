package menu

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/structs"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/ludos"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
)

type sceneSettings struct {
	entry
}

// Don't display settings flagged with hide:"always"
// If we're in Desktop Environment mode, hide settings flagged with hide:"ludos"
// If we're in program mode, hide settings flagged with hide:"program"
func isHidden(f *structs.Field) bool {
	return f.Tag("hide") == "always" ||
		(state.Global.LudOS && f.Tag("hide") == "ludos") ||
		(!state.Global.LudOS && f.Tag("hide") == "app")
}

func buildSettings() Scene {
	var list sceneSettings
	list.label = "Settings"

	if state.Global.LudOS {
		list.children = append(list.children, entry{
			label:       "Wi-Fi",
			icon:        "subsetting",
			stringValue: func() string { return ludos.CurrentNetwork.SSID },
			callbackOK: func() {
				list.segueNext()
				menu.Push(buildWiFi())
			},
		})
	}

	fields := structs.Fields(&settings.Current)
	for _, f := range fields {
		f := f

		if isHidden(f) {
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
					menu.Push(buildExplorer(
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
	info, err := os.Stat(path)
	if err != nil {
		ntf.DisplayAndLog(ntf.Error, "Settings", err.Error())
		return
	}
	if !info.IsDir() {
		ntf.DisplayAndLog(ntf.Error, "Settings", "Not a directory")
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
		vid.DrawImage(menu.icons[icon],
			float32(w)-128*menu.ratio-128*menu.ratio,
			float32(h)*e.yp-64*1.25*menu.ratio,
			128*menu.ratio, 128*menu.ratio,
			1.25, textColor.Alpha(e.iconAlpha))
	},

	// Range widget for audio volume and similat float settings
	"range": func(e *entry) {
		fbw, fbh := vid.Window.GetFramebufferSize()
		x := float32(fbw) - 128*menu.ratio - 175*menu.ratio
		y := float32(fbh)*e.yp - 4*menu.ratio
		w := 175 * menu.ratio
		h := 8 * menu.ratio
		vid.DrawRect(x, y, w, h, 0.9, textColor.Alpha(e.iconAlpha/4))
		w = 175 * menu.ratio * e.value().(float32)
		vid.DrawRect(x, y, w, h, 0.9, textColor.Alpha(e.iconAlpha))
		vid.DrawCircle(x+w, y+4*menu.ratio, 38*menu.ratio, textColor.Alpha(e.iconAlpha))
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
		filters := []string{"Raw", "Smooth", "Pixel Perfect", "CRT"}
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
	"VideoDarkMode": func(f *structs.Field, direction int) {
		v := f.Value().(bool)
		v = !v
		f.Set(v)
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
	"MenuAudioVolume": func(f *structs.Field, direction int) {
		v := f.Value().(float32)
		v += 0.1 * float32(direction)
		if v < 0 {
			v = 0
		}
		if v > 1 {
			v = 1
		}
		f.Set(v)
		audio.SetEffectsVolume(v)
		settings.Save()
	},
	"ShowHiddenFiles": func(f *structs.Field, direction int) {
		v := f.Value().(bool)
		v = !v
		f.Set(v)
		settings.Save()
	},
	"SSHService":       ludos.ServiceSettingIncrCallback,
	"SambaService":     ludos.ServiceSettingIncrCallback,
	"BluetoothService": ludos.ServiceSettingIncrCallback,
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
	vid.DrawRect(0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 0, lightGrey)

	_, upDown, leftRight, a, b, _, _, _, _, guide := hintIcons()

	var stack float32
	list := menu.stack[len(menu.stack)-1].Entry()
	if state.Global.CoreRunning {
		stackHint(&stack, guide, "RESUME", h)
	}
	stackHint(&stack, upDown, "NAVIGATE", h)
	stackHint(&stack, b, "BACK", h)
	if list.children[list.ptr].callbackOK != nil {
		stackHint(&stack, a, "SET", h)
	} else {
		stackHint(&stack, leftRight, "SET", h)
	}
}
