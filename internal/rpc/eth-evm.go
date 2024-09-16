package rpc

import (
	"context"
	"log/slog"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ethEvmRpcs = make(map[string]EthEvmRpc)
)

const retriesKey = contextKey("retries")

type contextKey string

// EthEvmRpc is the interface for the Ethereum RPC client.
type EthEvmRpc interface {
	BlockNumber(ctx context.Context) (uint64, error)
	FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
}

type ethEvmRpc struct {
	network string
	mutex   sync.Mutex
	client  *ethclient.Client
	retries int
}

func NewEthEvmRpc(network string, url string, retries int) (EthEvmRpc, error) {
	if _, ok := ethEvmRpcs[network]; ok {
		return ethEvmRpcs[network], nil
	}
	slog.Debug("initializing new eth-evm rpc |", "network", network)
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	ethEvmRpc := &ethEvmRpc{network: network, client: client, retries: retries}
	ethEvmRpcs[network] = ethEvmRpc

	return ethEvmRpc, nil
}

func (e *ethEvmRpc) BlockNumber(ctx context.Context) (uint64, error) {
	// Manage lock
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Init retries counter
	retries := ctx.Value(retriesKey)
	if retries == nil {
		ctx = context.WithValue(ctx, retriesKey, e.retries)
		retries = e.retries
	}

	slog.Debug("getting block number |", "rpcNetwork", e.network)
	blockNumber, err := e.client.BlockNumber(ctx)
	if err != nil {
		if retries.(int) > 0 {
			slog.Error("failed to get block number, retrying...", "rpc-network", e.network, "error", err)
			ctx = context.WithValue(ctx, retriesKey, retries.(int)-1)
			return e.BlockNumber(ctx)
		}
		slog.Error("failed to get block number", "rpc-network", e.network, "error", err)
		return 0, err
	}

	return blockNumber, nil
}

func (e *ethEvmRpc) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	// Manage lock
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Init retries counter
	retries := ctx.Value(retriesKey)
	if retries == nil {
		ctx = context.WithValue(ctx, retriesKey, e.retries)
		retries = e.retries
	}

	logs, err := e.client.FilterLogs(ctx, query)
	if err != nil {
		if retries.(int) > 0 {
			slog.Error("failed to filter logs, retrying... |", "rpc-network", e.network, "error", err)
			ctx = context.WithValue(ctx, retriesKey, retries.(int)-1)
			return e.FilterLogs(ctx, query)
		}
		slog.Error("failed to filter logs |", "rpc-network", e.network, "error", err)
		return nil, err
	}

	return logs, nil
}

func (e *ethEvmRpc) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	// Manage lock
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Init retries counter
	retries := ctx.Value(retriesKey)
	if retries == nil {
		ctx = context.WithValue(ctx, retriesKey, e.retries)
		retries = e.retries
	}

	slog.Debug("getting transaction by hash |", "rpc-network", e.network, "hash", hash)
	tx, isPending, err := e.client.TransactionByHash(ctx, hash)
	if err != nil {
		if retries.(int) > 0 {
			slog.Error("failed to get transaction by hash, retrying... |", "rpc-network", e.network, "error", err)
			ctx = context.WithValue(ctx, retriesKey, retries.(int)-1)
			return e.TransactionByHash(ctx, hash)
		}
		slog.Error("failed to get transaction by hash |", "rpc-network", e.network, "error", err)
		return nil, false, err
	}

	return tx, isPending, nil
}
