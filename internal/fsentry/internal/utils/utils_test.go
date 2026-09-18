package utils

import (
	"os"
	"path"
	"testing"
)

func FuzzNameToIDFile(f *testing.F) {
	// Check before starting fuzzing
	tmp, err := os.MkdirTemp("", "fsentry")
	if err != nil {
		f.Fatal("error creating temp dir", err)
	}
	err = os.Remove(tmp)
	if err != nil {
		f.Fatal("error removing temp dir", err)
	}

	tests := []string{
		"con",
		"prn",
		"aux",
		"nul",
		"com1", "com2", "com3", "com4", "com5", "com6", "com7", "com8", "com9", "com0",
		"lpt1", "lpt2", "lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9", "lpt0",
	}
	for _, tc := range tests {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, name string) {
		resName := NameToID(name)
		if resName == "" {
			t.Logf("Empty name: %q", name)
			return
		}

		tmpDirName, err := os.MkdirTemp("", "fsentry")
		if err != nil {
			t.Fatal("error creating temp dir", err)
		}
		defer os.RemoveAll(tmpDirName)

		filePath := path.Join(tmpDirName, resName)

		f, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("error create file %q: %s", resName, err.Error())
		}
		err = f.Close()
		if err != nil {
			t.Fatalf("error close file %q: %s", resName, err.Error())
		}
		err = os.Remove(filePath)
		if err != nil {
			t.Fatalf("error removing file %q: %s", resName, err.Error())
		}
	})
}

func FuzzNameToIDFolder(f *testing.F) {
	// Check before starting fuzzing
	tmp, err := os.MkdirTemp("", "fsentry")
	if err != nil {
		f.Fatal("error creating temp dir", err)
	}
	err = os.Remove(tmp)
	if err != nil {
		f.Fatal("error removing temp dir", err)
	}

	tests := []string{
		"con",
		"prn",
		"aux",
		"nul",
		"com1", "com2", "com3", "com4", "com5", "com6", "com7", "com8", "com9", "com0",
		"lpt1", "lpt2", "lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9", "lpt0",
	}
	for _, tc := range tests {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, name string) {
		resName := NameToID(name)
		if resName == "" {
			t.Logf("Empty name: %q", name)
			return
		}

		tmpDirName, err := os.MkdirTemp("", "fsentry")
		if err != nil {
			t.Fatal("error creating temp dir", err)
		}
		defer os.RemoveAll(tmpDirName)

		filePath := path.Join(tmpDirName, resName)

		err = os.Mkdir(filePath, 0755)
		if err != nil {
			t.Fatalf("error create folder %q: %s", resName, err.Error())
		}
		err = os.Remove(filePath)
		if err != nil {
			t.Fatalf("error removing folder %q: %s", resName, err.Error())
		}
	})
}
