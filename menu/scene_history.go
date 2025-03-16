package menu

import (
	"os"
	"path/filepath"

	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/history"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/state"
)

func loadHistoryEntry(list Scene, game history.Game) {
	if _, err := os.Stat(game.Path); os.IsNotExist(err) {
		ntf.DisplayAndLog(ntf.Error, "Menu", "Game not found.")
		return
	}
	corePath := game.CorePath
	if _, err := os.Stat(corePath); os.IsNotExist(err) {
		ntf.DisplayAndLog(ntf.Error, "Menu", "Core not found: %s", filepath.Base(corePath))
		return
	}
	if state.CorePath != corePath {
		err := core.Load(corePath)
		if err != nil {
			ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			return
		}
	}
	if state.GamePath != game.Path {
		err := core.LoadGame(game.Path)
		if err != nil {
			ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			return
		}
		history.Push(history.Game{
			Path:     game.Path,
			Name:     game.Name,
			System:   game.System,
			CorePath: corePath,
		})
		history.Load()
		menu.WarpToQuickMenu()
		state.MenuActive = false
	} else {
		menu.WarpToQuickMenu()
		state.MenuActive = false
	}
}
