package main

import (
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type menuCallback func()
type menuCallbackIncr func(int)
type menuCallbackGetValue func() string

type entry struct {
	x, y          float32
	label         string
	labelAlpha    float32
	icon          string
	iconAlpha     float32
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

	for e, t := range menu.tweens {
		var finished bool
		*e, finished = t.Update(1.0 / 60.0)
		if finished {
			delete(menu.tweens, e)
		}
	}

	currentMenu.render()
}

func initEntries(list entry) {
	for i := range list.children {
		e := &list.children[i]

		if i == list.ptr {
			e.y = 200 + float32(menu.spacing*(i-list.ptr))
			e.labelAlpha = 1.0
			e.iconAlpha = 1.0
		} else if i < list.ptr {
			e.y = 0 + float32(menu.spacing*(i-list.ptr))
			e.labelAlpha = 0.5
			e.iconAlpha = 0.5
		} else if i > list.ptr {
			e.y = 250 + float32(menu.spacing*(i-list.ptr))
			e.labelAlpha = 0.5
			e.iconAlpha = 0.5
		}
	}
}

func animateEntries() {
	currentMenu := &menu.stack[len(menu.stack)-1]

	for i := range currentMenu.children {
		e := &currentMenu.children[i]

		var y, la, ia float32
		if i == currentMenu.ptr {
			y = 200 + float32(menu.spacing*(i-currentMenu.ptr))
			la = 1.0
			ia = 1.0
		} else if i < currentMenu.ptr {
			y = 0 + float32(menu.spacing*(i-currentMenu.ptr))
			la = 0.5
			ia = 0.5
		} else if i > currentMenu.ptr {
			y = 250 + float32(menu.spacing*(i-currentMenu.ptr))
			la = 0.5
			ia = 0.5
		}

		menu.tweens[&e.y] = gween.New(e.y, y, 0.15, ease.OutSine)
		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, la, 0.15, ease.OutSine)
		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, ia, 0.15, ease.OutSine)
	}
}

func verticalRender() {
	w, h := window.GetFramebufferSize()
	currentMenu := &menu.stack[len(menu.stack)-1]

	video.font.SetColor(1, 1, 1, 1.0)
	video.font.Printf(60, 20+60, 0.5, currentMenu.label)

	for _, e := range currentMenu.children {
		if e.y < -128 || e.y > float32(h+128) {
			continue
		}

		video.font.SetColor(1.0, 1.0, 1.0, e.labelAlpha)
		video.font.Printf(110, e.y, 0.5, e.label)

		drawImage(menu.icons[e.icon], 45, int32(e.y)-44, 64, 64, color{1, 1, 1, 1})

		if e.callbackValue != nil {
			video.font.Printf(float32(w)-250, e.y, 0.5, e.callbackValue())
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
