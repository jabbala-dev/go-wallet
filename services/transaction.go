package services

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var ethClient *ethclient.Client

func init() {
	var err error
	ethClient, err = ethclient.Dial("https://mainnet.infura.io/v3/"+os.Getenv("INFURA_PROJECT_ID"))
	if err != nil {
		log.Fatal(err)
	}
}

func CreateAndSendTransaction(toAddress string, value int64)(string, error){

	privateKey, err := loadKey()
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	gasLimit := uint64(21000)
	gasprice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	to:= common.HexToAddress(toAddress)
	chainID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		return "", err
	}

	tx := types.NewTransaction(nonce, to, big.NewInt(value), gasLimit, gasprice, nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	err = ethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}
