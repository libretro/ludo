package main

import (
	"github.com/kivutar/go-playthemall/libretro"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

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
