package system

import (
	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	renderprogress "github.com/HardDie/DeckBuilder/internal/render/progress"
	servicesSystem "github.com/HardDie/DeckBuilder/internal/services/system"
)

type SettingsResult struct {
	Data dto.Settings `json:"data"`
}

type StatusResult struct {
	Data dto.Status `json:"data"`
}

type VersionResult struct {
	Data string `json:"data"`
}

type UpdateSettingsRequest struct {
	Lang string `json:"lang"`
}

type System struct {
	cfg config.Config
	svc servicesSystem.System
}

func New(cfg config.Config, svc servicesSystem.System) *System {
	return &System{cfg: cfg, svc: svc}
}

func (s *System) GetSettings() (*SettingsResult, error) {
	setting, err := s.svc.GetSettings()
	if err != nil {
		return nil, err
	}
	return &SettingsResult{Data: settingsDTO(*setting)}, nil
}

func (s *System) UpdateSettings(req UpdateSettingsRequest) (*SettingsResult, error) {
	setting, err := s.svc.UpdateSettings(servicesSystem.UpdateSettingsRequest{
		Lang: req.Lang,
	})
	if err != nil {
		return nil, err
	}
	return &SettingsResult{Data: settingsDTO(*setting)}, nil
}

func (s *System) Status() *StatusResult {
	status := renderprogress.Get()
	if status.Status == renderprogress.Error || status.Status == renderprogress.Done {
		renderprogress.Reset()
	}
	kind := "No process"
	if status.Status != renderprogress.Empty {
		kind = "Image generation"
	}

	return &StatusResult{Data: dto.Status{
		Type:     kind,
		Progress: status.Percent,
		Status:   status.Status,
	}}
}

func (s *System) GetVersion() *VersionResult {
	return &VersionResult{Data: s.cfg.Version}
}

func settingsDTO(setting entitiesSettings.Settings) dto.Settings {
	return dto.Settings{
		Lang:             setting.Lang,
		EnableBackShadow: setting.EnableBackShadow,
		CardSize: dto.CardSize{
			ScaleX: setting.CardSize.ScaleX,
			ScaleY: setting.CardSize.ScaleY,
			ScaleZ: setting.CardSize.ScaleZ,
		},
	}
}
