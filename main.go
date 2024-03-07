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

	// Define varaibles used in code
	var inDir string
	var outDir string
	var drun bool
	var help bool
	var err error
	var cmdStr string
	cli := true

	// List file names for the code to ignore
	ignore := []string{".DS_Store", "._.DS_Store"}

	// Define CLI flags in shrot and long form
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

	// Print the Help docuemntation to the terminal if user passes help flag
	if help {
		flag.Usage()
		return
	}

	// Set the input directory for the application to search through (non-recursive)
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

	// Set the output directory for encoded files
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

	// Encode files
	app.Encode(inDir, outDir)

	// Finish processing
	if !cli {
		zenity.Info("Files converted!",
			zenity.Title("Complete"),
			zenity.InfoIcon,
		)
	} else {
		fmt.Print("Complete.\n")
	}
}

