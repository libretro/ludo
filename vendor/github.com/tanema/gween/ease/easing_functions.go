// Package ease provides default easing functions to be used in a Tween.
package ease

import (
	"math"
)

const backS float32 = 1.70158

var pi = float32(math.Pi)

// TweenFunc provides an interface used for the easing equation. You can use
// one of the provided easing functions or provide your own.
type TweenFunc func(t, b, c, d float32) float32

func Linear(t, b, c, d float32) float32 {
	return c*t/d + b
}

func InQuad(t, b, c, d float32) float32 {
	return c*pow(t/d, 2) + b
}

func OutQuad(t, b, c, d float32) float32 {
	t = t / d
	return -c*t*(t-2) + b
}

func InOutQuad(t, b, c, d float32) float32 {
	t = t / d * 2
	if t < 1 {
		return c/2*pow(t, 2) + b
	}
	return -c/2*((t-1)*(t-3)-1) + b
}

func OutInQuad(t, b, c, d float32) float32 {
	if t < d/2 {
		return OutQuad(t*2, b, c/2, d)
	}
	return InQuad((t*2)-d, b+c/2, c/2, d)
}

func InCubic(t, b, c, d float32) float32 {
	return c*pow(t/d, 3) + b
}

func OutCubic(t, b, c, d float32) float32 {
	return c*(pow(t/d-1, 3)+1) + b
}

func InOutCubic(t, b, c, d float32) float32 {
	t = t / d * 2
	if t < 1 {
		return c/2*t*t*t + b
	}
	t = t - 2
	return c/2*(t*t*t+2) + b
}

func OutInCubic(t, b, c, d float32) float32 {
	if t < d/2 {
		return OutCubic(t*2, b, c/2, d)
	}
	return InCubic((t*2)-d, b+c/2, c/2, d)
}

func InQuart(t, b, c, d float32) float32 {
	return c*pow(t/d, 4) + b
}

func OutQuart(t, b, c, d float32) float32 {
	return -c*(pow(t/d-1, 4)-1) + b
}

func InOutQuart(t, b, c, d float32) float32 {
	t = t / d * 2
	if t < 1 {
		return c/2*pow(t, 4) + b
	}
	return -c/2*(pow(t-2, 4)-2) + b
}

func OutInQuart(t, b, c, d float32) float32 {
	if t < d/2 {
		return OutQuart(t*2, b, c/2, d)
	}
	return InQuart((t*2)-d, b+c/2, c/2, d)
}

func InQuint(t, b, c, d float32) float32 {
	return c*pow(t/d, 5) + b
}

func OutQuint(t, b, c, d float32) float32 {
	return c*(pow(t/d-1, 5)+1) + b
}

func InOutQuint(t, b, c, d float32) float32 {
	t = t / d * 2
	if t < 1 {
		return c/2*pow(t, 5) + b
	}
	return c/2*(pow(t-2, 5)+2) + b
}

func OutInQuint(t, b, c, d float32) float32 {
	if t < d/2 {
		return OutQuint(t*2, b, c/2, d)
	}
	return InQuint((t*2)-d, b+c/2, c/2, d)
}

func InSine(t, b, c, d float32) float32 {
	return -c*cos(t/d*(pi/2)) + c + b
}

func OutSine(t, b, c, d float32) float32 {
	return c*sin(t/d*(pi/2)) + b
}

func InOutSine(t, b, c, d float32) float32 {
	return -c/2*(cos(pi*t/d)-1) + b
}

func OutInSine(t, b, c, d float32) float32 {
	if t < d/2 {
		return OutSine(t*2, b, c/2, d)
	}
	return InSine((t*2)-d, b+c/2, c/2, d)
}

func InExpo(t, b, c, d float32) float32 {
	if t == 0 {
		return b
	}
	return c*pow(2, 10*(t/d-1)) + b - c*0.001
}

func OutExpo(t, b, c, d float32) float32 {
	if t == d {
		return b + c
	}
	return c*1.001*(-pow(2, -10*t/d)+1) + b
}

func InOutExpo(t, b, c, d float32) float32 {
	if t == 0 {
		return b
	}
	if t == d {
		return b + c
	}
	t = t / d * 2
	if t < 1 {
		return c/2*pow(2, 10*(t-1)) + b - c*0.0005
	}
	return c/2*1.0005*(-pow(2, -10*(t-1))+2) + b
}

func OutInExpo(t, b, c, d float32) float32 {
	if t < d/2 {
		return OutExpo(t*2, b, c/2, d)
	}
	return InExpo((t*2)-d, b+c/2, c/2, d)
}

func InCirc(t, b, c, d float32) float32 {
	return (-c*(sqrt(1-pow(t/d, 2))-1) + b)
}

func OutCirc(t, b, c, d float32) float32 {
	return (c*sqrt(1-pow(t/d-1, 2)) + b)
}

func InOutCirc(t, b, c, d float32) float32 {
	t = t / d * 2
	if t < 1 {
		return -c/2*(sqrt(1-t*t)-1) + b
	}
	t = t - 2
	return c/2*(sqrt(1-t*t)+1) + b
}

func OutInCirc(t, b, c, d float32) float32 {
	if t < d/2 {
		return OutCirc(t*2, b, c/2, d)
	}
	return InCirc((t*2)-d, b+c/2, c/2, d)
}

func InElastic(t, b, c, d float32) float32 {
	if t == 0 {
		return b
	}
	t = t / d
	if t == 1 {
		return b + c
	}
	p, a, s := calculatePAS(c, d)
	t = t - 1
	return -(a * pow(2, 10*t) * sin((t*d-s)*(2*pi)/p)) + b
}

func OutElastic(t, b, c, d float32) float32 {
	if t == 0 {
		return b
	}
	t = t / d
	if t == 1 {
		return b + c
	}
	p, a, s := calculatePAS(c, d)
	return a*pow(2, -10*t)*sin((t*d-s)*(2*pi)/p) + c + b
}

func InOutElastic(t, b, c, d float32) float32 {
	if t == 0 {
		return b
	}
	t = t / d * 2
	if t == 2 {
		return b + c
	}
	p, a, s := calculatePAS(c, d)
	t = t - 1
	if t < 0 {
		return -0.5*(a*pow(2, 10*t)*sin((t*d-s)*(2*pi)/p)) + b
	}
	return a*pow(2, -10*t)*sin((t*d-s)*(2*pi)/p)*0.5 + c + b
}

func OutInElastic(t, b, c, d float32) float32 {
	if t < d/2 {
		return OutElastic(t*2, b, c/2, d)
	}
	return InElastic((t*2)-d, b+c/2, c/2, d)
}

func InBack(t, b, c, d float32) float32 {
	t = t / d
	return c*t*t*((backS+1)*t-backS) + b
}

func OutBack(t, b, c, d float32) float32 {
	t = t/d - 1
	return c*(t*t*((backS+1)*t+backS)+1) + b
}

func InOutBack(t, b, c, d float32) float32 {
	s := backS * 1.525
	t = t / d * 2
	if t < 1 {
		return c/2*(t*t*((s+1)*t-s)) + b
	}
	t = t - 2
	return c/2*(t*t*((s+1)*t+s)+2) + b
}

func OutInBack(t, b, c, d float32) float32 {
	if t < (d / 2) {
		return OutBack(t*2, b, c/2, d)
	}
	return InBack((t*2)-d, b+c/2, c/2, d)
}

func OutBounce(t, b, c, d float32) float32 {
	t = t / d
	if t < 1/2.75 {
		return c*(7.5625*t*t) + b
	}
	if t < 2/2.75 {
		t = t - (1.5 / 2.75)
		return c*(7.5625*t*t+0.75) + b
	} else if t < 2.5/2.75 {
		t = t - (2.25 / 2.75)
		return c*(7.5625*t*t+0.9375) + b
	}
	t = t - (2.625 / 2.75)
	return c*(7.5625*t*t+0.984375) + b
}

func InBounce(t, b, c, d float32) float32 {
	return c - OutBounce(d-t, 0, c, d) + b
}

func InOutBounce(t, b, c, d float32) float32 {
	if t < d/2 {
		return InBounce(t*2, 0, c, d)*0.5 + b
	}
	return OutBounce(t*2-d, 0, c, d)*0.5 + c*.5 + b
}

func OutInBounce(t, b, c, d float32) float32 {
	if t < d/2 {
		return OutBounce(t*2, b, c/2, d)
	}
	return InBounce((t*2)-d, b+c/2, c/2, d)
}

func calculatePAS(c, d float32) (p, a, s float32) {
	p = d * 0.3
	if a < abs(c) {
		return p, c, p / 4
	}
	return p, a, p / (2 * pi) * asin(c/a)
}

func abs(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

func pow(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}

func cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func asin(x float32) float32 {
	return float32(math.Asin(float64(x)))
}

func sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}
