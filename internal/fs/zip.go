package fs

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
)

func ArchiveFolder(gamePath, gameID string) (data []byte, err error) {
	// Create buffer for result archive
	buf := bytes.Buffer{}

	// Create zip processor
	zipWriter := zip.NewWriter(&buf)
	defer func() {
		if zipWriter != nil {
			errors.IfErrorLog(zipWriter.Close())
		}
	}()

	// Add all filed to archive
	err = recursiveWalk(zipWriter, gamePath, gameID)
	if err != nil {
		return
	}

	// Flush data from zip writer into buffer
	err = zipWriter.Close()
	if err != nil {
		errors.IfErrorLog(err)
		return nil, errors.InternalError.AddMessage("Error closing zip archive: " + err.Error())
	}
	// Bypassing double closing
	zipWriter = nil

	// Return zip file
	return buf.Bytes(), nil
}

func UnarchiveFolder(data []byte, gameID string) (err error) {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage("Error open zip archive: " + err.Error())
	}

	var importGameID string

	// Basic validation
	for _, file := range zipReader.File {
		pathList := strings.Split(file.Name, string(filepath.Separator))
		if importGameID == "" {
			// Get the original title of the game
			importGameID = pathList[0]
			continue
		}
		if importGameID != pathList[0] {
			// All files should be located inside one folder
			return errors.BadArchive
		}
	}

	// Unzip files
	for _, file := range zipReader.File {
		zipFilePath := file.Name
		if gameID != "" {
			// If the user passed a different title of the game, replace the original title with the new one
			pathList := strings.Split(zipFilePath, string(filepath.Separator))
			pathList[0] = gameID
			zipFilePath = filepath.Join(pathList...)
		}

		// Building the path in the FS for the file from the archive
		resultPath := filepath.Join(config.GetConfig().Games(), zipFilePath)

		if file.FileInfo().IsDir() {
			err = CreateFolderIfNotExist(resultPath)
			if err != nil {
				return err
			}
		} else {
			// Create file in the new game directory
			err = createFileFromArchive(resultPath, file)
			if err != nil {
				return
			}
		}
	}
	return
}

func recursiveWalk(zipWriter *zip.Writer, dirPath string, relatePath string) error {
	// Open dir
	dirFiles, err := os.ReadDir(dirPath)
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage("Error open folder: " + dirPath + "; " + err.Error())
	}

	// Go through all files
	for _, file := range dirFiles {
		filePath := filepath.Join(dirPath, file.Name())
		newRelatePath := filepath.Join(relatePath, file.Name())

		if file.IsDir() {
			err = recursiveWalk(zipWriter, filePath, newRelatePath)
			if err != nil {
				return err
			}
		} else {
			// Copy file inside archive
			err = addFileIntoArchive(filePath, newRelatePath, zipWriter)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Adding single file into archive
func addFileIntoArchive(filePath, zipPath string, w *zip.Writer) error {
	// Open file for reading
	f, err := os.Open(filePath)
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage("Error open file: " + err.Error())
	}
	defer func() {
		errors.IfErrorLog(f.Close())
	}()

	// Create path inside zip archive
	zipW, err := w.Create(zipPath)
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage("Error create path in zip: " + err.Error())
	}

	// Copy file inside archive
	_, err = io.Copy(zipW, f)
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage("Error copy file insize zip: " + err.Error())
	}
	return nil
}

func createFileFromArchive(filePath string, f *zip.File) error {
	folder := filepath.Dir(filePath)
	err := CreateFolderIfNotExist(folder)
	if err != nil {
		return err
	}

	// Open file in zip archive
	fileInArchive, err := f.Open()
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage("Error open file in zip archive: " + err.Error())
	}
	defer func() { errors.IfErrorLog(fileInArchive.Close()) }()

	// Create new file
	createdFile, err := os.Create(filePath)
	if err != nil {
		return errors.InternalError.AddMessage("Error creating new file: " + err.Error())
	}

	// Write data into created file
	_, err = io.Copy(createdFile, fileInArchive)
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage("Error copy data from zip to file: " + err.Error())
	}

	return nil
}
