package main

import "os/user"

func buildMainMenu() entry {
	var list entry
	list.label = "Main Menu"
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
