package fs

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"tts_deck_build/internal/errors"
)

const (
	DirPerm = 0755
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

	// check if it is a folder
	if !stat.IsDir() {
		err = errors.InternalError.AddMessage("there should be a folder, but it's file")
		return false, err
	}

	// folder exists
	return true, nil
}
func IsFileExist(path string) (isExist bool, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// file not exist
			return false, nil
		}

		// other error
		errors.IfErrorLog(err)
		err = errors.InternalError.AddMessage(err.Error())
		return false, err
	}

	// check if it is a file
	if stat.IsDir() {
		err = errors.InternalError.AddMessage("there should be a file, but it's folder")
		return false, err
	}

	// file exists
	return true, nil
}

func CreateFolder(path string) error {
	err := os.Mkdir(path, DirPerm)
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage(err.Error())
	}
	return nil
}
func CreateFolderIfNotExist(path string) error {
	isExists, err := IsFolderExist(path)
	if err != nil || isExists {
		return err
	}
	err = CreateFolder(path)
	if err != nil {
		return err
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
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.InternalError.AddMessage(err.Error())
	}

	folders := make([]string, 0)
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

	data, err := io.ReadAll(file)
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
func ListOfFiles(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.InternalError.AddMessage(err.Error())
	}

	listFiles := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if []rune(file.Name())[0] == '.' {
			// Skip hidden files
			continue
		}
		if filepath.Ext(file.Name()) != ".json" {
			// Skip non json files
			continue
		}
		listFiles = append(listFiles, file.Name())
	}

	return listFiles, nil
}

func GetFilenameWithoutExt(name string) string {
	return name[:len(name)-len(filepath.Ext(name))]
}
