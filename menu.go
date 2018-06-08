package main

import (
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type menuCallback func()
type menuCallbackIncr func(int)
type menuCallbackGetValue func() string

type entry struct {
	x, y, scale     float32
	width           float32
	label, subLabel string
	labelAlpha      float32
	icon            string
	iconAlpha       float32
	ptr             int
	callback        menuCallback
	callbackValue   menuCallbackGetValue
	callbackIncr    menuCallbackIncr
	children        []entry
}

type screen interface {
	update()
	render()
}

var menu struct {
	stack         []screen
	icons         map[string]uint32
	inputCooldown int
	spacing       int
	tweens        map[*float32]*gween.Tween
	scroll        float32
}

func menuRender() {
	fullscreenViewport()

	for e, t := range menu.tweens {
		var finished bool
		*e, finished = t.Update(1.0 / 60.0)
		if finished {
			delete(menu.tweens, e)
		}
	}

	cur := len(menu.stack) - 1
	for i := cur - 1; i <= cur+1; i++ {
		if i < 0 || i > cur {
			continue
		}

		menu := menu.stack[i]
		menu.render()
	}
}

func initEntries(list entry) {
	for i := range list.children {
		e := &list.children[i]

		if i == list.ptr {
			e.y = 200 + float32(menu.spacing*(i-list.ptr))
			e.labelAlpha = 1.0
			e.iconAlpha = 1.0
			e.scale = 1.0
		} else if i < list.ptr {
			e.y = 0 + float32(menu.spacing*(i-list.ptr))
			e.labelAlpha = 0.5
			e.iconAlpha = 0.5
			e.scale = 0.5
		} else if i > list.ptr {
			e.y = 250 + float32(menu.spacing*(i-list.ptr))
			e.labelAlpha = 0.5
			e.iconAlpha = 0.5
			e.scale = 0.5
		}
	}
}

func animateEntries(list *entry) {
	for i := range list.children {
		e := &list.children[i]

		var y, la, a, s float32
		if i == list.ptr {
			y = 200 + float32(menu.spacing*(i-list.ptr))
			la = 1.0
			a = 1.0
			s = 1.0
		} else if i < list.ptr {
			y = 0 + float32(menu.spacing*(i-list.ptr))
			la = 0.5
			a = 0.5
			s = 0.5
		} else if i > list.ptr {
			y = 250 + float32(menu.spacing*(i-list.ptr))
			la = 0.5
			a = 0.5
			s = 0.5
		}

		menu.tweens[&e.y] = gween.New(e.y, y, 0.15, ease.OutSine)
		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, la, 0.15, ease.OutSine)
		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, a, 0.15, ease.OutSine)
		menu.tweens[&e.scale] = gween.New(e.scale, s, 0.15, ease.OutSine)
	}
}

func verticalRender(list *entry) {
	w, h := window.GetFramebufferSize()

	video.font.SetColor(1, 1, 1, 1.0)
	video.font.Printf(60, 20+60, 0.5, list.label)

	for _, e := range list.children {
		if e.y < -128 || e.y > float32(h+128) {
			continue
		}

		drawImage(menu.icons[e.icon], 88-64*e.scale, e.y-12-64*e.scale, 128, 128, e.scale, color{1, 1, 1, e.iconAlpha})
		video.font.SetColor(1.0, 1.0, 1.0, e.labelAlpha)
		video.font.Printf(150, e.y, 0.5, e.label)

		if e.callbackValue != nil {
			lw := video.font.Width(0.5, e.callbackValue())
			video.font.Printf(float32(w)-lw-70, e.y, 0.5, e.callbackValue())
		}
	}
}

func contextReset() {
	menu.spacing = 70
	video.white = newWhite()

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
