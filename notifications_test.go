package main

import (
	"reflect"
	"testing"
)

func Test_notify(t *testing.T) {
	clearNotifications()
	t.Run("Stacks notifications correctly", func(t *testing.T) {
		notify("Test1", 240)
		notify("Test2", 240)
		notify("Test3", 240)
		got := notifications
		want := []notification{
			notification{message: "Test1", frames: 240},
			notification{message: "Test2", frames: 240},
			notification{message: "Test3", frames: 240},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func Test_notifyAndLog(t *testing.T) {
	clearNotifications()
	t.Run("Format message properly", func(t *testing.T) {
		notifyAndLog("Tests", "Joypad #%d loaded with name %s.", 3, "Foo")
		got := notifications[0].message
		want := "Joypad #3 loaded with name Foo."
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	clearNotifications()
	t.Run("Logs to stdout if verbose", func(t *testing.T) {
		g.verbose = true
		got := captureOutput(func() { notifyAndLog("Test", "Joypad #%d loaded with name %s.", 3, "Foo") })
		want := "[Test]: Joypad #3 loaded with name Foo.\n"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	clearNotifications()
	t.Run("Logs nothing if not verbose", func(t *testing.T) {
		g.verbose = false
		got := captureOutput(func() { notifyAndLog("Test", "Joypad #%d loaded with name %s.", 3, "Foo") })
		want := ""
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func Test_processNotifications(t *testing.T) {
	clearNotifications()
	t.Run("Delete outdated notification", func(t *testing.T) {
		notify("Test1", 5)
		notify("Test1", 4)
		notify("Test1", 3)
		notify("Test2", 2)
		notify("Test3", 1)
		processNotifications()
		processNotifications()
		got := len(notifications)
		want := 3
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func Test_clearNotifications(t *testing.T) {
	clearNotifications()
	t.Run("Empties the notification list", func(t *testing.T) {
		notify("Test1", 240)
		notify("Test2", 240)
		notify("Test3", 240)
		clearNotifications()
		got := len(notifications)
		want := 0
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}
