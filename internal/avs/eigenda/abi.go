package eigenda

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/avs/eigenda/contracts"
)

type confirmBatchInput struct {
	/* The confirmBatch input also has a 0th input which is not used. A new
	type could be created and added here if it becomes necessary.*/

	NonSignerStakesAndSignature nonSignerStakesAndSignature // input 1
}

type nonSignerStakesAndSignature struct {
	NonSignerQuorumBitmapIndices []uint32   `json:"nonSignerQuorumBitmapIndices"`
	NonSignerPubkeys             []g1Point  `json:"nonSignerPubkeys"`
	QuorumApks                   []g1Point  `json:"quorumApks"`
	ApkG2                        g2Point    `json:"apkG2"`
	Sigma                        g1Point    `json:"sigma"`
	QuorumApkIndices             []uint32   `json:"quorumApkIndices"`
	TotalStakeIndices            []uint32   `json:"totalStakeIndices"`
	NonSignerStakeIndices        [][]uint32 `json:"nonSignerStakeIndices"`
}

type g1Point struct {
	X *big.Int `json:"X"`
	Y *big.Int `json:"Y"`
}

type g2Point struct {
	X [2]*big.Int `json:"x"`
	Y [2]*big.Int `json:"y"`
}

func unpackConfirmBatchInput(avsEnv string, data []byte) (*confirmBatchInput, error) {
	contract, err := contracts.GetServiceManagerContract(avsEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to get service manager contract: %v", err)
	}

	method, exists := contract.Abi.Methods["confirmBatch"]
	if !exists {
		return nil, fmt.Errorf("confirmBatch method not found in ABI")
	}

	inputs, err := method.Inputs.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack confirmBatch input: %v", err)
	}

	// Unpack nonSignerStakesAndSignature
	var nonSignerStakesAndSignature nonSignerStakesAndSignature
	jsonRaw, err := json.Marshal(inputs[1])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal nonSignerStakesAndSignature: %v", err)
	}
	err = json.Unmarshal(jsonRaw, &nonSignerStakesAndSignature)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal nonSignerStakesAndSignature: %v", err)
	}

	return &confirmBatchInput{
		NonSignerStakesAndSignature: nonSignerStakesAndSignature,
	}, nil
}
