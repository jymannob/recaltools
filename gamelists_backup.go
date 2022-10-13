package recaltools

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/antchfx/xmlquery"
	"github.com/jymannob/recaltools/utils"
	"github.com/jymannob/recaltools/xml"
)

type FavBackup struct {
	RomsDir    []string
	Gamelists  []string
	FormatJson bool
	Verbose    bool
	RestoreBkp bool // unused in reclatools version
	wg         sync.WaitGroup
}

type SystemBackup struct {
	Games map[string]*Game `json:"games"`
}

type Game struct {
	RomPath    string `json:"path"`
	Favorite   bool   `json:"favorite,omitempty"`
	Playcount  string `json:"playcount,omitempty"`
	Lastplayed string `json:"lastplayed,omitempty"`
}

var fileBackupName string = "gamelist-backup.json"

func (s *SystemBackup) AddGame(node *xmlquery.Node) {

	//create Game and add path from Xml (always present)
	g := Game{
		RomPath: node.SelectElement("path").InnerText(),
	}

	favorite := node.SelectElement("favorite")
	if favorite != nil {
		b, err := strconv.ParseBool(favorite.InnerText())
		if err == nil {
			g.Favorite = b
		}
	}

	playcount := node.SelectElement("playcount")
	if playcount != nil {
		g.Playcount = playcount.InnerText()
	}

	lastplayed := node.SelectElement("lastplayed")
	if lastplayed != nil {
		g.Lastplayed = lastplayed.InnerText()
	}

	s.Games[g.RomPath] = &g
}

func (fb *FavBackup) PopulateGamelists(path string, di fs.DirEntry, err error) error {

	if filepath.Base(path) == "gamelist.xml" {
		fb.Gamelists = append(fb.Gamelists, path)
	}

	return err
}

//get all `game` nodes with child node (`lastplayed` OR `playcount` OR `favorite` contains "true" (insensitive))
var gameXpath string = "//game[./favorite[matches(text(), \"(?i)^true$\")]|./lastplayed|./playcount]"

func (fb *FavBackup) Backup() error {

	for _, romsdir := range fb.RomsDir {
		err := filepath.WalkDir(romsdir, fb.PopulateGamelists)
		if err != nil {
			return err
		}

		for _, gamelist := range fb.Gamelists {

			fb.wg.Add(1)
			go fb.backupSystem(gamelist)
		}
	}

	fb.wg.Wait()
	log.Println("Backup Done !")
	return nil
}

func (fb *FavBackup) backupSystem(gamelist string) {
	defer fb.wg.Done()

	systemPath := filepath.Dir(gamelist)
	if fb.Verbose {
		log.Printf("%s Found\n", gamelist)
	}

	doc, err := xml.OpenXml(gamelist)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = xmlquery.QueryAll(doc, gameXpath)
	if err != nil {
		log.Println(err)
		return
	}
	if nodes, err := xmlquery.QueryAll(doc, gameXpath); err == nil {

		if len(nodes) > 0 {

			systemBkp := SystemBackup{
				Games: make(map[string]*Game),
			}

			for _, node := range nodes {

				systemBkp.AddGame(node)
				if fb.Verbose {
					log.Printf("Backup game : %s\n", node.SelectElement("name").InnerText())
				}
			}

			if fb.Verbose {
				j, _ := json.MarshalIndent(systemBkp, "", "  ")
				log.Println(string(j))
				log.Printf("Write Json file : %v", filepath.Join(systemPath, fileBackupName))
			}

			utils.WriteJsonFile(filepath.Join(systemPath, fileBackupName), systemBkp, fb.FormatJson)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (fb *FavBackup) Restore() error {

	for _, romsdir := range fb.RomsDir {

		err := filepath.WalkDir(romsdir, fb.PopulateGamelists)
		if err != nil {
			return err
		}

		for _, gamelist := range fb.Gamelists {

			fb.wg.Add(1)
			go fb.restoreSystem(gamelist)
		}
	}

	fb.wg.Wait()
	log.Println("Restore Done !")
	return nil
}

//get `game` nodes with child node (`path` contains specific test)
var gameFindXpath string = "//game[./path[contains(text(), \"%s\")]]"

func (fb *FavBackup) restoreSystem(gamelist string) {
	defer fb.wg.Done()

	systemPath := filepath.Dir(gamelist)

	doc, err := xml.OpenXml(gamelist)
	if err != nil {
		log.Println(err)
		return
	}

	// Read Json
	var backup SystemBackup
	err = utils.ReadJsonFile(filepath.Join(systemPath, fileBackupName), &backup)
	if err != nil {
		return
	}

	if fb.Verbose {
		log.Printf("%s Found backup\n", filepath.Join(systemPath, fileBackupName))
	}

	for _, v := range backup.Games {

		if fb.Verbose {
			log.Printf("Restore game : %s \n", v.RomPath)
		}

		// get `game` Node
		a, err := xmlquery.Query(doc, fmt.Sprintf(gameFindXpath, v.RomPath))
		if err != nil {
			return // no `game` node
		}

		// update childs Nodes
		if v.Favorite {
			favorite := xml.NewNode("favorite", fmt.Sprintf("%t", v.Favorite))
			xml.ReplaceChildNode(a, favorite)
		}

		if v.Lastplayed != "" {
			lastplayed := xml.NewNode("lastplayed", v.Lastplayed)
			xml.ReplaceChildNode(a, lastplayed)
		}

		if v.Playcount != "" {
			playcount := xml.NewNode("playcount", v.Playcount)
			xml.ReplaceChildNode(a, playcount)
		}
	}

	if fb.Verbose {
		log.Printf("Write Xml file : %s", gamelist)
	}
	_, err = xml.WriteXml(gamelist, doc)
	if err != nil {
		log.Println(err)
		return
	}

}
