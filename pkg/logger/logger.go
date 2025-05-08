package logger

import (
	"fmt"
	"log/slog"
	"sync"
)

const (
	DefaultComponent = "default"
)

// componentLeveler maps component names to their respective slog.LevelVar instance
var componentLeveler sync.Map

func init() {
	defaultLogger := New(DefaultComponent)
	slog.SetDefault(defaultLogger)
}

// New returns a new slog.Logger instance for the given component with default Options.
// If the component is empty, it returns the default logger.
func New(component string) *slog.Logger {
	return NewWithOptions(component, Options{})
}

// NewWithOptions returns a new slog.Logger instance for the given component with the provided Options
// If the component is empty, it returns the default logger.
func NewWithOptions(component string, opts Options) *slog.Logger {
	if component == "" {
		return slog.Default()
	}

	opts.Default()

	level := &slog.LevelVar{}
	if opts.Level != nil {
		level.Set(*opts.Level)
	} else {
		defaultLvl, ok := componentLeveler.Load(DefaultComponent)
		if ok {
			level.Set(defaultLvl.(*slog.LevelVar).Level())
		}
	}

	handlerOpts := &slog.HandlerOptions{
		AddSource:   opts.AddSource,
		Level:       level,
		ReplaceAttr: slogLevelReplacer,
	}

	attrs := []slog.Attr{{Key: "component", Value: slog.StringValue(component)}}

	componentLeveler.Store(component, level)
	var slogHandler slog.Handler
	switch opts.Format {
	case TextFormat:
		slogHandler = slog.NewTextHandler(opts.Writer, handlerOpts).WithAttrs(attrs)
	case JSONFormat:
		slogHandler = slog.NewJSONHandler(opts.Writer, handlerOpts).WithAttrs(attrs)
	default:
		slogHandler = slog.NewTextHandler(opts.Writer, handlerOpts).WithAttrs(attrs)
	}

	return slog.New(slogHandler)
}

// DeleteLeveler deletes the leveler instance for the given component
func DeleteLeveler(component string) error {
	if component == "" {
		return fmt.Errorf("component unspecified")
	}
	componentLeveler.Delete(component)
	return nil
}
