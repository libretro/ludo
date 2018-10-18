# gween [![](https://godoc.org/github.com/tanema/gween?status.svg)](http://godoc.org/github.com/tanema/gween)

gween (go-between) is a small library to perform [tweening](http://en.wikipedia.org/wiki/Tweening) in Go. It has a minimal
interface, and it comes with several easing functions.

# Examples

```golang
package gween

import (
  "github.com/tanema/gween/ease"
  "github.com/tanema/gween/gween"
)

// increase the value from 0 to 5 in 10 seconds
var tweenLinear = gween.New(0, 5, 10, ease.Linear)
current, isFinished := tweenLinear.Update(dt)

// make some text fall from the top of the screen, bouncing on y=300, in 4 seconds
var tweenLabel = gween.new(0, 300, 4, ease.OutBounce)
label.Y, _ = tweenLabel.Update(dt)

// fade background from white to black and foregrond from black to red in 2 seconds
currentBGColor = [4]float32{255, 255, 255, 255}
currentColor = [4]float32{0, 0, 0, 0}
var tweenBackground = gween.new(255, 0, 2, ease.Linear)
var tweenRed = gween.new(255, 0, 2, ease.Linear)
currentBG, _ := tweenBackground.Update(dt)
currentBGColor = [4]float32{currentBG, currentBG, currentBG, currentBG}
currentColor[0], _ = tweenRed.Update(dt)
```

# Interface

## Tween creation

```golang
t := gween.New(begin, end, duration, easingFunction)
```

Creates a new tween.

* `begin` is the start value
* `end` is the ending value
* `duration` means how much the change will take until it's finished. It must be a positive number.
* `easingFunction` can be either a function or a function name (see the easing section below).

This function only creates and returns the tween. It must be captured in a variable
and updated via `t.Update(dt)` in order for the changes to take place.

## Tween methods

```golang
currentValue, isFinished := t.Update(dt)
```

Gradually changes the `currentValue` toward the `end` value as time passes.

* `t` is a tween returned by `gween.New`
* `dt` is the difference in time. It will be added to the internal time counter of
  the tween. The current value at the current value will be returned using selected
  easing function.
* `currentValue` is the current eased value for the current time.
* `isFinised` is `true` if the tween has reached its limit (its *internal clock* is `>= duration`). It is false otherwise.

When the tween is complete, the `currentValue` will be equal to the `end` value.
The way they change over time will depend on the chosen easing function.

If `dt` is positive, the easing will be applied until the internal clock equals
`duration`, at which point the easing will stop. If it is negative,
the easing will play "backwards", until it reaches the initial value.


```golang
currentValue, isFinished := t.Set(clock)
```

Moves a tween's internal clock to a particular moment.

* `t` is a tween returned by `gween.New`
* `clock` is a positive number or 0. It's the new value of the tween's internal clock.
* `currentValue` is the value of the tween at the time set.
* `isFinished` works like in `t.Update`; it's `true` if the tween has reached its
  end, and `false` otherwise.

# Easing functions

Easing functions are functions that express how slow/fast the interpolation happens in tween.

Gween comes with 45 default easing functions already built-in (adapted from [Enrique García Cota's easing library](https://github.com/kikito/tween.lua)).

![tween families](https://kikito.github.io/tween.lua/img/tween-families.png)

The easing functions can be found in the `ease` package.

They can be divided into several families:

* `linear` is the simplest easing function, straight from one value to the other.
* `quad`, `cubic`, `quart`, `quint`, `expo`, `sine` and `circle` are all "smooth" curves that will make transitions look natural.
* The `back` family starts by moving the interpolation slightly "backwards" before moving it forward.
* The `bounce` family simulates the motion of an object bouncing.
* The `elastic` family simulates inertia in the easing, like an elastic gum.

Each family (except `linear`) has 4 variants:
* `In` starts slow, and accelerates at the end
* `Out` starts fast, and decelerates at the end
* `InOut` starts and ends slow, but it's fast in the middle
* `OutIn` starts and ends fast, but it's slow in the middle

| family      | in        | out        | inOut        | outIn        |
|-------------|-----------|------------|--------------|--------------|
| **Linear**  | Linear    | Linear     | Linear       | Linear       |
| **Quad**    | InQuad    | OutQuad    | InOutQuad    | OutInQuad    |
| **Cubic**   | InCubic   | OutCubic   | InOutCubic   | OutInCubic   |
| **Quart**   | InQuart   | OutQuart   | InOutQuart   | OutInQuart   |
| **Quint**   | InQuint   | OutQuint   | InOutQuint   | OutInQuint   |
| **Expo**    | InExpo    | OutExpo    | InOutExpo    | OutInExpo    |
| **Sine**    | InSine    | OutSine    | InOutSine    | OutInSine    |
| **Circ**    | InCirc    | OutCirc    | InOutCirc    | OutInCirc    |
| **Back**    | InBack    | OutBack    | InOutBack    | OutInBack    |
| **Bounce**  | InBounce  | OutBounce  | InOutBounce  | OutInBounce  |
| **Elastic** | InElastic | OutElastic | InOutElastic | OutInElastic |

## Custom easing functions

You are not limited to gween's easing functions; if you pass a function parameter
in the easing, it will be used.

The passed function will need to suite the TweenFunc interface: `func(t, b, c, d float32) float32`

* `t` (time): starts in 0 and usually moves towards duration
* `b` (begin): initial value of the of the property being eased.
* `c` (change): ending value of the property - starting value of the property
* `d` (duration): total duration of the tween

And must return the new value after the interpolation occurs.

Here's an example using a custom easing.

```golang
labelTween := tween.new(0, 300, 4, func(t, b, c, d) float32 {
  return c*t/d + b // linear ease
})
```

# Credits

The easing functions have been translated from Enrique García Cota's project in

https://github.com/kikito/tween.lua

See the LICENSE file for details.
