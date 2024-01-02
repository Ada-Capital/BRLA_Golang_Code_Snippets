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
const amount = 10000 // In cents of BRL
const chain = "Polygon"
const fixOutput = false
const outputCoin = "USDC"
const taxId = "55444879603" // Must've passed KYC
const receiverWalletAddress = "0x140f10EecF8C2090450F4f9c68Abe0bE78F5A84C"

const privateKeyFileName = "key/keypair.pem"
const apiKey = "fc58f55a-7f0a-45f1-9616-acdaa9e77bb2"

// Sandbox endpoint
const apiQuoteEndpoint = "https://api.brla.digital:4567/v1/superuser/fast-quote"
const apiPixToUsdEndpoint = "https://api.brla.digital:4567/v1/superuser/buy/pix-to-usd"

func main() {

	privateKey, err := utils.OpenPrivateKeyFile(privateKeyFileName)
	if err != nil {
		log.Fatal(err)
	}

	quoteInput := models.QuoteInput{
		Operation:  string(models.PIXTOUSD),
		Amount:     amount,
		Chain:      chain,
		FixOutput:  fixOutput,
		InputCoin:  "BRLA",
		OutputCoin: outputCoin,
	}

	queryData, err := query.Values(quoteInput)
	if err != nil {
		log.Fatal(err)
	}
	queryString := queryData.Encode()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s?%s", apiQuoteEndpoint, queryString), nil)
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

	var respData models.QuoteOutput
	json.Unmarshal(bodyOutput, &respData)

	log.Println("Quote data\n", respData)

	pixToUsdInput := models.PixToUsdInput{
		TaxId:           taxId,
		Token:           respData.Token,
		ReceiverAddress: receiverWalletAddress,
		MarkupAddress:   "",
	}

	bodyInput, _ := json.Marshal(pixToUsdInput)

	queryData2, err := query.Values(pixToUsdInput)
	if err != nil {
		log.Fatal(err)
	}
	queryString2 := queryData2.Encode()

	req2, _ := http.NewRequest("POST", fmt.Sprintf("%s?%s", apiPixToUsdEndpoint, queryString2), bytes.NewReader(bodyInput))
	req2.Header.Add("accept", "application/json")

	err = utils.SignRequest(req2, privateKey, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	res2, _ := http.DefaultClient.Do(req2)

	defer res2.Body.Close()
	bodyOutput2, _ := io.ReadAll(res2.Body)

	if res2.StatusCode != http.StatusOK {
		log.Fatal(fmt.Sprintf("Response returned status %d: %s", res2.StatusCode, string(bodyOutput2)))
	}

	var respData2 models.PixToUsdOutput
	json.Unmarshal(bodyOutput2, &respData2)

	log.Println("Pix to usd data\n", respData2)

}
