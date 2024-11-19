package req

import (
	"arch/ikeppu/github.com/pkg/response"
	"net/http"
)

func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		response.Json(*w, err.Error(), 400)
		return nil, err
	}

	err = IsValid(body)
	if err != nil {
		response.Json(*w, err.Error(), 400)
		return nil, err
	}

	return &body, nil
}
