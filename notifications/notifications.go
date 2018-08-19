package notifications

import (
	"fmt"
	"log"

	"github.com/libretro/go-playthemall/state"
)

// Notification is a message that will be displayed on the screen during a number of frames.
type Notification struct {
	Message string
	Frames  int
}

var notifications []Notification

// List lists the current notifications.
func List() []Notification {
	return notifications
}

// Display creates a new notification.
func Display(message string, frames int) {
	n := Notification{
		message,
		frames,
	}

	notifications = append(notifications, n)
}

// DisplayAndLog creates a new notification and also logs the message to stdout.
func DisplayAndLog(prefix, message string, vars ...interface{}) {
	var msg string
	if len(vars) > 0 {
		msg = fmt.Sprintf(message, vars...)
	} else {
		msg = message
	}
	if state.Global.Verbose {
		log.Print("[" + prefix + "]: " + msg + "\n")
	}
	Display(msg, 240)
}

// Process iterates over the notifications, update them, delete the old ones.
func Process() {
	deleted := 0
	for i := range notifications {
		j := i - deleted
		notifications[j].Frames--
		if notifications[j].Frames <= 0 {
			notifications = append(notifications[:j], notifications[j+1:]...)
			deleted++
		}
	}
}

// Clear empties the notification list
func Clear() {
	notifications = []Notification{}
}
