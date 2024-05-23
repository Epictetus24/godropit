package dropfmt

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/kopp0ut/go-util/pkg/enc"
)

// Struct for Handling Encrypted Shellcode
type DropFmt struct {
	Buf    []byte //
	EncBuf []byte
	DecBuf []byte
	Key    []byte
}

// Convert .EncBuf to base64
func (d *DropFmt) ToB64() string {
	return base64.StdEncoding.EncodeToString(d.EncBuf)

}

// Convert .Key to base64 string
func (d *DropFmt) KeyB64() string {
	return base64.StdEncoding.EncodeToString(d.Key)

}

// converts to hex string
func (d *DropFmt) ToHex() string {
	return hex.EncodeToString(d.EncBuf)

}

/*
func (d *DropFmt) SGN(arch int) error {
	var err error
	// Create a new SGN encoder
	encoder := sgn.NewEncoder()
	// Set the proper architecture
	encoder.SetArchitecture(arch)
	sc := d.Buf
	// Encode the binary
	d.Buf, err = encoder.Encode(sc)
	if err != nil {
		return err
	}
	return nil
}
*/
// Encrypt Buf with AES to EncBuf
func (d *DropFmt) AESEncrypt() (string, error) {
	var err error

	if d.Key == nil {
		d.NewAESKey()
	}

	d.EncBuf, err = enc.AESCBCEncrypt(d.Key, d.Buf)
	if err != nil {
		return "", err
	}

	encstr := base64.StdEncoding.EncodeToString(d.DecBuf)
	return encstr, nil

}

func (d *DropFmt) AESDecrypt() (string, string, error) {
	var err error

	d.DecBuf, err = enc.AESCBCDecrypt(d.Key, d.Buf)
	if err != nil {
		return "", "", err
	}

	encstr := base64.StdEncoding.EncodeToString(d.EncBuf)
	return hex.EncodeToString(d.Key), encstr, nil

}

// Generate an AES Key for use with encryption/decryption.
func (d *DropFmt) NewAESKey() ([]byte, error) {
	var err error
	d.Key, err = enc.GenKey(32)
	return d.Key, err

}
