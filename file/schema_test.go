package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApev2IsCoverTag(t *testing.T) {
	assert.True(t, apev2IsCoverTag("COVER ART (FRONT)"))
	assert.False(t, apev2IsCoverTag("FOLDER PICTURE"))
}

// TODO: Создать минимальные тестовые данные для функции.
func TestApev2PictMetadata(t *testing.T) {

}

func TestID3v2FrameSize(t *testing.T) {
	d := map[int][]byte{201666: {0, 12, 39, 66}}
	for _, bytes := range d {
		frameSize := parseBlockSize(bytes)
		assert.Contains(t, d, int(frameSize))
	}
}
