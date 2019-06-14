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
	fw := float32(w)
	fh := float32(h)

	width := 1000 * menu.ratio
	height := 400 * menu.ratio

	white := video.Color{R: 1, G: 1, B: 1, A: 1}
	black := video.Color{R: 0, G: 0, B: 0, A: 1}
	warningTitle := video.Color{R: 0.8, G: 0.4, B: 0.1, A: 1}

	// Background
	Box(&Props{Width: fw, Height: fh, Color: video.Color{R: 0, G: 0, B: 0, A: 0.85}},
		// Dialog
		VBox(&Props{
			X:            fw/2 - width/2,
			Y:            fh/2 - height/2,
			Width:        width,
			Height:       height,
			BorderRadius: 0.05,
			Color:        white,
		},
			// Title
			Label(&Props{
				TextAlign: "center",
				Scale:     0.7 * menu.ratio,
				Color:     warningTitle,
				Height:    150 * menu.ratio,
			}, "A game is currently running"),
			// Messages
			Label(&Props{
				TextAlign: "center",
				Scale:     0.5 * menu.ratio,
				Color:     black,
				Height:    60 * menu.ratio,
			}, "If you have not saved yet, your progress will be lost."),
			Label(&Props{
				TextAlign: "center",
				Scale:     0.5 * menu.ratio,
				Color:     black,
				Height:    60 * menu.ratio,
			}, "Do you want to exit Ludo anyway?"),
			Box(&Props{Height: 40 * menu.ratio}),
			Box(&Props{},
				// The NO Hint
				Hint(&Props{}, "key-z", "NO"),
				// The YES Hint
				Hint(&Props{X: width - 175*menu.ratio}, "key-x", "YES"),
			),
		),
	)()
}

func (s *sceneDialog) drawHintBar() {
}
