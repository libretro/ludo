package menu

import (
	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/video"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type sceneKeyboard struct {
	entry
	index        int
	layout       int
	value        string
	y            float32
	alpha        float32
	callbackDone func(string)
}

var layouts = [][]string{
	{
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "0",
		"q", "w", "e", "r", "t", "y", "u", "i", "o", "p",
		"a", "s", "d", "f", "g", "h", "j", "k", "l", "@",
		"z", "x", "c", "v", "b", "n", "m", " ", "-", ".",
	},
	{
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "0",
		"Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P",
		"A", "S", "D", "F", "G", "H", "J", "K", "L", "+",
		"Z", "X", "C", "V", "B", "N", "M", " ", "_", "/",
	},
	{
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "0",
		"!", "\"", "#", "$", "%%", "&", "'", "*", "(", ")",
		"+", ",", "-", "~", "/", ":", ";", "=", "<", ">",
		"?", "@", "[", "\\", "]", "^", "_", "|", "{", "}",
	},
}

func buildKeyboard(label string, callbackDone func(string)) Scene {
	var list sceneKeyboard
	list.label = label
	list.callbackDone = callbackDone

	list.segueMount()

	return &list
}

func (s *sceneKeyboard) Entry() *entry {
	return &s.entry
}

func (s *sceneKeyboard) segueMount() {
	_, h := vid.Window.GetFramebufferSize()
	s.y = float32(h)
	s.alpha = 0
	menu.tweens[&s.y] = gween.New(s.y, 0, 0.15, ease.OutSine)
	menu.tweens[&s.alpha] = gween.New(s.alpha, 1, 0.15, ease.OutSine)
}

func (s *sceneKeyboard) segueNext() {
}

func (s *sceneKeyboard) segueBack() {
}

func (s *sceneKeyboard) update(dt float32) {
	// Right
	repeatRight(dt, input.NewState[0][libretro.DeviceIDJoypadRight], func() {
		audio.PlayEffect(audio.Effects["up"])
		if (s.index+1)%10 == 0 {
			s.index -= 9
		} else {
			s.index++
		}
	})

	// Left
	repeatLeft(dt, input.NewState[0][libretro.DeviceIDJoypadLeft], func() {
		audio.PlayEffect(audio.Effects["down"])
		if s.index%10 == 0 {
			s.index += 9
		} else {
			s.index--
		}
	})

	// Up
	repeatUp(dt, input.NewState[0][libretro.DeviceIDJoypadUp], func() {
		audio.PlayEffect(audio.Effects["up"])
		if s.index < 10 {
			s.index += len(layouts[s.layout]) - 10
		} else {
			s.index -= 10
		}
	})

	// Down
	repeatDown(dt, input.NewState[0][libretro.DeviceIDJoypadDown], func() {
		audio.PlayEffect(audio.Effects["down"])
		if s.index >= len(layouts[s.layout])-10 {
			s.index -= len(layouts[s.layout]) - 10
		} else {
			s.index += 10
		}
	})

	// OK
	if input.Released[0][libretro.DeviceIDJoypadA] {
		audio.PlayEffect(audio.Effects["ok"])
		s.value += layouts[s.layout][s.index]
	}

	// Switch layout
	if input.Released[0][libretro.DeviceIDJoypadX] {
		audio.PlayEffect(audio.Effects["ok"])
		s.layout++
		if s.layout >= len(layouts) {
			s.layout = 0
		}
	}

	// Delete character
	repeatY(dt, input.NewState[0][libretro.DeviceIDJoypadY], func() {
		if len(s.value) > 0 {
			audio.PlayEffect(audio.Effects["cancel"])
			s.value = s.value[:len(s.value)-1]
		}
	})

	// Cancel
	if input.Released[0][libretro.DeviceIDJoypadB] && len(menu.stack) > 1 {
		audio.PlayEffect(audio.Effects["cancel"])
		menu.stack[len(menu.stack)-2].segueBack()
		menu.stack = menu.stack[:len(menu.stack)-1]
	}

	// Done
	if input.Released[0][libretro.DeviceIDJoypadStart] && s.value != "" {
		audio.PlayEffect(audio.Effects["notice"])
		s.callbackDone(s.value)
		menu.stack[len(menu.stack)-2].segueBack()
		menu.stack = menu.stack[:len(menu.stack)-1]
	}
}

func (s *sceneKeyboard) render() {
	w, h := vid.Window.GetFramebufferSize()
	lines := float32(4)
	kbh := float32(h) * 0.6
	ksp := (kbh - (50 * menu.ratio)) / (lines + 1)
	ksz := ksp * 0.9
	ttw := 10 * ksp

	// Background
	vid.DrawRect(0, 0, float32(w), float32(h), 0,
		white.Alpha(s.alpha))

	// Label
	vid.Font.SetColor(black)
	vid.Font.Printf(
		float32(w)/2-ttw/2,
		s.y+float32(h)*0.15-ksz/2+ksz*0.6,
		ksz/260, s.label)

	// Value
	vid.DrawRect(float32(w)/2-ttw/2, s.y+float32(h)*0.25-ksz/2, ttw, ksz, 0,
		video.Color{R: 0.95, G: 0.95, B: 0.95, A: 1})
	vid.Font.Printf(
		float32(w)/2-ttw/2+ksz/4,
		s.y+float32(h)*0.25-ksz/2+ksz*0.62,
		ksz/200, s.value+"|")

	// Keyboard

	vid.DrawRect(0, s.y+float32(h)-kbh, float32(w), kbh, 0, black)

	vid.Font.SetColor(white)

	for i, key := range layouts[s.layout] {
		x := float32(i%10)*ksp - ttw/2 + float32(w)/2
		y := s.y + float32(i/10)*ksp + ksp/2 + float32(h) - kbh
		gw := vid.Font.Width(ksz/200, key)

		c1 := video.Color{R: 0.15, G: 0.15, B: 0.15, A: 1}
		c2 := video.Color{R: 0.25, G: 0.25, B: 0.25, A: 1}
		if i == s.index {
			c1 = video.Color{R: 0.35, G: 0.35, B: 0.35, A: 1}
			c2 = video.Color{R: 0.45, G: 0.45, B: 0.45, A: 1}
		}

		vid.DrawRect(x, y, ksz, ksz, 0.2, c1)
		vid.DrawRect(x, y, ksz, ksz*0.95, 0.2, c2)

		vid.Font.Printf(
			x+ksz/2-gw/2,
			y+ksz*0.6,
			ksz/200, key)
	}
}

func (s *sceneKeyboard) drawHintBar() {
	w, h := vid.Window.GetFramebufferSize()
	vid.DrawRect(0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 0, lightGrey)

	arrows, _, _, a, b, x, y, start, _, _ := hintIcons()

	var stack float32
	stackHint(&stack, arrows, "SELECT", h)
	stackHint(&stack, b, "BACK", h)
	stackHint(&stack, x, "SHIFT", h)
	stackHint(&stack, y, "DELETE", h)
	stackHint(&stack, a, "INSERT", h)
	stackHint(&stack, start, "DONE", h)
}
