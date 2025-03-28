package aurorlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// LogLevel defines log level types
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	instance *Logger
	once     sync.Once
)

type Logger struct {
	mu         sync.Mutex
	logger     *log.Logger
	level      LogLevel
	logFile    string
	writer     io.Writer
	timeFormat string
	maxAge     time.Duration
	rotation   time.Duration
}

// GetLogger retrieves the logger instance (singleton pattern)
func GetLogger() *Logger {
	once.Do(func() {
		instance = &Logger{
			level:      DEBUG,
			timeFormat: "2006-01-02 15:04:05.000",
		}
		instance.init()
	})
	return instance
}

// init initializes the logger
func (l *Logger) init() {
	l.logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

// SetLogFile sets the log file (thread-safe)
func (l *Logger) SetLogFile(filename string, maxAge, rotation time.Duration) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Close the old writer
	if closer, ok := l.writer.(io.Closer); ok && closer != nil {
		_ = closer.Close()
	}

	l.maxAge = maxAge
	l.rotation = rotation

	writer, err := rotatelogs.New(
		filename+".%Y%m%d",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(l.maxAge),
		rotatelogs.WithRotationTime(l.rotation),
	)
	if err != nil {
		return err
	}

	l.writer = writer
	l.logFile = filename
	l.updateLogger()
	return nil
}

// updateLogger updates the log output configuration
func (l *Logger) updateLogger() {
	flags := log.Ldate | log.Ltime | log.Lshortfile
	if l.level == DEBUG {
		flags |= log.Lmicroseconds
	}
	l.logger = log.New(io.MultiWriter(os.Stdout, l.writer), "", flags)
	l.logger.SetOutput(l.writer)
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
	l.updateLogger()
}

// output handles log output
func (l *Logger) output(level LogLevel, format string, v ...interface{}) {
	if level < l.level {
		return
	}
	l.logger.Output(3, fmt.Sprintf("[%s] "+format, append([]interface{}{level.String()}, v...)...))
}

// Debug logs debug messages
func (l *Logger) Debug(format string, v ...interface{}) {
	l.output(DEBUG, format, v...)
}

// Info logs informational messages
func (l *Logger) Info(format string, v ...interface{}) {
	l.output(INFO, format, v...)
}

// Warn logs warning messages
func (l *Logger) Warn(format string, v ...interface{}) {
	l.output(WARN, format, v...)
}

// Error logs error messages
func (l *Logger) Error(format string, v ...interface{}) {
	l.output(ERROR, format, v...)
}

// Fatal logs fatal error messages and exits
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.output(FATAL, format, v...)
	os.Exit(1)
}

// String returns the string representation of the log level
func (lvl LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[lvl]
}
