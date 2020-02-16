package utils

import (
    "fmt"
    "log"
    "os"
    "path"
    "time"
)

// levels
type logLevelType uint32
const (
    DebugLevel logLevelType = iota
    InfoLevel
    ErrorLevel
    FatalLevel
)

const (
    printDebugLevel = "[debug ] "
    printInfoLevel  = "[info  ] "
    printErrorLevel = "[error ] "
    printFatalLevel = "[fatal ] "
)

type Logger struct {
    level      logLevelType
    baseLogger *log.Logger
    baseFile   *os.File
}

func NewLogger(level logLevelType, pathname string) (*Logger, error) {
    // logger
    var baseLogger *log.Logger
    var baseFile *os.File
    if pathname != "" {
        now := time.Now()

        filename := fmt.Sprintf("%d%02d%02d_%02d_%02d_%02d.log",
        now.Year(),
        now.Month(),
        now.Day(),
        now.Hour(),
        now.Minute(),
        now.Second())

        file, err := os.Create(path.Join(pathname, filename))
        if err != nil {
            return nil, err
        }

        baseLogger = log.New(file, "", log.LstdFlags)
        baseFile = file
    } else {
        baseLogger = log.New(os.Stdout, "", log.LstdFlags)
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

func (logger *Logger) doPrintf(level logLevelType, printLevel string, format string, a ...interface{}) {
    if level < logger.level {
        return
    }
    if logger.baseLogger == nil {
        panic("logger closed")
    }

    format = printLevel + format
    logger.baseLogger.Output(3, fmt.Sprintf(format, a...))

    // todo 这里直接退出？不好吧，后面再去掉吧
    if level == FatalLevel {
        os.Exit(1)
    }
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

var gLogger, _ = NewLogger(DebugLevel, "")

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
