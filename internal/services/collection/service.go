package collection

import (
	"strings"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
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

func (s *collection) Create(gameID string, req CreateRequest) (*entitiesCollection.Collection, error) {
	return s.repositoryCollection.Create(gameID, repositoriesCollection.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageFile,
	})
}
func (s *collection) Item(gameID, collectionID string) (*entitiesCollection.Collection, error) {
	return s.repositoryCollection.GetByID(gameID, collectionID)
}
func (s *collection) List(gameID, sortField, search string) ([]*entitiesCollection.Collection, error) {
	items, err := s.repositoryCollection.GetAll(gameID)
	if err != nil {
		return make([]*entitiesCollection.Collection, 0), err
	}

	// Filter
	var filteredItems []*entitiesCollection.Collection
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

	// Return empty array if no elements
	if filteredItems == nil {
		filteredItems = make([]*entitiesCollection.Collection, 0)
	}

	return filteredItems, nil
}
func (s *collection) Update(gameID, collectionID string, req UpdateRequest) (*entitiesCollection.Collection, error) {
	return s.repositoryCollection.Update(gameID, collectionID, repositoriesCollection.UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageFile,
	})
}
func (s *collection) Delete(gameID, collectionID string) error {
	return s.repositoryCollection.DeleteByID(gameID, collectionID)
}
func (s *collection) GetImage(gameID, collectionID string) ([]byte, string, error) {
	return s.repositoryCollection.GetImage(gameID, collectionID)
}
