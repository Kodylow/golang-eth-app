package main

import (
	"context"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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
	// keystorePath := os.Getenv("KEYSTORE_PATH")
	// password := os.Getenv("PASSWORD")
	// infuraURL := os.Getenv("INFURA_URL")
	ganacheURL := "http://localhost:8545"

	// Connect to a node
	ethClient, err := ethclient.DialContext(context.Background(), ganacheURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer ethClient.Close()

	pk1Bytes := common.FromHex(os.Getenv("PRIVKEY1HEX"))
	pk2Bytes := common.FromHex(os.Getenv("PRIVKEY2HEX"))
	pk1, err := crypto.ToECDSA(pk1Bytes)
	if err != nil {
		log.Fatalf("Failed to convert private key 1: %v", err)
	}
	pk2, err := crypto.ToECDSA(pk2Bytes)
	if err != nil {
		log.Fatalf("Failed to convert private key 2: %v", err)
	}

	a1 := crypto.PubkeyToAddress(pk1.PublicKey)
	a2 := crypto.PubkeyToAddress(pk2.PublicKey)

	log.Printf("Address 1: %s", a1.Hex())
	log.Printf("Address 2: %s", a2.Hex())

	// Check Balances
	a1Balance, err := ethClient.BalanceAt(context.Background(), a1, nil)
	if err != nil {
		log.Fatalf("Failed to get balance for addr1: %v", err)
	}
	a2Balance, err := ethClient.BalanceAt(context.Background(), a2, nil)
	if err != nil {
		log.Fatalf("Failed to get balance for addr2: %v", err)
	}

	log.Printf("Balance 1: %d", a1Balance)
	log.Printf("Balance 2: %d", a2Balance)

	// Create nonce
	nonce, err := ethClient.PendingNonceAt(context.Background(), a1)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	// Create transaction spending to a2
	tx := types.NewTransaction(nonce, a2, big.NewInt(1000000000000000000), 21000, big.NewInt(1000000000), nil)

	// Sign transaction from a1
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, pk1)

	// Send transaction
	err = ethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	// Check Balances
	a1Balance, err = ethClient.BalanceAt(context.Background(), a1, nil)
	if err != nil {
		log.Fatalf("Failed to get balance for addr1: %v", err)
	}

	a2Balance, err = ethClient.BalanceAt(context.Background(), a2, nil)
	if err != nil {
		log.Fatalf("Failed to get balance for addr2: %v", err)
	}

	log.Printf("Balance 1 after tx: %d", a1Balance)
	log.Printf("Balance 2 after tx: %d", a2Balance)
}
