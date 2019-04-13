package menu

import (
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/libretro/ludo/settings"

	"github.com/libretro/ludo/core"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/state"
)

type sceneMain struct {
	entry
}

func buildMainMenu() Scene {
	var list sceneMain
	list.label = "Main Menu"

	usr, _ := user.Current()

	if state.Global.CoreRunning {
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
			menu.stack = append(menu.stack, buildExplorer(
				settings.Current.CoresDirectory,
				[]string{".dll", ".dylib", ".so"},
				coreExplorerCb,
				nil,
			))
		},
	})

	list.children = append(list.children, entry{
		label: "Load Game",
		icon:  "subsetting",
		callbackOK: func() {
			if state.Global.Core != nil {
				list.segueNext()
				menu.stack = append(menu.stack, buildExplorer(
					usr.HomeDir,
					nil,
					gameExplorerCb,
					nil,
				))
			} else {
				ntf.DisplayAndLog(ntf.Warning, "Menu", "Please load a core first.")
			}
		},
	})

	list.children = append(list.children, entry{
		label: "Test Keyboard",
		icon:  "subsetting",
		callbackOK: func() {
			list.segueNext()
			menu.stack = append(menu.stack, buildKeyboard())
		},
	})

	if state.Global.LudOS {
		list.children = append(list.children, entry{
			label: "Updater",
			icon:  "subsetting",
			callbackOK: func() {
				list.segueNext()
				menu.stack = append(menu.stack, buildUpdater())
			},
		})

		list.children = append(list.children, entry{
			label: "Reboot",
			icon:  "subsetting",
			callbackOK: func() {
				cmd := exec.Command("/usr/sbin/shutdown", "-r", "now")
				core.UnloadGame()
				err := cmd.Run()
				if err != nil {
					ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
				}
			},
		})

		list.children = append(list.children, entry{
			label: "Shutdown",
			icon:  "subsetting",
			callbackOK: func() {
				cmd := exec.Command("/usr/sbin/shutdown", "-P", "now")
				core.UnloadGame()
				err := cmd.Run()
				if err != nil {
					ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
				}
			},
		})
	} else {
		list.children = append(list.children, entry{
			label: "Quit",
			icon:  "subsetting",
			callbackOK: func() {
				vid.Window.SetShouldClose(true)
			},
		})
	}

	list.segueMount()

	return &list
}

// triggered when a core is selected in the file explorer of Load Core
func coreExplorerCb(path string) {
	err := core.Load(path)
	if err != nil {
		ntf.DisplayAndLog(ntf.Error, "Core", err.Error())
		return
	}
	ntf.DisplayAndLog(ntf.Success, "Core", "Core loaded: %s", filepath.Base(path))
}

// triggered when a game is selected in the file explorer of Load Game
func gameExplorerCb(path string) {
	err := core.LoadGame(path)
	if err != nil {
		ntf.DisplayAndLog(ntf.Error, "Core", err.Error())
		return
	}
	menu.WarpToQuickMenu()
	state.Global.MenuActive = false
}

func (main *sceneMain) Entry() *entry {
	return &main.entry
}

func (main *sceneMain) segueMount() {
	genericSegueMount(&main.entry)
}

func (main *sceneMain) segueBack() {
	genericAnimate(&main.entry)
}

func (main *sceneMain) segueNext() {
	genericSegueNext(&main.entry)
}

func (main *sceneMain) update(dt float32) {
	genericInput(&main.entry, dt)
}

func (main *sceneMain) render() {
	genericRender(&main.entry)
}

func (main *sceneMain) drawHintBar() {
	genericDrawHintBar()
}
