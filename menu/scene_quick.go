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
			menu.Push(buildSavestates())
		},
	})

	list.children = append(list.children, entry{
		label: "Take Screenshot",
		icon:  "screenshot",
		callbackOK: func() {
			name := utils.DatedName(state.Global.GamePath)
			vid.TakeScreenshot(name)
			ntf.DisplayAndLog(ntf.Success, "Menu", "Took a screenshot.")
		},
	})

	list.children = append(list.children, entry{
		label: "Options",
		icon:  "subsetting",
		callbackOK: func() {
			list.segueNext()
			menu.Push(buildCoreOptions())
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
