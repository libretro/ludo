package menu

import (
	"github.com/libretro/ludo/state"
)

type sceneCoreDiscControl struct {
	entry
}

func buildCoreDiscControl() Scene {
	var list sceneCoreDiscControl
	list.label = "Core Disc Control"

	list.children = append(list.children, entry{
		label: "Eject State",
		icon:  "subsetting",
		stringValue: func() string {
			if state.Core.DiskControlCallback.GetEjectState() {
				return "True"
			}
			return "False"
		},
	})

	if !state.Core.DiskControlCallback.GetEjectState() {
		list.children = append(list.children, entry{
			label: "Eject Disc",
			icon:  "subsetting",
			callbackOK: func() {
				state.Core.DiskControlCallback.SetEjectState(true)
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
