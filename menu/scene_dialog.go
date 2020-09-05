package menu

import (
	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
)

type sceneDialog struct {
	entry
}

func buildDialog(callbackOK func()) Scene {
	var list sceneDialog
	list.label = "Exit Dialog"
	list.callbackOK = callbackOK
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
	if input.Released[0][libretro.DeviceIDJoypadA] {
		audio.PlayEffect(audio.Effects["ok"])
		s.callbackOK()
	}

	// Cancel
	if input.Released[0][libretro.DeviceIDJoypadB] {
		audio.PlayEffect(audio.Effects["cancel"])
		menu.stack[len(menu.stack)-2].segueBack()
		menu.stack = menu.stack[:len(menu.stack)-1]
	}
}

func (s *sceneDialog) render() {
	w, h := vid.Window.GetFramebufferSize()
	fw := float32(w)
	fh := float32(h)
	vid.DrawRect(0, 0, fw, fh, 0, black.Alpha(0.85))

	var width float32 = 1000
	var height float32 = 400

	vid.DrawRect(
		fw/2-width/2*menu.ratio,
		fh/2-height/2*menu.ratio,
		width*menu.ratio,
		height*menu.ratio,
		0.05,
		white,
	)

	vid.Font.SetColor(orange)
	msg1 := "A game is currently running."
	lw1 := vid.Font.Width(0.7*menu.ratio, msg1)
	vid.Font.Printf(fw/2-lw1/2, fh/2-120*menu.ratio+20*menu.ratio, 0.7*menu.ratio, msg1)
	vid.Font.SetColor(black)
	msg2 := "If you have not saved yet, your progress will be lost."
	lw2 := vid.Font.Width(0.5*menu.ratio, msg2)
	vid.Font.Printf(fw/2-lw2/2, fh/2-30*menu.ratio+20*menu.ratio, 0.5*menu.ratio, msg2)
	msg3 := "Do you want to exit Ludo anyway?"
	lw3 := vid.Font.Width(0.5*menu.ratio, msg3)
	vid.Font.Printf(fw/2-lw3/2, fh/2+30*menu.ratio+20*menu.ratio, 0.5*menu.ratio, msg3)

	vid.Font.SetColor(darkGrey)

	var margin float32 = 15

	_, _, _, a, b, _, _, _, _, _ := hintIcons()

	vid.DrawImage(
		b,
		fw/2-width/2*menu.ratio+margin*menu.ratio,
		fh/2+height/2*menu.ratio-70*menu.ratio-margin*menu.ratio,
		70*menu.ratio, 70*menu.ratio, 1.0, darkGrey)
	vid.Font.Printf(
		fw/2-width/2*menu.ratio+margin*menu.ratio+70*menu.ratio,
		fh/2+height/2*menu.ratio-23*menu.ratio-margin*menu.ratio,
		0.4*menu.ratio,
		"NO")

	vid.DrawImage(
		a,
		fw/2+width/2*menu.ratio-150*menu.ratio-margin*menu.ratio,
		fh/2+height/2*menu.ratio-70*menu.ratio-margin*menu.ratio,
		70*menu.ratio, 70*menu.ratio, 1.0, darkGrey)
	vid.Font.Printf(
		fw/2+width/2*menu.ratio-150*menu.ratio-margin*menu.ratio+70*menu.ratio,
		fh/2+height/2*menu.ratio-23*menu.ratio-margin*menu.ratio,
		0.4*menu.ratio,
		"YES")
}

func (s *sceneDialog) drawHintBar() {
}
