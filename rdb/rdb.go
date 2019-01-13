package rdb

import (
	"fmt"
	"math"
	"strconv"
	"sync"
)

// DB is a database that contains many RDB, mapped to their system name
type DB map[string]RDB

// RDB contains all the game descriptions for a system
type RDB []Entry

// Entry represents a game in the libretro database
type Entry struct {
	Path        string
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
	System      string
}

const (
	mpfFixMap   = 0x80
	mpfMap16    = 0xde
	mpfMap32    = 0xdf
	mpfFixArray = 0x90
	// mpfArray16  = 0xdc
	// mpfArray32 = 0xdd
	mpfFixStr = 0xa0
	mpfStr8   = 0xd9
	mpfStr16  = 0xda
	mpfStr32  = 0xdb
	mpfBin8   = 0xc4
	mpfBin16  = 0xc5
	mpfBin32  = 0xc6
	// mpfFalse = 0xc2
	// mpfTrue = 0xc3
	// mpfInt8  = 0xd0
	// mpfInt16 = 0xd1
	// mpfInt32 = 0xd2
	// mpfInt64 = 0xd3
	mpfUint8  = 0xcc
	mpfUint16 = 0xcd
	mpfUint32 = 0xce
	mpfUint64 = 0xcf
	mpfNil    = 0xc0
)

// SetField sets a field in the entry
func (g *Entry) SetField(key string, value string) {
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

// Parse parses a .rdb file content and returns an array of Entries
func Parse(rdb []byte) RDB {
	var output RDB
	pos := 0x10
	iskey := false
	key := ""
	g := Entry{}
	for int(rdb[pos]) != mpfNil {
		fieldtype := int(rdb[pos])
		var value []byte
		if fieldtype < mpfFixMap {
		} else if fieldtype < mpfFixArray {
			if (g != Entry{}) {
				output = append(output, g)
			}
			g = Entry{}
			pos++
			iskey = true
			continue
			// } else if fieldtype < mpfFixStr {
			// 	len := fieldtype - mpfFixArray
		} else if fieldtype < mpfNil {
			len := int(rdb[pos]) - mpfFixStr
			pos++
			value = rdb[pos : pos+len]
			pos += len
		}
		// else if fieldtype > mpfMap32 {
		// }
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
		}
		if iskey {
			key = string(value[:])
		} else {
			g.SetField(key, string(value[:]))
		}
		iskey = !iskey
	}
	// Don't forget to add the last rdb entry
	if (g != Entry{}) {
		output = append(output, g)
	}
	return output
}

// FindByCRC loops over the RDBs in the DB and concurrently matches CRC32 checksums.
func (db *DB) FindByCRC(romPath string, romName string, CRC32 uint32, games chan (Entry)) {
	var wg sync.WaitGroup
	wg.Add(len(*db))
	// For every RDB in the DB
	for system, rdb := range *db {
		go func(rdb RDB, CRC32 uint32, system string) {
			// For each game in the RDB
			for _, game := range rdb {
				// If the checksums match
				if CRC32 == game.CRC32 {
					games <- Entry{Path: romPath, ROMName: romName, Name: game.Name, CRC32: CRC32, System: system}
				}
			}
			wg.Done()
		}(rdb, CRC32, system)
	}
	// Synchronize all the goroutines
	wg.Wait()
}

// FindByROMName loops over the RDBs in the DB and concurrently matches ROM names.
func (db *DB) FindByROMName(romPath string, romName string, CRC32 uint32, games chan (Entry)) {
	var wg sync.WaitGroup
	wg.Add(len(*db))
	// For every RDB in the DB
	for system, rdb := range *db {
		go func(rdb RDB, CRC32 uint32, system string) {
			// For each game in the RDB
			for _, game := range rdb {
				// If the checksums match
				if romName == game.ROMName {
					games <- Entry{Path: romPath, ROMName: romName, Name: game.Name, CRC32: CRC32, System: system}
				}
			}
			wg.Done()
		}(rdb, CRC32, system)
	}
	// Synchronize all the goroutines
	wg.Wait()
}
