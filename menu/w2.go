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

// HBox

func mkHBox(props wProps, children ...Widget) Widget {
	return &hBox{
		Children: children,
		wProps:   props,
	}
}

type hBox struct {
	Children []Widget
	wProps
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

func mkVBox(props wProps, children ...Widget) Widget {
	return &vBox{
		Children: children,
		wProps:   props,
	}
}

type vBox struct {
	Children []Widget
	wProps
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
