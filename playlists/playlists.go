package playlists

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/libretro/ludo/settings"
)

// PlaylistEntry represents a game in a playlist.
type PlaylistEntry struct {
	Path  string // Absolute path of the game on the filesystem
	Name  string // Human readable name of the game, comes from the RDB
	CRC32 uint32 // Checksum of the game, used for deduplication
	LPL   string
}

// Playlist is a list of games, result of scanning for games on the filesystem.
type Playlist []PlaylistEntry

// Playlists is a map of playlists organized per system.
var Playlists = map[string]Playlist{}

// Load loops over lpl files in the playlists directory and loads them into
// memory.
func Load() {
	paths, _ := filepath.Glob(settings.Current.PlaylistsDirectory + "/*.lpl")

	Playlists = map[string]Playlist{}
	for _, path := range paths {
		path := path

		file, _ := os.Open(path)
		defer file.Close()
		reader := csv.NewReader(bufio.NewReader(file))
		reader.Comma = '\t'

		playlist := Playlist{}
		for {
			line, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Println(err)
				continue
			}
			var entry PlaylistEntry
			entry.Path = line[0]
			entry.Name = line[1]
			if line[2] != "DETECT" {
				u64, err := strconv.ParseUint(strings.Replace(line[2], "|crc", "", -1), 16, 64)
				if err != nil {
					log.Println(err)
				} else {
					entry.CRC32 = uint32(u64)
				}
			}
			entry.LPL = line[3]

			playlist = append(playlist, entry)
		}
		Playlists[path] = playlist
	}
}

// Contains checks if a game is already in a playlist.
func Contains(lplpath, path string, CRC32 uint32) bool {
	for _, entry := range Playlists[lplpath] {
		// Be careful, sometimes we don't have a CRC32
		if entry.Path == path || (CRC32 != 0 && entry.CRC32 == CRC32) {
			return true
		}
	}
	return false
}
