package menu

import (
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/savestates"
	"github.com/libretro/ludo/state"
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
		label: "Save State",
		icon:  "savestate",
		callbackOK: func() {
			err := savestates.Save()
			if err != nil {
				ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			} else {
				ntf.DisplayAndLog(ntf.Success, "Menu", "State saved.")
			}
		},
	})

	list.children = append(list.children, entry{
		label: "Load State",
		icon:  "loadstate",
		callbackOK: func() {
			err := savestates.Load()
			if err != nil {
				ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			} else {
				state.Global.MenuActive = false
				ntf.DisplayAndLog(ntf.Success, "Menu", "State loaded.")
			}
		},
	})

	list.children = append(list.children, entry{
		label: "Take Screenshot",
		icon:  "screenshot",
		callbackOK: func() {
			vid.TakeScreenshot()
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
