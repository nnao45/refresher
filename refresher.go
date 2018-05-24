package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	LIMIT      = 500
	LOG_FILE   = "./err.log"
	LOG_DIR    = "../"
	LOG_PREFIX = "PC-"
)

func addog(text string, filename string) {
	var writer *bufio.Writer
	textData := []byte(text)

	writeFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	writer = bufio.NewWriter(writeFile)
	writer.Write(textData)
	writer.Flush()
	if err != nil {
		panic(err)
	}
	defer writeFile.Close()
}

func main() {
	var files = dirwalk(LOG_DIR)

	sort.Slice(files, func(i, j int) bool {
		return files[i] > files[j]
	})

	for i, file := range files {
		if i > LIMIT {
			if err := os.RemoveAll(file); err != nil {
				addog(fmt.Sprint(err), LOG_FILE)
			}
		}
	}

	fmt.Println("\n", "Log refresh is done.")
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		addog(fmt.Sprint(err), LOG_FILE)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		if strings.Contains(file.Name(), LOG_PREFIX) {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}

	return paths
}
