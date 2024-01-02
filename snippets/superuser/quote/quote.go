package main

import (
	"BRLA_Golang_Code_Snippets/models"
	"BRLA_Golang_Code_Snippets/utils"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"log"
	"net/http"
)

// Fill those
const operation = models.SWAP // models.QuoteOperationType
const amount = 10000          // In cents of BRL
const chain = "Polygon"
const fixOutput = true
const inputCoin = "USDC"
const outputCoin = "BRLA"

const privateKeyFileName = "key/keypair.pem"
const apiKey = "fc58f55a-7f0a-45f1-9616-acdaa9e77bb2"

// Sandbox endpoint
const apiEndpoint = "https://api.brla.digital:4567/v1/superuser/fast-quote"

func main() {

	privateKey, err := utils.OpenPrivateKeyFile(privateKeyFileName)
	if err != nil {
		log.Fatal(err)
	}

	input := models.QuoteInput{
		Operation:  string(operation),
		Amount:     amount,
		Chain:      chain,
		FixOutput:  fixOutput,
		InputCoin:  inputCoin,
		OutputCoin: outputCoin,
	}

	queryData, err := query.Values(input)
	if err != nil {
		log.Fatal(err)
	}
	queryString := queryData.Encode()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s?%s", apiEndpoint, queryString), nil)
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

}
