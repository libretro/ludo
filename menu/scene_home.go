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

	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type sceneHome struct {
	entry
	yptr     int
	yscroll  float32
	xscrolls []float32
	xptrs    []int
}

func buildHome() Scene {
	var list sceneHome
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
			strippedName, tags := extractTags(game.Name)
			list.children[cat].children = append(list.children[cat].children, entry{
				label:    strippedName,
				gameName: game.Name,
				tags:     tags,
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

	if len(list.children) == 0 {
		menu.focus--
	}

	list.segueMount()

	return &list
}

func (s *sceneHome) Entry() *entry {
	return &s.entry
}

func (s *sceneHome) segueMount() {
	s.alpha = 0
	for j := range s.children {
		s.xscrolls[j] = 0
	}
	s.y = 300

	for j := range s.children {
		ve := &s.children[j]
		ve.labelAlpha = 0
		ve.height = 504 + 136

		for i := range ve.children {
			e := &s.children[j].children[i]

			if i == s.xptrs[j] {
				e.labelAlpha = 1
				e.iconAlpha = 1
				e.scale = 2.1
				e.borderAlpha = 0
			} else if i < s.xptrs[j] {
				e.labelAlpha = 0
				e.iconAlpha = 0
				e.scale = 1
				e.borderAlpha = 0
			} else {
				e.labelAlpha = 0
				e.iconAlpha = 1
				e.scale = 1
				e.borderAlpha = 0
			}
		}
	}

	s.animate()
}

func (s *sceneHome) segueBack() {
	s.animate()
}

func (s *sceneHome) animate() {
	for j := range s.children {
		ve := &s.children[j]

		labelAlpha := float32(1)
		if j < s.yptr {
			labelAlpha = 0
		}
		menu.tweens[&ve.labelAlpha] = gween.New(ve.labelAlpha, labelAlpha, 0.15, ease.OutSine)
		menu.tweens[&ve.height] = gween.New(ve.height, 504+136, 0.15, ease.OutSine)

		for i := range ve.children {
			e := &s.children[j].children[i]

			var labelAlpha, iconAlpha, scale, borderAlpha float32
			if i == s.xptrs[j] {
				labelAlpha = 1
				iconAlpha = 1
				scale = 2.1
				borderAlpha = 1
			} else {
				labelAlpha = 0
				iconAlpha = 1
				scale = 1
				borderAlpha = 0
			}
			if j < s.yptr {
				labelAlpha = 0
				iconAlpha = 0
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
	menu.tweens[&s.y] = gween.New(s.y, 0, 0.15, ease.OutSine)
}

func (s *sceneHome) segueNext() {
	menu.tweens[&s.alpha] = gween.New(s.alpha, 0, 0.15, ease.OutSine)
	menu.tweens[&s.y] = gween.New(s.y, -300, 0.15, ease.OutSine)

	for j := range s.children {
		ve := &s.children[j]
		for i := range ve.children {
			e := &s.children[j].children[i]
			menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, 0, 0.15, ease.OutSine)
		}
	}
}

func (s *sceneHome) update(dt float32) {
	// Empty state
	if len(s.children) == 0 {
		menu.focus--
		return
	}

	// Right
	repeatRight(dt, input.NewState[0][libretro.DeviceIDJoypadRight] == 1, func() {
		if s.xptrs[s.yptr] < len(s.children[s.yptr].children)-1 {
			s.xptrs[s.yptr]++
			audio.PlayEffect(audio.Effects["down"])
			menu.t = 0
			s.animate()
		}
	})

	// Left
	repeatLeft(dt, input.NewState[0][libretro.DeviceIDJoypadLeft] == 1, func() {
		if s.xptrs[s.yptr] > 0 {
			s.xptrs[s.yptr]--
			audio.PlayEffect(audio.Effects["up"])
			menu.t = 0
			s.animate()
		}
	})

	// Down
	repeatDown(dt, input.NewState[0][libretro.DeviceIDJoypadDown] == 1, func() {
		if s.yptr < len(s.children)-1 {
			s.yptr++
			audio.PlayEffect(audio.Effects["down"])
			menu.t = 0
			s.animate()
		}
	})

	// Up
	repeatUp(dt, input.NewState[0][libretro.DeviceIDJoypadUp] == 1, func() {
		if s.yptr > 0 {
			s.yptr--
			audio.PlayEffect(audio.Effects["up"])
			menu.t = 0
			s.animate()
		} else if s.yptr == 0 && len(menu.stack) > 1 {
			audio.PlayEffect(audio.Effects["cancel"])
			menu.stack[len(menu.stack)-2].segueBack()
			menu.focus--
			menu.t = 0
		}
	})

	// OK
	if input.Released[0][libretro.DeviceIDJoypadA] == 1 {
		if len(s.children) > 0 && s.children[s.yptr].children[s.xptrs[s.yptr]].callbackOK != nil {
			audio.PlayEffect(audio.Effects["ok"])
			s.segueNext()
			s.children[s.yptr].children[s.xptrs[s.yptr]].callbackOK()
		}
	}

	// Cancel
	if input.Released[0][libretro.DeviceIDJoypadB] == 1 {
		if len(menu.stack) > 1 {
			audio.PlayEffect(audio.Effects["cancel"])
			menu.stack[len(menu.stack)-2].segueBack()
			menu.focus--
		}
	}
}

func (s sceneHome) render() {
	vst := float32(0)
	for j, ve := range s.children {
		ve := ve

		menu.BoldFont.SetColor(titleColor.Alpha(ve.labelAlpha * s.alpha))
		menu.BoldFont.Printf(
			96*menu.ratio,
			s.y*menu.ratio+230*menu.ratio+vst*menu.ratio-s.yscroll*menu.ratio,
			0.5*menu.ratio, ve.label)

		y := s.y + 272 + vst - s.yscroll

		vst += ve.height

		// performance improvement
		if math.Abs(float64(j-s.yptr)) > 1 {
			continue
		}

		stackWidth := float32(96)
		for i, e := range ve.children {
			x := -s.xscrolls[j] + stackWidth

			stackWidth += 320*e.scale + e.margin + 32

			// performance improvement
			if math.Abs(float64(i-s.xptrs[j])) > 4 {
				freeThumbnail(&ve, i)
				continue
			}

			if menu.focus == 2 && j == s.yptr && i == s.xptrs[s.yptr] {
				blink := float32(math.Cos(menu.t))
				menu.DrawImage(
					menu.icons["selection"],
					x*menu.ratio-8*menu.ratio,
					y*menu.ratio-8*menu.ratio,
					320*e.scale*menu.ratio+16*menu.ratio, 240*e.scale*menu.ratio+16*menu.ratio,
					1, 0.1, white.Alpha((e.borderAlpha-blink)*s.alpha))
			}

			drawThumbnail(
				&ve, i,
				e.system, e.gameName,
				x*menu.ratio,
				y*menu.ratio,
				320*e.scale*menu.ratio, 240*e.scale*menu.ratio,
				1, white.Alpha(e.iconAlpha*s.alpha))

			menu.DrawImage(
				menu.icons["border"],
				x*menu.ratio,
				y*menu.ratio,
				320*e.scale*menu.ratio, 240*e.scale*menu.ratio,
				1, 0.07, white.Alpha(e.iconAlpha*s.alpha))

			menu.BoldFont.SetColor(textColor.Alpha(e.labelAlpha * s.alpha))
			menu.BoldFont.Printf(
				(x+672+32)*menu.ratio,
				(y+360)*menu.ratio,
				0.7*menu.ratio, e.label)

			menu.BoldFont.SetColor(mediumGrey.Alpha(e.labelAlpha * s.alpha))
			menu.BoldFont.Printf(
				(x+672+32)*menu.ratio,
				(y+430)*menu.ratio,
				0.5*menu.ratio, e.subLabel)

			stack := (x + 672 + 32) * menu.ratio
			for _, tag := range e.tags {
				if _, ok := menu.icons[tag]; ok {
					menu.DrawRect(stack-1*menu.ratio, (y+500-35)*menu.ratio-1*menu.ratio,
						48*menu.ratio+2*menu.ratio, 35*menu.ratio+2*menu.ratio, 0.22,
						mediumGrey.Alpha(e.labelAlpha*s.alpha))
					menu.DrawImage(
						menu.icons[tag],
						stack, (y+500-35)*menu.ratio,
						48*menu.ratio, 35*menu.ratio, 1.0, 0.2,
						white.Alpha(e.labelAlpha*s.alpha))
					stack += 48 * menu.ratio
					stack += 24 * menu.ratio
				}
			}
		}
	}

	if len(s.children) == 0 {
		w, h := menu.Window.GetFramebufferSize()
		menu.BoldFont.SetColor(black)
		msg := "Welcome to Ludo, please scan your collection."
		msgw := menu.BoldFont.Width(0.5*menu.ratio, msg)
		menu.BoldFont.Printf(
			float32(w)/2-msgw/2,
			float32(h)/2,
			0.5*menu.ratio, msg)
	}
}

func (s sceneHome) drawHintBar() {
	w, h := menu.Window.GetFramebufferSize()
	menu.DrawRect(0, float32(h)-88*menu.ratio, float32(w), 88*menu.ratio, 0, hintBgColor)
	menu.DrawRect(0, float32(h)-88*menu.ratio, float32(w), 2*menu.ratio, 0, sepColor)

	arrows, _, _, a, b, _, _, _, _, guide := hintIcons()

	lstack := float32(75) * menu.ratio
	rstack := float32(w) - 96*menu.ratio
	stackHintLeft(&lstack, arrows, "Navigate", h)
	stackHintRight(&rstack, a, "Run", h)
	stackHintRight(&rstack, b, "Back", h)
	if state.CoreRunning {
		stackHintRight(&rstack, guide, "Resume", h)
	}
}
