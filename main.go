package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ncruces/zenity"
)

func main() {

	var inDir string
	var outDir string
	var drun bool
	var help bool
	var err error
	var cmdStr string
	cli := true

	ignore := []string{".DS_Store", "._.DS_Store"}

	flag.StringVar(&inDir, "i", "", "Path to input directory where files are stored (shorthand)")
	flag.StringVar(&inDir, "indir", "", "Path to input directory where files are stored")
	flag.StringVar(&outDir, "o", "", "Path to directory where converted files will be stored (shorthand)")
	flag.StringVar(&outDir, "outdir", "", "Path to directory where converted files will be stored")
	flag.BoolVar(&drun, "d", false, "Job will run but not execute (shorthand)")
	flag.BoolVar(&drun, "dry-run", false, "Job will run but not execute")
	flag.BoolVar(&help, "h", false, "Show usage instructions (shorthand)")
	flag.BoolVar(&help, "help", false, "Show usage instructions")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------------------------")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Pass both the -i and -o parameters to run in CLI mode:\n")
		fmt.Fprintln(os.Stderr, "  -i\t\tstring\n  --indir\n")
		fmt.Fprintln(os.Stderr, "  \tPath to input directory where files are stored (shorthand)\n")
		fmt.Fprintln(os.Stderr, "  -o\t\tstring\n  --outdir\n")
		fmt.Fprintln(os.Stderr, "  \tPath to directory where converted files will be stored (shorthand)\n")
		fmt.Fprintln(os.Stderr, "  -h\t\tboolean\n  --help")
		fmt.Fprintln(os.Stderr, "  \tShow usage instructions (shorthand)\n")
		fmt.Fprintln(os.Stderr, "  -d\t\tboolean\n  --dry-run")
		fmt.Fprintln(os.Stderr, "  \tJob will run but not execute. Valid in CLI and GUI modes.")
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
			ext := strings.ToLower(filepath.Ext(file.Name()))

			if ext == ".iso" {
				fmt.Print("\n\n----------------------------------------------------------------------------------------\n")
				fmt.Printf("Running IOS file %s\n\n", file.Name())
				cmdStr = fmt.Sprintf("flatpak run --command=HandBrakeCLI fr.handbrake.ghb -i \"%s\" -t 0", inFile)

				fmt.Print("- Execute command:\n")
				fmt.Printf("\t%s\n\n", cmdStr)

				cmd := exec.Command("/bin/sh", "-c", cmdStr)

				var buf bytes.Buffer
				cmd.Stdout = &buf
				cmd.Stderr = &buf

				// Run the command
				err := cmd.Run()
				if err != nil {
					fmt.Printf("Failed to execute command: %v", err)
				}

				output := buf.String()
				// Define the regex pattern
				pattern := regexp.MustCompile(`\+ title [0-9]+:`)

				// Find all matches of the pattern
				matches := pattern.FindAllString(output, -1)

				// Count the number of matches
				count := len(matches)
				fmt.Printf("- Found %d valid titles.\n", count)

				for i := 1; i <= count; i++ {

					outFile := filepath.Join(outDir, fmt.Sprintf("%s%02d.m4v", fileNameWithoutExt, i))
					cmdStr = fmt.Sprintf("flatpak run --command=HandBrakeCLI fr.handbrake.ghb -v 1 -i \"%s\" -o \"%s\" -Z \"Apple 2160p60 4K HEVC Surround\" --encoder nvenc_h265", inFile, outFile)

					fmt.Print("  - Execute command:\n")
					fmt.Printf("\t%s\n\n", cmdStr)

					if !drun {
						// Create a new command
						cmd := exec.Command("/bin/sh", "-c", cmdStr)

						// Connect the command's stdout to the stdout of the Go program
						cmd.Stdout = os.Stdout

						// Run the command
						err := cmd.Run()
						if err != nil {
							fmt.Printf("Failed to execute command: %v", err)
						} else {
							tidy(inFile)
						}
					}
				}

			} else if !contains(ignore, file.Name()) {
				fmt.Print("\n\n----------------------------------------------------------------------------------------\n")
				fmt.Printf("Running file %s\n\n", file.Name())
				outFile := filepath.Join(outDir, fileNameWithoutExt+".m4v")
				cmdStr = fmt.Sprintf("flatpak run --command=HandBrakeCLI fr.handbrake.ghb -v 1 -i \"%s\" -o \"%s\" -Z \"Apple 2160p60 4K HEVC Surround\" --encoder nvenc_h265", inFile, outFile)

				fmt.Print("  - Execute command:\n")
				fmt.Printf("\t%s\n\n", cmdStr)

				if !drun {
					// Create a new command
					cmd := exec.Command("/bin/sh", "-c", cmdStr)

					// Connect the command's stdout to the stdout of the Go program
					cmd.Stdout = os.Stdout

					// Run the command
					err := cmd.Run()
					if err != nil {
						log.Fatalf("Failed to execute command: %v", err)
					} else {
						tidy(inFile)
					}
				}
			}
			fmt.Print("----------------------------------------------------------------------------------------\n\n")
		}
	}

	if !cli {
		zenity.Info("Files converted!",
			zenity.Title("Complete"),
			zenity.InfoIcon,
		)
	}
}

func tidy(fn string) {
	dirPath := filepath.Join(filepath.Dir(fn), "Complete")
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirPath, 0755)
		if errDir != nil {
			log.Fatal(errDir)
		}
	}

	// Move the file
	errMove := os.Rename(fn, filepath.Join(dirPath, filepath.Base(fn)))
	if errMove != nil {
		log.Fatal(errMove)
	}
}

func contains(list []string, target string) bool {
	for _, str := range list {
		if str == target {
			return true
		}
	}
	return false
}
