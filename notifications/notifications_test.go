package notifications

import (
	"reflect"
	"testing"

	"github.com/libretro/go-playthemall/state"
	"github.com/libretro/go-playthemall/utils"
)

func Test_List(t *testing.T) {
	Clear()
	t.Run("Returns the notifications", func(t *testing.T) {
		Display("Test1", 240)
		Display("Test2", 240)
		Display("Test3", 240)
		got := List()
		want := []Notification{
			Notification{Message: "Test1", Frames: 240},
			Notification{Message: "Test2", Frames: 240},
			Notification{Message: "Test3", Frames: 240},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func Test_Display(t *testing.T) {
	Clear()
	t.Run("Stacks notifications correctly", func(t *testing.T) {
		Display("Test1", 240)
		Display("Test2", 240)
		Display("Test3", 240)
		got := notifications
		want := []Notification{
			Notification{Message: "Test1", Frames: 240},
			Notification{Message: "Test2", Frames: 240},
			Notification{Message: "Test3", Frames: 240},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func Test_DisplayAndLog(t *testing.T) {
	Clear()
	t.Run("Format message properly", func(t *testing.T) {
		DisplayAndLog("Tests", "Joypad #%d loaded with name %s.", 3, "Foo")
		got := notifications[0].Message
		want := "Joypad #3 loaded with name Foo."
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	Clear()
	t.Run("Format simple message properly", func(t *testing.T) {
		DisplayAndLog("Tests", "Hello world.")
		got := notifications[0].Message
		want := "Hello world."
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	Clear()
	t.Run("Logs to stdout if verbose", func(t *testing.T) {
		state.Global.Verbose = true
		got := utils.CaptureOutput(func() { DisplayAndLog("Test", "Joypad #%d loaded with name %s.", 3, "Foo") })
		want := "[Test]: Joypad #3 loaded with name Foo.\n"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	Clear()
	t.Run("Logs nothing if not verbose", func(t *testing.T) {
		state.Global.Verbose = false
		got := utils.CaptureOutput(func() { DisplayAndLog("Test", "Joypad #%d loaded with name %s.", 3, "Foo") })
		want := ""
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func Test_processNotifications(t *testing.T) {
	Clear()
	t.Run("Delete outdated notification", func(t *testing.T) {
		Display("Test1", 5)
		Display("Test1", 4)
		Display("Test1", 3)
		Display("Test2", 2)
		Display("Test3", 1)
		Process()
		Process()
		got := len(notifications)
		want := 3
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func Test_Clear(t *testing.T) {
	Clear()
	t.Run("Empties the notification list", func(t *testing.T) {
		Display("Test1", 240)
		Display("Test2", 240)
		Display("Test3", 240)
		Clear()
		got := len(notifications)
		want := 0
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}
