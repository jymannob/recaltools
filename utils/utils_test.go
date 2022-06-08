package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func Test_fileExists(t *testing.T) {
	dir, _ := os.Getwd()
	dir = filepath.Join(dir, "..")

	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"directory", args{dir + "/utils"}, true},
		{"directory not exist", args{dir + "/utils/notExistingDir"}, false},
		{"file", args{dir + "/utils/utils.go"}, true},
		{"fileNotExist", args{dir + "utils/notExistingFile.txt"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileExists(tt.args.path); got != tt.want {
				t.Errorf("fileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dirExists(t *testing.T) {
	dir, _ := os.Getwd()
	dir = filepath.Join(dir, "..")

	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"directory", args{dir + "/utils"}, true},
		{"directory not exist", args{dir + "/utils/notExistingDir"}, false},
		{"file", args{dir + "/utils/utils.go"}, false},
		{"fileNotExist", args{dir + "utils/notExistingFile.txt"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dirExists(tt.args.path); got != tt.want {
				t.Errorf("dirExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoveFile(t *testing.T) {
	t.Parallel()
	dir, _ := os.Getwd()
	testdir := filepath.Join(dir, "../testdata")

	// 1st test : Moving File
	t1MoveFile := generateTempFile(testdir, t)
	defer os.Remove(t1MoveFile.Name())
	t1DestMovedFile := filepath.Join(testdir, "movedFile")
	defer os.Remove(t1DestMovedFile)

	// 2nd test : Moving file on existing file
	t2MoveFile := generateTempFile(testdir, t)
	defer os.Remove(t2MoveFile.Name())
	t2ExistingFile := generateTempFile(testdir, t)
	defer os.Remove(t2ExistingFile.Name())

	type args struct {
		from string
		to   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Moving file", args{t1MoveFile.Name(), t1DestMovedFile}, false},
		{"Moving file on existing file", args{t2MoveFile.Name(), t2ExistingFile.Name()}, false},
		{"File not exist", args{"fileNotExist", "fileNotExistToo"}, true},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if err := MoveFile(tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("MoveFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// It creates a temporary file in the given directory and returns the file
// do not forget to Delete It.
func generateTempFile(dir string, t *testing.T) (f *os.File) {
	tempfile, err := ioutil.TempFile(dir, "test_file")
	if err != nil {
		t.Fatal(err)
	}
	defer tempfile.Close()
	content := []byte("temporary file's content")
	if _, err := tempfile.Write(content); err != nil {
		log.Fatal(err)
	}

	return tempfile
}

func TestDeleteFile(t *testing.T) {
	t.Parallel()
	dir, _ := os.Getwd()
	testdir := filepath.Join(dir, "../testdata")

	// 1st test : Delete File
	t1DeleteFile := generateTempFile(testdir, t)
	defer os.Remove(t1DeleteFile.Name())

	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Delete file", args{t1DeleteFile.Name()}, false},
		{"File not exist", args{"filenotExist"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteFile(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadJsonFile(t *testing.T) {
	var data interface{}
	type args struct {
		fPath string
		data  interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Read good Json File", args{"../testdata/roms/megadrive/gamelist-backup.json", &data}, false},
		{"Read bad File (xml)", args{"../testdata/roms/megadrive/gamelist.xml", &data}, true},
		{"File not exist", args{"../testdata/roms/megadrive/notExist.file", &data}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadJsonFile(tt.args.fPath, &tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ReadJsonFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriteJsonFile(t *testing.T) {
	dir, _ := os.Getwd()
	testdir := filepath.Join(dir, "../testdata")
	// good data mock
	data := struct {
		String    string  `json:"string"`
		Bool      bool    `json:"bool"`
		Integer   int     `json:"integer"`
		Float     float32 `json:"float"`
		Null      interface{}
		isPrivate string // not in json (private)
	}{
		"testRomPath",
		true,
		-5,
		3.14,
		nil,
		"not in json",
	}

	t2Data := make(chan int)

	// 1st test : Write Json File
	t1JsonFile := generateTempFile(testdir, t)
	defer os.Remove(t1JsonFile.Name())
	// 2nd test : Bad data
	t2JsonFile := generateTempFile(testdir, t)
	defer os.Remove(t2JsonFile.Name())

	type args struct {
		fPath  string
		data   interface{}
		indent bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Write Json File", args{t1JsonFile.Name(), data, true}, false},
		{"Bad data", args{t2JsonFile.Name(), t2Data, true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteJsonFile(tt.args.fPath, tt.args.data, tt.args.indent); (err != nil) != tt.wantErr {
				t.Errorf("WriteJsonFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_bToMb(t *testing.T) {
	type args struct {
		b uint64
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{"Ok", args{625 * 1024 * 1024}, 625},
		{"Ok", args{0}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bToMb(tt.args.b); got != tt.want {
				t.Errorf("bToMb() = %v, want %v", got, tt.want)
			}
		})
	}
}
