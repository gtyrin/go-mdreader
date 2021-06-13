package file

import (
	"testing"

	md "github.com/ytsiuryn/ds-audiomd"
	collection "github.com/ytsiuryn/go-collection"
)

var tagMapTestData = map[TagKey]string{
	DiscID: "1234",
}

func TestSetDiscID(t *testing.T) {
	r := md.NewRelease()
	tr := md.NewTrack()
	setDiscID(tagMapTestData, r, tr)
	if len(r.Discs) > 0 {
		t.Fail()
	}
	tr.Position = "1"
	setDiscID(tagMapTestData, r, tr)
	if len(r.Discs) != 1 || r.Discs[0].Number != 1 {
		t.Fail()
	}
	if !r.Disc(1).IDs.Exists("discid") {
		t.Fail()
	}
}

func TestSetTrackPositionAndTotalTracks(t *testing.T) {
	r := md.NewRelease()
	tr := md.NewTrack()
	setTrackPositionAndTotalTracks("2", r, tr)
	if tr.Position != "02" || r.TotalTracks != 0 {
		t.Fail()
	}
	setTrackPositionAndTotalTracks("2/10", r, tr)
	if tr.Position != "02" || r.TotalTracks != 10 {
		t.Fail()
	}
}

func TestParseAndAddDiscFormat(t *testing.T) {
	r := md.NewRelease()
	tr := md.NewTrack()
	parseAndAddDiscFormat("LP", r, tr)
	if tr.Disc() != nil && tr.Disc().Number != 0 {
		t.Fail()
	}
	d := r.Disc(1)
	tr.LinkWithDisc(d)
	parseAndAddDiscFormat("LP", r, tr)
	if d.Format.Media != md.MediaLP {
		t.Fail()
	}
}

// TODO: реализовать!
func TestParseAndAddLabels(t *testing.T) {

}

func TestParseAndSetYears(t *testing.T) {
	r := md.NewRelease()
	parseAndSetYears("", r)
	if r.Year != 0 || r.Original.Year != 0 {
		t.Fail()
	}
	parseAndSetYears("2005", r)
	if r.Year != 2005 {
		t.Fail()
	}
	parseAndSetYears("1961,1962/2005", r)
	if r.Year != 2005 || r.Original.Year != 1961 {
		t.Fail()
	}
}

// TODO: реализовать!
func TestParseCopyrightAndAddLabels(t *testing.T) {

}

func TestSetCatno(t *testing.T) {
	r := md.NewRelease()
	setCatno("12345", r)
	if len(r.Publishing) != 1 && r.Publishing[0].Catno != "12345" {
		t.Fail()
	}
}

func TestParseAndAddActors(t *testing.T) {
	tr := md.NewTrack()
	parseAndAddActors("Karajan, conductor; BPO", tr)
	karajan := tr.Record.Actors.Actor("Karajan")
	if len(*tr.Record.Actors) != 2 || karajan == nil || tr.Record.Actors.Actor("BPO") == nil {
		t.Fail()
	}
	if !collection.ContainsStr("conductor", karajan.Roles) {
		t.Fail()
	}
}

func TestSetTrackDiscNumber(t *testing.T) {
	r := md.NewRelease()
	tr := md.NewTrack()
	if setTrackDiscNumber("1", r, tr) != nil {
		t.Fail()
	}
	if len(r.Discs) != 0 || tr.Disc() != nil {
		t.Fail()
	}
	tr.Position = "2"
	if setTrackDiscNumber("1", r, tr) != nil {
		t.Fail()
	}
	if tr.Disc().Number != 1 {
		t.Fail()
	}
	tr.Position = "2.2"
	if setTrackDiscNumber("1", r, tr) == nil || setTrackDiscNumber("Vol.1", r, tr) == nil {
		t.Fail()
	}
}
