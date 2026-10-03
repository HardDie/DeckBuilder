package deck

import (
	"encoding/json"
	"slices"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
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
			Exist:         apperr.ErrDeckExists,
			NotExist:      apperr.ErrDeckNotFound,
			ImageExist:    apperr.ErrDeckImageExists,
			ImageNotExist: apperr.ErrDeckImageNotFound,
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
	return r.toEntity(saved.Info, gameID, collectionID, saved.ImageError), nil
}

func (r *deck) GetByID(gameID, collectionID, deckID string) (*entitiesDeck.Deck, error) {
	info, err := r.folder.Get([]string{gameID, collectionID}, deckID)
	if err != nil {
		return nil, err
	}
	return r.toEntity(info, gameID, collectionID, nil), nil
}

func (r *deck) GetAll(gameID, collectionID string) ([]*entitiesDeck.Deck, error) {
	infos, err := r.folder.List([]string{gameID, collectionID})
	if err != nil {
		return nil, err
	}
	var decks []*entitiesDeck.Deck
	for _, info := range infos {
		decks = append(decks, r.toEntity(info, gameID, collectionID, nil))
	}
	return decks, nil
}

func (r *deck) Update(gameID, collectionID, deckID string, req UpdateRequest) (*entitiesDeck.Deck, error) {
	saved, err := r.folder.Update([]string{gameID, collectionID}, deckID, repositories.FolderWrite(req))
	if err != nil {
		return nil, err
	}
	return r.toEntity(saved.Info, gameID, collectionID, saved.ImageError), nil
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

func (r *deck) toEntity(info repositories.FolderInfo, gameID, collectionID string, imageErr error) *entitiesDeck.Deck {
	createdAt, updatedAt := utils.NormalizeTimestamps(info.CreatedAt, info.UpdatedAt)
	return &entitiesDeck.Deck{
		ID:                info.ID,
		Name:              info.Name,
		Description:       info.Data.Description.String(),
		Image:             info.Data.Image.String(),
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		GameID:            gameID,
		CollectionID:      collectionID,
		HasImage:          r.folder.HasImage([]string{gameID, collectionID}, info.ID),
		CardsMissingImage: r.cardsMissingImage(gameID, collectionID, info.ID),
		ImageError:        imageErr,
	}
}

// cardsMissingImage tells whether a card in the deck has no image file.
// The deck owns its "cards" folder: the card list is its data,
// and each card's image is a binary named by the card id. No image is read.
func (r *deck) cardsMissingImage(gameID, collectionID, deckID string) bool {
	cards, err := r.db.GetFolder[map[string]json.RawMessage]("cards", r.gamesPath, gameID, collectionID, deckID)
	if err != nil {
		return false
	}
	files, err := r.db.List(r.gamesPath, gameID, collectionID, deckID, "cards")
	if err != nil {
		return false
	}
	for id := range cards.Data {
		if !slices.Contains(files.Binaries, id) {
			return true
		}
	}
	return false
}
