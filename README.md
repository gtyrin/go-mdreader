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

*Пример использования команд приведен в тестовом клиенте в [mdreader.py](https://github.com/ytsiuryn/ds-mdreader/blob/main/mdreader.py)*.

Пример запуска микросервиса:
---
```go
    package main

    import (
	    "flag"
	    "fmt"

	    log "github.com/sirupsen/logrus"

	    mdreader "github.com/ytsiuryn/ds-mdreader"
	    srv "github.com/ytsiuryn/ds-microservice"
    )

    func main() {
	    connstr := flag.String(
		    "msg-server",
		    "amqp://guest:guest@localhost:5672/",
		    "Message server connection string")

		product := flag.Bool(
			"product",
			false,
			"product-режим запуска сервиса")

		flag.Parse()

	    log.Info(fmt.Sprintf("%s starting..", mdreader.ServiceName))

	    reader := mdreader.New()

		msgs := reader.ConnectToMessageBroker(*connstr)

		if *product {
			reader.Log.SetLevel(log.InfoLevel)
		} else {
			reader.Log.SetLevel(log.DebugLevel)
		}

		reader.Start(msgs)
	}
```

Пример клиента (Python тест):
---
См. файл [mdreader.py](https://github.com/ytsiuryn/ds-mdreader/blob/main/mdreader.py)