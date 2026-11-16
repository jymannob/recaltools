package recaltools

import (
	"io/fs"
	"log"
	"path/filepath"
	"sync"

	"github.com/jymannob/recaltools/utils"
)

type CleanScrapping struct {
	RomsDir    []string
	Gamelists  []string
	KeepMedias bool
	Verbose    bool
	wg         sync.WaitGroup
}

// Clean all scrapped data
func (cs *CleanScrapping) Clean() error {
	for _, romsdir := range cs.RomsDir {
		err := filepath.WalkDir(romsdir, cs.PopulateGamelists)
		if err != nil {
			return err
		}

		for _, gamelist := range cs.Gamelists {

			cs.wg.Add(1)
			go cs.backupGamelist(gamelist)

			if !cs.KeepMedias {
				cs.wg.Add(1)
				go cs.cleanMedias(gamelist)
			}
		}
	}

	cs.wg.Wait()
	log.Println("Clean Done !")
	return nil
}

// PopulateGamelists that is called by the filepath.WalkDir function.
func (cs *CleanScrapping) PopulateGamelists(path string, di fs.DirEntry, err error) error {

	if filepath.Base(path) == "gamelist.xml" {
		cs.Gamelists = append(cs.Gamelists, path)
	}

	return err
}

// backupGamelist move gamelist.xml to gamelist.old.xml
func (cs *CleanScrapping) backupGamelist(gamelist string) {
	defer cs.wg.Done()
	backup := filepath.Join(filepath.Dir(gamelist), "gamelist.old.xml")

	err := utils.MoveFile(gamelist, backup)
	if err != nil {
		log.Println(err)
		return
	}
	if cs.Verbose {
		log.Printf("%s Moved to %s", gamelist, backup)
	}
}

// Deleting the media directory.
func (cs *CleanScrapping) cleanMedias(gamelist string) {
	defer cs.wg.Done()

	// todo - remove only media data referenced by `gamelist.xml`
	mediaDirectory := filepath.Join(filepath.Dir(gamelist), "media")

	err := utils.DeleteFile(mediaDirectory)
	if err != nil {
		log.Println(err)
		return
	}
	if cs.Verbose {
		log.Printf("%s Deleted\n", mediaDirectory)
	}

}
