package main

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var infuraURL = "https://mainnet.infura.io/v3/360031b1b30f4b8b92b6a27850e11b8d"
var ganacheURL = "http://localhost:8545"

func main() {

	ethClient, err := ethclient.DialContext(context.Background(), ganacheURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer ethClient.Close()

	block, err := ethClient.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to get the latest block: %v", err)
	}
	log.Printf("Latest block: %d", block.NumberU64())

	addr := "0x0cd6f40fBceb4947749603cC069ed16D07FC548b"
	address := common.HexToAddress(addr)

	gweiBalance, err := ethClient.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}
	log.Printf("Balance Gwei: %d", gweiBalance)

	ethBalance := gweiToEth(gweiBalance)
	log.Printf("Balance ETH: %f", ethBalance)

}

func gweiToEth(balance *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1000000000000000000))
}
