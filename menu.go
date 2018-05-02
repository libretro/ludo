package main

import (
	"fmt"
	"io/ioutil"
	"libretro"
)

type menuCallback func()

type entry struct {
	label    string
	scroll   float32
	ptr      int
	callback menuCallback
	children []entry
}

var menuStack []entry

func buildExplorer(path string) entry {
	var menu entry
	menu.label = "Explorer"

	files, err := ioutil.ReadDir(path)
	if err != nil {
		notify(err.Error(), 120)
		fmt.Println(err)
	}

	for _, f := range files {
		f := f
		menu.children = append(menu.children, entry{
			label: f.Name(),
			callback: func() {
				if f.IsDir() {
					menuStack = append(menuStack, buildExplorer(path+"/"+f.Name()+"/"))
				}
			},
		})
	}

	return menu
}

func buildQuickMenu() entry {
	var menu entry
	menu.label = "Quick Menu"

	menu.children = append(menu.children, entry{
		label: "Resume",
		callback: func() {
			menuActive = !menuActive
		},
	})

	menu.children = append(menu.children, entry{
		label: "Save State",
		callback: func() {
			fmt.Println("[Menu]: Not implemented")
			notify("Not implemented", 120)
		},
	})

	menu.children = append(menu.children, entry{
		label: "Load State",
		callback: func() {
			fmt.Println("[Menu]: Not implemented")
			notify("Not implemented", 120)
		},
	})

	menu.children = append(menu.children, entry{
		label: "Take Screenshot",
		callback: func() {
			fmt.Println("[Menu]: Not implemented")
			notify("Not implemented", 120)
		},
	})

	menu.children = append(menu.children, entry{
		label: "Explorer",
		callback: func() {
			menuStack = append(menuStack, buildExplorer("/Users/kivutar/testroms"))
		},
	})

	menu.children = append(menu.children, entry{
		label: "Quit",
		callback: func() {
			window.SetShouldClose(true)
		},
	})

	return menu
}

func menuInput() {
	currentMenu := &menuStack[len(menuStack)-1]

	if pressed[0][libretro.DeviceIDJoypadDown] {
		currentMenu.ptr++
		if currentMenu.ptr >= len(currentMenu.children) {
			currentMenu.ptr = 0
		}
	}

	if pressed[0][libretro.DeviceIDJoypadUp] {
		currentMenu.ptr--
		if currentMenu.ptr < 0 {
			currentMenu.ptr = len(currentMenu.children) - 1
		}
	}

	if pressed[0][libretro.DeviceIDJoypadA] {
		if currentMenu.children[currentMenu.ptr].callback != nil {
			currentMenu.children[currentMenu.ptr].callback()
		}
	}

	if pressed[0][libretro.DeviceIDJoypadB] {
		if len(menuStack) > 1 {
			menuStack = menuStack[:len(menuStack)-1]
		}
	}
}

func renderMenuList() {
	vSpacing := 70
	currentMenu := menuStack[len(menuStack)-1]
	currentMenu.scroll = float32(currentMenu.ptr * vSpacing)

	video.font.SetColor(0, 0, 0, 1.0)
	video.font.Printf(60+2, 20+60+2, 0.5, currentMenu.label)
	video.font.SetColor(1, 1, 1, 1.0)
	video.font.Printf(60, 20+60, 0.5, currentMenu.label)

	for i, e := range currentMenu.children {
		y := -currentMenu.scroll + 20 + float32(vSpacing*(i+2))
		video.font.SetColor(0, 0, 0, 1.0)
		video.font.Printf(100+2, y+2, 0.5, e.label)
		if i == currentMenu.ptr {
			video.font.SetColor(0.0, 1.0, 0.0, 1.0)
		} else {
			video.font.SetColor(0.6, 0.6, 0.9, 1.0)
		}
		video.font.Printf(100, y, 0.5, e.label)
	}
}

func menuInit() {
	menuStack = append(menuStack, buildQuickMenu())
}
