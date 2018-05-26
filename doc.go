/*
Package logger provides an additional layer of interfacing on top of the Go
standard log package, to provide several features such as indentation and
consistent default log file formatting.

Due to its nature as a wrapper of the standart Go log package, it's stuctured
very similarly, complete with a standard logger that is accessible through
the package-level functions, which writes to std.log, as well as a Logger type
that can be used to create more loggers that point to different files and can
be configured differently, so as to structure the logs better.

Package logger produces log folders in a "log" directory located in the running
user's home directory. Inside the "log" directory exist folders describing the
program names that use logger (see func Init()). In each program directory are
produced timestamped folders for each run of the specified program, containing
a file for each logger that is spawned and used during that run. This is a tree
visualization of an example log directory structure:

	homedir
	└── log
			├── prog1
			│		├── std.log
			│		└── err.log
			├── prog2
			│		└── std.log
			└── prog3
					├── std.log
					├── stats.log
					├── urgent.log
					└── err.log


Package logger keeps a set of files that it operates on so as to avoid runtime
errors and instead return an error when a duplicate Logger is attempted. TODO:
Link detection.

Disabling logging for a session is as easy as not initializing the logger with
the logger.Init() function. Any log messages can remain in the source code, and
will promptly be ignored. No runtime disabling is supported. Logging should be
enabled or disable on a session-by-session basis.

Currently, indentation control is fairly useless, due to the varied length of
the timestamp and line call at the beginning of each log message. Normalizing
the length of the message is not possible without post-op editing or trimming
of some messages.
*/
package logger
