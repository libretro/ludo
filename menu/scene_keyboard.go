package menu

import (
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/video"
)

type sceneKeyboard struct {
	entry
	index int
	value string
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
	if input.NewState[0][libretro.DeviceIDJoypadRight] {
		if menu.inputCooldown == 0 {
			if (s.index+1)%11 == 0 {
				s.index -= 10
			} else {
				s.index++
			}
			menu.inputCooldown = 0.1
		}
	}

	// Left
	if input.NewState[0][libretro.DeviceIDJoypadLeft] {
		if menu.inputCooldown == 0 {
			if s.index%11 == 0 {
				s.index += 10
			} else {
				s.index--
			}
			menu.inputCooldown = 0.1
		}
	}

	// Up
	if input.NewState[0][libretro.DeviceIDJoypadUp] {
		if menu.inputCooldown == 0 {
			if s.index < 11 {
				s.index += len(layout) - 11
			} else {
				s.index -= 11
			}
			menu.inputCooldown = 0.1
		}
	}

	// Down
	if input.NewState[0][libretro.DeviceIDJoypadDown] {
		if menu.inputCooldown == 0 {
			if s.index >= len(layout)-11 {
				s.index -= len(layout) - 11
			} else {
				s.index += 11
			}
			menu.inputCooldown = 0.1
		}
	}

	// OK
	if input.Released[0][libretro.DeviceIDJoypadA] {
		s.value += layout[s.index]
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
	lines := float32(4)
	kbh := float32(h) * 0.65
	ksp := (kbh - (50 * menu.ratio)) / (lines + 1)
	ksz := ksp * 0.9
	ttw := 11 * ksp

	// Background
	vid.DrawRect(0, 0, float32(w), float32(h), 1,
		video.Color{R: 1, G: 1, B: 1, A: 1})

	// Value

	vid.DrawRect(float32(w)/2-ttw/2, float32(h)*0.175-ksp/2, ttw, ksp, 1,
		video.Color{R: 0.95, G: 0.95, B: 0.95, A: 1})
	vid.Font.SetColor(0, 0, 0, 1)
	vid.Font.Printf(
		float32(w)/2-ttw/2+ksp/2,
		float32(h)*0.175-ksp/2+ksp*0.6,
		ksz/150, s.value+"|")

	// Keyboard

	vid.DrawRect(0, float32(h)-kbh, float32(w), kbh, 1,
		video.Color{R: 0, G: 0, B: 0, A: 1})

	vid.Font.SetColor(1, 1, 1, 1)

	for i, key := range layout {
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
