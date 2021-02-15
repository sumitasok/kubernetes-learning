package store

import (
	"errors"
	"fmt"
)

// File holds checksum and name - creation of object is not allowed from outside this module.
type File struct {
	Checksum string
	Name     string
}

// NewFile factory for creating File instance
func NewFile(fname string, checksum string) (File, error) {
	return File{Name: fname, Checksum: checksum}, nil
}

// InMemory is a struct that acts as a database for holding File info
// Need to check how rsync works without a database https://serverfault.com/questions/556831/building-file-list-database-with-rsync
type InMemory struct {
	Files     map[string]*File
	CheckSums map[string]*File
}

// NewInMemory returna new instance of InMemory store
func NewInMemory() *InMemory {
	return &InMemory{
		Files:     map[string]*File{},
		CheckSums: map[string]*File{},
	}
}

// ErrFileExists error in case File already exist in the database
var ErrFileExists = errors.New("File exists")

// ErrSimilarFileExists error in case File already exist in the database
var ErrSimilarFileExists = "Similar File exists with a different name, File: %s"

// Add adds a File if already doesn't exists
func (inMemory InMemory) Add(fname string, checksum string) error {
	if _, ok := inMemory.Files[fname]; ok {
		return ErrFileExists
	}

	if f, ok := inMemory.CheckSums[checksum]; ok {
		return fmt.Errorf(ErrSimilarFileExists, f.Name)
	}

	f, err := NewFile(fname, checksum)

	inMemory.Files[fname] = &f
	inMemory.CheckSums[checksum] = &f

	return err
}
