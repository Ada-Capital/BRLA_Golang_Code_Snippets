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
const amount = 2000 // In cents
const chain = "Polygon"
const fixOutput = false // If true, amount specifies output amount. Specifies input amount otherwise
const inputCoin = "USDC"
const outputCoin = "BRLA"
const taxId = "55444879603"                                          // Must've passed KYC
const walletAddress = "0x140f10EecF8C2090450F4f9c68Abe0bE78F5A84C"   // Wallet holding the tokens to swap
const receiverAddress = "0x140f10EecF8C2090450F4f9c68Abe0bE78F5A84C" // Wallet that will receive tokens after swapping
const enforceAtomicSwap = true                                       // If true, request will fail in case of not enough liquidity to atomically execute the swap

const nonce = 0 // Can be calculated using an RPC client, but hardcoded for simplification

const privateKeyFileName = "key/keypair.pem"
const apiKey = "fc58f55a-7f0a-45f1-9616-acdaa9e77bb2"
const walletAddressPrivateKey = "240cde76f822b5af5ff1e1dc54af2a0d076bf9fa653e140ac3122f7adfa0e68c"

// Sandbox endpoint
const apiQuoteEndpoint = "https://api.brla.digital:4567/v1/superuser/fast-quote"
const apiSwapEndpoint = "https://api.brla.digital:4567/v1/superuser/swap/v2/place-order"

func main() {

	if inputCoin == "USDT" {
		log.Fatal("USDT does not support permit")
	}

	privateKey, err := utils.OpenPrivateKeyFile(privateKeyFileName)
	if err != nil {
		log.Fatal(err)
	}

	quoteInput := models.QuoteInput{
		Operation:  string(models.SWAP),
		Amount:     amount,
		Chain:      chain,
		FixOutput:  fixOutput,
		InputCoin:  inputCoin,
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

	permit, err := utils.SignPermit(inputCoin, walletAddressPrivateKey, amount, nonce)
	if err != nil {
		log.Fatal(err)
	}

	swapInput := models.SwapInput{
		TaxId:             taxId,
		Token:             respData.Token,
		WalletAddress:     walletAddress,
		ReceiverAddress:   receiverAddress,
		EnforceAtomicSwap: enforceAtomicSwap,
		Permit:            *permit,
	}

	bodyInput, _ := json.Marshal(swapInput)

	queryData2, err := query.Values(swapInput)
	if err != nil {
		log.Fatal(err)
	}
	queryString2 := queryData2.Encode()

	req2, _ := http.NewRequest("POST", fmt.Sprintf("%s?%s", apiSwapEndpoint, queryString2), bytes.NewReader(bodyInput))
	req2.Header.Add("accept", "application/json")

	err = utils.SignRequest(req2, privateKey, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	res2, _ := http.DefaultClient.Do(req2)

	defer res.Body.Close()
	bodyOutput2, _ := io.ReadAll(res2.Body)

	if res2.StatusCode != http.StatusOK {
		log.Fatal(fmt.Sprintf("Response returned status %d: %s", res2.StatusCode, string(bodyOutput)))
	}

	var respData2 models.SwapOutput
	json.Unmarshal(bodyOutput2, &respData2)

	log.Println("Swap data\n", respData2)

}
