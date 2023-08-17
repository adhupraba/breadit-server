package utils

import (
	"encoding/json"
	"io"
)

func BodyParser(body io.ReadCloser, v any) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(v)

	if err != nil {
		return err
	}

	return nil
}
