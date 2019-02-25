package menu

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/utils"
)

type sceneExplorer struct {
	entry
}

func matchesExtensions(f os.FileInfo, exts []string) bool {
	if len(exts) > 0 {
		var fileExtension = filepath.Ext(f.Name())
		for _, ext := range exts {
			if ext == fileExtension {
				return true
			}
		}
	}
	return false
}

func getWindowsDrives() (r []string) {
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		fd, err := os.Open(string(drive) + ":\\")
		if err == nil {
			r = append(r, string(drive))
			fd.Close()
		}
	}
	return r
}

func buildExplorer(path string, exts []string, cb func(string), dirAction *entry) Scene {
	var list sceneExplorer
	list.label = "Explorer"

	files, err := ioutil.ReadDir(path)
	if err != nil {
		ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
	}

	// Display the special directory action entry.
	if dirAction != nil && dirAction.label != "" {
		dirAction.callbackOK = func() { cb(path) }
		list.children = append(list.children, *dirAction)
	}

	if path == "/" {
		// Add windows drives.
		if runtime.GOOS == "windows" {
			drives := getWindowsDrives()
			for _, drive := range drives {
				list.children = append(list.children, entry{
					label: drive + ":\\",
					icon:  "folder",
					callbackOK: func() {
						list.segueNext()
						newPath := drive + ":\\"
						if dirAction != nil {
							dirAction.callbackOK = func() { cb(newPath) }
						}
						menu.stack = append(menu.stack, buildExplorer(newPath, exts, cb, dirAction))
					},
				})
			}
		}
	} else {
		// Add a first entry for the parent directory.
		list.children = append(list.children, entry{
			label: "..",
			icon:  "folder",
			callbackOK: func() {
				list.segueNext()
				newPath := filepath.Clean(path + "/..")
				if dirAction != nil {
					dirAction.callbackOK = func() { cb(newPath) }
				}
				menu.stack = append(menu.stack, buildExplorer(newPath, exts, cb, dirAction))
			},
		})
	}

	// Loop over files in the directory and add one entry for each.
	for _, f := range files {
		f := f
		icon := "file"

		// Check whether or not we are to display hidden files.
		if f.Name()[:1] == "." && settings.Current.ShowHiddenFiles {
			continue
		}

		// Filter files by extension.
		if exts != nil && !f.IsDir() && !matchesExtensions(f, exts) {
			continue
		}

		if f.IsDir() {
			icon = "folder"
		}

		list.children = append(list.children, entry{
			label: f.Name(),
			icon:  icon,
			callbackOK: func() {
				if f.IsDir() {
					list.segueNext()
					newPath := path + "/" + f.Name()
					if dirAction != nil {
						dirAction.callbackOK = func() { cb(newPath) }
					}
					menu.stack = append(menu.stack, buildExplorer(newPath, exts, cb, dirAction))
				} else if cb != nil && (exts == nil || utils.StringInSlice(filepath.Ext(f.Name()), exts)) {
					cb(path + "/" + f.Name())
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

func (explorer *sceneExplorer) Entry() *entry {
	return &explorer.entry
}

func (explorer *sceneExplorer) segueMount() {
	genericSegueMount(&explorer.entry)
}

func (explorer *sceneExplorer) segueNext() {
	genericSegueNext(&explorer.entry)
}

func (explorer *sceneExplorer) segueBack() {
	genericAnimate(&explorer.entry)
}

func (explorer *sceneExplorer) update(dt float32) {
	genericInput(&explorer.entry, dt)
}

func (explorer *sceneExplorer) render() {
	genericRender(&explorer.entry)
}

func (explorer *sceneExplorer) drawHintBar() {
	genericDrawHintBar()
}
