package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/jymannob/recaltools"
)

var (
	buildVersion string = "UNKNOWN"
	buildCommit  string = "UNKNOWN"
	buildDate    string = "UNKNOWN"
	toolName     string = "RecalTools"
)

type BackupCmd struct {
	FormatJson bool     `arg:"-f" help:"Format Json output"`
	RomsDir    []string `arg:"positional" help:"path/to/roms/dir default:/recalbox/share/roms"`
}
type RestoreCmd struct {
	RomsDir []string `arg:"positional" help:"path/to/roms/dir default:/recalbox/share/roms"`
}

type args struct {
	BackupCmd  *BackupCmd  `arg:"subcommand:backup"`
	RestoreCmd *RestoreCmd `arg:"subcommand:restore"`
	Verbose    bool        `arg:"--verbose, -v" default:"false" help:"Print debug logs"`
	Version    bool        `args:"--version" default:"false" help:"Print program Version"`
}

func main() {

	var args args
	arg.MustParse(&args)

	if args.Version {
		printVersion()
	}

	switch {
	case args.BackupCmd != nil:

		if len(args.BackupCmd.RomsDir) < 1 {
			args.BackupCmd.RomsDir = append(args.BackupCmd.RomsDir, "/recalbox/share/roms")
		}

		favBkp := recaltools.FavBackup{
			RomsDir:    args.BackupCmd.RomsDir,
			FormatJson: args.BackupCmd.FormatJson,
			Verbose:    args.Verbose,
		}
		err := favBkp.Backup()
		if err != nil {
			log.Println(err)
		}
	case args.RestoreCmd != nil:

		if len(args.RestoreCmd.RomsDir) < 1 {
			args.RestoreCmd.RomsDir = append(args.RestoreCmd.RomsDir, "/recalbox/share/roms")
		}

		favBkp := recaltools.FavBackup{
			RomsDir:    args.RestoreCmd.RomsDir,
			FormatJson: false,
			Verbose:    args.Verbose,
		}
		err := favBkp.Restore()
		if err != nil {
			log.Println(err)
		}
	}

}

// printVersion prints the tool name, build commit, build version, and build date
func printVersion() {

	fmt.Printf("%s by Jymannob\n\nBuild : %s\nVersion : %s\nDate : %s\n", toolName, buildCommit, buildVersion, buildDate)
	os.Exit(0)
}
