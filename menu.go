package main

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"

	"github.com/fatih/structs"
	"github.com/kivutar/go-playthemall/libretro"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type menuCallback func()
type menuCallbackIncr func(int)
type menuCallbackGetValue func() string
type entry struct {
	label         string
	icon          string
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
}

func buildExplorer(path string) entry {
	var list entry
	list.label = "Explorer"
	list.input = verticalInput
	list.render = verticalRender

	files, err := ioutil.ReadDir(path)
	if err != nil {
		notifyAndLog("Menu", err.Error())
	}

	for _, f := range files {
		f := f
		icon := "file"
		if f.IsDir() {
			icon = "folder"
		}
		list.children = append(list.children, entry{
			label: f.Name(),
			icon:  icon,
			callback: func() {
				if f.IsDir() {
					menu.stack = append(menu.stack, buildExplorer(path+"/"+f.Name()+"/"))
				} else if stringInSlice(filepath.Ext(f.Name()), []string{".so", ".dll", ".dylib"}) {
					g.corePath = path + "/" + f.Name()
					coreLoad(g.corePath)
				} else {
					coreLoadGame(path + "/" + f.Name())
				}
			},
		})
	}

	return list
}

func buildSettings() entry {
	var list entry
	list.label = "Settings"
	list.input = verticalInput
	list.render = verticalRender

	fields := structs.Fields(&settings)
	for _, f := range fields {
		f := f
		list.children = append(list.children, entry{
			label: f.Tag("label"),
			icon:  "subsetting",
			callbackIncr: func(direction int) {
				incrCallbacks[f.Name()](f, direction)
			},
			callbackValue: func() string {
				return fmt.Sprintf(f.Tag("fmt"), f.Value())
			},
		})
	}

	return list
}

func buildTabs() entry {
	var list entry
	list.label = "Play Them All"
	list.input = horizontalInput
	list.render = tabsRender

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

func buildMainMenu() entry {
	var list entry
	list.label = "Main list"
	list.input = verticalInput
	list.render = verticalRender

	usr, _ := user.Current()

	if g.coreRunning {
		list.children = append(list.children, entry{
			label: "Quick Menu",
			icon:  "subsetting",
			callback: func() {
				menu.stack = append(menu.stack, buildQuickMenu())
			},
		})
	}

	list.children = append(list.children, entry{
		label: "Load Core",
		icon:  "subsetting",
		callback: func() {
			menu.stack = append(menu.stack, buildExplorer(usr.HomeDir))
		},
	})

	list.children = append(list.children, entry{
		label: "Load Game",
		icon:  "subsetting",
		callback: func() {
			menu.stack = append(menu.stack, buildExplorer(usr.HomeDir))
		},
	})

	list.children = append(list.children, entry{
		label: "Settings",
		icon:  "subsetting",
		callback: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label: "Help",
		icon:  "subsetting",
		callback: func() {
			notifyAndLog("Menu", "Not implemented yet.")
		},
	})

	list.children = append(list.children, entry{
		label: "Quit",
		icon:  "subsetting",
		callback: func() {
			window.SetShouldClose(true)
		},
	})

	return list
}

func buildQuickMenu() entry {
	var list entry
	list.label = "Quick Menu"
	list.input = verticalInput
	list.render = verticalRender

	list.children = append(list.children, entry{
		label: "Resume",
		icon:  "resume",
		callback: func() {
			g.menuActive = !g.menuActive
		},
	})

	list.children = append(list.children, entry{
		label: "Reset",
		icon:  "reset",
		callback: func() {
			g.core.Reset()
			g.menuActive = false
		},
	})

	list.children = append(list.children, entry{
		label: "Save State",
		icon:  "savestate",
		callback: func() {
			err := saveState()
			if err != nil {
				notifyAndLog("Menu", err.Error())
			} else {
				notifyAndLog("Menu", "State saved.")
			}
		},
	})

	list.children = append(list.children, entry{
		label: "Load State",
		icon:  "loadstate",
		callback: func() {
			err := loadState()
			if err != nil {
				notifyAndLog("Menu", err.Error())
			} else {
				g.menuActive = false
				notifyAndLog("Menu", "State loaded.")
			}
		},
	})

	list.children = append(list.children, entry{
		label: "Take Screenshot",
		icon:  "screenshot",
		callback: func() {
			takeScreenshot()
			notifyAndLog("Menu", "Took a screenshot.")
		},
	})

	return list
}

func menuInput() {
	currentMenu := &menu.stack[len(menu.stack)-1]
	currentMenu.input()
}

func verticalInput() {
	currentMenu := &menu.stack[len(menu.stack)-1]

	if menu.inputCooldown > 0 {
		menu.inputCooldown--
	}

	if newState[0][libretro.DeviceIDJoypadDown] && menu.inputCooldown == 0 {
		currentMenu.ptr++
		if currentMenu.ptr >= len(currentMenu.children) {
			currentMenu.ptr = 0
		}
		currentMenu.scrollTween = gween.New(currentMenu.scroll, float32(currentMenu.ptr*menu.spacing), 0.15, ease.OutSine)
		menu.inputCooldown = 10
	}

	if newState[0][libretro.DeviceIDJoypadUp] && menu.inputCooldown == 0 {
		currentMenu.ptr--
		if currentMenu.ptr < 0 {
			currentMenu.ptr = len(currentMenu.children) - 1
		}
		currentMenu.scrollTween = gween.New(currentMenu.scroll, float32(currentMenu.ptr*menu.spacing), 0.10, ease.OutSine)
		menu.inputCooldown = 10
	}

	commonInput()
}

func horizontalInput() {
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

func commonInput() {
	currentMenu := &menu.stack[len(menu.stack)-1]

	// OK
	if released[0][libretro.DeviceIDJoypadA] {
		if currentMenu.children[currentMenu.ptr].callback != nil {
			currentMenu.children[currentMenu.ptr].callback()
		}
	}

	// Right
	if released[0][libretro.DeviceIDJoypadRight] {
		if currentMenu.children[currentMenu.ptr].callbackIncr != nil {
			currentMenu.children[currentMenu.ptr].callbackIncr(1)
		}
	}

	// Left
	if released[0][libretro.DeviceIDJoypadLeft] {
		if currentMenu.children[currentMenu.ptr].callbackIncr != nil {
			currentMenu.children[currentMenu.ptr].callbackIncr(-1)
		}
	}

	// Cancel
	if released[0][libretro.DeviceIDJoypadB] {
		if len(menu.stack) > 1 {
			menu.stack = menu.stack[:len(menu.stack)-1]
		}
	}
}

func menuRender() {
	fullscreenViewport()

	currentMenu := &menu.stack[len(menu.stack)-1]
	if currentMenu.scrollTween != nil {
		currentMenu.scroll, _ = currentMenu.scrollTween.Update(1.0 / 60.0)
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

		drawImage(menu.icons[e.icon], 45, int32(y)-44, 64, 64)

		if e.callbackValue != nil {
			video.font.Printf(float32(w)-250, y, 0.5, e.callbackValue())
		}
	}
}

func tabsRender() {
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
	if g.coreRunning {
		menu.stack = append(menu.stack, buildTabs())
		menu.stack = append(menu.stack, buildMainMenu())
		menu.stack = append(menu.stack, buildQuickMenu())
	} else {
		menu.stack = append(menu.stack, buildTabs())
	}
}
