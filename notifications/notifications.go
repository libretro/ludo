// Package notifications exposes functions to display messages in toast
// widgets.
package notifications

import (
	"fmt"
	"log"

	"github.com/libretro/ludo/state"
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
	Severity Severity
	Message  string
	Duration float32
}

// Medium is the standard duration for a notification
const Medium float32 = 4

var notifications []*Notification

// List lists the current notifications.
func List() []*Notification {
	return notifications
}

// Display creates a new notification.
func Display(severity Severity, message string, duration float32) *Notification {
	n := &Notification{
		severity,
		message,
		duration,
	}

	notifications = append(notifications, n)

	return n
}

// DisplayAndLog creates a new notification and also logs the message to stdout.
func DisplayAndLog(severity Severity, prefix, message string, vars ...interface{}) *Notification {
	msg := fmt.Sprintf(message, vars...)
	if state.Verbose {
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
	notifications = []*Notification{}
}

// Update the message of a given notification. Also resets the delay before
// disapearing.
func (n *Notification) Update(severity Severity, message string, vars ...interface{}) {
	msg := fmt.Sprintf(message, vars...)

	n.Duration = Medium
	n.Message = msg
	n.Severity = severity
}
