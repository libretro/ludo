package main

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
