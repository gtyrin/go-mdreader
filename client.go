package mdreader

import (
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"

	md "github.com/ytsiuryn/ds-audiomd"
	srv "github.com/ytsiuryn/ds-microservice"
)

type AudioReaderRequest struct {
	Cmd  string `json:"cmd"`
	Path string `json:"path"`
}

type AudioReaderResponse struct {
	Assumption *md.Assumption    `json:"assumption,omitempty"`
	Error      srv.ErrorResponse `json:"error,omitempty"`
}

// CreateDirRequest формирует данные запроса поиска релиза по указанным метаданным.
func CreateDirRequest(dir string) (string, []byte, error) {
	correlationID, _ := uuid.NewV4()
	req := AudioReaderRequest{
		Cmd:  "release",
		Path: dir}
	data, err := json.Marshal(&req)
	if err != nil {
		return "", nil, err
	}
	return correlationID.String(), data, nil
}

// ParseDirAnswer разбирает ответ с предложением метаданных релиза.
func ParseDirAnswer(data []byte) (*AudioReaderResponse, error) {
	resp := AudioReaderResponse{}
	fmt.Println(string(data))
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
