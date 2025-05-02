package logger

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	tests := []struct {
		name           string
		components     []string
		query          string
		setLevel       map[string]slog.Level
		wantStatusCode int
		wantBody       string
		wantLevels     map[string]slog.Level
	}{
		{
			name:           "only default logger",
			wantStatusCode: http.StatusOK,
			wantLevels: map[string]slog.Level{
				DefaultComponent: GlobalLevel.Level(),
			},
		},
		{
			name:           "update default level to debug",
			query:          "level=debug",
			wantStatusCode: http.StatusOK,
			wantLevels: map[string]slog.Level{
				DefaultComponent: slog.LevelDebug,
			},
		},
		{
			name:           "update all loggers to debug level",
			components:     []string{"c1", "c2", "c3"},
			query:          "level=debug",
			wantStatusCode: http.StatusOK,
			wantLevels: map[string]slog.Level{
				DefaultComponent: slog.LevelDebug,
				"c1":             slog.LevelDebug,
				"c2":             slog.LevelDebug,
				"c3":             slog.LevelDebug,
			},
		},
		{
			name:           "ignore component levels when updating specific logger levels",
			components:     []string{"c1", "c2", "c3"},
			query:          "level=debug&c1=error&c2=warn&c3=trace",
			wantStatusCode: http.StatusOK,
			wantLevels: map[string]slog.Level{
				DefaultComponent: slog.LevelDebug,
				"c1":             slog.LevelDebug,
				"c2":             slog.LevelDebug,
				"c3":             slog.LevelDebug,
			},
		},
		{
			name:           "update default and component levels",
			components:     []string{"c1", "c2", "c3"},
			query:          "default=debug&c1=error&c2=warn&c3=trace",
			wantStatusCode: http.StatusOK,
			wantLevels: map[string]slog.Level{
				DefaultComponent: slog.LevelDebug,
				"c1":             slog.LevelError,
				"c2":             slog.LevelWarn,
				"c3":             LevelTrace,
			},
		},
		{
			name:           "incorrect global log level should error and preserve current level",
			query:          "level=foo",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "unknown log level foo",
			wantLevels: map[string]slog.Level{
				DefaultComponent: slog.LevelInfo,
			},
		},
		{
			name:           "incorrect component log level should error and preserve current level",
			components:     []string{"c1"},
			query:          "c1=foo",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "component c1: unknown log level foo",
			wantLevels: map[string]slog.Level{
				"c1": slog.LevelInfo,
			},
		},
		{
			name:           "update default and component levels using SetLevel",
			components:     []string{"c1", "c2", "c3"},
			setLevel:       map[string]slog.Level{"default": slog.LevelDebug, "c1": slog.LevelError, "c2": slog.LevelWarn, "c3": LevelTrace},
			wantStatusCode: http.StatusOK,
			wantLevels: map[string]slog.Level{
				DefaultComponent: slog.LevelDebug,
				"c1":             slog.LevelError,
				"c2":             slog.LevelWarn,
				"c3":             LevelTrace,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a := assert.New(t)

			// Reset component levels to default level
			Reset(slog.LevelInfo)

			loggers := map[string]*slog.Logger{DefaultComponent: slog.Default()}
			for _, component := range tc.components {
				logger := New(component)
				a.NotNil(logger)
				loggers[component] = logger
			}

			// Test HTTP handler
			path := "/logging"
			if tc.query != "" {
				path += "?" + tc.query
			}
			req := httptest.NewRequest(http.MethodPost, path, nil)
			w := httptest.NewRecorder()
			HTTPLevelHandler(w, req)
			resp := w.Result()
			a.Equal(tc.wantStatusCode, resp.StatusCode)
			data, err := io.ReadAll(resp.Body)
			a.NoError(err)
			a.NotEmpty(data)
			a.Contains(string(data), tc.wantBody)

			// Test SetLevel
			for component, level := range tc.setLevel {
				err := SetLevel(component, level)
				a.NoError(err)
			}

			for component, level := range tc.wantLevels {
				a.Equal(level, MustGetLevel(component), component)
				a.True(loggers[component].Enabled(context.TODO(), level), component)
			}
		})
	}
}
