// FLAC processing module.
// Specification link: https://xiph.org/flac/format.html

package file

import (
	"errors"
	"io"
	"math"
	"net/url"
	"strings"

	encb "encoding/binary"

	md "github.com/ytsiuryn/ds-audiomd"
	binary "github.com/ytsiuryn/go-binary"
	intutils "github.com/ytsiuryn/go-intutils"
)

const (
	flacSign           = "fLaC"
	isLastBlock        = 1
	streamInfoBlock    = 0
	vorbisCommentBlock = 4
	cueSheetBlock      = 5
	pictureBlock       = 6
)

// Public errors
var (
	ErrFLACNoSign                    = errors.New("has no FLAC sign mark")
	ErrFLACInfoblockSize             = errors.New("incorrect streamInfoBlock section size")
	ErrFLACIncorrectVorbisComment    = errors.New("vorbis comment has illegal structure")
	ErrFLACIncorrectPictureblockSize = errors.New("incorrect mdBlockPicture size")
)

// Flac is type for FLAC audio files processing.
type Flac struct {
	*md.Track
	release *md.Release
	r       *binary.Reader
}

// TrackMetadata gatheres Metadata info for Flac file
func (flac *Flac) TrackMetadata(f io.ReadSeeker, release *md.Release, track *md.Track) error {
	flac.release = release
	flac.Track = track
	flac.r = binary.NewReader(f)
	if ID3v2CheckSign(flac.r) {
		tagsToProcess, err := ID3v2Metadata(flac.r, flac.Track, flac.release)
		if err != nil {
			return err
		}
		if err = ProcessTags(tagsToProcess, flac.release, flac.Track); err != nil {
			return err
		}
	}
	if string(flac.r.ReadBytes(4)) != flacSign {
		return ErrFLACNoSign
	}
	if err := flac.mdBlocks(); err != nil {
		return err
	}
	track.LinkWithDisc(release.Disc(md.DiscNumberByTrackPos(track.Position)))
	return nil
}

// Common block processing
func (flac *Flac) mdBlocks() error {
	var processedTags map[TagKey]string
	var err error
	var x uint32
	for {
		x = flac.r.ReadBEUint32()
		b := (x >> 24)
		blDataLen := int64(x & 0xffffff)
		switch b & 0x7f { // blType
		case streamInfoBlock:
			err = flac.mdBlockStreamInfo(blDataLen)
		case vorbisCommentBlock:
			if processedTags, err = flac.mdBlockVorbisComment(blDataLen); err == nil {
				err = ProcessTags(processedTags, flac.release, flac.Track)
			}
		case cueSheetBlock:
			err = flac.mdBlockCueSheet(blDataLen)
		case pictureBlock:
			err = flac.mdBlockPicture(blDataLen)
		default:
			flac.r.SkipBytes(blDataLen)
		}
		if err != nil {
			return err
		}
		if (b&0x80)>>7 == isLastBlock {
			break
		}
	}
	return nil
}

// Audio properties section data processing
func (flac *Flac) mdBlockStreamInfo(blDataLen int64) error {
	if blDataLen != 34 {
		return ErrFLACInfoblockSize
	}
	d := flac.r.ReadBytes(blDataLen)
	// Skip 10 bytes: MinBlSize:16, MaxBlSize:16, MinFrameSize:24, MaxFrameSize:24
	data := encb.BigEndian.Uint64(d[10:18])
	totalSamples := data & 0xfffffffff                      // 36 bits
	flac.AudioInfo.SampleSize = int((data>>36)&0x1f) + 1    // 5 bits
	flac.AudioInfo.Channels = int((data>>41)&0x7) + 1       // 3 bits
	flac.AudioInfo.Samplerate = int((data >> 44) & 0xfffff) // 20 bits
	flac.Duration = intutils.Duration(math.Round(
		1000 * float64(totalSamples) / float64(flac.AudioInfo.Samplerate)))
	flac.AudioInfo.AvgBitrate = int(math.Round(
		.008 * float64(flac.FileInfo.FileSize) / float64(flac.Duration/1000)))
	// The last 16 bytes: Md5Sum:128
	return nil
}

// Vorbis metadata processing
func (flac *Flac) mdBlockVorbisComment(blDataLen int64) (map[TagKey]string, error) {
	var frameID, val string
	processedTags := make(map[TagKey]string)
	d := flac.r.ReadBytes(blDataLen)
	x := encb.LittleEndian.Uint32(d[:4])
	pos := x + 4 // skip LibData
	fldCounter := encb.LittleEndian.Uint32(d[pos : pos+4])
	pos += 4
	for fldCounter > 0 {
		x = encb.LittleEndian.Uint32(d[pos : pos+4])
		pos += 4
		fldData := d[pos : pos+x]
		fields := strings.SplitN(string(fldData), "=", 2)
		if len(fields) != 2 {
			return nil, ErrFLACIncorrectVorbisComment
		}
		frameID = strings.ToUpper(fields[0])
		val = strings.TrimSpace(fields[1])
		if tag, ok := SchemaTagToUniKey[VorbisComment][frameID]; ok {
			processedTags[tag] = val
		} else {
			flac.Unprocessed[frameID] = val
		}
		fldCounter--
		pos += x
	}
	return processedTags, nil
}

// for CD-DA track ISRC extraction
func (flac *Flac) mdBlockCueSheet(blDataLen int64) error {
	flac.r.SkipBytes(9) // Track offset in samples, Track number
	flac.Track.SetISRC(string(flac.r.ReadBytes(12)))
	flac.r.SkipBytes(15) // other fields
	return nil
}

// Cover picture data processing
func (flac *Flac) mdBlockPicture(blDataLen int64) error {
	if flac.release.Cover() != nil {
		flac.r.SkipBytes(blDataLen)
		return nil
	}
	picture := md.PictureInAudio{PictureMetadata: &md.PictureMetadata{}}
	d := flac.r.ReadBytes(blDataLen)
	picture.PictType = md.PictType(encb.BigEndian.Uint32(d[:4]))
	x := encb.BigEndian.Uint32(d[4:8])
	picture.MimeType = string(d[8 : 8+x])
	pos := 8 + x
	x = encb.BigEndian.Uint32(d[pos : 4+pos])
	pos += 4
	description := strings.TrimSpace(string(d[pos : pos+x]))
	_, err := url.ParseRequestURI(description)
	if err == nil {
		picture.CoverURL = description
	} else {
		picture.Notes = description
	}
	pos += x
	picture.Width = encb.BigEndian.Uint32(d[pos : pos+4])
	pos += 4
	picture.Height = encb.BigEndian.Uint32(d[pos : pos+4])
	pos += 4
	picture.ColorDepth = encb.BigEndian.Uint32(d[pos : pos+4])
	pos += 4
	picture.Colors = encb.BigEndian.Uint32(d[pos : pos+4])
	pos += 4
	picture.Size = encb.BigEndian.Uint32(d[pos : pos+4])
	pos += 4
	if blDataLen != 32+int64(len(picture.MimeType)+len(description))+int64(picture.Size) {
		return ErrFLACIncorrectPictureblockSize
	}
	picture.Data = d[pos:]
	flac.release.Pictures = append(flac.release.Pictures, &picture)
	return nil
}
