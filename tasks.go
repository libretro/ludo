package main

import (
	"os"
	"os/user"
	"strconv"

	"github.com/kivutar/go-playthemall/rdb"
)

type task struct {
	update func()
}

func scanDir(dir string) {
	nid := notifyAndLog("Menu", "Scanning %s", dir)
	usr, _ := user.Current()
	roms := allFilesIn(dir)
	scannedGames := make(chan (rdb.Game))
	go rdb.Scan(roms, scannedGames, g.db.Find)
	task := task{
		update: func() {
			i := 0
			for game := range scannedGames {
				i++
				lpl, _ := os.OpenFile(usr.HomeDir+"/.playthemall/playlists/"+game.System+".lpl", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
				lpl.WriteString(game.Path + "#" + game.ROMName + "\n")
				lpl.WriteString(game.Name + "\n")
				lpl.WriteString("DETECT\n")
				lpl.WriteString("DETECT\n")
				lpl.WriteString(strconv.FormatUint(uint64(game.CRC32), 10) + "|crc\n")
				lpl.WriteString(game.System + ".lpl\n")
				lpl.Close()
				notifications[nid].frames = 240
				notifications[nid].message = strconv.Itoa(i) + "/" + strconv.Itoa(len(roms)) + " " + game.Name
			}
		},
	}
	go task.update()
	g.tasks = append(g.tasks, task)
}
