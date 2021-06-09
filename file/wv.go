package file

// Specification: https://www.wavpack.com/WavPack5FileFormat.pdf
// Apev2 specification: https://wiki.hydrogenaud.io/index.php?title=APEv2_specification

import (
	encb "encoding/binary"
	"errors"
	"io"

	md "github.com/gtyrin/go-audiomd"
	binary "github.com/gtyrin/go-binary"
)

var (
	wvBlockSign = [4]byte{'w', 'v', 'p', 'k'}
)

// Public errors.
var (
	ErrWvBlockNotFound = errors.New("has no Wavpack block sign mark")
)

type wavpackHeader struct {
	CkID           [4]byte // "wvpk"
	CkSize         uint32  // size of entire block (minus 8)
	Version        uint16  // 0x402 to 0x410 are valid for decode
	BlockIndexU8   byte    // upper 8 bits of 40-bit block_index (since v.5)
	TotalSamplesU8 byte    // upper 8 bits of 40-bit total_samples (since v.5)
	// lower 32 bits of total samples for
	// entire file, but this is only valid
	// if block_index == 0 and a value of -1
	// indicates an unknown length
	TotalSamples uint32
	// lower 32 bit index of the first sample
	// in the block relative to file start,
	// normally this is zero in first block
	BlockIndex uint32
	// number of samples in this block, 0 =
	// non-audio block
	BlockSamples uint32
	// various flags for id and decoding
	// bits 1,0:   00 = 1 byte / sample (1-8 bits / sample)
	//             01 = 2 bytes / sample (9-16 bits / sample)
	//             10 = 3 bytes / sample (15-24 bits / sample)
	//             11 = 4 bytes / sample (25-32 bits / sample)
	// bit 2:      0 = stereo output; 1 = mono output
	// bit 3:      0 = lossless mode; 1 = hybrid mode
	// bit 4:      0 = true stereo; 1 = joint stereo (mid/side)
	// bit 5:      0 = independent channels; 1 = cross-channel decorrelation
	// bit 6:      0 = flat noise spectrum in hybrid; 1 = hybrid noise shaping
	// bit 7:      0 = integer data; 1 = floating point data
	// bit 8:      1 = extended size integers (> 24-bit) or shifted integers
	// bit 9:      0 = hybrid mode parameters control noise level (not used yet)
	//             1 = hybrid mode parameters control bitrate
	// bit 10:     1 = hybrid noise balanced between channels
	// bit 11:     1 = initial block in sequence (for multichannel)
	// bit 12:     1 = final block in sequence (for multichannel)
	// bits 17-13: amount of data left-shift after decode (0-31 places)
	// bits 22-18: maximum magnitude of decoded data (number of bits integers require minus 1)
	// bits 26-23: sampling rate (1111 = unknown/custom)
	// bit 27:     reserved (but decoders should ignore if set)
	// bit 28:     block contains checksum in last 2 or 4 bytes (ver 5.0+)
	// bit 29:     1 = use IIR for negative hybrid noise shaping
	// bit 30:     1 = false stereo (data is mono but output is stereo)
	// bit 31:     0 = PCM audio; 1 = DSD audio (ver 5.0+)
	Flags uint32
	Crc   uint32 // crc for actual decoded data
}

// type metadataSubblock struct {
// 	// 0x3f - unique metadata function id
// 	// 0x20 - decoder needn't understand metadata
// 	// 0x40 - actual data byte length is 1 less
// 	// 0x80 - large block (> 255 words)
// 	id byte
// 	// if small block: data size in words
// 	// if large block: data size in words (le)
// 	ws   [3]byte
// 	data uint16 // data, padded to an even number of bytes
// }

// Wv is type for Wavpack audio files processing.
type Wv struct {
	*md.Track
	release *md.Release
	r       *binary.Reader
}

// TrackMetadata gatheres Metadata info for Wavpack file
func (wv *Wv) TrackMetadata(f io.ReadSeeker, release *md.Release, track *md.Track) error {
	var err error
	wv.release = release
	wv.Track = track
	wv.r = binary.NewReader(f)
	if err = wv.readAudioProps(); err != nil {
		return err
	}
	if err = APEv2Metadata(wv.r, wv.Track, wv.release); err != nil {
		return err
	}
	track.LinkWithDisc(release.Disc(md.DiscNumberByTrackPos(track.Position)))
	return nil
}

func (wv *Wv) readAudioProps() error {
	header := wavpackHeader{}
	wv.r.ReadInto(int64(encb.Size(header)), encb.LittleEndian, &header)
	if header.CkID != wvBlockSign {
		return ErrWvBlockNotFound
	}
	// header.TotalSamples
	wv.r.SeekBytes(int64(header.CkSize)+8-32, io.SeekStart)
	// TODO: parse flags
	// header.Flags
	// bit 2 - 0 = stereo output; 1 = mono output
	// bits 26-23 - sampling rate
	// totalSamples := data & 0xfffffffff                        // 36 bits
	// flac.AudioInfo.SampleSize = int32((data>>36)&0x1f) + 1    // 5 bits
	// flac.AudioInfo.Channels = int32((data>>41)&0x7) + 1       // 3 bits
	// flac.AudioInfo.Samplerate = int32((data >> 44) & 0xfffff) // 20 bits
	// flac.Duration = int32(math.Round(
	// 	1000 * float64(totalSamples) / float64(flac.AudioInfo.Samplerate)))
	// flac.AudioInfo.AvgBitrate = int32(math.Round(
	// 	.008 * float64(flac.size) / float64(flac.Duration/1000)))
	return nil
}
