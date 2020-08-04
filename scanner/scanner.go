// Package scanner generates game playlists by scanning your game collection
// against the database. It uses CRC checksums for No-Intro zip files and
// name matching for Redump cue files.
package scanner

import (
	"archive/zip"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/playlists"
	"github.com/libretro/ludo/rdb"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
)

// LoadDB loops over the RDBs in a given directory and parses them
func LoadDB(dir string) (rdb.DB, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return rdb.DB{}, err
	}
	db := make(rdb.DB)
	for _, f := range files {
		name := f.Name()
		if !strings.Contains(name, ".rdb") {
			continue
		}
		system := name[0 : len(name)-4]
		bytes, _ := ioutil.ReadFile(filepath.Join(dir, name))
		db[system] = rdb.Parse(bytes)
	}
	return db, nil
}

// ScanDir scans a full directory, report progress and generate playlists
func ScanDir(dir string, doneCb func()) {
	n := ntf.DisplayAndLog(ntf.Info, "Menu", "Scanning %s", dir)
	roms, err := utils.AllFilesIn(dir)
	if err != nil {
		n.Update(ntf.Error, err.Error())
		return
	}
	games := make(chan (rdb.Game))
	go Scan(dir, roms, games, n)
	go func() {
		i := 0
		for game := range games {
			os.MkdirAll(settings.Current.PlaylistsDirectory, os.ModePerm)
			CSVPath := filepath.Join(settings.Current.PlaylistsDirectory, game.System+".csv")
			if playlists.Contains(CSVPath, game.Path, game.CRC32) {
				continue
			}
			f, _ := os.OpenFile(CSVPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			f.WriteString(game.Path + "\t")
			f.WriteString(game.Name + "\t")
			if game.CRC32 > 0 {
				f.WriteString(strconv.FormatUint(uint64(game.CRC32), 16))
			}
			f.WriteString("\n")
			f.Close()
			i++
		}
		if doneCb != nil {
			doneCb()
		}
		n.Update(ntf.Success, "Done scanning. %d new games found.", i)
	}()
}

// Scan scans a list of roms against the database
func Scan(dir string, roms []string, games chan (rdb.Game), n *ntf.Notification) {
	for i, f := range roms {
		ext := filepath.Ext(f)
		switch ext {
		case ".zip":
			// Open the ZIP archive
			z, err := zip.OpenReader(f)
			if err != nil {
				n.Update(ntf.Error, err.Error())
				continue
			}
			for _, rom := range z.File {
				if rom.CRC32 > 0 {
					// Look for a matching game entry in the database
					state.Global.DB.FindByCRC(f, rom.Name, rom.CRC32, games)
					n.Update(ntf.Info, strconv.Itoa(i)+"/"+strconv.Itoa(len(roms))+" "+f)
				}
			}
			z.Close()
		case ".cue":
			// Look for a matching game entry in the database
			state.Global.DB.FindByROMName(f, filepath.Base(f), 0, games)
			n.Update(ntf.Info, strconv.Itoa(i)+"/"+strconv.Itoa(len(roms))+" "+f)
		case ".32x", "a52", ".a78", ".col", ".crt", ".d64", ".pce", ".fds", ".gb", ".gba", ".gbc", ".gen", ".gg", ".ipf", ".j64", ".jag", ".lnx", ".md", ".n64", ".nes", ".ngc", ".nds", ".rom", ".sfc", ".sg", ".smc", ".smd", ".sms", ".ws", ".wsc":
			bytes, err := ioutil.ReadFile(f)
			if err != nil {
				n.Update(ntf.Error, err.Error())
				continue
			}
			CRC32 := crc32.ChecksumIEEE(bytes)
			state.Global.DB.FindByCRC(f, utils.FileName(f), CRC32, games)
			n.Update(ntf.Info, strconv.Itoa(i)+"/"+strconv.Itoa(len(roms))+" "+f)
		}
	}
	close(games)
}
