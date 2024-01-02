package main

import (
	"BRLA_Golang_Code_Snippets/models"
	"BRLA_Golang_Code_Snippets/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"log"
	"net/http"
)

// Fill those
const taxId = "55444879603"  // Must've passed KYC
const pixKey = "55444879603" // In sandbox, must be the same value of taxId
const walletAddress = "0x140f10EecF8C2090450F4f9c68Abe0bE78F5A84C"
const chain = "Polygon"
const amount = 8000 // In cents of BRL

const nonce = 0 // Can be calculated using an RPC client, but hardcoded for simplification

const privateKeyFileName = "key/keypair.pem"
const apiKey = "fc58f55a-7f0a-45f1-9616-acdaa9e77bb2"
const walletAddressPrivateKey = "240cde76f822b5af5ff1e1dc54af2a0d076bf9fa653e140ac3122f7adfa0e68c"

// Sandbox endpoint
const apiEndpoint = "https://api.brla.digital:4567/v1/superuser/sell"

func main() {

	privateKey, err := utils.OpenPrivateKeyFile(privateKeyFileName)
	if err != nil {
		log.Fatal(err)
	}

	permit, err := utils.SignPermit(utils.BRLA, walletAddressPrivateKey, amount, nonce)
	if err != nil {
		log.Fatal(err)
	}

	input := models.SellBRLAInput{
		TaxId:         taxId,
		PixKey:        pixKey,
		WalletAddress: walletAddress,
		Chain:         chain,
		Amount:        amount,
		Permit:        *permit,
	}

	bodyInput, _ := json.Marshal(input)

	queryData, err := query.Values(input)
	if err != nil {
		log.Fatal(err)
	}
	queryString := queryData.Encode()

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s?%s", apiEndpoint, queryString), bytes.NewReader(bodyInput))
	req.Header.Add("accept", "application/json")

	err = utils.SignRequest(req, privateKey, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	bodyOutput, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		log.Fatal(fmt.Sprintf("Response returned status %d: %s", res.StatusCode, string(bodyOutput)))
	}

	var respData models.SellBRLAOutput
	json.Unmarshal(bodyOutput, &respData)

	log.Println("Order Id", respData.Id)

}
