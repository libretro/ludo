package menu

import (
	"fmt"
	"os/user"
	"sort"

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

type sceneTags struct {
	entry
}

func buildTabs() Scene {
	var list sceneTags
	list.label = "Ludo"

	list.children = append(list.children, entry{
		label:    "Main Menu",
		subLabel: "Load cores and games manually",
		icon:     "main",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildMainMenu())
		},
	})

	list.children = append(list.children, entry{
		label:    "Settings",
		subLabel: "Configure Ludo",
		icon:     "setting",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, getPlaylists()...)

	list.children = append(list.children, entry{
		label:    "Add games",
		subLabel: "Scan your collection",
		icon:     "add",
		callbackOK: func() {
			usr, _ := user.Current()
			menu.stack = append(menu.stack, buildExplorer(usr.HomeDir, nil,
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

	// This assumes that the two first tabs are not playlists, and that the last
	// tab is the scanner.
	e.children = append(e.children[:2], append(pls, e.children[l-1:]...)...)

	// Update which tab is the active tab after the refresh
	if e.ptr >= 2 {
		e.ptr += len(pls) - (l - 3)
	}

	// Ensure new icons are styled properly
	for i := range e.children {
		if i == e.ptr {
			e.children[i].yp = 0.5
			e.children[i].iconAlpha = 1
			e.children[i].scale = 0.75
			e.children[i].width = 500
		} else if i < e.ptr {
			e.children[i].yp = 0.5
			e.children[i].iconAlpha = 1
			e.children[i].scale = 0.25
			e.children[i].width = 128
		} else if i > e.ptr {
			e.children[i].yp = 0.5
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
		pls = append(pls, entry{
			label:    playlists.ShortName(filename),
			subLabel: fmt.Sprintf("%d Games - 0 Favorites", count),
			icon:     filename,
			callbackOK: func() {
				menu.stack = append(menu.stack, buildPlaylist(path))
			},
		})
	}
	return pls
}

func (tabs *sceneTags) Entry() *entry {
	return &tabs.entry
}

func (tabs *sceneTags) segueMount() {
	for i := range tabs.children {
		e := &tabs.children[i]

		if i == tabs.ptr {
			e.yp = 0.5
			e.labelAlpha = 1
			e.iconAlpha = 1
			e.scale = 0.75
			e.width = 500
		} else if i < tabs.ptr {
			e.yp = 0.5
			e.labelAlpha = 0
			e.iconAlpha = 1
			e.scale = 0.25
			e.width = 128
		} else if i > tabs.ptr {
			e.yp = 0.5
			e.labelAlpha = 0
			e.iconAlpha = 1
			e.scale = 0.25
			e.width = 128
		}
	}

	tabs.animate()
}

func (tabs *sceneTags) segueBack() {
	tabs.animate()
}

func (tabs *sceneTags) animate() {
	for i := range tabs.children {
		e := &tabs.children[i]

		var yp, labelAlpha, iconAlpha, scale, width float32
		if i == tabs.ptr {
			yp = 0.5
			labelAlpha = 1
			iconAlpha = 1
			scale = 0.75
			width = 500
		} else if i < tabs.ptr {
			yp = 0.5
			labelAlpha = 0
			iconAlpha = 1
			scale = 0.25
			width = 128
		} else if i > tabs.ptr {
			yp = 0.5
			labelAlpha = 0
			iconAlpha = 1
			scale = 0.25
			width = 128
		}

		menu.tweens[&e.yp] = gween.New(e.yp, yp, 0.15, ease.OutSine)
		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, labelAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, iconAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.scale] = gween.New(e.scale, scale, 0.15, ease.OutSine)
		menu.tweens[&e.width] = gween.New(e.width, width, 0.15, ease.OutSine)
		menu.tweens[&e.margin] = gween.New(e.margin, 0, 0.15, ease.OutSine)
	}
	menu.tweens[&menu.scroll] = gween.New(menu.scroll, float32(tabs.ptr*128), 0.15, ease.OutSine)
}

func (tabs *sceneTags) segueNext() {
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

func (tabs *sceneTags) update(dt float32) {
	// Right
	repeatRight(dt, input.NewState[0][libretro.DeviceIDJoypadRight], func() {
		tabs.ptr++
		if tabs.ptr >= len(tabs.children) {
			tabs.ptr = 0
		}
		tabs.animate()
	})

	// Left
	repeatLeft(dt, input.NewState[0][libretro.DeviceIDJoypadLeft], func() {
		tabs.ptr--
		if tabs.ptr < 0 {
			tabs.ptr = len(tabs.children) - 1
		}
		tabs.animate()
	})

	// OK
	if input.Released[0][libretro.DeviceIDJoypadA] {
		if tabs.children[tabs.ptr].callbackOK != nil {
			tabs.segueNext()
			tabs.children[tabs.ptr].callbackOK()
		}
	}
}

func (tabs sceneTags) render() {
	_, h := vid.Window.GetFramebufferSize()

	stackWidth := 710 * menu.ratio
	for i, e := range tabs.children {

		c := colorful.Hcl(float64(i)*20, 0.5, 0.5)

		x := -menu.scroll*menu.ratio + stackWidth + e.width/2*menu.ratio

		stackWidth += e.width*menu.ratio + e.margin*menu.ratio

		if e.labelAlpha > 0 {
			vid.Font.SetColor(float32(c.R), float32(c.B), float32(c.G), e.labelAlpha)
			lw := vid.Font.Width(0.6*menu.ratio, e.label)
			vid.Font.Printf(x-lw/2, float32(h)*e.yp+250*menu.ratio, 0.6*menu.ratio, e.label)
			lw = vid.Font.Width(0.4*menu.ratio, e.subLabel)
			vid.Font.Printf(x-lw/2, float32(h)*e.yp+330*menu.ratio, 0.4*menu.ratio, e.subLabel)
		}

		vid.DrawImage(menu.icons["hexagon"],
			x-220*e.scale*menu.ratio, float32(h)*e.yp-220*e.scale*menu.ratio,
			440*menu.ratio, 440*menu.ratio, e.scale, video.Color{R: float32(c.R), G: float32(c.B), B: float32(c.G), A: e.iconAlpha})

		vid.DrawImage(menu.icons[e.icon],
			x-128*e.scale*menu.ratio, float32(h)*e.yp-128*e.scale*menu.ratio,
			256*menu.ratio, 256*menu.ratio, e.scale, video.Color{R: 1, G: 1, B: 1, A: e.iconAlpha})
	}
}

func (tabs sceneTags) drawHintBar() {
	w, h := vid.Window.GetFramebufferSize()
	menu.ratio = float32(w) / 1920
	vid.DrawRect(0.0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 1.0, video.Color{R: 0.75, G: 0.75, B: 0.75, A: 1})

	var stack float32
	if state.Global.CoreRunning {
		stackHint(&stack, "key-p", "RESUME", h)
	}
	stackHint(&stack, "key-left-right", "NAVIGATE", h)
	stackHint(&stack, "key-x", "OPEN", h)
}
