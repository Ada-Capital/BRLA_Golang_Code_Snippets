package models

import (
	"time"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	AccessToken string `json:"accessToken"`
}

type RegisterApiKeyInput struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
}

type RegisterApiKeyOutput struct {
	ApiKey string `json:"apiKey"`
}

type RegisterWebhookInput struct {
	Url string `json:"url"`
}

type RegisterWebhookOutput struct {
	WebhookId string `json:"webhookId"`
}

type QuoteOperationType string

const (
	SWAP     QuoteOperationType = "swap"
	PIXTOUSD QuoteOperationType = "pix-to-usd"
)

type QuoteInput struct {
	Operation  string `url:"operation"`
	Amount     int64  `url:"amount"`
	Chain      string `url:"chain"`
	FixOutput  bool   `url:"fixOutput,omitempty"`
	InputCoin  string `url:"inputCoin"`
	OutputCoin string `url:"outputCoin"`
	Markup     string `url:"markup,omitempty"`
}

type QuoteOutput struct {
	Token      string `json:"token"`
	BasePrice  string `json:"basePrice"`
	Error      string `json:"error"`
	Sub        string `json:"sub"`
	Operation  string `json:"operation"`
	AmountBrl  string `json:"amountBrl"`
	AmountUsd  string `json:"amountUsd"`
	BaseFee    string `json:"baseFee"`
	GasFee     string `json:"gasFee"`
	MarkupFee  string `json:"markupFee"`
	InputCoin  string `json:"inputCoin"`
	OutputCoin string `json:"outputCoin"`
	Chain      string `json:"chain"`
}

type PixToUsdInput struct {
	TaxId           string `url:"taxId" json:"-"`
	Token           string `json:"token"`
	ReceiverAddress string `json:"receiverAddress"`
	MarkupAddress   string `json:"markupAddress,omitempty"`
}

type PixToUsdOutput struct {
	Id     string    `json:"id"`
	Due    time.Time `json:"due"`
	BrCode string    `json:"brCode"`
}

type SandboxApproveKycLevel1Input struct {
	Cpf       string `json:"cpf"`
	BirthDate string `json:"birthDate"`
	FullName  string `json:"fullName"`
}

type SandboxApproveKycLevel1Output struct {
	Id string `json:"id"`
}

type Permit struct {
	Deadline int64  `json:"deadline"`
	Nonce    int64  `json:"nonce"`
	R        string `json:"r"` //[32]byte
	S        string `json:"s"` //[32]byte
	V        uint8  `json:"v"`
}

type SellBRLAInput struct {
	TaxId             string `json:"-" url:"taxId"`
	PixKey            string `json:"pixKey"`
	WalletAddress     string `json:"walletAddress"`
	Chain             string `json:"chain"`
	Amount            int64  `json:"amount"`
	Signature         string `json:"signature,omitempty"`
	SignatureDeadline int64  `json:"signatureDeadline,omitempty"`
	Permit            Permit `json:"permit,omitempty"`
}

type SellBRLAOutput struct {
	Id string `json:"id"`
}

type SwapInput struct {
	TaxId               string `json:"-" url:"taxId"`
	Token               string `json:"token"`
	WalletAddress       string `json:"walletAddress"`
	ReceiverAddress     string `json:"receiverAddress,omitempty"`
	EnforceAtomicSwap   bool   `json:"enforceAtomicSwap"`
	Permit              Permit `json:"permit,omitempty"`
	ReturnAuthorization bool   `json:"returnAuthorization,omitempty"`
	OperatorWallet      string `json:"operatorWallet,omitempty"`
	MarkupAddress       string `json:"markupAddress,omitempty"`
	FunctionSignature   string `json:"functionSignature,omitempty"`
	Signature           string `json:"signature,omitempty"`
	SignatureDeadline   int64  `json:"signatureDeadline,omitempty"`
}

type SwapOutput struct {
	Id            string `json:"id"`
	Authorization string `json:"authorization,omitempty"`
}
