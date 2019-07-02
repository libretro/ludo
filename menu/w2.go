package menu

import (
	"github.com/libretro/ludo/video"
)

type Widget interface {
	Draw(x, y float32)
	Layout() (w float32, h float32)
	Size() (w float32, h float32)
}

type wProps struct {
	Padding       float32
	BorderRadius  float32
	Width, Height float32
	Scale         float32
	Color         video.Color
	Hidden        bool
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
	vid.DrawRect(x, y, b.Width+b.Padding*2, b.Height+b.Padding*2, b.BorderRadius, b.Color)
	var advance float32
	for _, child := range b.Children {
		w, h := child.Size()
		switch b.Direction {
		case Horizontal:
			child.Draw(x+advance+b.Padding, y+b.Padding)
			advance += w
		case Vertical:
			child.Draw(x+b.Padding, y+advance+b.Padding)
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
	return b.Width + b.Padding*2, b.Height + b.Padding*2
}

func (b *box) Size() (float32, float32) {
	return b.Width + b.Padding*2, b.Height + b.Padding*2
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
	vid.Font.Printf(x, y+lb.Height*0.67, lb.Scale, lb.Text)
}

func (lb *label) Layout() (float32, float32) {
	lb.Width = vid.Font.Width(lb.Scale, lb.Text)
	return lb.Width, lb.Height
}

func (lb *label) Size() (float32, float32) {
	return lb.Width, lb.Height
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
	vid.DrawImage(img.Texture, x, y, img.Width, img.Height, img.Scale, img.Color)
}

func (img *image) Layout() (float32, float32) {
	return img.Width, img.Height
}

func (img *image) Size() (float32, float32) {
	return img.Width, img.Height
}

func mkButton(icon, txt string, c video.Color) Widget {
	return mkHBox(wProps{
		Color:        c,
		BorderRadius: 0.2,
	},
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
