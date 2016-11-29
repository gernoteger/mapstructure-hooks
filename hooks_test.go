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

// plugin config über beliebiges unmarshal + github.com/mitchellh/mapstructure
// Problemstellung: Typ des Objekts ist abhängig von Context (field "kind")
package hooks

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

// Plugin is just an arbitrary marker interface
type Plugin interface {
	// execute it..just for testing
	Init() string

	// DefaultConfig returns the default config.
	// It is also uswed to determine type for unmarshalling configs
	//NewInstance() interface{}
}

//-----------------------------------------

type PlugA struct {
	A string
}

func NewPlugA() interface{} {
	return &PlugA{A: "default a"}
}

func (p *PlugA) Init() string {
	fmt.Println("a:", p.A)
	return p.A
}

//-----------------------------------------
type PlugB struct {
	B     string
	T     time.Duration
	Extra map[string]int
}

func NewPlugB() Plugin {
	return &PlugB{
		B:     "default",
		Extra: map[string]int{"dflt": 42},
	}
}

func (p *PlugB) Init() string {
	fmt.Println("b:", p.B)

	fmt.Printf("%#v", p.Extra)
	fmt.Println()
	return p.B
}

//-----------------------------------------

type Config struct {
	GlobalName        string
	Freq              time.Duration
	ReconnectInterval time.Duration
	Items             map[string]Plugin
}

type RawKind map[string]interface{}

var t1 = `
globalname: hugo
freq: 10ms
reconnectinterval: 100us

items:
  aaa:
    kind1: kindA
    a: "Aa"
  bbb:
    a: "Ab"
    kind1: kindA
  ccc:
#    a: "B"
    t: 42s
    kind1: kindB
    extra:
      X1: 1
      x2: 2
`

func TestPlugConfigWithMapstructure(t *testing.T) {
	assert := assert.New(t)

	pluginType := reflect.TypeOf((*Plugin)(nil)).Elem()

	c := Config{}

	RegisterInterface(pluginType, "kind1")

	Register(pluginType, "kindA", NewPlugA)                                 // defaults will be overruled
	Register(pluginType, "kindB", func() interface{} { return NewPlugB() }) // Example of reuse of function returning concrete type

	ci := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(t1), &ci)
	assert.Nil(err)

	//logging.Dump(ci,"raw input")

	//err = decode(ci, defaultDecoderConfig(&c))
	err = Decode(ci, &c)
	if err != nil {
		t.Error(err)
	}
	assert.Equal("hugo", c.GlobalName)

	//logging.Dump(c,"config")

	// times are set..
	assert.Equal(MustParseDuration("10ms"), c.Freq)
	assert.Equal(MustParseDuration("42s"), c.Items["ccc"].(*PlugB).T)

	// defaults too
	ccc := c.Items["ccc"].(*PlugB)
	assert.Equal("default", ccc.B)

	// and dmap
	assert.Equal(42, ccc.Extra["dflt"])
	assert.Equal(1, ccc.Extra["X1"])

	// execute them
	for _, item := range c.Items {
		item.Init()
	}

}

// MustParseDuration parses duration or panics.
// Desigend as a helper for tests
func MustParseDuration(ds string) time.Duration {
	d, err := time.ParseDuration(ds)
	if err != nil {
		panic(err)
	}
	return d
}

func TestExtractFromMapByString(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	m := map[string]interface{}{"alice": "foo", "other": "bar"}
	f := reflect.TypeOf(&m).Elem()
	k, m1, ismap, err := extractFromMap("alice", f, m)
	require.Nil(err)
	require.True(ismap)

	assert.Equal("foo", k)
	assert.EqualValues(1, len(m1))
	assert.EqualValues("bar", m1["other"])
}
func TestExtractFromMapByInterface(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	m := map[interface{}]interface{}{"alice": "foo", "other": "bar"}
	f := reflect.TypeOf(&m).Elem()
	k, m1, ismap, err := extractFromMap("alice", f, m)
	require.Nil(err)
	require.True(ismap)

	assert.Equal("foo", k)
	assert.EqualValues(1, len(m1))
	assert.EqualValues("bar", m1["other"])
}
func TestExtractFromMapNoMap(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	m := "something else"
	f := reflect.TypeOf(&m).Elem()
	k, m1, ismap, err := extractFromMap("alice", f, m)
	require.Nil(err)
	assert.False(ismap)
	assert.Empty(k)

	assert.EqualValues(0, len(m1))
}
