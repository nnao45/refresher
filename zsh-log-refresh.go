package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/codeskyblue/go-sh"
	"gopkg.in/alecthomas/kingpin.v2"
)

var version string

var (
	app    = kingpin.New("zsh-log-refresh", "A zsh-log-refresh application.")
	limit  = app.Flag("l", "log refresh limit").Default("500").Int()
	prefix = app.Flag("s", "search log name").Required().String()
)

const (
	LOG_DIR         = "~/Documents/term_logs"
	RUNTIME_LOGFILE = "goruntime_err.log"
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

func zshLogger() {
	if os.Getenv("TERM") == "screen" || os.Getenv("TERM") == "screen-256color" {
		sh.Command("tmux", "set-option", "default-terminal", "\"screen\"").Run()
		sh.Command("pipe-pane", "cat", ">>", "$LOGDIR/$LOGFILE").Run()
		sh.Command("display-message", "💾Started logging to $LOGDIR/$LOGFILE").Run()
	}
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		addog(fmt.Sprint(err), RUNTIME_LOGFILE)
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

func init() {
	var err error
	if _, err = os.Stat(LOG_DIR); err != nil {
		if err := os.MkdirAll(LOG_DIR, 0700); err != nil {
			panic(err)
		}
	}

	zshLogger()
}

func main() {
	app.HelpFlag.Short('h')
	app.Version(fmt.Sprint("zsh-log-refresh's version: ", version))
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	//
	}

	var files = dirwalk(LOG_DIR)

	sort.Slice(files, func(i, j int) bool {
		return files[i] > files[j]
	})

	for i, file := range files {
		if i > *limit {
			if err := os.RemoveAll(file); err != nil {
				addog(fmt.Sprint(err), RUNTIME_LOGFILE)
			}
		}
	}

	fmt.Println("\n", "Log refresh is done.")
}
