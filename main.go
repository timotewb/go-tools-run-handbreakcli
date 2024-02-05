package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ncruces/zenity"
)

func main() {

	var inDir string
	var outDir string
	var help bool
	var err error
	cli := true

	flag.StringVar(&inDir, "indir", "", "Path to input directory where files are stored.")
	flag.StringVar(&inDir, "i", "", "Path to input directory where files are stored.")
	flag.StringVar(&outDir, "outdir", "", "Path to directory where converted files will be stored.")
	flag.StringVar(&outDir, "o", "", "Path to directory where converted files will be stored.")
	flag.BoolVar(&help, "help", false, "Show usage instructions")
	flag.BoolVar(&help, "h", false, "Show usage instructions (shorthand)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------------------------")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Pass both the -i and -o parameters to run in CLI mode:\n")
		fmt.Fprintln(os.Stderr, "  -i string\n\tPath to input directory where files are stored (shorthand)")
		fmt.Fprintln(os.Stderr, "  --indir string\n")
		fmt.Fprintln(os.Stderr, "  -o string\n\tPath to directory where converted files will be stored (shorthand)")
		fmt.Fprintln(os.Stderr, "  --outdir string\n\n")
		fmt.Fprintln(os.Stderr, "  -h\n\tShow usage instructions (shorthand)")
		fmt.Fprintln(os.Stderr, "  --help")
		fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------------------------\n")
	}
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if inDir == "" {
		cli = false
		inDir, err = zenity.SelectFile(
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
	}

	if outDir == "" {
		cli = false
		outDir, err = zenity.SelectFile(
			zenity.Filename(""),
			zenity.Directory(),
			zenity.DisallowEmpty(),
			zenity.Title("Select output directory"),
		)
		if err != nil {
			zenity.Error(
				err.Error(),
				zenity.Title("Error"),
				zenity.ErrorIcon,
			)
			log.Fatal(err)
		}
	}

	files, err := os.ReadDir(inDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Type().IsRegular() {
			fileNameWithoutExt := strings.TrimSuffix(filepath.Base(file.Name()), filepath.Ext(file.Name()))
			inFile := filepath.Join(inDir, file.Name())
			outFile := filepath.Join(outDir, fileNameWithoutExt+".m4v")
			cmdStr := fmt.Sprintf("flatpak run --command=HandBrakeCLI fr.handbrake.ghb -v 1 -i \"%s\" -o \"%s\" -Z \"Apple 2160p60 4K HEVC Surround\" --encoder nvenc_h265", inFile, outFile)
			fmt.Println(cmdStr)

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

	if !cli {
		zenity.Info("Files converted!",
			zenity.Title("Complete"),
			zenity.InfoIcon,
		)
	}
}
