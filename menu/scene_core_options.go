package menu

import (
	"strings"
)

type screenCoreOptions struct {
	entry
}

func buildCoreOptions() Scene {
	var list screenCoreOptions
	list.label = "Core Options"

	for i, v := range opts.Vars {
		i := i
		v := v
		list.children = append(list.children, entry{
			label: strings.Replace(v.Desc(), "%", "%%", -1),
			icon:  "subsetting",
			stringValue: func() string {
				val := v.Choices()[opts.Choices[i]]
				return strings.Replace(val, "%", "%%", -1)
			},
			incr: func(direction int) {
				opts.Choices[i] += direction
				if opts.Choices[i] < 0 {
					opts.Choices[i] = opts.NumChoices(i) - 1
				} else if opts.Choices[i] > opts.NumChoices(i)-1 {
					opts.Choices[i] = 0
				}
				opts.Updated = true
				opts.Save()
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
