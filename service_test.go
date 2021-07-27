package mdreader

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	md "github.com/ytsiuryn/ds-audiomd"
	srv "github.com/ytsiuryn/ds-microservice"
	"github.com/ytsiuryn/go-collection"
)

type MdreaderTestSuite struct {
	suite.Suite
	cl *srv.RPCClient
}

func (suite *MdreaderTestSuite) SetupSuite() {
	suite.startTestService()
	suite.cl = srv.NewRPCClient()
}

func (suite *MdreaderTestSuite) TearDownSuite() {
	suite.cl.Close()
}

func (suite *MdreaderTestSuite) TestBaseServiceCommands() {
	correlationID, data, err := srv.CreateCmdRequest("ping")
	require.NoError(suite.T(), err)
	suite.cl.Request(ServiceName, correlationID, data)
	respData := suite.cl.Result(correlationID)
	suite.Empty(respData)

	correlationID, data, err = srv.CreateCmdRequest("x")
	require.NoError(suite.T(), err)
	suite.cl.Request(ServiceName, correlationID, data)
	vInfo, err := srv.ParseErrorAnswer(suite.cl.Result(correlationID))
	require.NoError(suite.T(), err)
	// {"error": "Unknown command: x", "context": "Message dispatcher"}
	suite.Equal(vInfo.Error, "Unknown command: x")
}

func (suite *MdreaderTestSuite) TestDirRequest() {
	for _, subdir := range []string{"dsf", "flac", "mp3", "wavpack"} {
		correlationID, data, _ := CreateDirRequest("testdata/" + subdir)
		suite.cl.Request(ServiceName, correlationID, data)

		resp, _ := ParseDirAnswer(suite.cl.Result(correlationID))
		checkResp(&suite.Suite, resp.Assumption)
	}
}

// func (suite *MdreaderTestSuite) TestRepoDir() {
// 	home, err := os.UserHomeDir()
// 	require.NoError(suite.T(), err)
// 	correlationID, data, err := CreateDirRequest(filepath.Join(home, "Downloads/TEST/!AFTER FOREVER [2005] [SACD] Remagine [DSD-CD-FLAC]"))
// 	require.NoError(suite.T(), err)
// 	suite.cl.Request(ServiceName, correlationID, data)
// 	data = suite.cl.Result(correlationID)
// 	err = ioutil.WriteFile(filepath.Join(home, "Downloads/TEST/test.json"), data, 0644)
// 	require.NoError(suite.T(), err)
// }

func checkResp(suite *suite.Suite, assumption *md.Assumption) {
	r := assumption.Release
	tr := assumption.Release.Tracks[0]
	suite.Equal(r.Title, "test_album_title")
	suite.Equal(r.ActorRoles.First(), "test_performer")
	suite.Equal(r.Discs[0].Number, 1)
	suite.Equal(r.TotalTracks, 10)
	suite.Equal(r.Tracks[0].Position, "03")
	suite.Equal(tr.Composition.ActorRoles.First(), "test_composer")
	suite.Equal(tr.Record.Actors.First(), "test_track_artist")
	suite.Equal(tr.Record.Genres[0], "test_genre")
	suite.Equal(tr.Title, "test_track_title")
	ext := filepath.Ext(tr.FileName)
	if !collection.ContainsStr(ext, []string{".mp3", ".wv"}) { // TODO
		suite.Equal(int64(tr.Duration), int64(500))
	}
	if !collection.ContainsStr(ext, []string{".dsf", ".wv"}) { // TODO
		suite.Equal(tr.AudioInfo.Samplerate, 44100)
		suite.Equal(tr.AudioInfo.SampleSize, 16)
	}
	if ext != ".wv" { // TODO
		suite.Equal(tr.AudioInfo.Channels, 1)
		suite.Equal(assumption.Pictures[0].PictureMetadata.MimeType, "image/jpeg")
		suite.Equal(assumption.Pictures[0].PictType, md.PictTypeCoverFront)
	}
	suite.Equal(r.Publishing[0].Name, "test_label")
	suite.Equal(r.Publishing[0].Catno, "test_catno")
	suite.Equal(r.Country, "test_country")
	suite.Equal(r.Year, 2000)
	suite.Equal(r.Notes, "test_notes")
}

func (suite *MdreaderTestSuite) startTestService() {
	testService := New()
	msgs := testService.ConnectToMessageBroker("amqp://guest:guest@localhost:5672/")
	// testService.Log.SetLevel(log.DebugLevel)
	go testService.Start(msgs)
}

func TestMdreaderService(t *testing.T) {
	suite.Run(t, new(MdreaderTestSuite))
}
