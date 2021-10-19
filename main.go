package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
)

const CWD = "."

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		args = append(args, CWD)
	}

	fileName := args[0]
	total := countLines(fileName)
	fmt.Printf("total: %d\n", total)
}

func countLines(fileName string) int {
	totalLines := 0

	fi, err := os.Lstat(fileName)
	if err != nil {
		log.Fatal(err)
	}

	fileMode := fi.Mode()

	if (fileMode&fs.ModeSymlink != 0) || (fileMode&fs.ModeNamedPipe != 0) {
		log.Fatalf("%s not a file, skipping...\n", fileName)
	}

	if fileMode.IsRegular() {
		lines := getLinesCountRegularFile(fileName)
		totalLines += lines
	}

	if fileMode.IsDir() {
		lines := getLinesCountDir(fileName)
		totalLines += lines
	}
	return totalLines
}

func getLinesCountRegularFile(fileName string) int {
	linesCount, err := getLinesCount(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: %d\n", fileName, linesCount)
	return linesCount
}

func getLinesCountDir(dirName string) int {
	totalLines := 0
	dir, err := os.Open(dirName)
	if err != nil {
		log.Fatal(err)
	}

	dirEntries, err := dir.Readdirnames(-1)
	if err != nil {
		log.Fatal(err)
	}

	for _, fileInDir := range dirEntries {
		fileNameInDir := path.Join(dirName, fileInDir)
		linesCount := countLines(fileNameInDir)
		totalLines += linesCount
	}
	return totalLines
}

func getLinesCount(fileName string) (int, error) {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return -1, err
	}

	lines := bytes.Count(fileContent, []byte("\n"))
	return lines, nil
}
