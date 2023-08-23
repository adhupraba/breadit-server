package utils

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/adhupraba/breadit-server/lib"
)

func BodyParser(body io.ReadCloser, v any) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(v)

	if err != nil {
		return fmt.Errorf("Unable to parse request")
	}

	err = lib.Validate.Struct(v)

	if err != nil {
		fmt.Println("validation err =>", err)
		return fmt.Errorf("Received invalid data")
	}

	return nil
}
