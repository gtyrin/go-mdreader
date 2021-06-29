# go-mdreader #

Микросервис с парсером метаданных аудиофайлов. Обмен сообщениями с микросервисом реализован с использованием [RabbitMQ](https://www.rabbitmq.com).

Поддержка аудиоформатов:
---
- mp3 (id3v1/id3v2)
- flac (id3v2/vorbis comments)
- dsf (id3v2)
- wavpack (id3v2/apev2; без аудиосвойств треков)

Команды микросервиса:
---
|Команда|            Назначение                |
|-------|--------------------------------------|
|release|чтение метаданных альбома в каталоге  |
|ping   |проверка жизнеспособности микросервиса|
|info   |информация о микросервисе             |

*Пример использования команд приведен в тестовом клиенте в [mdreader.py](https://github.com/ytsiuryn/ds-mdreader/blob/main/mdreader.py)*.

Файл настроек (YAML):
---
|  Секция/параметр  |                         Назначение                       |
|-------------------|----------------------------------------------------------|
|product            |параметр-признак работы микросервиса в неотладочном режиме|

Пример запуска микросервиса:
---
```go
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

	    log.Info(fmt.Sprintf("%s starting..", mdreader.ServiceName))

	    cl, err := mdreader.NewAudioMetadataReader(*connstr)
	    srv.FailOnError(err, "Failed to create metadata reader")

	    defer cl.Close()

	    cl.Dispatch(cl)
	}
```

Пример клиента (Python тест):
---
См. файл [mdreader.py](https://github.com/ytsiuryn/ds-mdreader/blob/main/mdreader.py)