// TODO:
// Faster JSON https://journal.paul.querna.org/articles/2014/03/31/ffjson-faster-json-in-go/

package logx

import (
	// "encoding/json"
	//"fmt"
	"net/url"
	"os"
	"strings"
	//"time"
)

type Options struct {
	Hostname string
	App      string
	Remote   string
	logger   Logger
}

type Record struct {
	level   string
	message string
}

var config = new(Options)

func init() {

	full_app_path := strings.Split(os.Args[0], "/")

	config.App = full_app_path[len(full_app_path)-1]
	config.Hostname, _ = os.Hostname()

}

func SetOptions(opts Options) {
	if opts.Hostname != "" {
		config.Hostname = opts.Hostname
	}

	if opts.App != "" {
		config.App = opts.App
	}

	if opts.Remote != "" {
		config.Remote = opts.Remote
		remote, err := url.Parse(opts.Remote)
		if err != nil {
			panic(err)
		}

		switch strings.ToLower(remote.Scheme) {
		case "fluent", "fluentd", "td-agent":
			config.logger = FluentLoggerInit(*config)

		default:
			panic("Unrecognized remote logging URL, try fluent://server")

		}

	}
}

func log(level string, message string, metadata interface{}) {
	if config.logger == nil {
		config.logger = FileLoggerInit(*config)
	}

	config.logger.Log(level, message, metadata)
}

func Info(message string, metadata interface{}) {
	log("info", message, metadata)
}

func Warn(message string, metadata interface{}) {
	log("warn", message, metadata)
}

func Error(message string, metadata interface{}) {
	log("error", message, metadata)
}

func Fatal(message string, metadata interface{}) {
	log("fatal", message, metadata)
}

func Debug(message string, metadata interface{}) {
	log("debug", message, metadata)
}
