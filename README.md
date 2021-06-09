# go-mdreader #

Микросервис с парсером метаданных аудиофайлов. Обмен сообщениями реализован с использованием
[RabbitMQ](https://www.rabbitmq.com).

Поддержка аудиоформатов:
- mp3
- flac
- dsf
- wavpack (без аудиосвойств треков)

Пример использования:

    package main

    import (
	    "flag"
	    "fmt"

	    log "github.com/sirupsen/logrus"

	    mdreader "github.com/gtyrin/go-mdreader"
	    srv "github.com/gtyrin/go-service"
    )

    func main() {
	    connstr := flag.String(
		    "msg-server",
		    "amqp://guest:guest@localhost:5672/",
		    "Message server connection string")
	    idle := flag.Bool(
		    "idle",
		    false,
		    "Free-running mode of the service to the message queue cleaning")
	    flag.Parse()

	    log.Info(
		    fmt.Sprintf("%s starting in %s mode..",
			    mdreader.ServiceName, srv.RunModeName(*idle)))

	    cl, err := mdreader.NewAudioMetadataReader(*connstr)
	    srv.FailOnError(err, "Failed to create metadata reader")

	    cl.Idle = *idle

	    defer cl.Close()

	    cl.Dispatch(cl)
    }
