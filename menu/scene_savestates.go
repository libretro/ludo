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
	list.entryHeight = 160

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
	w, h := vid.Window.GetFramebufferSize()

	vid.BoldFont.SetColor(blue.Alpha(list.cursor.alpha))
	vid.BoldFont.Printf(
		360*menu.ratio,
		230*menu.ratio,
		0.5*menu.ratio, list.label)

	vid.DrawRect(
		360*menu.ratio,
		270*menu.ratio,
		float32(w)-720*menu.ratio,
		2*menu.ratio,
		0, lightGrey,
	)

	vid.ScissorStart(
		int32(360*menu.ratio-8*menu.ratio), 0,
		int32(float32(w)-720*menu.ratio+16*menu.ratio), int32(h)-int32(272*menu.ratio))

	fontOffset := 12 * menu.ratio

	for i, e := range list.children {
		vid.DrawRect(
			360*menu.ratio,
			(270+32-1)*menu.ratio+list.scroll*menu.ratio+list.entryHeight*float32(i)*menu.ratio+list.entryHeight*menu.ratio,
			float32(w)-720*menu.ratio,
			2*menu.ratio,
			0, lightGrey,
		)

		if i == list.ptr {
			genericDrawCursor(list, i)
		}

		y := (270+32)*menu.ratio +
			list.scroll*menu.ratio +
			list.entryHeight*float32(i)*menu.ratio +
			list.entryHeight/2*menu.ratio

		if e.labelAlpha > 0 {
			drawSavestateThumbnail(
				list, i,
				filepath.Join(settings.Current.ScreenshotsDirectory, utils.FileName(e.path)+".png"),
				480*menu.ratio-85*1*menu.ratio,
				y-64*menu.ratio,
				170*menu.ratio, 128*menu.ratio,
				1, white.Alpha(e.iconAlpha),
			)
			if i == 0 {
				vid.DrawImage(menu.icons["savestate"],
					480*menu.ratio-25*1*menu.ratio,
					y-64*menu.ratio,
					50*menu.ratio, 50*menu.ratio,
					1, 0, white.Alpha(e.iconAlpha))
			}

			vid.Font.SetColor(black.Alpha(e.labelAlpha))
			vid.Font.Printf(
				600*menu.ratio,
				y+fontOffset,
				0.5*menu.ratio, e.label)
		}
	}

	vid.ScissorEnd()
}

func (s *sceneSavestates) drawHintBar() {
	w, h := vid.Window.GetFramebufferSize()
	vid.DrawRect(0, float32(h)-88*menu.ratio, float32(w), 88*menu.ratio, 0, white)
	vid.DrawRect(0, float32(h)-88*menu.ratio, float32(w), 2*menu.ratio, 0, lightGrey)

	ptr := menu.stack[len(menu.stack)-1].Entry().ptr

	_, upDown, _, a, b, _, _, _, _, guide := hintIcons()

	lstack := float32(75) * menu.ratio
	rstack := float32(w) - 96*menu.ratio
	stackHintLeft(&lstack, upDown, "Navigate", h)
	if ptr == 0 {
		stackHintRight(&rstack, a, "Save", h)
	} else {
		stackHintRight(&rstack, a, "Load", h)
	}
	stackHintRight(&rstack, b, "Back", h)
	if state.Global.CoreRunning {
		stackHintRight(&rstack, guide, "Resume", h)
	}
}
