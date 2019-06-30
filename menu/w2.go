package menu

import (
	"github.com/libretro/ludo/video"
)

type Widget interface {
	Draw(x, y float32)
	Layout() (w float32, h float32)
	Size() (w float32, h float32)
}

// HBox

type hBox struct {
	Width, Height float32
	Padding       float32
	BorderRadius  float32
	Color         video.Color
	Hidden        bool
	Children      []Widget
}

func (hb *hBox) Draw(x, y float32) {
	hb.Layout()
	vid.DrawRect(x, y, hb.Width+hb.Padding*2, hb.Height+hb.Padding*2, hb.BorderRadius, hb.Color)
	var advance float32
	for _, child := range hb.Children {
		w, _ := child.Size()
		child.Draw(x+advance+hb.Padding, y+hb.Padding)
		advance += w
	}
}

func (hb *hBox) Layout() (float32, float32) {
	hb.Width = 0
	for _, child := range hb.Children {
		w, h := child.Layout()
		hb.Width += w
		if h > hb.Height {
			hb.Height = h
		}
	}
	return hb.Width + hb.Padding*2, hb.Height + hb.Padding*2
}

func (hb *hBox) Size() (float32, float32) {
	return hb.Width + hb.Padding*2, hb.Height + hb.Padding*2
}

// VBox

type vBox struct {
	Width, Height float32
	Padding       float32
	BorderRadius  float32
	Color         video.Color
	Hidden        bool
	Children      []Widget
}

func (vb *vBox) Draw(x, y float32) {
	vb.Layout()
	vid.DrawRect(x, y, vb.Width+vb.Padding*2, vb.Height+vb.Padding*2, vb.BorderRadius, vb.Color)
	var advance float32
	for _, child := range vb.Children {
		_, h := child.Size()
		child.Draw(x+vb.Padding, y+advance+vb.Padding)
		advance += h
	}
}

func (vb *vBox) Layout() (float32, float32) {
	vb.Width = 0
	for _, child := range vb.Children {
		w, h := child.Layout()
		vb.Height += h
		if w > vb.Width {
			vb.Width = w
		}
	}
	return vb.Width + vb.Padding*2, vb.Height + vb.Padding*2
}

func (vb *vBox) Size() (float32, float32) {
	return vb.Width + vb.Padding*2, vb.Height + vb.Padding*2
}

// Label

type label struct {
	Width, Height float32
	Scale         float32
	Color         video.Color
	Hidden        bool
	Text          string
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
	Width, Height float32
	Scale         float32
	Color         video.Color
	Hidden        bool
	Image         uint32
}

func (img *image) Draw(x, y float32) {
	vid.DrawImage(img.Image, x, y, img.Width, img.Height, img.Scale, img.Color)
}

func (img *image) Layout() (float32, float32) {
	return img.Width, img.Height
}

func (img *image) Size() (float32, float32) {
	return img.Width, img.Height
}

func mkButton(icon, txt string, c video.Color) Widget {
	return &hBox{
		Color:        c,
		BorderRadius: 0.2,
		Children: []Widget{
			&image{
				Width:  70,
				Height: 70,
				Color:  video.Color{1, 1, 1, 1},
				Scale:  1,
				Image:  menu.icons[icon],
			},
			&label{
				Height: 70,
				Color:  video.Color{1, 1, 1, 1},
				Scale:  0.6 * menu.ratio,
				Text:   txt,
			},
		},
	}
}
