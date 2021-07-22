package mdreader

import (
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"

	md "github.com/ytsiuryn/ds-audiomd"
)

type AudioReaderRequest struct {
	Cmd  string `json:"cmd"`
	Path string `json:"path"`
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
func ParseDirAnswer(data []byte) (*md.Assumption, error) {
	assumption := md.Assumption{}
	fmt.Println(string(data))
	err := json.Unmarshal(data, &assumption)
	if err != nil {
		return nil, err
	}
	return &assumption, nil
}
