package main

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const epsilon = 1e-9

func floatEqual(a, b *big.Float) bool {
	diff := new(big.Float).Sub(a, b)
	return diff.Abs(diff).Cmp(big.NewFloat(epsilon)) < 0
}

func TestGweiToEth(t *testing.T) {
	testCases := []struct {
		name     string
		balance  *big.Int
		expected *big.Float
	}{
		{
			name:     "Converts 1 Gwei to 0.000000001 ETH",
			balance:  big.NewInt(1),
			expected: big.NewFloat(0.000000001),
		},
		{
			name:     "Converts 1000000000000000000 Gwei to 1000000000 ETH",
			balance:  big.NewInt(1000000000000000000),
			expected: big.NewFloat(1000000000),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := gweiToEth(tc.balance)
			if !floatEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestMain(t *testing.T) {
	// Skip the main function, since it requires an active Ethereum node.
	t.SkipNow()
}

// Mocking the Ethereum client
type mockEthClient struct{}

func (m *mockEthClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(42),
	}), nil
}

func (m *mockEthClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return big.NewInt(1000000000000000000), nil
}

func TestEthereumClient(t *testing.T) {
	ethClient := &mockEthClient{}

	block, err := ethClient.BlockByNumber(context.Background(), nil)
	if err != nil {
		t.Fatalf("Failed to get the latest block: %v", err)
	}
	if block.NumberU64() == 0 {
		t.Error("Block number should not be zero")
	}

	addr := "0x0cd6f40fBceb4947749603cC069ed16D07FC548b"
	address := common.HexToAddress(addr)
	gweiBalance, err := ethClient.BalanceAt(context.Background(), address, nil)
	if err != nil {
		t.Fatalf("Failed to get balance: %v", err)
	}
	if gweiBalance.Cmp(big.NewInt(0)) == 0 {
		t.Error("Balance should not be zero")
	}

	ethBalance := gweiToEth(gweiBalance)
	if ethBalance.Cmp(big.NewFloat(0)) == 0 {
		t.Error("ETH balance should not be zero")
	}
}
