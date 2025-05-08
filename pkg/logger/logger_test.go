package logger

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteLeveler(t *testing.T) {
	r := require.New(t)

	l := New("foo")
	r.True(l.Enabled(context.TODO(), slog.LevelInfo))
	r.False(l.Enabled(context.TODO(), slog.LevelDebug))
	err := SetLevel("foo", slog.LevelDebug)
	r.NoError(err)
	r.True(l.Enabled(context.TODO(), slog.LevelDebug))
	err = DeleteLeveler("foo")
	r.NoError(err)
	r.True(l.Enabled(context.TODO(), slog.LevelDebug))
	err = SetLevel("foo", slog.LevelDebug)
	r.ErrorContains(err, "logger not found")
}
