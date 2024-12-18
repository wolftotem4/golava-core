package session

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register(map[string]string{})
}

func marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	var base64Encoder = base64.NewEncoder(base64.StdEncoding, &buf)
	enc := gob.NewEncoder(base64Encoder)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	base64Encoder.Close()
	return buf.Bytes(), nil
}

func unmarshal(data []byte, v interface{}) error {
	var base64Decoder = base64.NewDecoder(base64.StdEncoding, bytes.NewReader(data))
	dec := gob.NewDecoder(base64Decoder)
	err := dec.Decode(v)
	if err != nil {
		return err
	}

	return nil
}
