package menu

import "github.com/libretro/ludo/video"

type sceneKeyboard struct {
	entry
}

var layout = []string{
	"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", " ",
	"q", "w", "e", "r", "t", "y", "u", "i", "o", "p", " ",
	"a", "s", "d", "f", "g", "h", "j", "k", "l", "@", " ",
	"z", "x", "c", "v", "b", "n", "m", " ", "-", ".", " ",
}

func buildKeyboard() Scene {
	var list sceneKeyboard
	list.label = "Keyboard"

	list.children = append(list.children, entry{
		label: "Placeholder",
		icon:  "reload",
	})

	list.segueMount()

	return &list
}

func (s *sceneKeyboard) Entry() *entry {
	return &s.entry
}

func (s *sceneKeyboard) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneKeyboard) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneKeyboard) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneKeyboard) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *sceneKeyboard) render() {
	w, h := vid.Window.GetFramebufferSize()
	menu.ratio = float32(w) / 1920
	kbh := float32(h) * 0.65
	lines := float32(4)

	vid.DrawRect(0, float32(h)-kbh, float32(w), kbh, 1,
		video.Color{R: 0, G: 0, B: 0, A: 1})

	vid.Font.SetColor(1, 1, 1, 1)

	for i, key := range layout {
		ksp := (kbh - (50 * menu.ratio)) / (lines + 1)
		ksz := ksp * 0.9
		ttw := 11 * ksp
		x := float32(i%11)*ksp - ttw/2 + float32(w)/2
		y := float32(i/11)*ksp + ksp/2 + float32(h) - kbh
		gw := vid.Font.Width(ksz/150, key)

		vid.DrawRoundedRect(x, y, ksz, ksz, 0.2,
			video.Color{R: 0.15, G: 0.15, B: 0.15, A: 1})

		vid.DrawRoundedRect(x, y, ksz, ksz*0.95, 0.2,
			video.Color{R: 0.25, G: 0.25, B: 0.25, A: 1})

		vid.Font.Printf(
			x+ksz/2-gw/2,
			y+ksz*0.6,
			ksz/150, key)
	}
}

func (s *sceneKeyboard) drawHintBar() {
	genericDrawHintBar()
}
