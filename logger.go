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
	logger *Logger

	logPath     string
	programName string

	enabled bool = false
)

func Init(progName string) {
	programName = progName

	usr, err := user.Current()
	if err != nil {
		fmt.Println("Logger couldn't get user")
		panic(err)
	}

	t := time.Now()

	currRun := fmt.Sprintf("%04d-%02d-%02d_%02d-%02d-%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	logPath = filepath.Join(usr.HomeDir, "logs", programName, currRun)

	err = os.Mkdir(logPath, 0777)
	if !os.IsExist(err) {
		panic(err)
	}

	logger = New("std")
	enabled = true
}

func timestampedFile(logName string) (*os.File, error) {
	file, err := os.Create(filepath.Join(logPath, logName+".log"))

	if err != nil {
		fmt.Println("Logger couldn't create file")
		return nil, err
	}

	return file, nil
}

type Logger struct {
	gologger *log.Logger
	tabLevel int
}

func New(filename string) *Logger {
	file, err := timestampedFile(filename)
	if err != nil {
		panic(err)
	}

	gologger := log.New(file, "", log.LstdFlags|log.Lshortfile)

	return &Logger{gologger, 0}
}

func (l *Logger) IncTab() {
	l.SetTab(l.tabLevel + 1)
}

func (l *Logger) DecTab() {
	l.SetTab(l.tabLevel - 1)
}

func (l *Logger) SetTab(i int) {
	if i < 0 {
		l.tabLevel = 0
	} else {
		l.tabLevel = i
	}
}

func (l *Logger) TabLevel() int {
	return l.tabLevel
}

func (l *Logger) Log(s ...interface{}) {
	tabs := l.tabs()
	str := fmt.Sprintln(tabs, s)
	l.gologger.Output(2, str)
}

func (l *Logger) tabs() string {
	slice := make([]rune, l.tabLevel)
	for i := range slice {
		slice[i] = '\t'
	}
	return string(slice)
}
