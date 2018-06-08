package main

import (
	"github.com/kivutar/go-playthemall/libretro"
)

func menuInput() {
	currentMenu := menu.stack[len(menu.stack)-1]
	currentMenu.update()
}

func verticalInput(list *entry) {
	if menu.inputCooldown > 0 {
		menu.inputCooldown--
	}

	if newState[0][libretro.DeviceIDJoypadDown] && menu.inputCooldown == 0 {
		list.ptr++
		if list.ptr >= len(list.children) {
			list.ptr = 0
		}
		animateEntries(list)
		menu.inputCooldown = 10
	}

	if newState[0][libretro.DeviceIDJoypadUp] && menu.inputCooldown == 0 {
		list.ptr--
		if list.ptr < 0 {
			list.ptr = len(list.children) - 1
		}
		animateEntries(list)
		menu.inputCooldown = 10
	}

	commonInput(list)
}

func commonInput(list *entry) {
	// OK
	if released[0][libretro.DeviceIDJoypadA] {
		if list.children[list.ptr].callback != nil {
			list.children[list.ptr].callback()
		}
	}

	// Right
	if released[0][libretro.DeviceIDJoypadRight] {
		if list.children[list.ptr].callbackIncr != nil {
			list.children[list.ptr].callbackIncr(1)
		}
	}

	// Left
	if released[0][libretro.DeviceIDJoypadLeft] {
		if list.children[list.ptr].callbackIncr != nil {
			list.children[list.ptr].callbackIncr(-1)
		}
	}

	// Cancel
	if released[0][libretro.DeviceIDJoypadB] {
		if len(menu.stack) > 1 {
			menu.stack = menu.stack[:len(menu.stack)-1]
		}
	}
}
