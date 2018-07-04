package main

import (
	"strings"
)

type screenCoreOptions struct {
	entry
}

func buildCoreOptions() scene {
	var list screenCoreOptions
	list.label = "Core Options"

	for i, v := range g.options.Vars {
		i := i
		v := v
		list.children = append(list.children, entry{
			label: strings.Replace(v.Desc(), "%", "%%", -1),
			icon:  "subsetting",
			stringValue: func() string {
				val := v.Choices()[g.options.Choices[i]]
				return strings.Replace(val, "%", "%%", -1)
			},
			incr: func(direction int) {
				g.options.Choices[i] += direction
				if g.options.Choices[i] < 0 {
					g.options.Choices[i] = g.options.numChoices(i) - 1
				} else if g.options.Choices[i] > g.options.numChoices(i)-1 {
					g.options.Choices[i] = 0
				}
				g.options.Updated = true
				g.options.save()
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
