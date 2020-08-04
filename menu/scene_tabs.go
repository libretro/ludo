package menu

import (
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"

	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"
)

type sceneTabs struct {
	entry
}

func buildTabs() Scene {
	var list sceneTabs
	list.label = "Tabs"

	list.children = append(list.children, entry{
		icon: "tab-collection",
		callbackOK: func() {
			menu.stack = menu.stack[:len(menu.stack)-1]
			menu.Push(buildHome())
			menu.focus--
		},
	})

	list.children = append(list.children, entry{
		icon: "tab-settings",
		callbackOK: func() {
			menu.stack = menu.stack[:len(menu.stack)-1]
			menu.Push(buildSettings())
			menu.focus--
		},
	})

	if state.Global.LudOS {
		list.children = append(list.children, entry{
			icon: "updater",
			callbackOK: func() {
				list.segueNext()
				menu.Push(buildUpdater())
			},
		})

		list.children = append(list.children, entry{
			icon: "reboot",
			callbackOK: func() {
				askConfirmation(func() { cleanReboot() })
			},
		})

		list.children = append(list.children, entry{
			icon: "tab-quit",
			callbackOK: func() {
				askConfirmation(func() { cleanShutdown() })
			},
		})
	} else {
		list.children = append(list.children, entry{
			icon: "tab-quit",
			callbackOK: func() {
				askConfirmation(func() {
					vid.Window.SetShouldClose(true)
				})
			},
		})
	}

	list.segueMount()

	return &list
}

func (s *sceneTabs) Entry() *entry {
	return &s.entry
}

func (s *sceneTabs) segueMount() {
	for i := range s.children {
		e := &s.children[i]

		e.iconAlpha = 0.25
		if i == s.ptr {
			e.iconAlpha = 1
		}
	}

	s.animate()
}

func (s *sceneTabs) segueBack() {
	s.animate()
	s.winFocus()
}

func (s *sceneTabs) winFocus() {
	for i := range s.children {
		e := &s.children[i]

		labelAlpha := float32(0.5)
		if i == s.ptr {
			labelAlpha = 1
		}

		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, labelAlpha, 0.15, ease.OutSine)
	}

	menu.tweens[&s.alpha] = gween.New(s.alpha, 1, 0.15, ease.OutSine)
	menu.tweens[&s.width] = gween.New(s.width, 450, 0.15, ease.OutSine)
}

func (s *sceneTabs) loseFocus() {
	for i := range s.children {
		e := &s.children[i]
		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, 0, 0.15, ease.OutSine)
	}

	menu.tweens[&s.alpha] = gween.New(s.alpha, 0, 0.15, ease.OutSine)
	menu.tweens[&s.width] = gween.New(s.width, 150, 0.15, ease.OutSine)
}

func (s *sceneTabs) animate() {
	for i := range s.children {
		e := &s.children[i]

		var iconAlpha, labelAlpha float32
		if i == s.ptr {
			iconAlpha = 1
			labelAlpha = 1
		} else if i < s.ptr {
			iconAlpha = 0.5
			labelAlpha = 0.5
		} else if i > s.ptr {
			iconAlpha = 0.5
			labelAlpha = 0.5
		}
		if menu.focus != 1 {
			labelAlpha = 0
		}

		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, iconAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, labelAlpha, 0.15, ease.OutSine)
	}
}

// left tabs are never removed, we don't need to implement this callback
func (s *sceneTabs) segueNext() {
}

func (s *sceneTabs) update(dt float32) {
	// Down
	repeatDown(dt, input.NewState[0][libretro.DeviceIDJoypadDown], func() {
		s.ptr++
		if s.ptr >= len(s.children) {
			s.ptr = len(s.children) - 1
		} else {
			audio.PlayEffect(audio.Effects["down"])
			menu.t = 0
			s.animate()
			s.children[s.ptr].callbackOK()
		}
	})

	// Up
	repeatUp(dt, input.NewState[0][libretro.DeviceIDJoypadUp], func() {
		s.ptr--
		if s.ptr < 0 {
			s.ptr = 0
		} else {
			audio.PlayEffect(audio.Effects["up"])
			menu.t = 0
			s.animate()
			s.children[s.ptr].callbackOK()
		}
	})

	// Right
	repeatRight(dt, input.NewState[0][libretro.DeviceIDJoypadRight], func() {
		audio.PlayEffect(audio.Effects["ok"])
		menu.t = 0
		menu.focus++
		s.loseFocus()
	})

	// OK
	if input.Released[0][libretro.DeviceIDJoypadA] {
		audio.PlayEffect(audio.Effects["ok"])
		menu.t = 0
		menu.focus++
		s.loseFocus()
	}
}

func (s sceneTabs) render() {
	w, _ := vid.Window.GetFramebufferSize()

	spacing := float32(96 + 32)
	totalWidth := spacing * float32(len(s.children)) * menu.ratio

	for i, tab := range s.children {
		c := video.Color{R: 1, G: 1, B: 1, A: tab.iconAlpha}
		vid.DrawImage(menu.icons["circle"],
			float32(w)/2-totalWidth/2+float32(i)*spacing*menu.ratio+96*menu.ratio/2,
			32*menu.ratio,
			96*menu.ratio, 96*menu.ratio, 1, 0, c)
		vid.DrawImage(menu.icons[tab.icon],
			float32(w)/2-totalWidth/2+float32(i)*spacing*menu.ratio+96*menu.ratio/2+24*menu.ratio,
			56*menu.ratio,
			48*menu.ratio, 48*menu.ratio, 1, 0, c)
	}
}

func (s *sceneTabs) drawHintBar() {
}
