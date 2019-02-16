package menu

import (
	"time"

	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
)

type screenQuick struct {
	entry
}

func buildQuickMenu() Scene {
	var list screenQuick
	list.label = "Quick Menu"

	list.children = append(list.children, entry{
		label: "Resume",
		icon:  "resume",
		callbackOK: func() {
			state.Global.MenuActive = false
		},
	})

	list.children = append(list.children, entry{
		label: "Reset",
		icon:  "reset",
		callbackOK: func() {
			state.Global.Core.Reset()
			state.Global.MenuActive = false
		},
	})

	list.children = append(list.children, entry{
		label: "Savestates",
		icon:  "states",
		callbackOK: func() {
			list.segueNext()
			menu.stack = append(menu.stack, buildSavestates())
		},
	})

	list.children = append(list.children, entry{
		label: "Take Screenshot",
		icon:  "screenshot",
		callbackOK: func() {
			name := utils.Filename(state.Global.GamePath)
			date := time.Now().Format("2006-01-02-15-04-05")
			vid.TakeScreenshot(name + "@" + date)
			ntf.DisplayAndLog(ntf.Success, "Menu", "Took a screenshot.")
		},
	})

	list.children = append(list.children, entry{
		label: "Options",
		icon:  "subsetting",
		callbackOK: func() {
			list.segueNext()
			menu.stack = append(menu.stack, buildCoreOptions())
		},
	})

	list.segueMount()

	return &list
}

func (s *screenQuick) Entry() *entry {
	return &s.entry
}

func (s *screenQuick) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *screenQuick) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *screenQuick) segueBack() {
	genericAnimate(&s.entry)
}

func (s *screenQuick) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *screenQuick) render() {
	genericRender(&s.entry)
}

func (s *screenQuick) drawHintBar() {
	genericDrawHintBar()
}
