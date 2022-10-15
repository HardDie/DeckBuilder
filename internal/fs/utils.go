package fs

import (
	"io"
	"os"
	"path/filepath"

	"github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/logger"
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
	err := os.MkdirAll(path, DirPerm)
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

func CreateAndProcess[T any](path string, in T, cb func(w io.Writer, in T) error) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { errors.IfErrorLog(file.Close()) }()

	return cb(file, in)
}
func OpenAndProcess[T any](path string, cb func(r io.Reader) (T, error)) (res T, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() { errors.IfErrorLog(file.Close()) }()

	return cb(file)
}

func PathToAbsolutePath(path string) string {
	res, err := filepath.Abs(path)
	if err != nil {
		logger.Error.Printf("Can't transform path %q to absolute path. %q", path, err.Error())
		return path
	}
	return res
}
