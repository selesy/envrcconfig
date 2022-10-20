package generator

import (
	"errors"
	"flag"
	"sort"
	"strings"
)

// Format provides an enumeration listing the various output formats that
// envrcconfig can generate.
type Format int

const (
	DirEnv Format = iota
	DotEnv
	Kubernetes
	Terraform
)

// ErrUnsupportedOutputFormat is returned while building a configuration if
// the user creates a format flag that is unknown.
var ErrUnsupportedOutputFormat = errors.New("unsupported output format")

// ParseFormat attempts to match strings passed via the format flag to valid
// values of the enumeration.
func ParseFormat(v string) (Format, error) {
	fmt, ok := map[string]Format{
		"direnv":     DirEnv,
		"dotenv":     DotEnv,
		"kubernetes": Kubernetes,
		"terraform":  Terraform,
	}[strings.TrimSpace(strings.ToLower(v))]

	if !ok {
		return DirEnv, ErrUnsupportedOutputFormat
	}

	return fmt, nil
}

// String returns the textual value associated with the enumeration's
// ordinal value.
func (f Format) String() string {
	if f < 0 || f > Terraform {
		return ""
	}

	return []string{
		"direnv",
		"dotenv",
		"kubernetes",
		"terraform",
	}[int(f)]
}

// Formats is a simple type to adapt a comma-separated list of format
// values, which may contain duplicates, into a set of unique valid
// Format values.
type Formats map[Format]struct{}

func (f Formats) Set(v string) error {
	for _, val := range strings.Split(strings.ToLower(v), ",") {
		fmt, err := ParseFormat(val)
		if err != nil {
			return err
		}

		f[fmt] = struct{}{}
	}
	return nil
}

// String returns an ascending, sorted list of the formats that were
// passed on the command line.
func (f Formats) String() string {
	str := []string{}

	for fmt := range f {
		str = append(str, fmt.String())
	}

	sort.StringSlice(str).Sort()

	return strings.Join(str, ", ")
}

// Config contains the parsed configuration passed to the application
// via the command line.  After parsing CLI arguments, this type is
// effectively immuntable.
type Config struct {
	formats  Formats
	help     bool
	logLevel string
	targets  []string
	version  bool
}

// Formats returns the value of the configuration's formats field.
func (c *Config) Formats() Formats {
	return c.formats
}

// Help returns the value of the configuration's help field.
func (c *Config) Help() bool {
	return c.help
}

// LogLevel returns the value of the configuration's logLevel field.
func (c *Config) LogLevel() string {
	return c.logLevel
}

// Targets returns the value of the configuration's targets field which
// contains the name of each envconfig annotated type in the form
// <package>,<name>.
func (c *Config) Targets() []string {
	return c.targets
}

// Version returns the value of the configuration's version field.
func (c *Config) Version() bool {
	return c.version
}

// NewConfig parse the command line's args and attempts to build a valid
// Config.
func NewConfig(args []string) (*Config, error) {
	cfg := &Config{}

	cli := flag.NewFlagSet(args[0], flag.ContinueOnError)
	cli.Var(cfg.formats, "formats", "formats help")
	cli.BoolVar(&cfg.help, "help", false, "help help")
	cli.StringVar(&cfg.logLevel, "logging", "INFO", "log level help")
	cli.BoolVar(&cfg.version, "version", false, "version help")

	err := cli.Parse(args[1:])
	if err != nil {
		return nil, err
	}

	cfg.targets = cli.Args()

	// TODO: add validation

	return cfg, nil
}
