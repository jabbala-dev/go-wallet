package services

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jabbala/go-wallet/utils"
)

var privateKeyFile = "private_key.txt"

func GenerateKeyPair() (string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey).Hex()

	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))

	err = ioutil.WriteFile(privateKeyFile, []byte(privateKeyHex), 0600)
	if err != nil {
		return "", "", err
	}

	return privateKeyHex, address, nil
}

func GetAddress() (string, error) {
	privateKey, err := loadKey()
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey).Hex()

	return address, nil
}

func SignMessage(message string) (string, error) {
	privateKey, err := loadKey()
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(message))
	signature, err := crypto.Sign(hash[:], privateKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(signature), nil
}

func VerifyMessage(message, signatureHex string) (bool, error) {
	privateKey, err := loadKey()
	if err != nil {
		return false, err
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	hash := sha256.Sum256([]byte(message))
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, err
	}

	return crypto.VerifySignature(crypto.FromECDSAPub(publicKey), hash[:], signature[:len(signature)-1]), nil
}

func loadKey() (*ecdsa.PrivateKey, error) {
	if _, err := os.Stat(privateKeyFile); os.IsNotExist(err) {
		return nil, errors.New("private key file does not exist")
	}

	privateKeyHex, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, err
	}

	privateKeyBytes, err := hex.DecodeString(string(privateKeyHex))
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
