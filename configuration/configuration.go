// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package configuration

import (
	"fmt"
	"io"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Definition is a configuration definition
type Definition struct {
	defaults interface{}
	pf       *pflag.FlagSet
	viper    *viper.Viper
}

// Define a new config struct with the structure of the type of the passed
// interface{} that used as default values of the fields of the passed value
// the passed interface
func Define(name string, defaults interface{}) (*Definition, error) {
	c := &Definition{
		defaults: defaults,
		pf:       pflag.NewFlagSet(name, pflag.ExitOnError),
		viper:    viper.New(),
	}

	return c, c.define("", c.defaults)
}

// Parse the arguments
func (c *Definition) Parse(args []string) {
	c.pf.Parse(args)
	c.viper.AutomaticEnv()
	c.viper.SetTypeByDefaultValue(true)
	c.viper.BindPFlags(c.pf)
}

// Bind the parsed arguments into the passed struct
// this struct must be of the same type as the struct used
// to Define this configuration
func (c *Definition) Bind(v interface{}) error {
	typ := reflect.ValueOf(v)
	if typ.Kind() != reflect.Ptr || typ.IsNil() {
		return &InvalidBindError{typ: reflect.TypeOf(v)}
	}

	d := reflect.ValueOf(c.defaults)

	if typ.Type() != d.Type() {
		return fmt.Errorf("Mismatched types of definition (%s) and bound (%s)", typ.Type(), d.Type())
	}

	ve := typ.Elem()
	vt := ve.Type()

	// find the key that is the configfile
	for i := 0; i < vt.NumField(); i++ {
		et := vt.Field(i)
		name := et.Tag.Get("name")

		if name == "" {
			continue
		}

		config := et.Tag.Get("config")
		if config == "" {
			continue
		}

		// read the config var to get to now the filename
		file, ok := c.Get(name).(string)
		if ok && file != "" {
			c.viper.SetConfigType(config)
			c.viper.SetConfigFile(file)
			err := c.viper.ReadInConfig()
			if err != nil {
				return err
			}
		}
	}

	return c.bind("", ve)
}

func (c *Definition) bind(prefix string, ve reflect.Value) error {
	vt := ve.Type()

	// bind all the keys
	for i := 0; i < vt.NumField(); i++ {
		ef := ve.Field(i)
		et := vt.Field(i)

		name := et.Tag.Get("name")

		// skip if no name is defined
		if name == "" {
			continue
		}

		if prefix != "" {
			name = prefix + "." + name
		}

		if ef.Kind() == reflect.Struct {
			// recurse into struct
			err := c.bind(name, ef)
			if err != nil {
				return err
			}
		} else if et.Type.Kind() == reflect.Ptr && et.Type.Elem().Kind() == reflect.Struct {
			// recurse into struct pointer
			if ef.IsNil() {
				ef.Set(reflect.New(et.Type.Elem()))
			}
			err := c.bind(name, ef.Elem())
			if err != nil {
				return err
			}
		} else {
			got := c.get(name)
			if got.IsValid() {
				ef.Set(got)
			}
		}
	}

	return nil
}

// Get returns the value of the key name
func (c *Definition) Get(name string) interface{} {
	return c.get(name).Interface()
}

// SetEnvPrefix sets the prefix to use for environment variables
func (c *Definition) SetEnvPrefix(prefix string) {
	c.viper.SetEnvPrefix(prefix)
}

// SetEnvKeyReplacer sets the replacer to use for environment varialbe lookups
func (c *Definition) SetEnvKeyReplacer(replacer *strings.Replacer) {
	c.viper.SetEnvKeyReplacer(replacer)
}

// SetConfigType sets the type of the config file (if any)
func (c *Definition) SetConfigType(typ string) {
	c.viper.SetConfigType(typ)
}

// AddConfigPath adds a config path to look for config files
func (c *Definition) AddConfigPath(path string) {
	c.viper.AddConfigPath(path)
}

// ReadConfig reads the config from an io.Reader
func (c *Definition) ReadConfig(in io.Reader) error {
	return c.viper.ReadConfig(in)
}

// ConfigFileUsed returns the path to the config file that was used
func (c *Definition) ConfigFileUsed() string {
	return c.viper.ConfigFileUsed()
}

func (c *Definition) define(prefix string, defaults interface{}) error {
	dv := reflect.ValueOf(defaults)
	if dv.Type().Kind() == reflect.Ptr {
		dv = dv.Elem()
	}

	dt := dv.Type()

	for i := 0; i < dt.NumField(); i++ {
		ft := dt.Field(i)
		fv := dv.Field(i)

		name := ft.Tag.Get("name")
		shorthand := ft.Tag.Get("shorthand")
		usage := ft.Tag.Get("description")

		// skip if no name is defined
		if name == "" {
			continue
		}

		if prefix != "" {
			name = prefix + "." + name
		}

		c.viper.SetDefault(name, fv.Interface())

		switch t := fv.Interface().(type) {
		case bool:
			c.pf.BoolP(name, shorthand, t, usage)
		case time.Duration:
			c.pf.DurationP(name, shorthand, t, usage)
		case float32:
			c.pf.Float32P(name, shorthand, t, usage)
		case float64:
			c.pf.Float64P(name, shorthand, t, usage)
		case net.IP:
			c.pf.IPP(name, shorthand, t, usage)
		case net.IPNet:
			c.pf.IPNetP(name, shorthand, t, usage)
		case net.IPMask:
			c.pf.IPMaskP(name, shorthand, t, usage)
		case int:
			c.pf.IntP(name, shorthand, t, usage)
		case int32:
			c.pf.Int32P(name, shorthand, t, usage)
		case int64:
			c.pf.Int64P(name, shorthand, t, usage)
		case int8:
			c.pf.Int8P(name, shorthand, t, usage)
		case []int:
			c.pf.IntSliceP(name, shorthand, t, usage)
		case string:
			c.pf.StringP(name, shorthand, t, usage)
		case []string:
			c.pf.StringSliceP(name, shorthand, t, usage)
		case uint32:
			c.pf.Uint32P(name, shorthand, t, usage)
		case uint64:
			c.pf.Uint64P(name, shorthand, t, usage)
		case uint8:
			c.pf.Uint8P(name, shorthand, t, usage)
		default:
			switch {
			case fv.Kind() == reflect.Struct:
				// recurse into struct
				err := c.define(name, fv.Interface())
				if err != nil {
					return err
				}
				continue

			case fv.Type().Kind() == reflect.Ptr && fv.Type().Elem().Kind() == reflect.Struct:
				// recurse into struct pointer
				err := c.define(name, fv.Elem().Interface())
				if err != nil {
					return err
				}
				continue

			default:
				return &UnsupportedTypeError{
					Type: fv.Type(),
				}
			}
		}
	}

	return nil
}

func (c *Definition) get(name string) reflect.Value {
	v := c.viper.Get(name)
	return reflect.ValueOf(v)
}
