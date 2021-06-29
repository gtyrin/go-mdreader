package file

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	md "github.com/ytsiuryn/ds-audiomd"
	collection "github.com/ytsiuryn/go-collection"
	intutils "github.com/ytsiuryn/go-intutils"

	"github.com/stretchr/testify/assert"
)

// Тестовые файлы.
const (
	testFLAC = "testdata/flac/440_hz_mono.flac"
	testDSF  = "testdata/dsf/440_hz_mono.dsf"
	testMP3  = "testdata/mp3/440_hz_mono.mp3"
	testWV   = "testdata/wavpack/440_hz_mono.wv"
)

var testFileData = map[string]TrackMetadataReader{
	testFLAC: new(Flac),
	testMP3:  new(Mp3),
	testDSF:  new(Dsf),
	testWV:   new(Wv),
}

func TestTrackMetadataReader(t *testing.T) {
	testData := map[string]TrackMetadataReader{
		"test.Dsf":  new(Dsf),
		"test.fLaC": new(Flac),
		"test.WV":   new(Wv),
		"test.mp3":  new(Mp3),
		"test.txt":  nil,
	}
	for filename, reader := range testData {
		if reflect.TypeOf(Reader(filename)) != reflect.TypeOf(reader) {
			t.Fail()
		}
	}
}

func TestAudioFiles(t *testing.T) {
	for fn, r := range testFileData {
		var release = md.NewRelease()
		var track = md.NewTrack()
		f := initTestFileEnvironment(fn, track)
		defer f.Close()
		if err := r.TrackMetadata(f, release, track); err != nil {
			t.Error(err)
		}
		checkMdReaderResult(t, release, track)
	}
}

func firstActor(actors md.ActorIDs) md.ActorName {
	for k := range actors {
		return k
	}
	return ""
}

func checkMdReaderResult(t *testing.T, r *md.Release, tr *md.Track) {
	ext := filepath.Ext(tr.FileInfo.FileName)
	basename := strings.TrimSuffix(tr.FileInfo.FileName, ext)
	assert.Equal(t, r.Title, "test_album_title")
	assert.Equal(t, firstActor(r.Actors), md.ActorName("test_performer"))
	assert.Equal(t, tr.Record.Genres[0], "test_genre")
	assert.Equal(t, r.Discs[0].Number, 1)
	assert.Equal(t, r.TotalTracks, 10)
	assert.Equal(t, firstActor(tr.Record.Actors), md.ActorName("test_track_artist"))
	assert.Equal(t, tr.Position, "03")
	assert.Equal(t, tr.Title, "test_track_title")
	if !collection.ContainsStr(ext, []string{".mp3", ".wv"}) { // TODO
		assert.Equal(t, tr.Duration, intutils.Duration(500))
	}
	assert.Equal(t, basename, "440_hz_mono")
	if !collection.ContainsStr(ext, []string{".dsf", ".wv"}) { // TODO
		assert.Equal(t, tr.AudioInfo.Samplerate, 44100)
		assert.Equal(t, tr.AudioInfo.SampleSize, 16)
	}
	if ext != ".wv" { // TODO
		assert.Equal(t, tr.AudioInfo.Channels, 1)
		assert.Equal(t, r.Pictures[0].PictureMetadata.MimeType, "image/jpeg")
		assert.Equal(t, r.Pictures[0].PictType, md.PictType(3))
	}
	assert.Equal(t, r.Publishing[0].Name, "test_label")
	assert.Equal(t, r.Publishing[0].Catno, "test_catno")
	assert.Equal(t, r.Country, "test_country")
	assert.Equal(t, r.Year, 2000)
	assert.Equal(t, tr.Notes, "test_notes")
}

func initTestFileEnvironment(fn string, track *md.Track) *os.File {
	f, _ := os.OpenFile(fn, os.O_RDONLY, 0444)
	fi, _ := f.Stat()
	track.FileInfo.FileName = fi.Name()
	track.FileInfo.ModTime = fi.ModTime().Unix()
	track.FileInfo.FileSize = fi.Size()
	return f
}
