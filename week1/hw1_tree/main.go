package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func printTree(out io.Writer, path string, depth int, openDirs map[int]bool, printFiles bool, isLast bool) {
	file, _ := os.Open(path)
	fileInfo, _ := os.Lstat(path)

	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		{
			dirContents, _ := file.Readdir(0)

			sort.Slice(dirContents, func(i, j int) bool {
				return dirContents[i].Name() < dirContents[j].Name()
			})

			if depth > 0 {
				_, dirName := filepath.Split(file.Name())
				indent := ""

				for i := 0; i < depth-1; i++ {
					if openDirs[i] {
						indent += "│\t"
					} else {
						indent += "\t"
					}
				}

				fmt.Fprint(out, indent)
				if isLast {
					fmt.Fprint(out, "└───")
				} else {
					fmt.Fprint(out, "├───")
				}

				fmt.Fprintln(out, dirName)
			}

			openDirs[depth] = true
			last := 0

			if printFiles {
				last = len(dirContents) - 1
			} else {
				for i := range dirContents {
					tmp, _ := os.Lstat(path + "/" + dirContents[i].Name())
					if tmp.IsDir() {
						last = i
					}
				}
			}

			for i := range dirContents {
				name := dirContents[i].Name()

				isLast := i == last
				if isLast {
					delete(openDirs, depth)
				}

				printTree(out, path+"/"+name, depth+1, openDirs, printFiles, isLast)
			}
		}

	case mode.IsRegular():
		if printFiles {
			_, relName := filepath.Split(file.Name())

			tabStr := ""

			for i := 0; i < depth-1; i++ {
				if openDirs[i] {
					tabStr += "│\t"
				} else {
					tabStr += "\t"
				}
			}

			fmt.Fprint(out, tabStr)

			if isLast {
				fmt.Fprint(out, "└───")
			} else {
				fmt.Fprint(out, "├───")
			}

			fmt.Fprint(out, relName)

			size := fileInfo.Size()
			if size == 0 {
				fmt.Fprintf(out, " (empty)\n")
			} else {
				fmt.Fprintf(out, " (%db)\n", size)
			}
		}
	}
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	mainDir, err := os.Lstat(path)

	if err != nil {
		return err
	}

	if mainDir.Mode().IsDir() {
		openDirs := map[int]bool{}
		printTree(out, path, 0, openDirs, printFiles, false)
	}

	return nil
}

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
