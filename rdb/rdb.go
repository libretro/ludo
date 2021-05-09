// Package rdb is a parser for rdb files, a binary database of games with
// metadata also used by RetroArch.
package rdb

import (
	"encoding/xml"
	"log"
	"strconv"
	"sync"
)

// DB is a database that contains many Dats, mapped to their system name
type DB map[string]Dat

type Dat struct {
	XMLName xml.Name `xml:"datafile"`
	Games   []Game   `xml:"game"`
}

type Game struct {
	XMLName     xml.Name `xml:"game"`
	Name        string   `xml:"name,attr"`
	Description string   `xml:"description"`
	ROMs        []ROM    `xml:"rom"`

	Path   string
	System string
}

type CRC uint32

type ROM struct {
	XMLName xml.Name `xml:"rom"`
	Name    string   `xml:"name,attr"`
	CRC     CRC      `xml:"crc,attr"`
}

func (s *CRC) UnmarshalXMLAttr(attr xml.Attr) error {
	u64, err := strconv.ParseUint(attr.Value, 16, 64)
	if err != nil {
		log.Println(err)
	} else {
		*s = CRC(uint32(u64))
	}
	return nil
}

// Parse parses a .rdb file content and returns an array of Entries
func Parse(rdb []byte) Dat {
	var output Dat

	err := xml.Unmarshal(rdb, &output)
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
