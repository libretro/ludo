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
	"sort"
	"strconv"

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

// Gets a list of full paths to playlists
func getPaths() (paths []string) {
	paths, err := filepath.Glob(settings.Current.PlaylistsDirectory + "/*.csv")
	if err != nil {
		log.Println(err)
	}
	return
}

// Load loops over lpl files in the playlists directory and loads them into
// memory.
func Load() {
	for _, path := range getPaths() {
		path := path

		file, err := os.Open(path)
		if err != nil {
			log.Println(err)
			continue
		}
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
			entry.Path = filepath.Clean(line[0])
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
		sort.Slice(playlist, func(i, j int) bool {
			return playlist[i].Name < playlist[j].Name
		})
		Playlists[path] = playlist
	}
}

// Contains checks if a game is already in a playlist.
func Contains(CSVPath, path string, CRC32 uint32) bool {
	for _, entry := range Playlists[filepath.Clean(CSVPath)] {
		// Be careful, sometimes we don't have a CRC32
		if filepath.Clean(entry.Path) == filepath.Clean(path) || (CRC32 != 0 && entry.CRC32 == CRC32) {
			return true
		}
	}
	return false
}

// Count is a quick way of knowing how many games are in a playlist
func Count(path string) int {
	return len(Playlists[filepath.Clean(path)])
}

// Save will write a playlist to the filesystem
func Save(path string) {
	f, _ := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()
	for _, game := range Playlists[path] {
		f.WriteString(game.Path + "\t")
		f.WriteString(game.Name + "\t")
		f.WriteString(strconv.FormatUint(uint64(game.CRC32), 16) + "\n")
	}
}

// ShortName shortens the name of some game systems that are too long to be
// displayed in the menu
func ShortName(in string) string {
	shortNames := map[string]string{
		"Atari - 2600":                                   "Atari 2600",
		"Atari - 5200":                                   "Atari 5200",
		"Atari - 7800":                                   "Atari 7800",
		"Atari - Jaguar":                                 "Atari Jaguar",
		"Atari - Lynx":                                   "Atari Lynx",
		"Atari - ST":                                     "Atari ST",
		"Bandai - WonderSwan Color":                      "WonderSwan Color",
		"Bandai - WonderSwan":                            "WonderSwan",
		"Coleco - ColecoVision":                          "ColecoVision",
		"Commodore - 64":                                 "Commodore 64",
		"FBNeo - Arcade Games":                           "Arcade (FBNeo)",
		"GCE - Vectrex":                                  "Vectrex",
		"Magnavox - Odyssey2":                            "Magnavox Odyssey²",
		"Microsoft - MSX":                                "MSX",
		"Microsoft - MSX2":                               "MSX2",
		"NEC - PC Engine - TurboGrafx 16":                "TurboGrafx-16",
		"NEC - PC Engine CD - TurboGrafx-CD":             "TurboGrafx-CD",
		"NEC - PC Engine SuperGrafx":                     "SuperGrafx",
		"NEC - PC-FX":                                    "PC-FX",
		"Nintendo - Family Computer Disk System":         "Famicom Disk System",
		"Nintendo - Game Boy Advance":                    "Game Boy Advance",
		"Nintendo - Game Boy Color":                      "Game Boy Color",
		"Nintendo - Game Boy":                            "Game Boy",
		"Nintendo - Nintendo Entertainment System":       "NES / Famicom",
		"Nintendo - Pokemon Mini":                        "Pokemon Mini",
		"Nintendo - Super Nintendo Entertainment System": "Super Nintendo",
		"Nintendo - Virtual Boy":                         "Virtual Boy",
		"Sega - 32X":                                     "32X",
		"Sega - Game Gear":                               "Game Gear",
		"Sega - Master System - Mark III":                "Master System",
		"Sega - Mega Drive - Genesis":                    "Mega Drive / Genesis",
		"Sega - PICO":                                    "Pico",
		"Sega - Saturn":                                  "Saturn",
		"Sega - SG-1000":                                 "SG-1000",
		"Sharp - X68000":                                 "X68000",
		"Sinclair - ZX 81":                               "ZX81",
		"Sinclair - ZX Spectrum +3":                      "ZX Spectrum +3",
		"Sinclair - ZX Spectrum":                         "ZX Spectrum",
		"SNK - Neo Geo CD":                               "Neo Geo CD",
		"SNK - Neo Geo Pocket Color":                     "Neo Geo Pocket Color",
		"SNK - Neo Geo Pocket":                           "Neo Geo Pocket",
		"Sony - PlayStation":                             "PlayStation",
		"The 3DO Company - 3DO":                          "3DO",
		"Uzebox":                                         "Uzebox",
	}

	out, ok := shortNames[in]
	if !ok {
		return in
	}
	return out
}
