package menu

import (
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/history"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
)

type sceneMain struct {
	entry
}

var coreNames = map[string]string{
	"atari800_libretro":            "Atari 800",
	"bluemsx_libretro":             "BlueMSX",
	"fbneo_libretro":               "FBNeo",
	"fceumm_libretro":              "Fceumm",
	"gambatte_libretro":            "Gambatte",
	"genesis_plus_gx_libretro":     "Genesis Plus GX",
	"handy_libretro":               "Handy",
	"lutro_libretro":               "Lutro",
	"mednafen_ngp_libretro":        "Beetle NeoGeo Pocket",
	"mednafen_pce_fast_libretro":   "Beetle PC-Engine Fast",
	"mednafen_pce_libretro":        "Beetle PC-Engine",
	"mednafen_pcfx_libretro":       "Beetle PC-FX",
	"mednafen_psx_libretro":        "Beetle PlayStation",
	"mednafen_saturn_libretro":     "Beetle Saturn",
	"mednafen_supergrafx_libretro": "Beetle SuperGrafx",
	"mednafen_vb_libretro":         "Beetle VirtualBoy",
	"mednafen_wswan_libretro":      "Beetle WonderSwan",
	"melonds_libretro":             "MelonDS",
	"mgba_libretro":                "mGBA",
	"np2kai_libretro":              "NP2Kai",
	"o2em_libretro":                "O2EM",
	"pcsx_rearmed_libretro":        "PCSX Rearmed",
	"picodrive_libretro":           "PicoDrive",
	"pokemini_libretro":            "PokeMini",
	"prosystem_libretro":           "Prosystem",
	"sameboy_libretro":             "SameBoy",
	"snes9x_libretro":              "Snes9x",
	"stella2014_libretro":          "Stella 2014",
	"swanstation_libretro":         "SwanStation",
	"vecx_libretro":                "VecX",
	"virtualjaguar_libretro":       "Virtual Jaguar",
}

func coreName(in string) string {
	name, ok := coreNames[in]
	if ok {
		return name
	}
	return in
}

func buildMainMenu() Scene {
	var list sceneMain
	list.label = "Main Menu"

	usr, _ := user.Current()

	if state.CoreRunning {
		list.children = append(list.children, entry{
			label: "Quick Menu",
			icon:  "subsetting",
			callbackOK: func() {
				list.segueNext()
				menu.Push(buildQuickMenu())
			},
		})
	}

	list.children = append(list.children, entry{
		label: "Load Core",
		icon:  "subsetting",
		callbackOK: func() {
			list.segueNext()
			menu.Push(buildExplorer(
				settings.Current.CoresDirectory,
				[]string{".dll", ".dylib", ".so"},
				coreExplorerCb,
				nil,
				coreName,
			))
		},
	})

	list.children = append(list.children, entry{
		label: "Load Game",
		icon:  "subsetting",
		callbackOK: func() {
			if state.Core != nil {
				list.segueNext()
				menu.Push(buildExplorer(
					usr.HomeDir,
					nil,
					gameExplorerCb,
					nil,
					nil,
				))
			} else {
				ntf.DisplayAndLog(ntf.Warning, "Menu", "Please load a core first.")
			}
		},
	})

	if state.LudOS {
		list.children = append(list.children, entry{
			label: "Updater",
			icon:  "subsetting",
			callbackOK: func() {
				list.segueNext()
				menu.Push(buildUpdater())
			},
		})

		list.children = append(list.children, entry{
			label: "Reboot",
			icon:  "subsetting",
			callbackOK: func() {
				askQuitConfirmation(func() { cleanReboot() })
			},
		})

		list.children = append(list.children, entry{
			label: "Shutdown",
			icon:  "subsetting",
			callbackOK: func() {
				askQuitConfirmation(func() { cleanShutdown() })
			},
		})
	} else {
		list.children = append(list.children, entry{
			label: "Quit",
			icon:  "subsetting",
			callbackOK: func() {
				askQuitConfirmation(func() {
					menu.SetShouldClose(true)
				})
			},
		})
	}

	list.segueMount()

	return &list
}

// triggered when a core is selected in the file explorer of Load Core
func coreExplorerCb(path string) {
	if err := core.Load(path); err != nil {
		ntf.DisplayAndLog(ntf.Error, "Core", err.Error())
		return
	}
	ntf.DisplayAndLog(ntf.Success, "Core", "Core loaded: %s", filepath.Base(path))
}

// triggered when a game is selected in the file explorer of Load Game
func gameExplorerCb(path string) {
	if err := core.LoadGame(path); err != nil {
		ntf.DisplayAndLog(ntf.Error, "Core", err.Error())
		return
	}
	history.Push(history.Game{
		Path:     path,
		Name:     utils.FileName(path),
		CorePath: state.CorePath,
	})
	menu.WarpToQuickMenu()
	state.MenuActive = false
}

// Shutdown the operating system
func cleanShutdown() {
	core.UnloadGame()
	if err := exec.Command("/usr/sbin/shutdown", "-P", "now").Run(); err != nil {
		ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
	}
}

// Reboots the operating system
func cleanReboot() {
	core.UnloadGame()
	if err := exec.Command("/usr/sbin/shutdown", "-r", "now").Run(); err != nil {
		ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
	}
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
