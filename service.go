package mdreader

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/streadway/amqp"

	md "github.com/ytsiuryn/ds-audiomd"
	"github.com/ytsiuryn/ds-mdreader/file"
	srv "github.com/ytsiuryn/ds-microservice"
)

// Описание сервиса
const (
	ServiceSubsystem   = "audio"
	ServiceName        = "mdreader"
	ServiceDescription = "Audio Metadata Reader"
)

// AudioMdReader содержит состояние сервиса чтения метаданных.
type AudioMdReader struct {
	*srv.Service
}

type Assumption struct {
	ServiceName string `json:"service"`
	*md.Release `json:"release"`
}

// New создает объект нового клиента AudioMetadataReader.
func New() *AudioMdReader {
	return &AudioMdReader{
		Service: srv.NewService(ServiceName)}
}

// AnswerWithError заполняет структуру ответа информацией об ошибке.
func (ar *AudioMdReader) AnswerWithError(delivery *amqp.Delivery, err error, context string) {
	ar.LogOnError(err, context)
	req := &AudioReaderResponse{
		Error: &srv.ErrorResponse{
			Error:   err.Error(),
			Context: context,
		},
	}
	data, err := json.Marshal(req)
	if err != nil {
		srv.FailOnError(err, "Answer marshalling error")
	}
	ar.Answer(delivery, data)
}

// Start запускает Web Poller и цикл обработки взодящих запросов.
// Контролирует сигнал завершения цикла и последующего освобождения ресурсов микросервиса.
func (ar *AudioMdReader) Start(msgs <-chan amqp.Delivery) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		for delivery := range msgs {
			var req AudioReaderRequest
			if err := json.Unmarshal(delivery.Body, &req); err != nil {
				ar.AnswerWithError(&delivery, err, "Message dispatcher")
				continue
			}
			ar.logRequest(&req)
			ar.RunCmd(&req, &delivery)
		}
	}()

	ar.Log.Info("Awaiting RPC requests")
	<-c

	ar.cleanup()
}

func (ar *AudioMdReader) cleanup() {
	ar.Service.Cleanup()
}

// Отображение сведений о выполняемом запросе.
func (ar *AudioMdReader) logRequest(req *AudioReaderRequest) {
	if len(req.Path) > 0 {
		ar.Log.WithField("args", req.Path).Info(req.Cmd + "()")
	} else {
		ar.Log.Info(req.Cmd + "()")
	}
}

// RunCmdByName выполняет команды и возвращает результат клиенту в виде JSON-сообщения.
func (ar *AudioMdReader) RunCmd(req *AudioReaderRequest, delivery *amqp.Delivery) {
	var data []byte
	var err error

	switch req.Cmd {
	case "release":
		data, err = ar.releaseInfo(req)
	default:
		ar.Service.RunCmd(req.Cmd, delivery)
		return
	}

	if err != nil {
		ar.AnswerWithError(delivery, err, req.Cmd)
	} else {
		ar.Log.Debug(string(data))
		ar.Answer(delivery, data)
	}
}

func (ar *AudioMdReader) releaseInfo(req *AudioReaderRequest) ([]byte, error) {

	fileinfo, err := ioutil.ReadDir(req.Path)
	if err != nil {
		return nil, err
	}

	r := md.NewRelease()
	for _, fi := range fileinfo {
		fn := filepath.Join(req.Path, fi.Name())
		track, err := ar.readTrackFile(fn, r)
		if err != nil {
			return nil, err
		}
		if track == nil { // not audio file
			continue
		}
		r.Tracks = append(r.Tracks, track)
	}
	if len(r.Tracks) == 0 {
		return nil, err
	}

	assumption := md.NewAssumption(r)
	assumption.Optimize()

	return json.Marshal(AudioReaderResponse{Assumption: assumption})
}

func (ar *AudioMdReader) readTrackFile(fn string, r *md.Release) (*md.Track, error) {
	if reader := file.Reader(fn); reader != nil {
		f, err := os.OpenFile(fn, os.O_RDONLY, 0444)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		fi, err := f.Stat()
		if err != nil {
			return nil, err
		}
		track := md.NewTrack()
		track.FileInfo.FileName = fi.Name()
		track.FileInfo.ModTime = fi.ModTime().Unix()
		track.FileInfo.FileSize = fi.Size()
		if err := reader.TrackMetadata(f, r, track); err != nil {
			return nil, err
		}
		return track, nil
	}
	return nil, nil
}
