package repository

import (
	"errors"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/db"
	"github.com/HardDie/DeckBuilder/internal/entity"
	errors2 "github.com/HardDie/DeckBuilder/internal/errors"
)

type ISystemRepository interface {
	GetSettings() (*entity.SettingInfo, error)
	SaveSettings(set *entity.SettingInfo) error
}
type SystemRepository struct {
	cfg *config.Config
	db  *db.DB
}

func NewSystemRepository(cfg *config.Config, db *db.DB) *SystemRepository {
	return &SystemRepository{
		cfg: cfg,
		db:  db,
	}
}

func (s *SystemRepository) GetSettings() (*entity.SettingInfo, error) {
	resp, err := s.db.SettingsGet()
	if err != nil {
		if errors.Is(err, errors2.SettingsNotExists) {
			return entity.NewSettings(), nil
		} else {
			return nil, err
		}
	}
	return resp, nil
}
func (s *SystemRepository) SaveSettings(set *entity.SettingInfo) error {
	return s.db.SettingsSet(set)
}
