// Package playlists is the playlist manager of Ludo. In Ludo, playlists are
// CSV files containing the ROM path, name, and CRC32 checksum.
// Playlists are kept into memory for fast lookup of entries and deduplication.
package playlists

import (
	"path/filepath"
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
			filepath.Join("testdata", "Sega - Master System - Mark III.csv"): Playlist{
				{
					filepath.Clean("/Users/kivutar/testroms/Sega - Master System - Mark III/Aleste (Japan).zip"),
					"Aleste (Japan)",
					3636729435,
				},
				{
					filepath.Clean("/Users/kivutar/testroms/Sega - Master System - Mark III/Alex Kidd in Miracle World (USA, Europe) (Rev 1).zip"),
					"Alex Kidd in Miracle World (USA, Europe, Brazil) (Rev 1)",
					2933500612,
				},
				{
					filepath.Clean("/Users/kivutar/testroms/Sega - Master System - Mark III/Aztec Adventure - The Golden Road to Paradise (World).zip"),
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

	t.Run("Should return the number of playlist entries", func(t *testing.T) {
		got := Count("testdata/Sega - Master System - Mark III.csv")
		want := 3
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func TestShortName(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should remove vendor",
			args: args{
				in: "Sega - 32X",
			},
			want: "32X",
		},
		{
			name: "Should specify vendor as additional information",
			args: args{
				in: "FB Alpha - Arcade Games",
			},
			want: "Arcade (FB Alpha)",
		},
		{
			name: "Should remove vendor and alternative name",
			args: args{
				in: "NEC - PC Engine - TurboGrafx 16",
			},
			want: "TurboGrafx-16",
		},
		{
			name: "Should replace with prefered name",
			args: args{
				in: "Nintendo - Super Nintendo Entertainment System",
			},
			want: "Super Nintendo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShortName(tt.args.in); got != tt.want {
				t.Errorf("ShortName() = %v, want %v", got, tt.want)
			}
		})
	}
}
