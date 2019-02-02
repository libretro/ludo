package menu

import (
	"fmt"
	"os"
	"os/user"
	"regexp"
	"sort"
	"strings"

	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/playlists"
	"github.com/libretro/ludo/scanner"
	"github.com/libretro/ludo/utils"
	"github.com/libretro/ludo/video"
	colorful "github.com/lucasb-eyer/go-colorful"

	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type screenTabs struct {
	entry
}

func buildTabs() Scene {
	var list screenTabs
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
		e.children[e.ptr].width = 5200
		menu.scroll = float32(e.ptr*128 + 3030)
	}
}

// getPlaylists browse the filesystem for lpl files, parse them and returns
// a list of menu entries. It is used in the tabs, but could be used somewhere
// else too.
func getPlaylists() []entry {
	playlists.LoadPlaylists()

	// To store the keys in slice in sorted order
	var keys []string
	for k := range playlists.Playlists {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pls []entry
	for _, path := range keys {
		path := path
		filename := utils.Filename(path)
		count := playlistCount(path)
		pls = append(pls, entry{
			label:    playlistShortName(filename),
			subLabel: fmt.Sprintf("%d Games - 0 Favorites", count),
			icon:     filename,
			callbackOK: func() {
				menu.stack = append(menu.stack, buildPlaylist(path))
			},
		})
	}
	return pls
}

// playlistShortName shortens the name of some game systems that are too long
// to be displayed in the menu
func playlistShortName(in string) string {
	if len(in) < 20 {
		return in
	}
	r, _ := regexp.Compile(`(.*?) - (.*)`)
	out := r.ReplaceAllString(in, "$2")
	out = strings.Replace(out, "Nintendo Entertainment System", "NES", -1)
	out = strings.Replace(out, "PC Engine", "PCE", -1)
	return out
}

// Quick way of knowing how many games are in a playlist
func playlistCount(path string) int {
	file, _ := os.Open(path)
	c, _ := utils.LinesInFile(file)
	if c > 0 {
		return c / 6
	}
	return 0
}

func (tabs *screenTabs) Entry() *entry {
	return &tabs.entry
}

func (tabs *screenTabs) segueMount() {
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

func (tabs *screenTabs) segueBack() {
	tabs.animate()
}

func (tabs *screenTabs) animate() {
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
	}
	menu.tweens[&menu.scroll] = gween.New(menu.scroll, float32(tabs.ptr*128), 0.15, ease.OutSine)
}

func (tabs *screenTabs) segueNext() {
	cur := &tabs.children[tabs.ptr]
	menu.tweens[&cur.width] = gween.New(cur.width, 5200, 0.15, ease.OutSine)
	menu.tweens[&menu.scroll] = gween.New(menu.scroll, menu.scroll+3030, 0.15, ease.OutSine)
}

func (tabs *screenTabs) update(dt float32) {
	menu.inputCooldown -= dt
	if menu.inputCooldown < 0 {
		menu.inputCooldown = 0
	}

	// Right
	if input.NewState[0][libretro.DeviceIDJoypadRight] && menu.inputCooldown == 0 {
		tabs.ptr++
		if tabs.ptr >= len(tabs.children) {
			tabs.ptr = 0
		}
		tabs.animate()
		menu.inputCooldown = 0.15
	}

	// Left
	if input.NewState[0][libretro.DeviceIDJoypadLeft] && menu.inputCooldown == 0 {
		tabs.ptr--
		if tabs.ptr < 0 {
			tabs.ptr = len(tabs.children) - 1
		}
		tabs.animate()
		menu.inputCooldown = 0.15
	}

	// OK
	if input.Released[0][libretro.DeviceIDJoypadA] {
		if tabs.children[tabs.ptr].callbackOK != nil {
			tabs.segueNext()
			tabs.children[tabs.ptr].callbackOK()
		}
	}
}

func (tabs screenTabs) render() {
	_, h := vid.Window.GetFramebufferSize()

	stackWidth := 710 * menu.ratio
	for i, e := range tabs.children {

		c := colorful.Hcl(float64(i)*20, 0.5, 0.5)

		stackWidth += e.width * menu.ratio

		x := -menu.scroll*menu.ratio + stackWidth - e.width/2*menu.ratio

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

func (tabs screenTabs) drawHintBar() {
	w, h := vid.Window.GetFramebufferSize()
	c := video.Color{R: 0.25, G: 0.25, B: 0.25, A: 1}
	menu.ratio = float32(w) / 1920
	vid.DrawRect(0.0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 1.0, video.Color{R: 0.75, G: 0.75, B: 0.75, A: 1})
	vid.Font.SetColor(0.25, 0.25, 0.25, 1.0)

	stack := 30 * menu.ratio
	vid.DrawImage(menu.icons["key-left-right"], stack, float32(h)-70*menu.ratio, 70*menu.ratio, 70*menu.ratio, 1.0, c)
	stack += 70 * menu.ratio
	stack += 10 * menu.ratio
	vid.Font.Printf(stack, float32(h)-23*menu.ratio, 0.5*menu.ratio, "NAVIGATE")
	stack += vid.Font.Width(0.5*menu.ratio, "NAVIGATE")

	stack += 30 * menu.ratio
	vid.DrawImage(menu.icons["key-x"], stack, float32(h)-70*menu.ratio, 70*menu.ratio, 70*menu.ratio, 1.0, c)
	stack += 70 * menu.ratio
	stack += 10 * menu.ratio
	vid.Font.Printf(stack, float32(h)-23*menu.ratio, 0.5*menu.ratio, "OPEN")
	stack += vid.Font.Width(0.5*menu.ratio, "OPEN")
}
