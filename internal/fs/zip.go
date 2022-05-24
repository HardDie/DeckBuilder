package fs

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"

	"tts_deck_build/internal/errors"
)

func ArchiveFolder(gamePath, gameId string) (data []byte, err error) {
	// Create buffer for result archive
	buf := bytes.Buffer{}

	// Create zip processor
	zipWriter := zip.NewWriter(&buf)
	defer func() {
		if zipWriter != nil {
			errors.IfErrorLog(zipWriter.Close())
		}
	}()

	// Open game dir
	gamePathFiles, err := os.ReadDir(gamePath)
	if err != nil {
		errors.IfErrorLog(err)
		return nil, errors.InternalError.AddMessage("Error open game: " + err.Error())
	}

	// Go through all collections
	for _, gamePathFile := range gamePathFiles {
		// If dir - it's collection
		if gamePathFile.IsDir() {
			collectionPath := filepath.Join(gamePath, gamePathFile.Name())
			collectionPathFiles, err := os.ReadDir(collectionPath)
			if err != nil {
				errors.IfErrorLog(err)
				return nil, errors.InternalError.AddMessage("Error open collection: " + err.Error())
			}

			for _, collectionPathFile := range collectionPathFiles {
				// Path to file inside collection
				deckPathFile := filepath.Join(collectionPath, collectionPathFile.Name())
				// Path to file inside zip archive
				zipPath := filepath.Join(gameId, gamePathFile.Name(), collectionPathFile.Name())
				log.Println("Add file to archive:", deckPathFile, "->", zipPath)
				// Copy file inside archive
				err = addFileIntoArchive(deckPathFile, zipPath, zipWriter)
				if err != nil {
					return nil, err
				}
			}
			continue
		}

		// If file - it's game related
		// Path to game info file
		gamePathInfo := filepath.Join(gamePath, gamePathFile.Name())
		// Path to file inside zip archive
		zipPath := filepath.Join(gameId, gamePathFile.Name())
		log.Println("Add file to archive:", gamePathInfo, "->", zipPath)
		// Copy file inside archive
		err = addFileIntoArchive(gamePathInfo, zipPath, zipWriter)
		if err != nil {
			return nil, err
		}
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
