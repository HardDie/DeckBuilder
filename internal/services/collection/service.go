package collection

import (
	"strings"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
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

func (s *collection) Create(gameID string, dtoObject *dto.CreateCollectionDTO) (*entity.CollectionInfo, error) {
	collection, err := s.repositoryCollection.Create(gameID, dtoObject)
	if err != nil {
		return nil, err
	}
	collection.FillCachedImage(s.cfg, gameID)
	return collection, nil
}
func (s *collection) Item(gameID, collectionID string) (*entity.CollectionInfo, error) {
	collection, err := s.repositoryCollection.GetByID(gameID, collectionID)
	if err != nil {
		return nil, err
	}
	collection.FillCachedImage(s.cfg, gameID)
	return collection, nil
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
func (s *collection) Update(gameID, collectionID string, dtoObject *dto.UpdateCollectionDTO) (*entity.CollectionInfo, error) {
	collection, err := s.repositoryCollection.Update(gameID, collectionID, dtoObject)
	if err != nil {
		return nil, err
	}
	collection.FillCachedImage(s.cfg, gameID)
	return collection, nil
}
func (s *collection) Delete(gameID, collectionID string) error {
	return s.repositoryCollection.DeleteByID(gameID, collectionID)
}
func (s *collection) GetImage(gameID, collectionID string) ([]byte, string, error) {
	return s.repositoryCollection.GetImage(gameID, collectionID)
}
