package utils

import (
    "fmt"
    "log"
    "os"
    "path"
    "time"
)

// levels
type LogLevelType uint32
const (
    DebugLevel LogLevelType = iota
    InfoLevel
    ErrorLevel
    FatalLevel
)

const (
    printDebugLevel = "|debug|"
    printInfoLevel  = "|info|"
    printErrorLevel = "|error|"
    printFatalLevel = "|fatal|"
)

type Logger struct {
    level      LogLevelType
    baseLogger *log.Logger
    baseFile   *os.File
}

func NewLogger(level LogLevelType, pathname string, prefix string) (*Logger, error) {
    var baseLogger *log.Logger
    var baseFile *os.File
    if pathname != "" {
        now := time.Now()

        err := os.MkdirAll(pathname, os.ModeDir | os.ModePerm)
        if err != nil {
            return nil, err
        }

        filename := fmt.Sprintf("%s_%d%02d%02d%02d.log",
        prefix,
        now.Year(),
        now.Month(),
        now.Day(),
        now.Hour())

        file, err := os.OpenFile(path.Join(pathname, filename), os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
        if err != nil {
            return nil, err
        }

        pid_str := fmt.Sprintf("%d|", os.Getpid())

        baseLogger = log.New(file, pid_str, log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
        baseFile = file
    } else {
        baseLogger = log.New(os.Stdout, "", log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
    }

    logger := new(Logger)
    logger.level = level
    logger.baseLogger = baseLogger
    logger.baseFile = baseFile

    return logger, nil
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
    if logger.baseFile != nil {
        logger.baseFile.Close()
    }

    logger.baseLogger = nil
    logger.baseFile = nil
}

func (logger *Logger) doPrintf(level LogLevelType, printLevel string, format string, a ...interface{}) {
    if level < logger.level {
        return
    }
    if logger.baseLogger == nil {
        panic("logger closed")
    }

    format = printLevel + format
    // 这个call depth参数有点蛋疼，竟然要使用的人去算究竟调用栈有多少层
    // 因为需要打印的位置是调用下面的Debug之类的函数的位置，所以最终到里面的runtime.Caller调用栈是三层
    logger.baseLogger.Output(3, fmt.Sprintf(format, a...))
}

func (logger *Logger) Debug(format string, a ...interface{}) {
    logger.doPrintf(DebugLevel, printDebugLevel, format, a...)
}

func (logger *Logger) Info(format string, a ...interface{}) {
    logger.doPrintf(InfoLevel, printInfoLevel, format, a...)
}

func (logger *Logger) Error(format string, a ...interface{}) {
    logger.doPrintf(ErrorLevel, printErrorLevel, format, a...)
}

func (logger *Logger) Fatal(format string, a ...interface{}) {
    logger.doPrintf(FatalLevel, printFatalLevel, format, a...)
}

var gLogger, _ = NewLogger(DebugLevel, "", "__global__")

// It's dangerous to call the method on logging
func ExportLogger(logger *Logger) {
    if logger != nil {
        gLogger = logger
    }
}

func Debug(format string, a ...interface{}) {
    gLogger.doPrintf(DebugLevel, printDebugLevel, format, a...)
}

func Info(format string, a ...interface{}) {
    gLogger.doPrintf(InfoLevel, printInfoLevel, format, a...)
}

func Error(format string, a ...interface{}) {
    gLogger.doPrintf(ErrorLevel, printErrorLevel, format, a...)
}

func Fatal(format string, a ...interface{}) {
    gLogger.doPrintf(FatalLevel, printFatalLevel, format, a...)
}

func Close() {
    gLogger.Close()
}
