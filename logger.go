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

	err = os.MkdirAll(logPath, 0777)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	enabled = true

	logger = New("std")
}

/** Logging wrappers **/

func Fatal(v ...interface{}) {
	logger.fatal(fmt.Sprint(v), 2)
}

func Fatalf(format string, v ...interface{}) {
	logger.fatal(fmt.Sprintf(format, v), 2)
}

func Fatalln(v ...interface{}) {
	logger.fatal(fmt.Sprintln(v), 2)
}

func Panic(v ...interface{}) {
	logger.panic(fmt.Sprint(v), 2)
}

func Panicf(format string, v ...interface{}) {
	logger.panic(fmt.Sprintf(format, v), 2)
}

func Panicln(v ...interface{}) {
	logger.panic(fmt.Sprintln(v), 2)
}

func Print(v ...interface{}) {
	logger.print(fmt.Sprint(v), 2)
}

func Printf(format string, v ...interface{}) {
	logger.print(fmt.Sprintf(format, v), 2)
}

func Println(v ...interface{}) {
	logger.print(fmt.Sprintln(v), 2)
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
	if !enabled {
		return nil
	}

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

func (l *Logger) Fatal(v ...interface{}) {
	l.fatal(fmt.Sprint(v), 2)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.fatal(fmt.Sprintf(format, v), 2)
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.fatal(fmt.Sprintln(v), 2)
}

func (l *Logger) Panic(v ...interface{}) {
	l.panic(fmt.Sprint(v), 2)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.panic(fmt.Sprintf(format, v), 2)
}

func (l *Logger) Panicln(v ...interface{}) {
	l.panic(fmt.Sprintln(v), 2)
}

func (l *Logger) Print(v ...interface{}) {
	l.print(fmt.Sprint(v), 2)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.print(fmt.Sprintf(format, v), 2)
}

func (l *Logger) Println(v ...interface{}) {
	l.print(fmt.Sprintln(v), 2)
}

func (l *Logger) fatal(v string, calldepth int) {
	l.print(v, calldepth+1)
	os.Exit(1)
}

func (l *Logger) panic(v string, calldepth int) {
	l.print(v, calldepth+1)
	panic(v)
}
func (l *Logger) print(v string, calldepth int) {
	str := fmt.Sprint(v)
	l.gologger.Output(calldepth-1, str)
}
