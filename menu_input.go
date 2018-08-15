package main

import (
	"github.com/libretro/go-playthemall/input"
	"github.com/libretro/go-playthemall/libretro"
)

func menuInput() {
	currentMenu := menu.stack[len(menu.stack)-1]
	currentMenu.update()
}

func genericInput(list *entry) {
	if menu.inputCooldown > 0 {
		menu.inputCooldown--
	}

	if input.NewState[0][libretro.DeviceIDJoypadDown] && menu.inputCooldown == 0 {
		list.ptr++
		if list.ptr >= len(list.children) {
			list.ptr = 0
		}
		genericAnimate(list)
		menu.inputCooldown = 10
	}

	if input.NewState[0][libretro.DeviceIDJoypadUp] && menu.inputCooldown == 0 {
		list.ptr--
		if list.ptr < 0 {
			list.ptr = len(list.children) - 1
		}
		genericAnimate(list)
		menu.inputCooldown = 10
	}

	// OK
	if input.Released[0][libretro.DeviceIDJoypadA] {
		if list.children[list.ptr].callbackOK != nil {
			list.children[list.ptr].callbackOK()
		}
	}

	// Right
	if input.Released[0][libretro.DeviceIDJoypadRight] {
		if list.children[list.ptr].incr != nil {
			list.children[list.ptr].incr(1)
		}
	}

	// Left
	if input.Released[0][libretro.DeviceIDJoypadLeft] {
		if list.children[list.ptr].incr != nil {
			list.children[list.ptr].incr(-1)
		}
	}

	// Cancel
	if input.Released[0][libretro.DeviceIDJoypadB] {
		if len(menu.stack) > 1 {
			menu.stack[len(menu.stack)-2].segueBack()
			menu.stack = menu.stack[:len(menu.stack)-1]
		}
	}
}
