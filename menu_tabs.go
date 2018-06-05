package main

import (
	"github.com/kivutar/go-playthemall/libretro"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

func buildTabs() entry {
	var list entry
	list.label = "Play Them All"
	list.input = inputTabs
	list.render = renderTabs

	list.children = append(list.children, entry{
		label: "Main Menu",
		icon:  "setting",
		callback: func() {
			menu.stack = append(menu.stack, buildMainMenu())
		},
	})

	list.children = append(list.children, entry{
		label: "Settings",
		icon:  "setting",
		callback: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label: "Nintendo - Super Nintendo Entertainment System",
		icon:  "Nintendo - Super Nintendo Entertainment System",
		callback: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label: "Sega - Mega Drive - Genesis",
		icon:  "Sega - Mega Drive - Genesis",
		callback: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	w, _ := window.GetFramebufferSize()
	for i := range list.children {
		e := &list.children[i]

		if i == list.ptr {
			e.x = float32(w / 2)
			e.labelAlpha = 1.0
		} else if i < list.ptr {
			e.x = float32(w/2) + float32(menu.spacing*2*(i-list.ptr)-menu.spacing*2)
			e.labelAlpha = 0
		} else if i > list.ptr {
			e.x = float32(w/2) + float32(menu.spacing*2*(i-list.ptr)+menu.spacing*2)
			e.labelAlpha = 0
		}
	}

	return list
}

func animateTabs() {
	w, _ := window.GetFramebufferSize()
	currentMenu := &menu.stack[len(menu.stack)-1]

	for i := range currentMenu.children {
		e := &currentMenu.children[i]

		var nx float32
		var nla float32
		if i == currentMenu.ptr {
			nx = float32(w / 2)
			nla = 1.0
		} else if i < currentMenu.ptr {
			nx = float32(w/2) + float32(menu.spacing*2*(i-currentMenu.ptr)-menu.spacing*2)
			nla = 0
		} else if i > currentMenu.ptr {
			nx = float32(w/2) + float32(menu.spacing*2*(i-currentMenu.ptr)+menu.spacing*2)
			nla = 0
		}

		menu.tweens[&e.x] = gween.New(e.x, nx, 0.15, ease.OutSine)
		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, nla, 0.15, ease.OutSine)
	}
}

func inputTabs() {
	currentMenu := &menu.stack[len(menu.stack)-1]

	if menu.inputCooldown > 0 {
		menu.inputCooldown--
	}

	if newState[0][libretro.DeviceIDJoypadRight] && menu.inputCooldown == 0 {
		currentMenu.ptr++
		if currentMenu.ptr >= len(currentMenu.children) {
			currentMenu.ptr = 0
		}
		animateTabs()
		menu.inputCooldown = 10
	}

	if newState[0][libretro.DeviceIDJoypadLeft] && menu.inputCooldown == 0 {
		currentMenu.ptr--
		if currentMenu.ptr < 0 {
			currentMenu.ptr = len(currentMenu.children) - 1
		}
		animateTabs()
		menu.inputCooldown = 10
	}

	commonInput()
}

func renderTabs() {
	w, h := window.GetFramebufferSize()
	currentMenu := &menu.stack[len(menu.stack)-1]

	for _, e := range currentMenu.children {
		if e.x < -128 || e.x > float32(w+128) {
			continue
		}

		video.font.SetColor(1.0, 1.0, 1.0, e.labelAlpha)
		video.font.Printf(e.x, float32(h/2)+100, 0.5, e.label)

		drawImage(menu.icons[e.icon], int32(e.x)-64, int32(h/2)-64, 128, 128)
	}
}
