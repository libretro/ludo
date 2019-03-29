package menu

import (
	"path/filepath"

	"github.com/libretro/ludo/deskenv"
	ntf "github.com/libretro/ludo/notifications"
)

type sceneUpdater struct {
	entry
}

func buildUpdater(releases []deskenv.GHRelease) Scene {
	var list sceneUpdater
	list.label = "Updater Menu"

	for _, rel := range releases {
		list.children = append(list.children, entry{
			label: rel.Name,
			icon:  "menu_saving",
			callbackOK: func() {
				asset := deskenv.FilterAssets(rel.Assets)
				if asset == nil {
					ntf.DisplayAndLog(ntf.Error, "Menu", "No matching asset")
					return
				}
				go func() {
					path := filepath.Join(deskenv.UpdatesDir, asset.Name)
					deskenv.DownloadRelease(asset.Name, path, asset.BrowserDownloadURL)
				}()
			},
		})
	}

	if len(list.children) == 0 {
		list.children = append(list.children, entry{
			label: "Empty",
			icon:  "subsetting",
		})
	}

	list.segueMount()

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
