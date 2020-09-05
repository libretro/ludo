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

type sceneSavestates struct {
	entry
}

func buildSavestates() Scene {
	var list sceneSavestates
	list.label = "Savestates"

	list.children = append(list.children, entry{
		label: "Save State",
		icon:  "savestate",
		callbackOK: func() {
			name := utils.DatedName(state.Global.GamePath)
			err := vid.TakeScreenshot(name)
			if err != nil {
				ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			}
			err = savestates.Save(name)
			if err != nil {
				ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			} else {
				menu.stack[len(menu.stack)-1] = buildSavestates()
				menu.tweens.FastForward()
				ntf.DisplayAndLog(ntf.Success, "Menu", "State saved.")
			}
		},
	})

	gameName := utils.FileName(state.Global.GamePath)
	paths, _ := filepath.Glob(settings.Current.SavestatesDirectory + "/" + gameName + "@*.state")
	sort.Sort(sort.Reverse(sort.StringSlice(paths)))
	for _, path := range paths {
		path := path
		date := strings.Replace(utils.FileName(path), gameName+"@", "", 1)
		list.children = append(list.children, entry{
			label: "Load " + date,
			icon:  "loadstate",
			path:  path,
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

func (s *sceneSavestates) Entry() *entry {
	return &s.entry
}

func (s *sceneSavestates) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneSavestates) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneSavestates) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneSavestates) update(dt float32) {
	genericInput(&s.entry, dt)
}

// Override rendering
func (s *sceneSavestates) render() {
	list := &s.entry

	_, h := vid.Window.GetFramebufferSize()

	thumbnailDrawCursor(list)

	for i, e := range list.children {
		if e.yp < -0.1 || e.yp > 1.1 {
			continue
		}

		fontOffset := 64 * 0.7 * menu.ratio * 0.3

		if e.labelAlpha > 0 {
			drawSavestateThumbnail(
				list, i,
				filepath.Join(settings.Current.ScreenshotsDirectory, utils.FileName(e.path)+".png"),
				680*menu.ratio-85*e.scale*menu.ratio,
				float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
				170*menu.ratio, 128*menu.ratio,
				e.scale, textColor.Alpha(e.iconAlpha),
			)
			vid.DrawBorder(
				680*menu.ratio-85*e.scale*menu.ratio,
				float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
				170*menu.ratio*e.scale, 128*menu.ratio*e.scale, 0.02/e.scale,
				textColor.Alpha(e.iconAlpha))
			if i == 0 {
				vid.DrawImage(menu.icons["savestate"],
					680*menu.ratio-25*e.scale*menu.ratio,
					float32(h)*e.yp-14*menu.ratio-25*e.scale*menu.ratio+fontOffset,
					50*menu.ratio, 50*menu.ratio,
					e.scale, textColor.Alpha(e.iconAlpha))
			}

			vid.Font.SetColor(textColor.Alpha(e.labelAlpha))
			vid.Font.Printf(
				840*menu.ratio,
				float32(h)*e.yp+fontOffset,
				0.5*menu.ratio, e.label)
		}
	}
}

func (s *sceneSavestates) drawHintBar() {
	w, h := vid.Window.GetFramebufferSize()
	vid.DrawRect(0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 0, lightGrey)

	ptr := menu.stack[len(menu.stack)-1].Entry().ptr

	_, upDown, _, a, b, _, _, _, _, guide := hintIcons()

	var stack float32
	if state.Global.CoreRunning {
		stackHint(&stack, guide, "RESUME", h)
	}
	stackHint(&stack, upDown, "NAVIGATE", h)
	stackHint(&stack, b, "BACK", h)
	if ptr == 0 {
		stackHint(&stack, a, "SAVE", h)
	} else {
		stackHint(&stack, a, "LOAD", h)
	}
}
