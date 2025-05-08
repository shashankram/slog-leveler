package logger

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/shashankram/slog-leveler/pkg/utils/ptr"
)

func TestDeleteLeveler(t *testing.T) {
	r := require.New(t)
	l := New("delete")
	err := SetLevel("delete", slog.LevelInfo)
	r.NoError(err)
	r.True(l.Enabled(context.TODO(), slog.LevelInfo))
	r.False(l.Enabled(context.TODO(), slog.LevelDebug))
	err = SetLevel("delete", slog.LevelDebug)
	r.NoError(err)
	r.True(l.Enabled(context.TODO(), slog.LevelDebug))
	err = DeleteLeveler("delete")
	r.NoError(err)
	r.True(l.Enabled(context.TODO(), slog.LevelDebug))
	err = SetLevel("delete", slog.LevelDebug)
	r.ErrorContains(err, "logger not found")
}

func TestDefaultLevelInheritence(t *testing.T) {
	r := require.New(t)

	l1 := New("l1")
	l2 := NewWithOptions("l2", Options{Level: ptr.To(slog.LevelDebug)})

	r.True(slog.Default().Enabled(context.TODO(), slog.LevelInfo))
	r.True(l1.Enabled(context.TODO(), slog.LevelInfo))
	r.True(l2.Enabled(context.TODO(), slog.LevelDebug))

	Reset(slog.LevelError)
	r.True(slog.Default().Enabled(context.TODO(), slog.LevelError))
	r.True(l1.Enabled(context.TODO(), slog.LevelError))
	r.True(l2.Enabled(context.TODO(), slog.LevelError))

	l3 := NewWithOptions("l3", Options{Level: ptr.To(slog.LevelDebug)})
	r.True(l3.Enabled(context.TODO(), slog.LevelDebug))
	l4 := New("l4")
	r.True(l4.Enabled(context.TODO(), slog.LevelError))
}
