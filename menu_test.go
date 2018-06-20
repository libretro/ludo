package main

import (
	"reflect"
	"testing"
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
}
