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
		callbackOK: func() {
			menu.stack = append(menu.stack, buildMainMenu())
		},
	})

	list.children = append(list.children, entry{
		label:    "Settings",
		subLabel: "Configure Play Them All",
		icon:     "setting",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label:    "Super NES",
		subLabel: "10 Games - 5 Favorites",
		icon:     "Nintendo - Super Nintendo Entertainment System",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label:    "Mega Drive - Genesis",
		subLabel: "10 Games - 5 Favorites",
		icon:     "Sega - Mega Drive - Genesis",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label:    "Super NES",
		subLabel: "10 Games - 5 Favorites",
		icon:     "Nintendo - Super Nintendo Entertainment System",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label:    "Mega Drive - Genesis",
		subLabel: "10 Games - 5 Favorites",
		icon:     "Sega - Mega Drive - Genesis",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label:    "Super NES",
		subLabel: "10 Games - 5 Favorites",
		icon:     "Nintendo - Super Nintendo Entertainment System",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label:    "Mega Drive - Genesis",
		subLabel: "10 Games - 5 Favorites",
		icon:     "Sega - Mega Drive - Genesis",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label:    "Super NES",
		subLabel: "10 Games - 5 Favorites",
		icon:     "Nintendo - Super Nintendo Entertainment System",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label:    "Mega Drive - Genesis",
		subLabel: "10 Games - 5 Favorites",
		icon:     "Sega - Mega Drive - Genesis",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label:    "Super NES",
		subLabel: "10 Games - 5 Favorites",
		icon:     "Nintendo - Super Nintendo Entertainment System",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label:    "Mega Drive - Genesis",
		subLabel: "10 Games - 5 Favorites",
		icon:     "Sega - Mega Drive - Genesis",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label:    "Super NES",
		subLabel: "10 Games - 5 Favorites",
		icon:     "Nintendo - Super Nintendo Entertainment System",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label:    "Mega Drive - Genesis",
		subLabel: "10 Games - 5 Favorites",
		icon:     "Sega - Mega Drive - Genesis",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.open()

	return &list
}

func (tabs *screenTabs) open() {
	_, h := window.GetFramebufferSize()

	for i := range tabs.children {
		e := &tabs.children[i]

		if i == tabs.ptr {
			e.y = float32(h / 2)
			e.labelAlpha = 1
			e.iconAlpha = 1
			e.scale = 1
			e.width = 1256
		} else if i < tabs.ptr {
			e.y = 64
			e.labelAlpha = 0
			e.iconAlpha = 0.5
			e.scale = 0.25
			e.width = 128
		} else if i > tabs.ptr {
			e.y = float32(h) - 64
			e.labelAlpha = 0
			e.iconAlpha = 0.5
			e.scale = 0.25
			e.width = 128
		}
	}
}

func (tabs *screenTabs) animate() {
	_, h := window.GetFramebufferSize()

	for i := range tabs.children {
		e := &tabs.children[i]

		var y, labelAlpha, iconAlpha, scale, width float32
		if i == tabs.ptr {
			y = float32(h / 2)
			labelAlpha = 1
			iconAlpha = 1
			scale = 1
			width = 1256
		} else if i < tabs.ptr {
			y = 64
			labelAlpha = 0
			iconAlpha = 0.5
			scale = 0.25
			width = 128
		} else if i > tabs.ptr {
			y = float32(h) - 64
			labelAlpha = 0
			iconAlpha = 0.5
			scale = 0.25
			width = 128
		}

		menu.tweens[&e.y] = gween.New(e.y, y, 0.15, ease.OutSine)
		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, labelAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, iconAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.scale] = gween.New(e.scale, scale, 0.15, ease.OutSine)
		menu.tweens[&e.width] = gween.New(e.width, width, 0.15, ease.OutSine)
	}
	menu.tweens[&menu.scroll] = gween.New(menu.scroll, float32(tabs.ptr*128), 0.15, ease.OutSine)
}

func (tabs *screenTabs) animateNext() {
	cur := &tabs.children[tabs.ptr]
	menu.tweens[&cur.width] = gween.New(cur.width, 4000, 0.15, ease.OutSine)
	menu.tweens[&menu.scroll] = gween.New(menu.scroll, menu.scroll+700, 0.15, ease.OutSine)
}

func (tabs *screenTabs) update() {
	if menu.inputCooldown > 0 {
		menu.inputCooldown--
	}

	// Right
	if newState[0][libretro.DeviceIDJoypadRight] && menu.inputCooldown == 0 {
		tabs.ptr++
		if tabs.ptr >= len(tabs.children) {
			tabs.ptr = 0
		}
		tabs.animate()
		menu.inputCooldown = 10
	}

	// Left
	if newState[0][libretro.DeviceIDJoypadLeft] && menu.inputCooldown == 0 {
		tabs.ptr--
		if tabs.ptr < 0 {
			tabs.ptr = len(tabs.children) - 1
		}
		tabs.animate()
		menu.inputCooldown = 10
	}

	// OK
	if released[0][libretro.DeviceIDJoypadA] {
		if tabs.children[tabs.ptr].callbackOK != nil {
			tabs.animateNext()
			tabs.children[tabs.ptr].callbackOK()
		}
	}

	// Cancel
	if released[0][libretro.DeviceIDJoypadB] {
		if len(menu.stack) > 1 {
			menu.stack = menu.stack[:len(menu.stack)-1]
		}
	}
}

func (tabs screenTabs) render() {
	_, h := window.GetFramebufferSize()

	var stackWidth float32 = 132
	for i, e := range tabs.children {

		c := colorful.Hcl(float64(i%12)*30, 0.5, 0.5)

		drawPolygon(
			-menu.scroll+stackWidth, 0,
			-menu.scroll+stackWidth+e.width, 0,
			-menu.scroll+stackWidth+400, float32(h),
			-menu.scroll+stackWidth+400+e.width, float32(h),
			color{float32(c.R), float32(c.G), float32(c.B), 1})

		stackWidth += e.width

		x := -menu.scroll + stackWidth - e.width/2 + 400 - 400/(float32(h)/e.y)

		video.font.SetColor(1.0, 1.0, 1.0, e.labelAlpha)
		lw := video.font.Width(0.7, e.label)
		video.font.Printf(x-lw/2, e.y+180, 0.7, e.label)
		lw = video.font.Width(0.4, e.subLabel)
		video.font.Printf(x-lw/2, e.y+260, 0.4, e.subLabel)

		drawImage(menu.icons[e.icon],
			x-128*e.scale, e.y-128*e.scale,
			256, 256, e.scale, color{1, 1, 1, e.iconAlpha})
	}
}
