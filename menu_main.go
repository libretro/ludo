package main

import "os/user"

type screenMain struct {
	entry
}

func buildMainMenu() screen {
	var list screenMain
	list.label = "Main Menu"

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

	initEntries(list.entry)

	return &list
}

func (main *screenMain) update() {
	verticalInput(&main.entry)
}

func (main *screenMain) render() {
	verticalRender(&main.entry)
}
