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
				slug := "Generic.x86_64"
				asset := deskenv.FilterAssets(slug, rel.Assets)
				if asset == nil {
					ntf.DisplayAndLog(ntf.Error, "Menu", "Couldn't find asset matching %s", slug)
					return
				}
				nid := ntf.DisplayAndLog(ntf.Info, "Menu", "Downloading %s", asset.Name)
				go func() {
					path := filepath.Join(deskenv.UpdatesDir, asset.Name)
					err := deskenv.DownloadRelease(path, asset.BrowserDownloadURL)
					if err != nil {
						ntf.Update(nid, ntf.Error, err.Error())
						return
					}
					ntf.Update(nid, ntf.Success, "Done downloading. You can now reboot your system.")
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
