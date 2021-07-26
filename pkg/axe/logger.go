package axe

import (
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Print(args ...interface{})

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Printf(template string, args ...interface{})
}

type DefaultLogger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	panicLogger *log.Logger
	fatalLogger *log.Logger
	printLogger *log.Logger
}

func New() *DefaultLogger {
	flags := log.Ldate | log.Ltime | log.Lshortfile
	return NewWithLoggers(
		log.New(os.Stdout, "DEBUG\t", flags),
		log.New(os.Stdout, "INFO\t", flags),
		log.New(os.Stdout, "WARN\t", flags),
		log.New(os.Stdout, "ERROR\t", flags),
		log.New(os.Stdout, "PANIC\t", flags),
		log.New(os.Stdout, "FATAL\t", flags),
		log.New(os.Stdout, "", flags),
	)
}

func NewWithLoggers(debugLogger *log.Logger, infoLogger *log.Logger, warnLogger *log.Logger, errorLogger *log.Logger, panicLogger *log.Logger, fatalLogger *log.Logger, printLogger *log.Logger) *DefaultLogger {
	if panicLogger == nil {
		panic(fmt.Errorf("cannot create a Logger with a nil panicLogger"))
	}
	if debugLogger == nil {
		panicLogger.Panic(fmt.Errorf("cannot create a Logger with a nil debugLogger"))
	}
	if infoLogger == nil {
		panicLogger.Panic(fmt.Errorf("cannot create a Logger with a nil infoLogger"))
	}
	if warnLogger == nil {
		panicLogger.Panic(fmt.Errorf("cannot create a Logger with a nil warnLogger"))
	}
	if errorLogger == nil {
		panicLogger.Panic(fmt.Errorf("cannot create a Logger with a nil errorLogger"))
	}
	if fatalLogger == nil {
		panicLogger.Panic(fmt.Errorf("cannot create a Logger with a nil fatalLogger"))
	}
	if printLogger == nil {
		panicLogger.Panic(fmt.Errorf("cannot create a Logger with a nil fatalLogger"))
	}
	return &DefaultLogger{
		debugLogger: debugLogger,
		infoLogger:  infoLogger,
		warnLogger:  warnLogger,
		errorLogger: errorLogger,
		panicLogger: panicLogger,
		fatalLogger: fatalLogger,
		printLogger: printLogger,
	}
}

func (l *DefaultLogger) Debug(args ...interface{}) {
	l.debugLogger.Output(2, fmt.Sprintln(args...))
}
func (l *DefaultLogger) Info(args ...interface{}) {
	l.infoLogger.Output(2, fmt.Sprintln(args...))
}
func (l *DefaultLogger) Warn(args ...interface{}) {
	l.warnLogger.Output(2, fmt.Sprintln(args...))
}
func (l *DefaultLogger) Error(args ...interface{}) {
	l.errorLogger.Output(2, fmt.Sprintln(args...))
}
func (l *DefaultLogger) Panic(args ...interface{}) {
	l.panicLogger.Output(2, fmt.Sprintln(args...))
}
func (l *DefaultLogger) Fatal(args ...interface{}) {
	l.fatalLogger.Output(2, fmt.Sprintln(args...))
}
func (l *DefaultLogger) Print(args ...interface{}) {
	l.printLogger.Output(2, fmt.Sprintln(args...))
}

func (l *DefaultLogger) Debugf(template string, args ...interface{}) {
	l.debugLogger.Output(2, fmt.Sprintf(template+"\n", args...))
}
func (l *DefaultLogger) Infof(template string, args ...interface{}) {
	l.infoLogger.Output(2, fmt.Sprintf(template+"\n", args...))
}
func (l *DefaultLogger) Warnf(template string, args ...interface{}) {
	l.warnLogger.Output(2, fmt.Sprintf(template+"\n", args...))
}
func (l *DefaultLogger) Errorf(template string, args ...interface{}) {
	l.errorLogger.Output(2, fmt.Sprintf(template+"\n", args...))
}
func (l *DefaultLogger) Panicf(template string, args ...interface{}) {
	l.panicLogger.Output(2, fmt.Sprintf(template+"\n", args...))
}
func (l *DefaultLogger) Fatalf(template string, args ...interface{}) {
	l.fatalLogger.Output(2, fmt.Sprintf(template+"\n", args...))
}
func (l *DefaultLogger) Printf(template string, args ...interface{}) {
	l.printLogger.Output(2, fmt.Sprintf(template+"\n", args...))
}
