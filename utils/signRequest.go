package utils

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func SignRequest(r *http.Request, privateKey *rsa.PrivateKey, apiKey string) error {

	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
		r.Body.Close() //  must close
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	timestamp := time.Now().UnixMilli()
	requestMethod := r.Method
	endpointPath := r.URL.RequestURI()

	content := fmt.Sprintf("%s%s%s%s", strconv.FormatInt(timestamp, 10), requestMethod, endpointPath, string(bodyBytes))

	hasher := sha256.New()
	hasher.Write([]byte(content))
	hash := hasher.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash)
	if err != nil {
		return err
	}

	base64Signature := base64.StdEncoding.EncodeToString(signature)

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-API-Key", apiKey)
	r.Header.Set("X-API-Timestamp", strconv.FormatInt(timestamp, 10))
	r.Header.Set("X-API-Signature", base64Signature)

	return nil

}
