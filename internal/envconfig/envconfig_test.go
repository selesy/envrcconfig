package envconfig_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/selesy/envrcconfig/internal/envconfig"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update the golden files")

func TestProcess(t *testing.T) {
	const goldenFilePath = "testdata/process-golden.txt"
	type test struct {
		Str      string
		Int      int
		Uint     uint
		Float32  float32
		Float64  float64
		Bool     bool
		Default  string `required:"true"`
		Required string `default:"default"`
		Desc     string `desc:"desc"`
		Alt      string `envconfig:"Alternate"`
	}

	v := &test{}
	vars, err := envconfig.Process("TEST", v)
	require.NoError(t, err)

	act := fmt.Sprintf("%v", vars)
	t.Log(act)

	if *update {
		require.NoError(t, ioutil.WriteFile(goldenFilePath, []byte(act), 0644))
	}

	exp, err := ioutil.ReadFile(goldenFilePath)
	require.NoError(t, err)
	require.Equal(t, string(exp), act)
}
