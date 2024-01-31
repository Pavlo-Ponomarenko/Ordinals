package requests

import (
	"encoding/json"
	"net/http"
)

type CreateInscriptionRequest struct {
	Hash    string `json:"hash"`
	Output  uint32 `json:"output"`
	Key     string `json:"key"`
	Content string `json:"content"`
	Address string `json:"address"`
}

func NewCreateInscriptionRequest(r *http.Request) (*CreateInscriptionRequest, error) {
	request := new(CreateInscriptionRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		return nil, err
	}
	return request, nil
}
