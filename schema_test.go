package service

import "testing"

func TestApev2IsCoverTag(t *testing.T) {
	if !apev2IsCoverTag("COVER ART (FRONT)") {
		t.Fail()
	}
	if apev2IsCoverTag("FOLDER PICTURE") {
		t.Fail()
	}
}

// TODO: Создать минимальные тестовые данные для функции.
func TestApev2PictMetadata(t *testing.T) {

}

func TestID3v2FrameSize(t *testing.T) {
	d := map[int][]byte{201666: {0, 12, 39, 66}}
	for size, bytes := range d {
		frameSize := parseBlockSize(bytes)
		if _, ok := d[int(frameSize)]; !ok {
			t.Errorf("wants %d, has %d", size, frameSize)
		}
	}
}
