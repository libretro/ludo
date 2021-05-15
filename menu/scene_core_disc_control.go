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

	list.children = append(list.children, entry{
		label: "Disc State",
		icon:  "subsetting",
		stringValue: func() string {
			if state.Core.DiskControlCallback.GetEjectState() {
				return "Ejected"
			}
			return "Inserted"
		},
	})

	if !state.Core.DiskControlCallback.GetEjectState() {
		list.children = append(list.children, entry{
			label: "Eject Disc",
			icon:  "close",
			callbackOK: func() {
				state.Core.DiskControlCallback.SetEjectState(true)
				menu.stack[len(menu.stack)-1] = buildCoreDiscControl()
				menu.tweens.FastForward()
				ntf.DisplayAndLog(ntf.Success, "Menu", "Disc Ejected.")
			},
		})

		list.children = append(list.children, entry{
			label: "Num Images",
			icon:  "subsetting",
			stringValue: func() string {
				return fmt.Sprintf("%d", state.Core.DiskControlCallback.GetNumImages())
			},
		})

		list.children = append(list.children, entry{
			label: "Image Index",
			icon:  "subsetting",
			stringValue: func() string {
				return fmt.Sprintf("%d", state.Core.DiskControlCallback.GetImageIndex())
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
