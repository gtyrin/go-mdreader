package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/gtyrin/go-audio/mdreader/file"
	md "github.com/gtyrin/go-audiomd"
	srv "github.com/gtyrin/go-service"
)

// Описание сервиса
const (
	ServiceSubsystem   = "audio"
	ServiceName        = "mdreader"
	ServiceDescription = "Audio Metadata Reader"
	ServiceVersion     = "1.0.11"
)

type config struct {
	Product bool `yaml:"product"`
}

// AudioMetadataReader содержит состояние сервиса чтения метаданных.
type AudioMetadataReader struct {
	*srv.Service
	conf *config
}

type assumption struct {
	ServiceName string `json:"service"`
	*md.Release `json:"release"`
}

// NewAudioMetadataReader создает объект нового клиента AudioMetadataReader.
func NewAudioMetadataReader(connstr string) (*AudioMetadataReader, error) {
	conf := config{}
	srv.ReadConfig("mdreader.yml", &conf)

	log.SetLevel(srv.LogLevel(conf.Product))

	cl := &AudioMetadataReader{}
	cl.conf = &conf
	cl.Service = srv.NewService()
	cl.ConnectToMessageBroker(connstr, ServiceName)

	return cl, nil
}

// Cleanup ..
func (ar *AudioMetadataReader) Cleanup() {
	ar.Service.Cleanup()
}

// RunCmdByName выполняет команды и возвращает результат клиенту в виде JSON-сообщения.
func (ar *AudioMetadataReader) RunCmdByName(cmd string, delivery *amqp.Delivery) {
	switch cmd {
	case "release":
		ar.releaseInfo(delivery)
	case "info":
		version := srv.Version{
			Subsystem:   ServiceSubsystem,
			Name:        ServiceName,
			Description: ServiceDescription,
			Version:     ServiceVersion,
		}
		go ar.Service.Info(delivery, &version)
	default:
		ar.Service.RunCommonCmd(cmd, delivery)
	}
}

func (ar *AudioMetadataReader) releaseInfo(delivery *amqp.Delivery) {
	if ar.Idle {
		res := assumption{
			ServiceName: ServiceName,
			Release:     md.NewRelease(),
		}
		assumptionJSON, err := json.Marshal(res)
		if err != nil {
			ar.ErrorResult(delivery, err, "Response")
			return
		}
		ar.Answer(delivery, assumptionJSON)
		return
	}
	// прием входного запроса
	var request srv.Request
	err := json.Unmarshal(delivery.Body, &request)
	if err != nil {
		ar.ErrorResult(delivery, err, "Request")
		return
	}
	dir, ok := request.Params["dir"]
	if !ok {
		ar.ErrorResult(delivery, errors.New("parameter 'dir' has absent"), "Request")
	}
	fileinfo, err := ioutil.ReadDir(dir)
	if err != nil {
		ar.ErrorResult(delivery, err, "Directory reading")
		return
	}
	r := md.NewRelease()
	for _, fi := range fileinfo {
		fn := filepath.Join(dir, fi.Name())
		track, err := ar.readTrackFile(fn, r)
		if err != nil {
			ar.ErrorResult(delivery, err, fmt.Sprintf("File %s parsing", fn))
			return
		}
		if track == nil { // not audio file
			continue
		}
		r.Tracks = append(r.Tracks, track)
	}
	r.Optimize()
	res := assumption{
		ServiceName: ServiceName,
		Release:     r,
	}
	// отправка ответа
	assumptionJSON, err := json.Marshal(res)
	if err != nil {
		ar.ErrorResult(delivery, err, "Response")
		return
	}
	if !ar.conf.Product {
		log.Println(string(assumptionJSON))
	}
	ar.Answer(delivery, assumptionJSON)
}

func (ar *AudioMetadataReader) readTrackFile(fn string, r *md.Release) (*md.Track, error) {
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
