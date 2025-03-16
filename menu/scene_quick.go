package menu

import (
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
)

type sceneQuick struct {
	entry
}

func buildQuickMenu() Scene {
	var list sceneQuick
	list.label = "Quick Menu"

	list.children = append(list.children, entry{
		label: "Resume",
		icon:  "resume",
		callbackOK: func() {
			state.MenuActive = false
			state.FastForward = false
		},
	})

	list.children = append(list.children, entry{
		label: "Reset",
		icon:  "reset",
		callbackOK: func() {
			state.Core.Reset()
			state.MenuActive = false
			state.FastForward = false
		},
	})

	list.children = append(list.children, entry{
		label: "Savestates",
		icon:  "states",
		callbackOK: func() {
			list.segueNext()
			menu.Push(buildSavestates())
		},
	})

	list.children = append(list.children, entry{
		label: "Take Screenshot",
		icon:  "screenshot",
		callbackOK: func() {
			name := utils.DatedName(state.GamePath)
			err := menu.TakeScreenshot(name)
			if err != nil {
				ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			} else {
				ntf.DisplayAndLog(ntf.Success, "Menu", "Took a screenshot.")
			}
		},
	})

	list.children = append(list.children, entry{
		label: "Options",
		icon:  "options",
		callbackOK: func() {
			list.segueNext()
			menu.Push(buildCoreOptions())
		},
	})

	if state.Core != nil && state.Core.DiskControlCallback != nil {
		list.children = append(list.children, entry{
			label: "Disk Control",
			icon:  "core-disk-options",
			callbackOK: func() {
				list.segueNext()
				menu.Push(buildCoreDiskControl())
			},
		})
	}

	list.children = append(list.children, entry{
		label: "Settings",
		icon:  "subsetting",
		callbackOK: func() {
			list.segueNext()
			menu.Push(buildSettings())
		},
	})

	list.segueMount()

	return &list
}

func (s *sceneQuick) Entry() *entry {
	return &s.entry
}

func (s *sceneQuick) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneQuick) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneQuick) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneQuick) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *sceneQuick) render() {
	genericRender(&s.entry)
}

func (s *sceneQuick) drawHintBar() {
	genericDrawHintBar()
}
