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
