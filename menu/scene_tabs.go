package menu

import (
	"math"
	"sort"

	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/history"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/playlists"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
	"github.com/libretro/ludo/video"

	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type sceneTabs struct {
	entry
	yptr     int
	yscroll  float32
	xscrolls []float32
	xptrs    []int
}

func buildTabs() Scene {
	var list sceneTabs
	list.label = "Home"

	cat := 0
	history.Load()
	if len(history.List) > 0 {
		list.children = append(list.children, entry{
			label: "Recently played",
		})
		list.xscrolls = append(list.xscrolls, 0)
		list.xptrs = append(list.xptrs, 0)

		for _, game := range history.List {
			game := game
			strippedName, _ := extractTags(game.Name)
			list.children[cat].children = append(list.children[cat].children, entry{
				label:    strippedName,
				gameName: game.Name,
				subLabel: game.System,
				system:   game.System,
				callbackOK: func() {
					loadHistoryEntry(&list, game)
				},
			})
		}
		cat++
	}

	playlists.Load()

	// To store the keys in slice in sorted order
	var keys []string
	for k := range playlists.Playlists {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, path := range keys {
		path := path
		filename := utils.FileName(path)

		list.children = append(list.children, entry{
			label: filename,
		})
		list.xscrolls = append(list.xscrolls, 0)
		list.xptrs = append(list.xptrs, 0)

		for _, game := range playlists.Playlists[path] {
			game := game
			strippedName, tags := extractTags(game.Name)
			list.children[cat].children = append(list.children[cat].children, entry{
				label:      strippedName,
				gameName:   game.Name,
				path:       game.Path,
				tags:       tags,
				icon:       utils.FileName(path) + "-content",
				subLabel:   filename,
				system:     filename,
				callbackOK: func() { loadPlaylistEntry(&list, filename, game) },
			})
		}
		cat++
	}

	list.segueMount()

	return &list
}

func (s *sceneTabs) Entry() *entry {
	return &s.entry
}

func (s *sceneTabs) segueMount() {
	s.alpha = 0
	for j := range s.children {
		s.xscrolls[j] = 0
	}
	s.yscroll = -500

	for j := range s.children {
		ve := &s.children[j]
		ve.labelAlpha = 0
		ve.height = 504 + 136
		//if j == s.yptr {
		ve.height = 240 + 136
		//}

		for i := range ve.children {
			e := &s.children[j].children[i]

			if i == s.xptrs[j] {
				e.labelAlpha = 0
				e.iconAlpha = 0
				e.scale = 2.1
				e.borderAlpha = 0
			} else if i < s.xptrs[j] {
				e.labelAlpha = 0
				e.iconAlpha = 0
				e.scale = 1
				e.borderAlpha = 0
			} else {
				e.labelAlpha = 0
				e.iconAlpha = 0
				e.scale = 1
				e.borderAlpha = 0
			}
		}
	}

	s.animate()
}

func (s *sceneTabs) segueBack() {
	s.animate()
}

func (s *sceneTabs) animate() {
	for j := range s.children {
		ve := &s.children[j]

		labelAlpha := float32(1)
		if j < s.yptr {
			labelAlpha = 0
		}
		menu.tweens[&ve.labelAlpha] = gween.New(ve.labelAlpha, labelAlpha, 0.15, ease.OutSine)
		height := float32(240 + 136)
		//if j == s.yptr {
		height = 504 + 136
		//}
		menu.tweens[&ve.height] = gween.New(ve.height, height, 0.15, ease.OutSine)

		for i := range ve.children {
			e := &s.children[j].children[i]

			var labelAlpha, iconAlpha, scale, borderAlpha float32
			if i == s.xptrs[j] {
				labelAlpha = 1
				iconAlpha = 1
				scale = 2.1
				borderAlpha = 1
			} else if j < s.yptr {
				labelAlpha = 0
				iconAlpha = 0
				scale = 1
				borderAlpha = 0
			} else {
				labelAlpha = 0
				iconAlpha = 1
				scale = 1
				borderAlpha = 0
			}

			menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, labelAlpha, 0.15, ease.OutSine)
			menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, iconAlpha, 0.15, ease.OutSine)
			menu.tweens[&e.borderAlpha] = gween.New(e.borderAlpha, borderAlpha, 0.15, ease.OutSine)
			menu.tweens[&e.scale] = gween.New(e.scale, scale, 0.15, ease.OutSine)
		}
	}

	for j := range s.children {
		menu.tweens[&s.xscrolls[j]] = gween.New(s.xscrolls[j], float32(s.xptrs[j]*(320+32)), 0.15, ease.OutSine)
	}

	vst := float32(0)
	for j := range s.children {
		if j == s.yptr {
			break
		}
		vst += 504 + 136
	}

	menu.tweens[&s.yscroll] = gween.New(s.yscroll, vst, 0.15, ease.OutSine)
	menu.tweens[&s.alpha] = gween.New(s.alpha, 1, 0.15, ease.OutSine)
}

func (s *sceneTabs) segueNext() {
	menu.tweens[&s.alpha] = gween.New(s.alpha, 0, 0.15, ease.OutSine)
	menu.tweens[&s.yscroll] = gween.New(s.yscroll, s.yscroll+300, 0.15, ease.OutSine)

	for j := range s.children {
		ve := &s.children[j]
		for i := range ve.children {
			e := &s.children[j].children[i]
			menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, 0, 0.15, ease.OutSine)
		}
	}
}

func (s *sceneTabs) update(dt float32) {
	// Right
	repeatRight(dt, input.NewState[0][libretro.DeviceIDJoypadRight], func() {
		if s.xptrs[s.yptr] < len(s.children[s.yptr].children)-1 {
			s.xptrs[s.yptr]++
			audio.PlayEffect(audio.Effects["down"])
			menu.t = 0
			s.animate()
		}
	})

	// Left
	repeatLeft(dt, input.NewState[0][libretro.DeviceIDJoypadLeft], func() {
		if s.xptrs[s.yptr] > 0 {
			s.xptrs[s.yptr]--
			audio.PlayEffect(audio.Effects["up"])
			menu.t = 0
			s.animate()
		}
	})

	// Down
	repeatDown(dt, input.NewState[0][libretro.DeviceIDJoypadDown], func() {
		if s.yptr < len(s.children)-1 {
			s.yptr++
			audio.PlayEffect(audio.Effects["down"])
			menu.t = 0
			s.animate()
		}
	})

	// Up
	repeatUp(dt, input.NewState[0][libretro.DeviceIDJoypadUp], func() {
		if s.yptr > 0 {
			s.yptr--
			audio.PlayEffect(audio.Effects["up"])
			menu.t = 0
			s.animate()
		}
	})

	// OK
	if input.Released[0][libretro.DeviceIDJoypadA] {
		if s.children[s.yptr].children[s.xptrs[s.yptr]].callbackOK != nil {
			audio.PlayEffect(audio.Effects["ok"])
			s.segueNext()
			s.children[s.yptr].children[s.xptrs[s.yptr]].callbackOK()
		}
	}
}

func (s sceneTabs) render() {
	vst := float32(0)
	for j, ve := range s.children {
		ve := ve

		vid.BoldFont.SetColor(0.129, 0.441, 0.684, ve.labelAlpha*s.alpha)
		vid.BoldFont.Printf(
			96*menu.ratio,
			230*menu.ratio+vst*menu.ratio-s.yscroll*menu.ratio,
			0.5*menu.ratio, ve.label)

		y := 272 + vst - s.yscroll

		vst += ve.height

		if y < -400 || y > 1080 {
			continue
		}

		stackWidth := float32(96)
		for i, e := range ve.children {
			x := -s.xscrolls[j] + stackWidth

			stackWidth += 320*e.scale + e.margin + 32

			if x < -400 || x > 1920 {
				continue
			}

			blink := float32(1)
			if j == s.yptr && i == s.xptrs[s.yptr] {
				blink = float32(math.Cos(menu.t))
			}

			vid.DrawImage(
				menu.icons["selection"],
				x*menu.ratio-8*menu.ratio,
				y*menu.ratio-8*menu.ratio,
				320*e.scale*menu.ratio+16*menu.ratio, 240*e.scale*menu.ratio+16*menu.ratio,
				1, 0.1, video.Color{R: 1, G: 1, B: 1, A: (e.borderAlpha - blink) * s.alpha})

			drawThumbnail(
				&ve, i,
				e.system, e.gameName,
				x*menu.ratio,
				y*menu.ratio,
				320*e.scale*menu.ratio, 240*e.scale*menu.ratio,
				1, video.Color{R: 1, G: 1, B: 1, A: e.iconAlpha * s.alpha})

			vid.DrawImage(
				menu.icons["border"],
				x*menu.ratio,
				y*menu.ratio,
				320*e.scale*menu.ratio, 240*e.scale*menu.ratio,
				1, 0.07, video.Color{R: 1, G: 1, B: 1, A: e.iconAlpha * s.alpha})

			vid.BoldFont.SetColor(0, 0, 0, e.labelAlpha*s.alpha)
			vid.BoldFont.Printf(
				(x+672+32)*menu.ratio,
				(y+360)*menu.ratio,
				0.7*menu.ratio, e.label)

			vid.BoldFont.SetColor(0.56, 0.56, 0.56, e.labelAlpha*s.alpha)
			vid.BoldFont.Printf(
				(x+672+32)*menu.ratio,
				(y+430)*menu.ratio,
				0.5*menu.ratio, e.subLabel)
		}
	}
}

func (s sceneTabs) drawHintBar() {
	w, h := vid.Window.GetFramebufferSize()
	vid.DrawRect(0, float32(h)-88*menu.ratio, float32(w), 88*menu.ratio, 0, video.Color{R: 1, G: 1, B: 1, A: 1})
	vid.DrawRect(0, float32(h)-88*menu.ratio, float32(w), 2*menu.ratio, 0, video.Color{R: 0.85, G: 0.85, B: 0.85, A: 1})

	arrows, _, _, a, _, _, _, _, _, guide := hintIcons()

	lstack := float32(75) * menu.ratio
	rstack := float32(w) - 96*menu.ratio
	stackHintLeft(&lstack, arrows, "Navigate", h)
	stackHintRight(&rstack, a, "Open", h)
	if state.Global.CoreRunning {
		stackHintRight(&rstack, guide, "Resume", h)
	}
}
