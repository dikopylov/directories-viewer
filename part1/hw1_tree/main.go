package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	defaultFile = "├───"
	endFile     = "└───"
	childMargin = "│\t"
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
	Path     string
	FileInfo os.FileInfo
	Children []*Node
}

func dirTree(out io.Writer, path string, needPrintFiles bool) error {
	nodes, err := buildDirTree([]*Node{}, path, path, needPrintFiles)

	if err != nil {
		return err
	}

	draw(out, nodes, path)

	return nil
}

func buildDirTree(nodes []*Node, root string, path string, needPrintFiles bool) ([]*Node, error) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && !needPrintFiles {
			continue
		}

		isDir := file.IsDir()

		var childrensNodes []*Node = nil

		currentPath := path + string(os.PathSeparator) + file.Name()

		if isDir {
			childrensNodes, _ = buildDirTree(nodes, root, currentPath, needPrintFiles)
		}

		node := Node{
			Path:     currentPath,
			FileInfo: file,
			Children: childrensNodes,
		}

		nodes = append(nodes, &node)
	}

	return nodes, nil
}

func draw(out io.Writer, nodes []*Node, rootPath string) {
	lenNodes := len(nodes)

	for index, node := range nodes {
		path := strings.Replace(node.Path, rootPath+string(os.PathSeparator), "", 1)
		splitPath := strings.Split(path, string(os.PathSeparator))

		depth := len(splitPath)

		for i := 1; i < depth; i++ {
			fmt.Fprint(out, childMargin)
		}

		sizeInformation := ""

		if !node.FileInfo.IsDir() {
			var size string

			if node.FileInfo.Size() > 0 {
				size = strconv.Itoa(int(node.FileInfo.Size())) + "b"
			} else {
				size = "empty"
			}

			sizeInformation = " (" + size + ")"
		}

		fileInfo := node.FileInfo.Name() + sizeInformation + "\n"

		if lenNodes == index {
			fmt.Fprint(out, endFile+fileInfo)
		} else {
			fmt.Fprint(out, defaultFile+fileInfo)
		}

		if node.Children != nil {
			draw(out, node.Children, rootPath)
		}
	}
}
