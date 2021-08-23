package menu

import (
	"math"
	"os/user"
	"time"

	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/scanner"
	"github.com/libretro/ludo/state"
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
		icon: "file",
		callbackOK: func() {
			menu.stack = menu.stack[:len(menu.stack)-1]
			menu.Push(buildMainMenu())
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

	list.children = append(list.children, entry{
		icon: "tab-scan",
		callbackOK: func() {
			usr, _ := user.Current()
			menu.stack = menu.stack[:len(menu.stack)-1]
			menu.Push(buildExplorer(usr.HomeDir, nil,
				func(path string) { scanner.ScanDir(path, nil) },
				&entry{
					label: "<Scan this directory>",
					icon:  "scan",
				},
				nil,
			))
			menu.focus--
		},
	})

	if state.LudOS {
		list.children = append(list.children, entry{
			icon: "tab-updater",
			callbackOK: func() {
				menu.stack = menu.stack[:len(menu.stack)-1]
				menu.Push(buildUpdater())
				menu.focus--
			},
		})

		list.children = append(list.children, entry{
			icon: "tab-reboot",
			callbackOK: func() {
				askQuitConfirmation(func() { cleanReboot() })
			},
		})

		list.children = append(list.children, entry{
			icon: "tab-shutdown",
			callbackOK: func() {
				askQuitConfirmation(func() { cleanShutdown() })
			},
		})
	} else {
		list.children = append(list.children, entry{
			icon: "tab-quit",
			callbackOK: func() {
				askQuitConfirmation(func() {
					menu.SetShouldClose(true)
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
	s.animate()
}

func (s *sceneTabs) segueBack() {
	s.animate()
}

func (s *sceneTabs) animate() {
}

// left tabs are never removed, we don't need to implement this callback
func (s *sceneTabs) segueNext() {
}

func (s *sceneTabs) update(dt float32) {
	// Left
	repeatLeft(dt, input.NewState[0][libretro.DeviceIDJoypadLeft] == 1, func() {
		s.ptr--
		if s.ptr < 0 {
			s.ptr = 0
		} else {
			audio.PlayEffect(audio.Effects["up"])
			menu.t = 0
			s.animate()
		}
	})

	// Right
	repeatRight(dt, input.NewState[0][libretro.DeviceIDJoypadRight] == 1, func() {
		s.ptr++
		if s.ptr >= len(s.children) {
			s.ptr = len(s.children) - 1
		} else {
			audio.PlayEffect(audio.Effects["down"])
			menu.t = 0
			s.animate()
		}
	})

	// Down
	repeatDown(dt, input.NewState[0][libretro.DeviceIDJoypadDown] == 1, func() {
		audio.PlayEffect(audio.Effects["ok"])
		menu.t = 0
		menu.focus++
		s.animate()
	})

	// OK
	if input.Released[0][libretro.DeviceIDJoypadA] == 1 {
		audio.PlayEffect(audio.Effects["ok"])
		s.children[s.ptr].callbackOK()
		menu.t = 0
		menu.focus++
		s.animate()
	}
}

func (s sceneTabs) render() {
	w, _ := menu.Window.GetFramebufferSize()

	now := time.Now().Format("3:04PM")
	nowWidth := menu.BoldFont.Width(0.5*menu.ratio, now)
	menu.BoldFont.SetColor(black)
	menu.BoldFont.Printf(
		float32(w)-96*menu.ratio-nowWidth,
		90*menu.ratio, 0.5*menu.ratio, now)

	if menu.focus > 2 {
		return
	}

	spacing := float32(96 + 32)
	totalWidth := spacing * float32(len(s.children)) * menu.ratio

	for i, e := range s.children {
		if i == s.ptr && menu.focus == 1 {
			blink := float32(math.Cos(menu.t))
			menu.DrawImage(menu.icons["selection"],
				float32(w)/2-totalWidth/2+float32(i)*spacing*menu.ratio+96*menu.ratio/2-8*menu.ratio,
				32*menu.ratio-8*menu.ratio,
				96*menu.ratio+16*menu.ratio, 96*menu.ratio+16*menu.ratio, 1, 1,
				white.Alpha(1-blink))
		}
		menu.DrawImage(menu.icons["circle"],
			float32(w)/2-totalWidth/2+float32(i)*spacing*menu.ratio+96*menu.ratio/2,
			32*menu.ratio,
			96*menu.ratio, 96*menu.ratio, 1, 0, white)
		menu.DrawImage(menu.icons[e.icon],
			float32(w)/2-totalWidth/2+float32(i)*spacing*menu.ratio+96*menu.ratio/2+24*menu.ratio,
			56*menu.ratio,
			48*menu.ratio, 48*menu.ratio, 1, 0, blue)
	}
}

func (s *sceneTabs) drawHintBar() {
	w, h := menu.Window.GetFramebufferSize()
	menu.DrawRect(0, float32(h)-88*menu.ratio, float32(w), 88*menu.ratio, 0, white)
	menu.DrawRect(0, float32(h)-88*menu.ratio, float32(w), 2*menu.ratio, 0, lightGrey)

	_, _, leftRight, a, _, _, _, _, _, guide := hintIcons()

	lstack := float32(75) * menu.ratio
	rstack := float32(w) - 96*menu.ratio
	stackHintLeft(&lstack, leftRight, "Navigate", h)
	stackHintRight(&rstack, a, "Ok", h)
	if state.CoreRunning {
		stackHintRight(&rstack, guide, "Resume", h)
	}
}
