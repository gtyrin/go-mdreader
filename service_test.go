package mdreader

import (
	"context"
	"path/filepath"
	"sync"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	md "github.com/ytsiuryn/ds-audiomd"
	srv "github.com/ytsiuryn/ds-microservice"
	"github.com/ytsiuryn/go-collection"
)

var mut sync.Mutex
var testService *AudioMdReader

func TestBaseServiceCommands(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startTestService(ctx)

	cl := srv.NewRPCClient()
	defer cl.Close()

	correlationID, data, err := srv.CreateCmdRequest("ping")
	require.NoError(t, err)
	cl.Request(ServiceName, correlationID, data)
	respData := cl.Result(correlationID)
	assert.Empty(t, respData)

	correlationID, data, err = srv.CreateCmdRequest("x")
	require.NoError(t, err)
	cl.Request(ServiceName, correlationID, data)
	vInfo, err := srv.ParseErrorAnswer(cl.Result(correlationID))
	require.NoError(t, err)
	// {"error": "Unknown command: x", "context": "Message dispatcher"}
	assert.Equal(t, vInfo.Error, "Unknown command: x")
}

func TestDirRequest(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startTestService(ctx)

	cl := srv.NewRPCClient()
	defer cl.Close()

	for _, subdir := range []string{"dsf", "flac", "mp3", "wavpack"} {
		correlationID, data, _ := CreateDirRequest("testdata/" + subdir)
		cl.Request(ServiceName, correlationID, data)

		suggestion, _ := ParseDirAnswer(cl.Result(correlationID))
		checkResp(t, suggestion)
	}
}

// func TestRepoDir(t *testing.T) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	startTestService(ctx)

// 	cl := srv.NewRPCClient()
// 	defer cl.Close()
// 	correlationID, data, err := CreateDirRequest("/home/me/Downloads/TEST/!AFTER FOREVER [2005] [SACD] Remagine [DSD-CD-FLAC]")
// 	require.NoError(t, err)
// 	cl.Request(ServiceName, correlationID, data)

// 	data = cl.Result(correlationID)
// 	log.Fatal(string(data))
// }

func checkResp(t *testing.T, suggestion *md.Suggestion) {
	r := suggestion.Release
	tr := suggestion.Release.Tracks[0]
	assert.Equal(t, r.Title, "test_album_title")
	assert.Equal(t, r.ActorRoles.First(), "test_performer")
	assert.Equal(t, r.Discs[0].Number, 1)
	assert.Equal(t, r.TotalTracks, 10)
	assert.Equal(t, r.Tracks[0].Position, "03")
	assert.Equal(t, tr.Composition.ActorRoles.First(), "test_composer")
	assert.Equal(t, tr.Record.Actors.First(), "test_track_artist")
	assert.Equal(t, tr.Record.Genres[0], "test_genre")
	assert.Equal(t, tr.Title, "test_track_title")
	ext := filepath.Ext(tr.FileName)
	if !collection.ContainsStr(ext, []string{".mp3", ".wv"}) { // TODO
		assert.Equal(t, tr.Duration, 500)
	}
	if !collection.ContainsStr(ext, []string{".dsf", ".wv"}) { // TODO
		assert.Equal(t, tr.AudioInfo.Samplerate, 44100)
		assert.Equal(t, tr.AudioInfo.SampleSize, 16)
	}
	if ext != ".wv" { // TODO
		assert.Equal(t, tr.AudioInfo.Channels, 1)
		assert.Equal(t, r.Pictures[0].PictureMetadata.MimeType, "image/jpeg")
		assert.Equal(t, r.Pictures[0].PictType, md.PictTypeCoverFront)
	}
	assert.Equal(t, r.Publishing[0].Name, "test_label")
	assert.Equal(t, r.Publishing[0].Catno, "test_catno")
	assert.Equal(t, r.Country, "test_country")
	assert.Equal(t, r.Year, 2000)
	assert.Equal(t, r.Notes, "test_notes")
}

func startTestService(ctx context.Context) {
	mut.Lock()
	defer mut.Unlock()
	if testService == nil {
		testService = New()
		msgs := testService.ConnectToMessageBroker("amqp://guest:guest@localhost:5672/")
		testService.Log.SetLevel(log.DebugLevel)
		go testService.Start(msgs)
	}
}
