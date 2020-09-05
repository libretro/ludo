package menu

import (
	"github.com/libretro/ludo/video"
)

// Direction of the children widgets in a container
type Direction uint8

const (
	// Vertical means that the children will stack on top of each others
	Vertical Direction = iota
	// Horizontal means that the children will stack from left to right
	Horizontal
)

// Props are the properties of a widget
type Props struct {
	X, Y          float32
	Width, Height float32
	BorderRadius  float32
	Scale         float32
	TextAlign     string
	ContentWidth  float32
	ContentHeight float32
	Direction     Direction
	Color         video.Color
	Hidden        bool
}

type propsStack []*Props

func (s *propsStack) Push(v *Props) *propsStack {
	*s = append(*s, v)
	return s
}

func (s *propsStack) Pop() {
	l := len(*s)
	*s = (*s)[:l-1]
}

func (s *propsStack) Last() *Props {
	if len(*s) == 0 {
		return nil
	}
	return (*s)[len(*s)-1]
}

func (s *propsStack) Sum() *Props {
	sum := Props{}
	for _, t := range *s {
		sum.X += t.X
		sum.Y += t.Y
	}
	return &sum
}

var wStack = &propsStack{}
var tStack = &propsStack{}

// Box is a basic widget container
func Box(props *Props, children ...func()) func() {
	return func() {
		if props.Hidden {
			return
		}
		parent := wStack.Last()
		self := wStack.Push(props).Last()
		tsum := tStack.Push(&Props{}).Sum()
		sum := wStack.Sum()
		vid.DrawRect(sum.X+tsum.X, sum.Y+tsum.Y, props.Width, props.Height, props.BorderRadius, props.Color)
		for _, child := range children {
			child()
		}
		updateParent(parent, maxf32(props.Width, self.ContentWidth), maxf32(props.Height, self.ContentHeight))
		tStack.Pop()
		wStack.Pop()
	}
}

// HBox is an horizontal widget container
func HBox(props *Props, children ...func()) func() {
	return func() {
		if props.Hidden {
			return
		}
		props.Direction = Horizontal
		parent := wStack.Last()
		self := wStack.Push(props).Last()
		tsum := tStack.Push(&Props{}).Sum()
		sum := wStack.Sum()
		vid.DrawRect(sum.X+tsum.X, sum.Y+tsum.Y, props.Width, props.Height, props.BorderRadius, props.Color)
		for _, child := range children {
			tStack.Push(&Props{X: wStack.Last().ContentWidth})
			child()
			tStack.Pop()
		}
		updateParent(parent, maxf32(props.Width, self.ContentWidth), maxf32(props.Height, self.ContentHeight))
		tStack.Pop()
		wStack.Pop()
	}
}

// VBox is a vertical widget container
func VBox(props *Props, children ...func()) func() {
	return func() {
		if props.Hidden {
			return
		}
		props.Direction = Vertical
		parent := wStack.Last()
		self := wStack.Push(props).Last()
		tsum := tStack.Push(&Props{}).Sum()
		sum := wStack.Sum()
		vid.DrawRect(sum.X+tsum.X, sum.Y+tsum.Y, props.Width, props.Height, props.BorderRadius, props.Color)
		for _, child := range children {
			tStack.Push(&Props{Y: wStack.Last().ContentHeight})
			child()
			tStack.Pop()
		}
		updateParent(parent, maxf32(props.Width, self.ContentWidth), maxf32(props.Height, self.ContentHeight))
		tStack.Pop()
		wStack.Pop()
	}
}

// Label a widget is used to draw some text
func Label(props *Props, msg string) func() {
	return func() {
		if props.Hidden || props.Color.A == 0 {
			return
		}
		parent := wStack.Last()
		self := wStack.Push(props).Last()
		tsum := tStack.Push(&Props{}).Sum()
		sum := wStack.Sum()
		vid.Font.SetColor(props.Color)
		textWidth := vid.Font.Width(props.Scale, msg)
		switch props.TextAlign {
		case "center":
			vid.Font.Printf(sum.X+tsum.X+parent.Width/2-textWidth/2, sum.Y+tsum.Y+self.Height*0.67, props.Scale, msg)
		default:
			vid.Font.Printf(sum.X+tsum.X, sum.Y+tsum.Y+self.Height*0.67, props.Scale, msg)
		}
		updateParent(parent, textWidth, maxf32(props.Height, self.ContentHeight))
		tStack.Pop()
		wStack.Pop()
	}
}

// Image is a widget used to display an image
func Image(props *Props, image uint32) func() {
	return func() {
		if props.Hidden {
			return
		}
		parent := wStack.Last()
		sum := wStack.Push(props).Sum()
		tsum := tStack.Push(&Props{}).Sum()
		vid.DrawImage(image, sum.X+tsum.X, sum.Y+tsum.Y, props.Width, props.Height, props.Scale, props.Color)
		updateParent(parent, props.Width, props.Height)
		tStack.Pop()
		wStack.Pop()
	}
}

func updateParent(parent *Props, w, h float32) {
	if parent != nil {
		switch parent.Direction {
		case Horizontal:
			parent.ContentWidth += w
		case Vertical:
			parent.ContentHeight += h
		}
	}
}

func maxf32(a, b float32) float32 {
	if b > a {
		return b
	}
	return a
}

func minf32(a, b float32) float32 {
	if b < a {
		return b
	}
	return a
}
