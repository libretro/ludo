package rdb

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"strconv"
	"sync"
)

// DB is a database that contains many RDB, mapped to their system name
type DB map[string]RDB

// RDB contains all the game descriptions for a system
type RDB []Game

// Game represents a game in the libretro database
type Game struct {
	Name        string
	Description string
	Genre       string
	Developer   string
	Publisher   string
	Franchise   string
	Serial      string
	ROMName     string
	Size        uint64
	CRC32       uint32
}

const (
	mpfFixMap = 0x80
	mpfMap16  = 0xde
	mpfMap32  = 0xdf

	mpfFixArray = 0x90
	mpfArray16  = 0xdc
	mpfArray32  = 0xdd

	mpfFixStr = 0xa0
	mpfStr8   = 0xd9
	mpfStr16  = 0xda
	mpfStr32  = 0xdb

	mpfBin8  = 0xc4
	mpfBin16 = 0xc5
	mpfBin32 = 0xc6

	mpfFalse = 0xc2
	mpfTrue  = 0xc3

	mpfInt8  = 0xd0
	mpfInt16 = 0xd1
	mpfInt32 = 0xd2
	mpfInt64 = 0xd3

	mpfUint8  = 0xcc
	mpfUint16 = 0xcd
	mpfUint32 = 0xce
	mpfUint64 = 0xcf

	mpfNil = 0xc0
)

func setField(g *Game, key string, value string) {
	switch key {
	case "name":
		g.Name = string(value[:])
	case "description":
		g.Description = string(value[:])
	case "genre":
		g.Genre = string(value[:])
	case "developer":
		g.Developer = string(value[:])
	case "publisher":
		g.Publisher = string(value[:])
	case "franchise":
		g.Franchise = string(value[:])
	case "serial":
		g.Serial = string(value[:])
	case "rom_name":
		g.ROMName = string(value[:])
	case "size":
		value2 := fmt.Sprintf("%x", string(value[:]))
		u64, _ := strconv.ParseUint(value2, 16, 32)
		g.Size = u64
	case "crc":
		value2 := fmt.Sprintf("%x", string(value[:]))
		u64, _ := strconv.ParseUint(value2, 16, 32)
		g.CRC32 = uint32(u64)
	}
}

// ParseRDB parses a .rdb file content and returns an array of Games
func ParseRDB(rdb []byte) RDB {
	var output RDB
	pos := 0x10
	iskey := false
	key := ""

	g := Game{}

	for int(rdb[pos]) != mpfNil {
		fieldtype := int(rdb[pos])

		var value []byte

		if fieldtype < mpfFixMap {
		} else if fieldtype < mpfFixArray {
			if (g != Game{}) {
				output = append(output, g)
			}
			g = Game{}
			pos++
			iskey = true
			continue
		} else if fieldtype < mpfFixStr {
			// len := fieldtype - mpfFixArray
		} else if fieldtype < mpfNil {
			len := int(rdb[pos]) - mpfFixStr
			pos++
			value = rdb[pos : pos+len]
			pos += len
		} else if fieldtype > mpfMap32 {
		}

		switch fieldtype {
		case mpfStr8, mpfStr16, mpfStr32:
			pos++
			lenlen := fieldtype - mpfStr8 + 1
			lenhex := fmt.Sprintf("%x", string(rdb[pos:pos+lenlen]))
			i64, _ := strconv.ParseInt(lenhex, 16, 32)
			len := int(i64)
			pos += lenlen
			value = rdb[pos : pos+len]
			pos += len
		case mpfUint8, mpfUint16, mpfUint32, mpfUint64:
			pow := float64(rdb[pos]) - 0xC9
			len := int(math.Pow(2, pow)) / 8
			pos++
			value = rdb[pos : pos+len]
			pos += len
		case mpfBin8, mpfBin16, mpfBin32:
			pos++
			len := int(rdb[pos])
			pos++
			value = rdb[pos : pos+len]
			pos += len
		case mpfMap16, mpfMap32:
			len := 2
			if int(rdb[pos]) == mpfMap32 {
				len = 4
			}
			pos++
			value = rdb[pos : pos+len]
			pos += len
			iskey = true
			continue
		}

		if iskey {
			key = string(value[:])
		} else {
			setField(&g, key, string(value[:]))
		}
		iskey = !iskey
	}

	return output
}

// Find loops over the RDBs in the DB and concurrently matches CRC32 checksums.
func (db *DB) Find(rompath string, romname string, CRC32 uint32) {
	var wg sync.WaitGroup
	wg.Add(len(*db))
	// For every RDB in the DB
	for system, rdb := range *db {
		go func(rdb RDB, CRC32 uint32, system string) {
			// For each game in the RDB
			for _, game := range rdb {
				// If the checksums match
				if CRC32 == game.CRC32 {
					// Write the playlist entry
					// writePlaylistEntry(rompath, romname, game.Name, CRC32, system)
					fmt.Println(rompath, romname, game.Name, CRC32, system)
				}
			}
			wg.Done()
		}(rdb, CRC32, system)
	}
	// Synchronize all the goroutines
	wg.Wait()
}

// Scan scans a list of roms against the database
func Scan(roms []string, cb func(rompath string, romname string, CRC32 uint32)) {
	for _, f := range roms {
		ext := filepath.Ext(f)
		switch ext {
		case ".zip":
			// Open the ZIP archive
			z, _ := zip.OpenReader(f)
			for _, rom := range z.File {
				if rom.CRC32 > 0 {
					// Look for a matching game entry in the database
					cb(f, rom.Name, rom.CRC32)
				}
			}
			z.Close()
		}
	}
}

// LoadDB loops over the RDBs in a given directory and parses them
func LoadDB(dir string) (DB, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return DB{}, err
	}

	db := make(DB)
	for _, f := range files {
		filename := f.Name()
		system := filename[0 : len(filename)-4]
		bytes, _ := ioutil.ReadFile(dir + f.Name())
		db[system] = ParseRDB(bytes)
	}

	return db, nil
}
