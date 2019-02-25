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
