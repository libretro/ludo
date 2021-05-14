// Package dat is a parser for dat files, a binary database of games with
// metadata also used by RetroArch.
package dat

import (
	"encoding/xml"
	"log"
	"strconv"
	"sync"
)

// DB is a database that contains many Dats, mapped to their system name
type DB map[string]Dat

// Dat is a list of the games of a system
type Dat struct {
	XMLName xml.Name `xml:"datafile"`
	Games   []Game   `xml:"game"`
}

// Game represents a game and can contain a list of ROMs
type Game struct {
	XMLName     xml.Name `xml:"game"`
	Name        string   `xml:"name,attr"`
	Description string   `xml:"description"` // The human readable name of the game
	ROMs        []ROM    `xml:"rom"`

	Path   string
	System string
}

// CRC is the CRC32 checksum of a ROM
type CRC uint32

// ROM can be a game file or part of a game
type ROM struct {
	XMLName xml.Name `xml:"rom"`
	Name    string   `xml:"name,attr"`
	CRC     CRC      `xml:"crc,attr"`
}

// UnmarshalXMLAttr is used to parse a hex number in string form to uint
func (s *CRC) UnmarshalXMLAttr(attr xml.Attr) error {
	u64, err := strconv.ParseUint(attr.Value, 16, 64)
	if err != nil {
		log.Println(err)
	} else {
		*s = CRC(uint32(u64))
	}
	return nil
}

// Parse parses a .dat file content and returns an array of Entries
func Parse(dat []byte) Dat {
	var output Dat

	err := xml.Unmarshal(dat, &output)
	if err != nil {
		log.Println(err)
	}

	return output
}

// FindByCRC loops over the Dats in the DB and concurrently matches CRC checksums.
func (db *DB) FindByCRC(romPath string, romName string, crc uint32, games chan (Game)) {
	var wg sync.WaitGroup
	wg.Add(len(*db))
	// For every Dat in the DB
	for system, dat := range *db {
		go func(dat Dat, crc uint32, system string) {
			// For each game in the Dat
			for _, game := range dat.Games {
				if len(game.ROMs) == 0 {
					continue
				}
				// If the checksums match
				if crc == uint32(game.ROMs[0].CRC) {
					game.Path = romPath
					game.System = system
					games <- game
				}
			}
			wg.Done()
		}(dat, crc, system)
	}
	// Synchronize all the goroutines
	wg.Wait()
}

// FindByROMName loops over the Dats in the DB and concurrently matches ROM names.
func (db *DB) FindByROMName(romPath string, romName string, crc uint32, games chan (Game)) {
	var wg sync.WaitGroup
	wg.Add(len(*db))
	// For every Dat in the DB
	for system, dat := range *db {
		go func(dat Dat, crc uint32, system string) {
			// For each game in the Dat
			for _, game := range dat.Games {
				if len(game.ROMs) == 0 {
					continue
				}
				// If the checksums match
				if romName == game.ROMs[0].Name {
					game.Path = romPath
					game.System = system
					games <- game
				}
			}
			wg.Done()
		}(dat, crc, system)
	}
	// Synchronize all the goroutines
	wg.Wait()
}
