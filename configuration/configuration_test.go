package configuration

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/smartystreets/assertions"
)

type StrctA struct {
	Foo string `name:"foo" shorthand:"A" description:"strct Foo"`
}

type StrctB struct {
	Foo string `name:"foo" shorthand:"B" description:"strct Foo"`
}

type test struct {
	Bool        bool          `name:"bool"     shorthand:"b" description:"A boolean"`
	String      string        `name:"string"   shorthand:"s" description:"A string"`
	Int         int           `name:"int"      shorthand:"i" description:"An int"`
	StringSlice []string      `name:"strings"  shorthand:"S" description:"A slice of strings"`
	Duration    time.Duration `name:"duration" shorthand:"d" description:"A duration"`
	ConfigFile  string        `name:"config"   shorthand:"c" description:"The location of the config file" config:"yaml"`
	Struct      StrctA        `name:"strct"`
	StructPtr   *StrctB       `name:"strctptr"`
	NotUsed     string
}

var defaults = &test{
	Bool:        true,
	String:      "string",
	Int:         42,
	StringSlice: []string{"a", "b", "c"},
	Duration:    42 * time.Second,
	Struct: StrctA{
		Foo: "foo",
	},
	StructPtr: &StrctB{
		Foo: "bar",
	},
}

func TestBindDefaults(t *testing.T) {
	a := assertions.New(t)

	conf, err := Define("Test", defaults)
	a.So(err, assertions.ShouldBeNil)

	// parse no args
	conf.Parse([]string{})

	res := new(test)

	err = conf.Bind(res)
	a.So(err, assertions.ShouldBeNil)
	a.So(res, assertions.ShouldResemble, defaults)
}

func TestBindArgs(t *testing.T) {
	a := assertions.New(t)

	conf, err := Define("Test", defaults)
	a.So(err, assertions.ShouldBeNil)

	conf.Parse([]string{
		"--bool=false",
		"--string=foo",
		"--int=33",
		"--strings=q",
		"--strings=r",
		"--strings=s",
		"--duration=12m",
		"--strct.foo=baz",
		"--strctptr.foo=qux",
	})

	res := new(test)

	expected := &test{
		Bool:        false,
		String:      "foo",
		Int:         33,
		StringSlice: []string{"q", "r", "s"},
		Duration:    12 * time.Minute,
		Struct: StrctA{
			Foo: "baz",
		},
		StructPtr: &StrctB{
			Foo: "qux",
		},
	}

	err = conf.Bind(res)
	a.So(err, assertions.ShouldBeNil)
	a.So(res, assertions.ShouldResemble, expected)
}

func TestBindEnv(t *testing.T) {
	a := assertions.New(t)

	conf, err := Define("Test", defaults)
	a.So(err, assertions.ShouldBeNil)
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	os.Setenv("BOOL", "0")
	os.Setenv("STRING", "foo")
	os.Setenv("INT", "33")
	os.Setenv("STRINGS", "q r s")
	os.Setenv("DURATION", "12m")
	os.Setenv("STRCT_FOO", "baz")
	os.Setenv("STRCTPTR_FOO", "qux")

	conf.Parse([]string{})

	res := new(test)

	expected := &test{
		Bool:        false,
		String:      "foo",
		Int:         33,
		StringSlice: []string{"q", "r", "s"},
		Duration:    12 * time.Minute,
		Struct: StrctA{
			Foo: "baz",
		},
		StructPtr: &StrctB{
			Foo: "qux",
		},
	}

	err = conf.Bind(res)
	a.So(err, assertions.ShouldBeNil)
	a.So(res, assertions.ShouldResemble, expected)

	os.Setenv("BOOL", "")
	os.Setenv("STRING", "")
	os.Setenv("INT", "")
	os.Setenv("STRINGS", "")
	os.Setenv("DURATION", "")
	os.Setenv("STRCT_FOO", "")
	os.Setenv("STRCTPTR_FOO", "")
}

func TestBindShorthand(t *testing.T) {
	a := assertions.New(t)

	conf, err := Define("Test", defaults)
	a.So(err, assertions.ShouldBeNil)

	conf.Parse([]string{
		"-b=false",
		"-s=foo",
		"-i=33",
		"-S=q",
		"-S=r",
		"-S=s",
		"-d=12m",
		"-A=baz",
		"-B=qux",
	})

	res := new(test)

	expected := &test{
		Bool:        false,
		String:      "foo",
		Int:         33,
		StringSlice: []string{"q", "r", "s"},
		Duration:    12 * time.Minute,
		Struct: StrctA{
			Foo: "baz",
		},
		StructPtr: &StrctB{
			Foo: "qux",
		},
	}

	err = conf.Bind(res)
	a.So(err, assertions.ShouldBeNil)
	a.So(res, assertions.ShouldResemble, expected)
}

func TestBindConfig(t *testing.T) {
	a := assertions.New(t)

	conf, err := Define("Test", defaults)
	a.So(err, assertions.ShouldBeNil)

	conf.SetConfigType("yaml")
	err = conf.ReadConfig(strings.NewReader(`
bool: false
string: foo
int: 41
strings:
  - q
  - r
  - s
duration: 13h
strct:
  foo: baz
strctptr:
  foo: qux
`))

	conf.Parse([]string{})
	a.So(err, assertions.ShouldBeNil)

	res := new(test)

	expected := &test{
		Bool:        false,
		String:      "foo",
		Int:         41,
		StringSlice: []string{"q", "r", "s"},
		Duration:    13 * time.Hour,
		Struct: StrctA{
			Foo: "baz",
		},
		StructPtr: &StrctB{
			Foo: "qux",
		},
	}

	err = conf.Bind(res)
	a.So(err, assertions.ShouldBeNil)
	a.So(res, assertions.ShouldResemble, expected)
}

func TestBindConfigFile(t *testing.T) {
	a := assertions.New(t)

	conf, err := Define("Test", defaults)
	a.So(err, assertions.ShouldBeNil)

	conf.Parse([]string{
		"--config=./foo.yaml",
	})

	res := new(test)

	err = conf.Bind(res)
	a.So(err, assertions.ShouldNotBeNil)
	a.So(err.Error(), assertions.ShouldContainSubstring, "./foo.yaml")
	a.So(err.Error(), assertions.ShouldContainSubstring, "no such file")
}

func TestUnsupported(t *testing.T) {
	a := assertions.New(t)

	type Config struct {
		Unsupported complex64 `name:"unsupported"`
	}

	_, err := Define("Unsupported", &Config{
		Unsupported: complex(1, 0),
	})

	a.So(err, assertions.ShouldNotBeNil)
}
