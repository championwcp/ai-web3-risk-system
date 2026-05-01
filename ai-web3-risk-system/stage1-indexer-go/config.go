package main

import (
	"fmt"
	"math/big"
	"os"
)

// LoadBlockRangeFromEnv 从环境变量读取本次要扫描的区块范围。
// 这样不用改代码就能调整查询范围，也能避免 indexer 永远只查固定区块。
func LoadBlockRangeFromEnv() (*big.Int, *big.Int, error) {
	fromBlockStr := os.Getenv("FROM_BLOCK")
	if fromBlockStr == "" {
		return nil, nil, fmt.Errorf("missing FROM_BLOCK")
	}

	toBlockStr := os.Getenv("TO_BLOCK")
	if toBlockStr == "" {
		return nil, nil, fmt.Errorf("missing TO_BLOCK")
	}

	fromBlock := new(big.Int)
	if _, ok := fromBlock.SetString(fromBlockStr, 10); !ok {
		return nil, nil, fmt.Errorf("invalid FROM_BLOCK: %s", fromBlockStr)
	}

	toBlock := new(big.Int)
	if _, ok := toBlock.SetString(toBlockStr, 10); !ok {
		return nil, nil, fmt.Errorf("invalid TO_BLOCK: %s", toBlockStr)
	}

	if fromBlock.Cmp(toBlock) > 0 {
		return nil, nil, fmt.Errorf("FROM_BLOCK must be <= TO_BLOCK")
	}

	return fromBlock, toBlock, nil
}
