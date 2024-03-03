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

var prettyCoreNames = map[string]string{
	"atari800_libretro":            "Atari - 5200 (Atari800)",
	"bluemsx_libretro":             "MSX/SVI/ColecoVision/SG-1000 (blueMSX)",
	"fbneo_libretro":               "Arcade (FinalBurn Neo)",
	"fceumm_libretro":              "Nintendo - NES / Famicom (FCEUmm)",
	"gambatte_libretro":            "Nintendo - Game Boy / Color (Gambatte)",
	"gearsystem_libretro":          "Sega - MS/GG/SG-1000 (Gearsystem)",
	"genesis_plus_gx_libretro":     "Sega - MS/GG/MD/CD (Genesis Plus GX)",
	"handy_libretro":               "Atari - Lynx (Handy)",
	"lutro_libretro":               "Lua Engine (Lutro)",
	"mednafen_ngp_libretro":        "SNK - Neo Geo Pocket / Color (Beetle NeoPop)",
	"mednafen_pce_fast_libretro":   "NEC - PC Engine / CD (Beetle PCE FAST)",
	"mednafen_pce_libretro":        "NEC - PC Engine / SuperGrafx / CD (Beetle PCE)",
	"mednafen_pcfx_libretro":       "NEC - PC-FX (Beetle PC-FX)",
	"mednafen_psx_libretro":        "Sony - PlayStation (Beetle PSX)",
	"mednafen_saturn_libretro":     "Sega - Saturn (Beetle Saturn)",
	"mednafen_supergrafx_libretro": "NEC - PC Engine SuperGrafx (Beetle SuperGrafx)",
	"mednafen_vb_libretro":         "Nintendo - Virtual Boy (Beetle VB)",
	"mednafen_wswan_libretro":      "Bandai - WonderSwan/Color (Beetle Cygne)",
	"melonds_libretro":             "Nintendo - DS (melonDS)",
	"mgba_libretro":                "Nintendo - Game Boy Advance (mGBA)",
	"np2kai_libretro":              "NEC - PC-98 (Neko Project II Kai)",
	"o2em_libretro":                "Magnavox - Odyssey2 / Phillips Videopac+ (O2EM)",
	"pcsx_rearmed_libretro":        "Sony - PlayStation (PCSX ReARMed)",
	"picodrive_libretro":           "Sega - MS/GG/MD/CD/32X (PicoDrive)",
	"pokemini_libretro":            "Nintendo - Pokemon Mini (PokeMini)",
	"prosystem_libretro":           "Atari - 7800 (ProSystem)",
	"sameboy_libretro":             "Nintendo - Game Boy / Color (SameBoy)",
	"snes9x_libretro":              "Nintendo - SNES / SFC (Snes9x - Current)",
	"stella2014_libretro":          "Atari - 2600 (Stella 2014)",
	"swanstation_libretro":         "Sony - PlayStation (SwanStation)",
	"vecx_libretro":                "GCE - Vectrex (vecx)",
	"virtualjaguar_libretro":       "Atari - Jaguar (Virtual Jaguar)",
}

func prettifyCoreName(in string) string {
	name, ok := prettyCoreNames[in]
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
				prettifyCoreName,
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
