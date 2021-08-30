package main

import (
	"io"
	"io/ioutil"
	"os"
)

const (
	defaultFile = "├───"
	endFile     = "└───"
	childFile   = "│\t"
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
	builderDirTree, _ := buildDirTree([]*Node{}, path, path, needPrintFiles)

	if len(builderDirTree) > 2 {
		//out.Write()
	}

	//out.Write(builderDirTree)

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

func draw(nodes []Node) {
	//var result []string
	//
	//for index, node := range nodes {
	//
	//	if node
	//	//result = append(result, )
	//
	//}
}
