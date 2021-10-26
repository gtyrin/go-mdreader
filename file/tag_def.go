package file

// TagName ..
type TagName = string

// TagValue - строковое значение тега.
type TagValue = string

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
var SchemaTagToUniKey = map[TagScheme]map[TagName]TagKey{
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
