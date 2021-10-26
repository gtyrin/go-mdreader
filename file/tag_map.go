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
	binary "github.com/ytsiuryn/go-binary"
	intutils "github.com/ytsiuryn/go-intutils"
	stringutils "github.com/ytsiuryn/go-stringutils"
)

// ProcessTags обрабатывает переданные теги, обновляя метаданные трека, альбома, релиза.
// Необработанные теги возвращаются функцией обратно.
func ProcessTags(tags map[TagKey]TagValue, r *md.Release, t *md.Track) error {
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
		case AlbumArtist, Performer:
			r.ActorRoles.Add(v, "performer")
		case TrackArtist, InvolvedPeople:
			parseAndAddActors(v, t)
		case Arranger:
			t.Record.ActorRoles.Add(v, "arranger")
		case AuthorWriter, Writer:
			t.Composition.ActorRoles.Add(v, "writer")
		case Composer:
			t.Composition.ActorRoles.Add(v, "composer")
		case Conductor:
			t.Record.ActorRoles.Add(v, "conductor")
		case Engineer:
			t.Record.ActorRoles.Add(v, "engineer")
		case Ensemble:
			t.Record.ActorRoles.Add(v, "ensemble")
		case Lyricist:
			t.Composition.ActorRoles.Add(v, "lyricist")
		case MixDJ:
			t.Record.ActorRoles.Add(v, "mix-DJ")
		case MixEngineer:
			t.Record.ActorRoles.Add(v, "mix-engineer")
		// MusicianCredits
		// Organisation
		// OriginalArtist
		case Producer:
			t.Record.ActorRoles.Add(v, "producer")
		case Publisher, Label:
			setLabels(v, r)
		case RemixedBy:
			t.Record.ActorRoles.Add(v, "remixer")
		case Soloists:
			t.Record.ActorRoles.Add(v, "soloist")
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
			setBarcode(v, r)
		case CatalogueNumber, LabelNumber:
			setCatno(v, r)
		case UPC:
			setBarcode(v, r) // ведь UPC=barcode?
		case AccurateRipDiscID:
			r.IDs[md.AccurateRip] = v
		case DiscogsReleaseID:
			r.IDs[md.DiscogsReleaseID] = v
		case MusicbrainzAlbumID:
			r.IDs[md.MusicbrainzAlbumID] = v
		case RutrackerID:
			r.IDs[md.Rutracker] = v
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

func setDiscID(tags map[TagKey]TagValue, r *md.Release, t *md.Track) {
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
		t.Disc().IDs[md.ID] = tags[DiscID]
	}
}

func setTrackPositionAndTotalTracks(trackStr string, r *md.Release, t *md.Track) {
	flds := strings.Split(trackStr, "/")
	fldsLen := len(flds)
	if 1 <= fldsLen && fldsLen <= 2 {
		t.SetPosition(flds[0])
	}
	if fldsLen == 2 {
		r.TotalTracks = stringutils.NaiveStringToInt(flds[1])
	}
}

func parseAndAddDiscFormat(mediaStr string, r *md.Release, t *md.Track) {
	if t.Disc() == nil || t.Disc().Number == 0 {
		return
	}
	if media := md.DecodeMedia(mediaStr); media != 0 {
		discInd := t.Disc().Number - 1
		r.Discs[discInd].Format = &md.DiscFormat{Media: media}
	}
}

// ----- Release processing -----

// case "RELEASECOUNTRY", "DISCOGS_COUNTRY", "COUNTRY":
func parseAndAddCountries(country string, r *md.Release) {
	r.Country = country
}

// Разбор числового значения даты
func parseAndSetYears(yearStr string, r *md.Release) {
	flds := stringutils.SplitIntoRegularFields(yearStr)
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
func parseAndSetYearFromDate(dateStr string, r *md.Release) {
	r.Year = stringutils.NaiveStringToInt(dateStr)
}

// Разбор timestamp ISO 8601 (yyyy-MM-ddTHH:mm:ss) или ее подстроки для Release.Album.Year.
// TODO: вынести разбор в отдельную функцию.
func parseAndSetOriginalYearFromDate(dateStr string, r *md.Release) {
	r.Original.Year = stringutils.NaiveStringToInt(dateStr)
}

func setCopyright(cr string, r *md.Release) {
	if r.Publishing == nil { // отдавать приоритет тегу Labels, а не Copyright
		setLabels(cr, r)
	}
}

func setLabels(label string, r *md.Release) {
	if r.Publishing == nil {
		r.Publishing = append(r.Publishing, &md.Publishing{})
	}
	r.Publishing[0].Name = label
}

func setCatno(catno string, r *md.Release) {
	if r.Publishing == nil {
		r.Publishing = append(r.Publishing, &md.Publishing{})
	}
	r.Publishing[0].Catno = catno
}

func setBarcode(barcode string, r *md.Release) {
	if r.Publishing == nil {
		r.Publishing = append(r.Publishing, &md.Publishing{})
	}
	r.Publishing[0].IDs[md.Barcode] = barcode
}

// ----- Track processing -----

// Possible format is a list of {soloists,conductor,orchestra}, separated with ';'.
func parseAndAddActors(names string, track *md.Track) {
	binary.PanicIfNonUtf8(names)
	for _, name := range stringutils.SplitIntoRegularFieldsWithDelimiters(names, []rune{';'}) {
		flds := stringutils.SplitIntoRegularFieldsWithDelimiters(name, []rune{'-', ',', '(', ')'})
		if len(flds) > 1 {
			for _, role := range flds[1:] {
				track.Record.ActorRoles.Add(strings.TrimSpace(flds[0]), strings.TrimSpace(role))
			}
		} else {
			track.Record.Actors.Add(name, 0, "")
		}
	}
}

// Обработка строк "hh:mm:ss" для записи длительности трека в миллисекундах.
func parseAndSetTrackDuration(durationStr string, t *md.Track) {
	t.Duration = intutils.NewDurationFromString(durationStr)
}

func setMood(moods string, track *md.Track) {
	for _, moodName := range stringutils.SplitIntoRegularFields(moods) {
		track.Record.Moods = append(track.Record.Moods, md.MoodFromName(moodName))
	}
}

func setTrackDiscNumber(discNumStr string, r *md.Release, t *md.Track) error {
	if t.Position != "" {
		dn, err := strconv.Atoi(discNumStr)
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
