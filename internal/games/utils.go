package games

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/utils"
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
	_, err := os.Stat(gameDir)
	if err != nil && !os.IsNotExist(err) {
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	isExist = true
	return
}
func GameCreate(name string) (e *errors.Error) {
	gameDir := filepath.Join(config.GetConfig().Games(), name)
	err := os.Mkdir(gameDir, 0755)
	if err != nil {
		utils.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	return
}
func GameRename(oldName, newName string) (e *errors.Error) {
	oldGameDir := filepath.Join(config.GetConfig().Games(), oldName)
	newGameDir := filepath.Join(config.GetConfig().Games(), newName)
	err := os.Rename(oldGameDir, newGameDir)
	if err != nil {
		utils.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
	}
	return
}
func GameDelete(name string) (e *errors.Error) {
	gameDir := filepath.Join(config.GetConfig().Games(), name)
	err := os.RemoveAll(gameDir)
	if err != nil {
		utils.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	return
}

// Info
func GameIsInfoExist(name string) (isExist bool, e *errors.Error) {
	infoFile := filepath.Join(config.GetConfig().Games(), name, config.GetConfig().InfoFilename)
	_, err := os.Stat(infoFile)
	if err != nil && !os.IsNotExist(err) {
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	isExist = true
	return
}
func GameAddInfo(name string, info GameInfo) (e *errors.Error) {
	data, err := json.Marshal(info)
	if err != nil {
		utils.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	err = ioutil.WriteFile(filepath.Join(config.GetConfig().Games(), name, "info.json"), data, 0644)
	if err != nil {
		utils.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	return
}
func GameGetInfo(name string) (result *GameInfo, e *errors.Error) {
	infoFile := filepath.Join(config.GetConfig().Games(), name, config.GetConfig().InfoFilename)
	file, err := os.Open(infoFile)
	if err != nil {
		utils.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	defer func() { utils.IfErrorLog(file.Close()) }()

	err = json.NewDecoder(file).Decode(&result)
	if err != nil {
		utils.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	return
}
