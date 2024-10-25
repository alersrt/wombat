package log

import (
	"log"
	"os"
)

var (
	infoLog    = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLog = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog   = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func Info(v ...any) {
	infoLog.Print(v)
}

func Infoln(v ...any) {
	infoLog.Println(v)
}

func Infof(format string, v ...any) {
	infoLog.Printf(format, v...)
}

func Warn(v ...any) {
	warningLog.Print(v)
}
func Warnln(v ...any) {
	warningLog.Println(v)
}

func Warnf(format string, v ...any) {
	warningLog.Printf(format, v...)
}

func Error(v ...any) {
	errorLog.Print(v)
}

func Errorln(v ...any) {
	errorLog.Println(v)
}

func Errorf(format string, v ...any) {
	errorLog.Printf(format, v)
}
