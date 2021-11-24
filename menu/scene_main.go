package menu

import (
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/dat"
	"github.com/libretro/ludo/history"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/scanner"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
)

type sceneMain struct {
	entry
}

var prettyCoreNames = map[string]string{
	"atari800_libretro":            "Atari 800 (Atari 5200)",
	"bluemsx_libretro":             "BlueMSX (MSX)",
	"fbneo_libretro":               "FBNeo (Arcade)",
	"fceumm_libretro":              "Fceumm (NES)",
	"gambatte_libretro":            "Gambatte (Game Boy)",
	"gearsystem_libretro":          "GearSystem (Master System)",
	"genesis_plus_gx_libretro":     "Genesis Plus GX (Genesis, Master System, Sega CD)",
	"handy_libretro":               "Handy (Atari Lynx)",
	"lutro_libretro":               "Lutro (LÖVE)",
	"mednafen_ngp_libretro":        "Beetle NeoGeo Pocket",
	"mednafen_pce_fast_libretro":   "Beetle PC-Engine Fast",
	"mednafen_pce_libretro":        "Beetle PC-Engine",
	"mednafen_pcfx_libretro":       "Beetle PC-FX",
	"mednafen_psx_libretro":        "Beetle PlayStation",
	"mednafen_saturn_libretro":     "Beetle Saturn",
	"mednafen_supergrafx_libretro": "Beetle SuperGrafx",
	"mednafen_vb_libretro":         "Beetle VirtualBoy",
	"mednafen_wswan_libretro":      "Beetle WonderSwan",
	"melonds_libretro":             "MelonDS (Nintendo DS)",
	"mgba_libretro":                "mGBA (Game Boy Advance)",
	"np2kai_libretro":              "NP2Kai (PC-98)",
	"o2em_libretro":                "O2EM (Odyssey²)",
	"pcsx_rearmed_libretro":        "PCSX Rearmed (PLayStation)",
	"picodrive_libretro":           "PicoDrive (Genesis, 32X)",
	"pokemini_libretro":            "PokeMini (Pokemon Mini)",
	"prosystem_libretro":           "Prosystem (Atari 7800)",
	"sameboy_libretro":             "SameBoy (Game Boy)",
	"snes9x_libretro":              "Snes9x (SNES)",
	"stella2014_libretro":          "Stella 2014 (Atari 2600)",
	"swanstation_libretro":         "SwanStation (PlayStation)",
	"vecx_libretro":                "VecX (Vectrex)",
	"virtualjaguar_libretro":       "Virtual Jaguar (Atari Jaguar)",
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
	list.label = "Manual Menu"

	usr, _ := user.Current()

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
	} else {
		scanner.ScanFile(path, func(game dat.Game) {
			name := game.Name
			if name == "" {
				name = utils.FileName(path)
			}
			history.Push(history.Game{
				Path:     path,
				Name:     name,
				System:   game.System,
				CorePath: state.CorePath,
			})
			history.Load()
			menu.WarpToQuickMenu()
		})
		state.MenuActive = false
	}
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
