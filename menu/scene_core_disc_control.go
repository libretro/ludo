package menu

import (
	"fmt"

	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/state"
)

type sceneCoreDiscControl struct {
	entry
}

func buildCoreDiscControl() Scene {
	var list sceneCoreDiscControl
	list.label = "Core Disc Control"

	for i := uint(0); i < state.Core.DiskControlCallback.GetNumImages(); i++ {
		index := i
		list.children = append(list.children, entry{
			label: fmt.Sprintf("Disc %d", index+1),
			icon:  "subsetting",
			stringValue: func() string {
				if index == state.Core.DiskControlCallback.GetImageIndex() {
					return "Active"
				}
				return ""
			},
			callbackOK: func() {
				if index == state.Core.DiskControlCallback.GetImageIndex() {
					return
				}
				state.Core.DiskControlCallback.SetEjectState(true)
				state.Core.DiskControlCallback.SetImageIndex(index)
				state.Core.DiskControlCallback.SetEjectState(false)
				ntf.DisplayAndLog(ntf.Success, "Menu", "Switched to disk %d.", index+1)
				state.MenuActive = false
			},
		})
	}

	list.segueMount()

	return &list
}

func (s *sceneCoreDiscControl) Entry() *entry {
	return &s.entry
}

func (s *sceneCoreDiscControl) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneCoreDiscControl) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneCoreDiscControl) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneCoreDiscControl) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *sceneCoreDiscControl) render() {
	genericRender(&s.entry)
}

func (s *sceneCoreDiscControl) drawHintBar() {
	w, h := menu.GetFramebufferSize()
	menu.DrawRect(0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 0, lightGrey)

	_, upDown, leftRight, _, b, _, _, _, _, guide := hintIcons()

	var stack float32
	if state.CoreRunning {
		stackHint(&stack, guide, "RESUME", h)
	}
	stackHint(&stack, upDown, "NAVIGATE", h)
	stackHint(&stack, b, "BACK", h)
	stackHint(&stack, leftRight, "SET", h)
}
