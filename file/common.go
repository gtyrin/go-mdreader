package file

import (
	"io"
	"path/filepath"
	"strings"

	// 	"regexp"
	// 	"strconv"

	md "github.com/gtyrin/go-audiomd"
)

// TrackMetadataReader - общий интерфейс читателей метаданных аудиотреков различных расширений.
type TrackMetadataReader interface {
	// Извлечь метаданные трекфайла.
	TrackMetadata(f io.ReadSeeker, release *md.Release, track *md.Track) error
}

// // RutrackerRegexp is a regexp for Rutracker.org URL
// var RutrackerRegexp = regexp.MustCompile(`^http[s]?:\/\/rutracker\.org\/forum\/viewtopic\.php\?t=(\d+)\s*`)

var infoLoaders = map[string]TrackMetadataReader{
	".dsf":  new(Dsf),
	".flac": new(Flac),
	".wv":   new(Wv),
	".mp3":  new(Mp3),
}

// Reader returns TrackMetadataReader of the appropriate type or nil.
func Reader(fn string) TrackMetadataReader {
	if cls, ok := infoLoaders[filepath.Ext(strings.ToLower(fn))]; ok {
		return cls
	}
	return nil
}
