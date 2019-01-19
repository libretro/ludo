package menu

import (
	"strings"

	"github.com/libretro/ludo/core"
)

type screenCoreOptions struct {
	entry
}

func buildCoreOptions() Scene {
	var list screenCoreOptions
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

func (s *screenCoreOptions) Entry() *entry {
	return &s.entry
}

func (s *screenCoreOptions) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *screenCoreOptions) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *screenCoreOptions) segueBack() {
	genericAnimate(&s.entry)
}

func (s *screenCoreOptions) update() {
	genericInput(&s.entry)
}

func (s *screenCoreOptions) render() {
	genericRender(&s.entry)
}

func (s *screenCoreOptions) drawHintBar() {
	genericDrawHintBar()
}
