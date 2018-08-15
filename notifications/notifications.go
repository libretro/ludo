package notifications

import (
	"fmt"
)

type Notification struct {
	Message string
	Frames  int
}

var notifications []Notification

func List() []Notification {
	return notifications
}

func Display(message string, frames int) {
	n := Notification{
		message,
		frames,
	}

	notifications = append(notifications, n)
}

func DisplayAndLog(prefix, message string, vars ...interface{}) {
	var msg string
	if len(vars) > 0 {
		msg = fmt.Sprintf(message, vars...)
	} else {
		msg = message
	}
	// if g.verbose {
	// 	log.Print("[" + prefix + "]: " + msg + "\n")
	// }
	Display(msg, 240)
}

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

func Clear() {
	notifications = []Notification{}
}
