package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	defaultNode   = "├───"
	endNode       = "└───"
	childMargin   = "│\t"
	space         = "\t"
	emptySizeNode = " (empty)"
	sizeNode      = " (%vb)"
	newLine       = "\n"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

type Node struct {
	path     string
	fileInfo os.FileInfo
	children []Node
}

func dirTree(out io.Writer, path string, needPrintFiles bool) error {
	nodes, err := buildDirTree([]Node{}, path, path, needPrintFiles)

	if err != nil {
		return err
	}

	drawNodes(out, nodes, needPrintFiles, "")

	return nil
}

func buildDirTree(nodes []Node, root string, path string, needPrintFiles bool) ([]Node, error) {
	nodes = nil
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && !needPrintFiles {
			continue
		}

		isDir := file.IsDir()

		var childrensNodes []Node = nil

		currentPath := path + string(os.PathSeparator) + file.Name()

		if isDir {
			childrensNodes, _ = buildDirTree(nodes, root, currentPath, needPrintFiles)
		}

		node := Node{
			path:     currentPath,
			fileInfo: file,
			children: childrensNodes,
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func drawNodes(out io.Writer, fileNodes []Node, needPrintFiles bool, levelPrefix string) {
	for key, node := range fileNodes {
		length := len(fileNodes)
		nodePathName := levelPrefix

		if length > key+1 {
			nodePathName += defaultNode
		}

		if length == key+1 {
			nodePathName += endNode
		}

		nodePathName += node.fileInfo.Name()
		nodeSize := int(node.fileInfo.Size())

		if !node.fileInfo.IsDir() {
			if needPrintFiles && nodeSize == 0 {
				nodePathName += emptySizeNode
			}

			if needPrintFiles && nodeSize > 0 {
				nodePathName += fmt.Sprintf(sizeNode, strconv.Itoa(nodeSize))
			}
		}

		nodePathName += newLine
		out.Write([]byte(nodePathName))

		if len(node.children) > 0 {
			childrenPrefix := levelPrefix

			if length > key+1 {
				childrenPrefix += childMargin
			}

			if length == key+1 {
				childrenPrefix += space
			}

			drawNodes(out, node.children, needPrintFiles, childrenPrefix)
		}
	}
}
