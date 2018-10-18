// Package gween provides the Tween struct that allows an easing function to be
// animated over time. This can be used in tandem with the ease package to provide
// the easing functions.
package gween

import "github.com/tanema/gween/ease"

type (
	// Tween encapsulates the easing function along with timing data. This allows
	// a ease.TweenFunc to be used to be easily animated.
	Tween struct {
		duration float32
		time     float32
		begin    float32
		end      float32
		change   float32
		easing   ease.TweenFunc
	}
)

// New will return a new Tween when passed a begining and end value, the duration
// of the tween and the easing function to anumate between the two values. The
// easing function can be one of the provided easing functions from the ease package
// or you can provide one of your own.
func New(begin, end, duration float32, easing ease.TweenFunc) *Tween {
	return &Tween{
		begin:    begin,
		end:      end,
		change:   end - begin,
		duration: duration,
		easing:   easing,
	}
}

// Set will set the current time along the duration of the tween. It will then return
// the current value as well as a boolean to determine if the tween is finished.
func (tween *Tween) Set(time float32) (current float32, isFinished bool) {
	if time <= 0 {
		tween.time = 0
		current = tween.begin
	} else if time >= tween.duration {
		tween.time = tween.duration
		current = tween.end
	} else {
		tween.time = time
		current = tween.easing(tween.time, tween.begin, tween.change, tween.duration)
	}

	return current, tween.time >= tween.duration
}

// Reset will set the Tween to the beginning of the two values.
func (tween *Tween) Reset() {
	tween.Set(0)
}

// Update will increment the timer of the Tween and ease the value. It will then
// return the current value as well as a bool to mark if the tween is finished or not.
func (tween *Tween) Update(dt float32) (current float32, isFinished bool) {
	return tween.Set(tween.time + dt)
}
