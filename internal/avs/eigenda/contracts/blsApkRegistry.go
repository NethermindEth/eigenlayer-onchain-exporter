package contracts

import (
	"bytes"
	_ "embed"
	"fmt"

	"github.com/NethermindEth/eigenlayer-onchain-exporter/internal/config"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// Holesky variables
var (
	//go:embed abi/holesky-bls-apk-registry.json
	holeskyBlsApkRegistryABIBytes []byte
	// TODO: add a configuration option for the address
	holeskyBlsApkRegistryAddress  = common.HexToAddress("0x066cF95c1bf0927124DFB8B02B401bc23A79730D")
	holeskyBlsApkRegistryContract *BlsApkRegistryContract
)

// Mainnet variables
var (
	//go:embed abi/mainnet-bls-apk-registry.json
	mainnetBlsApkRegistryABIBytes []byte
	// TODO: add a configuration option for the address
	mainnetBlsApkRegistryAddress  = common.HexToAddress("0x00A5Fd09F6CeE6AE9C8b0E5e33287F7c82880505")
	mainnetBlsApkRegistryContract *BlsApkRegistryContract
)

type BlsApkRegistryContract struct {
	Address common.Address
	Abi     abi.ABI
}

func GetBlsApkRegistryContract(avsEnv string) (*BlsApkRegistryContract, error) {
	switch avsEnv {
	case config.AVSEnvEigenDAHolesky:
		if holeskyBlsApkRegistryContract != nil {
			return holeskyBlsApkRegistryContract, nil
		}
		return getHoleskyBlsApkRegistryContract()
	case config.AVSEnvEigenDAMainnet:
		if mainnetBlsApkRegistryContract != nil {
			return mainnetBlsApkRegistryContract, nil
		}
		return getMainnetBlsApkRegistryContract()
	default:
		return nil, fmt.Errorf("invalid avs environment: %s", avsEnv)
	}
}

func getHoleskyBlsApkRegistryContract() (*BlsApkRegistryContract, error) {
	abi, err := abi.JSON(bytes.NewReader(holeskyBlsApkRegistryABIBytes))
	if err != nil {
		return nil, err
	}
	holeskyBlsApkRegistryContract = &BlsApkRegistryContract{
		Address: holeskyBlsApkRegistryAddress,
		Abi:     abi,
	}
	return holeskyBlsApkRegistryContract, nil
}

func getMainnetBlsApkRegistryContract() (*BlsApkRegistryContract, error) {
	abi, err := abi.JSON(bytes.NewReader(mainnetBlsApkRegistryABIBytes))
	if err != nil {
		return nil, err
	}
	mainnetBlsApkRegistryContract = &BlsApkRegistryContract{
		Address: mainnetBlsApkRegistryAddress,
		Abi:     abi,
	}
	return mainnetBlsApkRegistryContract, nil
}
