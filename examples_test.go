package hooks_test

import (
	"reflect"
	"github.com/inconshreveable/log15"
	"gopkg.in/yaml.v2"
	"github.com/gernoteger/mapstructure-hooks"
	"fmt"
)

type LoggerConfig struct {
	Level    string
	Handlers []HandlerConfig
}


type HandlerConfig interface {
	NewHandler() (log15.Handler, error)
}

// use for registry functions
var HandlerConfigType = reflect.TypeOf((*HandlerConfig)(nil)).Elem()


type FileConfig struct {
	Path	string
}

func (c *FileConfig) NewHandler() (log15.Handler, error) {
	return nil,nil
}

func NewFileConfig() interface {} {
	return &FileConfig{}
}


type GelfConfig struct {
	Url	string
}

func (c *GelfConfig) NewHandler() (log15.Handler, error) {
	return nil,nil
}

func NewGelfConfig() interface {} {
	return &GelfConfig{}
}


func ExampleDecodeYaml() {
	var loggerConfig = `
   level: INFO
   handlers:
    - kind: gelf
      url: udp://myawesomehost:12201
    - kind: file
      path: /var/log/awesomeapp.log
`

	// registers all handlers
	// put into init()
	hooks.RegisterInterface(HandlerConfigType,"kind")

	hooks.Register(HandlerConfigType, "gelf", NewGelfConfig)
	hooks.Register(HandlerConfigType, "file", NewFileConfig)

	// and init your config

	ci := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(loggerConfig), &ci)
	if err != nil {
		panic(err)
	}


	c := LoggerConfig{}
	err = hooks.Decode(ci, &c)

	fmt.Println(c.Handlers[0].(*GelfConfig).Url)
	fmt.Println(c.Handlers[1].(*FileConfig).Path)
	// Output:
	//
	// udp://myawesomehost:12201
	// /var/log/awesomeapp.log

}