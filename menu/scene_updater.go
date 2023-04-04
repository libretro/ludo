package menu

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/libretro/ludo/core"
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

			if rel.Name[1:] == ludos.Version {
				list.children[0].label = "Up to date"
				list.children[0].icon = "subsetting"
				return
			}

			list.children[0].label = "Upgrade to " + rel.Name
			list.children[0].icon = "menu_saving"
			list.children[0].callbackOK = func() {
				asset := ludos.FilterAssets(rel.Assets)
				if asset == nil {
					ntf.DisplayAndLog(ntf.Error, "Menu", "No matching asset")
					return
				}
				go func() {
					path := filepath.Join(ludos.UpdatesDir, asset.Name)
					ludos.DownloadRelease(path, asset.BrowserDownloadURL)
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
	if ludos.IsDownloading() {
		s.children[0].label = fmt.Sprintf(
			"Downloading update %.0f%%%%", ludos.GetProgress()*100)
		s.children[0].icon = "reload"
		s.children[0].callbackOK = nil
	} else if ludos.IsDone() {
		s.children[0].label = "Reboot and upgrade"
		s.children[0].icon = "reload"
		s.children[0].callbackOK = func() {
			cmd := exec.Command("/usr/sbin/shutdown", "-r", "now")
			core.UnloadGame()
			err := cmd.Run()
			if err != nil {
				ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			}
		}
	}

	genericInput(&s.entry, dt)
}

func (s *sceneUpdater) render() {
	genericRender(&s.entry)
}

func (s *sceneUpdater) drawHintBar() {
	genericDrawHintBar()
}
