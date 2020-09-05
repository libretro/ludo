// Package menu is the graphical interface allowing to browse games, launch
// games, configure settings, and display a contextual menu to interract with
// the running game.
package menu

import (
	"path/filepath"

	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
	"github.com/libretro/ludo/video"
)

var vid *video.Video
var menu *Menu

// Menu is a type holding the menu state, the stack of scenes, tweens, etc
type Menu struct {
	stack  []Scene
	icons  map[string]uint32
	tweens Tweens
	scroll float32
	ratio  float32
	t      float64
}

// Init initializes the menu.
// If a game is already running, it will warp the user to the quick menu.
// If not, it will display the menu tabs.
func Init(v *video.Video) *Menu {
	vid = v

	w, _ := v.Window.GetFramebufferSize()

	menu = &Menu{}
	menu.stack = []Scene{}
	menu.tweens = make(Tweens)
	menu.ratio = float32(w) / 1920
	menu.icons = map[string]uint32{}

	menu.Push(buildTabs())

	menu.ContextReset()

	return menu
}

// Push will navigate to a new scene. It usually happen when the user presses
// OK on a menu entry.
func (m *Menu) Push(s Scene) {
	m.stack = append(m.stack, s)
}

// Render takes care of rendering the menu
func (m *Menu) Render(dt float32) {
	// Early return to not render the menu, in case MenuActive is set to false
	// during the same mainloop iteration
	if !state.Global.MenuActive {
		return
	}

	m.t += float64(dt * 8)
	w, h := vid.Window.GetFramebufferSize()
	m.ratio = float32(w) / 1920

	if state.Global.CoreRunning {
		vid.DrawRect(0, 0, float32(w), float32(h), 0, bgColor.Alpha(0.85))
	} else {
		vid.DrawRect(0, 0, float32(w), float32(h), 0, bgColor)
	}

	m.tweens.Update(dt)

	currentScreenIndex := len(m.stack) - 1
	for i := 0; i <= currentScreenIndex+1; i++ {
		if i < 0 || i > currentScreenIndex {
			continue
		}

		m.stack[i].render()
	}
	m.stack[currentScreenIndex].drawHintBar()
}

// ContextReset uploads the UI images to the GPU.
// It should be called after each time the window is recreated.
func (m *Menu) ContextReset() {
	assets := settings.Current.AssetsDirectory

	paths, _ := filepath.Glob(assets + "/*.png")
	for _, path := range paths {
		path := path
		filename := utils.FileName(path)
		m.icons[filename] = video.NewImage(path)
	}

	paths, _ = filepath.Glob(assets + "/flags/*.png")
	for _, path := range paths {
		path := path
		filename := utils.FileName(path)
		m.icons[filename] = video.NewImage(path)
	}

	currentScreenIndex := len(m.stack) - 1
	curList := m.stack[currentScreenIndex].Entry()
	for i := range curList.children {
		curList.children[i].thumbnail = 0
	}
}

// WarpToQuickMenu loads the contextual menu for games that are launched from
// the command line interface or from 'Load Game'.
func (m *Menu) WarpToQuickMenu() {
	m.scroll = 0
	m.stack = []Scene{}
	m.Push(buildTabs())
	m.stack[0].segueNext()
	m.Push(buildMainMenu())
	m.stack[1].segueNext()
	m.Push(buildQuickMenu())
	m.tweens.FastForward()
}
