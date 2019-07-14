package menu

import (
	"github.com/libretro/ludo/video"
)

type Widget interface {
	Draw(x, y float32)
	Layout() (w float32, h float32)
}

const (
	top int = iota
	bottom
	left
	right
)

type dirs [4]float32

type wProps struct {
	Margin        dirs
	Padding       dirs
	BorderRadius  float32
	Width, Height float32
	Scale         float32
	Color         video.Color
	Hidden        bool
}

// Dirs returns a padding or margin array, allowing specifying the 4 directions
// as {top, bottom, left, right}, but also taking shortcuts like
// {vertical, horizontal} and {allfour}
func Dirs(ds ...float32) dirs {
	if len(ds) >= 4 {
		return dirs{ds[0], ds[1], ds[2], ds[3]}
	} else if len(ds) >= 2 {
		return dirs{ds[0], ds[0], ds[1], ds[1]}
	}
	return dirs{ds[0], ds[0], ds[0], ds[0]}
}

// Box

type box struct {
	Children  []Widget
	Direction Direction
	wProps
}

func mkHBox(props wProps, children ...Widget) Widget {
	return &box{
		Children:  children,
		Direction: Horizontal,
		wProps:    props,
	}
}

func mkVBox(props wProps, children ...Widget) Widget {
	return &box{
		Children:  children,
		Direction: Vertical,
		wProps:    props,
	}
}

func (b *box) Draw(x, y float32) {
	b.Layout()
	vid.DrawRect(
		x+b.Margin[left],
		y+b.Margin[top],
		b.Width+b.Padding[left]+b.Padding[right],
		b.Height+b.Padding[top]+b.Padding[bottom],
		b.BorderRadius,
		b.Color,
	)
	var advance float32
	for _, child := range b.Children {
		w, h := child.Layout()
		switch b.Direction {
		case Horizontal:
			child.Draw(x+advance+b.Padding[left]+b.Margin[left], y+b.Padding[top]+b.Margin[top])
			advance += w
		case Vertical:
			child.Draw(x+b.Padding[left]+b.Margin[left], y+advance+b.Padding[top]+b.Margin[top])
			advance += h
		}
	}
}

func (b *box) Layout() (float32, float32) {
	b.Width = 0
	for _, child := range b.Children {
		w, h := child.Layout()
		switch b.Direction {
		case Horizontal:
			b.Width += w
			if h > b.Height {
				b.Height = h
			}
		case Vertical:
			b.Height += h
			if w > b.Width {
				b.Width = w
			}
		}
	}
	return b.Width + b.Padding[left] + b.Padding[right] + b.Margin[left] + b.Margin[right],
		b.Height + b.Padding[top] + b.Padding[bottom] + b.Margin[top] + b.Margin[bottom]
}

// Label

type label struct {
	Text string
	wProps
}

func mkLabel(props wProps, text string) Widget {
	return &label{
		Text:   text,
		wProps: props,
	}
}

func (lb *label) Draw(x, y float32) {
	lb.Layout()
	vid.Font.SetColor(lb.Color.R, lb.Color.G, lb.Color.B, lb.Color.A)
	vid.Font.Printf(
		x+lb.Padding[left]+lb.Margin[left],
		y+lb.Padding[top]+lb.Margin[top]+lb.Height*0.67,
		lb.Scale,
		lb.Text,
	)
}

func (lb *label) Layout() (float32, float32) {
	lb.Width = vid.Font.Width(lb.Scale, lb.Text)
	return lb.Width + lb.Padding[left] + lb.Padding[right] + lb.Margin[left] + lb.Margin[right],
		lb.Height + lb.Padding[top] + lb.Padding[bottom] + lb.Margin[top] + lb.Margin[bottom]
}

// Image

type image struct {
	Texture uint32
	wProps
}

func mkImage(props wProps, texture uint32) Widget {
	return &image{
		Texture: texture,
		wProps:  props,
	}
}

func (img *image) Draw(x, y float32) {
	vid.DrawImage(
		img.Texture,
		x+img.Padding[left]+img.Margin[left],
		y+img.Padding[left]+img.Margin[left],
		img.Width,
		img.Height,
		img.Scale,
		img.Color,
	)
}

func (img *image) Layout() (float32, float32) {
	return img.Width + img.Padding[left] + img.Padding[right] + img.Margin[left] + img.Margin[right],
		img.Height + img.Padding[top] + img.Padding[bottom] + img.Margin[top] + img.Margin[bottom]
}

func mkButton(props wProps, icon, txt string) Widget {
	return mkHBox(props,
		mkImage(wProps{
			Width:  70,
			Height: 70,
			Color:  video.Color{1, 1, 1, 1},
			Scale:  1,
		}, menu.icons[icon]),
		mkLabel(wProps{
			Height: 70,
			Color:  video.Color{1, 1, 1, 1},
			Scale:  0.6 * menu.ratio,
		}, txt),
	)
}
