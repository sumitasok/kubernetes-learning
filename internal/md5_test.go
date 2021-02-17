package internal_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sumitasok/kubernetes-learning/internal"
)

func TestMd5(t *testing.T) {
	assert.NotEmpty(t, os.Getenv("GOPATH"), "please set GOPATH in environment for tests to run smoothly")
	md5Value, err := internal.Md5(os.Getenv("GOPATH") + "/src/github.com/sumitasok/kubernetes-learning/samples/download.jpeg")

	assert.Equal(t, md5Value, "49ec43e9d11781ddc9ac1a9a95c71801", "Md5 doesn't match")
	assert.NoError(t, err)
}
