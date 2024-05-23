package box

import (
	"encoding/base64"

	"github.com/kopp0ut/go-util/pkg/enc"
)

// Decrypt AES Buffer
func AESDecrypt(key string, buf string) ([]byte, error) {

	bitkey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	dat, err := base64.StdEncoding.DecodeString(buf)
	if err != nil {
		return nil, err
	}

	decbuf, err := enc.AESCBCDecrypt(bitkey, dat)
	if err != nil {
		return nil, err
	}

	return decbuf, nil

}
