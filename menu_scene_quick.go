package main

import (
	"github.com/libretro/go-playthemall/notifications"
	"github.com/libretro/go-playthemall/savestates"
	"github.com/libretro/go-playthemall/state"
)

type screenQuick struct {
	entry
}

func buildQuickMenu() scene {
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
				notifications.DisplayAndLog("Menu", err.Error())
			} else {
				notifications.DisplayAndLog("Menu", "State saved.")
			}
		},
	})

	list.children = append(list.children, entry{
		label: "Load State",
		icon:  "loadstate",
		callbackOK: func() {
			err := savestates.Load()
			if err != nil {
				notifications.DisplayAndLog("Menu", err.Error())
			} else {
				state.Global.MenuActive = false
				notifications.DisplayAndLog("Menu", "State loaded.")
			}
		},
	})

	list.children = append(list.children, entry{
		label: "Take Screenshot",
		icon:  "screenshot",
		callbackOK: func() {
			takeScreenshot()
			notifications.DisplayAndLog("Menu", "Took a screenshot.")
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

func (s *screenQuick) update() {
	genericInput(&s.entry)
}

func (s *screenQuick) render() {
	genericRender(&s.entry)
}
