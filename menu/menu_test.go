package menu

import (
	"os"
	"reflect"
	"testing"

	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"

	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type WindowMock struct{}

func (m WindowMock) GetFramebufferSize() (width, height int)     { return 320, 240 }
func (m WindowMock) Destroy()                                    {}
func (m WindowMock) MakeContextCurrent()                         {}
func (m WindowMock) SetSizeLimits(minw, minh, maxw, maxh int)    {}
func (m WindowMock) SetInputMode(mode glfw.InputMode, value int) {}
func (m WindowMock) GetKey(key glfw.Key) glfw.Action             { return 0 }
func (m WindowMock) SetShouldClose(bool)                         {}
func (m WindowMock) ShouldClose() bool                           { return false }
func (m WindowMock) SetTitle(string)                             {}
func (m WindowMock) SwapBuffers()                                {}

func Test_Init(t *testing.T) {

	var vid = &video.Video{
		Window: &WindowMock{},
	}

	Init(vid)

	t.Run("Starts with a single scene if no game is running", func(t *testing.T) {
		got := len(menu.stack)
		want := 1
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("Starts on the tabs scene if no game is running", func(t *testing.T) {
		got := menu.stack[0].Entry().label
		want := "Ludo"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	state.Global.CoreRunning = true
	Init(vid)

	t.Run("Warps at the quick menu if a game is launched", func(t *testing.T) {
		got := len(menu.stack)
		want := 3
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("Warps at the quick menu if a game is launched", func(t *testing.T) {
		got := menu.stack[len(menu.stack)-1].Entry().label
		want := "Quick Menu"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("No tweens are left after warping to quick menu", func(t *testing.T) {
		got := len(menu.tweens)
		want := 0
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func Test_fastForwardTweens(t *testing.T) {
	foo := float32(5)
	menu.tweens[&foo] = gween.New(foo, 25, 0.015, ease.OutSine)
	bar := float32(-10)
	menu.tweens[&bar] = gween.New(bar, 0, 0.9, ease.OutSine)

	fastForwardTweens()

	t.Run("No tweens are left", func(t *testing.T) {
		got := len(menu.tweens)
		want := 0
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("Tweened vars have the right value", func(t *testing.T) {
		got := foo
		want := float32(25)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("Tweened vars have the right value", func(t *testing.T) {
		got := bar
		want := float32(0)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func Test_buildExplorer(t *testing.T) {

	menu.stack = []Scene{}

	exec := 0
	cbMock := func(str string) {
		exec++
	}

	dirActionMock := &entry{
		label: "<Scan this directory>",
		icon:  "scan",
	}

	tmp := os.TempDir() + "/Test_buildExplorer/"

	os.RemoveAll(tmp)
	os.Mkdir(tmp, 0777)
	defer os.RemoveAll(tmp)
	os.Create(tmp + "File 1.txt")
	os.Create(tmp + "File 2.img")
	os.Create(tmp + "File 3.txt")
	os.Create(tmp + "File 4.img")
	os.Mkdir(tmp+"Folder 1", 0777)

	scene := buildExplorer(os.TempDir()+"/Test_buildExplorer/", []string{".img"}, cbMock, dirActionMock)
	menu.stack = append(menu.stack, scene)

	children := scene.Entry().children

	t.Run("Should display the right number of menu entries", func(t *testing.T) {
		if !(len(children) == 6) {
			t.Errorf("buildExplorer = %v, want %v", len(children), 6)
		}
	})

	t.Run("Inserts the directory Action as first entry", func(t *testing.T) {
		if children[0].label != dirActionMock.label {
			t.Errorf("buildExplorer = %v, want %v", children[0].label, dirActionMock.label)
		}
	})

	t.Run("Normal files have no OK callbacks", func(t *testing.T) {
		children[1].callbackOK()
		if exec != 0 {
			t.Errorf("buildExplorer = %v, want %v", exec, 0)
		}
	})

	t.Run("Files have file icon", func(t *testing.T) {
		if children[1].icon != "file" {
			t.Errorf("buildExplorer = %v, want %v", children[1].icon, "file")
		}
	})

	t.Run("Targeted files have OK callbacks", func(t *testing.T) {
		children[2].callbackOK()
		if exec != 1 {
			t.Errorf("buildExplorer = %v, want %v", exec, 1)
		}
	})

	t.Run("Folders have folder icon", func(t *testing.T) {
		if children[5].icon != "folder" {
			t.Errorf("buildExplorer = %v, want %v", children[1].icon, "folder")
		}
	})

	t.Run("Folder callback opens the folder", func(t *testing.T) {
		if len(menu.stack) != 1 {
			t.Errorf("buildExplorer = %v, want %v", len(menu.stack), 1)
		}
		children[5].callbackOK()
		if len(menu.stack) != 2 {
			t.Errorf("buildExplorer = %v, want %v", len(menu.stack), 2)
		}
	})
}
