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
			callbackOK: func() {
				menu.stack = append(menu.stack, buildQuickMenu())
			},
		})
	}

	list.children = append(list.children, entry{
		label: "Load Core",
		icon:  "subsetting",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildExplorer(usr.HomeDir))
		},
	})

	list.children = append(list.children, entry{
		label: "Load Game",
		icon:  "subsetting",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildExplorer(usr.HomeDir))
		},
	})

	list.children = append(list.children, entry{
		label: "Settings",
		icon:  "subsetting",
		callbackOK: func() {
			menu.stack = append(menu.stack, buildSettings())
		},
	})

	list.children = append(list.children, entry{
		label: "Help",
		icon:  "subsetting",
		callbackOK: func() {
			notifyAndLog("Menu", "Not implemented yet.")
		},
	})

	list.children = append(list.children, entry{
		label: "Quit",
		icon:  "subsetting",
		callbackOK: func() {
			window.SetShouldClose(true)
		},
	})

	list.open()

	return &list
}

func (main *screenMain) open() {
	initEntries(&main.entry)
}

func (main *screenMain) close() {
	animateEntries(&main.entry)
}

func (main *screenMain) update() {
	verticalInput(&main.entry)
}

func (main *screenMain) render() {
	verticalRender(&main.entry)
}
