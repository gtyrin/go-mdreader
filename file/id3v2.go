package file

import (
	"errors"
	"net/url"
	"strings"

	md "github.com/ytsiuryn/ds-audiomd"
	binary "github.com/ytsiuryn/go-binary"
)

const id3Sign = "ID3"

var (
	errID3NotFound = errors.New("ID3v2 section has incorrect sign mark")

	// excludingTags = []string{"ALBUM DYNAMIC RANGE", "ENCODER", "ENCODED BY",
	// 	"HDTRACKS", "RATING", "REPLAYGAIN_ALBUM_GAIN", "REPLAYGAIN_ALBUM_PEAK",
	// 	"TENC", "TOOL NAME", "TOOL VERSION", "TSSE"}
)

// ID3v2CheckSign reads 3 bytes and compares them with ID3 sign mark
func ID3v2CheckSign(r *binary.Reader) bool {
	return string(r.CheckBytes(3)) == id3Sign
}

// $00 ISO-8859-1 [ISO-8859-1]. Terminated with $00.
// $01 UTF-16 [UTF-16] encoded Unicode [UNICODE] with BOM. All strings in the same frame SHALL have the same byteorder. Terminated with $00 00.
// $02 UTF-16BE [UTF-16] encoded Unicode [UNICODE] without BOM. Terminated with $00 00.
// $03 UTF-8 [UTF-8] encoded Unicode [UNICODE]. Terminated with $00.”
func id3v2DecodeString(b []byte) (string, error) {
	var err error
	ret := b
	switch b[0] {
	case 0:
		ret = b[1:]
	case 1:
		ret, err = binary.FromUTF16LE(b[1:])
	case 2:
		ret, err = binary.FromUTF16BE(b[1:])
	}
	if err != nil {
		return "", err
	}
	return string(binary.FromASCIIZ(ret)), nil
}

// ID3v2Metadata is main fuction to read ID3 section data
func ID3v2Metadata(r *binary.Reader, track *md.Track, release *md.Release) (
	map[TagKey]string, error) {
	if !ID3v2CheckSign(r) {
		return nil, errID3NotFound
	}
	r.SkipBytes(6) // Sign mark(3), ID3Info.version(2), flags (1)
	sectionSize := parseBlockSize(r.ReadBytes(4))
	d := r.ReadBytes(sectionSize)
	var pos, frameSize int64
	var frameID string
	processedTags := make(map[TagKey]string)
	for pos < sectionSize {
		frameID = string(d[pos : pos+4])
		pos += 4
		frameSize = parseBlockSize(d[pos : pos+4])
		pos += 4 + 2 // 2(frame flags)
		if frameID == "APIC" {
			id3v2PictMetadata(r, d[pos:pos+frameSize], release)
		} else {
			frameValue, err := id3v2DecodeString(d[pos : pos+frameSize])
			if err != nil {
				return nil, err
			}
			if frameID == "COMM" {
				frameValue = frameValue[4:] // 3(lang)+1(0x0)
			} else if frameID == "TXXX" {
				flds := strings.SplitN(frameValue, "\x00", 2)
				frameID = "TXXX:" + flds[0]
				frameValue = flds[1]
			}
			if tag, ok := SchemaTagToUniKey[ID3v2][frameID]; ok {
				processedTags[tag] = frameValue
			} else {
				track.Unprocessed[frameID] = frameValue
			}
		}
		pos += frameSize
		// format alignment
		for ; pos < sectionSize && d[pos] == 0; pos++ {
		}
	}
	return processedTags, nil
}

// APIC tag processing
func id3v2PictMetadata(r *binary.Reader, frame []byte, release *md.Release) {
	if release.Cover() != nil {
		return
	}
	var pos, x uint32
	pict := md.PictureInAudio{PictureMetadata: &md.PictureMetadata{}}
	for ; frame[pos] == 0; pos++ {
	}
	for ; frame[pos+x] != 0; x++ {
	}
	pict.MimeType = string(frame[pos : pos+x])
	pos += x + 1
	pict.PictType = md.PictType(frame[pos])
	pos++
	for x = 0; frame[pos+x] != 0; x++ {
	}
	description := string(frame[pos : pos+x])
	pos += x + 1
	_, err := url.ParseRequestURI(description)
	if err == nil {
		pict.CoverURL = description
	} else {
		pict.Notes = description
	}
	pict.Size = uint32(len(frame)) - pos
	pict.Data = frame[pos:]
	release.Pictures = append(release.Pictures, &pict)
}

func parseBlockSize(b []byte) int64 {
	var n int64
	for _, x := range b {
		n <<= 7
		n |= int64(x)
	}
	return n
}
