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

func (hb *box) Draw(x, y float32) {
	hb.Layout()
	vid.DrawRect(x, y, hb.Width+hb.Padding*2, hb.Height+hb.Padding*2, hb.BorderRadius, hb.Color)
	var advance float32
	for _, child := range hb.Children {
		w, h := child.Size()
		switch hb.Direction {
		case Horizontal:
			child.Draw(x+advance+hb.Padding, y+hb.Padding)
			advance += w
		case Vertical:
			child.Draw(x+hb.Padding, y+advance+hb.Padding)
			advance += h
		}
	}
}

func (hb *box) Layout() (float32, float32) {
	hb.Width = 0
	for _, child := range hb.Children {
		w, h := child.Layout()
		switch hb.Direction {
		case Horizontal:
			hb.Width += w
			if h > hb.Height {
				hb.Height = h
			}
		case Vertical:
			hb.Height += h
			if w > hb.Width {
				hb.Width = w
			}
		}
	}
	return hb.Width + hb.Padding*2, hb.Height + hb.Padding*2
}

func (hb *box) Size() (float32, float32) {
	return hb.Width + hb.Padding*2, hb.Height + hb.Padding*2
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
