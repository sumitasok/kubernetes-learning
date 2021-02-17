package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sumitasok/kubernetes-learning/internal/store"
)

func TestInMemory_Add(t *testing.T) {
	assert := assert.New(t)

	iM := store.NewInMemory()

	assert.NoError(iM.Add("fileA.txt", "check-sum-1"))
	assert.Error(iM.Add("fileB.txt", "check-sum-1"), "New error")
	assert.NoError(iM.Add("fileB.txt", "check-sum-2"))
	assert.Error(iM.Add("fileB.txt", "check-sum-3"), "New error 2")
}

func TestInMemory_GetFileByName(t *testing.T) {
	assert := assert.New(t)

	iM := store.NewInMemory()

	assert.NoError(iM.Add("fileA.txt", "check-sum-1"))
	_, err := iM.GetFileByName("fileA.txt")
	assert.NoError(err)
	_, err = iM.GetFileByName("fileB.txt")
	assert.Error(err)
}

func TestInMemory_DeleteFileByName(t *testing.T) {
	assert := assert.New(t)

	iM := store.NewInMemory()

	assert.NoError(iM.Add("fileA.txt", "check-sum-1"))
	f, err := iM.DeleteFileByName("fileA.txt")
	assert.NoError(err)
	assert.Equal(f.Checksum, "check-sum-1")
	_, err = iM.DeleteFileByName("fileA.txt")
	assert.Error(err)
}
