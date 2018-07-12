package main

import (
	"strconv"

	"github.com/kivutar/go-playthemall/rdb"
)

type task struct {
	update func()
}

func scanDir(dir string) {
	nid := notifyAndLog("Menu", "Scanning %s", dir)
	roms := allFilesIn(dir)
	scannedGames := make(chan (rdb.Game))
	go rdb.Scan(roms, scannedGames, g.db.Find)
	task := task{
		update: func() {
			i := 0
			for game := range scannedGames {
				i++
				notifications[nid].frames = 240
				notifications[nid].message = strconv.Itoa(i) + "/" + strconv.Itoa(len(roms)) + " " + game.Name
			}
		},
	}
	go task.update()
	g.tasks = append(g.tasks, task)
}
