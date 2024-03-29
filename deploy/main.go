package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	todo "github.com/habfem/learn/11-solidity/gen"
)

func main() {
	b, err := os.ReadFile("wallet/UTC--2023-11-07T08-26-52.205876600Z--33cf0dc27a3334ac2d9c579b7e2ee56cceaa4b62")
	if err != nil {
		log.Fatal(err)
	}

	key, err := keystore.DecryptKey(b, "password")
	if err != nil {
		log.Fatal(err)
	}

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/c76438a754b54ef2bea366c87b7500fb")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	add := crypto.PubkeyToAddress(key.PrivateKey.PublicKey)

	nonce, err := client.PendingNonceAt(context.Background(), add)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(3000000)
	auth.Nonce = big.NewInt(int64(nonce))

	a, tx, _, err := todo.DeployTodo(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("-----------------------------")
	fmt.Println(a.Hex())
	fmt.Println(tx.Hash().Hex())
	fmt.Println("-----------------------------")
}
