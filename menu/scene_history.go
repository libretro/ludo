package menu

import (
	"os"
	"path/filepath"

	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/history"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/state"
)

type sceneHistory struct {
	entry
}

func buildHistory() Scene {
	var list sceneHistory
	list.label = "History"

	history.Load()
	for _, game := range history.List {
		game := game // needed for callbackOK
		strippedName, tags := extractTags(game.Name)
		list.children = append(list.children, entry{
			label:      strippedName,
			subLabel:   game.System,
			gameName:   game.Name,
			path:       game.Path,
			system:     game.System,
			tags:       tags,
			callbackOK: func() { loadHistoryEntry(&list, game) },
		})
	}

	if len(history.List) == 0 {
		list.children = append(list.children, entry{
			label: "Empty history",
			icon:  "subsetting",
		})
	}

	list.segueMount()
	return &list
}

func loadHistoryEntry(list Scene, game history.Game) {
	if _, err := os.Stat(game.Path); os.IsNotExist(err) {
		ntf.DisplayAndLog(ntf.Error, "Menu", "Game not found.")
		return
	}
	corePath := game.CorePath
	if _, err := os.Stat(corePath); os.IsNotExist(err) {
		ntf.DisplayAndLog(ntf.Error, "Menu", "Core not found: %s", filepath.Base(corePath))
		return
	}
	if state.Global.CorePath != corePath {
		err := core.Load(corePath)
		if err != nil {
			ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			return
		}
	}
	if state.Global.GamePath != game.Path {
		err := core.LoadGame(game.Path)
		if err != nil {
			ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			return
		}
		history.Push(history.Game{
			Path:     game.Path,
			Name:     game.Name,
			System:   game.System,
			CorePath: corePath,
		})
		list.segueNext()
		menu.Push(buildQuickMenu())
		menu.tweens.FastForward() // position the elements without animating
		state.Global.MenuActive = false
	} else {
		list.segueNext()
		menu.Push(buildQuickMenu())
	}
}

// Generic stuff
func (s *sceneHistory) Entry() *entry {
	return &s.entry
}

func (s *sceneHistory) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneHistory) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneHistory) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneHistory) update(dt float32) {
	genericInput(&s.entry, dt)
}

// Override rendering
func (s *sceneHistory) render() {
	list := &s.entry

	_, h := vid.Window.GetFramebufferSize()

	thumbnailDrawCursor(list)

	vid.ScissorStart(int32(510*menu.ratio), 0, int32(1310*menu.ratio), int32(h))

	for i, e := range list.children {
		if e.yp < -0.1 || e.yp > 1.1 {
			freeThumbnail(list, i)
			continue
		}

		fontOffset := 64 * 0.7 * menu.ratio * 0.3

		if e.labelAlpha > 0 {
			drawThumbnail(
				list, i,
				e.system, e.gameName,
				680*menu.ratio-85*e.scale*menu.ratio,
				float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
				170*menu.ratio, 128*menu.ratio,
				e.scale, white.Alpha(e.iconAlpha),
			)
			vid.DrawBorder(
				680*menu.ratio-85*e.scale*menu.ratio,
				float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
				170*menu.ratio*e.scale, 128*menu.ratio*e.scale, 0.02/e.scale,
				textColor.Alpha(e.iconAlpha))
			if e.path == state.Global.GamePath && e.path != "" {
				vid.DrawCircle(
					680*menu.ratio,
					float32(h)*e.yp-14*menu.ratio+fontOffset,
					90*menu.ratio*e.scale,
					black.Alpha(e.iconAlpha))
				vid.DrawImage(menu.icons["resume"],
					680*menu.ratio-25*e.scale*menu.ratio,
					float32(h)*e.yp-14*menu.ratio-25*e.scale*menu.ratio+fontOffset,
					50*menu.ratio, 50*menu.ratio,
					e.scale, white.Alpha(e.iconAlpha))
			}

			// Offset on Y to vertically center label + sublabel if there is a sublabel
			slOffset := float32(0)
			if e.subLabel != "" {
				slOffset = 30 * menu.ratio * e.subLabelAlpha
			}

			vid.Font.SetColor(textColor.Alpha(e.labelAlpha))
			stack := 840 * menu.ratio
			vid.Font.Printf(
				840*menu.ratio,
				float32(h)*e.yp+fontOffset-slOffset,
				0.5*menu.ratio, e.label)
			stack += float32(int(vid.Font.Width(0.5*menu.ratio, e.label)))
			stack += 10

			for _, tag := range e.tags {
				if _, ok := menu.icons[tag]; ok {
					stack += 20
					vid.DrawImage(
						menu.icons[tag],
						stack, float32(h)*e.yp-22*menu.ratio-slOffset,
						60*menu.ratio, 44*menu.ratio, 1.0, white.Alpha(e.tagAlpha))
					vid.DrawBorder(stack, float32(h)*e.yp-22*menu.ratio-slOffset,
						60*menu.ratio, 44*menu.ratio, 0.05/menu.ratio, black.Alpha(e.tagAlpha/4))
					stack += 60 * menu.ratio
				}
			}

			vid.Font.SetColor(mediumGrey.Alpha(e.subLabelAlpha))
			vid.Font.Printf(
				840*menu.ratio,
				float32(h)*e.yp+fontOffset+60*menu.ratio-slOffset,
				0.5*menu.ratio, e.subLabel)
		}
	}

	vid.ScissorEnd()
}

func (s *sceneHistory) drawHintBar() {
	w, h := vid.Window.GetFramebufferSize()
	vid.DrawRect(0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 0, lightGrey)

	_, upDown, _, a, b, _, _, _, _, guide := hintIcons()

	var stack float32
	if state.Global.CoreRunning {
		stackHint(&stack, guide, "RESUME", h)
	}
	stackHint(&stack, upDown, "NAVIGATE", h)
	stackHint(&stack, b, "BACK", h)
	stackHint(&stack, a, "RUN", h)
}
