package recaltools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strconv"
	"strings"
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

// errNoPath is returned by AddGame for a `game` without `path` (or with an
// empty one): such a game cannot be matched on restore.
var errNoPath = errors.New("game without path")

// AddGame adds the `game` node to the backup. Games without `path` are
// ignored and errNoPath is returned.
func (s *SystemBackup) AddGame(node *xmlquery.Node) error {

	path := node.SelectElement("path")
	if path == nil || strings.TrimSpace(path.InnerText()) == "" {
		return errNoPath
	}

	//create Game and add path from Xml
	g := Game{
		RomPath: path.InnerText(),
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
	return nil
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

	gamelists, err := findGamelists(fb.RomsDir)
	if err != nil {
		return err
	}
	fb.Gamelists = gamelists

	for _, gamelist := range fb.Gamelists {

		fb.wg.Add(1)
		go fb.backupSystem(gamelist)
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

	nodes, err := xmlquery.QueryAll(doc, gameXpath)
	if err != nil {
		log.Println(err)
		return
	}
	if len(nodes) == 0 {
		return
	}

	systemBkp := SystemBackup{
		Games: make(map[string]*Game),
	}

	for _, node := range nodes {

		if err := systemBkp.AddGame(node); err != nil {
			if fb.Verbose {
				log.Printf("%s: %v, skipped : %s\n", gamelist, err, node.OutputXML(true))
			}
			continue
		}
		if fb.Verbose {
			label := node.SelectElement("path").InnerText()
			if name := node.SelectElement("name"); name != nil {
				label = name.InnerText()
			}
			log.Printf("Backup game : %s\n", label)
		}
	}

	if fb.Verbose {
		j, _ := json.MarshalIndent(systemBkp, "", "  ")
		log.Println(string(j))
		log.Printf("Write Json file : %v", filepath.Join(systemPath, fileBackupName))
	}

	err = utils.WriteJsonFile(filepath.Join(systemPath, fileBackupName), systemBkp, fb.FormatJson)
	if err != nil {
		log.Println(err)
	}
}

// findGamelists walks every roms directory and returns the `gamelist.xml`
// files found, each one only once (even if roms directories overlap).
func findGamelists(romsDirs []string) ([]string, error) {
	var gamelists []string
	seen := make(map[string]bool)

	for _, romsdir := range romsDirs {
		err := filepath.WalkDir(romsdir, func(path string, di fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if di.IsDir() || di.Name() != "gamelist.xml" {
				return nil
			}

			key, absErr := filepath.Abs(path)
			if absErr != nil {
				key = filepath.Clean(path)
			}
			if !seen[key] {
				seen[key] = true
				gamelists = append(gamelists, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("cannot walk roms directory %s: %w", romsdir, err)
		}
	}

	return gamelists, nil
}

func (fb *FavBackup) Restore() error {

	gamelists, err := findGamelists(fb.RomsDir)
	if err != nil {
		return err
	}
	fb.Gamelists = gamelists

	for _, gamelist := range fb.Gamelists {

		fb.wg.Add(1)
		go fb.restoreSystem(gamelist)
	}

	fb.wg.Wait()
	log.Println("Restore Done !")
	return nil
}

// indexGames returns the `game` nodes of doc indexed by the exact text of
// their `path` child. When several games share the same path, the first one
// in document order is kept.
//
// Comparing in Go instead of building an XPath query avoids any escaping
// issue with paths containing `"` or `'`.
func indexGames(doc *xmlquery.Node) map[string]*xmlquery.Node {
	games := make(map[string]*xmlquery.Node)
	for _, game := range xmlquery.Find(doc, "//game") {
		path := game.SelectElement("path")
		if path == nil {
			continue
		}
		if _, found := games[path.InnerText()]; !found {
			games[path.InnerText()] = game
		}
	}
	return games
}

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

	games := indexGames(doc)

	for _, v := range backup.Games {

		if fb.Verbose {
			log.Printf("Restore game : %s \n", v.RomPath)
		}

		// get `game` Node
		a, found := games[v.RomPath]
		if !found {
			if fb.Verbose {
				log.Printf("Game not found in %s, skipped : %s\n", gamelist, v.RomPath)
			}
			continue
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
