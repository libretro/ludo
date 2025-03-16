package menu

import (
	"strings"

	"github.com/libretro/ludo/core"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/state"
)

type sceneCoreOptions struct {
	entry
}

func buildCoreOptions() Scene {
	var list sceneCoreOptions
	list.label = "Core Options"

	if core.Options == nil {
		list.children = append(list.children, entry{
			label: "No options",
			icon:  "subsetting",
		})
		list.segueMount()
		return &list
	}

	for _, v := range core.Options.Vars {
		v := v
		list.children = append(list.children, entry{
			label: strings.Replace(v.Desc, "%", "%%", -1),
			icon:  "subsetting",
			stringValue: func() string {
				val := v.Choices[v.Choice]
				return strings.Replace(val, "%", "%%", -1)
			},
			incr: func(direction int) {
				v.Choice += direction
				if v.Choice < 0 {
					v.Choice = len(v.Choices) - 1
				} else if v.Choice > len(v.Choices)-1 {
					v.Choice = 0
				}
				core.Options.Updated = true
				err := core.Options.Save()
				if err != nil {
					ntf.DisplayAndLog(ntf.Error, "Core", "Error saving core options: %v", err.Error())
				}
			},
		})
	}

	list.segueMount()

	return &list
}

func (s *sceneCoreOptions) Entry() *entry {
	return &s.entry
}

func (s *sceneCoreOptions) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneCoreOptions) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneCoreOptions) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneCoreOptions) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *sceneCoreOptions) render() {
	genericRender(&s.entry)
}

func (s *sceneCoreOptions) drawHintBar() {
	w, h := menu.GetFramebufferSize()
	menu.DrawRect(0, float32(h)-88*menu.ratio, float32(w), 88*menu.ratio, 0, hintBgColor)
	menu.DrawRect(0, float32(h)-88*menu.ratio, float32(w), 2*menu.ratio, 0, sepColor)

	_, upDown, leftRight, _, b, _, _, _, _, guide := hintIcons()

	lstack := float32(75) * menu.ratio
	rstack := float32(w) - 96*menu.ratio
	stackHintLeft(&lstack, upDown, "Navigate", h)
	stackHintRight(&rstack, leftRight, "Set", h)
	stackHintRight(&rstack, b, "Back", h)
	if state.CoreRunning {
		stackHintRight(&rstack, guide, "Resume", h)
	}
}
