package games

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
)

func FullGameCheck(name string) (e *errors.Error) {
	// Check game exist
	exist, e := GameIsExist(name)
	if e != nil {
		return
	}
	if !exist {
		e = errors.GameNotExists
		return
	}

	// Check game info exist
	exist, e = GameIsInfoExist(name)
	if e != nil {
		return
	}
	if !exist {
		e = errors.GameInfoNotExists
		return
	}
	return
}

// Game
func GameIsExist(name string) (isExist bool, e *errors.Error) {
	gameDir := filepath.Join(config.GetConfig().Games(), name)

	// Check game
	isExist, isDir, e := fs.IsDir(gameDir)
	if e != nil {
		return
	}

	// Game folder not exist
	if !isExist {
		return
	}

	// Is not folder
	if !isDir {
		e = errors.GameInvalid
		return
	}

	// Game exist
	return
}
func GameCreate(name string) (e *errors.Error) {
	gameDir := filepath.Join(config.GetConfig().Games(), name)
	return fs.CreateDir(gameDir)
}
func GameRename(oldName, newName string) (e *errors.Error) {
	oldGameDir := filepath.Join(config.GetConfig().Games(), oldName)
	newGameDir := filepath.Join(config.GetConfig().Games(), newName)
	err := os.Rename(oldGameDir, newGameDir)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
	}
	return
}
func GameDelete(name string) (e *errors.Error) {
	gameDir := filepath.Join(config.GetConfig().Games(), name)
	if e = fs.RemoveDir(gameDir); e != nil {
		return
	}
	return
}

// Info
func GameIsInfoExist(name string) (isExist bool, e *errors.Error) {
	infoFile := filepath.Join(config.GetConfig().Games(), name, config.GetConfig().InfoFilename)
	return fs.FileExist(infoFile)
}
func GameAddInfo(name string, info GameInfo) (e *errors.Error) {
	data, err := json.Marshal(info)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	err = ioutil.WriteFile(filepath.Join(config.GetConfig().Games(), name, "info.json"), data, 0644)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	return
}
func GameGetInfo(name string) (result *GameInfo, e *errors.Error) {
	infoFile := filepath.Join(config.GetConfig().Games(), name, config.GetConfig().InfoFilename)
	file, err := os.Open(infoFile)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	defer func() { errors.IfErrorLog(file.Close()) }()

	err = json.NewDecoder(file).Decode(&result)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	return
}
