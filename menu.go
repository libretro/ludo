package main

import (
	"math"

	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type menuCallback func()
type menuCallbackIncr func(int)
type menuCallbackGetValue func() string

type entry struct {
	yp, scale       float32
	width           float32
	label, subLabel string
	labelAlpha      float32
	icon            string
	iconAlpha       float32
	cursor          struct {
		alpha float32
		yp    float32
	}
	ptr           int
	callbackOK    menuCallback
	callbackValue menuCallbackGetValue
	callbackIncr  menuCallbackIncr
	children      []entry
}

type scene interface {
	segueMount()
	segueNext()
	segueBack()
	update()
	render()
	Entry() *entry
}

var menu struct {
	stack         []scene
	icons         map[string]uint32
	inputCooldown int
	tweens        map[*float32]*gween.Tween
	scroll        float32
	ratio         float32
	t             float64
}

func menuRender() {
	menu.t += 0.1
	w, _ := window.GetFramebufferSize()
	menu.ratio = float32(w) / 1920

	fullscreenViewport()

	for e, t := range menu.tweens {
		var finished bool
		*e, finished = t.Update(1.0 / 60.0)
		if finished {
			delete(menu.tweens, e)
		}
	}

	currentScreenIndex := len(menu.stack) - 1
	for i := 0; i <= currentScreenIndex+1; i++ {
		if i < 0 || i > currentScreenIndex {
			continue
		}

		menu := menu.stack[i]
		menu.render()
	}
}

func genericSegueMount(list *entry) {
	for i := range list.children {
		e := &list.children[i]

		if i == list.ptr {
			e.yp = 0.5 + 0.3
			e.labelAlpha = 0
			e.iconAlpha = 0
			e.scale = 1.0
		} else if i < list.ptr {
			e.yp = 0.4 + 0.3 + 0.08*float32(i-list.ptr)
			e.labelAlpha = 0
			e.iconAlpha = 0
			e.scale = 0.5
		} else if i > list.ptr {
			e.yp = 0.6 + 0.3 + 0.08*float32(i-list.ptr)
			e.labelAlpha = 0
			e.iconAlpha = 0
			e.scale = 0.5
		}
		e.cursor.alpha = 0
		e.cursor.yp = 0.5 + 0.3
	}

	genericAnimate(list)
}

func genericAnimate(list *entry) {
	for i := range list.children {
		e := &list.children[i]

		var yp, la, a, s float32
		if i == list.ptr {
			yp = 0.5
			la = 1.0
			a = 1.0
			s = 1.0
		} else if i < list.ptr {
			yp = 0.4 + 0.08*float32(i-list.ptr)
			la = 0.5
			a = 0.5
			s = 0.5
		} else if i > list.ptr {
			yp = 0.6 + 0.08*float32(i-list.ptr)
			la = 0.5
			a = 0.5
			s = 0.5
		}

		menu.tweens[&e.yp] = gween.New(e.yp, yp, 0.15, ease.OutSine)
		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, la, 0.15, ease.OutSine)
		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, a, 0.15, ease.OutSine)
		menu.tweens[&e.scale] = gween.New(e.scale, s, 0.15, ease.OutSine)
	}
	menu.tweens[&list.cursor.alpha] = gween.New(list.cursor.alpha, 0.1, 0.15, ease.OutSine)
	menu.tweens[&list.cursor.yp] = gween.New(list.cursor.yp, 0.5, 0.15, ease.OutSine)
}

func genericSegueNext(list *entry) {
	for i := range list.children {
		e := &list.children[i]

		var yp, la, a, s float32
		if i == list.ptr {
			yp = 0.5 - 0.3
			la = 0
			a = 0
			s = 1.0
		} else if i < list.ptr {
			yp = 0.4 - 0.3 + 0.08*float32(i-list.ptr)
			la = 0
			a = 0
			s = 0.5
		} else if i > list.ptr {
			yp = 0.6 - 0.3 + 0.08*float32(i-list.ptr)
			la = 0
			a = 0
			s = 0.5
		}

		menu.tweens[&e.yp] = gween.New(e.yp, yp, 0.15, ease.OutSine)
		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, la, 0.15, ease.OutSine)
		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, a, 0.15, ease.OutSine)
		menu.tweens[&e.scale] = gween.New(e.scale, s, 0.15, ease.OutSine)
	}
	menu.tweens[&list.cursor.alpha] = gween.New(list.cursor.alpha, 0, 0.15, ease.OutSine)
	menu.tweens[&list.cursor.yp] = gween.New(list.cursor.yp, 0.5+0.3, 0.15, ease.OutSine)
}

func drawCursor(list *entry) {
	w, h := window.GetFramebufferSize()
	alpha := list.cursor.alpha - float32(math.Cos(menu.t))*0.025 - 0.025
	drawQuad(
		60*menu.ratio, float32(h)*list.cursor.yp-50*menu.ratio,
		float32(w)-610*menu.ratio, float32(h)*list.cursor.yp-50*menu.ratio,
		60*menu.ratio, float32(h)*list.cursor.yp+50*menu.ratio,
		float32(w)-610*menu.ratio, float32(h)*list.cursor.yp+50*menu.ratio,
		color{1, 1, 1, alpha},
	)
}

func genericRender(list *entry) {
	w, h := window.GetFramebufferSize()

	drawCursor(list)

	for _, e := range list.children {
		if e.yp < -1.1 || e.yp > 1.1 {
			continue
		}

		fontOffset := 64 * 0.7 * menu.ratio * 0.3

		drawImage(menu.icons[e.icon],
			120*menu.ratio-64*e.scale*menu.ratio,
			float32(h)*e.yp-16*menu.ratio-64*e.scale*menu.ratio+fontOffset,
			128*menu.ratio, 128*menu.ratio,
			e.scale, color{1, 1, 1, e.iconAlpha})

		if e.labelAlpha > 0 {

			video.font.SetColor(1, 1, 1, e.labelAlpha)
			video.font.Printf(
				200*menu.ratio,
				float32(h)*e.yp+fontOffset,
				0.7*menu.ratio, e.label)

			if e.callbackValue != nil {
				lw := video.font.Width(0.7*menu.ratio, e.callbackValue())
				video.font.Printf(
					float32(w)-lw-650*menu.ratio,
					float32(h)*e.yp+fontOffset,
					0.7*menu.ratio, e.callbackValue())
			}
		}
	}
}

func contextReset() {
	video.white = newWhite()

	menu.icons = map[string]uint32{
		"main":       newImage("assets/main.png"),
		"file":       newImage("assets/file.png"),
		"folder":     newImage("assets/folder.png"),
		"subsetting": newImage("assets/subsetting.png"),
		"setting":    newImage("assets/setting.png"),
		"resume":     newImage("assets/resume.png"),
		"reset":      newImage("assets/reset.png"),
		"loadstate":  newImage("assets/loadstate.png"),
		"savestate":  newImage("assets/savestate.png"),
		"screenshot": newImage("assets/screenshot.png"),
		"add":        newImage("assets/add.png"),
		"scan":       newImage("assets/scan.png"),
		"Nintendo - Super Nintendo Entertainment System": newImage("assets/Nintendo - Super Nintendo Entertainment System.png"),
		"Sega - Mega Drive - Genesis":                    newImage("assets/Sega - Mega Drive - Genesis.png"),
	}
}

func fastForwardTweens() {
	for e, t := range menu.tweens {
		var finished bool
		*e, finished = t.Update(1)
		if finished {
			delete(menu.tweens, e)
		}
	}
}

func menuInit() {
	w, _ := window.GetFramebufferSize()
	menu.stack = []scene{}
	menu.tweens = make(map[*float32]*gween.Tween)
	menu.ratio = float32(w) / 1920

	if g.coreRunning {
		menu.stack = append(menu.stack, buildTabs())
		menu.stack[0].segueNext()
		menu.stack = append(menu.stack, buildMainMenu())
		menu.stack[1].segueNext()
		menu.stack = append(menu.stack, buildQuickMenu())
		fastForwardTweens()
	} else {
		menu.stack = append(menu.stack, buildTabs())
	}
}
