package main

import (
	"BRLA_Golang_Code_Snippets/models"
	"BRLA_Golang_Code_Snippets/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Fill those
const webhookUrl = "https://myendpoint.com"
const privateKeyFileName = "key/keypair.pem"
const apiKey = "fc58f55a-7f0a-45f1-9616-acdaa9e77bb2"

// Sandbox endpoint
const apiEndpoint = "https://api.brla.digital:4567/v1/superuser/webhooks"

func main() {

	privateKey, err := utils.OpenPrivateKeyFile(privateKeyFileName)
	if err != nil {
		log.Fatal(err)
	}

	bodyInput, _ := json.Marshal(models.RegisterWebhookInput{
		Url: webhookUrl,
	})

	req, _ := http.NewRequest("POST", apiEndpoint, bytes.NewReader(bodyInput))
	req.Header.Add("accept", "application/json")

	err = utils.SignRequest(req, privateKey, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	bodyOutput, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusCreated {
		log.Fatal(fmt.Sprintf("Response returned status %d: %s", res.StatusCode, string(bodyOutput)))
	}

	var respData models.RegisterWebhookOutput
	json.Unmarshal(bodyOutput, &respData)

	log.Println("Webhook Id", respData.WebhookId)

}
