// Package playlists is the playlist manager of Ludo. In Ludo, playlists are
// CSV files containing the ROM path, name, and CRC32 checksum.
// Playlists are kept into memory for fast lookup of entries and deduplication.
package playlists

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/libretro/ludo/settings"
)

// Game represents a game in a playlist.
type Game struct {
	Path  string // Absolute path of the game on the filesystem
	Name  string // Human readable name of the game, comes from the RDB
	CRC32 uint32 // Checksum of the game, used for deduplication
}

// Playlist is a list of games, result of scanning for games on the filesystem.
type Playlist []Game

// Playlists is a map of playlists organized per system.
var Playlists = map[string]Playlist{}

// Load loops over lpl files in the playlists directory and loads them into
// memory.
func Load() {
	paths, _ := filepath.Glob(settings.Current.PlaylistsDirectory + "/*.csv")

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
			var entry Game
			entry.Path = line[0]
			entry.Name = line[1]
			if line[2] != "" {
				u64, err := strconv.ParseUint(line[2], 16, 64)
				if err != nil {
					log.Println(err)
				} else {
					entry.CRC32 = uint32(u64)
				}
			}

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

// Count is a quick way of knowing how many games are in a playlist
func Count(path string) int {
	return len(Playlists[path])
}

// ShortName shortens the name of some game systems that are too long to be
// displayed in the menu
func ShortName(in string) string {
	if len(in) < 20 {
		return in
	}
	r, _ := regexp.Compile(`(.*?) - (.*)`)
	out := r.ReplaceAllString(in, "$2")
	out = strings.Replace(out, "Nintendo Entertainment System", "NES", -1)
	out = strings.Replace(out, "PC Engine", "PCE", -1)
	return out
}
