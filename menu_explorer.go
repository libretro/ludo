package main

import (
	"io/ioutil"
	"path/filepath"
)

type screenExplorer struct {
	entry
}

func buildExplorer(path string) scene {
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
					list.segueNext()
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

	if len(files) == 0 {
		list.children = append(list.children, entry{
			label: "Empty",
			icon:  "subsetting",
		})
	}

	list.segueMount()

	return &list
}

func (explorer *screenExplorer) Entry() *entry {
	return &explorer.entry
}

func (explorer *screenExplorer) segueMount() {
	genericSegueMount(&explorer.entry)
}

func (explorer *screenExplorer) segueNext() {
	genericSegueNext(&explorer.entry)
}

func (explorer *screenExplorer) segueBack() {
	genericAnimate(&explorer.entry)
}

func (explorer *screenExplorer) update() {
	genericInput(&explorer.entry)
}

func (explorer *screenExplorer) render() {
	genericRender(&explorer.entry)
}
