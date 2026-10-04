package system

import (
	"log/slog"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	repositoriesSettings "github.com/HardDie/DeckBuilder/internal/repositories/settings"
)

type system struct {
	repositorySettings repositoriesSettings.Settings
}

func New(repositorySettings repositoriesSettings.Settings) System {
	return &system{
		repositorySettings: repositorySettings,
	}
}

// GetSettings returns the stored settings, or the defaults when none are stored.
// Missing or invalid values read as their defaults (see Settings.Normalize).
func (s *system) GetSettings() (*entitiesSettings.Settings, error) {
	set, err := s.repositorySettings.Get()
	if err != nil {
		return nil, err
	}
	normalized := set.Normalize()
	return &normalized, nil
}

func (s *system) UpdateSettings(req UpdateSettingsRequest) (*entitiesSettings.Settings, error) {
	if !entitiesSettings.ValidCardScale(req.CardScale) {
		return nil, apperr.ErrBadCardScale
	}
	set, err := s.GetSettings()
	if err != nil {
		return nil, err
	}
	updated := *set
	switch req.Lang {
	case "en", "ru":
		updated.Lang = req.Lang
	}
	updated.EnableBackShadow = req.EnableBackShadow
	updated.CardScale = req.CardScale
	if updated == *set {
		return set, nil
	}
	slog.Info("update settings", "lang", updated.Lang, "back_shadow", updated.EnableBackShadow, "card_scale", updated.CardScale)
	if err := s.repositorySettings.Save(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}
