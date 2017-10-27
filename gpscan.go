package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
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

	dirInfo, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			usage("Scan directory does not exist")
		} else {
			usage(err.Error())
		}
	} else {
		if !dirInfo.IsDir() {
			usage("Provided scan path is not a directory")
		}
	}

	dir, err = filepath.Abs(dir)
	if err != nil {
		usage(err.Error())
	}
	dirLen := len(dir)

	var outFile *os.File

	if out == "-" {
		outFile = os.Stdout
	} else {
		outFile, err = os.Create(out)
		if err != nil {
			usage(err.Error())
		}
		defer outFile.Close()
	}

	encoder := xml.NewEncoder(outFile)
	encoder.Indent("", "  ")

	outFile.WriteString(xml.Header)
	encoder.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "GrandPerspectiveScanDump"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "appVersion"}, Value: "1.8.1"},
			{Name: xml.Name{Local: "formatVersion"}, Value: "5"},
		},
	})
	encoder.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "ScanInfo"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "volumePath"}, Value: dir},
			{Name: xml.Name{Local: "volumeSize"}, Value: "0"},
			{Name: xml.Name{Local: "freeSpace"}, Value: "0"},
			{Name: xml.Name{Local: "scanTime"}, Value: "1970-01-01T00:00:00Z"},
			{Name: xml.Name{Local: "fileSizeMeasure"}, Value: "logical"},
		},
	})
	encoder.Flush()

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
					encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "Folder"}})
				}

			}
			lastDirSet = curDirSet

			encoder.EncodeToken(xml.StartElement{
				Name: xml.Name{Local: "Folder"},
				Attr: []xml.Attr{
					{Name: xml.Name{Local: "name"}, Value: info.Name()},
					{Name: xml.Name{Local: "created"}, Value: "1970-01-01T00:00:00Z"},
					{Name: xml.Name{Local: "modified"}, Value: "1970-01-01T00:00:00Z"},
					{Name: xml.Name{Local: "accessed"}, Value: "1970-01-01T00:00:00Z"},
				},
			})
		} else {
			encoder.EncodeToken(xml.StartElement{
				Name: xml.Name{Local: "File"},
				Attr: []xml.Attr{
					{Name: xml.Name{Local: "name"}, Value: info.Name()},
					{Name: xml.Name{Local: "size"}, Value: fmt.Sprintf("%d",info.Size())},
					{Name: xml.Name{Local: "created"}, Value: "1970-01-01T00:00:00Z"},
					{Name: xml.Name{Local: "modified"}, Value: "1970-01-01T00:00:00Z"},
					{Name: xml.Name{Local: "accessed"}, Value: "1970-01-01T00:00:00Z"},
				},
			})
			encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "File"}})
		}
		return nil
	}

	filepath.Walk(dir, walkFunc)

	upLevels := lastDirSet.Cardinality()
	for i := 0; i < upLevels; i++ {
		encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "Folder"}})
	}
	encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "ScanInfo"}})
	encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "GrandPerspectiveScanDump"}})
	encoder.Flush()
}

func usage(message string) {
	cmd := os.Args[0]
	fmt.Print("Error: " + message + "\n\nUsage: " + cmd + " dir outfile\n")
	os.Exit(1)
}
