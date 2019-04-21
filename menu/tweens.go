package menu

import "github.com/tanema/gween"

// Tweens are the current animations of the menu components
type Tweens map[*float32]*gween.Tween

// Update loops over the animation queue and updade them so we can see progress
func (ts Tweens) Update(dt float32) {
	for e, t := range ts {
		var finished bool
		*e, finished = t.Update(dt)
		if finished {
			delete(ts, e)
		}
	}
}

// FastForward finishes all the current animations in the queue.
func (ts Tweens) FastForward() {
	ts.Update(10)
}
