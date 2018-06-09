package main

import (
	"io/ioutil"
	"path/filepath"
)

type screenExplorer struct {
	entry
}

func buildExplorer(path string) screen {
	var list screenExplorer
	list.label = "Explorer"

	files, err := ioutil.ReadDir(path)
	if err != nil {
		notifyAndLog("Menu", err.Error())
	}

	for _, f := range files {
		f := f
		icon := "file"
		if f.IsDir() {
			icon = "folder"
		}
		list.children = append(list.children, entry{
			label: f.Name(),
			icon:  icon,
			callbackOK: func() {
				if f.IsDir() {
					list.makeRoomForChildren()
					menu.stack = append(menu.stack, buildExplorer(path+"/"+f.Name()+"/"))
				} else if stringInSlice(filepath.Ext(f.Name()), []string{".so", ".dll", ".dylib"}) {
					g.corePath = path + "/" + f.Name()
					coreLoad(g.corePath)
				} else {
					coreLoadGame(path + "/" + f.Name())
				}
			},
		})
	}

	list.present()

	return &list
}

func (explorer *screenExplorer) present() {
	initEntries(&explorer.entry)
}

func (explorer *screenExplorer) makeRoomForChildren() {
	genericMakeRoomForChildren(&explorer.entry)
}

func (explorer *screenExplorer) getFocusBack() {
	animateEntries(&explorer.entry)
}

func (explorer *screenExplorer) update() {
	verticalInput(&explorer.entry)
}

func (explorer *screenExplorer) render() {
	verticalRender(&explorer.entry)
}
