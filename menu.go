package main

import (
	"github.com/tanema/gween"
)

type menuCallback func()
type menuCallbackIncr func(int)
type menuCallbackGetValue func() string

type entry struct {
	x             float32
	label         string
	labelAlpha    float32
	icon          string
	iconAlpha     float32
	scroll        float32
	scrollTween   *gween.Tween
	ptr           int
	callback      menuCallback
	callbackValue menuCallbackGetValue
	callbackIncr  menuCallbackIncr
	children      []entry
	input         func()
	render        func()
}

var menu struct {
	stack         []entry
	icons         map[string]uint32
	inputCooldown int
	spacing       int
	tweens        map[*float32]*gween.Tween
}

func menuRender() {
	fullscreenViewport()

	currentMenu := &menu.stack[len(menu.stack)-1]
	if currentMenu.scrollTween != nil {
		currentMenu.scroll, _ = currentMenu.scrollTween.Update(1.0 / 60.0)
	}

	for e, t := range menu.tweens {
		var finished bool
		*e, finished = t.Update(1.0 / 60.0)
		if finished {
			delete(menu.tweens, e)
		}
	}

	currentMenu.render()
}

func verticalRender() {
	w, h := window.GetFramebufferSize()
	currentMenu := &menu.stack[len(menu.stack)-1]

	video.font.SetColor(1, 1, 1, 1.0)
	video.font.Printf(60, 20+60, 0.5, currentMenu.label)

	for i, e := range currentMenu.children {
		y := -currentMenu.scroll + 20 + float32(menu.spacing*(i+2))

		if y < 0 || y > float32(h) {
			continue
		}

		if i == currentMenu.ptr {
			video.font.SetColor(0.0, 1.0, 0.0, 1.0)
		} else {
			video.font.SetColor(0.6, 0.6, 0.9, 1.0)
		}
		video.font.Printf(110, y, 0.5, e.label)

		drawImage(menu.icons[e.icon], 45, int32(y)-44, 64, 64, color{1, 1, 1, 1})

		if e.callbackValue != nil {
			video.font.Printf(float32(w)-250, y, 0.5, e.callbackValue())
		}
	}
}

func contextReset() {
	menu.spacing = 70

	menu.icons = map[string]uint32{
		"file":       newImage("assets/file.png"),
		"folder":     newImage("assets/folder.png"),
		"subsetting": newImage("assets/subsetting.png"),
		"setting":    newImage("assets/setting.png"),
		"resume":     newImage("assets/resume.png"),
		"reset":      newImage("assets/reset.png"),
		"loadstate":  newImage("assets/loadstate.png"),
		"savestate":  newImage("assets/savestate.png"),
		"screenshot": newImage("assets/screenshot.png"),
		"Nintendo - Super Nintendo Entertainment System": newImage("assets/Nintendo - Super Nintendo Entertainment System.png"),
		"Sega - Mega Drive - Genesis":                    newImage("assets/Sega - Mega Drive - Genesis.png"),
	}
}

func menuInit() {
	menu.tweens = make(map[*float32]*gween.Tween)

	if g.coreRunning {
		menu.stack = append(menu.stack, buildTabs())
		menu.stack = append(menu.stack, buildMainMenu())
		menu.stack = append(menu.stack, buildQuickMenu())
	} else {
		menu.stack = append(menu.stack, buildTabs())
	}
}
