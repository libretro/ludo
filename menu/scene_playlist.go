package menu

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/history"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/playlists"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
)

/*
func buildPlaylist(path string) Scene {
	var list scenePlaylist
	list.label = utils.FileName(path)

	re := regexp.MustCompile(`\(([Dd]isc [1-9]?)\)`)
	for _, game := range playlists.Playlists[path] {
		game := game // needed for callbackOK
		strippedName, tags := extractTags(game.Name)
		if match := re.FindStringSubmatch(game.Name); len(match) >= 2 {
			strippedName = strippedName + " (" + match[1] + ")"
		}
		list.children = append(list.children, entry{
			label:      strippedName,
			gameName:   game.Name,
			path:       game.Path,
			tags:       tags,
			icon:       utils.FileName(path) + "-content",
			callbackOK: func() { loadPlaylistEntry(&list, list.label, game) },
			callbackX:  func() { askDeleteGameConfirmation(func() { deletePlaylistEntry(&list, path, game) }) },
		})
	}

	if len(playlists.Playlists[path]) == 0 {
		list.children = append(list.children, entry{
			label: "Empty playlist",
			icon:  "subsetting",
		})
	}

	buildIndexes(&list.entry)

	list.segueMount()
	return &list
}
*/

// Index first letters of entries to allow quick jump to the next or previous
// letter
func buildIndexes(list *entry) {
	var last byte
	for i := 0; i < len(list.children); i++ {
		char := list.children[i].label[0]
		if char != last {
			list.indexes = append(list.indexes, struct {
				Char  byte
				Index int
			}{char, i})
			last = char
		}
	}
}

func extractTags(name string) (string, []string) {
	re := regexp.MustCompile(`\(.*?\)`)
	pars := re.FindAllString(name, -1)
	var tags []string
	for _, par := range pars {
		name = strings.Replace(name, par, "", -1)
		par = strings.Replace(par, "(", "", -1)
		par = strings.Replace(par, ")", "", -1)
		results := strings.Split(par, ",")
		for _, result := range results {
			tags = append(tags, strings.TrimSpace(result))
		}
	}
	name = strings.TrimSpace(name)
	return name, tags
}

func loadPlaylistEntry(list Scene, playlist string, game playlists.Game) {
	if _, err := os.Stat(game.Path); os.IsNotExist(err) {
		ntf.DisplayAndLog(ntf.Error, "Menu", "Game not found.")
		return
	}
	corePath, err := settings.CoreForPlaylist(playlist)
	if err != nil {
		ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
		return
	}
	if _, err := os.Stat(corePath); os.IsNotExist(err) {
		ntf.DisplayAndLog(ntf.Error, "Menu", "Core not found: %s", filepath.Base(corePath))
		return
	}
	if state.CorePath != corePath {
		if err := core.Load(corePath); err != nil {
			ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			return
		}
	}
	if state.GamePath != game.Path {
		if err := core.LoadGame(game.Path); err != nil {
			ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			return
		}
		history.Push(history.Game{
			Path:     game.Path,
			Name:     game.Name,
			System:   playlist,
			CorePath: corePath,
		})
		history.Load()
		menu.WarpToQuickMenu()
		state.MenuActive = false
	} else {
		list.segueNext()
		menu.WarpToQuickMenu()
		state.MenuActive = false
	}
}
