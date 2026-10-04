package system

import (
	"context"
	"io"
	"log/slog"
	"math"
	"testing"

	"github.com/HardDie/fsentry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	repositoriesSettings "github.com/HardDie/DeckBuilder/internal/repositories/settings"
	"github.com/HardDie/DeckBuilder/pkg/logger"
)

func newSystem(t *testing.T) System {
	t.Helper()
	s, _ := newSystemDB(t)
	return s
}

func newSystemDB(t *testing.T) (System, *fsentry.DB) {
	t.Helper()
	cfg := config.Get("")
	cfg.SetDataPath(t.TempDir())
	db := fsentry.New(cfg.Data, fsentry.WithPretty(), fsentry.WithNoLockFile())
	require.NoError(t, db.Init())
	return New(repositoriesSettings.New(cfg, db)), db
}

func TestSettings(t *testing.T) {
	s := newSystem(t)

	got, err := s.GetSettings()
	require.NoError(t, err)
	assert.Equal(t, entitiesSettings.Default(), *got, "no settings file gives the defaults")

	updated, err := s.UpdateSettings(UpdateSettingsRequest{Lang: "ru", CardScale: 1})
	require.NoError(t, err)
	assert.Equal(t, "ru", updated.Lang)
	got, err = s.GetSettings()
	require.NoError(t, err)
	assert.Equal(t, "ru", got.Lang, "the saved language reads back")
	assert.Equal(t, entitiesSettings.Default().CardScale, got.CardScale, "other fields keep their values")

	updated, err = s.UpdateSettings(UpdateSettingsRequest{Lang: "fr", CardScale: 1})
	require.NoError(t, err)
	assert.Equal(t, "ru", updated.Lang, "an unknown language is ignored")

	// The dialog sends no language: the stored one stays.
	updated, err = s.UpdateSettings(UpdateSettingsRequest{EnableBackShadow: true, CardScale: 2.5})
	require.NoError(t, err)
	want := entitiesSettings.Settings{Lang: "ru", EnableBackShadow: true, CardScale: 2.5, LogLevel: entitiesSettings.LogLevelInfo}
	assert.Equal(t, want, *updated)
	got, err = s.GetSettings()
	require.NoError(t, err)
	assert.Equal(t, want, *got, "scale and shadow read back")
}

func TestUpdateSettingsBadCardScale(t *testing.T) {
	s := newSystem(t)
	for _, scale := range []float64{0, 0.05, 10.5, -1, math.NaN(), math.Inf(1)} {
		_, err := s.UpdateSettings(UpdateSettingsRequest{CardScale: scale})
		assert.ErrorIs(t, err, apperr.ErrBadCardScale, "scale %v", scale)
	}
	got, err := s.GetSettings()
	require.NoError(t, err)
	assert.Equal(t, entitiesSettings.Default(), *got, "a rejected update saves nothing")
}

// The stored file keeps the older x/y/z shape: X and Z carry the scale, Y is 1.
func TestSettingsStoredShape(t *testing.T) {
	s, db := newSystemDB(t)
	_, err := s.UpdateSettings(UpdateSettingsRequest{CardScale: 1.5})
	require.NoError(t, err)

	type stored struct {
		CardSize map[string]float64 `json:"card_size"`
	}
	entry, err := db.GetEntry[stored]("settings")
	require.NoError(t, err)
	assert.Equal(t, map[string]float64{"scaleX": 1.5, "scaleY": 1, "scaleZ": 1.5}, entry.Data.CardSize)
}

// An older or hand-edited settings.json may miss fields; they read as defaults.
func TestSettingsMissingFields(t *testing.T) {
	s, db := newSystemDB(t)
	_, err := db.CreateEntry("settings", map[string]any{"lang": "ru"})
	require.NoError(t, err)

	got, err := s.GetSettings()
	require.NoError(t, err)
	assert.Equal(t, "ru", got.Lang)
	assert.Equal(t, entitiesSettings.Default().CardScale, got.CardScale, "missing card size is the default, not 0")
}

// A saved log level is stored, reads back, and applies to the logger at once.
// An empty or unknown level keeps the stored one.
func TestUpdateSettingsLogLevel(t *testing.T) {
	_, err := logger.Init(logger.Options{Console: io.Discard})
	require.NoError(t, err)
	t.Cleanup(func() { _ = logger.Close() })
	debugOn := func() bool { return slog.Default().Enabled(context.Background(), slog.LevelDebug) }

	s := newSystem(t)
	require.False(t, debugOn(), "Info by default")

	updated, err := s.UpdateSettings(UpdateSettingsRequest{CardScale: 1, LogLevel: entitiesSettings.LogLevelDebug})
	require.NoError(t, err)
	assert.Equal(t, entitiesSettings.LogLevelDebug, updated.LogLevel)
	assert.True(t, debugOn(), "debug applies without a restart")

	for _, level := range []string{"", "trace"} {
		updated, err = s.UpdateSettings(UpdateSettingsRequest{CardScale: 1, LogLevel: level})
		require.NoError(t, err)
		assert.Equal(t, entitiesSettings.LogLevelDebug, updated.LogLevel, "level %q keeps the stored one", level)
	}
	got, err := s.GetSettings()
	require.NoError(t, err)
	assert.Equal(t, entitiesSettings.LogLevelDebug, got.LogLevel, "the level reads back")

	_, err = s.UpdateSettings(UpdateSettingsRequest{CardScale: 1, LogLevel: entitiesSettings.LogLevelError})
	require.NoError(t, err)
	ctx := context.Background()
	assert.False(t, slog.Default().Enabled(ctx, slog.LevelWarn), "error hides warn")
	assert.True(t, slog.Default().Enabled(ctx, slog.LevelError), "error keeps error")

	_, err = s.UpdateSettings(UpdateSettingsRequest{CardScale: 1, LogLevel: entitiesSettings.LogLevelInfo})
	require.NoError(t, err)
	assert.False(t, debugOn(), "back to Info")
}
