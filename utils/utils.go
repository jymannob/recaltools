package utils

import (
	"encoding/json"
	"fmt"
	"os"
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

// WriteJsonFile encodes the data into JSON, and writes it to the file
func WriteJsonFile(fPath string, data interface{}, indent bool) error {

	// Write Only, Create file if not exist, truncate if exist
	f, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return fmt.Errorf("file cannot be write : %s | %v", fPath, err)
	}
	defer f.Close() // close file at end
	enc := json.NewEncoder(f)

	// format json
	if indent {
		enc.SetIndent("", "  ")
	}

	mu.Lock()         // Lock file
	defer mu.Unlock() // Unlock file at end
	{
		if err := enc.Encode(data); err != nil {
			fmt.Println(err)
			return fmt.Errorf("data cannot be convert to Json : %v", data)
		}
	}

	return nil
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
