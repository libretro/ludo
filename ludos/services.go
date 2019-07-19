// Package ludos encapsulates logic used in LudOS, the Operating System
// version of Ludo. Here Ludo is used as a Desktop Environment and can
// perform actions like rebooting the system or enabling a daemon.
package ludos

import (
	"os"
	"os/exec"

	"github.com/fatih/structs"
	ntf "github.com/libretro/ludo/notifications"
)

// exists naively checks if a file exists. We're not handling cases like, the
// file is a directory.
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// SystemdServiceToggle can enable and start, or disable and stop a systemd
// service in LudOS.
func SystemdServiceToggle(path string, serviceName string, enable bool) error {
	action := "stop"
	if enable {
		action = "start"
		if !exists(path) {
			var file, err = os.Create(path)
			if err != nil {
				return err
			}
			file.Close()
		}
	} else if exists(path) {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("/usr/bin/systemctl", action, serviceName)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// ServiceSettingIncrCallback is executed when a service settings is toggled.
// It enables or disables the daemon corresponding to the current setting
// field.
func ServiceSettingIncrCallback(f *structs.Field, direction int) {
	v := f.Value().(bool)
	v = !v
	err := SystemdServiceToggle(f.Tag("path"), f.Tag("service"), v)
	if err != nil {
		ntf.DisplayAndLog(ntf.Error, "Settings", err.Error())
	} else {
		f.Set(v)
	}
}

// InitializeServiceSettingsValues is called after settings.json is loaded.
// It sets the values of SSHService, SambaService and BluetoothService that
// don't depend on settings.json but on the presence of files in the system.
func InitializeServiceSettingsValues(fields []*structs.Field) {
	for _, f := range fields {
		switch f.Name() {
		case "SSHService":
			f.Set(exists(f.Tag("path")))
		case "SambaService":
			f.Set(exists(f.Tag("path")))
		case "BluetoothService":
			f.Set(exists(f.Tag("path")))
		}
	}
}
