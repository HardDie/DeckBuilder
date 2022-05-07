package collections

import (
	"os"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/games"
)

func FullCollectionCheck(gameName, collectionName string) (e error) {
	// Check game exist
	_, e = games.NewService().Item(gameName)
	if e != nil {
		return
	}

	// Check collection exist
	exist, e := CollectionIsExist(gameName, collectionName)
	if e != nil {
		return
	}
	if !exist {
		e = errors.CollectionNotExists
		return
	}

	// Check collection info exist
	exist, e = CollectionIsInfoExist(gameName, collectionName)
	if e != nil {
		return
	}
	if !exist {
		e = errors.CollectionInfoNotExists
		return
	}
	return
}

// Collection
func CollectionIsExist(gameName, collectionName string) (isExist bool, e error) {
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
func CollectionCreate(gameName, collectionName string) (e error) {
	collectionDir := filepath.Join(config.GetConfig().Games(), gameName, collectionName)
	return fs.CreateDir(collectionDir)
}
func CollectionRename(gameName, oldName, newName string) (e error) {
	oldCollectionDir := filepath.Join(config.GetConfig().Games(), gameName, oldName)
	newCollectionDir := filepath.Join(config.GetConfig().Games(), gameName, newName)
	err := os.Rename(oldCollectionDir, newCollectionDir)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
	}
	return
}
func CollectionDelete(gameName, collectionName string) (e error) {
	gameDir := filepath.Join(config.GetConfig().Games(), gameName, collectionName)
	return fs.Remove(gameDir)
}

// Info
func CollectionIsInfoExist(gameName, collectionName string) (isExist bool, e error) {
	infoFile := filepath.Join(config.GetConfig().Games(), gameName, collectionName, config.GetConfig().InfoFilename)
	return fs.FileExist(infoFile)
}
func CollectionAddInfo(gameName, collectionName string, info CollectionInfo) (e error) {
	infoPath := filepath.Join(config.GetConfig().Games(), gameName, collectionName, config.GetConfig().InfoFilename)
	return fs.WriteDataToFile(infoPath, info)
}
func CollectionGetInfo(gameName, collectionName string) (result *CollectionInfo, e error) {
	infoFile := filepath.Join(config.GetConfig().Games(), gameName, collectionName, config.GetConfig().InfoFilename)
	return fs.ReadDataFromFile[CollectionInfo](infoFile)
}
