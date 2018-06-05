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

	return list
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
		currentMenu.scrollTween = gween.New(currentMenu.scroll, float32(currentMenu.ptr*menu.spacing*3), 0.15, ease.OutSine)
		menu.inputCooldown = 10
	}

	if newState[0][libretro.DeviceIDJoypadLeft] && menu.inputCooldown == 0 {
		currentMenu.ptr--
		if currentMenu.ptr < 0 {
			currentMenu.ptr = len(currentMenu.children) - 1
		}
		currentMenu.scrollTween = gween.New(currentMenu.scroll, float32(currentMenu.ptr*menu.spacing*3), 0.10, ease.OutSine)
		menu.inputCooldown = 10
	}

	commonInput()
}

func renderTabs() {
	w, h := window.GetFramebufferSize()
	currentMenu := &menu.stack[len(menu.stack)-1]

	for i, e := range currentMenu.children {
		x := float32(w/2) - currentMenu.scroll + float32(menu.spacing*3*i)

		if x < 0 || x > float32(w) {
			continue
		}

		if i == currentMenu.ptr {
			video.font.SetColor(0.0, 1.0, 0.0, 1.0)
		} else {
			video.font.SetColor(0.6, 0.6, 0.9, 1.0)
		}
		video.font.Printf(x, float32(h/2), 0.5, e.label)

		drawImage(menu.icons[e.icon], int32(x)-64, int32(h/2)-64, 64, 64)
	}
}
