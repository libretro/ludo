package main

type screenQuick struct {
	entry
}

func buildQuickMenu() scene {
	var list screenQuick
	list.label = "Quick Menu"

	list.children = append(list.children, entry{
		label: "Resume",
		icon:  "resume",
		callbackOK: func() {
			g.menuActive = !g.menuActive
		},
	})

	list.children = append(list.children, entry{
		label: "Reset",
		icon:  "reset",
		callbackOK: func() {
			g.core.Reset()
			g.menuActive = false
		},
	})

	list.children = append(list.children, entry{
		label: "Save State",
		icon:  "savestate",
		callbackOK: func() {
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
		callbackOK: func() {
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
		callbackOK: func() {
			takeScreenshot()
			notifyAndLog("Menu", "Took a screenshot.")
		},
	})

	list.present()

	return &list
}

func (s *screenQuick) present() {
	initEntries(&s.entry)
}

func (s *screenQuick) makeRoomForChildren() {
	genericMakeRoomForChildren(&s.entry)
}

func (s *screenQuick) getFocusBack() {
	animateEntries(&s.entry)
}

func (s *screenQuick) update() {
	verticalInput(&s.entry)
}

func (s *screenQuick) render() {
	verticalRender(&s.entry)
}
