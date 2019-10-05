package menu

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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

func isWindowsDrive(path string) bool {
	path, _ = filepath.Abs(path)
	validWindowsDrive := regexp.MustCompile(`^[A-Z]\:\\$`)
	return validWindowsDrive.MatchString(path)
}

func getWindowsDrives() (drives []string) {
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		fd, err := os.Open(string(drive) + ":\\")
		if err == nil {
			drives = append(drives, string(drive))
			fd.Close()
		}
	}
	return drives
}

func appendFolder(list *sceneExplorer, label, newPath string, exts []string, cb func(string), dirAction *entry) {
	list.children = append(list.children, entry{
		label: label,
		icon:  "folder",
		callbackOK: func() {
			list.segueNext()
			newPath := newPath
			if dirAction != nil {
				dirAction.callbackOK = func() { cb(newPath) }
			}
			menu.Push(buildExplorer(newPath, exts, cb, dirAction))
		},
	})
}

func explorerIcon(f os.FileInfo) string {
	icon := "file"
	if f.IsDir() {
		icon = "folder"
	}
	return icon
}

func appendNode(list *sceneExplorer, path string, f os.FileInfo, exts []string, cb func(string), dirAction *entry) {
	// Check whether or not we are to display hidden files.
	if f.Name()[:1] == "." && !settings.Current.ShowHiddenFiles {
		return
	}

	// Filter files by extension.
	if exts != nil && !f.IsDir() && !matchesExtensions(f, exts) {
		return
	}

	list.children = append(list.children, entry{
		label: f.Name(),
		icon:  explorerIcon(f),
		callbackOK: func() {
			if f.IsDir() {
				list.segueNext()
				newPath := filepath.Clean(filepath.Join(path, f.Name()))
				if dirAction != nil {
					dirAction.callbackOK = func() { cb(newPath) }
				}
				menu.Push(buildExplorer(newPath, exts, cb, dirAction))
			} else if cb != nil && (exts == nil || utils.StringInSlice(filepath.Ext(f.Name()), exts)) {
				cb(filepath.Clean(filepath.Join(path, f.Name())))
			}
		},
	})
}

func buildExplorer(path string, exts []string, cb func(string), dirAction *entry) Scene {
	var list sceneExplorer
	list.label = "Explorer"

	// Display the special directory action entry.
	if dirAction != nil && dirAction.label != "" {
		dirAction.callbackOK = func() { cb(path) }
		list.children = append(list.children, *dirAction)
	}

	if path == "/" {
		// On Windows there is no / and we want to display a list of drives instead
		if runtime.GOOS == "windows" {
			drives := getWindowsDrives()
			for _, drive := range drives {
				appendFolder(&list, drive+":\\", drive+":\\", exts, cb, dirAction)
			}
			list.segueMount()
			return &list
		}
	} else if isWindowsDrive(path) {
		// Special .. entry pointing to the list of drives on Windows
		appendFolder(&list, "..", "/", exts, cb, dirAction)
	} else {
		// Add a first entry for the parent directory.
		appendFolder(&list, "..", filepath.Clean(path+"/.."), exts, cb, dirAction)
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
	}

	// Loop over files in the directory and add one entry for each.
	for _, f := range files {
		f := f
		appendNode(&list, path, f, exts, cb, dirAction)
	}
	buildIndexes(&list.entry)

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
