package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func OpenPrivateKeyFile(privateKeyFileName string) (*rsa.PrivateKey, error) {

	privateKeyFile, err := os.ReadFile(privateKeyFileName)
	if err != nil {
		return nil, err
	}

	privateKeyBlock, _ := pem.Decode(privateKeyFile)
	if privateKeyBlock == nil {
		return nil, errors.New("error decoding private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil

}