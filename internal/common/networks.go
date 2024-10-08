package common

import (
	"fmt"
	"math/big"
)

const (
	NetworkHolesky = "holesky"
	NetworkMainnet = "mainnet"
)

func AssertChainID(network string, chainId *big.Int) error {
	switch network {
	case NetworkHolesky:
		if chainId.Cmp(big.NewInt(17000)) != 0 {
			return fmt.Errorf("invalid chain id for network: %s", network)
		}
	case NetworkMainnet:
		if chainId.Cmp(big.NewInt(1)) != 0 {
			return fmt.Errorf("invalid chain id for network: %s", network)
		}
	default:
		return fmt.Errorf("invalid network: %s", network)
	}
	return nil
}
