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
	w, _ := window.GetFramebufferSize()

	for i := range list.children {
		e := &list.children[i]

		if i == list.ptr {
			e.x = float32(w / 2)
			e.labelAlpha = 1
			e.iconAlpha = 1
			e.scale = 1
		} else if i < list.ptr {
			e.x = float32(w/2) + float32(menu.spacing*2*(i-list.ptr)-menu.spacing*2)
			e.labelAlpha = 0
			e.iconAlpha = 0.5
			e.scale = 0.25
		} else if i > list.ptr {
			e.x = float32(w/2) + float32(menu.spacing*2*(i-list.ptr)+menu.spacing*2)
			e.labelAlpha = 0
			e.iconAlpha = 0.5
			e.scale = 0.25
		}
	}
}

func animateTabs() {
	w, _ := window.GetFramebufferSize()
	currentMenu := &menu.stack[len(menu.stack)-1]

	for i := range currentMenu.children {
		e := &currentMenu.children[i]

		var x, la, a, s float32
		if i == currentMenu.ptr {
			x = float32(w / 2)
			la = 1
			a = 1
			s = 1
		} else if i < currentMenu.ptr {
			x = float32(w/2) + float32(menu.spacing*2*(i-currentMenu.ptr)-menu.spacing*2)
			la = 0
			a = 0.5
			s = 0.25
		} else if i > currentMenu.ptr {
			x = float32(w/2) + float32(menu.spacing*2*(i-currentMenu.ptr)+menu.spacing*2)
			la = 0
			a = 0.5
			s = 0.25
		}

		menu.tweens[&e.x] = gween.New(e.x, x, 0.15, ease.OutSine)
		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, la, 0.15, ease.OutSine)
		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, a, 0.15, ease.OutSine)
		menu.tweens[&e.scale] = gween.New(e.scale, s, 0.15, ease.OutSine)
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

	drawPolygon(
		0, 0,
		float32(w)/2, 0,
		0, float32(h),
		float32(w), float32(h),
		color{1, 1, 0, 1})

	for _, e := range currentMenu.children {
		if e.x < -128 || e.x > float32(w+128) {
			continue
		}

		video.font.SetColor(1.0, 1.0, 1.0, e.labelAlpha)
		lw := video.font.Width(0.75, e.label)
		video.font.Printf(e.x-lw/2, float32(h/2)+180, 0.75, e.label)
		lw = video.font.Width(0.5, e.subLabel)
		video.font.Printf(e.x-lw/2, float32(h/2)+260, 0.5, e.subLabel)

		drawImage(menu.icons[e.icon],
			e.x-128*e.scale, float32(h)/2-128*e.scale,
			256, 256, e.scale, color{1, 1, 1, e.iconAlpha})
	}
}
