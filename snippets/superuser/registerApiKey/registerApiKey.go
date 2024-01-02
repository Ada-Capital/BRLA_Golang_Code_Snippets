package main

import (
	"BRLA_Golang_Code_Snippets/models"
	"BRLA_Golang_Code_Snippets/utils"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Fill those
const privateKeyFileName = "key/keypair.pem"
const jwtToken = ""
const keyName = "My API Key"

// Sandbox endpoint
const apiEndpoint = "https://api.brla.digital:4567/v1/superuser/api-keys"

func main() {

	privateKey, err := utils.OpenPrivateKeyFile(privateKeyFileName)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := &privateKey.PublicKey

	publicKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatal(err)
	}

	// Encode the public key to PEM format
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDER,
	})

	hasher := sha256.New()
	hasher.Write([]byte(keyName))
	hash := hasher.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash)
	if err != nil {
		log.Fatal(err)
	}

	base64Signature := base64.StdEncoding.EncodeToString(signature)

	log.Println("Signature:", base64Signature)

	bodyInput, _ := json.Marshal(models.RegisterApiKeyInput{
		Name:      keyName,
		PublicKey: string(publicKeyPEM),
		Signature: base64Signature,
	})

	req, _ := http.NewRequest("POST", apiEndpoint, bytes.NewReader(bodyInput))

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	bodyOutput, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusCreated {
		log.Fatal(fmt.Sprintf("Response returned status %d: %s", res.StatusCode, string(bodyOutput)))
	}

	var respData models.RegisterApiKeyOutput
	json.Unmarshal(bodyOutput, &respData)

	log.Println("Api key:", respData.ApiKey)

}
