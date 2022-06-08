package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jymannob/recaltools"
)

var (
	buildVersion string            = "UNKNOWN"
	buildCommit  string            = "UNKNOWN"
	buildDate    string            = "UNKNOWN"
	toolName     string            = "RecalTools"
	toolsDesc    map[string]string = make(map[string]string)
	showVersion  bool              = false
	showHelp     bool              = false
)

func main() {
	// If no sub-command is passed set to "" for show usage.
	if len(os.Args[1:]) < 1 {
		fmt.Println("Expected Sub-command")
		os.Args = append(os.Args, "")
	}

	favBkp := recaltools.FavBackup{}

	// Creating a new flagset for sub-command backup.
	backupMode := flag.NewFlagSet("backup", flag.ExitOnError)
	toolsDesc[backupMode.Name()] = "Backup user metadata for each gamelist.xml found in paths"
	backupMode.BoolVar(&favBkp.FormatJson, "f", false, "Format Json output")
	backupMode.BoolVar(&favBkp.Verbose, "verbose", false, "Print debug logs")
	initCommonFlags(backupMode)

	// Creating a new flagset for sub-command restore.
	restoreMode := flag.NewFlagSet("restore", flag.ExitOnError)
	toolsDesc[restoreMode.Name()] = "Restore user metadata saved with backup command found in paths"
	restoreMode.BoolVar(&favBkp.Verbose, "verbose", false, "Print debug logs")
	initCommonFlags(restoreMode)

	switch os.Args[1] {

	case backupMode.Name():
		backupMode.Parse(os.Args[2:])
		testCommonFlags(backupMode)

		if len(backupMode.Args()) > 0 {
			favBkp.RomsDir = backupMode.Args()
		} else {
			// for multimount add ":/recalbox/share/externals" --- @todo check before if /recalbox/share/externals/usbX is relevant
			favBkp.RomsDir = filepath.SplitList("/recalbox/share/roms")
		}

		err := favBkp.Backup()
		if err != nil {
			log.Println(err)
		}

	case restoreMode.Name():
		restoreMode.Parse(os.Args[2:])
		testCommonFlags(restoreMode)

		if len(backupMode.Args()) > 0 {
			favBkp.RomsDir = backupMode.Args()
		} else {
			// for multimount add ":/recalbox/share/externals" --- @todo check before if /recalbox/share/externals/usbX is relevant
			favBkp.RomsDir = filepath.SplitList("/recalbox/share/roms")
		}

		err := favBkp.Restore()
		if err != nil {
			log.Println(err)
		}
	default:
		printHelp(backupMode, restoreMode)
	}

}

// printVersion prints the tool name, build commit, build version, and build date
func printVersion() {

	fmt.Printf("%s by Jymannob\n\nBuild : %s\nVersion : %s\nDate : %s\n", toolName, buildCommit, buildVersion, buildDate)
	os.Exit(0)
}

// printHelp prints the tool's name, description, and usage
func printHelp(flagsets ...*flag.FlagSet) {
	fmt.Printf("%s by Jymannob\n\n", toolName)
	for _, f := range flagsets {

		fmt.Printf("Usage : %s %s <path/to/roms/directory>\n", os.Args[0], f.Name())
		fmt.Printf("    %s\n", toolsDesc[f.Name()])
		f.PrintDefaults()
	}
	os.Exit(0)
}

// initCommonFlags initializes the common flags for FlagSet
func initCommonFlags(f *flag.FlagSet) {
	f.BoolVar(&showVersion, "v", false, "Print version")
	f.BoolVar(&showHelp, "h", false, "Print this help")
}

// testCommonFlags If the user has requested the version or help, print it and exit
func testCommonFlags(f *flag.FlagSet) {
	if showVersion {
		printVersion()
	}
	if showHelp {
		printHelp(f)
	}
}
