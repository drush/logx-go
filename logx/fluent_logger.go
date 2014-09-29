// TODO
// Ensure fluent://server:port style Remotes can be parsed
package logx

import (
	"encoding/json"
	"fmt"
	"github.com/t-k/fluent-logger-golang/fluent"
	"net/url"
	"reflect"
	"strings"
)

type FluentLogger struct {
	fluent *fluent.Fluent
	config *Options
}

func FluentLoggerInit(opts Options) FluentLogger {
	remote, _ := url.Parse(opts.Remote)

	f, err := fluent.New(fluent.Config{FluentHost: remote.Host, FluentPort: 24224})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Logging to %s\n", opts.Remote)

	//defer logger.Close()
	//logger.Post(tag, data)

	return FluentLogger{fluent: f, config: &opts}
}

func (f FluentLogger) Log(level string, message string, metadata interface{}) {

	// timestamp := time.Now().Format(time.RFC3339)

	var data = map[string]string{
		"host":    f.config.Hostname,
		"app":     f.config.App,
		"level":   level,
		"message": message}

	if metadata != nil {
		if metamap, ok := metadata.(map[string]string); ok {
			for key, _ := range metamap {
				data[key] = metamap[key]
			}
		}

		// if reflect.TypeOf(metadata).Kind() == reflect.Struct {
		// 	j, _ := json.Marshal(metadata)
		// 	data["meta"] = string(j)
		// }

		a := attributes(metadata)
		for key, _ := range a {
			data[key] = a[key]
		}
	}

	f.fluent.Post("tag", data)

	// if err != nil {
	// 	panic(err)
	// }
}

func attributes(m interface{}) map[string]string {
	typ := reflect.TypeOf(m)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// create an attribute data structure as a map of types keyed by a string.
	attrs := make(map[string]string)
	// Only structs are supported so return an empty result if the passed object
	// isn't a struct
	if typ.Kind() != reflect.Struct {
		fmt.Printf("%v type can't have attributes inspected\n", typ.Kind())
		return attrs
	}

	// loop through the struct's fields and set the map
	for i := 0; i < typ.NumField(); i++ {
		p := typ.Field(i)

		if !p.Anonymous && p.Type.Kind() != reflect.Ptr {
			name := p.Name
			tag := p.Tag.Get("json")
			if tag != "" && tag != "-" {
				name = parseTag(tag)
			}
			val, _ := json.Marshal(reflect.ValueOf(m).Field(i).Interface())
			attrs[name] = string(val)
		}
	}

	return attrs
}

func parseTag(tag string) string {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}
