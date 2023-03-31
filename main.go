package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

type EthClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
}

func gweiToEth(gwei *big.Int) *big.Float {
	eth := new(big.Float).SetInt(gwei)
	eth = eth.Quo(eth, big.NewFloat(1e9))
	return eth
}

func generateAccount(password string) (accounts.Account, error) {
	key := keystore.NewKeyStore("./keystore", keystore.StandardScryptN, keystore.StandardScryptP)

	account, err := key.NewAccount(password)
	if err != nil {
		log.Fatalf("Failed to create account: %v", err)
		return accounts.Account{}, err
	}
	return account, nil

}

func check(ethClient *ethclient.Client) {
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

func main() {

	// Read in .env and set environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	keystorePath := os.Getenv("KEYSTORE_PATH")
	password := os.Getenv("PASSWORD")
	infuraURL := os.Getenv("INFURA_URL")
	// ganacheURL := "http://localhost:8545"

	// Connect to a node
	ethClient, err := ethclient.DialContext(context.Background(), infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer ethClient.Close()

	// Load keystore
	b, err := ioutil.ReadFile(keystorePath)
	if err != nil {
		log.Fatalf("Failed to read keystore: %v", err)
	}

	// Decrypt keystore
	key, err := keystore.DecryptKey(b, password)
	if err != nil {
		log.Fatalf("Failed to decrypt keystore: %v", err)
	}

	// Get account
	account := accounts.Account{
		Address: key.Address,
	}

	// Get balance
	balance, err := ethClient.BalanceAt(context.Background(), account.Address, nil)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}

	log.Printf("Balance Gwei: %d", balance)
}
