package files

import (
	"encoding/gob"
	"errors"
	"home/pkg/lib/e"
	"home/pkg/lib/storage"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultPerm = 0774
)

type Storage struct {
	BasePath string
}

// New ...
func New(basePath string) *Storage {
	return &Storage{
		BasePath: basePath,
	}
}

// Save ...
func (s *Storage) Save(p *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save page", err) }()

	fPath := filepath.Join(s.BasePath, p.UserName)
	// вне зависимости от ОС будет использован правильный разделитель эелементов

	if err = os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(p)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := gob.NewEncoder(file).Encode(p); err != nil {
		return err
	}
	return nil
}

// PickRandom ...
func (s *Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random", err) }()

	path := filepath.Join(s.BasePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}
	rand.Seed(time.Now().UnixNano())

	n := rand.Intn(len(files))
	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

// Remove ...
func (s *Storage) Remove(p *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't remove file", err) }()
	fileName, err := fileName(p)
	if err != nil {
		return err
	}
	fPath := filepath.Join(s.BasePath, p.UserName, fileName)
	if err = os.Remove(fPath); err != nil {
		return e.Wrap(fileName, err)
	}

	return nil
}

// IsExist ...
func (s *Storage) IsExist(p *storage.Page) (isExist bool, err error) {
	defer func() { err = e.WrapIfErr("can't find file", err) }()
	fileName, err := fileName(p)
	if err != nil {
		return false, err
	}
	fPath := filepath.Join(s.BasePath, p.UserName, fileName)
	switch _, err = os.Stat(fPath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, e.Wrap(fPath, err)
	}
	return true, nil
}

// decodePage ...
func (s *Storage) decodePage(filePath string) (p *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't decode file "+filePath, err) }()
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	if err := gob.NewDecoder(f).Decode(p); err != nil {
		return nil, err
	}
	return
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
