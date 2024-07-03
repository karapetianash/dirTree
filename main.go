package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	return walkDir(out, path, printFiles, "")
}

func walkDir(out io.Writer, path string, printFiles bool, prefix string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	files, err := f.Readdir(0)
	if err != nil {
		return err
	}

	var entries []os.FileInfo
	if printFiles {
		entries = files
	} else {
		for _, file := range files {
			if file.IsDir() {
				entries = append(entries, file)
			}
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for i, entry := range entries {
		isLast := i == len(entries)-1
		if entry.IsDir() {
			fmt.Fprintf(out, "%s%s───%s\n", prefix, getPrefix(isLast), entry.Name())
			newPrefix := prefix + getNewPrefix(isLast)
			err = walkDir(out, filepath.Join(path, entry.Name()), printFiles, newPrefix)
			if err != nil {
				return err
			}
		} else if printFiles {
			size := ""
			if entry.Size() == 0 {
				size = " (empty)"
			} else {
				size = fmt.Sprintf(" (%db)", entry.Size())
			}
			fmt.Fprintf(out, "%s%s───%s%s\n", prefix, getPrefix(isLast), entry.Name(), size)
		}
	}
	return nil
}

func getPrefix(isLast bool) string {
	if isLast {
		return "└"
	}
	return "├"
}

func getNewPrefix(isLast bool) string {
	if isLast {
		return "\t"
	}
	return "│\t"
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
