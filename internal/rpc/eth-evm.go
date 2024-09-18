package rpc

import (
	"context"
	"log/slog"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ethEvmRpcs = make(map[string]EthEvmRpc)
)

// EthEvmRpc is the interface for the Ethereum RPC client.
type EthEvmRpc interface {
	BlockNumber(ctx context.Context) (uint64, error)
	FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
}

type ethEvmRpc struct {
	network        string
	client         *ethclient.Client
	maxElapsedTime time.Duration
}

func NewEthEvmRpc(network string, url string, maxElapsedTime time.Duration) (EthEvmRpc, error) {
	if _, ok := ethEvmRpcs[network]; ok {
		return ethEvmRpcs[network], nil
	}
	slog.Debug("initializing new eth-evm rpc |", "network", network)
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	ethEvmRpc := &ethEvmRpc{network: network, client: client, maxElapsedTime: maxElapsedTime}
	ethEvmRpcs[network] = ethEvmRpc

	return ethEvmRpc, nil
}

func (e *ethEvmRpc) BlockNumber(ctx context.Context) (uint64, error) {
	operation := func() (uint64, error) {
		slog.Debug("getting block number |", "rpcNetwork", e.network)
		return e.client.BlockNumber(ctx)
	}
	notify := func(err error, duration time.Duration) {
		slog.Error("failed to get block number, retrying... |", "rpc-network", e.network, "duration", duration, "error", err)
	}

	return backoff.RetryNotifyWithData(
		operation,
		backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(e.maxElapsedTime)),
		notify,
	)
}

func (e *ethEvmRpc) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	operation := func() ([]types.Log, error) {
		slog.Debug("filtering logs |", "rpc-network", e.network)
		return e.client.FilterLogs(ctx, query)
	}
	notify := func(err error, duration time.Duration) {
		slog.Error("failed to filter logs, retrying... |", "rpc-network", e.network, "duration", duration, "error", err)
	}

	return backoff.RetryNotifyWithData(
		operation,
		backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(e.maxElapsedTime)),
		notify,
	)
}

func (e *ethEvmRpc) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	type result struct {
		tx        *types.Transaction
		isPending bool
	}
	operation := func() (result, error) {
		slog.Debug("getting transaction by hash |", "rpc-network", e.network, "hash", hash)
		tx, isPending, err := e.client.TransactionByHash(ctx, hash)
		return result{tx: tx, isPending: isPending}, err
	}
	notify := func(err error, duration time.Duration) {
		slog.Error("failed to get transaction by hash, retrying... |", "rpc-network", e.network, "error", err)
	}

	out, err := backoff.RetryNotifyWithData(
		operation,
		backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(e.maxElapsedTime)),
		notify,
	)
	return out.tx, out.isPending, err
}
