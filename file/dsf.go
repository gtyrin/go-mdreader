// DSF-files processing module.
// Specification link: http://dsd-guide.com/sites/default/files/white-papers/DSFFileFormatSpec_E.pdf

package file

import (
	"bytes"
	encb "encoding/binary"
	"errors"
	"io"
	"math"

	md "github.com/gtyrin/go-audiomd"
	binary "github.com/gtyrin/go-binary"
	intutils "github.com/gtyrin/go-intutils"
)

// DSF Sign marks
const (
	DSFSign = "DSD "
	FmtSign = "fmt "
)

// Public errors.
var (
	ErrDSFNoSignMark        = errors.New("DSF has no sign mark")
	ErrIncorrectDSFChunk    = errors.New("incorrect DSF chunk")
	ErrDSFIncorrectFileSize = errors.New("incorrect file size")
	ErrDSFIncorrectFMTChunk = errors.New("incorrect FMT chunk")
)

// ChannelType describes actual channel count.
type ChannelType byte

// ChannelType constants
const (
	Mono ChannelType = 1
	Stereo
	Channels3
	Quad
	Channels4
	Channels5
	Channels51
)

// Dsf is type for DSF audio files processing.
type Dsf struct {
	*md.Track
	release *md.Release
	r       *binary.Reader
}

// TrackMetadata читает метаданные трек-файла и жобавляет объект Track в коллекцию треков релиза.
func (dsf *Dsf) TrackMetadata(f io.ReadSeeker, release *md.Release, track *md.Track) error {
	dsf.release = release
	dsf.Track = track
	dsf.r = binary.NewReader(f)
	mdChunkOffset, err := dsf.chunk(f)
	if err != nil {
		return err
	}
	if err = dsf.fmtChunk(); err != nil {
		return err
	}
	// MetadataBlockSize = FileSize - MdChunkOffset
	dsf.r.SeekBytes(mdChunkOffset, io.SeekStart)
	processedTags, err := ID3v2Metadata(dsf.r, dsf.Track, dsf.release)
	if err != nil {
		return err
	}
	if err = ProcessTags(processedTags, release, track); err != nil {
		return err
	}
	track.LinkWithDisc(release.Disc(md.DiscNumberByTrackPos(track.Position)))
	return nil
}

// DSF chunk 28 bytes (4 + 8 + 8 + 8)
// returns metadata chunk offset (or -1) and error.
func (dsf *Dsf) chunk(f io.ReadSeeker) (int64, error) {
	data := dsf.r.ReadBytes(28)
	if string(data[:4]) != DSFSign {
		return -1, ErrDSFNoSignMark
	}
	if !bytes.Equal(data[4:12], []byte{28, 0, 0, 0, 0, 0, 0, 0}) {
		return -1, ErrIncorrectDSFChunk
	}
	if int64(encb.LittleEndian.Uint64(data[12:20])) != dsf.FileInfo.FileSize {
		return -1, ErrDSFIncorrectFileSize
	}
	return int64(encb.LittleEndian.Uint64(data[20:28])), nil
}

// FMT chunk 52 bytes
// Saves audio stream info
func (dsf *Dsf) fmtChunk() error {
	data := dsf.r.ReadBytes(52)
	if string(data[:4]) != FmtSign {
		return ErrDSFIncorrectFMTChunk
	}
	// skip 20 bytes: ChunkSize:64, FmtVer:32, FmtID:32, ChannelType:32
	dsf.AudioInfo.Channels = int(encb.LittleEndian.Uint32(data[24:28]))
	dsf.AudioInfo.Samplerate = int(encb.LittleEndian.Uint32(data[28:32]))
	dsf.AudioInfo.SampleSize = int(encb.LittleEndian.Uint32(data[32:36]))
	sampleCount := int64(encb.LittleEndian.Uint64(data[36:44]))
	dsf.Duration = intutils.Duration(sampleCount * 1000 / int64(dsf.AudioInfo.Samplerate))
	dsf.AudioInfo.AvgBitrate = int(
		math.Round(8 * float64(dsf.FileInfo.FileSize) / float64(dsf.Duration)))
	// dsf.TechInfo.SetDuration(int32(math.Round(
	// 	.008 * float64(sampleCount) / float64(dsf.TechInfo.Samplerate))))
	// file.SeekBytes(f, 8, io.SeekCurrent) // BlSizePerChannel:32, Reserved:32
	return nil
}
