package repositories

import (
	"errors"
	"slices"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
)

const (
	// gamesPath is the catalog root inside the fsentry store.
	gamesPath = "games"
	// imageName is the binary that holds an entity image.
	imageName = "image"
)

// FolderModel is the sidecar data of a game, collection, or deck folder.
type FolderModel struct {
	Description fsentry.QuotedString `json:"description"`
	Image       fsentry.QuotedString `json:"image"`
}

// FolderInfo is one stored game, collection, or deck.
type FolderInfo = fsentry.FolderInfo[FolderModel]

// FolderErrors are the errors of one catalog level.
type FolderErrors struct {
	Exist         *apperr.Error
	NotExist      *apperr.Error
	ImageExist    *apperr.Error
	ImageNotExist *apperr.Error
}

// FolderWrite is a create or update of one folder.
type FolderWrite struct {
	Name        string
	Description string
	Image       string
	ImageFile   []byte
}

// Saved is the result of a create or update.
// ImageError says why a new image was not applied; the save itself succeeded.
type Saved struct {
	Info       FolderInfo
	ImageError error
}

// Folder stores one catalog level: folders under a parent, each with an optional image.
// The parent is the path below "games": nil for a game,
// {gameID} for a collection, {gameID, collectionID} for a deck.
type Folder struct {
	db   *fsentry.DB
	errs FolderErrors
}

func NewFolder(db *fsentry.DB, errs FolderErrors) *Folder {
	return &Folder{db: db, errs: errs}
}

// Create makes a folder and stores its image.
// The image is downloaded and validated first; a bad one is not applied.
func (f *Folder) Create(parent []string, req FolderWrite) (Saved, error) {
	change, imageErr := ResolveImage("", req.Image, req.ImageFile)

	info, err := f.db.CreateFolder(req.Name, FolderModel{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(change.URL),
	}, f.path(parent)...)
	if err != nil {
		return Saved{}, MapFsentry(err, FsentrySentinels{Exist: f.errs.Exist})
	}
	if change.Data != nil {
		if err = f.imageCreate(parent, info.ID, change.Data); err != nil {
			imageErr = err
		}
	}
	return Saved{Info: info, ImageError: imageErr}, nil
}

// Get reads one folder.
func (f *Folder) Get(parent []string, id string) (FolderInfo, error) {
	info, err := f.db.GetFolder[FolderModel](id, f.path(parent)...)
	if err != nil {
		return FolderInfo{}, MapFsentry(err, FsentrySentinels{NotExist: f.errs.NotExist})
	}
	return info, nil
}

// List reads every folder under parent.
// A folder that cannot be read, or whose id does not match its name, is logged and skipped.
func (f *Folder) List(parent []string) ([]FolderInfo, error) {
	list, err := f.db.List(f.path(parent)...)
	if err != nil {
		return nil, MapFsentry(err, FsentrySentinels{})
	}

	var infos []FolderInfo
	for _, folder := range list.Folders {
		info, err := f.Get(parent, folder)
		if err != nil {
			logger.Error.Println(folder, err.Error())
			continue
		}
		if folder != info.ID {
			logger.Error.Println("Corrupted folder:", folder)
			continue
		}
		infos = append(infos, info)
	}
	return infos, nil
}

// Update renames the folder, saves its data, and swaps its image.
// The new image is downloaded and validated before any write;
// a bad one keeps the old image and URL.
func (f *Folder) Update(parent []string, id string, req FolderWrite) (Saved, error) {
	old, err := f.Get(parent, id)
	if err != nil {
		return Saved{}, err
	}
	change, imageErr := ResolveImage(old.Data.Image.String(), req.Image, req.ImageFile)

	info := old
	if old.Name != req.Name {
		info, err = f.move(parent, old.Name, req.Name)
		if err != nil {
			return Saved{}, err
		}
	}

	if old.Data.Description.String() != req.Description ||
		old.Data.Image.String() != change.URL ||
		change.Data != nil {
		info, err = f.update(parent, req.Name, req.Description, change.URL)
		if err != nil {
			return Saved{}, err
		}
	}

	if change.Data != nil || change.Clear {
		err = f.imageDelete(parent, info.ID)
		if err != nil && !errors.Is(err, f.errs.ImageNotExist) {
			return Saved{}, err
		}
	}
	if change.Data != nil {
		if err = f.imageCreate(parent, info.ID, change.Data); err != nil {
			imageErr = err
		}
	}
	return Saved{Info: info, ImageError: imageErr}, nil
}

// Delete removes the folder and everything under it.
func (f *Folder) Delete(parent []string, id string) error {
	err := f.db.RemoveFolder(id, f.path(parent)...)
	if err != nil {
		return MapFsentry(err, FsentrySentinels{NotExist: f.errs.NotExist})
	}
	return nil
}

// HasImage tells whether the folder holds an image file. It lists the folder; no image is read.
func (f *Folder) HasImage(parent []string, id string) bool {
	list, err := f.db.List(f.path(parent, id)...)
	if err != nil {
		return false
	}
	return slices.Contains(list.Binaries, imageName)
}

// Image returns the image bytes and their type ("png", "jpeg", "gif").
// Only the header is read to learn the type.
func (f *Folder) Image(parent []string, id string) ([]byte, string, error) {
	data, err := f.imageGet(parent, id)
	if err != nil {
		return nil, "", err
	}
	imgType, err := images.ImageType(data)
	if err != nil {
		return nil, "", err
	}
	return data, imgType, nil
}

func (f *Folder) move(parent []string, oldName, newName string) (FolderInfo, error) {
	info, err := f.db.MoveFolder[FolderModel](oldName, newName, f.path(parent)...)
	if err != nil {
		return FolderInfo{}, MapFsentry(err, FsentrySentinels{NotExist: f.errs.NotExist})
	}
	return info, nil
}

func (f *Folder) update(parent []string, name, description, image string) (FolderInfo, error) {
	info, err := f.db.UpdateFolder(name, FolderModel{
		Description: fsentry.QuotedString(description),
		Image:       fsentry.QuotedString(image),
	}, f.path(parent)...)
	if err != nil {
		return FolderInfo{}, MapFsentry(err, FsentrySentinels{NotExist: f.errs.NotExist})
	}
	return info, nil
}

func (f *Folder) imageCreate(parent []string, id string, data []byte) error {
	info, err := f.Get(parent, id)
	if err != nil {
		return err
	}
	err = f.db.CreateBinary(imageName, data, f.path(parent, info.ID)...)
	if err != nil {
		return MapFsentry(err, FsentrySentinels{Exist: f.errs.ImageExist})
	}
	return nil
}

func (f *Folder) imageGet(parent []string, id string) ([]byte, error) {
	info, err := f.Get(parent, id)
	if err != nil {
		return nil, err
	}
	data, err := f.db.GetBinary(imageName, nil, f.path(parent, info.ID)...)
	if err != nil {
		return nil, MapFsentry(err, FsentrySentinels{NotExist: f.errs.ImageNotExist})
	}
	return data, nil
}

func (f *Folder) imageDelete(parent []string, id string) error {
	info, err := f.Get(parent, id)
	if err != nil {
		return err
	}
	err = f.db.RemoveBinary(imageName, f.path(parent, info.ID)...)
	if err != nil {
		return MapFsentry(err, FsentrySentinels{NotExist: f.errs.ImageNotExist})
	}
	return nil
}

// path is "games", then parent, then rest.
func (f *Folder) path(parent []string, rest ...string) []string {
	path := make([]string, 0, 1+len(parent)+len(rest))
	path = append(path, gamesPath)
	path = append(path, parent...)
	return append(path, rest...)
}
