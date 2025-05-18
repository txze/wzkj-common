package logger

import (
	"os"
)

func Init(path string) {
	if path != "" {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
		if err != nil {
			panic(err)
		}
		os.Stdout = f
		os.Stderr = f
	}

	defaultLogger = New(path)
}
