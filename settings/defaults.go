package settings

import (
	"log"
	"os/user"
	"path/filepath"
)

func defaultSettings() Settings {
	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	return Settings{
		VideoFullscreen:   false,
		VideoMonitorIndex: 0,
		VideoFilter:       "sharp-bilinear",
		AudioVolume:       0.5,
		ShowHiddenFiles:   false,
		GLVersion:         "3.2",
		CoreForPlaylist: map[string]string{
			"Atari - 2600":                                   "stella_libretro",
			"Atari - 5200":                                   "atari800_libretro",
			"Atari - 7800":                                   "prosystem_libretro",
			"Atari - Jaguar":                                 "virtualjaguar_libretro",
			"Atari - Lynx":                                   "handy_libretro",
			"Atari - ST":                                     "hatari_libretro",
			"Bandai - WonderSwan Color":                      "mednafen_wswan_libretro",
			"Bandai - WonderSwan":                            "mednafen_wswan_libretro",
			"Cave Story":                                     "nxengine_libretro",
			"ChaiLove":                                       "chailove_libretro",
			"Coleco - ColecoVision":                          "bluemsx_libretro",
			"FB Alpha - Arcade Games":                        "fbneo_libretro",
			"GCE - Vectrex":                                  "vecx_libretro",
			"Magnavox - Odyssey2":                            "o2em_libretro",
			"Microsoft - MSX":                                "bluemsx_libretro",
			"Microsoft - MSX2":                               "bluemsx_libretro",
			"NEC - PC Engine SuperGrafx":                     "mednafen_supergrafx_libretro",
			"NEC - PC Engine - TurboGrafx 16":                "mednafen_pce_fast_libretro",
			"Nintendo - Family Computer Disk System":         "fceumm_libretro",
			"Nintendo - Game Boy Advance":                    "mgba_libretro",
			"Nintendo - Game Boy Color":                      "gambatte_libretro",
			"Nintendo - Game Boy":                            "gambatte_libretro",
			"Nintendo - Nintendo Entertainment System":       "fceumm_libretro",
			"Nintendo - Pokemon Mini":                        "pokemini_libretro",
			"Nintendo - Super Nintendo Entertainment System": "snes9x_libretro",
			"Nintendo - Virtual Boy":                         "mednafen_vb_libretro",
			"Sega - 32X":                                     "picodrive_libretro",
			"Sega - Game Gear":                               "genesis_plus_gx_libretro",
			"Sega - Master System - Mark III":                "genesis_plus_gx_libretro",
			"Sega - Mega Drive - Genesis":                    "genesis_plus_gx_libretro",
			"Sega - PICO":                                    "picodrive_libretro",
			"Sega - Saturn":                                  "mednafen_saturn_libretro",
			"Sega - SG-1000":                                 "genesis_plus_gx_libretro",
			"SNK - Neo Geo Pocket Color":                     "mednafen_ngp_libretro",
			"SNK - Neo Geo Pocket":                           "mednafen_ngp_libretro",
			"Sony - PlayStation":                             playstationCore,
		},
		CoresDirectory:       "./cores",
		AssetsDirectory:      "./assets",
		DatabaseDirectory:    "./database",
		SavestatesDirectory:  filepath.Join(usr.HomeDir, ".ludo", "savestates"),
		SavefilesDirectory:   filepath.Join(usr.HomeDir, ".ludo", "savefiles"),
		ScreenshotsDirectory: filepath.Join(usr.HomeDir, ".ludo", "screenshots"),
		SystemDirectory:      filepath.Join(usr.HomeDir, ".ludo", "system"),
		PlaylistsDirectory:   filepath.Join(usr.HomeDir, ".ludo", "playlists"),
		ThumbnailsDirectory:  filepath.Join(usr.HomeDir, ".ludo", "thumbnails"),
	}
}
