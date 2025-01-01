package hashing

import (
	"fmt"
	"testing"
)

func TestMd5Hasher_Make(t *testing.T) {
	h := Md5Hasher{}
	hash, err := h.Make("secret")
	if err != nil {
		t.Error(err)
	}
	if hash != "$md5$5ebe2294ecd0e0f08eab7690d2a6ee69" {
		t.Error("hash not match")
	}

	h.Callback = func(value string) (string, error) {
		return fmt.Sprintf("%s%s", value, "salt"), nil
	}
	hash, err = h.Make("secret")
	if err != nil {
		t.Error(err)
	}
	if hash != "$md5$99cd2e5a95d555ee7be3b038a4a84625" {
		t.Error("hash not match")
	}
}
