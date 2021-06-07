package booster

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	_ "embed"
	"errors"
	"io"
	"log"
	"os"
	"runtime"
)

//embed assets/openshift-install-linux.tar.gz
//go:embed assets/openshift-install-linux.tar.gz
var installerBibLinux []byte

//go:embed assets/openshift-install-mac.tar.gz
var installerBibMac []byte

func init() {
	switch runtime.GOOS {
	case "linux":
		installerBin = installerBibLinux
	case "darwin":
		installerBin = installerBibMac
	default:
		log.Fatalln("unsupported installer platform")
	}
}

var installerBin []byte

const (
	OCPVersion               = "4.7.11"
	InstallerFileName        = "openshift-install"
	InstallerArchiveFileName = "openshift-install-linux.tar.gz"
)

var errIsDir = errors.New("it seems like openshift-install is a directory, we won't do anything with it, replace or remove it manually")
var errEmptyPath = errors.New("empty openshift-install path")

func ExtractInstaller(extractToFilePath string) error {
	extractToFile, err := os.Open(extractToFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			extractToFile, err = os.Create(extractToFilePath)
			if err != nil {
				return err
			}
			err := extractToFile.Chmod(0755)
			if err != nil {
				return err
			}
			var installerArchive bytes.Buffer
			if _, err := installerArchive.Write(installerBin); err != nil {
				return err
			}
			openArchive, err := gzip.NewReader(&installerArchive)
			if err != nil {
				if err := os.Remove(extractToFilePath); err != nil {
					return err
				}
				return err
			}
			tarReader := tar.NewReader(openArchive)
			for {
				header, err := tarReader.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}
				if header.Name == InstallerFileName {
					_, err := io.Copy(extractToFile, tarReader)
					if err != nil {
						return err
					}
				}
			}
			defer extractToFile.Close()

		} else {
			return err
		}
	}
	extractToFileStat, err := extractToFile.Stat()
	if err != nil {
		return err
	}
	if extractToFileStat.IsDir() {
		return errIsDir
	} else {
		return nil
	}

}
