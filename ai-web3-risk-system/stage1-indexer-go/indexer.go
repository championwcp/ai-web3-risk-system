package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ProcessTransferLogsResult struct {
	Total     int
	Succeeded int
	Failed    int
	Skipped   int
}

// ProcessTransferLog 负责处理单条 raw log。
// 主流程是：先判断是不是 Transfer，再解析成 TransferEvent，最后写入数据库。
func ProcessTransferLog(
	ctx context.Context,
	db *sql.DB,
	rawLog types.Log,
) error {
	if len(rawLog.Topics) == 0 || rawLog.Topics[0] != transferEventSigHash {
		return nil
	}

	event, err := ParseTransferEvent(rawLog)
	if err != nil {
		return err
	}

	return InsertTransferEvent(ctx, db, event)
}

// ProcessTransferLogs 负责批量处理日志。
// 它复用 ProcessTransferLog，并统计成功、跳过、失败，方便观察真实运行结果。
func ProcessTransferLogs(
	ctx context.Context,
	db *sql.DB,
	logs []types.Log,
) ProcessTransferLogsResult {
	result := ProcessTransferLogsResult{
		Total: len(logs),
	}

	for _, rawLog := range logs {

		err := ProcessTransferLog(ctx, db, rawLog)
		if err != nil {
			// skipped 表示“当前阶段不支持这种日志布局”，不是程序真正失败。
			if errors.Is(err, ErrUnsupportedTransferLog) {
				result.Skipped++
				continue
			}

			result.Failed++
			log.Printf(
				"failed to process log tx=%s index=%d: %v",
				rawLog.TxHash.Hex(),
				rawLog.Index,
				err,
			)
			continue
		}

		result.Succeeded++
	}

	return result
}

// FetchTransferLogs 从 Sepolia 的一小段历史区块里查询 Transfer 签名日志。
// 现在先固定小范围，原因是学习阶段要先跑通闭环，而且免费 RPC 通常限制日志查询范围。
func FetchTransferLogs(ctx context.Context, client *ethclient.Client) ([]types.Log, error) {
	fromBlock, toBlock, err := LoadBlockRangeFromEnv()
	if err != nil {
		return nil, err
	}

	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Topics: [][]common.Hash{
			{transferEventSigHash},
		},
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("filter logs failed: %w", err)
	}

	fmt.Printf("queried logs count: %d\n", len(logs))
	if len(logs) == 0 {
		return nil, fmt.Errorf(
			"no transfer logs found in block range [%d, %d]",
			fromBlock.Uint64(),
			toBlock.Uint64(),
		)
	}

	first := logs[0]
	fmt.Println("first transfer log in fetched batch:")
	fmt.Println("Address:", first.Address.Hex())
	fmt.Println("BlockNumber:", first.BlockNumber)
	fmt.Println("TxHash:", first.TxHash.Hex())
	fmt.Println("Topics count:", len(first.Topics))
	fmt.Println("Data length:", len(first.Data))

	return logs, nil
}
