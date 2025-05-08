package logger

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/shashankram/slog-leveler/pkg/utils/ptr"
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
				Format: TextFormat,
				Writer: os.Stderr,
			},
		},
		{
			name: "custom options",
			opts: Options{
				Level:     ptr.To(slog.LevelDebug),
				Format:    JSONFormat,
				Writer:    nil,
				AddSource: true,
			},
			want: Options{
				Level:     ptr.To(slog.LevelDebug),
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
