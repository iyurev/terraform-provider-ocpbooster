package booster

import (
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("Extract openshift-install")
	if err := ExtractInstaller(boosterTestsInstallerPath()); err != nil {
		log.Fatal(err)
	}
	m.Run()

}
