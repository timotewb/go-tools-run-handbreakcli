package app

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Encode processes files in the specified input directory, encoding them according to the specified output directory.
// It supports encoding of ISO files to extract and encode individual titles, as well as encoding of other video files.
// If the encoding of titles is successful the input file will be moved to a `Complete` folder which will be created
// inside the input directory if it does not already exist.
// The encoding process uses HandBrakeCLI via a flatpak command for compatibility and performance.
//
// Parameters:
//   - inDir: The path to the input directory where the files to be encoded are stored.
//   - outDir: The path to the output directory where the encoded files will be stored.
//
// Returns:
//   None. The function performs the encoding process and moves the original files to a "Complete" directory within the same directory as the file.

func Encode(inDir, outDir string, ignore []string, drun bool){

	var cmdStr string
	var tidyCheck bool

	// Get a list of files in the input directory
	files, err := os.ReadDir(inDir)
	if err != nil {
		log.Fatal(err)
	}

	// Process each file
	for _, file := range files {
		// Check file is a file
		if file.Type().IsRegular() {
			// Set filename vars
			fileNameWithoutExt := strings.TrimSuffix(filepath.Base(file.Name()), filepath.Ext(file.Name()))
			inFile := filepath.Join(inDir, file.Name())
			ext := strings.ToLower(filepath.Ext(file.Name()))

			// Search for titles in iso files and for each title encode.
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
					cmdStr = fmt.Sprintf("flatpak run --command=HandBrakeCLI fr.handbrake.ghb -v 1 -i \"%s\" -o \"%s\" -Z \"Apple 2160p60 4K HEVC Surround\" --encoder nvenc_h265 --title %d", inFile, outFile, i)

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
							fmt.Printf("Failed to execute command: %v\n", err)
							tidyCheck = false
						}
					}
				}
			// Encode video files
			} else if !contains(ignore, file.Name()) {

				tidyCheck = true
				fmt.Print("\n\n----------------------------------------------------------------------------------------\n")
				fmt.Printf("Running file %s\n\n", file.Name())
				outFile := filepath.Join(outDir, fileNameWithoutExt+".m4v")
				cmdStr = fmt.Sprintf("flatpak run --command=HandBrakeCLI fr.handbrake.ghb -i \"%s\" -o \"%s\" -Z \"Apple 2160p60 4K HEVC Surround\" --encoder nvenc_h265", inFile, outFile)

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
						log.Fatalf("Failed to execute command: %v\n", err)
						tidyCheck = false
					}
				}
			}
			if tidyCheck {
				tidy(inFile)
			}
			fmt.Print("----------------------------------------------------------------------------------------\n\n")
		}
	}
}

// tidy moves a file to a "Complete" directory within the same directory as the file.
// If the "Complete" directory does not exist, it creates it.
//
// Parameters:
//   - fn: The path to the file to be moved.
//
// Returns:
//   None.
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

// contains checks if a given string is present in a slice of strings.
//
// Parameters:
//   - list: A slice of strings to search through.
//   - target: The string to search for in the list.
//
// Returns:
//   - bool: True if the target string is found in the list, false otherwise.
func contains(list []string, target string) bool {
	for _, str := range list {
		if str == target {
			return true
		}
	}
	return false
}
