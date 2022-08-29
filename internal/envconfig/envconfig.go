package envconfig

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/kelseyhightower/envconfig"
)

const defaultTableFormat = `{{range .}}{{.Name}}	{{.Alt}}	{{usage_key .}}	{{usage_type .}}	{{usage_default .}}	{{usage_required .}}	{{usage_description .}}
{{end}}`

// VarInfo contains the parsed specification for each environment variable
// envconfig would parse.
//
// This is somewhat analogous to the unexported varInfo struct in the
// envconfig project but since the Usage() methods already parse the
// spec's field tags, we'll capture them at the same time.
type VarInfo struct {
	Name     string
	Alt      string
	Key      string
	Type     string
	Default  string
	Required bool
	Desc     string
}

// Process returns VarInfo data for the given prefix and spec.
//
// This function is named, and has arguments, to purposely reflect the
// equivalent functionality in the envconfig package.  The difference
// is that this version returns the information for each environment
// variable instead of processing environment variables into the spec
// struct.
//
// This information is eventually used to: a) use the AST to look up the
// struct and field comments for the spec (for use as documentation if no
// desc tag is provided) and b) to provide data to the template writers
// that generate the sample output files.
func Process(prefix string, spec any) ([]VarInfo, error) {
	out := &bytes.Buffer{}
	tabs := tabwriter.NewWriter(out, 1, 256, 0, '\t', 0)

	err := envconfig.Usagef(prefix, spec, tabs, defaultTableFormat)
	if err != nil {
		return []VarInfo{}, err
	}

	tabs.Flush()

	vars := []VarInfo{}

	for _, line := range strings.Split(out.String(), "\n") {
		if line == "" {
			continue
		}

		v, err := parse(line)
		if err != nil {
			return nil, err
		}

		vars = append(vars, v)
	}

	return vars, nil
}

var ErrIncompatibleEnvConfig = errors.New("incompatible envconfig (this should not happen unless the version has changed")
var ErrExpectFiveTokens = fmt.Errorf("%w - each line from envconfig Usage() should contain five tokens", ErrIncompatibleEnvConfig)

func parse(line string) (VarInfo, error) {
	tkns := strings.Split(line, "\t")
	if len(tkns) != 7 {
		return VarInfo{}, ErrExpectFiveTokens
	}

	if tkns[5] == "" {
		tkns[5] = "false"
	}

	req, err := strconv.ParseBool(tkns[5])
	if err != nil {
		return VarInfo{}, err
	}

	return VarInfo{
		Name:     tkns[0],
		Alt:      tkns[1],
		Key:      tkns[2],
		Type:     tkns[3],
		Default:  tkns[4],
		Required: req,
		Desc:     tkns[6],
	}, nil
}
