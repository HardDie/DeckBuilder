package collections

import (
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
)

// Collection
func CollectionIsExist(gameName, collectionName string) (isExist bool, e *errors.Error) {
	collectionDir := filepath.Join(config.GetConfig().Games(), gameName, collectionName)

	// Check collection
	isExist, isDir, e := fs.IsDir(collectionDir)
	if e != nil {
		return
	}

	// Collection folder not exist
	if !isExist {
		return
	}

	// Is not folder
	if !isDir {
		e = errors.CollectionInvalid
		return
	}

	// Collection exist
	return
}
func CollectionCreate(gameName, collectionName string) (e *errors.Error) {
	collectionDir := filepath.Join(config.GetConfig().Games(), gameName, collectionName)
	return fs.CreateDir(collectionDir)
}

// Info
func CollectionIsInfoExist(gameName, collectionName string) (isExist bool, e *errors.Error) {
	infoFile := filepath.Join(config.GetConfig().Games(), gameName, collectionName, config.GetConfig().InfoFilename)
	return fs.FileExist(infoFile)
}
func CollectionAddInfo(gameName, collectionName string, info CollectionInfo) (e *errors.Error) {
	infoPath := filepath.Join(config.GetConfig().Games(), gameName, collectionName, config.GetConfig().InfoFilename)
	return fs.WriteDataToFile(infoPath, info)
}
func CollectionGetInfo(gameName, collectionName string) (result *CollectionInfo, e *errors.Error) {
	infoFile := filepath.Join(config.GetConfig().Games(), gameName, collectionName, config.GetConfig().InfoFilename)
	return fs.ReadDataFromFile[CollectionInfo](infoFile)
}
