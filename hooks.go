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

// Package hooks helps to configure hook configurations from a map.
// It's a wrapper around mapstructure.
//
// The main functions to be used are:
//   RegisterInterface: just use once to register an interface type
//   Register:          used to register possible kinds of an interface
//   Decode:            used once to decode the config
// All other functions are for advanced uses.
package hooks

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// this file holds everything to initialize & config hook structures, e.g log writers, meters.

// NewHookFunc creates a new instance for an registered interface.
type NewHookFunc func() interface{}

type interfaceMapper struct {
	kindKey  string // map
	newfuncs map[string]NewHookFunc
}

// InitRegistry resets the registry. It is intended mainly for testing.
func InitRegistry() {
	registry = make(map[reflect.Type]*interfaceMapper)
}

// the registry maps all types to their rules
var registry = make(map[reflect.Type]*interfaceMapper)

// RegisterInterface registers a new interface type; must only be called once!
func RegisterInterface(forType reflect.Type, key string) {
	_, found := registry[forType]
	if found {
		panic(fmt.Sprintf("interface already registered: %#v", forType))
	}

	registry[forType] = &interfaceMapper{
		kindKey:  key,
		newfuncs: make(map[string]NewHookFunc, 2),
	}
}

// Register registers the factory function for all instances of a given type
// hint: provide a type for consumers with MyHookType := reflect.TypeOf((*MyHook)(nil)).Elem()
func Register(forType reflect.Type, kind string, f func() interface{}) {
	//fmt.Printf("hooks.Register kind='%v' for '%#v'\n", kind, forType)

	im := registry[forType]
	im.newfuncs[kind] = f
}

// DefaultDecoderConfig returns default mapstructure.DecoderConfig with support
// of time.Duration values and Plugins
func DefaultDecoderConfig(output interface{}) *mapstructure.DecoderConfig {
	return &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: false,

		ErrorUnused: true,

		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			DecodeElementsHookFunc(),
			StringToStringUnmarshallerHookFunc(),
		),
	}
}

// A wrapper around mapstructure.Decode that mimics the WeakDecode functionality
func decode(input interface{}, config *mapstructure.DecoderConfig) error {
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

// Decode transkripts map[interface{}]interface{} to the target struct
func Decode(in, out interface{}) error {
	return decode(in, DefaultDecoderConfig(out))
}

var mapTypeIf = reflect.TypeOf((*map[interface{}]interface{})(nil)).Elem()
var mapTypeString = reflect.TypeOf((*map[string]interface{})(nil)).Elem()

// DecodeElementsHookFunc returns a DecodeHookFunc that converts
// maps to a Hook config derived from a registry.
func DecodeElementsHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {

		// not obvious, but, hey, it works...

		// check target
		im, registered := registry[t]
		if !registered {
			return data, nil
		}

		kind, m1, ismap, err := extractFromMap(im.kindKey, f, data)
		if !ismap {
			return data, nil
		}
		if err != nil {
			return data, fmt.Errorf("no kind with key '%v' found: %v", im.kindKey, err)
		}

		if kind == "" {
			return data, errors.New("no kind found")
		}

		newInstanceFunc := im.newfuncs[kind]
		if newInstanceFunc == nil {
			return data, fmt.Errorf("no registered kind '%v' for type %v", kind, t)
		}

		hook := newInstanceFunc() // get new instance

		// Convert it by parsing
		err = Decode(m1, hook)
		if err != nil {
			return data, err
		}

		return hook, nil
	}
}

// RemoveFromMap removes a key from the map; the map can have strings or interfaces as key
// if it's not  amp, an error is returned
func extractFromMap(key string, f reflect.Type, data interface{}) (rkey string, newmap map[string]interface{}, ismap bool, err error) {

	m1 := make(map[string]interface{})

	var ok bool
	switch f {
	case mapTypeIf:
		mm := data.(map[interface{}]interface{})
		rkey, ok = mm[key].(string)
		// a case for generics!!!
		for k, v := range mm {
			if k != key {
				m1[k.(string)] = v
			}
		}
	case mapTypeString:
		mm := data.(map[string]interface{})
		rkey, ok = mm[key].(string)
		// a case for generics!!!
		for k, v := range mm {
			if k != key {
				m1[k] = v
			}
		}
	default: //  not a map: tell outside to ignore
		return "", m1, false, nil
	}
	if !ok {
		return "", m1, true, fmt.Errorf("no element '%v' found", key)
	}

	return rkey, m1, true, nil
}

type stringUnmarshaller interface {
	// Unmarshal into struct
	UnmarshalString(from string) (interface{}, error)
}

var stringUnmarshallerType = reflect.TypeOf((*stringUnmarshaller)(nil)).Elem()

// StringToStringUnmarshallerHookFunc returns a DecodeHookFunc that converts
// strings by an unmarshaller. Can be used to construct custom DecodeHookFunctions.
func StringToStringUnmarshallerHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {

		if f.Kind() != reflect.String {
			return data, nil
		}

		if t.Implements(stringUnmarshallerType) {
			// Convert it by unmarshaller
			//fmt.Println("=========== by struct 2")
			e := reflect.New(t).Interface()
			e1, err := e.(stringUnmarshaller).UnmarshalString(data.(string))
			return e1, err
		}

		return data, nil
	}
}
