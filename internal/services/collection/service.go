package collection

import (
	"strings"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
	repositoriesCollection "github.com/HardDie/DeckBuilder/internal/repositories/collection"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type collection struct {
	cfg                  *config.Config
	repositoryCollection repositoriesCollection.Collection
}

func New(cfg *config.Config, repositoryCollection repositoriesCollection.Collection) Collection {
	return &collection{
		cfg:                  cfg,
		repositoryCollection: repositoryCollection,
	}
}

func (s *collection) Create(gameID string, req CreateRequest) (*entity.CollectionInfo, error) {
	c, err := s.repositoryCollection.Create(gameID, repositoriesCollection.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageFile,
	})
	if err != nil {
		return nil, err
	}
	c.FillCachedImage(s.cfg, gameID)
	return c, nil
}
func (s *collection) Item(gameID, collectionID string) (*entity.CollectionInfo, error) {
	c, err := s.repositoryCollection.GetByID(gameID, collectionID)
	if err != nil {
		return nil, err
	}
	c.FillCachedImage(s.cfg, gameID)
	return c, nil
}
func (s *collection) List(gameID, sortField, search string) ([]*entity.CollectionInfo, *network.Meta, error) {
	items, err := s.repositoryCollection.GetAll(gameID)
	if err != nil {
		return make([]*entity.CollectionInfo, 0), nil, err
	}

	// Filter
	var filteredItems []*entity.CollectionInfo
	if search != "" {
		search = strings.ToLower(search)
		for _, item := range items {
			if strings.Contains(strings.ToLower(item.Name), search) {
				filteredItems = append(filteredItems, item)
			}
		}
	} else {
		filteredItems = items
	}

	// Sorting
	utils.Sort(&filteredItems, sortField)

	// Generate field cachedImage
	for i := 0; i < len(filteredItems); i++ {
		filteredItems[i].FillCachedImage(s.cfg, gameID)
	}

	// Return empty array if no elements
	if filteredItems == nil {
		filteredItems = make([]*entity.CollectionInfo, 0)
	}

	meta := &network.Meta{
		Total: len(filteredItems),
	}
	return filteredItems, meta, nil
}
func (s *collection) Update(gameID, collectionID string, req UpdateRequest) (*entity.CollectionInfo, error) {
	c, err := s.repositoryCollection.Update(gameID, collectionID, repositoriesCollection.UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageFile,
	})
	if err != nil {
		return nil, err
	}
	c.FillCachedImage(s.cfg, gameID)
	return c, nil
}
func (s *collection) Delete(gameID, collectionID string) error {
	return s.repositoryCollection.DeleteByID(gameID, collectionID)
}
func (s *collection) GetImage(gameID, collectionID string) ([]byte, string, error) {
	return s.repositoryCollection.GetImage(gameID, collectionID)
}
