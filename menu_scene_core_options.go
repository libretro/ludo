package main

type screenCoreOptions struct {
	entry
}

func buildCoreOptions() scene {
	var list screenCoreOptions
	list.label = "Core Options"

	for _, v := range options_vars {
		v := v
		list.children = append(list.children, entry{
			label: v.Desc(),
			icon:  "subsetting",
			stringValue: func() string {
				return v.Choices()[0]
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
