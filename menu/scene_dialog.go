package menu

import (
	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
)

type sceneDialog struct {
	entry
	title, line1, line2 string
}

func buildYesNoDialog(title, line1, line2 string, callbackOK func()) Scene {
	var list sceneDialog
	list.label = "Confirm Dialog"
	list.callbackOK = callbackOK
	list.title = title
	list.line1 = line1
	list.line2 = line2
	audio.PlayEffect(audio.Effects["notice"])
	return &list
}

func (s *sceneDialog) Entry() *entry {
	return &s.entry
}

func (s *sceneDialog) segueMount() {
}

func (s *sceneDialog) segueNext() {
}

func (s *sceneDialog) segueBack() {
}

func (s *sceneDialog) update(dt float32) {
	// OK
	if input.Released[0][libretro.DeviceIDJoypadA] == 1 {
		audio.PlayEffect(audio.Effects["ok"])
		menu.stack[len(menu.stack)-2].segueBack()
		menu.stack = menu.stack[:len(menu.stack)-1]
		s.callbackOK()
	}

	// Cancel
	if input.Released[0][libretro.DeviceIDJoypadB] == 1 {
		audio.PlayEffect(audio.Effects["cancel"])
		menu.stack[len(menu.stack)-2].segueBack()
		menu.stack = menu.stack[:len(menu.stack)-1]
	}
}

func (s *sceneDialog) render() {
	w, h := menu.GetFramebufferSize()
	fw := float32(w)
	fh := float32(h)
	menu.DrawRect(0, 0, fw, fh, 0, black.Alpha(0.85))

	var width float32 = 1000
	var height float32 = 400

	menu.DrawRect(
		fw/2-width/2*menu.ratio,
		fh/2-height/2*menu.ratio,
		width*menu.ratio,
		height*menu.ratio,
		0.05,
		white,
	)

	menu.Font.SetColor(orange)
	lw1 := menu.Font.Width(0.7*menu.ratio, s.title)
	menu.Font.Printf(fw/2-lw1/2, fh/2-120*menu.ratio+20*menu.ratio, 0.7*menu.ratio, s.title)
	menu.Font.SetColor(black)
	lw2 := menu.Font.Width(0.5*menu.ratio, s.line1)
	menu.Font.Printf(fw/2-lw2/2, fh/2-30*menu.ratio+20*menu.ratio, 0.5*menu.ratio, s.line1)
	lw3 := menu.Font.Width(0.5*menu.ratio, s.line2)
	menu.Font.Printf(fw/2-lw3/2, fh/2+30*menu.ratio+20*menu.ratio, 0.5*menu.ratio, s.line2)

	menu.Font.SetColor(darkGrey)

	var margin float32 = 15

	_, _, _, a, b, _, _, _, _, _ := hintIcons()

	menu.DrawImage(
		b,
		fw/2-width/2*menu.ratio+margin*menu.ratio,
		fh/2+height/2*menu.ratio-70*menu.ratio-margin*menu.ratio,
		70*menu.ratio, 70*menu.ratio, 1.0, darkGrey)
	menu.Font.Printf(
		fw/2-width/2*menu.ratio+margin*menu.ratio+70*menu.ratio,
		fh/2+height/2*menu.ratio-23*menu.ratio-margin*menu.ratio,
		0.4*menu.ratio,
		"NO")

	menu.DrawImage(
		a,
		fw/2+width/2*menu.ratio-150*menu.ratio-margin*menu.ratio,
		fh/2+height/2*menu.ratio-70*menu.ratio-margin*menu.ratio,
		70*menu.ratio, 70*menu.ratio, 1.0, darkGrey)
	menu.Font.Printf(
		fw/2+width/2*menu.ratio-150*menu.ratio-margin*menu.ratio+70*menu.ratio,
		fh/2+height/2*menu.ratio-23*menu.ratio-margin*menu.ratio,
		0.4*menu.ratio,
		"YES")
}

func (s *sceneDialog) drawHintBar() {
}
