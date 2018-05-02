package main

import (
	"fmt"
	"io/ioutil"
	"libretro"
)

type menuCallback func()

type entry struct {
	label    string
	callback menuCallback
	ptr      int
	children []entry
}

var currentMenu entry

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
					fmt.Println(path + "/" + f.Name())
					currentMenu = buildExplorer(path + "/" + f.Name() + "/")
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
			currentMenu = buildExplorer("/Users/kivutar/testroms")
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

	}
}

func renderMenuList() {
	video.font.SetColor(0, 0, 0, 1.0)
	video.font.Printf(60+2, 20+60+2, 0.5, currentMenu.label)
	video.font.SetColor(1, 1, 1, 1.0)
	video.font.Printf(60, 20+60, 0.5, currentMenu.label)
	for i, e := range currentMenu.children {
		video.font.SetColor(0, 0, 0, 1.0)
		video.font.Printf(100+2, 20+float32(70*(i+2))+2, 0.5, e.label)
		if i == currentMenu.ptr {
			video.font.SetColor(0.0, 1.0, 0.0, 1.0)
		} else {
			video.font.SetColor(0.6, 0.6, 0.9, 1.0)
		}
		video.font.Printf(100, 20+float32(70*(i+2)), 0.5, e.label)
	}
}

func menuInit() {
	currentMenu = buildQuickMenu()
}
