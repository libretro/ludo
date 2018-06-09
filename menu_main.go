package main

import (
	"os/user"
)

type screenMain struct {
	entry
}

func buildMainMenu() scene {
	var list screenMain
	list.label = "Main Menu"

	usr, _ := user.Current()

	if g.coreRunning {
		list.children = append(list.children, entry{
			label: "Quick Menu",
			icon:  "subsetting",
			callbackOK: func() {
				list.segueNext()
				menu.stack = append(menu.stack, buildQuickMenu())
			},
		})
	}

	list.children = append(list.children, entry{
		label: "Load Core",
		icon:  "subsetting",
		callbackOK: func() {
			list.segueNext()
			menu.stack = append(menu.stack, buildExplorer(usr.HomeDir))
		},
	})

	list.children = append(list.children, entry{
		label: "Load Game",
		icon:  "subsetting",
		callbackOK: func() {
			list.segueNext()
			menu.stack = append(menu.stack, buildExplorer(usr.HomeDir))
		},
	})

	list.children = append(list.children, entry{
		label: "Settings",
		icon:  "subsetting",
		callbackOK: func() {
			list.segueNext()
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

	list.segueMount()

	return &list
}

func (main *screenMain) segueMount() {
	genericSegueMount(&main.entry)
}

func (main *screenMain) segueBack() {
	genericAnimate(&main.entry)
}

func (main *screenMain) segueNext() {
	genericSegueNext(&main.entry)
}

func (main *screenMain) update() {
	genericInput(&main.entry)
}

func (main *screenMain) render() {
	genericRender(&main.entry)
}
