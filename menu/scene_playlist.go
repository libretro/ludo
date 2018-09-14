package menu

import (
	"bufio"
	"os"

	"github.com/libretro/go-playthemall/core"
	"github.com/libretro/go-playthemall/notifications"
	"github.com/libretro/go-playthemall/settings"
	"github.com/libretro/go-playthemall/utils"
)

type screenPlaylist struct {
	entry
}

func buildPlaylist(path string) Scene {
	var list screenPlaylist
	list.label = utils.Filename(path)
	file, err := os.Open(path)
	if err != nil {
		notifications.DisplayAndLog("Menu", err.Error())
		list.segueMount()
		return &list
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for {
		more := scanner.Scan()
		if !more {
			break
		}
		gamePath := scanner.Text() // path
		scanner.Scan()
		name := scanner.Text()
		scanner.Scan() // unused
		scanner.Scan() // unused
		scanner.Scan() // CRC
		scanner.Scan() // lpl
		list.children = append(list.children, entry{
			label:      name,
			icon:       utils.Filename(path) + "-content",
			callbackOK: func() { loadEntry(list.label, gamePath) },
		})
	}
	list.segueMount()
	return &list
}

func loadEntry(playlist, gamePath string) {
	coreFullPath, err := settings.CoreForPlaylist(playlist)
	if err != nil {
		notifications.DisplayAndLog("Menu", err.Error())
		return
	}
	core.Load(coreFullPath)
	core.LoadGame(gamePath)
}

// Generic stuff
func (s *screenPlaylist) Entry() *entry {
	return &s.entry
}
func (s *screenPlaylist) segueMount() {
	genericSegueMount(&s.entry)
}
func (s *screenPlaylist) segueNext() {
	genericSegueNext(&s.entry)
}
func (s *screenPlaylist) segueBack() {
	genericAnimate(&s.entry)
}
func (s *screenPlaylist) update() {
	genericInput(&s.entry)
}
func (s *screenPlaylist) render() {
	genericRender(&s.entry)
}
