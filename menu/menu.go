// Package menu is the graphical interface allowing to browse games, launch
// games, configure settings, and display a contextual menu to interract with
// the running game.
package menu

import (
	"math"
	"os/user"
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
	labelAlpha      float32
	icon            string
	iconAlpha       float32
	ptr             int
	callbackOK      func()
	value           func() interface{}
	stringValue     func() string
	widget          func(*entry)
	incr            func(int)
	children        []entry
	cursor          struct {
		alpha float32
		yp    float32
	}
}

// Scene represents a page of the UI
// A Scene is typically an entry displaying its own children
type Scene interface {
	segueMount()
	segueNext()
	segueBack()
	update()
	render()
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
}

// genericSegueMount is the smooth transition of the menu entries first appearance
func genericSegueMount(list *entry) {
	for i := range list.children {
		e := &list.children[i]

		if i == list.ptr {
			e.yp = 0.5 + 0.3
			e.labelAlpha = 0
			e.iconAlpha = 0
			e.scale = 0.5
		} else if i < list.ptr {
			e.yp = 0.4 + 0.3 + 0.08*float32(i-list.ptr)
			e.labelAlpha = 0
			e.iconAlpha = 0
			e.scale = 0.5
		} else if i > list.ptr {
			e.yp = 0.6 + 0.3 + 0.08*float32(i-list.ptr)
			e.labelAlpha = 0
			e.iconAlpha = 0
			e.scale = 0.5
		}
		e.cursor.alpha = 0
		e.cursor.yp = 0.5 + 0.3
	}

	genericAnimate(list)
}

// genericAnimate is the generic animation of entries when the user scrolls up or down
func genericAnimate(list *entry) {
	for i := range list.children {
		e := &list.children[i]

		var yp, labelAlpha, iconAlpha, scale float32
		if i == list.ptr {
			yp = 0.5
			labelAlpha = 1.0
			iconAlpha = 1.0
			scale = 0.5
		} else if i < list.ptr {
			yp = 0.4 + 0.08*float32(i-list.ptr)
			labelAlpha = 0.5
			iconAlpha = 0.5
			scale = 0.5
		} else if i > list.ptr {
			yp = 0.6 + 0.08*float32(i-list.ptr)
			labelAlpha = 0.5
			iconAlpha = 0.5
			scale = 0.5
		}

		menu.tweens[&e.yp] = gween.New(e.yp, yp, 0.15, ease.OutSine)
		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, labelAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, iconAlpha, 0.15, ease.OutSine)
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

		var yp, la, a, s float32
		if i == list.ptr {
			yp = 0.5 - 0.3
			la = 0
			a = 0
			s = 1.0
		} else if i < list.ptr {
			yp = 0.4 - 0.3 + 0.08*float32(i-list.ptr)
			la = 0
			a = 0
			s = 0.5
		} else if i > list.ptr {
			yp = 0.6 - 0.3 + 0.08*float32(i-list.ptr)
			la = 0
			a = 0
			s = 0.5
		}

		menu.tweens[&e.yp] = gween.New(e.yp, yp, 0.15, ease.OutSine)
		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, la, 0.15, ease.OutSine)
		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, a, 0.15, ease.OutSine)
		menu.tweens[&e.scale] = gween.New(e.scale, s, 0.15, ease.OutSine)
	}
	menu.tweens[&list.cursor.alpha] = gween.New(list.cursor.alpha, 0, 0.15, ease.OutSine)
	menu.tweens[&list.cursor.yp] = gween.New(list.cursor.yp, 0.5+0.3, 0.15, ease.OutSine)
}

// drawCursor draws the blinking rectangular background of the active menu entry
func drawCursor(list *entry) {
	w, h := vid.Window.GetFramebufferSize()
	alpha := list.cursor.alpha - float32(math.Cos(menu.t))*0.025 - 0.025
	vid.DrawQuad(
		470*menu.ratio, float32(h)*list.cursor.yp-50*menu.ratio,
		float32(w)-100*menu.ratio, float32(h)*list.cursor.yp-50*menu.ratio,
		470*menu.ratio, float32(h)*list.cursor.yp+50*menu.ratio,
		float32(w)-100*menu.ratio, float32(h)*list.cursor.yp+50*menu.ratio,
		video.Color{R: 1, G: 1, B: 1, A: alpha},
	)
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

		vid.DrawImage(menu.icons[e.icon],
			520*menu.ratio-64*e.scale*menu.ratio,
			float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
			128*menu.ratio, 128*menu.ratio,
			e.scale, video.Color{R: 1, G: 1, B: 1, A: e.iconAlpha})

		if e.labelAlpha > 0 {

			vid.Font.SetColor(1, 1, 1, e.labelAlpha)
			vid.Font.Printf(
				600*menu.ratio,
				float32(h)*e.yp+fontOffset,
				0.7*menu.ratio, e.label)

			if e.widget != nil {
				e.widget(&e)
			} else if e.stringValue != nil {
				lw := vid.Font.Width(0.7*menu.ratio, e.stringValue())
				vid.Font.Printf(
					float32(w)-lw-128*menu.ratio,
					float32(h)*e.yp+fontOffset,
					0.7*menu.ratio, e.stringValue())
			}
		}
	}
}

// ContextReset uploads the UI images to the GPU.
// It should be called after each time the window is recreated.
func (menu *Menu) ContextReset() {
	menu.icons = map[string]uint32{
		"hexagon":    video.NewImage("assets/hexagon.png"),
		"main":       video.NewImage("assets/main.png"),
		"file":       video.NewImage("assets/file.png"),
		"folder":     video.NewImage("assets/folder.png"),
		"subsetting": video.NewImage("assets/subsetting.png"),
		"setting":    video.NewImage("assets/setting.png"),
		"resume":     video.NewImage("assets/resume.png"),
		"reset":      video.NewImage("assets/reset.png"),
		"loadstate":  video.NewImage("assets/loadstate.png"),
		"savestate":  video.NewImage("assets/savestate.png"),
		"screenshot": video.NewImage("assets/screenshot.png"),
		"add":        video.NewImage("assets/add.png"),
		"scan":       video.NewImage("assets/scan.png"),
		"on":         video.NewImage("assets/on.png"),
		"off":        video.NewImage("assets/off.png"),
	}

	usr, _ := user.Current()
	paths, _ := filepath.Glob(usr.HomeDir + "/.ludo/playlists/*.lpl")
	for _, path := range paths {
		path := path
		filename := utils.Filename(path)
		menu.icons[filename] = video.NewImage("assets/" + filename + ".png")
		menu.icons[filename+"-content"] = video.NewImage("assets/" + filename + "-content.png")
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

	if state.Global.CoreRunning {
		menu.stack = append(menu.stack, buildTabs())
		menu.stack[0].segueNext()
		menu.stack = append(menu.stack, buildMainMenu())
		menu.stack[1].segueNext()
		menu.stack = append(menu.stack, buildQuickMenu())
		fastForwardTweens()
	} else {
		menu.stack = append(menu.stack, buildTabs())
	}

	return menu
}
