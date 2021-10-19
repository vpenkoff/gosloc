package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
)

const CWD = "."
const ZERO_LINES_COUNT = 0

func main() {
	files := os.Args[1:]

	if len(files) == 0 {
		files = append(files, CWD)
	}

	total := 0
	for _, file := range files {
		lines := countLines(file)
		total += lines
	}
	fmt.Printf("total: %d\n", total)
}

func countLines(fileName string) int {
	totalLines := 0

	fi, err := os.Lstat(fileName)
	if err != nil {
		fmt.Printf("[Error] stat: %s\n", err)
		return ZERO_LINES_COUNT
	}

	fileMode := fi.Mode()

	if (fileMode&fs.ModeSymlink != 0) || (fileMode&fs.ModeNamedPipe != 0) {
		fmt.Printf("[Error] %s not a file, skipping...\n", fileName)
		return ZERO_LINES_COUNT
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
		fmt.Printf("[Error] %s\n", err)
		return ZERO_LINES_COUNT
	}
	fmt.Printf("%s: %d\n", fileName, linesCount)
	return linesCount
}

func getLinesCountDir(dirName string) int {
	totalLines := 0
	dir, err := os.Open(dirName)
	if err != nil {
		fmt.Printf("[Error] os: %s\n", err)
		return ZERO_LINES_COUNT
	}

	dirEntries, err := dir.Readdirnames(-1)
	if err != nil {
		fmt.Printf("[Error] %s\n", err)
		return ZERO_LINES_COUNT
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
		return ZERO_LINES_COUNT, err
	}

	lines := bytes.Count(fileContent, []byte("\n"))
	return lines, nil
}
