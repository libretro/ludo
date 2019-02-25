// Package playlists is the playlist manager of Ludo. In Ludo, playlists are
// CSV files containing the ROM path, name, and CRC32 checksum.
// Playlists are kept into memory for fast lookup of entries and deduplication.
package playlists

import (
	"reflect"
	"testing"

	"github.com/libretro/ludo/settings"
)

func TestLoad(t *testing.T) {
	settings.Current.PlaylistsDirectory = "./testdata"

	Load()

	t.Run("Should load playlists", func(t *testing.T) {
		got := Playlists
		want := map[string]Playlist{
			"testdata/Sega - Master System - Mark III.csv": Playlist{
				{
					"/Users/kivutar/testroms/Sega - Master System - Mark III/Aleste (Japan).zip",
					"Aleste (Japan)",
					3636729435,
				},
				{
					"/Users/kivutar/testroms/Sega - Master System - Mark III/Alex Kidd in Miracle World (USA, Europe) (Rev 1).zip",
					"Alex Kidd in Miracle World (USA, Europe, Brazil) (Rev 1)",
					2933500612,
				},
				{
					"/Users/kivutar/testroms/Sega - Master System - Mark III/Aztec Adventure - The Golden Road to Paradise (World).zip",
					"Aztec Adventure (World)",
					4284567219,
				},
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func TestContains(t *testing.T) {
	settings.Current.PlaylistsDirectory = "./testdata"

	Load()

	t.Run("Should find an existing entry by path", func(t *testing.T) {
		got := Contains("testdata/Sega - Master System - Mark III.csv", "/Users/kivutar/testroms/Sega - Master System - Mark III/Alex Kidd in Miracle World (USA, Europe) (Rev 1).zip", 0)
		want := true
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("Should find an existing entry by CRC", func(t *testing.T) {
		got := Contains("testdata/Sega - Master System - Mark III.csv", "", 2933500612)
		want := true
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("Should not generate false positive", func(t *testing.T) {
		got := Contains("testdata/Sega - Master System - Mark III.csv", "", 2933500613)
		want := false
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func TestCount(t *testing.T) {
	settings.Current.PlaylistsDirectory = "./testdata"

	Load()

	t.Run("Should find an existing entry by path", func(t *testing.T) {
		got := Count("testdata/Sega - Master System - Mark III.csv")
		want := 3
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}
