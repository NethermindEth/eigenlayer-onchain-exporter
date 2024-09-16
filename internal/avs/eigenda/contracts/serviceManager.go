package contracts

import (
	"bytes"
	_ "embed"
	"fmt"

	"github.com/NethermindEth/eigenda-blob-scrapper/internal/config"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// Holesky variables
var (
	//go:embed abi/holesky-service-manager.json
	holeskyServiceManagerABIBytes []byte
	// TODO: add a configuration option for the address
	holeskyServiceManagerAddress  = common.HexToAddress("0xD4A7E1Bd8015057293f0D0A557088c286942e84b")
	holeskyServiceManagerContract *ServiceManagerContract
)

var (
	//go:embed abi/mainnet-service-manager.json
	mainnetServiceManagerABIBytes []byte
	// TODO: add a configuration option for the address
	mainnetServiceManagerAddress  = common.HexToAddress("0x870679E138bCdf293b7Ff14dD44b70FC97e12fc0")
	mainnetServiceManagerContract *ServiceManagerContract
)

type ServiceManagerContract struct {
	Address common.Address
	Abi     abi.ABI
}

func GetServiceManagerContract(avsEnv string) (*ServiceManagerContract, error) {
	switch avsEnv {
	case config.AVSEnvEigenDAHolesky:
		if holeskyServiceManagerContract != nil {
			return holeskyServiceManagerContract, nil
		}
		return getHoleskyServiceManagerContract()
	case config.AVSEnvEigenDAMainnet:
		if mainnetServiceManagerContract != nil {
			return mainnetServiceManagerContract, nil
		}
		return getMainnetServiceManagerContract()
	default:
		return nil, fmt.Errorf("invalid avs environment: %s", avsEnv)
	}
}

func getHoleskyServiceManagerContract() (*ServiceManagerContract, error) {
	abi, err := abi.JSON(bytes.NewReader(holeskyServiceManagerABIBytes))
	if err != nil {
		return nil, err
	}
	holeskyServiceManagerContract = &ServiceManagerContract{
		Address: holeskyServiceManagerAddress,
		Abi:     abi,
	}
	return holeskyServiceManagerContract, nil
}

func getMainnetServiceManagerContract() (*ServiceManagerContract, error) {
	abi, err := abi.JSON(bytes.NewReader(mainnetServiceManagerABIBytes))
	if err != nil {
		return nil, err
	}
	mainnetServiceManagerContract = &ServiceManagerContract{
		Address: mainnetServiceManagerAddress,
		Abi:     abi,
	}
	return mainnetServiceManagerContract, nil
}
