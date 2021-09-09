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

	"github.com/libretro/ludo/dat"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/playlists"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
)

// LoadDB loops over the RDBs in a given directory and parses them
func LoadDB(dir string) (dat.DB, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return dat.DB{}, err
	}
	db := make(dat.DB)
	for _, f := range files {
		name := f.Name()
		if !strings.Contains(name, ".dat") {
			continue
		}
		system := name[0 : len(name)-4]
		bytes, _ := ioutil.ReadFile(filepath.Join(dir, name))
		db[system] = dat.Parse(bytes)
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
	games := make(chan (dat.Game))
	go Scan(dir, roms, games, n)
	go func() {
		i := 0
		for game := range games {
			os.MkdirAll(settings.Current.PlaylistsDirectory, os.ModePerm)
			CSVPath := filepath.Join(settings.Current.PlaylistsDirectory, game.System+".csv")
			if playlists.Contains(CSVPath, game.Path, uint32(game.ROMs[0].CRC)) {
				continue
			}
			f, _ := os.OpenFile(CSVPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if len(game.Description) == 0 {
				continue
			}
			f.WriteString(game.Path + "\t")
			f.WriteString(game.Description + "\t")
			if game.ROMs[0].CRC > 0 {
				f.WriteString(strconv.FormatUint(uint64(game.ROMs[0].CRC), 16))
			}
			f.WriteString("\n")
			f.Close()
			i++
		}
		doneCb()
		n.Update(ntf.Success, "Done scanning. %d new games found.", i)
	}()
}

// Returns the checksum and headerless checksum of a ROM
func checksumHeaderless(rom *zip.File, headerSize uint) (uint32, uint32, error) {
	h, err := rom.Open()
	if err != nil {
		return 0, 0, err
	}
	defer h.Close()
	bytes, err := ioutil.ReadAll(h)
	if err != nil {
		return 0, 0, err
	}
	crc := crc32.ChecksumIEEE(bytes)
	crcHeaderless := crc32.ChecksumIEEE(bytes[headerSize:])
	return crc, crcHeaderless, nil
}

// Some ROMs have a header that we will need to remove to calculare the checksum
// Our database has checksums of headerless ROMs
var headerSizes = map[string]uint{
	".nes": 16,
	".fds": 16,
	".a78": 128,
	".lnx": 64,
}

// Scan scans a list of roms against the database
func Scan(dir string, roms []string, games chan (dat.Game), n *ntf.Notification) {
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
				romExt := filepath.Ext(rom.Name)
				// these 4 systems might have headered or headerless roms and need special logic
				if headerSize, ok := headerSizes[romExt]; ok {
					crc, crcHeaderless, err := checksumHeaderless(rom, headerSize)
					if err != nil {
						n.Update(ntf.Error, err.Error())
						continue
					}
					state.DB.FindByCRC(f, rom.Name, crc, games)
					state.DB.FindByCRC(f, rom.Name, crcHeaderless, games)
					n.Update(ntf.Info, strconv.Itoa(i)+"/"+strconv.Itoa(len(roms))+" "+f)
				} else if rom.CRC32 > 0 {
					// Look for a matching game entry in the database
					state.DB.FindByCRC(f, rom.Name, rom.CRC32, games)
					n.Update(ntf.Info, strconv.Itoa(i)+"/"+strconv.Itoa(len(roms))+" "+f)
				}
			}
			z.Close()
		case ".cue":
			// Look for a matching game entry in the database
			state.DB.FindByROMName(f, filepath.Base(f), 0, games)
			n.Update(ntf.Info, strconv.Itoa(i)+"/"+strconv.Itoa(len(roms))+" "+f)
		case ".32x", "a52", ".a78", ".col", ".crt", ".d64", ".pce", ".fds", ".gb", ".gba", ".gbc", ".gen", ".gg", ".ipf", ".j64", ".jag", ".lnx", ".md", ".n64", ".nes", ".ngc", ".nds", ".rom", ".sfc", ".sg", ".smc", ".smd", ".sms", ".ws", ".wsc":
			bytes, err := ioutil.ReadFile(f)
			if err != nil {
				n.Update(ntf.Error, err.Error())
				continue
			}
			crc := crc32.ChecksumIEEE(bytes)
			state.DB.FindByCRC(f, utils.FileName(f), crc, games)
			if headerSize, ok := headerSizes[ext]; ok {
				crcHeaderless := crc32.ChecksumIEEE(bytes[headerSize:])
				state.DB.FindByCRC(f, utils.FileName(f), crcHeaderless, games)
			}
			n.Update(ntf.Info, strconv.Itoa(i)+"/"+strconv.Itoa(len(roms))+" "+f)
		}
	}
	close(games)
}
