package logger

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	tests := []struct {
		name string
		opts Options
		want Options
	}{
		{
			name: "default options",
			opts: Options{},
			want: Options{
				Level:  slog.LevelInfo,
				Format: TextFormat,
				Writer: os.Stderr,
			},
		},
		{
			name: "custom options",
			opts: Options{
				Level:     slog.LevelDebug,
				Format:    JSONFormat,
				Writer:    nil,
				AddSource: true,
			},
			want: Options{
				Level:     slog.LevelDebug,
				Format:    JSONFormat,
				Writer:    os.Stderr,
				AddSource: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			tt.opts.Default()
			a.Equal(tt.want, tt.opts)
		})
	}
}
