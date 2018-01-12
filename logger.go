/*
Package logger provides an additional layer of interfacing on top of the Go
standard log package, to provide several features such as indentation and
consistent default log file formatting.
*/
package logger

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

var (
	enabled  bool = false
	file     *os.File
	logger   *log.Logger
	tabLevel int
)

func Init() {
	if enabled {
		return
	}

	usr, err := user.Current()
	if err != nil {
		fmt.Println("Logger couldn't get user")
		return
	}

	t := time.Now()

	filename := fmt.Sprintf("hanoi_%04d-%02d-%02d_%02d-%02d-%02d.log",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	err = os.Mkdir(filepath.Join(usr.HomeDir, "logs"), 0777)
	if !os.IsExist(err) {
		return
	}

	file, err = os.Create(filepath.Join(usr.HomeDir, "logs", filename))
	if err != nil {
		fmt.Println("Logger couldn't create file")
		return
	}

	logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	enabled = true
}

func IncTab() {
	SetTab(tabLevel + 1)
}

func DecTab() {
	SetTab(tabLevel - 1)
}

func SetTab(i int) {
	if i < 0 {
		tabLevel = 0
	} else {
		tabLevel = i
	}
}

func Log(s ...interface{}) {
	if enabled {
		tabs := getTabs()
		str := fmt.Sprintln(tabs, s)
		logger.Output(2, str)
	}
}

func getTabs() string {
	slice := make([]rune, tabLevel)
	for i := range slice {
		slice[i] = '\t'
	}
	return string(slice)
}

func Close() {
	enabled = false
	file.Close()
}
