package fs

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"

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
