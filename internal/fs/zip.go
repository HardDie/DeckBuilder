package fs

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/utils"
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

func UnarchiveFolder(data []byte, gameID string, cfg *config.Config) (resultGameID string, err error) {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		errors.IfErrorLog(err)
		return "", errors.InternalError.AddMessage("Error open zip archive: " + err.Error())
	}

	// The original game ID as it is set in the zip-archive
	var importGameID string

	// Validation
	for _, file := range zipReader.File {
		file.Name = convertPathToPlatform(file.Name)
		// Split the full path to the file into parts
		pathList := strings.Split(file.Name, string(filepath.Separator))

		// The archive must not contain files outside the folder
		if !file.Mode().IsDir() && len(pathList) < 2 {
			return "", errors.BadArchive.AddMessage("The file is located outside the root folder")
		}

		// If this is the first file, extract the name of the root folder
		if importGameID == "" {
			// Get the original title of the game
			importGameID = pathList[0]

			// Check that the folder name matches the required format of the ID
			if utils.NameToID(importGameID) != importGameID {
				return "", errors.BadArchive.AddMessage("The root folder of the game has a bad ID")
			}
			continue
		}

		// The root folder for all files must be the same
		if importGameID != pathList[0] {
			return "", errors.BadArchive.AddMessage("There should only be one folder in the root of the archive")
		}
	}

	// Set the resulting game ID
	if gameID != "" {
		// If the user passed the game ID, the created folder will have the following name
		resultGameID = gameID
	} else {
		// If the user skips the game ID, the created folder will have the same name as before
		resultGameID = importGameID
	}

	// Build a full relative path to the root game folder
	gameRootPath := filepath.Join(cfg.Games(), resultGameID)
	// Create the root folder of the game
	err = CreateFolder(gameRootPath)
	if err != nil {
		return "", err
	}

	defer func() {
		if err != nil {
			// If an error occurs during unzipping, delete the created folder with the game
			errors.IfErrorLog(RemoveFolder(gameRootPath))
		}
	}()

	// Unzip files
	for _, file := range zipReader.File {
		file.Name = convertPathToPlatform(file.Name)

		// Backing up the original file path
		zipFilePath := file.Name

		if gameID != "" {
			// If the user passed a different name of the game, replace the name of the root folder
			pathList := strings.Split(zipFilePath, string(filepath.Separator))
			pathList[0] = gameID
			zipFilePath = filepath.Join(pathList...)
		}

		// Build a full relative path to the game folder
		resultPath := filepath.Join(cfg.Games(), zipFilePath)

		if file.FileInfo().IsDir() {
			err = CreateFolderIfNotExist(resultPath)
			if err != nil {
				return "", err
			}
		} else {
			// Create file in the new game directory
			err = createFileFromArchive(resultPath, file)
			if err != nil {
				return "", err
			}
		}
	}
	return resultGameID, nil
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

// Converting the file path inside a zip archive to platform specific.
// Because windows uses "\" and unix uses "/", we get errors when unzipping.
func convertPathToPlatform(path string) string {
	switch runtime.GOOS {
	case "linux", "darwin":
		return strings.ReplaceAll(path, "\\", "/")
	case "windows":
		return strings.ReplaceAll(path, "/", "\\")
	}
	return path
}
