package main

import (
	"BRLA_Golang_Code_Snippets/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Fill those
const email = ""
const password = ""

// Sandbox endpoint
const apiEndpoint = "https://api.brla.digital:4567/v1/superuser/login"

func main() {

	bodyInput, _ := json.Marshal(models.LoginInput{
		Email:    email,
		Password: password,
	})

	req, _ := http.NewRequest("POST", apiEndpoint, bytes.NewReader(bodyInput))

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	bodyOutput, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		log.Fatal(fmt.Sprintf("Response returned status %d: %s", res.StatusCode, string(bodyOutput)))
	}

	var respData models.LoginOutput
	json.Unmarshal(bodyOutput, &respData)

	log.Println("Jwt token:", respData.AccessToken)

}
