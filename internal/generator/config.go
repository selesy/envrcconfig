package generator

import (
	"errors"
	"flag"
	"sort"
	"strings"
)

type Format int

const (
	DirEnv Format = iota
	DotEnv
	Kubernetes
	Terraform
)

var ErrUnsupportedOutputFormat = errors.New("unsupported output format")

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

func (f Formats) String() string {
	str := []string{}

	for fmt := range f {
		str = append(str, fmt.String())
	}

	sort.StringSlice(str).Sort()

	return strings.Join(str, ", ")
}

type Config struct {
	formats  Formats
	help     bool
	logLevel string
	targets  []string
	version  bool
}

func (c *Config) Formats() Formats {
	return c.formats
}

func (c *Config) Help() bool {
	return c.help
}

func (c *Config) LogLevel() string {
	return c.logLevel
}

func (c *Config) Targets() []string {
	return c.targets
}

func (c *Config) Version() bool {
	return c.version
}

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
