package menu

import (
	"math"

	"github.com/libretro/ludo/state"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

// entry is a menu entry. It can also represent a scene.
// The menu data is a tree of entries.
type entry struct {
	alpha           float32
	scale           float32
	scroll          float32
	y               float32
	entryHeight     float32
	height          float32
	margin          float32
	label, subLabel string
	path            string // full path of the rom linked to the entry
	system          string // name of the game system
	labelAlpha      float32
	icon            string
	iconAlpha       float32
	tagAlpha        float32
	subLabelAlpha   float32
	borderAlpha     float32
	callbackOK      func() // callback executed when user presses OK
	value           func() interface{}
	stringValue     func() string
	widget          func(*entry, *entry, int) // widget draw callback used in settings
	incr            func(int)                 // increment callback used in settings
	tags            []string                  // flags extracted from game title
	thumbnail       uint32                    // thumbnail texture id
	gameName        string                    // title of the game in db, used for thumbnails
	children        []entry                   // children entries
	ptr             int                       // index of the active child
	indexes         []struct {
		Char  byte
		Index int
	}
}

// Scene represents a page of the UI
// A scene is typically an entry displaying its own children
// A segue is a smooth transition between two scenes.
type Scene interface {
	segueMount()
	segueNext()
	segueBack()
	update(dt float32)
	render()
	drawHintBar()
	Entry() *entry
}

// genericSegueMount is the smooth transition of the menu entries first appearance
func genericSegueMount(list *entry) {
	if list.entryHeight == 0 {
		list.entryHeight = 100
	}

	list.scroll = -float32(list.ptr) * list.entryHeight
	list.y = 300

	genericAnimate(list)
}

// genericAnimate is the generic animation of entries when the user scrolls up or down
func genericAnimate(list *entry) {
	for i := range list.children {
		e := &list.children[i]

		// performance improvement
		if math.Abs(float64(i-list.ptr)) > 8 {
			continue
		}

		var tagAlpha, subLabelAlpha float32
		if i == list.ptr {
			tagAlpha = 1
			subLabelAlpha = 1
		} else if i < list.ptr {
			tagAlpha = 0
			subLabelAlpha = 0
		} else if i > list.ptr {
			tagAlpha = 0
			subLabelAlpha = 0
		}

		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, 1, 0.15, ease.OutSine)
		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, 1, 0.15, ease.OutSine)
		menu.tweens[&e.tagAlpha] = gween.New(e.tagAlpha, tagAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.subLabelAlpha] = gween.New(e.subLabelAlpha, subLabelAlpha, 0.15, ease.OutSine)
	}

	margin := 32
	containerHeight := float32(1080 - 88 - 270 - margin*2)
	contentHeight := float32(len(list.children)) * list.entryHeight

	scroll := float32(0)
	if list.ptr >= 3 {
		scroll = -float32(list.ptr-3) * list.entryHeight
	}

	if -scroll > contentHeight-containerHeight {
		scroll = -(contentHeight - containerHeight)
	}

	if contentHeight < containerHeight {
		scroll = 0
	}

	menu.tweens[&list.scroll] = gween.New(list.scroll, scroll, 0.15, ease.OutSine)
	menu.tweens[&list.y] = gween.New(list.y, 0, 0.15, ease.OutSine)
	menu.tweens[&list.alpha] = gween.New(list.alpha, 1, 0.15, ease.OutSine)
}

// genericSegueNext is a smooth transition that fades out the current list
// to leave room for the next list to appear
func genericSegueNext(list *entry) {
	for i := range list.children {
		// performance improvement
		if math.Abs(float64(i-list.ptr)) > 8 {
			continue
		}

		e := &list.children[i]
		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, 0, 0.15, ease.OutSine)
		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, 0, 0.15, ease.OutSine)
		menu.tweens[&e.tagAlpha] = gween.New(e.tagAlpha, 0, 0.15, ease.OutSine)
		menu.tweens[&e.subLabelAlpha] = gween.New(e.subLabelAlpha, 0, 0.15, ease.OutSine)
	}
	menu.tweens[&list.alpha] = gween.New(list.alpha, 0, 0.15, ease.OutSine)
	menu.tweens[&list.y] = gween.New(list.y, -300, 0.15, ease.OutSine)
}

// genericDrawCursor draws the blinking rectangular background of the active
// menu entry
func genericDrawCursor(list *entry, i int) {
	w, _ := menu.Window.GetFramebufferSize()
	y := list.y + (270+32)*menu.ratio + list.scroll*menu.ratio + list.entryHeight*float32(i)*menu.ratio
	if menu.focus > 1 {
		blink := float32(math.Cos(menu.t))
		menu.DrawImage(
			menu.icons["selection"],
			360*menu.ratio-8*menu.ratio,
			y-8*menu.ratio,
			float32(w)-720*menu.ratio+16*menu.ratio,
			list.entryHeight*menu.ratio+16*menu.ratio,
			1, 0.15, white.Alpha(list.alpha-list.alpha*blink))
	}
	menu.DrawRect(
		360*menu.ratio,
		y,
		float32(w)-720*menu.ratio,
		list.entryHeight*menu.ratio, 0.1,
		white.Alpha(list.alpha))
}

// genericRender renders a vertical list of menu entries
// It also display values of settings if we are displaying a settings scene
func genericRender(list *entry) {
	w, h := menu.GetFramebufferSize()

	menu.BoldFont.SetColor(blue.Alpha(list.alpha))
	menu.BoldFont.Printf(
		360*menu.ratio,
		list.y*menu.ratio+230*menu.ratio,
		0.5*menu.ratio, list.label)

	menu.DrawRect(
		360*menu.ratio,
		list.y*menu.ratio+270*menu.ratio,
		float32(w)-720*menu.ratio,
		2*menu.ratio,
		0, lightGrey.Alpha(list.alpha),
	)

	menu.ScissorStart(
		int32(360*menu.ratio-8*menu.ratio), 0,
		int32(float32(w)-720*menu.ratio+16*menu.ratio), int32(h)-int32(272*menu.ratio+list.y*menu.ratio))

	fontOffset := 12 * menu.ratio

	for i, e := range list.children {
		// performance improvement
		if math.Abs(float64(i-list.ptr)) > 8 {
			continue
		}

		y := list.y*menu.ratio +
			(270+32)*menu.ratio +
			list.scroll*menu.ratio +
			list.entryHeight*float32(i)*menu.ratio +
			list.entryHeight/2*menu.ratio

		menu.DrawRect(
			360*menu.ratio,
			y-1*menu.ratio+list.entryHeight/2*menu.ratio,
			float32(w)-720*menu.ratio,
			2*menu.ratio,
			0, lightGrey.Alpha(e.iconAlpha),
		)

		if i == list.ptr {
			genericDrawCursor(list, i)
		}

		menu.DrawImage(menu.icons[e.icon],
			420*menu.ratio-64*0.35*menu.ratio,
			y-64*0.35*menu.ratio,
			128*menu.ratio, 128*menu.ratio,
			0.35, 0, black.Alpha(e.iconAlpha))

		if e.labelAlpha > 0 {
			menu.Font.SetColor(black.Alpha(e.labelAlpha))
			menu.Font.Printf(
				480*menu.ratio,
				y+fontOffset,
				0.5*menu.ratio, e.label)

			if e.widget != nil {
				e.widget(list, &e, i)
			} else if e.stringValue != nil {
				lw := menu.Font.Width(0.5*menu.ratio, e.stringValue())
				menu.Font.Printf(
					float32(w)-lw-400*menu.ratio,
					y+fontOffset,
					0.5*menu.ratio, e.stringValue())
			}
		}
	}

	menu.ScissorEnd()
}

func genericDrawHintBar() {
	w, h := menu.Window.GetFramebufferSize()
	menu.DrawRect(0, float32(h)-88*menu.ratio, float32(w), 88*menu.ratio, 0, white)
	menu.DrawRect(0, float32(h)-88*menu.ratio, float32(w), 2*menu.ratio, 0, lightGrey)

	_, upDown, _, a, b, _, _, _, _, guide := hintIcons()

	lstack := float32(75) * menu.ratio
	rstack := float32(w) - 96*menu.ratio
	stackHintLeft(&lstack, upDown, "Navigate", h)
	stackHintRight(&rstack, a, "Ok", h)
	stackHintRight(&rstack, b, "Back", h)
	if state.CoreRunning {
		stackHintRight(&rstack, guide, "Resume", h)
	}
}
