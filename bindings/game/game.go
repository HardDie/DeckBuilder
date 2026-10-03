package game

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/HardDie/DeckBuilder/bindings/catalog"
	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/network"
	servicesGame "github.com/HardDie/DeckBuilder/internal/services/game"
)

type ListResult struct {
	Data []*dto.Game   `json:"data"`
	Meta *network.Meta `json:"meta"`
}

type Result struct {
	Data dto.Game `json:"data"`
	// Warning explains why a new image was not applied. The save itself succeeded.
	Warning string `json:"warning,omitempty"`
}

type Game struct {
	cfg config.Config
	svc servicesGame.Game
	ctx context.Context
}

func New(cfg config.Config, svc servicesGame.Game) *Game {
	return &Game{cfg: cfg, svc: svc}
}

// BindWindow stores the Wails context used by native dialogs.
func BindWindow(g *Game, ctx context.Context) {
	if g == nil {
		return
	}
	g.ctx = ctx
}

func (g *Game) List(sort, search string) (*ListResult, error) {
	items, err := g.svc.List(sort, search)
	if err != nil {
		return nil, err
	}

	respItems := make([]*dto.Game, 0, len(items))
	for _, item := range items {
		d := catalog.GameDTO(g.cfg, *item)
		respItems = append(respItems, &d)
	}

	return &ListResult{
		Data: respItems,
		Meta: &network.Meta{Total: len(respItems)},
	}, nil
}

func (g *Game) Create(req catalog.WriteRequest) (*Result, error) {
	item, err := g.svc.Create(servicesGame.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageBytes(),
	})
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.GameDTO(g.cfg, *item), Warning: catalog.ImageWarning(item.ImageError)}, nil
}

func (g *Game) Read(gameID string) (*Result, error) {
	item, err := g.svc.Item(gameID)
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.GameDTO(g.cfg, *item)}, nil
}

func (g *Game) Update(gameID string, req catalog.WriteRequest) (*Result, error) {
	item, err := g.svc.Update(gameID, servicesGame.UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageBytes(),
	})
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.GameDTO(g.cfg, *item), Warning: catalog.ImageWarning(item.ImageError)}, nil
}

func (g *Game) Delete(gameID string) error {
	return g.svc.Delete(gameID)
}

func (g *Game) Duplicate(gameID, name string) (*Result, error) {
	item, err := g.svc.Duplicate(gameID, servicesGame.DuplicateRequest{Name: name})
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.GameDTO(g.cfg, *item)}, nil
}

// Export asks for a zip path, then writes that game archive there.
// A cancelled dialog returns nil.
func (g *Game) Export(gameID string) error {
	item, err := g.svc.Item(gameID)
	if err != nil {
		return err
	}
	if g.ctx == nil {
		return fmt.Errorf("window is not ready")
	}

	path, err := runtime.SaveFileDialog(g.ctx, runtime.SaveDialogOptions{
		Title:                "Export game",
		DefaultFilename:      exportFilename(item.Name, gameID),
		CanCreateDirectories: true,
		Filters: []runtime.FileFilter{
			{DisplayName: "Zip archive (*.zip)", Pattern: "*.zip"},
		},
	})
	if err != nil {
		return err
	}
	if path == "" {
		return nil
	}
	if !strings.HasSuffix(strings.ToLower(path), ".zip") {
		path += ".zip"
	}

	data, err := g.svc.Export(gameID)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Import creates a game from zip bytes. name is optional.
func (g *Game) Import(name string, data []byte) (*Result, error) {
	if len(data) == 0 {
		return nil, er.BadArchive.AddMessage("The file must be passed as an argument")
	}
	item, err := g.svc.Import(data, name)
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.GameDTO(g.cfg, *item)}, nil
}

func exportFilename(name, gameID string) string {
	name = strings.TrimSpace(strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	).Replace(name))
	name = strings.Trim(name, ". ")
	if name == "" {
		name = gameID
	}
	if name == "" {
		name = "game"
	}
	return name + ".zip"
}
