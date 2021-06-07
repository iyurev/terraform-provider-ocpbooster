package config

import (
	"fmt"
	"os"
)

var (
	InstallerPath            = ""
	InstallerFileName        = "openshift-install"
	InstallerArchiveFileName = ""
)

func init() {
	installerPath := os.Getenv("INSTALLER_PATH")
	if installerPath != "" {
		InstallerPath = installerPath
	} else {
		home := os.Getenv("HOME")
		if home != "" {
			InstallerPath = fmt.Sprintf("%s/%s", home, InstallerFileName)
		}
		if home == "" {
			InstallerPath = fmt.Sprintf("%s/%s", "/tmp", InstallerFileName)
		}
	}
}
