package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	sh "github.com/codeskyblue/go-sh"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version        string
	logDir         string
	runtimeLogFile string
)

var (
	app     = kingpin.New("zsh-log-refresh", "A zsh-log-refresh application.")
	limit   = app.Flag("limit", "log refresh limit").Default("500").Int()
	verbose = kingpin.Flag("verbose", "Verbose mode.").Default("false").Short('v').Bool()
)

const (
	LOG_DIR         = "/Documents/term_logs"
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

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		addog(fmt.Sprint(err, "\n"), runtimeLogFile)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		if strings.Contains(file.Name(), "zshlog_") {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}

	return paths
}

func init() {
	var err error

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	logDir = filepath.Join(usr.HomeDir, LOG_DIR)
	os.Setenv("LOGDIR", logDir)
	runtimeLogFile = filepath.Join(logDir, RUNTIME_LOGFILE)

	if _, err = os.Stat(logDir); err != nil {
		if err := os.MkdirAll(logDir, 0700); err != nil {
			panic(err)
		}
	}

	go func() {
		if _, err = os.Stat("./logger.sh"); err != nil {
			sh.Command("./logger.sh").Run()
		} else {
			sh.Command("./*/logger.sh").Run()
		}
	}()
}

func main() {
	app.HelpFlag.Short('h')
	app.Version(fmt.Sprint("zsh-log-refresh's version: ", version))
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	}

	var files = dirwalk(logDir)

	sort.Slice(files, func(i, j int) bool {
		return files[i] > files[j]
	})

	for i, file := range files {
		if i > *limit {
			if err := os.RemoveAll(file); err != nil {
				addog(fmt.Sprint(err, "\n"), runtimeLogFile)
			}
		}
	}

	if *verbose {
		fmt.Println("\n", "Log refresh is done.")
	}
}
