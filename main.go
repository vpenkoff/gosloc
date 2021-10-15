package main

import (
	"io/fs"
	"log"
	"os"
    "bytes"
    "fmt"
)

var supported_file_exts = []string{".go"}

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

    if ((fileMode&fs.ModeSymlink != 0 ) || (fileMode & fs.ModeNamedPipe != 0 )) {
        log.Fatalf("%s not a file, skipping...\n", fileName)
    }

    if fileMode.IsRegular() {
        linesCount, err := getLinesCount(fileName)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%s: %d\n", fileName, linesCount)
        totalLines += linesCount
    }

    if fileMode.IsDir() {
        dir, err := os.Open(fileName)
        if err != nil {
            log.Fatal(err)
        }
        dir_entries, err := dir.Readdirnames(-1)
        if err!= nil {
            log.Fatal(err);
        }
        for _, fileInDir := range dir_entries {
            linesCount :=  countLines(fileInDir)
            totalLines += linesCount
        }
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
