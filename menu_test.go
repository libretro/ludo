package main

import (
	"reflect"
	"testing"

	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type CtxMock struct{}

func (m CtxMock) GetFramebufferSize() (width, height int) {
	return 320, 240
}

func Test_menuInit(t *testing.T) {

	var ctx CtxMock
	menuInit(ctx)

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

	g.coreRunning = true
	menuInit(ctx)

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
