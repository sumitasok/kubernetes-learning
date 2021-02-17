package store

import (
	"errors"
	"fmt"
)

// Use of RDBMS instance would be better for persistance and easy query without having to duplicate values like done below (Filename and checksum based query)

// File holds checksum and name - creation of object is not allowed from outside this module.
// Making `File` exported as linted is giving a new error for having returned unexported store.file
// The intension is to have access to store.file via New factory method. However for this task, as I was not able to selectively suppress(Doesnt exist) this warning
// I am making this exported.
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

// GetFileByName gets a Files if exists
func (inMemory InMemory) GetFileByName(fname string) (File, error) {
	if f, ok := inMemory.Files[fname]; ok {
		return *f, nil
	}

	return File{}, errors.New("File not found")
}

// GetFileBySameData
// Update
// Delete

// DeleteFileByName deletes a Files if exists
func (inMemory InMemory) DeleteFileByName(fname string) (File, error) {
	f, ok := inMemory.Files[fname]
	if !ok {
		return File{}, errors.New("File not found")
	}

	delete(inMemory.Files, fname)
	// it is safe to do this even if checksum key is missing the map
	delete(inMemory.CheckSums, f.Checksum)

	return *f, nil
}
