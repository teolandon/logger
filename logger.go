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

	enabled = false
)

// Init has to be called before any logger can be initialized,
// including the default logger. The parameter progName specifies
// the name of the program to be ran, so as to place the log
// files in the correct folder.
//
// Init initializes the default logger to the current timestamp
// and given program name, pointing to the file std.log, the standard
// logging file.
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
	logger.fatal(fmt.Sprintf(format, v...), 2)
}

func Fatalln(v ...interface{}) {
	logger.fatal(fmt.Sprintln(v), 2)
}

func Panic(v ...interface{}) {
	logger.panic(fmt.Sprint(v), 2)
}

func Panicf(format string, v ...interface{}) {
	logger.panic(fmt.Sprintf(format, v...), 2)
}

func Panicln(v ...interface{}) {
	logger.panic(fmt.Sprintln(v), 2)
}

func Print(v ...interface{}) {
	logger.print(fmt.Sprint(v), 2)
}

func Printf(format string, v ...interface{}) {
	fmt.Println("Format string: ", format)
	fmt.Println("Arguments:  ", v)
	logger.print(fmt.Sprintf(format, v...), 2)
}

func Println(v ...interface{}) {
	logger.print(fmt.Sprintln(v), 2)
}

// newLogFile creates and returns a new log file with the given
// logName plus the ".log" extension, inside the current log path
// given by the timestamp of the initialization of the package and
// the name of the program that is being logged.
//
// The error returned can be any of the errors that os.Create()
// returns, returned when the file creation fails.
func newLogFile(logName string) (*os.File, error) {
	file, err := os.Create(filepath.Join(logPath, logName+".log"))

	if err != nil {
		fmt.Println("Logger couldn't create file")
		return nil, err
	}

	return file, nil
}

// A Logger can be used to log messages to a file using the standard Go Logger
// methods. Multiple loggers can be present during a program's run. In fact,
// the intended usage is to group log messages with similar purposes in
// different loggers, so as to avoid clutter and better organize logs.
//
// Furthermore, a tab level can be specified to indent lines. Common
// usage is increasing the tab level before calling an important function, and
// decreasing it back to the previous level after it returns:
//
//   // previous code
//
//   logger.Println("Calling functionThatLogsActions")
//   logger.IncTab()
//   functionThatLogsActions()
//   logger.DecTab()
//
//   // rest of the code
//
// It's also possible to increase the tab level at the beginning of all
// functions and defer decreasing it, but this can cause excessive
// indentation, and is not recommended.
//
// The tab characted can be set to any character to provide better visibility
// of indented log entries.
type Logger struct {
	gologger *log.Logger
	tabLevel int
}

// New initializes and returns a new Logger pointing to a file located in
// the current timestamped directory, with the given filename and the
// ".log" extension.
func New(filename string) *Logger {
	if !enabled {
		return nil
	}

	file, err := newLogFile(filename)
	if err != nil {
		panic(err)
	}

	gologger := log.New(file, "", log.LstdFlags|log.Lshortfile)

	return &Logger{gologger, 0}
}

// IncTab increases the indent level of the Logger l by 1 tab character.
func (l *Logger) IncTab() {
	l.SetTab(l.tabLevel + 1)
}

// DecTab decreases the indent level of the Logger l by 1 tab character.
func (l *Logger) DecTab() {
	l.SetTab(l.tabLevel - 1)
}

// SetTab sets the indent level of the Logger l to i tab characters. The given
// number i has to be non-negative.
func (l *Logger) SetTab(i int) {
	if i < 0 {
		l.tabLevel = 0
	} else {
		l.tabLevel = i
	}
}

// TabLevel returns the current indentation level of the Logger l.
func (l *Logger) TabLevel() int {
	return l.tabLevel
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
	l.panic(fmt.Sprintf(format, v...), 2)
}

func (l *Logger) Panicln(v ...interface{}) {
	l.panic(fmt.Sprintln(v), 2)
}

func (l *Logger) Print(v ...interface{}) {
	l.print(fmt.Sprint(v), 2)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.print(fmt.Sprintf(format, v...), 2)
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
	l.gologger.Output(calldepth+1, str)
}
