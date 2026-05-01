package main

import "math/big"

type TransferEvent struct {
	TransactionHash string
	LogIndex        uint64
	BlockNumber     uint64
	ContractAddress string
	FromAddress     string
	ToAddress       string
	Value           *big.Int
}

// TransferEventResponse 是 API 对外返回的 JSON 结构。
// Value 用字符串返回，避免 uint256 这类大整数在 JSON/前端里丢精度。
type TransferEventResponse struct {
	TransactionHash string `json:"transaction_hash"`
	LogIndex        uint64 `json:"log_index"`
	BlockNumber     uint64 `json:"block_number"`
	ContractAddress string `json:"contract_address"`
	FromAddress     string `json:"from_address"`
	ToAddress       string `json:"to_address"`
	Value           string `json:"value"`
}
