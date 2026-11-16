package recaltools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var cleanSystems = []string{"megadrive", "nes"}

// setupCleanRoms copies testdata into a temp dir and adds a `media` directory
// to each system that has a gamelist.xml.
func setupCleanRoms(t *testing.T) string {
	t.Helper()
	roms := copyTestdata(t)
	for _, s := range cleanSystems {
		writeFile(t, filepath.Join(roms, s, "media", "images", "game.png"), "png")
	}
	return roms
}

func assertCleaned(t *testing.T, roms string, original map[string]string, keepMedias bool) {
	t.Helper()
	for _, s := range cleanSystems {
		sys := filepath.Join(roms, s)
		if exists(filepath.Join(sys, "gamelist.xml")) {
			t.Errorf("%s/gamelist.xml should have been renamed", s)
		}
		old := filepath.Join(sys, "gamelist.old.xml")
		if !exists(old) {
			t.Errorf("%s/gamelist.old.xml missing", s)
		} else if readFile(t, old) != original[s] {
			t.Errorf("%s/gamelist.old.xml content differs from original gamelist.xml", s)
		}
		if got := exists(filepath.Join(sys, "media", "images", "game.png")); got != keepMedias {
			t.Errorf("%s/media exists = %v, want %v", s, got, keepMedias)
		}
		// files other than gamelist/media must not be touched
		if !exists(filepath.Join(sys, fileBackupName)) {
			t.Errorf("%s/%s must not be removed", s, fileBackupName)
		}
	}
	for _, f := range []string{"badFile.xml", "test.xml"} {
		if !exists(filepath.Join(roms, "testSystem", f)) {
			t.Errorf("testSystem/%s must not be removed", f)
		}
	}
}

func readGamelists(t *testing.T, roms string) map[string]string {
	t.Helper()
	m := map[string]string{}
	for _, s := range cleanSystems {
		m[s] = readFile(t, filepath.Join(roms, s, "gamelist.xml"))
	}
	return m
}

func TestCleanScrapping_Clean(t *testing.T) {
	for _, keep := range []bool{false, true} {
		name := "DeleteMedias"
		if keep {
			name = "KeepMedias"
		}
		t.Run(name, func(t *testing.T) {
			roms := setupCleanRoms(t)
			original := readGamelists(t, roms)

			cs := &CleanScrapping{RomsDir: []string{roms}, KeepMedias: keep, Verbose: true}
			if err := cs.Clean(); err != nil {
				t.Fatal(err)
			}
			assertCleaned(t, roms, original, keep)
		})
	}
}

// An existing gamelist.old.xml is replaced by the current gamelist.xml.
func TestCleanScrapping_Clean_ExistingOld(t *testing.T) {
	roms := setupCleanRoms(t)
	original := readGamelists(t, roms)
	for _, s := range cleanSystems {
		writeFile(t, filepath.Join(roms, s, "gamelist.old.xml"), "previous")
	}

	cs := &CleanScrapping{RomsDir: []string{roms}}
	if err := cs.Clean(); err != nil {
		t.Fatal(err)
	}
	assertCleaned(t, roms, original, false)
}

func TestCleanScrapping_Clean_MultipleRomsDir(t *testing.T) {
	roms1 := setupCleanRoms(t)
	roms2 := setupCleanRoms(t)
	orig1 := readGamelists(t, roms1)
	orig2 := readGamelists(t, roms2)

	logs := captureLog(t)
	cs := &CleanScrapping{RomsDir: []string{roms1, roms2}}
	if err := cs.Clean(); err != nil {
		t.Fatal(err)
	}
	assertCleaned(t, roms1, orig1, false)
	assertCleaned(t, roms2, orig2, false)

	if len(cs.Gamelists) != 4 {
		t.Errorf("Gamelists = %v, want 4 entries", cs.Gamelists)
	}
	assertUnique(t, cs.Gamelists)
	if strings.Contains(logs.String(), "do not exist") {
		t.Errorf("a gamelist was moved twice, logs:\n%s", logs.String())
	}

	// a 2nd call must not try to move the already moved gamelists
	if err := cs.Clean(); err != nil {
		t.Fatal(err)
	}
	if len(cs.Gamelists) != 0 {
		t.Errorf("after 2nd call Gamelists = %v, want none", cs.Gamelists)
	}
	if strings.Contains(logs.String(), "do not exist") {
		t.Errorf("2nd call tried to move already moved gamelists, logs:\n%s", logs.String())
	}
	assertCleaned(t, roms1, orig1, false)
	assertCleaned(t, roms2, orig2, false)
}

// Clean on a directory without gamelist must not fail nor delete anything.
func TestCleanScrapping_Clean_NoGamelist(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "media", "keep.png"), "png")

	cs := &CleanScrapping{RomsDir: []string{dir}}
	if err := cs.Clean(); err != nil {
		t.Fatal(err)
	}
	if !exists(filepath.Join(dir, "media", "keep.png")) {
		t.Error("media of a directory without gamelist.xml must not be deleted")
	}
	if _, err := os.Stat(dir); err != nil {
		t.Error(err)
	}
}
