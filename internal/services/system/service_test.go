package system

import (
	"testing"

	"github.com/HardDie/fsentry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	repositoriesSettings "github.com/HardDie/DeckBuilder/internal/repositories/settings"
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

	updated, err := s.UpdateSettings(UpdateSettingsRequest{Lang: "ru"})
	require.NoError(t, err)
	assert.Equal(t, "ru", updated.Lang)
	got, err = s.GetSettings()
	require.NoError(t, err)
	assert.Equal(t, "ru", got.Lang, "the saved language reads back")
	assert.Equal(t, entitiesSettings.Default().CardSize, got.CardSize, "other fields keep their values")

	updated, err = s.UpdateSettings(UpdateSettingsRequest{Lang: "fr"})
	require.NoError(t, err)
	assert.Equal(t, "ru", updated.Lang, "an unknown language is ignored")
}

// An older or hand-edited settings.json may miss fields; they read as defaults.
func TestSettingsMissingFields(t *testing.T) {
	s, db := newSystemDB(t)
	_, err := db.CreateEntry("settings", map[string]any{"lang": "ru"})
	require.NoError(t, err)

	got, err := s.GetSettings()
	require.NoError(t, err)
	assert.Equal(t, "ru", got.Lang)
	assert.Equal(t, entitiesSettings.Default().CardSize, got.CardSize, "missing card size is the default, not 0")
}
