package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
	md "github.com/ytsiuryn/ds-audiomd"
)

var tagMapTestData = map[TagKey]string{
	DiscID: "1234",
}

func TestSetDiscID(t *testing.T) {
	r := md.NewRelease()
	tr := md.NewTrack()
	setDiscID(tagMapTestData, r, tr)
	assert.Empty(t, r.Discs)
	tr.Position = "1"
	setDiscID(tagMapTestData, r, tr)
	assert.Len(t, r.Discs, 1)
	assert.Equal(t, r.Discs[0].Number, 1)
	_, ok := r.Disc(1).IDs[md.ID]
	assert.True(t, ok)
}

func TestSetTrackPositionAndTotalTracks(t *testing.T) {
	r := md.NewRelease()
	tr := md.NewTrack()
	setTrackPositionAndTotalTracks("2", r, tr)
	assert.Equal(t, tr.Position, "02")
	assert.Zero(t, r.TotalTracks)
	setTrackPositionAndTotalTracks("2/10", r, tr)
	assert.Equal(t, tr.Position, "02")
	assert.Equal(t, r.TotalTracks, 10)
}

func TestParseAndAddDiscFormat(t *testing.T) {
	r := md.NewRelease()
	tr := md.NewTrack()
	parseAndAddDiscFormat("LP", r, tr)
	assert.Nil(t, tr.Disc())
	d := r.Disc(1)
	tr.LinkWithDisc(d)
	parseAndAddDiscFormat("LP", r, tr)
	assert.Equal(t, d.Format.Media, md.MediaLP)
}

// TODO: реализовать!
func TestParseAndAddLabels(t *testing.T) {

}

func TestParseAndSetYears(t *testing.T) {
	r := md.NewRelease()
	parseAndSetYears("", r)
	assert.Zero(t, r.Year)
	assert.Zero(t, r.Original.Year)
	parseAndSetYears("2005", r)
	assert.Equal(t, r.Year, 2005)
	parseAndSetYears("1961,1962/2005", r)
	assert.Equal(t, r.Year, 2005)
	assert.Equal(t, r.Original.Year, 1961)
}

// TODO: реализовать!
func TestParseCopyrightAndAddLabels(t *testing.T) {

}

func TestSetCatno(t *testing.T) {
	r := md.NewRelease()
	setCatno("12345", r)
	assert.Len(t, r.Publishing, 1)
	assert.Equal(t, r.Publishing[0].Catno, "12345")
}

// func TestParseAndAddActors(t *testing.T) {
// 	tr := md.NewTrack()
// 	parseAndAddActors("Karajan, conductor; BPO", tr)
// 	karajan := tr.Record.Actors.Actor("Karajan")
// 	if len(*tr.Record.Actors) != 2 || karajan == nil || tr.Record.Actors.Actor("BPO") == nil {
// 		t.Fail()
// 	}
// 	if !collection.ContainsStr("conductor", karajan.Roles) {
// 		t.Fail()
// 	}
// }

func TestSetTrackDiscNumber(t *testing.T) {
	r := md.NewRelease()
	tr := md.NewTrack()
	assert.Nil(t, setTrackDiscNumber("1", r, tr))
	assert.Empty(t, r.Discs)
	assert.Nil(t, tr.Disc())
	tr.Position = "2"
	assert.Nil(t, setTrackDiscNumber("1", r, tr))
	assert.Equal(t, tr.Disc().Number, 1)
	tr.Position = "2.2"
	assert.NotNil(t, setTrackDiscNumber("1", r, tr))
	assert.NotNil(t, setTrackDiscNumber("Vol.1", r, tr))
}
