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
		VideoFullscreen:      false,
		VideoMonitorIndex:    0,
		AudioVolume:          0.5,
		ShowHiddenFiles:      true,
		CoreForPlaylist:      coreForPlaylist,
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
