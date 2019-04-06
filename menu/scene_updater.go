package menu

import (
	"path/filepath"

	"github.com/libretro/ludo/ludos"
	ntf "github.com/libretro/ludo/notifications"
)

type sceneUpdater struct {
	entry
}

func buildUpdater() Scene {
	var list sceneUpdater
	list.label = "Updater Menu"

	list.children = append(list.children, entry{
		label: "Checking updates",
		icon:  "reload",
	})

	list.segueMount()

	go func() {
		rels, err := ludos.GetReleases()
		if err != nil {
			ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			return
		}

		if len(*rels) > 0 {
			rel := (*rels)[0]
			list.children[0].label = "Upgrate to " + rel.Name
			list.children[0].icon = "menu_saving"
			list.children[0].callbackOK = func() {
				asset := ludos.FilterAssets(rel.Assets)
				if asset == nil {
					ntf.DisplayAndLog(ntf.Error, "Menu", "No matching asset")
					return
				}
				go func() {
					path := filepath.Join(ludos.UpdatesDir, asset.Name)
					ludos.DownloadRelease(asset.Name, path, asset.BrowserDownloadURL)
				}()
			}
		} else {
			list.children[0].label = "No updates found"
			list.children[0].icon = "menu_exit"
		}
	}()

	return &list
}

func (s *sceneUpdater) Entry() *entry {
	return &s.entry
}

func (s *sceneUpdater) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneUpdater) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneUpdater) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneUpdater) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *sceneUpdater) render() {
	genericRender(&s.entry)
}

func (s *sceneUpdater) drawHintBar() {
	genericDrawHintBar()
}
