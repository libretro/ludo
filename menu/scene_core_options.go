package menu

import (
	"strings"

	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/state"
)

type sceneCoreOptions struct {
	entry
}

func buildCoreOptions() Scene {
	var list sceneCoreOptions
	list.label = "Core Options"

	for i, v := range core.Options.Vars {
		i := i
		v := v
		list.children = append(list.children, entry{
			label: strings.Replace(v.Desc(), "%", "%%", -1),
			icon:  "subsetting",
			stringValue: func() string {
				val := v.Choices()[core.Options.Choices[i]]
				return strings.Replace(val, "%", "%%", -1)
			},
			incr: func(direction int) {
				core.Options.Choices[i] += direction
				if core.Options.Choices[i] < 0 {
					core.Options.Choices[i] = core.Options.NumChoices(i) - 1
				} else if core.Options.Choices[i] > core.Options.NumChoices(i)-1 {
					core.Options.Choices[i] = 0
				}
				core.Options.Updated = true
				core.Options.Save()
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
	HintBar(&Props{},
		Hint(&Props{Hidden: !state.Global.CoreRunning}, "key-p", "RESUME"),
		Hint(&Props{}, "key-up-down", "NAVIGATE"),
		Hint(&Props{}, "key-z", "BACK"),
		Hint(&Props{}, "key-left-right", "SET"),
	)()
}
