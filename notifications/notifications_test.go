package notifications

import (
	"reflect"
	"testing"

	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
)

func Test_List(t *testing.T) {
	Clear()
	t.Run("Returns the notifications", func(t *testing.T) {
		Display(Error, "Test1", 240)
		Display(Error, "Test2", 240)
		Display(Error, "Test3", 240)
		got := List()
		want := []Notification{
			Notification{ID: got[0].ID, Severity: Error, Message: "Test1", Frames: 240},
			Notification{ID: got[1].ID, Severity: Error, Message: "Test2", Frames: 240},
			Notification{ID: got[2].ID, Severity: Error, Message: "Test3", Frames: 240},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func Test_Display(t *testing.T) {
	Clear()
	t.Run("Stacks notifications correctly", func(t *testing.T) {
		Display(Error, "Test1", 240)
		Display(Info, "Test2", 240)
		Display(Warning, "Test3", 240)
		got := notifications
		want := []Notification{
			Notification{ID: got[0].ID, Severity: Error, Message: "Test1", Frames: 240},
			Notification{ID: got[1].ID, Severity: Info, Message: "Test2", Frames: 240},
			Notification{ID: got[2].ID, Severity: Warning, Message: "Test3", Frames: 240},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func Test_DisplayAndLog(t *testing.T) {
	Clear()
	t.Run("Format message properly", func(t *testing.T) {
		DisplayAndLog(Info, "Tests", "Joypad #%d loaded with name %s.", 3, "Foo")
		got := notifications[0].Message
		want := "Joypad #3 loaded with name Foo."
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	Clear()
	t.Run("Format simple message properly", func(t *testing.T) {
		DisplayAndLog(Info, "Tests", "Hello world.")
		got := notifications[0].Message
		want := "Hello world."
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	Clear()
	t.Run("Logs to stdout if verbose", func(t *testing.T) {
		state.Global.Verbose = true
		got := utils.CaptureOutput(func() { DisplayAndLog(Info, "Test", "Joypad #%d loaded with name %s.", 3, "Foo") })
		want := "[Test]: Joypad #3 loaded with name Foo.\n"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	Clear()
	t.Run("Logs nothing if not verbose", func(t *testing.T) {
		state.Global.Verbose = false
		got := utils.CaptureOutput(func() { DisplayAndLog(Info, "Test", "Joypad #%d loaded with name %s.", 3, "Foo") })
		want := ""
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func Test_processNotifications(t *testing.T) {
	Clear()
	t.Run("Delete outdated notification", func(t *testing.T) {
		Display(Error, "Test1", 5)
		Display(Error, "Test1", 4)
		Display(Error, "Test1", 3)
		Display(Error, "Test2", 2)
		Display(Error, "Test3", 1)
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
		Display(Error, "Test1", 240)
		Display(Error, "Test2", 240)
		Display(Error, "Test3", 240)
		Clear()
		got := len(notifications)
		want := 0
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}
