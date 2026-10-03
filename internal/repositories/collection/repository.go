package collection

import (
	"github.com/HardDie/fsentry"

	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/repositories"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type collection struct {
	folder *repositories.Folder
}

func New(db *fsentry.DB) Collection {
	return &collection{
		folder: repositories.NewFolder(db, repositories.FolderErrors{
			Exist:         er.CollectionExist,
			NotExist:      er.CollectionNotExists,
			ImageExist:    er.CollectionImageExist,
			ImageNotExist: er.CollectionImageNotExists,
		}),
	}
}

func (r *collection) Create(gameID string, req CreateRequest) (*entitiesCollection.Collection, error) {
	saved, err := r.folder.Create([]string{gameID}, repositories.FolderWrite(req))
	if err != nil {
		return nil, err
	}
	return toEntity(saved.Info, gameID, saved.ImageError), nil
}

func (r *collection) GetByID(gameID, collectionID string) (*entitiesCollection.Collection, error) {
	info, err := r.folder.Get([]string{gameID}, collectionID)
	if err != nil {
		return nil, err
	}
	return toEntity(info, gameID, nil), nil
}

func (r *collection) GetAll(gameID string) ([]*entitiesCollection.Collection, error) {
	infos, err := r.folder.List([]string{gameID})
	if err != nil {
		return nil, err
	}
	var collections []*entitiesCollection.Collection
	for _, info := range infos {
		collections = append(collections, toEntity(info, gameID, nil))
	}
	return collections, nil
}

func (r *collection) Update(gameID, collectionID string, req UpdateRequest) (*entitiesCollection.Collection, error) {
	saved, err := r.folder.Update([]string{gameID}, collectionID, repositories.FolderWrite(req))
	if err != nil {
		return nil, err
	}
	return toEntity(saved.Info, gameID, saved.ImageError), nil
}

func (r *collection) DeleteByID(gameID, collectionID string) error {
	return r.folder.Delete([]string{gameID}, collectionID)
}

func (r *collection) GetImage(gameID, collectionID string) ([]byte, string, error) {
	return r.folder.Image([]string{gameID}, collectionID)
}

func toEntity(info repositories.FolderInfo, gameID string, imageErr error) *entitiesCollection.Collection {
	createdAt, updatedAt := utils.NormalizeTimestamps(info.CreatedAt, info.UpdatedAt)
	return &entitiesCollection.Collection{
		ID:          info.ID,
		Name:        info.Name,
		Description: info.Data.Description.String(),
		Image:       info.Data.Image.String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		GameID:      gameID,
		ImageError:  imageErr,
	}
}
