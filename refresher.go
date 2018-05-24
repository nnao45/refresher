package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app            = kingpin.New("zsh-log-refresh", "A zsh-log-refresh application.")
	limit          = app.Flag("l", "log refresh limit").Default("500").Int()
	prefix         = app.Flag("s", "search log name").Required().String()
	logDir         = app.Flag("s", "term log dir").Required().String()
	runtimeLogFile = app.Flag("s", "this script runtime logging file").Default("./").String()
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
	var files = dirwalk(*logDir)

	sort.Slice(files, func(i, j int) bool {
		return files[i] > files[j]
	})

	for i, file := range files {
		if i > *limit {
			if err := os.RemoveAll(file); err != nil {
				addog(fmt.Sprint(err), *runtimeLogFile)
			}
		}
	}

	fmt.Println("\n", "Log refresh is done.")
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		addog(fmt.Sprint(err), *runtimeLogFile)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		if strings.Contains(file.Name(), *prefix) {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}

	return paths
}
