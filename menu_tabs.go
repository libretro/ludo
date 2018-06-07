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
		label:    "Settings",
		subLabel: "Configure Play Them All",
		icon:     "setting",
		callback: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label:    "Super NES",
		subLabel: "10 Games - 5 Favorites",
		icon:     "Nintendo - Super Nintendo Entertainment System",
		callback: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label:    "Mega Drive - Genesis",
		subLabel: "10 Games - 5 Favorites",
		icon:     "Sega - Mega Drive - Genesis",
		callback: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	initTabs(list)

	return list
}

func initTabs(list entry) {
	w, h := window.GetFramebufferSize()

	for i := range list.children {
		e := &list.children[i]

		if i == list.ptr {
			e.x = float32(w / 2)
			e.y = float32(h / 2)
			e.labelAlpha = 1
			e.iconAlpha = 1
			e.scale = 1
			e.width = 1000
		} else if i < list.ptr {
			e.x = float32(w/2) + float32(128*(i-list.ptr)-128*2)
			e.y = 64
			e.labelAlpha = 0
			e.iconAlpha = 0.5
			e.scale = 0.25
			e.width = 128
		} else if i > list.ptr {
			e.x = float32(w/2) + float32(128*(i-list.ptr)+128*2)
			e.y = float32(h) - 64
			e.labelAlpha = 0
			e.iconAlpha = 0.5
			e.scale = 0.25
			e.width = 128
		}
	}
}

func animateTabs() {
	w, h := window.GetFramebufferSize()
	currentMenu := &menu.stack[len(menu.stack)-1]

	for i := range currentMenu.children {
		e := &currentMenu.children[i]

		var x, y, labelAlpha, iconAlpha, scale, width float32
		if i == currentMenu.ptr {
			x = float32(w / 2)
			y = float32(h / 2)
			labelAlpha = 1
			iconAlpha = 1
			scale = 1
			width = 1000
		} else if i < currentMenu.ptr {
			x = float32(w/2) + float32(128*(i-currentMenu.ptr)-128*2)
			y = 64
			labelAlpha = 0
			iconAlpha = 0.5
			scale = 0.25
			width = 128
		} else if i > currentMenu.ptr {
			x = float32(w/2) + float32(128*(i-currentMenu.ptr)+128*2)
			y = float32(h) - 64
			labelAlpha = 0
			iconAlpha = 0.5
			scale = 0.25
			width = 128
		}

		menu.tweens[&e.x] = gween.New(e.x, x, 0.15, ease.OutSine)
		menu.tweens[&e.y] = gween.New(e.y, y, 0.15, ease.OutSine)
		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, labelAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, iconAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.scale] = gween.New(e.scale, scale, 0.15, ease.OutSine)
		menu.tweens[&e.width] = gween.New(e.width, width, 0.15, ease.OutSine)
	}
	menu.tweens[&menu.scroll] = gween.New(menu.scroll, float32(currentMenu.ptr*128), 0.15, ease.OutSine)
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

	var stackWidth float32 = 256
	for i, e := range currentMenu.children {

		drawPolygon(
			-menu.scroll+stackWidth, 0,
			-menu.scroll+stackWidth+e.width, 0,
			-menu.scroll+stackWidth+400, float32(h),
			-menu.scroll+stackWidth+400+e.width, float32(h),
			color{0.5, 0.5, float32(i) / 10, 1})

		stackWidth += e.width
	}

	for _, e := range currentMenu.children {
		if e.x < -128 || e.x > float32(w+128) {
			continue
		}

		video.font.SetColor(1.0, 1.0, 1.0, e.labelAlpha)
		lw := video.font.Width(0.75, e.label)
		video.font.Printf(e.x-lw/2, e.y+180, 0.75, e.label)
		lw = video.font.Width(0.5, e.subLabel)
		video.font.Printf(e.x-lw/2, e.y+260, 0.5, e.subLabel)

		drawImage(menu.icons[e.icon],
			e.x-128*e.scale, e.y-128*e.scale,
			256, 256, e.scale, color{1, 1, 1, e.iconAlpha})
	}
}
