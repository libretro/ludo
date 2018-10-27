package menu

import (
	"regexp"
	"strings"

	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/playlists"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
	"github.com/libretro/ludo/video"
)

type screenPlaylist struct {
	entry
}

func buildPlaylist(path string) Scene {
	var list screenPlaylist
	list.label = utils.Filename(path)

	for _, game := range playlists.Playlists[list.label] {
		strippedName, tags := extractTags(game.Name)
		list.children = append(list.children, entry{
			label:      strippedName,
			tags:       tags,
			icon:       utils.Filename(path) + "-content",
			callbackOK: func() { loadEntry(list.label, game.Path) },
		})
	}
	list.segueMount()
	return &list
}

func extractTags(name string) (string, []string) {
	re := regexp.MustCompile("\\(.*?\\)")
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
	list := &s.entry

	_, h := vid.Window.GetFramebufferSize()

	drawCursor(list)

	for _, e := range list.children {
		if e.yp < -0.1 || e.yp > 1.1 {
			continue
		}

		fontOffset := 64 * 0.7 * menu.ratio * 0.3

		color := video.Color{R: 0, G: 0, B: 0, A: e.iconAlpha}
		if state.Global.CoreRunning {
			color = video.Color{R: 1, G: 1, B: 1, A: e.iconAlpha}
		}

		vid.DrawImage(menu.icons[e.icon],
			610*menu.ratio-64*e.scale*menu.ratio,
			float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
			128*menu.ratio, 128*menu.ratio,
			e.scale, color)

		if e.labelAlpha > 0 {
			vid.Font.SetColor(color.R, color.G, color.B, e.labelAlpha)
			stack := 670 * menu.ratio
			vid.Font.Printf(
				670*menu.ratio,
				float32(h)*e.yp+fontOffset,
				0.7*menu.ratio, e.label)
			stack += float32(int(vid.Font.Width(0.7*menu.ratio, e.label)))
			stack += 10

			for _, tag := range e.tags {
				stack += 20
				vid.DrawImage(
					menu.icons[tag],
					stack, float32(h)*e.yp-22*menu.ratio,
					60*menu.ratio, 44*menu.ratio, 1.0, video.Color{R: 1, G: 1, B: 1, A: e.tagAlpha})
				vid.DrawBorder(stack, float32(h)*e.yp-22*menu.ratio,
					60*menu.ratio, 44*menu.ratio, 0.05/menu.ratio, video.Color{R: 0, G: 0, B: 0, A: e.tagAlpha / 4})
				stack += 60 * menu.ratio
			}
		}
	}
}
