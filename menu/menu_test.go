package menu

import (
	"reflect"
	"testing"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/libretro/go-playthemall/options"
	"github.com/libretro/go-playthemall/state"
	"github.com/libretro/go-playthemall/video"
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
	var opts *options.Options

	Init(vid, opts)

	t.Run("Starts with a single scene if no game is running", func(t *testing.T) {
		got := len(menu.stack)
		want := 1
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("Starts on the tabs scene if no game is running", func(t *testing.T) {
		got := menu.stack[0].Entry().label
		want := "Play Them All"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	state.Global.CoreRunning = true
	Init(vid, opts)

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
