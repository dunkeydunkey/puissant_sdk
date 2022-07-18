package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"time"

	"github.com/du5/puissant_demo/bnb48.sdk"
	"github.com/du5/puissant_demo/demo"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func main() {

	conf := demo.GetConf("config.yaml")

	client, err := bnb48.Dial("https://testnet-fonce-bsc.bnb48.club", "https://testnet-fonce-bsc.bnb48.club")
	if err != nil {
		log.Panicln(err.Error())
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Panicln(err.Error())
	}
	log.Printf("chainID: %s", chainID.String())

	privateKey, err := crypto.HexToECDSA(conf.Wallet[0])

	if err != nil {
		log.Panicln(err.Error())
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Panicln(err.Error())
	}

	value := big.NewInt(2e18)
	gasLimit := uint64(21000)
	gasPrice, _ := client.SuggestGasPrice(context.Background())

	var txs []string
	for k := range make([]int, 10) {
		tx := types.NewTransaction(nonce+uint64(k), fromAddress, value, gasLimit, gasPrice, nil)
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			log.Panicln(err.Error())
		}
		rawTxBytes, _ := rlp.EncodeToBytes(types.Transactions{signedTx}[0])
		rawTxHex := hexutil.Encode(rawTxBytes)

		txs = append(txs, rawTxHex)
	}

	// send puissant tx
	res, err := client.SendPuissant(context.Background(), txs, time.Now().Unix()+60, nil)
	if err != nil {
		log.Panicln(err.Error())
	}

	log.Println(res)
}