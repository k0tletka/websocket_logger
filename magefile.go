//go:build mage
// +build mage

package main

import (
	"bytes"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const (
	projectName   = "github.com/k0tletka/websocket_logger"
	artifactsBin  = "log_server"
	artifactsPath = "build"
)

type IndexHTMLData struct {
	WebsocketAddress string
}

var (
	buildLocation     = filepath.Join(artifactsPath, artifactsBin)
	indexHTMLTemplate = template.Must(template.ParseFiles("./static/index.html"))
)

func Build() error {
	return sh.RunV("go", "build", "-o", buildLocation, projectName)
}

func Install() error {
	mg.Deps(Build)
	mg.Deps(InstallStaticFiles)

	if err := os.MkdirAll(installPath, 0755); err != nil {
		return err
	}

	log.Println("Copying executable file...")

	sourceExec, err := ioutil.ReadFile(buildLocation)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(installPath, artifactsBin), sourceExec, 0755)
}

func InstallStaticFiles() error {
	// Create static directory in installation path
	if err := os.MkdirAll(filepath.Join(installPath, "static"), 0755); err != nil {
		return err
	}

	log.Println("Copying index.html file...")

	indexHTMLData := &IndexHTMLData{
		WebsocketAddress: websocketAddress,
	}

	indexHTMLContent := &bytes.Buffer{}
	if err := indexHTMLTemplate.Execute(indexHTMLContent, indexHTMLData); err != nil {
		return err
	}

	// Write index.html file
	if err := ioutil.WriteFile(filepath.Join(installPath, "static/index.html"), indexHTMLContent.Bytes(), 0755); err != nil {
		return err
	}

	log.Println("Copying style.css file...")

	styleCSS, err := ioutil.ReadFile("static/style.css")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(installPath, "static/style.css"), styleCSS, 0755)
}
