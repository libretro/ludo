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
		Display(Error, "Test1", Medium)
		Display(Error, "Test2", Medium)
		Display(Error, "Test3", Medium)
		got := List()
		want := []*Notification{
			&Notification{Severity: Error, Message: "Test1", Duration: Medium},
			&Notification{Severity: Error, Message: "Test2", Duration: Medium},
			&Notification{Severity: Error, Message: "Test3", Duration: Medium},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func Test_Display(t *testing.T) {
	Clear()
	t.Run("Stacks notifications correctly", func(t *testing.T) {
		Display(Error, "Test1", Medium)
		Display(Info, "Test2", Medium)
		Display(Warning, "Test3", Medium)
		got := notifications
		want := []*Notification{
			&Notification{Severity: Error, Message: "Test1", Duration: Medium},
			&Notification{Severity: Info, Message: "Test2", Duration: Medium},
			&Notification{Severity: Warning, Message: "Test3", Duration: Medium},
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
		state.Verbose = true
		got := utils.CaptureOutput(func() { DisplayAndLog(Info, "Test", "Joypad #%d loaded with name %s.", 3, "Foo") })
		want := "[Test]: Joypad #3 loaded with name Foo.\n"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	Clear()
	t.Run("Logs nothing if not verbose", func(t *testing.T) {
		state.Verbose = false
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
		Process(1)
		Process(1)
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
		Display(Error, "Test1", Medium)
		Display(Error, "Test2", Medium)
		Display(Error, "Test3", Medium)
		Clear()
		got := len(notifications)
		want := 0
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func Test_Update(t *testing.T) {
	Clear()
	t.Run("Is able to update a notification independently", func(t *testing.T) {
		Display(Error, "Test1", Medium/2)
		nid2 := Display(Error, "Test2", Medium)
		Display(Error, "Test3", Medium)

		Process(0.5)
		Process(0.5)
		Process(0.5)
		Process(0.5)
		nid2.Update(Success, "Test4")
		Process(0.5)

		got := List()
		want := []*Notification{
			&Notification{Severity: Success, Message: "Test4", Duration: Medium - 0.5},
			&Notification{Severity: Error, Message: "Test3", Duration: Medium - 0.5*5},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}
