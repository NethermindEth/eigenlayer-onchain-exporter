package eigenda

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"slices"
	"sort"
	"strconv"
	"time"

	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/avs/eigenda/contracts"
	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/avsexporter"
	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/config"
	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/rpc"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type eigenDAOnChainExporter struct {
	avsEnv    string
	network   string
	operators []config.OperatorConfig
	ethClient rpc.EthEvmRpc
}

func NewEigenDAOnChainExporter(avsEnv string, c *config.Config) (avsexporter.AVSExporter, error) {
	// Filter operators by AVS environment
	var operators []config.OperatorConfig
	for _, operator := range c.Operators {
		if slices.Contains(operator.AVSEnvs, avsEnv) {
			operators = append(operators, operator)
		}
	}
	// Get the network from the AVS environment
	var network string
	switch avsEnv {
	case config.AVSEnvEigenDAHolesky:
		network = "holesky"
	case config.AVSEnvEigenDAMainnet:
		network = "mainnet"
	default:
		return nil, fmt.Errorf("invalid AVS environment: %s", avsEnv)
	}
	e := &eigenDAOnChainExporter{
		avsEnv:    avsEnv,
		network:   network,
		operators: operators,
	}
	if err := e.init(c.RPCs); err != nil {
		return nil, fmt.Errorf("failed to initialize exporter: %v", err)
	}
	slog.Info("initialized exporter |", "avsEnv", e.avsEnv, "operators", len(e.operators))
	return e, nil
}

func (e *eigenDAOnChainExporter) Run(ctx context.Context, c *config.Config) error {
	// Set exporter status to DOWN by default
	metricExporterStatus.WithLabelValues(e.avsEnv).Set(0)

	// TODO: Add a configuration option for the ticker time
	tickerTime := time.Second * 30
	slog.Info("running exporter |", "avsEnv", e.avsEnv, "interval", tickerTime)

	// Load contracts
	serviceManagerContract, err := contracts.GetServiceManagerContract(e.avsEnv)
	if err != nil {
		return err
	}
	blsApkRegistryContract, err := contracts.GetBlsApkRegistryContract(e.avsEnv)
	if err != nil {
		return err
	}

	// Get current block to start from
	// TODO: Should we add a configuration option to start from a specific block?
	latestBlock, err := e.getLatestBlock()
	if err != nil {
		return err
	}

	// Set exporter status to UP
	metricExporterStatus.WithLabelValues(e.avsEnv).Set(1)

	ticker := time.Tick(tickerTime)
	for {
		select {
		case <-ctx.Done():
			slog.Info("exporter context done |", "avsEnv", e.avsEnv)
			return ctx.Err()
		case <-ticker:
			// Get the next block range
			fromBlock, toBlock, err := e.nextBlockRange(latestBlock, tickerTime)
			if err != nil {
				slog.Error("exporter error |", "avsEnv", e.avsEnv, "error", err)
				continue
			}
			if fromBlock == nil || toBlock == nil {
				continue
			}
			// Get logs from current block range
			logs, err := e.getLogs(fromBlock, toBlock)
			if err != nil {
				slog.Error("exporter error |", "avsEnv", e.avsEnv, "error", err)
				continue
			}

			for _, vLog := range logs {
				switch vLog.Topics[0].Hex() {
				case serviceManagerContract.Abi.Events["BatchConfirmed"].ID.Hex():
					if err := e.processBatchConfirmedLog(vLog); err != nil {
						slog.Error("exporter error |", "avsEnv", e.avsEnv, "error", err)
						continue
					}
				case blsApkRegistryContract.Abi.Events["OperatorRemovedFromQuorums"].ID.Hex():
					if err := e.processOperatorRemovedFromQuorumsLog(vLog); err != nil {
						slog.Error("exporter error |", "avsEnv", e.avsEnv, "error", err)
						continue
					}
				case blsApkRegistryContract.Abi.Events["OperatorAddedToQuorums"].ID.Hex():
					if err := e.processOperatorAddedToQuorumsLog(vLog); err != nil {
						slog.Error("exporter error |", "avsEnv", e.avsEnv, "error", err)
						continue
					}
				}
			}
			metricExporterLatestBlock.WithLabelValues(e.network).Set(float64(toBlock.Int64()))
			latestBlock = new(big.Int).Add(toBlock, big.NewInt(1))
		}
	}
}

func (e *eigenDAOnChainExporter) checkAVSEnv(avsEnv string) error {
	if avsEnv != config.AVSEnvEigenDAHolesky && avsEnv != config.AVSEnvEigenDAMainnet {
		return fmt.Errorf("invalid AVS environment: %s", avsEnv)
	}
	return nil
}

func (e *eigenDAOnChainExporter) init(rpcs map[string]string) error {
	if err := e.checkAVSEnv(e.avsEnv); err != nil {
		return fmt.Errorf("failed to check AVS environment: %v", err)
	}

	if err := e.initRPC(rpcs); err != nil {
		return fmt.Errorf("failed to initialize RPC: %v", err)
	}

	if err := e.initPrometheusMetrics(); err != nil {
		return fmt.Errorf("failed to initialize prometheus metrics: %v", err)
	}

	return nil
}

func (e *eigenDAOnChainExporter) initPrometheusMetrics() error {
	for _, operator := range e.operators {
		for quorum, in := range operator.EigenDAConfig.Quorums {
			if in {
				metricOnchainQuorumStatus.WithLabelValues(operator.Name, e.network, strconv.Itoa(quorum)).Set(1)
			} else {
				metricOnchainQuorumStatus.WithLabelValues(operator.Name, e.network, strconv.Itoa(quorum)).Set(0)
			}
		}
	}
	return nil
}

func (e *eigenDAOnChainExporter) initRPC(rpcs map[string]string) error {
	if rpcURL, ok := rpcs[e.network]; !ok {
		return fmt.Errorf("no RPC URL found for network: %s", e.network)
	} else {
		ethClient, err := rpc.NewEthEvmRpc(e.network, rpcURL, 3)
		if err != nil {
			return fmt.Errorf("failed to initialize RPC: %v", err)
		}
		e.ethClient = ethClient
	}
	return nil
}

func (e *eigenDAOnChainExporter) getLatestBlock() (*big.Int, error) {
	blockNumber, err := e.ethClient.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get the block number: %v", err)
	}
	return big.NewInt(int64(blockNumber)), nil
}

func (e *eigenDAOnChainExporter) nextBlockRange(latestBlock *big.Int, tickerTime time.Duration) (*big.Int, *big.Int, error) {
	toBlock, err := e.getLatestBlock()
	if err != nil {
		return nil, nil, err
	}
	if latestBlock.Cmp(toBlock) >= 0 {
		slog.Debug("latest RPC block is not greater than latest exporter block. Retrying after interval time |", "avsEnv", e.avsEnv, "rpcLatestBlock", toBlock, "exporterLatestBlock", latestBlock, "sleepTime", tickerTime)
		return nil, nil, nil
	}
	maxPaginationBlock := new(big.Int).Add(latestBlock, big.NewInt(1000))
	if toBlock.Cmp(maxPaginationBlock) > 0 {
		slog.Debug("latest block is greater than max pagination block. Using max pagination block instead |", "avsEnv", e.avsEnv, "latestBlock", toBlock, "maxPaginationBlock", maxPaginationBlock, "diff", new(big.Int).Sub(toBlock, maxPaginationBlock))
		toBlock = maxPaginationBlock
	}
	return latestBlock, toBlock, nil
}

func (e *eigenDAOnChainExporter) getLogs(fromBlock *big.Int, toBlock *big.Int) ([]types.Log, error) {
	// Load contracts
	serviceManagerContract, err := contracts.GetServiceManagerContract(e.avsEnv)
	if err != nil {
		return nil, err
	}
	blsApkRegistryContract, err := contracts.GetBlsApkRegistryContract(e.avsEnv)
	if err != nil {
		return nil, err
	}

	// Build the filter query
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []common.Address{
			serviceManagerContract.Address,
			blsApkRegistryContract.Address,
		},
		Topics: [][]common.Hash{
			{
				serviceManagerContract.Abi.Events["BatchConfirmed"].ID,
				blsApkRegistryContract.Abi.Events["OperatorAddedToQuorums"].ID,
				blsApkRegistryContract.Abi.Events["OperatorRemovedFromQuorums"].ID,
			},
		},
	}

	// Get the logs
	slog.Debug("filtering logs |", "avsEnv", e.avsEnv, "fromBlock", query.FromBlock, "toBlock", query.ToBlock)
	logs, err := e.ethClient.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, err
	}

	// Sort logs by block number and tx index
	sort.Slice(logs, func(i, j int) bool {
		if logs[i].BlockNumber == logs[j].BlockNumber {
			return logs[i].TxIndex < logs[j].TxIndex
		}
		return logs[i].BlockNumber < logs[j].BlockNumber
	})

	return logs, nil
}

func (e *eigenDAOnChainExporter) processBatchConfirmedLog(log types.Log) error {
	// Load contracts
	serviceManagerContract, err := contracts.GetServiceManagerContract(e.avsEnv)
	if err != nil {
		return err
	}

	// Increase the number of batches counter
	metricOnchainBatchesTotal.WithLabelValues(e.network).Inc()
	slog.Info("batch confirmed |", "avsEnv", e.avsEnv, "blockNumber", log.BlockNumber, "txHash", log.TxHash)

	// TODO: Ignoring the isPending output. Need to research more on this.
	tx, _, err := e.ethClient.TransactionByHash(context.Background(), log.TxHash)
	if err != nil {
		return fmt.Errorf("failed to get transaction by hash: %v", err)
	}

	// Get the function signature (first 4 bytes of the input data)
	funcSignature := tx.Data()[:4]
	if !bytes.Equal(funcSignature, serviceManagerContract.Abi.Methods["confirmBatch"].ID) {
		return nil
	}

	// Unpack the input data
	input, err := unpackConfirmBatchInput(e.avsEnv, tx.Data()[4:])
	if err != nil {
		return fmt.Errorf("failed to unpack confirmBatch input: %v", err)
	}

	// Iterate over the not signers and check if they are operators as not signers
	for _, operator := range e.operators {
		for _, pubkey := range input.NonSignerStakesAndSignature.NonSignerPubkeys {
			operatorBLSPubkeyX, operatorBLSPubkeyY, err := getOperatorBLSPubkey(operator)
			if err != nil {
				return fmt.Errorf("failed to get operator BLS public key: %v", err)
			}
			if operatorBLSPubkeyX.Cmp(pubkey.X) == 0 || operatorBLSPubkeyY.Cmp(pubkey.Y) == 0 {
				metricOnchainBatches.WithLabelValues(operator.Name, e.network, "missed").Inc()
				slog.Info("operator failed to sign batch |", "avsEnv", e.avsEnv, "blockNumber", log.BlockNumber, "txIndex", log.TxIndex, "operator", operator.Name)
			}
		}
	}

	return nil
}

func (e *eigenDAOnChainExporter) processOperatorRemovedFromQuorumsLog(log types.Log) error {
	// Load contracts
	blsApkRegistryContract, err := contracts.GetBlsApkRegistryContract(e.avsEnv)
	if err != nil {
		return err
	}

	logInputs, err := blsApkRegistryContract.Abi.Events["OperatorRemovedFromQuorums"].Inputs.Unpack(log.Data)
	if err != nil {
		return fmt.Errorf("failed to unpack operator removed from quorums log: %v", err)
	}
	operatorAddress := logInputs[0].(common.Address)
	quorumNumbers := logInputs[2].([]uint8)
	operatorIndex := slices.IndexFunc(e.operators, func(operator config.OperatorConfig) bool {
		return common.HexToAddress(operator.Address) == operatorAddress
	})
	if operatorIndex == -1 {
		return nil
	}
	for _, quorum := range quorumNumbers {
		slog.Info("operator removed from quorum |", "avsEnv", e.avsEnv, "blockNumber", log.BlockNumber, "txIndex", log.TxIndex, "operator", e.operators[operatorIndex].Name, "quorum", quorum)
		metricOnchainQuorumStatus.WithLabelValues(e.operators[operatorIndex].Name, e.network, strconv.Itoa(int(quorum))).Set(0)
	}

	return nil
}

func (e *eigenDAOnChainExporter) processOperatorAddedToQuorumsLog(log types.Log) error {
	// Load contracts
	blsApkRegistryContract, err := contracts.GetBlsApkRegistryContract(e.avsEnv)
	if err != nil {
		return err
	}

	logInputs, err := blsApkRegistryContract.Abi.Events["OperatorAddedToQuorums"].Inputs.Unpack(log.Data)
	if err != nil {
		return fmt.Errorf("failed to unpack operator added to quorums log: %v", err)
	}
	operatorAddress := logInputs[0].(common.Address)
	quorumNumbers := logInputs[2].([]uint8)
	operatorIndex := slices.IndexFunc(e.operators, func(operator config.OperatorConfig) bool {
		return common.HexToAddress(operator.Address) == operatorAddress
	})
	if operatorIndex == -1 {
		return nil
	}
	for _, quorum := range quorumNumbers {
		slog.Info("operator added to quorum |", "avsEnv", e.avsEnv, "blockNumber", log.BlockNumber, "txIndex", log.TxIndex, "operator", e.operators[operatorIndex].Name, "quorum", quorum)
		metricOnchainQuorumStatus.WithLabelValues(e.operators[operatorIndex].Name, e.network, strconv.Itoa(int(quorum))).Set(1)
	}
	return nil
}

func getOperatorBLSPubkey(operator config.OperatorConfig) (*big.Int, *big.Int, error) {
	x, ok := new(big.Int).SetString(operator.BLSPublicKey[0], 10)
	if !ok {
		return nil, nil, fmt.Errorf("failed to set string to big.Int: %s", operator.BLSPublicKey[0])
	}
	y, ok := new(big.Int).SetString(operator.BLSPublicKey[1], 10)
	if !ok {
		return nil, nil, fmt.Errorf("failed to set string to big.Int: %s", operator.BLSPublicKey[1])
	}
	return x, y, nil
}
