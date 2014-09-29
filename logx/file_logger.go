package logx

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	// "strings"
)

func FileLoggerInit(opts Options) FileLogger {
	basepaths := []string{"./log/", "/var/log/", "./"}

	for _, path := range basepaths {
		if pathExists(path) {
			dest := fmt.Sprintf("%s%s.logx", path, opts.App)
			fmt.Printf("Logging to %s\n", dest)
			fo, err := os.OpenFile(dest, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)

			if err != nil {
				panic(err)
			}

			return FileLogger{file: fo}
			// f.file = fo

			break
		}
	}

	panic("Could not initialize a FileLogger destination log file.")
	return FileLogger{}
}

type FileLogger struct {
	file *os.File
}

func (f FileLogger) Log(level string, message string, metadata interface{}) {
	metastring := ""
	timestamp := time.Now().Format(time.RFC3339)

	if metadata != nil {
		j, _ := json.Marshal(metadata)
		metastring = string(j)
	}

	record := fmt.Sprintf("%s %s %s [%s] %s %s\n", timestamp, config.Hostname, config.App, level, message, metastring)

	_, err := f.file.WriteString(record)

	if err != nil {
		panic(err)
	}
}

func pathExists(path string) bool {
	stat, _ := os.Stat(path) // if os.IsNotExist(err) { return false, nil }
	return stat != nil
}
