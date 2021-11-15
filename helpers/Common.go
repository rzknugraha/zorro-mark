package helpers

import (
	"bytes"
	"encoding/json"
	"io"
)

// SetBody func
func SetBody(Body map[string]interface{}) (io.Reader, error) {
	body, err := json.Marshal(Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(body), err
}
