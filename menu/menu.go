// Package menu is the graphical interface allowing to browse games, launch
// games, configure settings, and display a contextual menu to interract with
// the running game.
package menu

import (
	"math"
	"path/filepath"

	"github.com/libretro/ludo/options"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
	"github.com/libretro/ludo/video"

	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

var vid *video.Video
var opts *options.Options

// entry is a menu entry. It can also represent a scene.
// The menu data is a tree of entries.
type entry struct {
	yp, scale       float32
	width           float32
	label, subLabel string
	path            string // full path of the rom linked to the entry
	labelAlpha      float32
	icon            string
	iconAlpha       float32
	tagAlpha        float32
	callbackOK      func() // callback executed when user presses OK
	value           func() interface{}
	stringValue     func() string
	widget          func(*entry) // widget draw callback used in settings
	incr            func(int)    // increment callback used in settings
	tags            []string     // flags extracted from game title
	thumbnail       uint32       // thumbnail texture id
	gameName        string       // title of the game in db, used for thumbnails
	cursor          struct {
		alpha float32
		yp    float32
	}
	children []entry // children entries
	ptr      int     // index of the active child
}

// Scene represents a page of the UI
// A scene is typically an entry displaying its own children
// A segue is a smooth transition between two scenes.
type Scene interface {
	segueMount()
	segueNext()
	segueBack()
	update()
	render()
	drawHintBar()
	Entry() *entry
}

// Menu is a type holding the menu state, the stack of scenes, tweens, etc
type Menu struct {
	stack         []Scene
	icons         map[string]uint32
	inputCooldown int
	tweens        map[*float32]*gween.Tween
	scroll        float32
	ratio         float32
	t             float64
}

var menu *Menu

// updateTweens loops over the animation queue and updade them so we can see progress
func updateTweens(dt float32) {
	for e, t := range menu.tweens {
		var finished bool
		*e, finished = t.Update(dt)
		if finished {
			delete(menu.tweens, e)
		}
	}
}

// Render takes care of rendering the menu
func Render() {
	menu.t += 0.1
	w, _ := vid.Window.GetFramebufferSize()
	menu.ratio = float32(w) / 1920

	vid.FullViewport()

	updateTweens(1.0 / 60.0)

	currentScreenIndex := len(menu.stack) - 1
	for i := 0; i <= currentScreenIndex+1; i++ {
		if i < 0 || i > currentScreenIndex {
			continue
		}

		menu := menu.stack[i]
		menu.render()
	}
	menu.stack[currentScreenIndex].drawHintBar()
}

func genericDrawHintBar() {
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
	vid.Font.Printf(stack, float32(h)-23*menu.ratio, 0.5*menu.ratio, "OK")
	stack += vid.Font.Width(0.5*menu.ratio, "OK")
}

// genericSegueMount is the smooth transition of the menu entries first appearance
func genericSegueMount(list *entry) {
	for i := range list.children {
		e := &list.children[i]

		if i == list.ptr {
			e.yp = 0.5 + 0.3
			e.labelAlpha = 0
			e.iconAlpha = 0
			e.tagAlpha = 1
			e.scale = 1.5
		} else if i < list.ptr {
			e.yp = 0.4 + 0.3 + 0.08*float32(i-list.ptr)
			e.labelAlpha = 0
			e.iconAlpha = 0
			e.tagAlpha = 0
			e.scale = 0.5
		} else if i > list.ptr {
			e.yp = 0.6 + 0.3 + 0.08*float32(i-list.ptr)
			e.labelAlpha = 0
			e.iconAlpha = 0
			e.tagAlpha = 0
			e.scale = 0.5
		}
	}
	list.cursor.alpha = 0
	list.cursor.yp = 0.5 + 0.3

	genericAnimate(list)
}

// genericAnimate is the generic animation of entries when the user scrolls up or down
func genericAnimate(list *entry) {
	for i := range list.children {
		e := &list.children[i]

		var yp, labelAlpha, iconAlpha, tagAlpha, scale float32
		if i == list.ptr {
			yp = 0.5
			labelAlpha = 1
			iconAlpha = 1
			tagAlpha = 1
			scale = 1.5
		} else if i < list.ptr {
			yp = 0.4 + 0.08*float32(i-list.ptr)
			labelAlpha = 0.75
			iconAlpha = 0.75
			tagAlpha = 0
			scale = 0.5
		} else if i > list.ptr {
			yp = 0.6 + 0.08*float32(i-list.ptr)
			labelAlpha = 0.75
			iconAlpha = 0.75
			tagAlpha = 0
			scale = 0.5
		}

		menu.tweens[&e.yp] = gween.New(e.yp, yp, 0.15, ease.OutSine)
		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, labelAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, iconAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.tagAlpha] = gween.New(e.tagAlpha, tagAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.scale] = gween.New(e.scale, scale, 0.15, ease.OutSine)
	}
	menu.tweens[&list.cursor.alpha] = gween.New(list.cursor.alpha, 0.1, 0.15, ease.OutSine)
	menu.tweens[&list.cursor.yp] = gween.New(list.cursor.yp, 0.5, 0.15, ease.OutSine)
}

// genericSegueNext is a smooth transition that fades out the current list
// to leave room for the next list to appear
func genericSegueNext(list *entry) {
	for i := range list.children {
		e := &list.children[i]

		var yp, labelAlpha, iconAlpha, tagAlpha, scale float32
		if i == list.ptr {
			yp = 0.5 - 0.3
			labelAlpha = 0
			iconAlpha = 0
			tagAlpha = 1
			scale = 1.5
		} else if i < list.ptr {
			yp = 0.4 - 0.3 + 0.08*float32(i-list.ptr)
			labelAlpha = 0
			iconAlpha = 0
			tagAlpha = 0
			scale = 0.5
		} else if i > list.ptr {
			yp = 0.6 - 0.3 + 0.08*float32(i-list.ptr)
			labelAlpha = 0
			iconAlpha = 0
			tagAlpha = 0
			scale = 0.5
		}

		menu.tweens[&e.yp] = gween.New(e.yp, yp, 0.15, ease.OutSine)
		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, labelAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, iconAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.tagAlpha] = gween.New(e.tagAlpha, tagAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.scale] = gween.New(e.scale, scale, 0.15, ease.OutSine)
	}
	menu.tweens[&list.cursor.alpha] = gween.New(list.cursor.alpha, 0, 0.15, ease.OutSine)
	menu.tweens[&list.cursor.yp] = gween.New(list.cursor.yp, 0.5-0.3, 0.15, ease.OutSine)
}

// drawCursor draws the blinking rectangular background of the active menu entry
func drawCursor(list *entry) {
	w, h := vid.Window.GetFramebufferSize()
	alpha := list.cursor.alpha - float32(math.Cos(menu.t))*0.025 - 0.025
	c := video.Color{R: 0.25, G: 0.25, B: 0.25, A: alpha}
	if state.Global.CoreRunning {
		c = video.Color{R: 1, G: 1, B: 1, A: alpha}
	}
	vid.DrawRect(
		550*menu.ratio, float32(h)*list.cursor.yp-50*menu.ratio,
		float32(w)-630*menu.ratio, 100*menu.ratio, 1.0, c)
	vid.DrawBorder(
		550*menu.ratio, float32(h)*list.cursor.yp-50*menu.ratio,
		float32(w)-630*menu.ratio, 100*menu.ratio, 0.02,
		video.Color{R: c.R, G: c.G, B: c.B, A: alpha * 3})
}

// genericRender renders a vertical list of menu entries
// It also display values of settings if we are displaying a settings scene
func genericRender(list *entry) {
	w, h := vid.Window.GetFramebufferSize()

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
			610*menu.ratio-64*0.5*menu.ratio,
			float32(h)*e.yp-14*menu.ratio-64*0.5*menu.ratio+fontOffset,
			128*menu.ratio, 128*menu.ratio,
			0.5, color)

		if e.labelAlpha > 0 {
			vid.Font.SetColor(color.R, color.G, color.B, e.labelAlpha)
			vid.Font.Printf(
				670*menu.ratio,
				float32(h)*e.yp+fontOffset,
				0.6*menu.ratio, e.label)

			if e.widget != nil {
				e.widget(&e)
			} else if e.stringValue != nil {
				lw := vid.Font.Width(0.6*menu.ratio, e.stringValue())
				vid.Font.Printf(
					float32(w)-lw-128*menu.ratio,
					float32(h)*e.yp+fontOffset,
					0.6*menu.ratio, e.stringValue())
			}
		}
	}
}

// ContextReset uploads the UI images to the GPU.
// It should be called after each time the window is recreated.
func (menu *Menu) ContextReset() {
	paths, _ := filepath.Glob("assets/*.png")
	for _, path := range paths {
		path := path
		filename := utils.Filename(path)
		menu.icons[filename] = video.NewImage("assets/" + filename + ".png")
	}

	paths, _ = filepath.Glob("assets/flags/*.png")
	for _, path := range paths {
		path := path
		filename := utils.Filename(path)
		menu.icons[filename] = video.NewImage("assets/flags/" + filename + ".png")
	}

	currentScreenIndex := len(menu.stack) - 1
	curList := menu.stack[currentScreenIndex].Entry()
	for i := range curList.children {
		curList.children[i].thumbnail = 0
	}
}

// fastForwardTweens finishes all the current animations in the queue.
func fastForwardTweens() {
	updateTweens(10)
}

// UpdateOptions updates the menu with the core options of the newly loaded
// libretro core.
func (menu *Menu) UpdateOptions(o *options.Options) {
	opts = o
}

// WarpToQuickMenu loads the contextual menu for games that are launched from
// the command line interface.
func (menu *Menu) WarpToQuickMenu() {
	menu.stack = []Scene{}
	menu.stack = append(menu.stack, buildTabs())
	menu.stack[0].segueNext()
	menu.stack = append(menu.stack, buildMainMenu())
	menu.stack[1].segueNext()
	menu.stack = append(menu.stack, buildQuickMenu())
	fastForwardTweens()
}

// Init initializes the menu.
// If a game is already running, it will warp the user to the quick menu.
// If not, it will display the menu tabs.
func Init(v *video.Video) *Menu {
	vid = v

	w, _ := v.Window.GetFramebufferSize()
	menu = &Menu{}
	menu.stack = []Scene{}
	menu.tweens = make(map[*float32]*gween.Tween)
	menu.ratio = float32(w) / 1920
	menu.icons = map[string]uint32{}

	menu.stack = append(menu.stack, buildTabs())

	return menu
}
