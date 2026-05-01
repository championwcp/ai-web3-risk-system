package main

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// ERC20 Transfer 事件签名：
// Transfer(address indexed from, address indexed to, uint256 value)
var transferEventSigHash = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

// ErrUnsupportedTransferLog 表示“这条日志像 Transfer，但不是当前解析器支持的 ERC20 布局”。
// 外层可以用 errors.Is 判断它，然后把这类日志统计为 skipped。
var ErrUnsupportedTransferLog = errors.New("unsupported transfer log")

// ParseTransferEvent 把一条以太坊原始日志解析成项目里的 TransferEvent。
// 当前只支持标准 ERC20 Transfer：from/to 在 topics，value 在 data。
func ParseTransferEvent(rawLog types.Log) (TransferEvent, error) {
	// topics[0] 是事件签名，topics[1]/[2] 是 indexed 的 from/to。
	if len(rawLog.Topics) < 3 {
		return TransferEvent{}, errors.New("invalid transfer log: not enough topics")
	}

	if rawLog.Topics[0] != transferEventSigHash {
		return TransferEvent{}, errors.New("not a transfer event")
	}

	if len(rawLog.Data) != 32 {
		return TransferEvent{}, fmt.Errorf("%w: invalid data length got %d, want 32", ErrUnsupportedTransferLog, len(rawLog.Data))
	}

	// indexed 地址会被放进 32 字节 topic 槽位，真正的地址是最后 20 字节。
	fromAddress := common.BytesToAddress(rawLog.Topics[1].Bytes()[12:]).Hex()
	toAddress := common.BytesToAddress(rawLog.Topics[2].Bytes()[12:]).Hex()

	// value 是 ABI 编码后的 uint256，在 Go 里用 big.Int 承接，避免普通整数溢出。
	value := new(big.Int).SetBytes(rawLog.Data)

	event := TransferEvent{
		TransactionHash: rawLog.TxHash.Hex(),
		LogIndex:        uint64(rawLog.Index),
		BlockNumber:     rawLog.BlockNumber,
		ContractAddress: rawLog.Address.Hex(),
		FromAddress:     fromAddress,
		ToAddress:       toAddress,
		Value:           value,
	}

	return event, nil
}

// buildMockTransferLog 保留为本地学习/测试样本。
// 当前主流程已经改为使用真实 Sepolia 日志。
func buildMockTransferLog() types.Log {
	from := common.HexToAddress("0x1111111111111111111111111111111111111111")
	to := common.HexToAddress("0x2222222222222222222222222222222222222222")
	contract := common.HexToAddress("0x3333333333333333333333333333333333333333")
	value := big.NewInt(123456789)

	return types.Log{
		Address: contract,
		Topics: []common.Hash{
			transferEventSigHash,
			common.BytesToHash(common.LeftPadBytes(from.Bytes(), 32)),
			common.BytesToHash(common.LeftPadBytes(to.Bytes(), 32)),
		},
		Data:        common.LeftPadBytes(value.Bytes(), 32),
		BlockNumber: 12345678,
		TxHash:      common.HexToHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		Index:       0,
	}
}
