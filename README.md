# go-mdreader #

Микросервис с парсером метаданных аудиофайлов. Обмен сообщениями реализован с использованием
[RabbitMQ](https://www.rabbitmq.com).

## Поддержка аудиоформатов:
- mp3 (id3v1/id3v2)
- flac (id3v2/vorbis comments)
- dsf (id3v2)
- wavpack (id3v2/apev2; без аудиосвойств треков)

## Пример использования:

    package main

    import (
	    "flag"
	    "fmt"

	    log "github.com/sirupsen/logrus"

	    mdreader "github.com/ytsiuryn/ds-mdreader"
	    srv "github.com/ytsiuryn/ds-service"
    )

    func main() {
	    connstr := flag.String(
		    "msg-server",
		    "amqp://guest:guest@localhost:5672/",
		    "Message server connection string")
	    flag.Parse()

	    log.Info(
		    fmt.Sprintf("%s starting..", mdreader.ServiceName))

	    cl, err := mdreader.NewAudioMetadataReader(*connstr)
	    srv.FailOnError(err, "Failed to create metadata reader")

	    defer cl.Close()

	    cl.Dispatch(cl)
    }
