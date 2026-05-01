package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
)

// InsertTransferEvent 把一条已经解析好的 TransferEvent 写入数据库。
// transaction_hash + log_index 可以唯一定位一条日志，也用于避免重复插入。
func InsertTransferEvent(
	ctx context.Context,
	db *sql.DB,
	event TransferEvent,
) error {
	if event.Value == nil {
		return errors.New("event value is nil")
	}

	query := `
		INSERT INTO transfer_events (
			transaction_hash,
			log_index,
			block_number,
			contract_address,
			from_address,
			to_address,
			value
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		ON CONFLICT (transaction_hash, log_index) DO NOTHING
	`

	// PostgreSQL numeric 可以安全接收大整数的十进制字符串形式，不会丢精度。
	valueStr := event.Value.String()

	_, err := db.ExecContext(
		ctx,
		query,
		event.TransactionHash,
		event.LogIndex,
		event.BlockNumber,
		event.ContractAddress,
		event.FromAddress,
		event.ToAddress,
		valueStr,
	)
	if err != nil {
		return err
	}

	return nil
}

// QueryTransferEventsByAddress 是 Day 5 API 背后的查询函数。
// 一个地址既可能是转出方，也可能是接收方，所以需要同时查 from_address 和 to_address。
func QueryTransferEventsByAddress(
	ctx context.Context,
	db *sql.DB,
	address string,
	contractAddress string,
	limit int,
) ([]TransferEvent, error) {
	var query string
	var args []interface{}

	// 不传 contract 时保持原查询；传了 contract 时额外限制 token 合约地址。
	if contractAddress == "" {
		query = `
            SELECT
                transaction_hash,
                log_index,
                block_number,
                contract_address,
                from_address,
                to_address,
                value
            FROM transfer_events
            WHERE from_address = $1 OR to_address = $1
            ORDER BY block_number DESC
            LIMIT $2
        `
		args = []interface{}{address, limit}
	} else {
		query = `
            SELECT
                transaction_hash,
                log_index,
                block_number,
                contract_address,
                from_address,
                to_address,
                value
            FROM transfer_events
            WHERE (from_address = $1 OR to_address = $1)
              AND contract_address = $2
            ORDER BY block_number DESC
            LIMIT $3
        `
		args = []interface{}{address, contractAddress, limit}
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []TransferEvent

	for rows.Next() {
		var event TransferEvent
		var valueStr string

		err := rows.Scan(
			&event.TransactionHash,
			&event.LogIndex,
			&event.BlockNumber,
			&event.ContractAddress,
			&event.FromAddress,
			&event.ToAddress,
			&valueStr,
		)
		if err != nil {
			return nil, err
		}

		// 数据库里的 numeric 先按字符串读出，再转回 big.Int。
		value := new(big.Int)
		_, ok := value.SetString(valueStr, 10)
		if !ok {
			return nil, fmt.Errorf("invalid numeric value from database: %s", valueStr)
		}

		event.Value = value
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
