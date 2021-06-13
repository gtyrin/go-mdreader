package file

import (
	encb "encoding/binary"
	"errors"
	"io"
	"strings"

	md "github.com/ytsiuryn/ds-audiomd"
	binary "github.com/ytsiuryn/go-binary"
	collection "github.com/ytsiuryn/go-collection"
)

var (
	apeMetadataSign = [8]byte{'A', 'P', 'E', 'T', 'A', 'G', 'E', 'X'}
	pictureTags     = []string{"COVER ART (FRONT)", "COVER ART (BACK)"}
)

var (
	errApev2NotFound = errors.New("has no APEv2 metadata sign mark")
)

type apeTagsHeader struct {
	Preamble [8]byte // "APETAGEX"
	Version  uint32  // 1000 = Version 1.000 (old); 2000 = Version 2.000 (new)
	// Tag size in bytes including footer and all tag items excluding
	// the header to be as compatible as possible with APE Tags 1.000
	TagSize   uint32
	ItemCount uint32 // Number of items in the Tag (n)
	TagFlags  uint32 // Global flags of all items
	Reserved  uint64
}

// 0b0000 	get from STREAMINFO metadata block
// 0b0001 	88.2 kHz
// 0b0010 	176.4 kHz
// 0b0011 	192 kHz
// 0b0100 	8 kHz
// 0b0101 	16 kHz
// 0b0110 	22.05 kHz
// 0b0111 	24 kHz
// 0b1000 	32 kHz
// 0b1001 	44.1 kHz
// 0b1010 	48 kHz
// 0b1011 	96 kHz
// 0b1100 	get 8 bit sample rate (in kHz) from end of header
// 0b1101 	get 16 bit sample rate (in Hz) from end of header
// 0b1110 	get 16 bit sample rate (in tens of Hz) from end of header
// 0b1111 	invalid, to prevent sync-fooling string of 1s

// APEv2Metadata чтение и парсинг блока метаданных.
func APEv2Metadata(r *binary.Reader, track *md.Track, release *md.Release) error {
	header := apeTagsHeader{}
	headerSize := int64(encb.Size(header))
	pos := r.SeekBytes(-headerSize, io.SeekEnd)
	if r.ReadInto(headerSize, encb.LittleEndian, &header); header.Preamble != apeMetadataSign {
		return errApev2NotFound
	}
	pos += (-int64(header.TagSize) + headerSize)
	r.SeekBytes(pos, io.SeekStart)
	var tagName, tagVal string
	var itemLen int64
	m := make(map[TagKey]string)
	for i := 0; i < int(header.ItemCount); i++ {
		itemLen = int64(r.ReadLEUint32())
		r.SkipBytes(4) // flags
		tagName = strings.ToUpper(r.ReadString())
		if apev2IsCoverTag(tagName) {
			apev2PictMetadata(r, itemLen, release)
		} else {
			tagVal = string(r.ReadBytes(itemLen))
			if tag, ok := SchemaTagToUniKey[APEv2][tagName]; ok {
				m[tag] = tagVal
			} else {
				track.Unprocessed[tagName] = tagVal
			}
		}
	}
	if err := ProcessTags(m, release, track); err != nil {
		return err
	}
	return nil
}

func apev2IsCoverTag(tagName string) bool {
	return collection.ContainsStr(tagName, pictureTags)
}

func apev2PictMetadata(r *binary.Reader, nbytes int64, release *md.Release) {
	if release.Cover() != nil {
		r.SkipBytes(nbytes)
		return
	}
	data := r.ReadBytes(nbytes)
	picture := md.PictureInAudio{PictureMetadata: &md.PictureMetadata{}}
	picture.Data = data
	release.Pictures = append(release.Pictures, &picture)
}
