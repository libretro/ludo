package main

import "fmt"

type notification struct {
	message string
	frames  int
}

var notifications []notification

func notify(message string, frames int) {
	n := notification{
		message,
		frames,
	}

	notifications = append(notifications, n)
}

func notifyAndLog(prefix, message string, vars ...interface{}) {
	var msg string
	if len(vars) > 0 {
		msg = fmt.Sprintf(message, vars...)
	} else {
		msg = message
	}
	if g.verbose {
		fmt.Print("[" + prefix + "]: " + msg + "\n")
	}
	notify(msg, 240)
}

func processNotifications() {
	deleted := 0
	for i := range notifications {
		j := i - deleted
		notifications[j].frames--
		if notifications[j].frames <= 0 {
			notifications = append(notifications[:j], notifications[j+1:]...)
			deleted++
		}
	}
}

func clearNotifications() {
	notifications = []notification{}
}
