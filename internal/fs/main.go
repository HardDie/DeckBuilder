package fs

import (
	"encoding/json"
	"os"

	"tts_deck_build/internal/errors"
)

func FileExist(path string) (isExist bool, e *errors.Error) {
	isExist, _, e = IsDir(path)
	return
}
func IsDir(path string) (isExist, isDir bool, e *errors.Error) {
	stat, err := os.Stat(path)
	if err == nil {
		isExist = true
	}

	if err != nil {
		if !os.IsNotExist(err) {
			// Some error
			errors.IfErrorLog(err)
			e = errors.InternalError.AddMessage(err.Error())
		}
		// File is not exist
		return
	}

	if stat.IsDir() {
		isDir = true
	}
	return
}
func RemoveDir(path string) (e *errors.Error) {
	err := os.RemoveAll(path)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
	}
	return
}
func CreateDir(path string) (e *errors.Error) {
	err := os.Mkdir(path, 0755)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
	}
	return
}
func WriteDataToFile(path string, data interface{}) (e *errors.Error) {
	f, err := os.Create(path)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	defer func() { errors.IfErrorLog(f.Close()) }()

	err = json.NewEncoder(f).Encode(data)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	return
}
func ReadDataFromFile[T any](path string) (data *T, e *errors.Error) {
	file, err := os.Open(path)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	defer func() { errors.IfErrorLog(file.Close()) }()

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	return
}
