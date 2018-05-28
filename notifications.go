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
	if len(vars) > 0 {
		fmt.Printf("["+prefix+"]: "+message+"\n", vars)
	} else {
		fmt.Printf("[" + prefix + "]: " + message + "\n")
	}
	notify(message, 240)
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
