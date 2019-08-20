package persistence

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReadWriteTime(t *testing.T) {
	persistence := Dir{Path: "../../testdata"}
	err := persistence.Init()
	assert.Nil(t, err)

	now := time.Now()

	err = WriteTime(&persistence, "pr/key1", now)
	assert.Nil(t, err)

	now2, err := ReadTime(&persistence, "pr/key1");
	assert.Nil(t, err)
	assert.Equal(t, now.Unix(), now2.Unix())
}

func TestReadFirst(t *testing.T) {
	persistence := Dir{Path: "../../testdata"}
	err := persistence.Init()
	assert.Nil(t, err)

	tm, err := ReadTime(&persistence, "pr/key2")
	assert.Nil(t, err)

	assert.Equal(t, int64(0), tm.Unix())
}
