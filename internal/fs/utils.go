package fs

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"tts_deck_build/internal/errors"
)

func IsFolderExist(path string) (isExist bool, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// folder not exist
			return false, nil
		}

		// other error
		errors.IfErrorLog(err)
		err = errors.InternalError.AddMessage(err.Error())
		return false, err
	}

	// check if it folder
	if !stat.IsDir() {
		err = errors.InternalError.AddMessage("there should be a folder, but it's file")
		return false, err
	}

	// folder exist
	return true, nil
}
func IsFileExist(path string) (isExist bool, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		errors.IfErrorLog(err)
		err = errors.InternalError.AddMessage(err.Error())
		return
	}

	if stat.IsDir() {
		err = errors.InternalError.AddMessage("there should be a file, but it's folder")
		return
	}

	isExist = true
	return
}

func CreateFolder(path string) error {
	err := os.Mkdir(path, 0755)
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage(err.Error())
	}
	return nil
}
func RemoveFolder(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage(err.Error())
	}
	return nil
}
func MoveFolder(oldPath, newPath string) error {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage(err.Error())
	}
	return nil
}
func ListOfFolders(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, errors.InternalError.AddMessage(err.Error())
	}

	var folders []string
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		folders = append(folders, file.Name())
	}

	return folders, nil
}

func WriteFile[T any](path string, data T) error {
	// Creating a file
	f, err := os.Create(path)
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage(err.Error())
	}
	defer func() { errors.IfErrorLog(f.Close()) }()

	// Marshalling data to json and writing to file
	enc := json.NewEncoder(f)
	enc.SetIndent("", "	")
	if err = enc.Encode(data); err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage(err.Error())
	}
	return nil
}
func WriteBinaryFile(path string, data []byte) error {
	// Creating a file
	f, err := os.Create(path)
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage(err.Error())
	}
	defer func() { errors.IfErrorLog(f.Close()) }()

	// Write data to file
	_, err = f.Write(data)
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage(err.Error())
	}
	return nil
}
func ReadFile[T any](path string) (data *T, err error) {
	file, err := os.Open(path)
	if err != nil {
		errors.IfErrorLog(err)
		return nil, errors.InternalError.AddMessage(err.Error())
	}
	defer func() { errors.IfErrorLog(file.Close()) }()

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		errors.IfErrorLog(err)
		return nil, errors.InternalError.AddMessage(err.Error())
	}
	return
}
func ReadBinaryFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		errors.IfErrorLog(err)
		return nil, errors.InternalError.AddMessage(err.Error())
	}
	defer func() { errors.IfErrorLog(file.Close()) }()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		errors.IfErrorLog(err)
		return nil, errors.InternalError.AddMessage(err.Error())
	}

	return data, nil
}
func RemoveFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage(err.Error())
	}
	return nil
}
