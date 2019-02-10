package menu

import (
	"path/filepath"
	"sort"
	"strings"

	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/savestates"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
)

type screenSavestates struct {
	entry
}

func buildSavestates() Scene {
	var list screenSavestates
	list.label = "Savestates"

	list.children = append(list.children, entry{
		label: "Save State",
		icon:  "savestate",
		callbackOK: func() {
			err := savestates.Save()
			if err != nil {
				ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			} else {
				ntf.DisplayAndLog(ntf.Success, "Menu", "State saved.")
			}
		},
	})

	gameName := utils.Filename(state.Global.GamePath)
	paths, _ := filepath.Glob(settings.Current.SavestatesDirectory + "/" + gameName + "@*.state")
	sort.Sort(sort.Reverse(sort.StringSlice(paths)))
	for _, path := range paths {
		path := path
		date := strings.Replace(utils.Filename(path), gameName+"@", "", 1)
		list.children = append(list.children, entry{
			label: "Load " + date,
			icon:  "loadstate",
			callbackOK: func() {
				err := savestates.Load(path)
				if err != nil {
					ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
				} else {
					state.Global.MenuActive = false
					ntf.DisplayAndLog(ntf.Success, "Menu", "State loaded.")
				}
			},
		})
	}

	list.segueMount()

	return &list
}

func (s *screenSavestates) Entry() *entry {
	return &s.entry
}

func (s *screenSavestates) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *screenSavestates) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *screenSavestates) segueBack() {
	genericAnimate(&s.entry)
}

func (s *screenSavestates) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *screenSavestates) render() {
	genericRender(&s.entry)
}

func (s *screenSavestates) drawHintBar() {
	genericDrawHintBar()
}
