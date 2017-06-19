package configuration

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func ExampleConfigDefinition() {
	// Log is a nested struct
	type Log struct {
		Level string `name:"The log level"`
	}

	// Config is the shape of the configuration for this program
	type Config struct {
		// Foo is the most basic variable definition
		Foo string `name:"foo" description:"The Foo variable"`

		// Dur has a shorthand declared
		Dur time.Duration `name:"dur" description:"The Dur variable" shorthand:"d"`

		// Log is a nested struct, all variables defined inside of it will be
		// prefixed with log.
		Log *Log `name:"log"`

		// ConfigFile has the `config:"yaml"` tag on it and so will be interpreted
		// as the location of the config file, which will be read in automatically
		ConfigFile string `name:"config" description:"The location of the config file" config:"yaml"`

		// NotUsed does not have a name tag and will be ignored by the configuration
		// package
		NotUsed string
	}

	// Define the configuration defaults
	c, err := Define("Test program", &Config{
		Foo: "foo",
		Dur: 1 * time.Second,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Set the environment prefix, so foo would become EXAMPLE_FOO
	c.SetEnvPrefix("EXAMPLE")
	c.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Parse the arguments
	c.Parse(os.Args)

	// create a new config struct to load the configuration into
	config := new(Config)

	// load the configuration, be it from env variables, command line args or
	// config file
	err = c.Bind(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the variables
	fmt.Println("Foo", config.Foo)
	fmt.Println("Dur", config.Dur)
	fmt.Println("Log.Level", config.Log.Level)
}
