package main

type screenQuick struct {
	entry
}

func buildQuickMenu() screen {
	var list screenQuick
	list.label = "Quick Menu"

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

	initEntries(list.entry)

	return &list
}

func (s *screenQuick) update() {
	verticalInput(&s.entry)
}

func (s *screenQuick) render() {
	verticalRender(&s.entry)
}
