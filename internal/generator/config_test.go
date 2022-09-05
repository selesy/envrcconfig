package generator_test

import (
	"testing"

	"github.com/selesy/envrcconfig/internal/generator"
	"github.com/stretchr/testify/require"
)

func TestParseFormat(t *testing.T) {
	t.Run("with valid name (trimmed)", func(t *testing.T) {
		fmt, err := generator.ParseFormat("DirEnv ")
		require.NoError(t, err)
		require.Equal(t, generator.DirEnv, fmt)
	})

	t.Run("with unknown name", func(t *testing.T) {
		_, err := generator.ParseFormat("Unknown")
		require.EqualError(t, err, generator.ErrUnsupportedOutputFormat.Error())
	})
}

func TestFormat_String(t *testing.T) {
	t.Run("with valid format", func(t *testing.T) {
		require.Equal(t, "direnv", generator.DirEnv.String())
	})

	t.Run("with invalid format", func(t *testing.T) {
		require.Equal(t, "", generator.Format(42).String())
	})
}

func TestFormats_Set(t *testing.T) {
	t.Run("with valid input", func(t *testing.T) {
		fmts := generator.Formats{}
		err := fmts.Set("direnv,dotenv")
		require.NoError(t, err)
		require.Contains(t, fmts, generator.DirEnv)
		require.Contains(t, fmts, generator.DotEnv)
	})

	t.Run("with invalid input", func(t *testing.T) {
		fmts := generator.Formats{}
		err := fmts.Set("direnv,unknown")
		require.EqualError(t, err, generator.ErrUnsupportedOutputFormat.Error())
	})
}

func TestFormats_String(t *testing.T) {
	t.Run("with valid formats", func(t *testing.T) {
		fmts := generator.Formats{
			generator.DirEnv: struct{}{},
			generator.DotEnv: struct{}{},
		}

		require.Equal(t, "direnv, dotenv", fmts.String())
	})
}

func TestNewConfig(t *testing.T) {
	t.Run("output version", func(t *testing.T) {
		args := []string{
			"program",
			"-version",
		}

		cfg, err := generator.NewConfig(args)
		require.NoError(t, err)
		require.Nil(t, cfg.Formats())
		require.False(t, cfg.Help())
		require.Equal(t, "INFO", cfg.LogLevel())
		require.Empty(t, cfg.Targets())
		require.True(t, cfg.Version())
	})
}
