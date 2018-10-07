package menu

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/libretro/go-playthemall/input"
	"github.com/libretro/go-playthemall/libretro"
	"github.com/libretro/go-playthemall/scanner"
	"github.com/libretro/go-playthemall/utils"
	"github.com/libretro/go-playthemall/video"

	colorful "github.com/lucasb-eyer/go-colorful"

	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type screenTabs struct {
	entry
}

func buildTabs() Scene {
	var list screenTabs
	list.label = "Play Them All"

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
		subLabel: "Configure Play Them All",
		icon:     "setting",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	usr, _ := user.Current()
	paths, _ := filepath.Glob(usr.HomeDir + "/.playthemall/playlists/*.lpl")

	for _, path := range paths {
		path := path
		filename := utils.Filename(path)
		count := playlistCount(path)
		list.children = append(list.children, entry{
			label:    playlistShortName(filename),
			subLabel: fmt.Sprintf("%d Games - 0 Favorites", count),
			icon:     filename,
			callbackOK: func() {
				menu.stack = append(menu.stack, buildPlaylist(path))
			},
		})
	}

	list.children = append(list.children, entry{
		label:    "Add games",
		subLabel: "Scan your collection",
		icon:     "add",
		callbackOK: func() {
			usr, _ := user.Current()
			menu.stack = append(menu.stack, buildExplorer(usr.HomeDir, nil, scanner.ScanDir,
				&entry{
					label: "<Scan this directory>",
					icon:  "scan",
				}))
		},
	})

	list.segueMount()

	return &list
}

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
			e.scale = 1
			e.width = 1256
		} else if i < tabs.ptr {
			e.yp = 0.05
			e.labelAlpha = 0
			e.iconAlpha = 0.5
			e.scale = 0.25
			e.width = 128
		} else if i > tabs.ptr {
			e.yp = 0.95
			e.labelAlpha = 0
			e.iconAlpha = 0.5
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
			scale = 1
			width = 1256
		} else if i < tabs.ptr {
			yp = 0.05
			labelAlpha = 0
			iconAlpha = 0.5
			scale = 0.25
			width = 128
		} else if i > tabs.ptr {
			yp = 0.95
			labelAlpha = 0
			iconAlpha = 0.5
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
	menu.tweens[&cur.width] = gween.New(cur.width, 4000, 0.15, ease.OutSine)
	menu.tweens[&menu.scroll] = gween.New(menu.scroll, menu.scroll+2100, 0.15, ease.OutSine)
}

func (tabs *screenTabs) update() {
	if menu.inputCooldown > 0 {
		menu.inputCooldown--
	}

	// Right
	if input.NewState[0][libretro.DeviceIDJoypadRight] && menu.inputCooldown == 0 {
		tabs.ptr++
		if tabs.ptr >= len(tabs.children) {
			tabs.ptr = 0
		}
		tabs.animate()
		menu.inputCooldown = 10
	}

	// Left
	if input.NewState[0][libretro.DeviceIDJoypadLeft] && menu.inputCooldown == 0 {
		tabs.ptr--
		if tabs.ptr < 0 {
			tabs.ptr = len(tabs.children) - 1
		}
		tabs.animate()
		menu.inputCooldown = 10
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

	stackWidth := 132 * menu.ratio
	for i, e := range tabs.children {

		c := colorful.Hcl(float64(i)*20, 0.5, 0.5)
		var alpha float32 = 1
		if i == 0 {
			alpha = 0
		}

		vid.DrawQuad(
			-menu.scroll*menu.ratio+stackWidth, 0,
			-menu.scroll*menu.ratio+stackWidth+e.width*menu.ratio, 0,
			-menu.scroll*menu.ratio+stackWidth+400*menu.ratio, float32(h),
			-menu.scroll*menu.ratio+stackWidth+400*menu.ratio+e.width*menu.ratio, float32(h),
			video.Color{R: float32(c.R), G: float32(c.B), B: float32(c.G), A: alpha}) // intentional mix

		stackWidth += e.width * menu.ratio

		x := -menu.scroll*menu.ratio + stackWidth - e.width/2*menu.ratio + 400*menu.ratio - 400*e.yp*menu.ratio

		if e.labelAlpha > 0 {
			vid.Font.SetColor(1.0, 1.0, 1.0, e.labelAlpha)
			lw := vid.Font.Width(0.7*menu.ratio, e.label)
			vid.Font.Printf(x-lw/2, float32(h)*e.yp+180*menu.ratio, 0.7*menu.ratio, e.label)
			lw = vid.Font.Width(0.4*menu.ratio, e.subLabel)
			vid.Font.Printf(x-lw/2, float32(h)*e.yp+260*menu.ratio, 0.4*menu.ratio, e.subLabel)
		}

		vid.DrawImage(menu.icons[e.icon],
			x-128*e.scale*menu.ratio, float32(h)*e.yp-128*e.scale*menu.ratio,
			256*menu.ratio, 256*menu.ratio, e.scale, video.Color{R: 1, G: 1, B: 1, A: e.iconAlpha})
	}
}
