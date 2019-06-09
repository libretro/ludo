package menu

import (
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/video"
)

type sceneDialog struct {
	entry
}

func buildDialog(callbackOK func()) Scene {
	var list sceneDialog
	list.label = "Exit Dialog"
	list.callbackOK = callbackOK
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
		s.callbackOK()
	}

	// Cancel
	if input.Released[0][libretro.DeviceIDJoypadB] {
		menu.stack[len(menu.stack)-2].segueBack()
		menu.stack = menu.stack[:len(menu.stack)-1]
	}
}

func (s *sceneDialog) render() {
	w, h := vid.Window.GetFramebufferSize()
	vid.DrawRect(0, 0, float32(w), float32(h), 0, video.Color{R: 0.5, G: 0, B: 0, A: 1})
	vid.Font.SetColor(1, 1, 1, 1)
	msg1 := "A game is currently running."
	lw1 := vid.Font.Width(0.6*menu.ratio, msg1)
	vid.Font.Printf(float32(w)/2-lw1/2, float32(h)/2-60*menu.ratio, 0.6*menu.ratio, msg1)
	msg2 := "If you have not saved yet, your progress will be lost."
	lw2 := vid.Font.Width(0.6*menu.ratio, msg2)
	vid.Font.Printf(float32(w)/2-lw2/2, float32(h)/2, 0.6*menu.ratio, msg2)
	msg3 := "Do you want to exit Ludo anyway?"
	lw3 := vid.Font.Width(0.6*menu.ratio, msg3)
	vid.Font.Printf(float32(w)/2-lw3/2, float32(h)/2+60*menu.ratio, 0.6*menu.ratio, msg3)
}

func (s *sceneDialog) drawHintBar() {
	w, h := vid.Window.GetFramebufferSize()
	vid.DrawRect(0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 0, video.Color{R: 0.75, G: 0.75, B: 0.75, A: 1})

	var stack float32
	stackHint(&stack, "key-z", "NO", h)
	stackHint(&stack, "key-x", "YES", h)
}
