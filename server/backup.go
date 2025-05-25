package server

import (
	"crypto/sha256"
	u "dojer/utils"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

var noDirectoriesFound = errors.New("no directories found.")
var noChanges = errors.New("there was no changes since the last backup.")

func fileSum(name string) ([32]byte, error) {
	hasher := sha256.New()
	file, err := os.Open(name)
	if err != nil {
		return [32]byte{}, err
	}

	_, err = io.Copy(hasher, file)
	if err != nil {
		return [32]byte{}, err
	}

	return [32]byte(hasher.Sum(nil)), nil
}

func copyFile(name string, dest string) error {
	source, err := os.Open(name)
	if err != nil {
		return err
	}

	defer source.Close()
	destFolder := filepath.Dir(dest)
	_ = mkdirAll(destFolder)

	destination, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destination.Close()

	_, err = io.Copy(destination, source)

	if err != nil {
		return err
	}

	return nil
}

func getNewestDirectory(dir string) (string, error) {
	d, err := os.Open(dir)
	if err != nil {
		return "", err
	}
	defer d.Close()

	entries, err := d.Readdir(-1)
	if err != nil {
		return "", err
	}

	var dirs []os.FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry)
		}
	}

	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name() > dirs[j].Name()
	})

	if len(dirs) > 0 {
		result := filepath.Join(dir, dirs[0].Name())
		return result, nil
	}

	return "", noDirectoriesFound
}

func compareBackupHash(file1, file2 string) (bool, error) {
	f1Hash, err := fileSum(file1)
	if err != nil {
		return false, err
	}

	f2Hash, err := fileSum(file2)
	if err != nil {
		return false, err
	}

	return (f1Hash == f2Hash), nil
}

func backup(t time.Time) error {
	databaseFile := "database.sqlite3"
	backupsFolder := u.GetDataPath("backups")
	if err := mkdirAll(backupsFolder); err != nil {
		return err
	}

	dir, err := getNewestDirectory(backupsFolder)
	if err != nil {
		if errors.Is(noDirectoriesFound, err) {
			dir = "nil"
		} else {
			return err
		}
	}

	if dir != "nil" {
		olderBackup := filepath.Join(dir, databaseFile)

		sameContent, err := compareBackupHash(olderBackup, u.GetDataPath(databaseFile))
		if err != nil {
			return err
		}

		if sameContent {
			return noChanges
		}
	}

	backupName := t.Format("2006-01-02_15:04:05")
	backupPath := filepath.Join(backupsFolder, backupName)
	_ = mkdirAll(backupPath)
	backupFile := filepath.Join(backupPath, databaseFile) // datafolder/backups/2020_02_02_15:03:21/database.sqlite3

	err = copyFile(u.GetDataPath(databaseFile), backupFile)
	if err != nil {
		return err
	}

	return nil
}

func mkdirAll(path string) error {
	return os.MkdirAll(path, 0775)
}
