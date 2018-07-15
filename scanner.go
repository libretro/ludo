package main

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/kivutar/go-playthemall/rdb"
)

type task struct {
	update func()
}

// LoadDB loops over the RDBs in a given directory and parses them
func LoadDB(dir string) (rdb.DB, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return rdb.DB{}, err
	}

	db := make(rdb.DB)
	for _, f := range files {
		filename := f.Name()
		system := filename[0 : len(filename)-4]
		bytes, _ := ioutil.ReadFile(dir + f.Name())
		db[system] = rdb.Parse(bytes)
	}

	return db, nil
}

// ScanDir scans a full directory, report progress and generate playlists
func ScanDir(dir string) {
	nid := notifyAndLog("Menu", "Scanning %s", dir)
	usr, _ := user.Current()
	roms := allFilesIn(dir)
	scannedGames := make(chan (rdb.Entry))
	go Scan(roms, scannedGames, g.db.Find)
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

// Scan scans a list of roms against the database
func Scan(roms []string, games chan (rdb.Entry), cb func(rompath string, romname string, CRC32 uint32, games chan (rdb.Entry))) {
	for _, f := range roms {
		ext := filepath.Ext(f)
		switch ext {
		case ".zip":
			// Open the ZIP archive
			z, _ := zip.OpenReader(f)
			for _, rom := range z.File {
				if rom.CRC32 > 0 {
					// Look for a matching game entry in the database
					cb(f, rom.Name, rom.CRC32, games)
				}
			}
			z.Close()
		}
	}
}
