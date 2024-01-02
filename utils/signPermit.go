package utils

import (
	"BRLA_Golang_Code_Snippets/models"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"math/big"
	"strconv"
	"time"
)

// Info below is for sandbox. Check https://brla-superuser-api.readme.io/reference/addresses for production
const spenderContract = "0xc2B3031C0596ADf63E903E6b5dBeD074B98a79FB"

const brlaContractAddress = "0x658e5ea3c7690f0626aff87ced6fc30021a93657"
const brlaContractName = "BRLA Token"
const brlaContractVersion = "1"

const usdcContractAddress = "0x0FA8781a83E46826621b3BC094Ea2A0212e71B23"
const usdcContractName = "USD Coin (PoS)"
const usdcContractVersion = "1"

const brlaExtraDecimals = "0000000000000000"
const usdcExtraDecimals = "0000"

type PermitToken string

const (
	BRLA PermitToken = "BRLA"
	USDC PermitToken = "USDC"
)

func SignPermit(token PermitToken, walletPrivateKeyHex string, amountCents int64, nonce int64) (*models.Permit, error) {

	privateKey, err := crypto.HexToECDSA(walletPrivateKeyHex)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.Public()
	publicKeyOk, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	owner := crypto.PubkeyToAddress(*publicKeyOk)
	spender := common.HexToAddress(spenderContract)
	chainId := 80001 // 80001 for Mumbai, 137 for Polygon
	d := time.Now().Add(60 * 60 * time.Second).Unix()

	var verifyingContract common.Address
	var contractName string
	var contractVersion string
	var tokenExtraDecimals string
	var legacyPermit bool

	switch token {
	case BRLA:
		verifyingContract = common.HexToAddress(brlaContractAddress)
		contractName = brlaContractName
		contractVersion = brlaContractVersion
		tokenExtraDecimals = brlaExtraDecimals
		legacyPermit = false
	case USDC:
		verifyingContract = common.HexToAddress(usdcContractAddress)
		contractName = usdcContractName
		contractVersion = usdcContractVersion
		tokenExtraDecimals = usdcExtraDecimals
		legacyPermit = true
	}

	val := big.NewInt(0)
	val.SetString(fmt.Sprintf("%d%s", amountCents, tokenExtraDecimals), 10)

	r, s, v, deadline, _, err := GenerateSignedPermit(
		contractName,
		contractVersion,
		legacyPermit,
		owner,
		spender,
		verifyingContract,
		int64(chainId),
		val,
		nonce,
		d,
		privateKey,
	)
	if err != nil {
		return nil, err
	}

	return &models.Permit{
		Deadline: deadline,
		Nonce:    nonce,
		R:        r,
		S:        s,
		V:        v,
	}, nil

}

func GeneratePermitHash(
	contractName string,
	contractVersion string,
	legacyPermit bool,
	owner common.Address,
	spender common.Address,
	verifyingContract common.Address,
	chainId int64,
	value *big.Int,
	nonce int64,
	deadline int64,
) (common.Hash, error) {

	val := math.HexOrDecimal256(*value)

	var domain apitypes.TypedDataDomain
	var typesPermit apitypes.Types

	if !legacyPermit {
		domain = apitypes.TypedDataDomain{
			Name:              contractName,
			Version:           contractVersion,
			ChainId:           math.NewHexOrDecimal256(chainId),
			VerifyingContract: verifyingContract.String(),
		}
		typesPermit = apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"Permit": []apitypes.Type{
				{Name: "owner", Type: "address"},
				{Name: "spender", Type: "address"},
				{Name: "value", Type: "uint256"},
				{Name: "nonce", Type: "uint256"},
				{Name: "deadline", Type: "uint256"},
			},
		}
	} else {
		domain = apitypes.TypedDataDomain{
			Name:              contractName,
			Version:           contractVersion,
			Salt:              common.HexToHash(strconv.FormatInt(chainId, 16)).String(),
			VerifyingContract: verifyingContract.String(),
		}
		typesPermit = apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "verifyingContract", Type: "address"},
				{Name: "salt", Type: "bytes32"},
			},
			"Permit": []apitypes.Type{
				{Name: "owner", Type: "address"},
				{Name: "spender", Type: "address"},
				{Name: "value", Type: "uint256"},
				{Name: "nonce", Type: "uint256"},
				{Name: "deadline", Type: "uint256"},
			},
		}
	}

	signerData := apitypes.TypedData{
		Types:       typesPermit,
		PrimaryType: "Permit",
		Domain:      domain,
		Message: apitypes.TypedDataMessage{
			"owner":    owner.String(),
			"spender":  spender.String(),
			"value":    &val,
			"nonce":    math.NewHexOrDecimal256(nonce),
			"deadline": math.NewHexOrDecimal256(deadline),
		},
	}

	domainSeparator, err := signerData.HashStruct("EIP712Domain", signerData.Domain.Map())

	if err != nil {
		return common.Hash{}, err
	}

	typedDataHash, err := signerData.HashStruct(signerData.PrimaryType, signerData.Message)
	if err != nil {
		return common.Hash{}, err
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hash := common.BytesToHash(crypto.Keccak256(rawData))

	return hash, nil
}

func GenerateSignedPermit(
	contractName string,
	contractVersion string,
	legacyPermit bool,
	owner common.Address,
	spender common.Address,
	verifyingContract common.Address,
	chainId int64,
	value *big.Int,
	nonce int64,
	deadline int64,
	privateKey *ecdsa.PrivateKey,
) (r string, s string, v uint8, dl int64, hash common.Hash, err error) {

	hash, err = GeneratePermitHash(contractName, contractVersion, legacyPermit, owner, spender, verifyingContract, chainId, value, nonce, deadline)
	if err != nil {
		return
	}

	signatureBytes, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return
	}

	r = hexutil.Encode(signatureBytes[:32])
	s = hexutil.Encode(signatureBytes[32:64])
	v = uint8(int(signatureBytes[64])) + 27
	dl = deadline

	return

}
