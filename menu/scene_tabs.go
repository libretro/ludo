package menu

import (
	"fmt"
	"os/user"
	"sort"

	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/playlists"
	"github.com/libretro/ludo/scanner"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
	"github.com/libretro/ludo/video"
	colorful "github.com/lucasb-eyer/go-colorful"

	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type sceneTabs struct {
	entry
}

func buildTabs() Scene {
	var list sceneTabs
	list.label = "Ludo"

	list.children = append(list.children, entry{
		label:    "Main Menu",
		subLabel: "Load cores and games manually",
		icon:     "main",
		callbackOK: func() {
			menu.Push(buildMainMenu())
		},
	})

	list.children = append(list.children, entry{
		label:    "Settings",
		subLabel: "Configure Ludo",
		icon:     "setting",
		callbackOK: func() {
			menu.Push(buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label:    "History",
		subLabel: "Play again",
		icon:     "history",
		callbackOK: func() {
			menu.Push(buildHistory())
		},
	})

	list.children = append(list.children, getPlaylists()...)

	list.children = append(list.children, entry{
		label:    "Add games",
		subLabel: "Scan your collection",
		icon:     "add",
		callbackOK: func() {
			usr, _ := user.Current()
			menu.Push(buildExplorer(usr.HomeDir, nil,
				func(path string) {
					scanner.ScanDir(path, refreshTabs)
				},
				&entry{
					label: "<Scan this directory>",
					icon:  "scan",
				}))
		},
	})

	list.segueMount()

	return &list
}

// refreshTabs is called after playlist scanning is complete. It inserts the new
// playlists in the tabs, and makes sure that all the icons are positioned and
// sized properly.
func refreshTabs() {
	e := menu.stack[0].Entry()
	l := len(e.children)
	pls := getPlaylists()

	// This assumes that the 3 first tabs are not playlists, and that the last
	// tab is the scanner.
	e.children = append(e.children[:3], append(pls, e.children[l-1:]...)...)

	// Update which tab is the active tab after the refresh
	if e.ptr >= 3 {
		e.ptr += len(pls) - (l - 4)
	}

	// Ensure new icons are styled properly
	for i := range e.children {
		if i == e.ptr {
			e.children[i].iconAlpha = 1
			e.children[i].scale = 0.75
			e.children[i].width = 500
		} else if i < e.ptr {
			e.children[i].iconAlpha = 1
			e.children[i].scale = 0.25
			e.children[i].width = 128
		} else if i > e.ptr {
			e.children[i].iconAlpha = 1
			e.children[i].scale = 0.25
			e.children[i].width = 128
		}
	}

	// Adapt the tabs scroll value
	if len(menu.stack) == 1 {
		menu.scroll = float32(e.ptr * 128)
	} else {
		e.children[e.ptr].margin = 1360
		menu.scroll = float32(e.ptr*128 + 680)
	}
}

// getPlaylists browse the filesystem for CSV files, parse them and returns
// a list of menu entries. It is used in the tabs, but could be used somewhere
// else too.
func getPlaylists() []entry {
	playlists.Load()

	// To store the keys in slice in sorted order
	var keys []string
	for k := range playlists.Playlists {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pls []entry
	for _, path := range keys {
		path := path
		filename := utils.FileName(path)
		count := playlists.Count(path)
		label := playlists.ShortName(filename)
		pls = append(pls, entry{
			label:    label,
			subLabel: fmt.Sprintf("%d Games - 0 Favorites", count),
			icon:     filename,
			callbackOK: func() {
				menu.Push(buildPlaylist(path))
			},
		})
	}
	return pls
}

func (tabs *sceneTabs) Entry() *entry {
	return &tabs.entry
}

func (tabs *sceneTabs) segueMount() {
	for i := range tabs.children {
		e := &tabs.children[i]

		if i == tabs.ptr {
			e.labelAlpha = 1
			e.iconAlpha = 1
			e.scale = 0.75
			e.width = 500
		} else if i < tabs.ptr {
			e.labelAlpha = 0
			e.iconAlpha = 1
			e.scale = 0.25
			e.width = 128
		} else if i > tabs.ptr {
			e.labelAlpha = 0
			e.iconAlpha = 1
			e.scale = 0.25
			e.width = 128
		}
	}

	tabs.animate()
}

func (tabs *sceneTabs) segueBack() {
	tabs.animate()
}

func (tabs *sceneTabs) animate() {
	for i := range tabs.children {
		e := &tabs.children[i]

		var labelAlpha, scale, width float32
		if i == tabs.ptr {
			labelAlpha = 1
			scale = 0.75
			width = 500
		} else if i < tabs.ptr {
			labelAlpha = 0
			scale = 0.25
			width = 128
		} else if i > tabs.ptr {
			labelAlpha = 0
			scale = 0.25
			width = 128
		}

		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, labelAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, 1, 0.15, ease.OutSine)
		menu.tweens[&e.scale] = gween.New(e.scale, scale, 0.15, ease.OutSine)
		menu.tweens[&e.width] = gween.New(e.width, width, 0.15, ease.OutSine)
		menu.tweens[&e.margin] = gween.New(e.margin, 0, 0.15, ease.OutSine)
	}
	menu.tweens[&menu.scroll] = gween.New(menu.scroll, float32(tabs.ptr*128), 0.15, ease.OutSine)
}

func (tabs *sceneTabs) segueNext() {
	cur := &tabs.children[tabs.ptr]
	menu.tweens[&cur.margin] = gween.New(cur.margin, 1360, 0.15, ease.OutSine)
	menu.tweens[&menu.scroll] = gween.New(menu.scroll, menu.scroll+680, 0.15, ease.OutSine)
	for i := range tabs.children {
		e := &tabs.children[i]
		if i != tabs.ptr {
			menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, 0, 0.15, ease.OutSine)
		}
	}
}

func (tabs *sceneTabs) update(dt float32) {
	// Right
	repeatRight(dt, input.NewState[0][libretro.DeviceIDJoypadRight], func() {
		tabs.ptr++
		if tabs.ptr >= len(tabs.children) {
			tabs.ptr = 0
		}
		audio.PlayEffect(audio.Effects["down"])
		tabs.animate()
	})

	// Left
	repeatLeft(dt, input.NewState[0][libretro.DeviceIDJoypadLeft], func() {
		tabs.ptr--
		if tabs.ptr < 0 {
			tabs.ptr = len(tabs.children) - 1
		}
		audio.PlayEffect(audio.Effects["up"])
		tabs.animate()
	})

	// OK
	if input.Released[0][libretro.DeviceIDJoypadA] {
		if tabs.children[tabs.ptr].callbackOK != nil {
			audio.PlayEffect(audio.Effects["ok"])
			tabs.segueNext()
			tabs.children[tabs.ptr].callbackOK()
		}
	}
}

func (tabs sceneTabs) render() {
	_, h := vid.Window.GetFramebufferSize()

	stackWidth := 710 * menu.ratio
	for i, e := range tabs.children {

		cf := colorful.Hcl(float64(i)*20, 0.5, 0.5)
		c := video.Color{R: float32(cf.R), G: float32(cf.B), B: float32(cf.G), A: e.iconAlpha}

		x := -menu.scroll*menu.ratio + stackWidth + e.width/2*menu.ratio

		stackWidth += e.width*menu.ratio + e.margin*menu.ratio

		if e.labelAlpha > 0 {
			vid.Font.SetColor(c.Alpha(e.labelAlpha))
			lw := vid.Font.Width(0.5*menu.ratio, e.label)
			vid.Font.Printf(x-lw/2, float32(int(float32(h)/2+250*menu.ratio)), 0.5*menu.ratio, e.label)
			lw = vid.Font.Width(0.4*menu.ratio, e.subLabel)
			vid.Font.Printf(x-lw/2, float32(int(float32(h)/2+330*menu.ratio)), 0.4*menu.ratio, e.subLabel)
		}

		vid.DrawImage(menu.icons["hexagon"],
			x-220*e.scale*menu.ratio, float32(h)/2-220*e.scale*menu.ratio,
			440*menu.ratio, 440*menu.ratio, e.scale, c)

		vid.DrawImage(menu.icons[e.icon],
			x-128*e.scale*menu.ratio, float32(h)/2-128*e.scale*menu.ratio,
			256*menu.ratio, 256*menu.ratio, e.scale, white.Alpha(e.iconAlpha))
	}
}

func (tabs sceneTabs) drawHintBar() {
	w, h := vid.Window.GetFramebufferSize()
	vid.DrawRect(0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 0, lightGrey)

	_, _, leftRight, a, _, _, _, _, _, guide := hintIcons()

	var stack float32
	if state.Global.CoreRunning {
		stackHint(&stack, guide, "RESUME", h)
	}
	stackHint(&stack, leftRight, "NAVIGATE", h)
	stackHint(&stack, a, "OPEN", h)
}
