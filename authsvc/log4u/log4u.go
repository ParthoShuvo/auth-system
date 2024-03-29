// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package log implements a simple logging package. It defines a type, Logger,
// with methods for formatting output. It also has a predefined 'standard'
// Logger accessible through helper functions Print[f|ln], Fatal[f|ln], and
// Panic[f|ln], which are easier to use than creating a Logger manually.
// That logger writes to standard error and prints the date and time
// of each logged message.
// Every log message is output on a separate line: if the message being
// printed does not end in a newline, the logger will add one.
// The Fatal functions call os.Exit(1) after writing the log message.
// The Panic functions call panic after writing the log message.

// Copied from the original Go log package and added log level functionality
// much like log4j has. This to allow integration with the Kibana log tool.

package log4u

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// These flags define which text to prefix to each log entry generated by the Logger.
const (
	// Bits or'ed together to control what's printed.
	// There is no control over the order they appear (the order listed
	// here) or the format they present (as described in the comments).
	// The prefix is followed by a colon only when Llongfile or Lshortfile
	// is specified.
	// For example, flags Ldate | Ltime (or LstdFlags) produce,
	//	2009/01/23 01:23:23 message
	// while flags Ldate | Ltime | Lmilliseconds | Llongfile produce,
	//	2009/01/23 01:23:23.123 /a/b/c/d.go:23: message
	Ldate         = 1 << iota                                  // the date in the local time zone: 2009/01/23
	Ltime                                                      // the time in the local time zone: 01:23:23
	Lmilliseconds                                              // millisecond resolution: 01:23:23,123.  assumes Ltime.
	Llongfile                                                  // full file name and line number: /a/b/c/d.go:23
	Lshortfile                                                 // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                                                       // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime | Lmilliseconds | Lshortfile // initial values for the standard logger
)

// LogLevel defines the logging level.
type LogLevel int

// These constants define the available logging levels.
const (
	Ldebug LogLevel = iota
	Linfo
	Lwarn
	Lerror
	Lfatal
)

// A Logger represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	mu     sync.Mutex // ensures atomic writes; protects the following fields
	prefix string     // prefix to write at beginning of each line
	flag   int        // properties
	level  LogLevel   // the logging level
	out    io.Writer  // destination for output
	buf    []byte     // for accumulating text to write
}

var levelTags []string

func init() {
	levelTags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
}

// New creates a new Logger. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// The flag argument defines the logging properties.
func New(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{out: out, prefix: prefix, flag: flag, level: Linfo}
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

var std = New(os.Stderr, "", LstdFlags)

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// formatHeader writes log header to buf in following order:
//	 * the logging level
//   * l.prefix (if it's not blank),
//   * date and/or time (if corresponding flags are provided),
//   * file and line number (if corresponding flags are provided).
func (l *Logger) formatHeader(level LogLevel, buf *[]byte, t time.Time, method string, file string, line int) {
	*buf = append(*buf, '[')
	*buf = append(*buf, levelTags[level]...)
	*buf = append(*buf, ']')
	*buf = append(*buf, ' ')
	*buf = append(*buf, l.prefix...)
	if l.flag&(Ldate|Ltime|Lmilliseconds) != 0 {
		if l.flag&LUTC != 0 {
			t = t.UTC()
		}
		*buf = append(*buf, '[')
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '-')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '-')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmilliseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmilliseconds != 0 {
				*buf = append(*buf, ',')
				itoa(buf, t.Nanosecond()/1e6, 3)
			}
		}
		*buf = append(*buf, ']')
		*buf = append(*buf, ' ')
	}
	*buf = append(*buf, "[???] "...)
	if l.flag&(Lshortfile|Llongfile) != 0 {
		*buf = append(*buf, '[')
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		*buf = append(*buf, method...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ']')
		*buf = append(*buf, ' ')
	}
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *Logger) Output(calldepth int, level LogLevel, s string) error {
	now := time.Now() // get this early.
	var pc uintptr
	var method string
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(Lshortfile|Llongfile) != 0 {
		// Release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		pc, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			method = "???"
			file = "???"
			line = 0
		} else {
			f := runtime.FuncForPC(pc)
			if f != nil {
				method = f.Name()
			} else {
				method = "???"
			}
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(level, &l.buf, now, method, file, line)
	l.buf = append(l.buf, "==> "...)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return err
}

// Debugf calls l.Output to print to the logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.mustLog(Ldebug) {
		l.Output(2, Ldebug, fmt.Sprintf(format, v...))
	}
}

// Infof calls l.Output to print to the logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.mustLog(Linfo) {
		l.Output(2, Linfo, fmt.Sprintf(format, v...))
	}
}

// Warnf calls l.Output to print to the logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.mustLog(Lwarn) {
		l.Output(2, Lwarn, fmt.Sprintf(format, v...))
	}
}

// Errorf calls l.Output to print to the logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.mustLog(Lerror) {
		l.Output(2, Lerror, fmt.Sprintf(format, v...))
	}
}

// Debug calls l.Output to print to the logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Debug(v ...interface{}) {
	if l.mustLog(Ldebug) {
		l.Output(2, Ldebug, fmt.Sprint(v...))
	}
}

// Info calls l.Output to print to the logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Info(v ...interface{}) {
	if l.mustLog(Linfo) {
		l.Output(2, Linfo, fmt.Sprint(v...))
	}
}

// Warn calls l.Output to print to the logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Warn(v ...interface{}) {
	if l.mustLog(Lwarn) {
		l.Output(2, Lwarn, fmt.Sprint(v...))
	}
}

// Error calls l.Output to print to the logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Error(v ...interface{}) {
	if l.mustLog(Lerror) {
		l.Output(2, Lerror, fmt.Sprint(v...))
	}
}

// Debugln calls l.Output to print to the logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Debugln(v ...interface{}) {
	if l.mustLog(Ldebug) {
		l.Output(2, Ldebug, fmt.Sprintln(v...))
	}
}

// Infoln calls l.Output to print to the logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Infoln(v ...interface{}) {
	if l.mustLog(Linfo) {
		l.Output(2, Linfo, fmt.Sprintln(v...))
	}
}

// Warnln calls l.Output to print to the logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Warnln(v ...interface{}) {
	if l.mustLog(Lwarn) {
		l.Output(2, Lwarn, fmt.Sprintln(v...))
	}
}

// Errorln calls l.Output to print to the logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Errorln(v ...interface{}) {
	if l.mustLog(Lerror) {
		l.Output(2, Lerror, fmt.Sprintln(v...))
	}
}

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Println(v ...interface{}) {
	l.Output(2, Linfo, fmt.Sprintln(v...))
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *Logger) Fatal(v ...interface{}) {
	l.Output(2, Lfatal, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(2, Lfatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *Logger) Fatalln(v ...interface{}) {
	l.Output(2, Lfatal, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Output(2, Lfatal, s)
	panic(s)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Output(2, Lfatal, s)
	panic(s)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *Logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.Output(2, Lfatal, s)
	panic(s)
}

func (l *Logger) mustLog(level LogLevel) bool {
	return level >= l.getLevel()
}

// Flags returns the output flags for the logger.
func (l *Logger) Flags() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.flag
}

// SetFlags sets the output flags for the logger.
func (l *Logger) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flag = flag
}

// Prefix returns the output prefix for the logger.
func (l *Logger) Prefix() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.prefix
}

// SetPrefix sets the output prefix for the logger.
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

// Level return the current logging level
func (l *Logger) Level() string {
	return levelTags[l.getLevel()]
}

// level return the current logging level
func (l *Logger) getLevel() LogLevel {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

// SetLevel sets the logging level.
func (l *Logger) SetLevel(level string) {
	for i, tag := range levelTags {
		if strings.EqualFold(tag, level) {
			l.setLevel(LogLevel(i))
			return
		}
	}
}

// SetLevel sets the logging level
func (l *Logger) setLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetOutput sets the output destination for the standard logger.
func SetOutput(w io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.out = w
}

// Flags returns the output flags for the standard logger.
func Flags() int {
	return std.Flags()
}

// SetFlags sets the output flags for the standard logger.
func SetFlags(flag int) {
	std.SetFlags(flag)
}

// Level returns the current logging level for the standard logger.
func Level() string {
	return std.Level()
}

// SetLevel sets the logging level for the standard logger.
func SetLevel(level string) {
	std.SetLevel(level)
}

// Prefix returns the output prefix for the standard logger.
func Prefix() string {
	return std.Prefix()
}

// SetPrefix sets the output prefix for the standard logger.
func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

// These functions write to the standard logger.

// Debug calls Output to print to the standard logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Print.
func Debug(v ...interface{}) {
	if std.mustLog(Ldebug) {
		std.Output(2, Ldebug, fmt.Sprint(v...))
	}
}

// Info calls Output to print to the standard logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Print.
func Info(v ...interface{}) {
	if std.mustLog(Linfo) {
		std.Output(2, Linfo, fmt.Sprint(v...))
	}
}

// Warn calls Output to print to the standard logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Print.
func Warn(v ...interface{}) {
	if std.mustLog(Lwarn) {
		std.Output(2, Lwarn, fmt.Sprint(v...))
	}
}

// Error calls Output to print to the standard logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Print.
func Error(v ...interface{}) {
	if std.mustLog(Lerror) {
		std.Output(2, Lerror, fmt.Sprint(v...))
	}
}

// Debugf calls Output to print to the standard logger.
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, v ...interface{}) {
	if std.mustLog(Ldebug) {
		std.Output(2, Ldebug, fmt.Sprintf(format, v...))
	}
}

// Infof calls Output to print to the standard logger.
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, v ...interface{}) {
	if std.mustLog(Linfo) {
		std.Output(2, Linfo, fmt.Sprintf(format, v...))
	}
}

// Warnf calls Output to print to the standard logger.
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Printf.
func Warnf(format string, v ...interface{}) {
	if std.mustLog(Lwarn) {
		std.Output(2, Lwarn, fmt.Sprintf(format, v...))
	}
}

// Errorf calls Output to print to the standard logger.
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	if std.mustLog(Lerror) {
		std.Output(2, Lerror, fmt.Sprintf(format, v...))
	}
}

// Debugln calls Output to print to the standard logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Println.
func Debugln(v ...interface{}) {
	if std.mustLog(Ldebug) {
		std.Output(2, Ldebug, fmt.Sprintln(v...))
	}
}

// Infoln calls Output to print to the standard logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Println.
func Infoln(v ...interface{}) {
	if std.mustLog(Linfo) {
		std.Output(2, Linfo, fmt.Sprintln(v...))
	}
}

// Warnln calls Output to print to the standard logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Println.
func Warnln(v ...interface{}) {
	if std.mustLog(Lwarn) {
		std.Output(2, Lwarn, fmt.Sprintln(v...))
	}
}

// Errorln calls Output to print to the standard logger
// if the current logging level allows that.
// Arguments are handled in the manner of fmt.Println.
func Errorln(v ...interface{}) {
	if std.mustLog(Lerror) {
		std.Output(2, Lerror, fmt.Sprintln(v...))
	}
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	std.Output(2, Lfatal, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	std.Output(2, Lfatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	std.Output(2, Lfatal, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	std.Output(2, Lfatal, s)
	panic(s)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	std.Output(2, Lfatal, s)
	panic(s)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	std.Output(2, Lfatal, s)
	panic(s)
}

// Output writes the output for a logging event. The string s contains
// the text to print after the level tag and the prefix specified by
// the flags of the Logger. A newline is appended if the last character of
// s is not already a newline. Calldepth is the count of the number of
// frames to skip when computing the file name and line number
// if Llongfile or Lshortfile is set; a value of 1 will print the details
// for the caller of Output.
func Output(calldepth int, level LogLevel, s string) error {
	return std.Output(calldepth+1, level, s) // +1 for this frame.
}
