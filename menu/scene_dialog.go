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

	width := 1000 * menu.ratio
	height := 400 * menu.ratio

	_, _, _, a, b, _, _, _, _, _ := hintIcons()

	// Background
	Box(&Props{Width: fw, Height: fh, Color: black.Alpha(0.85)},
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
				Color:     orange,
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
				Hint(&Props{}, b, "NO"),
				// The YES Hint
				Hint(&Props{X: width - 175*menu.ratio}, a, "YES"),
			),
		),
	)()
}

func (s *sceneDialog) drawHintBar() {
}
