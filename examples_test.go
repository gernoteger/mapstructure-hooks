//Copyright 2016 Gernot Eger
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package hooks_test

import (
	"fmt"
	"reflect"

	"github.com/gernoteger/mapstructure-hooks"
	"github.com/inconshreveable/log15"
	"gopkg.in/yaml.v2"
)

type LoggerConfig struct {
	Level    string
	Handlers []HandlerConfig
}

// HandlerConfig is the common interface
type HandlerConfig interface {
	NewHandler() (log15.Handler, error)
}

// use for registry functions
var HandlerConfigType = reflect.TypeOf((*HandlerConfig)(nil)).Elem()

type FileConfig struct {
	Path string
}

func (c *FileConfig) NewHandler() (log15.Handler, error) {
	return nil, nil
}

func NewFileConfig() interface{} {
	return &FileConfig{}
}

type GelfConfig struct {
	URL string
}

func (c *GelfConfig) NewHandler() (log15.Handler, error) {
	return nil, nil
}

func NewGelfConfig() interface{} {
	return &GelfConfig{}
}

func ExampleDecode() {
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
	hooks.RegisterInterface(HandlerConfigType, "kind")

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
	if err != nil {
		panic(err)
	}

	fmt.Println(c.Handlers[0].(*GelfConfig).URL)
	fmt.Println(c.Handlers[1].(*FileConfig).Path)
	// Output:
	//
	// udp://myawesomehost:12201
	// /var/log/awesomeapp.log

}
