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
	"github.com/libretro/ludo/utils"
)

type scenePlaylist struct {
	entry
}

func buildPlaylist(path string) Scene {
	var list scenePlaylist
	list.label = utils.FileName(path)

	for _, game := range playlists.Playlists[path] {
		game := game // needed for callbackOK
		strippedName, tags := extractTags(game.Name)
		if strings.Contains(game.Name, "Disc") {
			re := regexp.MustCompile(`\((Disc [1-9]?)\)`)
			match := re.FindStringSubmatch(game.Name)
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

func loadPlaylistEntry(list *scenePlaylist, playlist string, game playlists.Game) {
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
		list.segueNext()
		menu.Push(buildQuickMenu())
		menu.tweens.FastForward() // position the elements without animating
		state.MenuActive = false
	} else {
		list.segueNext()
		menu.Push(buildQuickMenu())
	}
}

func removePlaylistGame(s []playlists.Game, game playlists.Game) []playlists.Game {
	l := []playlists.Game{}
	for _, g := range s {
		if g.Path != game.Path {
			l = append(l, g)
		}
	}
	return l
}

func removePlaylistEntry(s []entry, game playlists.Game) []entry {
	l := []entry{}
	for _, g := range s {
		if g.path != game.Path {
			l = append(l, g)
		}
	}

	return l
}

func deletePlaylistEntry(list *scenePlaylist, path string, game playlists.Game) {
	playlists.Playlists[path] = removePlaylistGame(playlists.Playlists[path], game)
	playlists.Save(path)
	refreshTabs()
	list.children = removePlaylistEntry(list.children, game)

	if len(playlists.Playlists[path]) == 0 {
		list.children = append(list.children, entry{
			label: "Empty playlist",
			icon:  "subsetting",
		})
	}

	if list.ptr >= len(list.children) {
		list.ptr = len(list.children) - 1
	}

	buildIndexes(&list.entry)
	genericAnimate(&list.entry)
}

// Generic stuff
func (s *scenePlaylist) Entry() *entry {
	return &s.entry
}

func (s *scenePlaylist) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *scenePlaylist) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *scenePlaylist) segueBack() {
	genericAnimate(&s.entry)
}

func (s *scenePlaylist) update(dt float32) {
	genericInput(&s.entry, dt)
}

// Override rendering
func (s *scenePlaylist) render() {
	list := &s.entry

	_, h := menu.GetFramebufferSize()

	thumbnailDrawCursor(list)

	menu.ScissorStart(int32(510*menu.ratio), 0, int32(1310*menu.ratio), int32(h))

	for i, e := range list.children {
		if e.yp < -0.1 || e.yp > 1.1 {
			freeThumbnail(list, i)
			continue
		}

		fontOffset := 64 * 0.7 * menu.ratio * 0.3

		if e.labelAlpha > 0 {
			drawThumbnail(
				list, i,
				list.label, e.gameName,
				680*menu.ratio-85*e.scale*menu.ratio,
				float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
				170*menu.ratio, 128*menu.ratio,
				e.scale, white.Alpha(e.iconAlpha),
			)
			menu.DrawBorder(
				680*menu.ratio-85*e.scale*menu.ratio,
				float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
				170*menu.ratio*e.scale, 128*menu.ratio*e.scale, 0.02/e.scale,
				textColor.Alpha(e.iconAlpha))
			if e.path == state.GamePath && e.path != "" {
				menu.DrawCircle(
					680*menu.ratio,
					float32(h)*e.yp-14*menu.ratio+fontOffset,
					90*menu.ratio*e.scale,
					black.Alpha(e.iconAlpha))
				menu.DrawImage(menu.icons["resume"],
					680*menu.ratio-25*e.scale*menu.ratio,
					float32(h)*e.yp-14*menu.ratio-25*e.scale*menu.ratio+fontOffset,
					50*menu.ratio, 50*menu.ratio,
					e.scale, white.Alpha(e.iconAlpha))
			}

			menu.Font.SetColor(textColor.Alpha(e.labelAlpha))
			stack := 840 * menu.ratio
			menu.Font.Printf(
				840*menu.ratio,
				float32(h)*e.yp+fontOffset,
				0.5*menu.ratio, e.label)
			stack += float32(int(menu.Font.Width(0.5*menu.ratio, e.label)))
			stack += 10

			for _, tag := range e.tags {
				if _, ok := menu.icons[tag]; ok {
					stack += 20
					menu.DrawImage(
						menu.icons[tag],
						stack, float32(h)*e.yp-22*menu.ratio,
						60*menu.ratio, 44*menu.ratio, 1.0, white.Alpha(e.tagAlpha))
					menu.DrawBorder(stack, float32(h)*e.yp-22*menu.ratio,
						60*menu.ratio, 44*menu.ratio, 0.05/menu.ratio, black.Alpha(e.tagAlpha/4))
					stack += 60 * menu.ratio
				}
			}
		}
	}

	menu.ScissorEnd()
}

func (s *scenePlaylist) drawHintBar() {
	w, h := menu.GetFramebufferSize()
	menu.DrawRect(0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 0, lightGrey)

	_, upDown, _, a, b, x, _, _, _, guide := hintIcons()

	var stack float32
	if state.CoreRunning {
		stackHint(&stack, guide, "RESUME", h)
	}
	stackHint(&stack, upDown, "NAVIGATE", h)

	if settings.Current.SwapConfirm {
		stackHint(&stack, b, "RUN", h)
        stackHint(&stack, a, "BACK", h)
	} else {
		stackHint(&stack, a, "RUN", h)
        stackHint(&stack, b, "BACK", h)
	}

	list := menu.stack[len(menu.stack)-1].Entry()
	if list.children[list.ptr].callbackX != nil {
		stackHint(&stack, x, "DELETE", h)
	}
}
