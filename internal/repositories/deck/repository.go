package deck

import (
	"github.com/HardDie/fsentry"

	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/repositories"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type deck struct {
	db        *fsentry.DB
	gamesPath string
	folder    *repositories.Folder
}

func New(db *fsentry.DB) Deck {
	return &deck{
		db:        db,
		gamesPath: "games",
		folder: repositories.NewFolder(db, repositories.FolderErrors{
			Exist:         er.DeckExist,
			NotExist:      er.DeckNotExists,
			ImageExist:    er.DeckImageExist,
			ImageNotExist: er.DeckImageNotExists,
		}),
	}
}

func (r *deck) Create(gameID, collectionID string, req CreateRequest) (*entitiesDeck.Deck, error) {
	saved, err := r.folder.Create([]string{gameID, collectionID}, repositories.FolderWrite(req))
	if err != nil {
		return nil, err
	}
	// The cards of a deck live in one "cards" folder under it.
	_, err = r.db.CreateFolder[any]("cards", nil, r.gamesPath, gameID, collectionID, saved.Info.ID)
	if err != nil {
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{})
	}
	return toEntity(saved.Info, gameID, collectionID, saved.ImageError), nil
}

func (r *deck) GetByID(gameID, collectionID, deckID string) (*entitiesDeck.Deck, error) {
	info, err := r.folder.Get([]string{gameID, collectionID}, deckID)
	if err != nil {
		return nil, err
	}
	return toEntity(info, gameID, collectionID, nil), nil
}

func (r *deck) GetAll(gameID, collectionID string) ([]*entitiesDeck.Deck, error) {
	infos, err := r.folder.List([]string{gameID, collectionID})
	if err != nil {
		return nil, err
	}
	var decks []*entitiesDeck.Deck
	for _, info := range infos {
		decks = append(decks, toEntity(info, gameID, collectionID, nil))
	}
	return decks, nil
}

func (r *deck) Update(gameID, collectionID, deckID string, req UpdateRequest) (*entitiesDeck.Deck, error) {
	saved, err := r.folder.Update([]string{gameID, collectionID}, deckID, repositories.FolderWrite(req))
	if err != nil {
		return nil, err
	}
	return toEntity(saved.Info, gameID, collectionID, saved.ImageError), nil
}

func (r *deck) DeleteByID(gameID, collectionID, deckID string) error {
	return r.folder.Delete([]string{gameID, collectionID}, deckID)
}

func (r *deck) GetImage(gameID, collectionID, deckID string) ([]byte, string, error) {
	return r.folder.Image([]string{gameID, collectionID}, deckID)
}

// GetAllDecksInGame lists the decks of every collection in the game.
// Decks with the same name and image URL are listed once.
func (r *deck) GetAllDecksInGame(gameID string) ([]*entitiesDeck.Deck, error) {
	list, err := r.db.List(r.gamesPath, gameID)
	if err != nil {
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{})
	}

	type key struct{ name, image string }
	seen := make(map[key]struct{})
	decks := make([]*entitiesDeck.Deck, 0)
	for _, collectionID := range list.Folders {
		collectionDecks, err := r.GetAll(gameID, collectionID)
		if err != nil {
			return nil, err
		}
		for _, d := range collectionDecks {
			k := key{d.Name, d.Image}
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			decks = append(decks, d)
		}
	}
	return decks, nil
}

func toEntity(info repositories.FolderInfo, gameID, collectionID string, imageErr error) *entitiesDeck.Deck {
	createdAt, updatedAt := utils.NormalizeTimestamps(info.CreatedAt, info.UpdatedAt)
	return &entitiesDeck.Deck{
		ID:           info.ID,
		Name:         info.Name,
		Description:  info.Data.Description.String(),
		Image:        info.Data.Image.String(),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		GameID:       gameID,
		CollectionID: collectionID,
		ImageError:   imageErr,
	}
}
