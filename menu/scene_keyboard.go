package menu

import (
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/video"
)

type sceneKeyboard struct {
	entry
	index int
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
}

func (s *sceneKeyboard) segueNext() {
}

func (s *sceneKeyboard) segueBack() {
}

func (s *sceneKeyboard) update(dt float32) {
	menu.inputCooldown -= dt
	if menu.inputCooldown < 0 {
		menu.inputCooldown = 0
	}

	// Right
	if input.Released[0][libretro.DeviceIDJoypadRight] {
		if (s.index+1)%11 == 0 {
			s.index -= 10
		} else {
			s.index++
		}
	}

	// Left
	if input.Released[0][libretro.DeviceIDJoypadLeft] {
		if s.index%11 == 0 {
			s.index += 10
		} else {
			s.index--
		}
	}

	// Up
	if input.Released[0][libretro.DeviceIDJoypadUp] {
		if s.index < 11 {
			s.index += len(layout) - 11
		} else {
			s.index -= 11
		}
	}

	// Down
	if input.Released[0][libretro.DeviceIDJoypadDown] {
		if s.index > len(layout)-11 {
			s.index -= len(layout) - 11
		} else {
			s.index += 11
		}
	}

	// Cancel
	if input.Released[0][libretro.DeviceIDJoypadB] {
		if len(menu.stack) > 1 {
			menu.stack = menu.stack[:len(menu.stack)-1]
		}
	}
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

		c1 := video.Color{R: 0.15, G: 0.15, B: 0.15, A: 1}
		c2 := video.Color{R: 0.25, G: 0.25, B: 0.25, A: 1}
		if i == s.index {
			c1 = video.Color{R: 0.35, G: 0.35, B: 0.35, A: 1}
			c2 = video.Color{R: 0.45, G: 0.45, B: 0.45, A: 1}
		}

		vid.DrawRoundedRect(x, y, ksz, ksz, 0.2, c1)
		vid.DrawRoundedRect(x, y, ksz, ksz*0.95, 0.2, c2)

		vid.Font.Printf(
			x+ksz/2-gw/2,
			y+ksz*0.6,
			ksz/150, key)
	}
}

func (s *sceneKeyboard) drawHintBar() {
	genericDrawHintBar()
}
