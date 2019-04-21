package menu

import "github.com/tanema/gween"

type tweens map[*float32]*gween.Tween

// updateTweens loops over the animation queue and updade them so we can see progress
func updateTweens(dt float32) {
	for e, t := range menu.tweens {
		var finished bool
		*e, finished = t.Update(dt)
		if finished {
			delete(menu.tweens, e)
		}
	}
}

// fastForwardTweens finishes all the current animations in the queue.
func fastForwardTweens() {
	updateTweens(10)
}
