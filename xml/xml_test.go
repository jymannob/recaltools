package xml

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/antchfx/xmlquery"
)

func TestOpenXml(t *testing.T) {

	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{"Open Xml", args{"../testdata/roms/testSystem/test.xml"}, false, false},
		{"Open innexistent file", args{"../testdata/roms/testSystem/notExist.xml"}, true, true},
		{"Open bad file", args{"../testdata/roms/testSystem/badFile.xml"}, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OpenXml(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got != nil) == tt.wantNil {
				t.Errorf("OpenXml() got = %v, wantNil %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewNode(t *testing.T) {
	t1Want := &xmlquery.Node{
		Type: xmlquery.ElementNode,
		Data: "test",
	}
	xmlquery.AddChild(t1Want, &xmlquery.Node{Type: xmlquery.TextNode, Data: "text"})
	type args struct {
		name string
		text string
	}
	tests := []struct {
		name string
		args args
		want *xmlquery.Node
	}{
		{"Create text node", args{"test", "text"}, t1Want},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNode(tt.args.name, tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplaceChildNode(t *testing.T) {
	// test 1
	t1New := &xmlquery.Node{
		Type: xmlquery.ElementNode,
		Data: "test",
	}
	xmlquery.AddChild(t1New, &xmlquery.Node{Type: xmlquery.TextNode, Data: "text"})
	t1Data := &xmlquery.Node{
		Type: xmlquery.ElementNode,
		Data: "test1",
	}
	xmlquery.AddChild(t1Data, &xmlquery.Node{Type: xmlquery.TextNode, Data: "text1"})
	t1Want := &xmlquery.Node{
		Type: xmlquery.ElementNode,
		Data: "test1",
	}
	xmlquery.AddChild(t1Want, &xmlquery.Node{Type: xmlquery.TextNode, Data: "text1"})
	xmlquery.AddChild(t1Want, t1New)

	// test 2
	t2Data := &xmlquery.Node{
		Type: xmlquery.ElementNode,
		Data: "test2",
	}
	xmlquery.AddChild(t2Data, &xmlquery.Node{Type: xmlquery.ElementNode, Data: "replace"})
	t2Child := &xmlquery.Node{Type: xmlquery.ElementNode, Data: "replace"}
	xmlquery.AddChild(t2Child, &xmlquery.Node{Type: xmlquery.TextNode, Data: "new data"})
	t2Want := &xmlquery.Node{
		Type: xmlquery.ElementNode,
		Data: "test2",
	}
	xmlquery.AddChild(t2Want, t2Child)

	type args struct {
		from    *xmlquery.Node
		replace *xmlquery.Node
	}
	tests := []struct {
		name string
		args args
		want *xmlquery.Node
	}{
		{"Add node", args{t1Data, t1New}, t1Want},
		{"Add and replace node", args{t2Data, t2Child}, t2Want},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReplaceChildNode(tt.args.from, tt.args.replace)
			if !reflect.DeepEqual(tt.args.from, tt.want) {
				t.Errorf("ReplaceChildNode() = %+v, want %+v", tt.args.from, tt.want)
			}
		})

	}
}

func TestWriteXml(t *testing.T) {
	// 1st test : Write Xml file
	t1XmlFile := generateTempFile("../testdata", t)
	defer os.Remove(t1XmlFile.Name())

	s := `<?xml version="1.0"?><gameList><game source="Recalbox" timestamp="0"><path>Werewolf - The Last Warrior (USA).zip</path><name>Werewolf - The Last Warrior (USA)</name></game></gameList>`
	t1Data, err := xmlquery.Parse(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	}

	// 2nd test : Write no data
	t2XmlFile := generateTempFile("../testdata", t)
	defer os.Remove(t2XmlFile.Name())

	t2Data, err := xmlquery.Parse(strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		filePath string
		data     *xmlquery.Node
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"Write Xml file", args{t1XmlFile.Name(), t1Data}, len(s), false},
		{"Write no data", args{t2XmlFile.Name(), t2Data}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WriteXml(tt.args.filePath, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WriteXml() = %v, want %v", got, tt.want)
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
