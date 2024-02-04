package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ncruces/zenity"
)

func main() {
	inDir, err := zenity.SelectFile(
		zenity.Filename(""),
		zenity.Directory(),
		zenity.DisallowEmpty(),
		zenity.Title("Select input directory"),
	)
	if err != nil {
		zenity.Error(
			err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}

	files, err := os.ReadDir(inDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Type().IsRegular() {
			fileNameWithoutExt := strings.TrimSuffix(filepath.Base(file.Name()), filepath.Ext(file.Name()))
			cmdStr := fmt.Sprintf("flatpak run --command=HandBrakeCLI fr.handbrake.ghb -v 1 -i \"%s\" -o \"/mnt/ns00/Complete/%s\" -Z \"Apple 2160p60 4K HEVC Surround\" --encoder nvenc_h265", filepath.Join(inDir, file.Name()), fileNameWithoutExt)

			// Create a new command
			cmd := exec.Command("/bin/sh", "-c", cmdStr)

			// Connect the command's stdout to the stdout of the Go program
			cmd.Stdout = os.Stdout

			// Run the command
			err := cmd.Run()
			if err != nil {
				log.Fatalf("Failed to execute command: %v", err)
			}
		}
	}

	zenity.Info("Folders renamed!",
		zenity.Title("Complete"),
		zenity.InfoIcon,
	)
}
