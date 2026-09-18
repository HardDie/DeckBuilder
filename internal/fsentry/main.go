package fsentry

import (
	"path/filepath"
	"sync"

	"github.com/HardDie/DeckBuilder/internal/fsentry/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/fsentry/internal/fsutils"
	"github.com/HardDie/DeckBuilder/internal/fsentry/internal/utils"
	"github.com/HardDie/DeckBuilder/internal/fsentry/pkg/fsentry_error"
	"github.com/HardDie/DeckBuilder/internal/fsentry/pkg/fsentry_types"
)

type IFSEntry interface {
	Init() error
	Drop() error
	List(path ...string) (*entity.List, error)

	CreateFolder(name string, data interface{}, path ...string) (*entity.FolderInfo, error)
	GetFolder(name string, path ...string) (*entity.FolderInfo, error)
	MoveFolder(oldName, newName string, path ...string) (*entity.FolderInfo, error)
	UpdateFolder(name string, data interface{}, path ...string) (*entity.FolderInfo, error)
	RemoveFolder(name string, path ...string) error
	DuplicateFolder(srcName, dstName string, path ...string) (*entity.FolderInfo, error)
	UpdateFolderNameWithoutTimestamp(name, newName string, path ...string) error

	CreateEntry(name string, data interface{}, path ...string) error
	GetEntry(name string, path ...string) (*entity.Entry, error)
	MoveEntry(oldName, newName string, path ...string) error
	UpdateEntry(name string, data interface{}, path ...string) error
	RemoveEntry(name string, path ...string) error
	DuplicateEntry(srcName, dstName string, path ...string) error

	CreateBinary(name string, data []byte, path ...string) error
	GetBinary(name string, path ...string) ([]byte, error)
	MoveBinary(oldName, newName string, path ...string) error
	UpdateBinary(name string, data []byte, path ...string) error
	RemoveBinary(name string, path ...string) error
}
type FSEntry struct {
	root string
	rwm  sync.RWMutex

	isPretty bool
}

var (
	// validate interface
	_ IFSEntry = &FSEntry{}
)

func WithPretty() func(fs *FSEntry) {
	return func(fs *FSEntry) {
		fs.isPretty = true
	}
}

func NewFSEntry(root string, ops ...func(fs *FSEntry)) IFSEntry {
	res := &FSEntry{
		root: root,
	}
	for _, op := range ops {
		op(res)
	}
	return res
}

// Basic

func (db *FSEntry) Init() error {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	// Check if db folder exist
	isExist, err := fsutils.IsFolderExist(db.root)
	if err != nil {
		return err
	}
	if isExist {
		return nil
	}
	err = fsutils.CreateAllFolder(db.root)
	if err != nil {
		return err
	}
	return nil
}
func (db *FSEntry) Drop() error {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	// Check if db folder exist
	isExist, err := fsutils.IsFolderExist(db.root)
	if err != nil {
		return err
	}
	if !isExist {
		return nil
	}

	// Remove db folder
	err = fsutils.RemoveFolder(db.root)
	if err != nil {
		return err
	}

	return nil
}
func (db *FSEntry) List(path ...string) (*entity.List, error) {
	db.rwm.RLock()
	defer db.rwm.RUnlock()

	fullPath := db.buildPath("", path...)
	return fsutils.List(fullPath)
}

// Folder

func (db *FSEntry) CreateFolder(name string, data interface{}, path ...string) (*entity.FolderInfo, error) {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	if utils.NameToID(name) == "" {
		return nil, fsentry_error.ErrorBadName
	}

	fullPath, err := db.isFolderNotExist(name, path...)
	if err != nil {
		return nil, err
	}

	// Create folder
	err = fsutils.CreateFolder(fullPath)
	if err != nil {
		return nil, err
	}

	// Create info file
	info := entity.NewFolderInfo(name, data, db.isPretty)
	err = fsutils.CreateInfo(fullPath, info, db.isPretty)
	if err != nil {
		return nil, err
	}

	return info, nil
}
func (db *FSEntry) GetFolder(name string, path ...string) (*entity.FolderInfo, error) {
	db.rwm.RLock()
	defer db.rwm.RUnlock()

	if utils.NameToID(name) == "" {
		return nil, fsentry_error.ErrorBadName
	}

	fullPath, err := db.isFolderExist(name, path...)
	if err != nil {
		return nil, err
	}

	// Get info from file
	info, err := fsutils.GetInfo(fullPath)
	if err != nil {
		return nil, err
	}

	return info, nil
}
func (db *FSEntry) MoveFolder(oldName, newName string, path ...string) (*entity.FolderInfo, error) {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	if utils.NameToID(oldName) == "" || utils.NameToID(newName) == "" {
		return nil, fsentry_error.ErrorBadName
	}

	// Check if source folder exist
	fullOldPath, err := db.isFolderExist(oldName, path...)
	if err != nil {
		return nil, err
	}

	// Get info from file
	info, err := fsutils.GetInfo(fullOldPath)
	if err != nil {
		return nil, err
	}

	info.SetName(newName).UpdatedNow()

	// Update info file
	err = fsutils.CreateInfo(fullOldPath, info, db.isPretty)
	if err != nil {
		return nil, err
	}

	// If folders have same ID
	if utils.NameToID(oldName) == utils.NameToID(newName) {
		return info, nil
	}

	// Check if destination folder not exist
	fullNewPath, err := db.isFolderNotExist(newName, path...)
	if err != nil {
		return nil, err
	}

	// Rename folder
	err = fsutils.MoveObject(fullOldPath, fullNewPath)
	if err != nil {
		return nil, err
	}

	return info, nil
}
func (db *FSEntry) UpdateFolder(name string, data interface{}, path ...string) (*entity.FolderInfo, error) {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	if utils.NameToID(name) == "" {
		return nil, fsentry_error.ErrorBadName
	}

	fullPath, err := db.isFolderExist(name, path...)
	if err != nil {
		return nil, err
	}

	// Get info from file
	info, err := fsutils.GetInfo(fullPath)
	if err != nil {
		return nil, err
	}

	err = info.UpdateData(data, db.isPretty)
	if err != nil {
		return nil, err
	}

	// Update info file
	err = fsutils.CreateInfo(fullPath, info, db.isPretty)
	if err != nil {
		return nil, err
	}

	return info, nil
}
func (db *FSEntry) RemoveFolder(name string, path ...string) error {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	if utils.NameToID(name) == "" {
		return fsentry_error.ErrorBadName
	}

	fullPath, err := db.isFolderExist(name, path...)
	if err != nil {
		return err
	}

	err = fsutils.RemoveFolder(fullPath)
	if err != nil {
		return err
	}

	return nil
}
func (db *FSEntry) DuplicateFolder(srcName, dstName string, path ...string) (*entity.FolderInfo, error) {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	if utils.NameToID(srcName) == "" || utils.NameToID(dstName) == "" {
		return nil, fsentry_error.ErrorBadName
	}

	// Check if source folder exist
	fullSrcPath, err := db.isFolderExist(srcName, path...)
	if err != nil {
		return nil, err
	}

	// Check if destination folder not exist
	fullDstPath, err := db.isFolderNotExist(dstName, path...)
	if err != nil {
		return nil, err
	}

	// Copy folder
	err = fsutils.CopyFolder(fullSrcPath, fullDstPath)
	if err != nil {
		return nil, err
	}

	// Get info from file
	info, err := fsutils.GetInfo(fullDstPath)
	if err != nil {
		return nil, err
	}

	info.SetName(dstName).FlushTime()

	// Update info file
	err = fsutils.CreateInfo(fullDstPath, info, db.isPretty)
	if err != nil {
		return nil, err
	}

	return info, nil
}
func (db *FSEntry) UpdateFolderNameWithoutTimestamp(name, newName string, path ...string) error {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	if utils.NameToID(name) == "" || utils.NameToID(newName) == "" {
		return fsentry_error.ErrorBadName
	}

	fullPath, err := db.isFolderExist(name, path...)
	if err != nil {
		return err
	}

	// Get info from file
	info, err := fsutils.GetInfo(fullPath)
	if err != nil {
		return err
	}

	info.Id = utils.NameToID(newName)
	info.Name = fsentry_types.QuotedString(newName)

	// Update info file
	err = fsutils.CreateInfo(fullPath, info, db.isPretty)
	if err != nil {
		return err
	}

	return nil
}

// Entry

func (db *FSEntry) CreateEntry(name string, data interface{}, path ...string) error {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	if utils.NameToID(name) == "" {
		return fsentry_error.ErrorBadName
	}

	fullPath, err := db.isEntryNotExist(name, path...)
	if err != nil {
		return err
	}

	entry := entity.NewEntry(name, data, db.isPretty)
	err = fsutils.CreateEntry(fullPath, entry, db.isPretty)
	if err != nil {
		return err
	}

	return nil
}
func (db *FSEntry) GetEntry(name string, path ...string) (*entity.Entry, error) {
	db.rwm.RLock()
	defer db.rwm.RUnlock()

	if utils.NameToID(name) == "" {
		return nil, fsentry_error.ErrorBadName
	}

	fullPath, err := db.isEntryExist(name, path...)
	if err != nil {
		return nil, err
	}

	// Get info from file
	entry, err := fsutils.GetEntry(fullPath)
	if err != nil {
		return nil, err
	}

	return entry, nil
}
func (db *FSEntry) MoveEntry(oldName, newName string, path ...string) error {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	if utils.NameToID(oldName) == "" || utils.NameToID(newName) == "" {
		return fsentry_error.ErrorBadName
	}

	// Check if source entry exist
	fullOldPath, err := db.isEntryExist(oldName, path...)
	if err != nil {
		return err
	}

	// Read old entry
	entry, err := fsutils.GetEntry(fullOldPath)
	if err != nil {
		return err
	}

	entry.SetName(newName).UpdatedNow()

	var fullNewPath string
	// If entries have same ID
	if utils.NameToID(oldName) != utils.NameToID(newName) {
		// Check if destination entry not exist
		fullNewPath, err = db.isEntryNotExist(newName, path...)
		if err != nil {
			return err
		}
	} else {
		fullNewPath = fullOldPath
	}

	// Remove old entry
	err = fsutils.RemoveEntry(fullOldPath)
	if err != nil {
		return err
	}

	// Create new entry
	err = fsutils.CreateEntry(fullNewPath, entry, db.isPretty)
	if err != nil {
		return err
	}

	return nil
}
func (db *FSEntry) UpdateEntry(name string, data interface{}, path ...string) error {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	if utils.NameToID(name) == "" {
		return fsentry_error.ErrorBadName
	}

	fullPath, err := db.isEntryExist(name, path...)
	if err != nil {
		return err
	}

	// Get entry from file
	entry, err := fsutils.GetEntry(fullPath)
	if err != nil {
		return err
	}

	err = entry.UpdateData(data, db.isPretty)
	if err != nil {
		return err
	}

	// Update entry file
	err = fsutils.CreateEntry(fullPath, entry, db.isPretty)
	if err != nil {
		return err
	}

	return nil
}
func (db *FSEntry) RemoveEntry(name string, path ...string) error {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	if utils.NameToID(name) == "" {
		return fsentry_error.ErrorBadName
	}

	fullPath, err := db.isEntryExist(name, path...)
	if err != nil {
		return err
	}

	err = fsutils.RemoveEntry(fullPath)
	if err != nil {
		return err
	}

	return nil
}
func (db *FSEntry) DuplicateEntry(srcName, dstName string, path ...string) error {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	if utils.NameToID(srcName) == "" || utils.NameToID(dstName) == "" {
		return fsentry_error.ErrorBadName
	}

	// Check if source entry exist
	fullSrcPath, err := db.isEntryExist(srcName, path...)
	if err != nil {
		return err
	}

	// Check if destination entry not exist
	fullDstPath, err := db.isEntryNotExist(dstName, path...)
	if err != nil {
		return err
	}

	// Get entry from file
	entry, err := fsutils.GetEntry(fullSrcPath)
	if err != nil {
		return err
	}

	entry.SetName(dstName).FlushTime()

	// Create entry file
	err = fsutils.CreateEntry(fullDstPath, entry, db.isPretty)
	if err != nil {
		return err
	}

	return nil
}

// Binary

func (db *FSEntry) CreateBinary(name string, data []byte, path ...string) error {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	if utils.NameToID(name) == "" {
		return fsentry_error.ErrorBadName
	}

	fullPath, err := db.isBinaryNotExist(name, path...)
	if err != nil {
		return err
	}

	err = fsutils.CreateBinary(fullPath, data)
	if err != nil {
		return err
	}

	return nil
}
func (db *FSEntry) GetBinary(name string, path ...string) ([]byte, error) {
	db.rwm.RLock()
	defer db.rwm.RUnlock()

	if utils.NameToID(name) == "" {
		return nil, fsentry_error.ErrorBadName
	}

	fullPath, err := db.isBinaryExist(name, path...)
	if err != nil {
		return nil, err
	}

	// Get data from file
	data, err := fsutils.GetBinary(fullPath)
	if err != nil {
		return nil, err
	}

	return data, nil
}
func (db *FSEntry) MoveBinary(oldName, newName string, path ...string) error {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	if utils.NameToID(oldName) == "" || utils.NameToID(newName) == "" {
		return fsentry_error.ErrorBadName
	}

	// Check if source binary exist
	fullOldPath, err := db.isBinaryExist(oldName, path...)
	if err != nil {
		return err
	}

	// Check if destination binary not exist
	fullNewPath, err := db.isBinaryNotExist(newName, path...)
	if err != nil {
		return err
	}

	// Rename binary
	err = fsutils.MoveObject(fullOldPath, fullNewPath)
	if err != nil {
		return err
	}

	return nil
}
func (db *FSEntry) UpdateBinary(name string, data []byte, path ...string) error {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	if utils.NameToID(name) == "" {
		return fsentry_error.ErrorBadName
	}

	fullPath, err := db.isBinaryExist(name, path...)
	if err != nil {
		return err
	}

	// Update binary file
	err = fsutils.CreateBinary(fullPath, data)
	if err != nil {
		return err
	}

	return nil
}
func (db *FSEntry) RemoveBinary(name string, path ...string) error {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	if utils.NameToID(name) == "" {
		return fsentry_error.ErrorBadName
	}

	fullPath, err := db.isBinaryExist(name, path...)
	if err != nil {
		return err
	}

	err = fsutils.RemoveBinary(fullPath)
	if err != nil {
		return err
	}

	return nil
}

// util

func (db *FSEntry) buildPath(id string, path ...string) string {
	pathSlice := append([]string{db.root}, path...)
	return filepath.Join(append(pathSlice, id)...)
}
func (db *FSEntry) isFolderExist(name string, path ...string) (string, error) {
	id := utils.NameToID(name)
	if id == "" {
		return "", fsentry_error.ErrorBadName
	}

	fullPath := db.buildPath(id, path...)

	// Check if root folder exist
	isExist, err := fsutils.IsFolderExist(filepath.Dir(fullPath))
	if err != nil {
		return "", err
	}
	if !isExist {
		return "", fsentry_error.ErrorBadPath
	}

	// Check if destination folder exist
	isExist, err = fsutils.IsFolderExist(fullPath)
	if err != nil {
		return "", err
	}
	if !isExist {
		return "", fsentry_error.ErrorNotExist
	}

	return fullPath, nil
}
func (db *FSEntry) isFolderNotExist(name string, path ...string) (string, error) {
	id := utils.NameToID(name)
	if id == "" {
		return "", fsentry_error.ErrorBadName
	}

	fullPath := db.buildPath(id, path...)

	// Check if root folder exist
	isExist, err := fsutils.IsFolderExist(filepath.Dir(fullPath))
	if err != nil {
		return "", err
	}
	if !isExist {
		return "", fsentry_error.ErrorBadPath
	}

	// Check if destination folder exist
	isExist, err = fsutils.IsFolderExist(fullPath)
	if err != nil {
		return "", err
	}
	if isExist {
		return "", fsentry_error.ErrorExist
	}

	return fullPath, nil
}
func (db *FSEntry) isEntryExist(name string, path ...string) (string, error) {
	return db.isFileExist(name, ".json", path...)
}
func (db *FSEntry) isEntryNotExist(name string, path ...string) (string, error) {
	return db.isFileNotExist(name, ".json", path...)
}
func (db *FSEntry) isBinaryExist(name string, path ...string) (string, error) {
	return db.isFileExist(name, ".bin", path...)
}
func (db *FSEntry) isBinaryNotExist(name string, path ...string) (string, error) {
	return db.isFileNotExist(name, ".bin", path...)
}

func (db *FSEntry) isFileExist(name, ext string, path ...string) (string, error) {
	id := utils.NameToID(name)
	if id == "" {
		return "", fsentry_error.ErrorBadName
	}
	id += ext

	fullPath := db.buildPath(id, path...)

	// Check if root folder exist
	isExist, err := fsutils.IsFolderExist(filepath.Dir(fullPath))
	if err != nil {
		return "", err
	}
	if !isExist {
		return "", fsentry_error.ErrorBadPath
	}

	// Check if destination entry exist
	isExist, err = fsutils.IsFileExist(fullPath)
	if err != nil {
		return "", err
	}
	if !isExist {
		return "", fsentry_error.ErrorNotExist
	}

	return fullPath, nil
}
func (db *FSEntry) isFileNotExist(name, ext string, path ...string) (string, error) {
	id := utils.NameToID(name)
	if id == "" {
		return "", fsentry_error.ErrorBadName
	}
	id += ext

	fullPath := db.buildPath(id, path...)

	// Check if root folder exist
	isExist, err := fsutils.IsFolderExist(filepath.Dir(fullPath))
	if err != nil {
		return "", err
	}
	if !isExist {
		return "", fsentry_error.ErrorBadPath
	}

	// Check if destination entry exist
	isExist, err = fsutils.IsFileExist(fullPath)
	if err != nil {
		return "", err
	}
	if isExist {
		return "", fsentry_error.ErrorExist
	}

	return fullPath, nil
}
