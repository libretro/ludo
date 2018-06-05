package main

import (
	"io/ioutil"
	"path/filepath"
)

func buildExplorer(path string) entry {
	var list entry
	list.label = "Explorer"
	list.input = verticalInput
	list.render = verticalRender

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
			callback: func() {
				if f.IsDir() {
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

	initEntries(list)

	return list
}
