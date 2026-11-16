package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

var mu sync.Mutex

// dirExists check if file exist and is a directory
func dirExists(path string) bool {
	dirInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return dirInfo.IsDir()
}

// fileExists check if file exist
func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// bToMb takes a number of bytes and returns a number of megabytes
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// defaultFilePerm is the permission of a file created by WriteFileAtomic.
const defaultFilePerm fs.FileMode = 0664

// WriteFileAtomic writes fPath through write without ever leaving a partially
// written file: data is written to a temporary file of the same directory,
// synced, then renamed over fPath. On error the temporary file is removed and
// fPath is left untouched.
//
// The permissions of an existing fPath are kept (0664 for a new file) when the
// filesystem supports it (FAT/exFAT does not, and it is ignored there).
func WriteFileAtomic(fPath string, write func(w io.Writer) error) (err error) {
	perm := defaultFilePerm
	if info, statErr := os.Stat(fPath); statErr == nil && info.Mode().IsRegular() {
		perm = info.Mode().Perm()
	}

	dir := filepath.Dir(fPath)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(fPath)+".tmp-*")
	if err != nil {
		return fmt.Errorf("cannot create temporary file for %s: %w", fPath, err)
	}
	tmpName := tmp.Name()
	closed := false
	defer func() {
		if err != nil {
			if !closed {
				tmp.Close()
			}
			os.Remove(tmpName)
		}
	}()

	if err = write(tmp); err != nil {
		return fmt.Errorf("cannot write %s: %w", fPath, err)
	}
	if err = tmp.Sync(); err != nil {
		return fmt.Errorf("cannot sync %s: %w", tmpName, err)
	}
	closed = true
	if err = tmp.Close(); err != nil {
		return fmt.Errorf("cannot close %s: %w", tmpName, err)
	}
	// best effort: chmod is not supported on every filesystem (FAT, exFAT)
	_ = os.Chmod(tmpName, perm)

	if err = os.Rename(tmpName, fPath); err != nil {
		return fmt.Errorf("cannot replace %s: %w", fPath, err)
	}

	// best effort: persist the rename itself (not supported everywhere)
	if d, dirErr := os.Open(dir); dirErr == nil {
		_ = d.Sync()
		d.Close()
	}

	return nil
}

// WriteJsonFile encodes the data into JSON, and writes it atomically to the file
func WriteJsonFile(fPath string, data interface{}, indent bool) error {
	mu.Lock()         // Lock file
	defer mu.Unlock() // Unlock file at end

	return WriteFileAtomic(fPath, func(w io.Writer) error {
		enc := json.NewEncoder(w)

		// format json
		if indent {
			enc.SetIndent("", "  ")
		}

		if err := enc.Encode(data); err != nil {
			return fmt.Errorf("data cannot be convert to Json: %w", err)
		}
		return nil
	})
}

// ReadJsonFile reads a json file and unmarshals it into a struct
func ReadJsonFile(fPath string, data interface{}) error {

	if !fileExists(fPath) {
		return fmt.Errorf("file do not exist : %s", fPath)
	}

	f, err := os.Open(fPath)
	if err != nil {
		return fmt.Errorf("file cannot be read : %s | %v", fPath, err)
	}
	defer f.Close() // close file at end
	dec := json.NewDecoder(f)

	mu.Lock()         // Lock file
	defer mu.Unlock() // Unlock file when finish
	{
		if err := dec.Decode(data); err != nil {
			fmt.Println(err)
			return fmt.Errorf("cannot parse Json : %v", data)
		}
	}

	return nil
}

// MoveFile moves a file from one location to another
func MoveFile(from, to string) error {
	if !fileExists(from) {
		return fmt.Errorf("file do not exist : %s", from)
	}

	if fileExists(to) {
		err := os.Remove(to)
		if err != nil {
			return fmt.Errorf("file %s already exist and cannot be deleted | %v", to, err)
		}
	}

	err := os.Rename(from, to)
	if err != nil {
		return fmt.Errorf("file %s cannot be move to %s | %v", from, to, err)
	}

	return nil
}

// DeleteFile deletes a file or directory
func DeleteFile(file string) error {
	err := os.RemoveAll(file)
	if err != nil {
		return fmt.Errorf("file or directory %s cannot be deleted | %v", file, err)
	}
	return nil
}
