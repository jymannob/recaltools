package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/jymannob/recaltools"
)

var (
	buildVersion string = "UNKNOWN"
	buildCommit  string = "UNKNOWN"
	buildDate    string = "UNKNOWN"
	toolName     string = "RecalTools"
)

func main() {

	favBkp := recaltools.FavBackup{}

	flag.BoolVar(&favBkp.FormatJson, "f", false, "Format Json output")
	restoreBkp := flag.Bool("R", false, "Restore backup to gamelist.xml")
	flag.BoolVar(&favBkp.Verbose, "verbose", false, "Print debug logs")
	directories := flag.String("d", "/recalbox/share/roms:/recalbox/share/externals", "Roms directory")
	v := flag.Bool("v", false, "Print version")
	h := flag.Bool("h", false, "Print this help")

	flag.Parse()

	favBkp.RomsDir = filepath.SplitList(*directories)

	if *v {
		fmt.Printf("%s by Jymannob\n\nBuild : %s\nVersion : %s\nDate : %s\n", toolName, buildCommit, buildVersion, buildDate)
		return
	}
	if *h {
		flag.PrintDefaults()
		return
	}

	if *restoreBkp {

		err := favBkp.Restore()
		if err != nil {
			log.Println(err)
		}
	} else {

		err := favBkp.Backup()
		if err != nil {
			log.Println(err)
		}
	}

}
