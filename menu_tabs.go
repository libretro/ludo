package main

import (
	"github.com/kivutar/go-playthemall/libretro"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type screenTabs struct {
	entry
}

func buildTabs() screen {
	var list screenTabs
	list.label = "Play Them All"

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

	list.init()

	return &list
}

func (tabs screenTabs) init() {
	w, h := window.GetFramebufferSize()

	for i := range tabs.children {
		e := &tabs.children[i]

		if i == tabs.ptr {
			e.x = float32(w / 2)
			e.y = float32(h / 2)
			e.labelAlpha = 1
			e.iconAlpha = 1
			e.scale = 1
			e.width = 1000
		} else if i < tabs.ptr {
			e.x = float32(w/2) + float32(128*(i-tabs.ptr)-128*2)
			e.y = 64
			e.labelAlpha = 0
			e.iconAlpha = 0.5
			e.scale = 0.25
			e.width = 128
		} else if i > tabs.ptr {
			e.x = float32(w/2) + float32(128*(i-tabs.ptr)+128*2)
			e.y = float32(h) - 64
			e.labelAlpha = 0
			e.iconAlpha = 0.5
			e.scale = 0.25
			e.width = 128
		}
	}
}

func (tabs screenTabs) animate() {
	w, h := window.GetFramebufferSize()

	for i := range tabs.children {
		e := &tabs.children[i]

		var x, y, labelAlpha, iconAlpha, scale, width float32
		if i == tabs.ptr {
			x = float32(w / 2)
			y = float32(h / 2)
			labelAlpha = 1
			iconAlpha = 1
			scale = 1
			width = 1000
		} else if i < tabs.ptr {
			x = float32(w/2) + float32(128*(i-tabs.ptr)-128*2)
			y = 64
			labelAlpha = 0
			iconAlpha = 0.5
			scale = 0.25
			width = 128
		} else if i > tabs.ptr {
			x = float32(w/2) + float32(128*(i-tabs.ptr)+128*2)
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
	menu.tweens[&menu.scroll] = gween.New(menu.scroll, float32(tabs.ptr*128), 0.15, ease.OutSine)
}

func (tabs *screenTabs) update() {
	if menu.inputCooldown > 0 {
		menu.inputCooldown--
	}

	if newState[0][libretro.DeviceIDJoypadRight] && menu.inputCooldown == 0 {
		tabs.ptr++
		if tabs.ptr >= len(tabs.children) {
			tabs.ptr = 0
		}
		tabs.animate()
		menu.inputCooldown = 10
	}

	if newState[0][libretro.DeviceIDJoypadLeft] && menu.inputCooldown == 0 {
		tabs.ptr--
		if tabs.ptr < 0 {
			tabs.ptr = len(tabs.children) - 1
		}
		tabs.animate()
		menu.inputCooldown = 10
	}

	commonInput(&tabs.entry)
}

func (tabs screenTabs) render() {
	w, h := window.GetFramebufferSize()

	var stackWidth float32 = 260
	for i, e := range tabs.children {

		c := colorful.Hcl(float64(i%12)*30, 0.5, 0.5)

		drawPolygon(
			-menu.scroll+stackWidth, 0,
			-menu.scroll+stackWidth+e.width, 0,
			-menu.scroll+stackWidth+400, float32(h),
			-menu.scroll+stackWidth+400+e.width, float32(h),
			color{float32(c.R), float32(c.G), float32(c.B), 1})

		stackWidth += e.width
	}

	for _, e := range tabs.children {
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
