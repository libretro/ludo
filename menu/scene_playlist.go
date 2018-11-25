package menu

import (
	"os"
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

	for _, game := range playlists.Playlists[path] {
		game := game // needed for callbackOK
		strippedName, tags := extractTags(game.Name)
		list.children = append(list.children, entry{
			label:      strippedName,
			gameName:   game.Name,
			path:       game.Path,
			tags:       tags,
			icon:       utils.Filename(path) + "-content",
			callbackOK: func() { loadEntry(&list, list.label, game.Path) },
		})
	}
	list.segueMount()
	return &list
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

func loadEntry(list *screenPlaylist, playlist, gamePath string) {
	corePath, err := settings.CoreForPlaylist(playlist)
	if err != nil {
		notifications.DisplayAndLog("Menu", err.Error())
		return
	}
	if _, err := os.Stat(corePath); os.IsNotExist(err) {
		notifications.DisplayAndLog("Menu", "Core not found.")
		return
	}
	if state.Global.CorePath != corePath {
		core.Load(corePath)
	}
	if state.Global.GamePath != gamePath {
		err := core.LoadGame(gamePath)
		if err != nil {
			notifications.DisplayAndLog("Menu", err.Error())
			return
		}
		list.segueNext()
		menu.stack = append(menu.stack, buildQuickMenu())
		fastForwardTweens() // position the elements without animating
	} else {
		list.segueNext()
		menu.stack = append(menu.stack, buildQuickMenu())
	}
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

	for i, e := range list.children {
		if e.yp < -0.1 || e.yp > 1.1 {
			continue
		}

		fontOffset := 64 * 0.7 * menu.ratio * 0.3

		color := video.Color{R: 0, G: 0, B: 0, A: e.iconAlpha}
		if state.Global.CoreRunning {
			color = video.Color{R: 1, G: 1, B: 1, A: e.iconAlpha}
		}

		if e.labelAlpha > 0 {
			drawThumbnail(
				list, i,
				list.label, e.gameName,
				680*menu.ratio-85*e.scale*menu.ratio,
				float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
				170*menu.ratio, 128*menu.ratio,
				e.scale,
			)
			vid.DrawBorder(
				680*menu.ratio-85*e.scale*menu.ratio,
				float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
				170*menu.ratio*e.scale, 128*menu.ratio*e.scale, 0.02/e.scale,
				video.Color{R: color.R, G: color.G, B: color.B, A: 0.75})
			if e.path == state.Global.GamePath {
				vid.DrawImage(menu.icons["resume"],
					680*menu.ratio-64*e.scale*menu.ratio,
					float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
					128*menu.ratio, 128*menu.ratio,
					e.scale, video.Color{R: 1, G: 1, B: 1, A: 1})
			}

			vid.Font.SetColor(color.R, color.G, color.B, e.labelAlpha)
			stack := 840 * menu.ratio
			vid.Font.Printf(
				840*menu.ratio,
				float32(h)*e.yp+fontOffset,
				0.6*menu.ratio, e.label)
			stack += float32(int(vid.Font.Width(0.6*menu.ratio, e.label)))
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

func (s *screenPlaylist) drawHintBar() {
	w, h := vid.Window.GetFramebufferSize()
	c := video.Color{R: 0.25, G: 0.25, B: 0.25, A: 1}
	menu.ratio = float32(w) / 1920
	vid.DrawRect(0.0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 1.0, video.Color{R: 0.75, G: 0.75, B: 0.75, A: 1})
	vid.Font.SetColor(0.25, 0.25, 0.25, 1.0)

	stack := 30 * menu.ratio
	vid.DrawImage(menu.icons["key-up-down"], stack, float32(h)-70*menu.ratio, 70*menu.ratio, 70*menu.ratio, 1.0, c)
	stack += 70 * menu.ratio
	stack += 10 * menu.ratio
	vid.Font.Printf(stack, float32(h)-23*menu.ratio, 0.5*menu.ratio, "NAVIGATE")
	stack += vid.Font.Width(0.5*menu.ratio, "NAVIGATE")

	stack += 30 * menu.ratio
	vid.DrawImage(menu.icons["key-z"], stack, float32(h)-70*menu.ratio, 70*menu.ratio, 70*menu.ratio, 1.0, c)
	stack += 70 * menu.ratio
	stack += 10 * menu.ratio
	vid.Font.Printf(stack, float32(h)-23*menu.ratio, 0.5*menu.ratio, "BACK")
	stack += vid.Font.Width(0.5*menu.ratio, "BACK")

	stack += 30 * menu.ratio
	vid.DrawImage(menu.icons["key-x"], stack, float32(h)-70*menu.ratio, 70*menu.ratio, 70*menu.ratio, 1.0, c)
	stack += 70 * menu.ratio
	stack += 10 * menu.ratio
	vid.Font.Printf(stack, float32(h)-23*menu.ratio, 0.5*menu.ratio, "RUN")
	stack += vid.Font.Width(0.5*menu.ratio, "RUN")
}
