package menu

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/utils"
)

type sceneExplorer struct {
	entry
}

// Prettifier processes a file name
type Prettifier func(string) string

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

func appendFolder(list *sceneExplorer, label, newPath string, exts []string, cb func(string), dirAction *entry, prettifier Prettifier) {
	list.children = append(list.children, entry{
		label: label,
		icon:  "folder",
		callbackOK: func() {
			list.segueNext()
			newPath := newPath
			if dirAction != nil {
				dirAction.callbackOK = func() { cb(newPath) }
			}
			menu.Push(buildExplorer(newPath, exts, cb, dirAction, prettifier))
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

func appendNode(list *sceneExplorer, fullPath string, name string, f os.FileInfo, exts []string, cb func(string), dirAction *entry, prettifier Prettifier) {
	// Check whether or not we are to display hidden files.
	if name[:1] == "." && !settings.Current.ShowHiddenFiles {
		return
	}

	// Filter files by extension.
	if exts != nil && !f.IsDir() && !matchesExtensions(f, exts) {
		return
	}

	// Process file name if needed, used for user friendly file names
	displayName := name
	if prettifier != nil {
		displayName = prettifier(utils.FileName(name))
	}

	list.children = append(list.children, entry{
		label: displayName,
		icon:  explorerIcon(f),
		callbackOK: func() {
			if f.IsDir() {
				list.segueNext()
				newPath := filepath.Clean(fullPath)
				if dirAction != nil {
					dirAction.callbackOK = func() { cb(newPath) }
				}
				menu.Push(buildExplorer(newPath, exts, cb, dirAction, prettifier))
			} else if cb != nil && (exts == nil || utils.StringInSlice(filepath.Ext(name), exts)) {
				cb(filepath.Clean(fullPath))
			}
		},
	})
}

func buildExplorer(path string, exts []string, cb func(string), dirAction *entry, prettifier Prettifier) Scene {
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
				appendFolder(&list, drive+":\\", drive+":\\", exts, cb, dirAction, prettifier)
			}
			list.segueMount()
			return &list
		}
	} else if isWindowsDrive(path) {
		// Special .. entry pointing to the list of drives on Windows
		appendFolder(&list, "..", "/", exts, cb, dirAction, prettifier)
	} else {
		// Add a first entry for the parent directory.
		appendFolder(&list, "..", filepath.Clean(path+"/.."), exts, cb, dirAction, prettifier)
	}

	files, err := ioutil.ReadDir(path)

	// Sort entries by their labels, ignoring case.
	sort.SliceStable(files, func(i, j int) bool {
		if prettifier != nil {
			return strings.ToLower(prettifier(utils.FileName(files[i].Name()))) < strings.ToLower(prettifier(utils.FileName(files[j].Name())))
		}
		return strings.ToLower(utils.FileName(files[i].Name())) < strings.ToLower(utils.FileName(files[j].Name()))
	})

	if err != nil {
		ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
	}

	// Loop over files in the directory and add one entry for each.
	for _, f := range files {
		f := f
		fullPath, err := filepath.EvalSymlinks(filepath.Join(path, f.Name()))
		if err != nil {
			log.Println(err)
			continue
		}
		fi, err := os.Stat(fullPath)
		if err != nil {
			log.Println(err)
			continue
		}
		appendNode(&list, fullPath, f.Name(), fi, exts, cb, dirAction, prettifier)
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
