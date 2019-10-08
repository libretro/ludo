// Package notifications exposes functions to display messages in toast
// widgets.
package notifications

import (
	"fmt"
	"log"

	"github.com/libretro/ludo/state"

	"github.com/rs/xid"
)

// Severity represents the severity of a notification message. It will affect
// the color of the notification text in the UI.
type Severity uint8

const (
	// Info is for informative message, when everything is fine
	Info Severity = iota
	// Success is for successful actions
	Success
	// Warning is also for informative messages, when something is not right
	// for example, if a menu entry has not been implemented.
	Warning
	// Error is for failed actions. For example, trying to load a game that
	// doesn't exists.
	Error
)

// Notification is a message that will be displayed on the screen during a
// certain time.
type Notification struct {
	ID       xid.ID
	Severity Severity
	Message  string
	Duration float32
}

// Medium is the standard duration for a notification
const Medium float32 = 4

var notifications []Notification

// List lists the current notifications.
func List() []Notification {
	return notifications
}

// Display creates a new notification.
func Display(severity Severity, message string, duration float32) xid.ID {
	id := xid.New()
	n := Notification{
		id,
		severity,
		message,
		duration,
	}

	notifications = append(notifications, n)

	return id
}

// DisplayAndLog creates a new notification and also logs the message to stdout.
func DisplayAndLog(severity Severity, prefix, message string, vars ...interface{}) xid.ID {
	var msg string
	if len(vars) > 0 {
		msg = fmt.Sprintf(message, vars...)
	} else {
		msg = message
	}
	if state.Global.Verbose {
		log.Print("[" + prefix + "]: " + msg + "\n")
	}
	return Display(severity, msg, Medium)
}

// Process iterates over the notifications, update them, delete the old ones.
func Process(dt float32) {
	deleted := 0
	for i := range notifications {
		j := i - deleted
		notifications[j].Duration -= dt
		if notifications[j].Duration <= 0 {
			notifications = append(notifications[:j], notifications[j+1:]...)
			deleted++
		}
	}
}

// Clear empties the notification list
func Clear() {
	notifications = []Notification{}
}

// find notification by unique ID
func find(id xid.ID) *Notification {
	for i := range notifications {
		if notifications[i].ID == id {
			return &notifications[i]
		}
	}
	return nil
}

// Update the message of a given notification. Also resets the delay before
// disapearing.
func Update(id xid.ID, severity Severity, message string, vars ...interface{}) {
	n := find(id)
	if n == nil {
		return
	}

	var msg string
	if len(vars) > 0 {
		msg = fmt.Sprintf(message, vars...)
	} else {
		msg = message
	}

	n.Duration = Medium
	n.Message = msg
	n.Severity = severity
}
