package recaltools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/antchfx/xmlquery"
	"github.com/jymannob/recaltools/xml"
)

// gameChild returns the text of the child `name` of the game whose path is
// exactly `romPath`, and whether the child exists.
func gameChild(t *testing.T, gamelist, romPath, name string) (string, bool) {
	t.Helper()
	doc, err := xml.OpenXml(gamelist)
	if err != nil {
		t.Fatal(err)
	}
	for _, g := range xmlquery.Find(doc, "//game") {
		p := g.SelectElement("path")
		if p == nil || p.InnerText() != romPath {
			continue
		}
		c := g.SelectElement(name)
		if c == nil {
			return "", false
		}
		return c.InnerText(), true
	}
	t.Fatalf("game %q not found in %s", romPath, gamelist)
	return "", false
}

func runRestoreSystem(fb *FavBackup, gamelist string) {
	fb.wg.Add(1)
	fb.restoreSystem(gamelist)
}

// A game present in the backup but missing from the gamelist must be skipped
// without panicking nor abandoning the other games of the system.
func TestRestoreSystem_MissingGame(t *testing.T) {
	dir := t.TempDir()
	gamelist := filepath.Join(dir, "gamelist.xml")
	writeFile(t, gamelist, `<?xml version="1.0"?>
<gameList>
	<game><path>present.zip</path><name>Present</name></game>
</gameList>`)
	writeFile(t, filepath.Join(dir, fileBackupName), `{"games":{
		"absent.zip":{"path":"absent.zip","favorite":true,"playcount":"1"},
		"present.zip":{"path":"present.zip","favorite":true,"playcount":"7"}
	}}`)

	logs := captureLog(t)
	fb := &FavBackup{Verbose: true}
	runRestoreSystem(fb, gamelist)

	if v, ok := gameChild(t, gamelist, "present.zip", "playcount"); !ok || v != "7" {
		t.Errorf("present.zip playcount = %q (%v), want 7", v, ok)
	}
	if v, ok := gameChild(t, gamelist, "present.zip", "favorite"); !ok || v != "true" {
		t.Errorf("present.zip favorite = %q (%v), want true", v, ok)
	}
	if !strings.Contains(logs.String(), "absent.zip") {
		t.Errorf("missing game should be logged in verbose mode, logs:\n%s", logs.String())
	}
}

// The game must be matched on its exact path, not on a substring, and paths
// containing quotes must be handled.
func TestRestoreSystem_StrictPathMatch(t *testing.T) {
	dir := t.TempDir()
	gamelist := filepath.Join(dir, "gamelist.xml")
	writeFile(t, gamelist, `<?xml version="1.0"?>
<gameList>
	<game><path>Super Game.zip</path><name>Super Game</name></game>
	<game><path>Game.zip</path><name>Game</name></game>
	<game><path>Bob's "Quest".zip</path><name>Bob's Quest</name></game>
	<game><path>It's.zip</path><name>Its</name></game>
</gameList>`)
	writeFile(t, filepath.Join(dir, fileBackupName), `{"games":{
		"Game.zip":{"path":"Game.zip","playcount":"3"},
		"Bob's \"Quest\".zip":{"path":"Bob's \"Quest\".zip","playcount":"5"},
		"It's.zip":{"path":"It's.zip","favorite":true}
	}}`)

	fb := &FavBackup{}
	runRestoreSystem(fb, gamelist)

	if v, ok := gameChild(t, gamelist, "Super Game.zip", "playcount"); ok {
		t.Errorf("Super Game.zip must not be modified, got playcount %q", v)
	}
	if v, ok := gameChild(t, gamelist, "Game.zip", "playcount"); !ok || v != "3" {
		t.Errorf("Game.zip playcount = %q (%v), want 3", v, ok)
	}
	if v, ok := gameChild(t, gamelist, `Bob's "Quest".zip`, "playcount"); !ok || v != "5" {
		t.Errorf(`Bob's "Quest".zip playcount = %q (%v), want 5`, v, ok)
	}
	if v, ok := gameChild(t, gamelist, "It's.zip", "favorite"); !ok || v != "true" {
		t.Errorf("It's.zip favorite = %q (%v), want true", v, ok)
	}
}

// With several RomsDir, each gamelist must be processed exactly once, and
// calling the method twice must not reprocess the gamelists of the 1st call.
func TestFavBackup_MultipleRomsDir(t *testing.T) {
	dir1 := copyTestdata(t)
	dir2 := copyTestdata(t)

	t.Run("Backup", func(t *testing.T) {
		logs := captureLog(t)
		fb := &FavBackup{RomsDir: []string{dir1, dir2}, Verbose: true}
		if err := fb.Backup(); err != nil {
			t.Fatal(err)
		}
		if len(fb.Gamelists) != 4 {
			t.Errorf("Gamelists = %v, want 4 entries", fb.Gamelists)
		}
		assertUnique(t, fb.Gamelists)
		assertProcessedOnce(t, logs.String(), fb.Gamelists, "%s Found\n")

		if err := fb.Backup(); err != nil {
			t.Fatal(err)
		}
		if len(fb.Gamelists) != 4 {
			t.Errorf("after 2nd call Gamelists = %v, want 4 entries", fb.Gamelists)
		}
		assertUnique(t, fb.Gamelists)
	})

	t.Run("Restore", func(t *testing.T) {
		logs := captureLog(t)
		fb := &FavBackup{RomsDir: []string{dir1, dir2}, Verbose: true}
		if err := fb.Restore(); err != nil {
			t.Fatal(err)
		}
		if len(fb.Gamelists) != 4 {
			t.Errorf("Gamelists = %v, want 4 entries", fb.Gamelists)
		}
		assertUnique(t, fb.Gamelists)
		assertProcessedOnce(t, logs.String(), fb.Gamelists, "Write Xml file : %s\n")

		if err := fb.Restore(); err != nil {
			t.Fatal(err)
		}
		if len(fb.Gamelists) != 4 {
			t.Errorf("after 2nd call Gamelists = %v, want 4 entries", fb.Gamelists)
		}
		assertUnique(t, fb.Gamelists)
	})

	t.Run("Overlapping RomsDir", func(t *testing.T) {
		fb := &FavBackup{RomsDir: []string{dir1, filepath.Join(dir1, "nes")}}
		if err := fb.Backup(); err != nil {
			t.Fatal(err)
		}
		assertUnique(t, fb.Gamelists)
	})
}

// The error returned by WriteJsonFile must be reported.
func TestBackupSystem_WriteJsonError(t *testing.T) {
	roms := copyTestdata(t)
	system := filepath.Join(roms, "nes")
	backupFile := filepath.Join(system, fileBackupName)
	// make the backup file impossible to write: replace it by a directory
	if err := os.Remove(backupFile); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(backupFile, 0755); err != nil {
		t.Fatal(err)
	}

	logs := captureLog(t)
	fb := &FavBackup{}
	fb.wg.Add(1)
	fb.backupSystem(filepath.Join(system, "gamelist.xml"))

	if !strings.Contains(logs.String(), backupFile) {
		t.Errorf("write error on %s should be logged, logs:\n%s", backupFile, logs.String())
	}
}

// assertProcessedOnce checks that the verbose log line `format` (applied to
// each gamelist) appears exactly once per gamelist.
func assertProcessedOnce(t *testing.T, logs string, gamelists []string, format string) {
	t.Helper()
	for _, g := range gamelists {
		if n := strings.Count(logs, fmt.Sprintf(format, g)); n != 1 {
			t.Errorf("%s processed %d times, want 1", g, n)
		}
	}
}
