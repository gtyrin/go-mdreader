// MP3-files processing module.
// Header structure image link: https://upload.wikimedia.org/wikipedia/commons/thumb/0/01/Mp3filestructure.svg/1257px-Mp3filestructure.svg.png
// Specification link: http://www.mp3-tech.org/programmer/docs/mp3_theory.pdf

package file

import (
	"errors"
	"io"

	md "github.com/ytsiuryn/ds-audiomd"
	binary "github.com/ytsiuryn/go-binary"
	intutils "github.com/ytsiuryn/go-intutils"
)

// ErrMP3Ver25NotSupport ..
var ErrMP3Ver25NotSupport = errors.New("MPEG2.5 is not supported yet")

// ErrMP3WrongSyncWord ..
var ErrMP3WrongSyncWord = errors.New("wrong synchronization word")

// MPEGVersion defines MPEG version
type MPEGVersion byte

// MPEG version (MPEG2.5 is omitted)
const (
	MPEG2 MPEGVersion = iota
	MPEG1
)

// LayerType defines Layer type
type LayerType byte

// Layer Type
const (
	Reserved LayerType = iota
	Layer3
	Layer2
	Layer1
)

// BitrateMap defines existing bitrates
var BitrateMap = map[MPEGVersion]map[LayerType]map[byte]int{
	MPEG1: {
		Layer1: {
			1:  32,
			2:  64,
			3:  96,
			4:  128,
			5:  160,
			6:  192,
			7:  224,
			8:  256,
			9:  288,
			10: 320,
			11: 352,
			12: 384,
			13: 416,
			14: 448,
		},
		Layer2: {
			1:  32,
			2:  48,
			3:  56,
			4:  64,
			5:  80,
			6:  96,
			7:  112,
			8:  128,
			9:  160,
			10: 192,
			11: 224,
			12: 256,
			13: 320,
			14: 384,
		},
		Layer3: {
			1:  32,
			2:  40,
			3:  48,
			4:  56,
			5:  64,
			6:  80,
			7:  96,
			8:  112,
			9:  128,
			10: 160,
			11: 192,
			12: 224,
			13: 256,
			14: 320,
		},
	},
	MPEG2: {
		Layer1: {
			1:  32,
			2:  64,
			3:  96,
			4:  128,
			5:  160,
			6:  192,
			7:  224,
			8:  256,
			9:  288,
			10: 320,
			11: 352,
			12: 384,
			13: 416,
			14: 448,
		},
		Layer2: {
			1:  32,
			2:  48,
			3:  56,
			4:  64,
			5:  80,
			6:  96,
			7:  112,
			8:  128,
			9:  160,
			10: 192,
			11: 224,
			12: 256,
			13: 320,
			14: 384,
		},
		Layer3: {
			1:  8,
			2:  16,
			3:  24,
			4:  32,
			5:  64,
			6:  80,
			7:  56,
			8:  64,
			9:  128,
			10: 160,
			11: 112,
			12: 128,
			13: 256,
			14: 320,
		},
	},
}

// FrequencyMap describes audio samlerate
var FrequencyMap = map[MPEGVersion]map[byte]int{
	MPEG1: {
		0: 44100,
		1: 48000,
		2: 32000,
	},
	MPEG2: {
		0: 22050,
		1: 24000,
		2: 16000,
	},
}

// Mp3 is type for MP3 audio files processing.
type Mp3 struct {
	*md.Track
	release *md.Release
	r       *binary.Reader
}

// TrackMetadata gatheres metadata info for MP3 file
func (mp3 *Mp3) TrackMetadata(f io.ReadSeeker, release *md.Release, track *md.Track) error {
	mp3.release = release
	mp3.Track = track
	mp3.r = binary.NewReader(f)
	if ID3v2CheckSign(mp3.r) {
		processedTags, err := ID3v2Metadata(mp3.r, mp3.Track, mp3.release)
		if err != nil {
			return err
		}
		if err = ProcessTags(processedTags, release, track); err != nil {
			return err
		}
	}
	ret := mp3.headerInfo(f)
	track.LinkWithDisc(release.Disc(md.DiscNumberByTrackPos(track.Position)))
	return ret
}

func (mp3 *Mp3) headerInfo(f io.ReadSeeker) error {
	bitMask2 := mp3.r.ReadBEUint16()
	switch (bitMask2 & 0xfff0) >> 4 {
	case 0xffe:
		return ErrMP3Ver25NotSupport
	case 0xfff:
	default:
		return ErrMP3WrongSyncWord
	}
	mpegVersion := MPEGVersion((bitMask2 & 0x8) >> 3)
	layer := LayerType((bitMask2 & 0x6) >> 1)
	// bitmask & 0x1 for protection bit
	bitMask := mp3.r.ReadUint8()
	bitRate := byte((bitMask & 0xf0) >> 4)
	mp3.AudioInfo.AvgBitrate = int(BitrateMap[mpegVersion][layer][bitRate])
	mp3.AudioInfo.Samplerate = int(FrequencyMap[mpegVersion][byte((bitMask&0xc)>>2)])
	mp3.AudioInfo.SampleSize = 16
	// bitmask & 0x2 for padding bit
	// bitmask & 0x1 for private bit
	bitMask = mp3.r.ReadUint8()
	if ((bitMask & 0xc0) >> 6) == 3 {
		mp3.AudioInfo.Channels = 1
	} else {
		mp3.AudioInfo.Channels = 2
	}
	// TODO: неверно считает для тестового примера
	mp3.Duration = intutils.Duration(
		float64(8*mp3.FileInfo.FileSize) / float64(mp3.AudioInfo.AvgBitrate))
	return nil
}
