package fs

import (
	"encoding/json"
	"os"

	"tts_deck_build/internal/errors"
)

const (
	ImagePngType  = "image/png"
	ImageJpegType = "image/jpeg"
	ImageGifType  = "image/gif"
)

func FileExist(path string) (isExist bool, e error) {
	isExist, _, e = IsDir(path)
	return
}
func IsDir(path string) (isExist, isDir bool, e error) {
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
func Remove(path string) (e error) {
	err := os.RemoveAll(path)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
	}
	return
}
func CreateDir(path string) (e error) {
	err := os.Mkdir(path, 0755)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
	}
	return
}
func WriteDataToFile(path string, data interface{}) (e error) {
	f, err := os.Create(path)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	defer func() { errors.IfErrorLog(f.Close()) }()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "	")
	if err = enc.Encode(data); err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
		return
	}
	return
}
func ReadDataFromFile[T any](path string) (data *T, e error) {
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
