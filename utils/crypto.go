package utils

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

func PublicKeyToAddress(publicKey *ecdsa.PublicKey) string {
	return crypto.PubkeyToAddress(*publicKey).Hex()
}

func PrivateKeyToHex(privateKey *ecdsa.PrivateKey) string {
	return hex.EncodeToString(crypto.FromECDSA(privateKey))
}

func HexToPrivateKey(hexKey string) (*ecdsa.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}
	return crypto.ToECDSA(privateKeyBytes)
}

func SignMessage(privateKey *ecdsa.PrivateKey, message string) (string, error) {
	hash := sha256.Sum256([]byte(message))
	signature, err := crypto.Sign(hash[:], privateKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signature), nil
}

func VerifyMessage(publicKey *ecdsa.PublicKey, message, signatureHex string) (bool, error) {
	hash := sha256.Sum256([]byte(message))
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, err
	}
	return crypto.VerifySignature(crypto.FromECDSAPub(publicKey), hash[:], signature[:len(signature)-1]), nil
}
