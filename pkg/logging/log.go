package logging

import (
	"log"
	"os"
)

var f *os.File

func Info(fmt string, v ...any) {
	log.Printf("[INFO] "+fmt, v...)
}

func Panic(fmt string, v ...any) {
	log.Printf("[PANIC] "+fmt, v...)
}

func Fail(fmt string, v ...any) {
	log.Printf("[FAIL] "+fmt, v...)
}

func Warn(fmt string, v ...any) {
	log.Printf("[WARN] "+fmt, v...)
}

func SetFile(file string) {
	log.SetFlags(log.LstdFlags)
	var err error
	f, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		Fail("Failed loading file")
		return
	}
	log.SetOutput(f)
}

func Close() {
	if f != nil {
		f.Close()
	} else {
		Fail("Log file is nil")
	}
}

func CapturePanic() {
	if err := recover(); err != nil {
		Panic("%v", err)
	}
}
