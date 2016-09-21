package main

import (
	"path/filepath"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/deckarep/golang-set"
)

func main() {
	flag.Parse()

	dir := flag.Arg(0)
	out := flag.Arg(1)

	if dir == "" {
		usage("Scan directory not provided")
	}

	if out == "" {
		usage("Output file not provided")
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		usage(err.Error())
	}
	dirLen := len(dir)

	outFile, err := os.Create(out)
	if err != nil {
		usage(err.Error())
	}
	defer outFile.Close()

	outFile.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	outFile.WriteString(`<ScanInfo volumePath="`+dir+`" volumeSize="0" freeSpace="0" scanTime="1970-01-01T00:00:00Z" fileSizeMeasure="logical">` + "\n")

	lastDirSet := mapset.NewSet()

	walkFunc := func(path string, info os.FileInfo, err error) error {
		path = path[dirLen:]
		if info.IsDir() {
			curDirSet := mapset.NewSet()
			dirParts := strings.Split(path, "/")
			_path := ""
			for _, dirName := range dirParts {
				_path += dirName + "/"
				curDirSet.Add(_path)
			}

			isSubDir := curDirSet.IsSuperset(lastDirSet)

			if !isSubDir {
				upLevels := lastDirSet.Difference(curDirSet).Cardinality()
				for i := 0; i < upLevels; i++ {
					outFile.WriteString("</Folder>\n")
				}

				fmt.Printf("%s | %s | %s | %s\n", curDirSet, lastDirSet, lastDirSet.Difference(curDirSet), upLevels)
			}
			lastDirSet = curDirSet

			outFile.WriteString(fmt.Sprintf(
				`<Folder name="%s" created="1970-01-01T00:00:00Z" modified="1970-01-01T00:00:00Z" accessed="1970-01-01T00:00:00Z">` + "\n", info.Name()))
		}
		fmt.Print(path + "\n")
		return nil
	}

	filepath.Walk(dir, walkFunc)

	upLevels := lastDirSet.Cardinality()
	for i := 0; i < upLevels; i++ {
		outFile.WriteString("</Folder>\n")
	}
	outFile.WriteString("</ScanInfo>\n")
}

func usage(message string) {
	cmd := os.Args[0]
	fmt.Print("error: " + message + "\n\nusage: " + cmd + " dir outfile\n")
	os.Exit(1)
}

