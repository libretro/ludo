package menu

import (
	"strings"

	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"
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
	w, h := vid.Window.GetFramebufferSize()
	vid.DrawRect(0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 0, video.Color{R: 0.75, G: 0.75, B: 0.75, A: 1})

	var stack float32
	if state.Global.CoreRunning {
		stackHint(&stack, "key-p", "RESUME", h)
	}
	stackHint(&stack, "key-up-down", "NAVIGATE", h)
	stackHint(&stack, "key-z", "BACK", h)
	stackHint(&stack, "key-left-right", "SET", h)
}
