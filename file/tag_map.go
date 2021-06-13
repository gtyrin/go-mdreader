// Модуль описания обобщенных тегов и их единообразного преобразования в систему метаданных,
// принятую в DS.
//
// По материалам из: https://wiki.hydrogenaud.io/index.php?title=Tag_Mapping

package file

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	md "github.com/ytsiuryn/ds-audiomd"
	intutils "github.com/ytsiuryn/go-intutils"
	"github.com/ytsiuryn/go-stringutils"
)

// TagKey - тип для обозначения обобщенных констант.
type TagKey uint8

// TagScheme - тип для обозначения схем кодирования метаданных.
type TagScheme uint8

// Описание поддерживаемых системой схем кодирования.
const (
	ID3v2         TagScheme = iota // Flac, Dsf, Wv
	VorbisComment                  // Flac
	APEv2                          // Wv
)

// 	Обобщенные теги для различных схем теггирования.
const (
	// Titles
	AlbumTitle TagKey = iota
	DiscSetSubtitle
	ContentGroup
	TrackTitle
	TrackSubtitle
	Version
	// People & Organizations
	AlbumArtist
	TrackArtist
	Arranger
	AuthorWriter
	Writer
	Composer
	Conductor
	Engineer
	Ensemble
	InvolvedPeople
	Lyricist
	MixDJ
	MixEngineer
	MusicianCredits
	Organisation
	OriginalArtist
	Performer
	Producer
	Publisher
	Label
	LabelNumber
	RemixedBy
	Soloists
	// Counts & Indexes
	DiscNumber
	DiscTotal
	TrackNumber
	TrackTotal
	PartNumber
	Length
	// Dates
	ReleaseDate
	Year
	OriginalReleaseDate
	RecordingDates
	// Identifiers
	ISRC
	Barcode
	CatalogueNumber
	UPC
	DiscID
	AccurateRipDiscID
	DiscogsReleaseID
	MusicbrainzAlbumID
	RutrackerID
	// Flags
	Compilation
	// Ripping & Encoding
	FileType
	MediaType
	SourceMedia
	Source
	// URLs
	AudioSourceWebpageURL
	CommercialInformationURL
	TrackArtistWebPageURL
	// Style
	Genre
	Mood
	Style
	// Miscellaneous
	Country
	Comments
	Description
	CopyrightMessage
	SyncedLyrics
	UnsyncedLyrics
	Language
)

// SchemaTagToUniKey - соответствие тега определенной схемы кодирования единому коду тега.
var SchemaTagToUniKey = map[TagScheme]map[string]TagKey{
	ID3v2: {
		"TALB":                     AlbumTitle,
		"TSST":                     DiscSetSubtitle,
		"TIT1":                     ContentGroup,
		"TIT2":                     TrackTitle,
		"TIT3":                     TrackSubtitle,
		"TPE2":                     AlbumArtist,
		"TPE1":                     TrackArtist,
		"IPLS:arranger":            Arranger,
		"TIPL:arranger":            Arranger,
		"TEXT":                     AuthorWriter,
		"TCOM":                     Composer,
		"TPE3":                     Conductor,
		"IPLS:engineer":            Engineer,
		"TIPL:engineer":            Engineer,
		"IPLS":                     InvolvedPeople,
		"TIPL":                     InvolvedPeople,
		"IPLS:DJ-mix":              MixDJ,
		"TIPL:DJ-mix":              MixDJ,
		"IPLS:mix":                 MixEngineer,
		"TIPL:mix":                 MixEngineer,
		"TOPE":                     OriginalArtist,
		"TMCL":                     Performer,
		"IPLS:producer":            Producer,
		"TIPL:producer":            Producer,
		"TPUB":                     Publisher,
		"TXXX:LABEL":               Publisher,
		"TPE4":                     RemixedBy,
		"TPOS":                     DiscNumber,
		"TRCK":                     TrackNumber,
		"TXXX:TRACKTOTAL":          TrackTotal,
		"TLEN":                     Length,
		"TDRC":                     ReleaseDate,
		"TDAT":                     ReleaseDate,
		"TYER":                     Year,
		"TORY":                     OriginalReleaseDate,
		"TDOR":                     OriginalReleaseDate,
		"TRDA":                     RecordingDates,
		"TSRC":                     ISRC,
		"DISCID":                   DiscID,
		"TXXX:BARCODE":             Barcode,
		"TXXX:CATALOGNUMBER":       CatalogueNumber,
		"TXXX:DISCOGS_RELEASE_ID":  DiscogsReleaseID,
		"TXXX:MUSICBRAINZ_ALBUMID": MusicbrainzAlbumID,
		"TXXX:RUTRACKER":           RutrackerID,
		"TCMP":                     Compilation,
		"TFLT":                     FileType,
		"TMED":                     MediaType,
		"WOAS":                     AudioSourceWebpageURL,
		"WCOM":                     CommercialInformationURL,
		"WOAR":                     TrackArtistWebPageURL,
		"TCON":                     Genre,
		"TMOO":                     Mood,
		"TXXX:RELEASECOUNTRY":      Country,
		"COMM":                     Comments,
		"TCOP":                     CopyrightMessage,
		"SYLT":                     SyncedLyrics,
		"USLT":                     UnsyncedLyrics,
		"TLAN":                     Language,
	},
	VorbisComment: {
		"ALBUM":               AlbumTitle,
		"DISCSUBTITLE":        DiscSetSubtitle,
		"GROUPING":            ContentGroup,
		"TITLE":               TrackTitle,
		"SUBTITLE":            TrackSubtitle,
		"VERSION":             Version,
		"ALBUMARTIST":         AlbumArtist,
		"ARTIST":              TrackArtist,
		"ARRANGER":            Arranger,
		"AUTHOR":              AuthorWriter,
		"WRITER":              Writer,
		"COMPOSER":            Composer,
		"CONDUCTOR":           Conductor,
		"ENGINEER":            Engineer,
		"ENSEMBLE":            Ensemble,
		"LYRICIST":            Lyricist,
		"LANGUAGE":            Language,
		"MIXER":               MixEngineer,
		"ORGANIZATION":        Organisation,
		"PERFORMER":           Performer,
		"PRODUCER":            Producer,
		"PUBLISHER":           Publisher,
		"LABEL":               Label,
		"LABELNO":             LabelNumber,
		"REMIXER":             RemixedBy,
		"SOLOISTS":            Soloists,
		"DISCNUMBER":          DiscNumber,
		"DISCTOTAL":           DiscTotal,
		"TOTALDISCS":          DiscTotal,
		"TRACKNUMBER":         TrackNumber,
		"TRACKTOTAL":          TrackTotal,
		"TOTALTRACKS":         TrackTotal,
		"PARTNUMBER":          PartNumber,
		"DATE":                ReleaseDate,
		"ORIGINALDATE":        OriginalReleaseDate,
		"ISRC":                ISRC,
		"BARCODE":             Barcode,
		"CATALOGNUMBER":       CatalogueNumber,
		"UPC":                 UPC,
		"DISCOGS_RELEASE_ID":  DiscogsReleaseID,
		"MUSICBRAINZ_ALBUMID": MusicbrainzAlbumID,
		"RUTRACKER":           RutrackerID,
		"DISCID":              DiscID,
		"ACCURATERIPDISCID":   AccurateRipDiscID,
		"COMPILATION":         Compilation,
		"MEDIA":               MediaType,
		"SOURCEMEDIA":         SourceMedia,
		"SOURCE":              Source,
		"GENRE":               Genre,
		"MOOD":                Mood,
		"STYLE":               Style,
		"RELEASECOUNTRY":      Country,
		"COMMENT":             Comments,
		"DESCRIPTION":         Description,
		"COPYRIGHT":           CopyrightMessage,
	},
	APEv2: {
		"ALBUM":               AlbumTitle,
		"DISCSUBTITLE":        DiscSetSubtitle,
		"GROUPING":            ContentGroup,
		"TITLE":               TrackTitle,
		"SUBTITLE":            TrackSubtitle,
		"ALBUMARTIST":         AlbumArtist,
		"ARTIST":              TrackArtist,
		"ARRANGER":            Arranger,
		"WRITER":              Writer,
		"COMPOSER":            Composer,
		"CONDUCTOR":           Conductor,
		"Enginee":             Engineer,
		"LYRICIST":            Lyricist,
		"LANGUAGE":            Language,
		"MIXER":               MixEngineer,
		"PERFORMER":           Performer,
		"PRODUCER":            Producer,
		"LABEL":               Label,
		"MIXARTIST":           RemixedBy,
		"DISC":                DiscNumber,
		"TRACK":               TrackNumber,
		"TRACKTOTAL":          TrackTotal,
		"YEAR":                Year,
		"ISRC":                ISRC,
		"BARCODE":             Barcode,
		"CATALOGNUMBER":       CatalogueNumber,
		"DISCOGS_RELEASE_ID":  DiscogsReleaseID,
		"MUSICBRAINZ_ALBUMID": MusicbrainzAlbumID,
		"RUTRACKER":           RutrackerID,
		"COMPILATION":         Compilation,
		"MEDIA":               MediaType,
		"SOURCEMEDIA":         SourceMedia,
		"GENRE":               Genre,
		"MOOD":                Mood,
		"RELEASECOUNTRY":      Country,
		"COMMENT":             Comments,
		"COPYRIGHT":           CopyrightMessage,
	},
}

// ProcessTags обрабатывает переданные теги, обновляя метаданные трека, альбома, релиза.
// Необработанные теги возвращаются функцией обратно.
func ProcessTags(tags map[TagKey]string, r *md.Release, t *md.Track) error {
	var err error
	for k, v := range tags {
		switch k {
		// --- Titles ---
		case AlbumTitle:
			r.Title = v
		// DiscSetSubtitle
		// ContentGroup
		case TrackTitle:
			t.Title = v
		case TrackSubtitle:
			if t.Title != "" {
				t.Title = md.ComplexTitle(t.Title, v)
			}
		// Version
		// --- People & Organizations ---
		case AlbumArtist:
			r.Actors.AddActorEntry(v)
		case TrackArtist, InvolvedPeople:
			parseAndAddActors(v, t)
		case Arranger:
			t.Record.Actors.AddRole(v, "arranger")
		case AuthorWriter, Writer:
			t.Composition.Actors.AddRole(v, "writer")
		case Composer:
			t.Composition.Actors.AddRole(v, "composer")
		case Conductor:
			t.Record.Actors.AddRole(v, "conductor")
		case Engineer:
			t.Record.Actors.AddRole(v, "engineer")
		case Ensemble:
			t.Record.Actors.AddRole(v, "ensemble")
		case Lyricist:
			t.Composition.Actors.AddRole(v, "lyricist")
		case MixDJ:
			t.Record.Actors.AddRole(v, "mix-DJ")
		case MixEngineer:
			t.Record.Actors.AddRole(v, "mix-engineer")
		// MusicianCredits
		// Organisation
		// OriginalArtist
		case Performer:
			r.Actors.AddRole(v, "performer")
		case Producer:
			t.Record.Actors.AddRole(v, "producer")
		case Publisher, Label:
			setLabels(v, r)
		case RemixedBy:
			t.Record.Actors.AddRole(v, "remixer")
		case Soloists:
			t.Record.Actors.AddRole(v, "soloist")
		// --- Counts & Indexes ---
		case DiscNumber:
			err = setTrackDiscNumber(v, r, t)
		case DiscTotal:
			r.TotalDiscs = stringutils.NaiveStringToInt(v)
		case TrackTotal:
			r.TotalTracks = stringutils.NaiveStringToInt(v)
		case TrackNumber:
			setTrackPositionAndTotalTracks(v, r, t)
		// PartNumber
		case Length:
			parseAndSetTrackDuration(v, t)
		// --- Dates ---
		case ReleaseDate:
			parseAndSetYearFromDate(v, r)
		case OriginalReleaseDate:
			parseAndSetOriginalYearFromDate(v, r)
		case Year:
			parseAndSetYears(v, r)
		case RecordingDates:
			t.AddComment(fmt.Sprintf("Recording: %s", v))
		// --- Identifiers ---
		case DiscID:
			setDiscID(tags, r, t)
		case ISRC:
			t.IDs["isrc"] = v
		case Barcode:
			r.IDs["barcode"] = v
		case CatalogueNumber, LabelNumber:
			setCatno(v, r)
		case UPC:
			r.IDs["upc"] = v
		case AccurateRipDiscID:
			r.IDs["accuraterip"] = v
		case DiscogsReleaseID:
			r.IDs["discogs"] = v
		case MusicbrainzAlbumID:
			r.IDs["musicbrainz"] = v
		case RutrackerID:
			r.IDs["rutracker"] = v
		// --- Flags ---
		case Compilation:
			r.ReleaseRepeat = md.ReleaseRepeatCompilation
		// --- Ripping & Encoding ---
		// FileType
		case MediaType:
			parseAndAddDiscFormat(v, r, t)
		// SourceMedia
		// Source
		// --- URLs ---
		// AudioSourceWebpageURL
		// CommercialInformationURL
		// TrackArtistWebPageURL
		// --- Style ---
		case Genre, Style:
			t.Record.Genres = append(t.Record.Genres, v)
		case Mood:
			setMood(v, t)
		// --- Miscellaneous ---
		case Country:
			parseAndAddCountries(v, r)
		case Comments, Description:
			t.AddComment(v)
		case CopyrightMessage:
			setCopyright(v, r)
		case SyncedLyrics:
			t.SetLyrics(v, true)
		case UnsyncedLyrics:
			t.SetLyrics(v, false)
		case Language:
			t.SetLyricsLanguage(v)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// ----- Compound processing -----

func setDiscID(tags map[TagKey]string, r *md.Release, t *md.Track) {
	var pos string
	if t.Position != "" {
		pos = t.Position
	} else if tn, ok := tags[TrackNumber]; ok {
		pos = tn
	}
	if pos != "" {
		if t.Disc() == nil {
			t.LinkWithDisc(r.Disc(md.DiscNumberByTrackPos(pos)))
		}
		t.Disc().IDs["discid"] = tags[DiscID]
	}
}

func setTrackPositionAndTotalTracks(v string, r *md.Release, t *md.Track) {
	flds := strings.Split(v, "/")
	fldsLen := len(flds)
	if 1 <= fldsLen && fldsLen <= 2 {
		t.SetPosition(flds[0])
	}
	if fldsLen == 2 {
		r.TotalTracks = stringutils.NaiveStringToInt(flds[1])
	}
}

func parseAndAddDiscFormat(v string, r *md.Release, t *md.Track) {
	if t.Disc() == nil || t.Disc().Number == 0 {
		return
	}
	if media := md.DecodeMedia(v); media != 0 {
		discInd := t.Disc().Number - 1
		r.Discs[discInd].Format = &md.DiscFormat{Media: media}
	}
}

// ----- Release processing -----

// case "RELEASECOUNTRY", "DISCOGS_COUNTRY", "COUNTRY":
func parseAndAddCountries(v string, r *md.Release) {
	r.Country = v
}

// Разбор числового значения даты
func parseAndSetYears(yearString string, r *md.Release) {
	flds := stringutils.SplitIntoRegularFields(yearString)
	fldLen := len(flds)
	if fldLen == 0 {
		return
	} else if fldLen == 1 {
		r.Year = stringutils.NaiveStringToInt(flds[len(flds)-1])
	} else {
		r.Original.Year = stringutils.NaiveStringToInt(flds[0])
		r.Year = stringutils.NaiveStringToInt(flds[len(flds)-1])
	}
}

// Разбор timestamp ISO 8601 (yyyy-MM-ddTHH:mm:ss) или ее подстроки для Release.Year.
// TODO: вынести разбор в отдельную функцию.
func parseAndSetYearFromDate(v string, r *md.Release) {
	r.Year = stringutils.NaiveStringToInt(v)
}

// Разбор timestamp ISO 8601 (yyyy-MM-ddTHH:mm:ss) или ее подстроки для Release.Album.Year.
// TODO: вынести разбор в отдельную функцию.
func parseAndSetOriginalYearFromDate(v string, r *md.Release) {
	r.Original.Year = stringutils.NaiveStringToInt(v)
}

func setCopyright(cr string, r *md.Release) {
	if r.Publishing == nil { // отдавать приоритет тегу Labels, а не Copyright
		setLabels(cr, r)
	}
}

func setLabels(v string, r *md.Release) {
	if r.Publishing == nil {
		r.Publishing = append(r.Publishing, &md.Publishing{})
	}
	r.Publishing[0].Name = v
}

func setCatno(v string, r *md.Release) {
	if r.Publishing == nil {
		r.Publishing = append(r.Publishing, &md.Publishing{})
	}
	r.Publishing[0].Catno = v
}

// ----- Track processing -----

// Possible format is a list of {soloists,conductor,orchestra}, separated with ';'.
func parseAndAddActors(names string, track *md.Track) {
	stringutils.PanicIfNonUtf8(names)
	for _, name := range stringutils.SplitIntoRegularFieldsWithDelimiters(names, []rune{';'}) {
		flds := stringutils.SplitIntoRegularFieldsWithDelimiters(name, []rune{'-', ',', '(', ')'})
		if len(flds) > 1 {
			for _, role := range flds[1:] {
				track.Record.Actors.AddRole(strings.TrimSpace(flds[0]), strings.TrimSpace(role))
			}
		} else {
			track.Record.Actors.AddActorEntry(name)
		}
	}
}

// Обработка строк "hh:mm:ss" для записи длительности трека в миллисекундах.
func parseAndSetTrackDuration(v string, t *md.Track) {
	t.Duration = intutils.NewDurationFromString(v)
}

func setMood(moods string, track *md.Track) {
	for _, moodName := range stringutils.SplitIntoRegularFields(moods) {
		track.Record.Moods = append(track.Record.Moods, md.MoodFromName(moodName))
	}
}

func setTrackDiscNumber(v string, r *md.Release, t *md.Track) error {
	if t.Position != "" {
		dn, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		if md.DiscNumberByTrackPos(t.Position) != dn {
			return errors.New("Incorrect disc number")
		}
		t.LinkWithDisc(r.Disc(dn))
	}
	return nil
}
