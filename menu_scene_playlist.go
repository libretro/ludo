package main

type screenPlaylist struct {
	entry
}

func buildPlaylist(path string) scene {
	var list screenPlaylist
	list.label = filename(path)

	list.segueMount()

	return &list
}

// Generic stuff

func (s *screenPlaylist) Entry() *entry {
	return &s.entry
}

func (s *screenPlaylist) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *screenPlaylist) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *screenPlaylist) segueBack() {
	genericAnimate(&s.entry)
}

func (s *screenPlaylist) update() {
	genericInput(&s.entry)
}

func (s *screenPlaylist) render() {
	genericRender(&s.entry)
}
