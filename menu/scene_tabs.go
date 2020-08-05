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
				}))
			menu.focus--
		},
	})

	if state.Global.LudOS {
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
				askConfirmation(func() { cleanReboot() })
			},
		})

		list.children = append(list.children, entry{
			icon: "tab-shutdown",
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
	repeatLeft(dt, input.NewState[0][libretro.DeviceIDJoypadLeft], func() {
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
	repeatRight(dt, input.NewState[0][libretro.DeviceIDJoypadRight], func() {
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
	repeatDown(dt, input.NewState[0][libretro.DeviceIDJoypadDown], func() {
		audio.PlayEffect(audio.Effects["ok"])
		menu.t = 0
		menu.focus++
		s.animate()
	})

	// OK
	if input.Released[0][libretro.DeviceIDJoypadA] {
		audio.PlayEffect(audio.Effects["ok"])
		s.children[s.ptr].callbackOK()
		menu.t = 0
		menu.focus++
		s.animate()
	}
}

func (s sceneTabs) render() {
	w, _ := vid.Window.GetFramebufferSize()

	spacing := float32(96 + 32)
	totalWidth := spacing * float32(len(s.children)) * menu.ratio

	now := time.Now().Format("3:04PM")
	nowWidth := vid.BoldFont.Width(0.5*menu.ratio, now)
	vid.BoldFont.SetColor(0, 0, 0, 1)
	vid.BoldFont.Printf(
		float32(w)-96*menu.ratio-nowWidth,
		90*menu.ratio, 0.5*menu.ratio, now)

	for i, e := range s.children {
		if i == s.ptr && menu.focus == 1 {
			blink := float32(math.Cos(menu.t))
			vid.DrawImage(menu.icons["selection"],
				float32(w)/2-totalWidth/2+float32(i)*spacing*menu.ratio+96*menu.ratio/2-8*menu.ratio,
				32*menu.ratio-8*menu.ratio,
				96*menu.ratio+16*menu.ratio, 96*menu.ratio+16*menu.ratio, 1, 1,
				video.Color{R: 1, G: 1, B: 1, A: 1 - blink})
		}
		c := video.Color{R: 1, G: 1, B: 1, A: 1}
		vid.DrawImage(menu.icons["circle"],
			float32(w)/2-totalWidth/2+float32(i)*spacing*menu.ratio+96*menu.ratio/2,
			32*menu.ratio,
			96*menu.ratio, 96*menu.ratio, 1, 0, c)
		vid.DrawImage(menu.icons[e.icon],
			float32(w)/2-totalWidth/2+float32(i)*spacing*menu.ratio+96*menu.ratio/2+24*menu.ratio,
			56*menu.ratio,
			48*menu.ratio, 48*menu.ratio, 1, 0, c)
	}
}

func (s *sceneTabs) drawHintBar() {
}
