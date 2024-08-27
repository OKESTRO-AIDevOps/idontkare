package comm

import (
	"encoding/hex"
	"fmt"

	pkgutils "github.com/OKESTRO-AIDevOps/idontkare/pkg/utils"
)

func CommDataEncrypt(data []byte, key []byte) (string, error) {

	enc_b, err := pkgutils.EncryptWithSymmetricKey(key, data)

	if err != nil {

		return "", fmt.Errorf("failed to encrypt data: %s", err.Error())
	}

	enc_hex := hex.EncodeToString(enc_b)

	return enc_hex, nil

}

func CommDataDecrypt(data string, key []byte) ([]byte, error) {

	hex_b, err := hex.DecodeString(data)

	if err != nil {

		return nil, fmt.Errorf("failed to decode: %s", err.Error())

	}

	dec_b, err := pkgutils.DecryptWithSymmetricKey(key, hex_b)

	if err != nil {

		return nil, fmt.Errorf("failed to decrypt data: %s", err.Error())
	}

	return dec_b, nil

}
