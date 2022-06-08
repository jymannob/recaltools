package xml

import (
	"fmt"
	"os"
	"sync"

	"github.com/antchfx/xmlquery"
)

var mu sync.RWMutex

// OpenXml open Xml file and return a xmlquery.Node
func OpenXml(filePath string) (*xmlquery.Node, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("impossible d'ouvrir le fichier : %s | %v", filePath, err)
	}
	defer file.Close()

	mu.RLock()
	defer mu.RUnlock()
	doc, err := xmlquery.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("impossible de parser le xml : %s | %v", filePath, err)
	}

	return doc, nil
}

// WriteXml opens a file and writes the XML data to it
func WriteXml(filePath string, data *xmlquery.Node) (int, error) {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return 0, fmt.Errorf("impossible d'écrire dans le fichier : %s | %v", filePath, err)
	}
	defer f.Close()

	return f.WriteString(data.OutputXML(true))
}

// It creates a new XML node with the given name and text
func NewNode(name, text string) *xmlquery.Node {

	node := &xmlquery.Node{Type: xmlquery.ElementNode, Data: name}
	xmlquery.AddChild(node, &xmlquery.Node{Type: xmlquery.TextNode, Data: text})

	return node
}

// ReplaceChildNode replace a node by a new node on Xml tree"
func ReplaceChildNode(from *xmlquery.Node, replace *xmlquery.Node) {
	xpath := fmt.Sprintf("./%s", replace.Data)

	playcount, _ := xmlquery.Query(from, xpath)
	if playcount != nil {
		xmlquery.RemoveFromTree(playcount)
	}

	xmlquery.AddChild(from, replace)
}
