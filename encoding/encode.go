package encoding

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

func stringInSlice(search string, slice []string) bool {
	for _, i := range slice {
		if i == search {
			return true
		}
	}
	return false
}

type tagOptions []string

// Has returns true if opt is one of the options
func (t tagOptions) Has(opt string) bool {
	for _, opt := range t {
		if opt == opt {
			return true
		}
	}
	return false
}

func parseTag(tag string) (string, tagOptions) {
	res := strings.Split(tag, ",")
	return res[0], res[1:]
}

type isZeroer interface {
	IsZero() bool
}

type isEmptier interface {
	IsEmpty() bool
}

// ToStringStringMap encodes fields tagged with tagName in input into map[string]string. Optional argument properties specifies fields to encode.
func ToStringStringMap(tagName string, input interface{}, properties ...string) (map[string]string, error) {
	vmap := make(map[string]string)
	s := structs.New(input)
	s.TagName = tagName
	if len(properties) == 0 {
		properties = s.Names()
	}

	for _, field := range s.Fields() {
		if !field.IsExported() {
			continue
		}

		if !stringInSlice(field.Name(), properties) {
			continue
		}

		tagName, opts := parseTag(field.Tag(tagName))
		if tagName == "" || tagName == "-" {
			continue
		}

		val := field.Value()

		if opts.Has("omitempty") {
			if field.IsZero() {
				continue
			}
			if z, ok := val.(isZeroer); ok && z.IsZero() {
				continue
			}
			if z, ok := val.(isEmptier); ok && z.IsEmpty() {
				continue
			}
		}

		if v, ok := val.(string); ok {
			vmap[tagName] = v
			continue
		} else if v, ok := val.(*string); ok {
			vmap[tagName] = *v
			continue
		}

		if !field.IsZero() {
			if m, ok := val.(encoding.TextMarshaler); ok {
				txt, err := m.MarshalText()
				if err != nil {
					return nil, err
				}
				vmap[tagName] = string(txt)
				continue
			}
			if m, ok := val.(json.Marshaler); ok {
				txt, err := m.MarshalJSON()
				if err != nil {
					return nil, err
				}
				vmap[tagName] = string(txt)
				continue
			}
		}

		if field.Kind() == reflect.String {
			vmap[tagName] = fmt.Sprint(val)
			continue
		}

		if txt, err := json.Marshal(val); err == nil {
			vmap[tagName] = string(txt)
			if vmap[tagName] == `""` || vmap[tagName] == "null" {
				vmap[tagName] = ""
			}
			continue
		}

		vmap[tagName] = fmt.Sprintf("%v", val)
	}
	return vmap, nil
}

// ToStringInterfaceMap encodes fields tagged with tagName in input into map[string]interface{}. Optional argument properties specifies fields to encode.
func ToStringInterfaceMap(tagName string, input interface{}, properties ...string) (map[string]interface{}, error) {
	vmap := make(map[string]interface{})
	s := structs.New(input)
	s.TagName = tagName
	if len(properties) == 0 {
		properties = s.Names()
	}

	for _, field := range s.Fields() {
		if !field.IsExported() {
			continue
		}

		if !stringInSlice(field.Name(), properties) {
			continue
		}

		tagName, opts := parseTag(field.Tag(tagName))
		if tagName == "" || tagName == "-" {
			continue
		}

		val := field.Value()

		if opts.Has("omitempty") {
			if field.IsZero() {
				continue
			}
			if z, ok := val.(isZeroer); ok && z.IsZero() {
				continue
			}
			if z, ok := val.(isEmptier); ok && z.IsEmpty() {
				continue
			}
		}

		vmap[tagName] = val
	}
	return vmap, nil
}
