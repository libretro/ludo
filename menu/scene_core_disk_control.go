package menu

import (
	"fmt"

	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/state"
)

type sceneCoreDiskControl struct {
	entry
}

func buildCoreDiskControl() Scene {
	var list sceneCoreDiskControl
	list.label = "Core Disk Control"

	for i := uint(0); i < state.Core.DiskControlCallback.GetNumImages(); i++ {
		index := i
		list.children = append(list.children, entry{
			label: fmt.Sprintf("Disk %d", index+1),
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

	if len(list.children) == 0 {
		list.children = append(list.children, entry{
			label: "No disk",
			icon:  "subsetting",
		})
	}

	list.segueMount()

	return &list
}

func (s *sceneCoreDiskControl) Entry() *entry {
	return &s.entry
}

func (s *sceneCoreDiskControl) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneCoreDiskControl) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneCoreDiskControl) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneCoreDiskControl) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *sceneCoreDiskControl) render() {
	genericRender(&s.entry)
}

func (s *sceneCoreDiskControl) drawHintBar() {
	w, h := menu.GetFramebufferSize()
	menu.DrawRect(0, float32(h)-88*menu.ratio, float32(w), 88*menu.ratio, 0, hintBgColor)
	menu.DrawRect(0, float32(h)-88*menu.ratio, float32(w), 2*menu.ratio, 0, sepColor)

	_, upDown, leftRight, a, b, _, _, _, _, guide := hintIcons()

	lstack := float32(75) * menu.ratio
	rstack := float32(w) - 96*menu.ratio
	list := menu.stack[len(menu.stack)-1].Entry()
	stackHintLeft(&lstack, upDown, "Navigate", h)
	if list.children[list.ptr].callbackOK != nil {
		stackHintRight(&rstack, a, "Set", h)
	} else {
		stackHintLeft(&lstack, leftRight, "Set", h)
	}
	stackHintRight(&rstack, b, "Back", h)
	if state.CoreRunning {
		stackHintRight(&rstack, guide, "Resume", h)
	}
}
