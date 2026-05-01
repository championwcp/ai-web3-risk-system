# 学习记录

用于持续记录本项目中的学习收获、未理解点和待补课内容。

## 记录规则

每次记录尽量使用下面的格式：

### 日期

`YYYY-MM-DD`

### 本次学到的内容

- 学会了什么
- 理解了什么

### 仍然不清楚的点

- 哪个概念还不明白
- 哪段代码还看不懂

### 后续要补课的内容

- 下次需要补的知识点

## 当前记录

### 2026-04-12

本次学到的内容：

- 明确了项目第一阶段的目标是“链上数据采集 -> PostgreSQL -> 基础查询 API”
- 明确了当前应优先选择 Go 建立最小闭环，而不是同时学习 Go 和 Rust
- 理清了后续文档需要长期维护，不能只写代码不记录过程
- 在 `ai-web3-risk-system/day1-go-rpc` 中完成了最小 Go 程序，能够连接 Sepolia RPC 并获取最新区块高度
- 理解了 `ethclient` 是与以太坊节点通信的客户端，本地不需要自己先运行节点
- 学会了通过环境变量 `SEPOLIA_RPC_URL` 管理 RPC 配置，而不是把地址硬编码在代码里
- 学会了使用 `context.WithTimeout` 为 `client.BlockNumber(...)` 这类 RPC 请求增加超时控制
- 能区分 `client` 和 `ctx` 的职责：`client` 负责通信，`ctx` 负责控制单次请求的时间边界

仍然不清楚的点：

- Sepolia 上如何高效获取适合练习的 ERC20 `Transfer` 数据
- ERC20 `Transfer` 事件的 `topics` 与 `data` 在实际日志中的映射细节
- 区块、交易、日志三者在链上数据读取流程中的关系还需要继续巩固
- `godotenv` 和系统环境变量分别适合什么场景，还需要在后续工程实践中进一步理解

后续要补课的内容：

- 区块、交易、日志三者关系
- RPC 与 WebSocket 的使用场景差异
- 事件 ABI 编码与解码基础
- ERC20 `Transfer` 事件为什么最终体现在日志里

### 2026-04-23

本次学到的内容：

- 理清了区块、交易、日志三者的关系，明确了 ERC20 `Transfer` 事件最终体现在日志里
- 理解了 `Transfer(address indexed from, address indexed to, uint256 value)` 中 `from`、`to`、`value` 的职责划分
- 理解了 `topics[0]` 是事件签名哈希，`topics[1]` 和 `topics[2]` 对应 `from`、`to`，`data` 对应 `value`
- 完成了 `transfer_events` 的最小表结构设计，明确了联合主键候选为 `transaction_hash + log_index`
- 明确了第一版索引优先考虑 `from_address`、`to_address`、`contract_address`
- 设计了 `TransferEvent` 结构体，并为字段选择了第一版 Go 类型
- 设计了 `InsertTransferEvent(ctx, db, event)` 的第一版流程，明确了 `Value` 需要先从 `*big.Int` 转成十进制字符串再写入 PostgreSQL `numeric`
- 设计并补全了 `ParseTransferEvent(rawLog)` 的第一版解析思路，完成了 `topics -> address` 和 `data -> big.Int` 的基本映射
- 理解了在解析日志前要先做最小安全检查，包括 `topics` 数量、事件签名、`data` 长度等

仍然不清楚的点：

- 如何在外层主流程中把“识别 Transfer 日志 -> 调用 ParseTransferEvent -> 调用 InsertTransferEvent”稳定串起来
- Sepolia 上如何选择合适的真实日志样本做验证
- `Transfer` 日志的真实监听方式是更适合从历史日志抓取开始，还是直接从 WebSocket 订阅开始

后续要补课的内容：

- 外层日志处理主流程设计
- Go 中 PostgreSQL 驱动与最小连接初始化
- 如何获取和验证一条真实的 ERC20 `Transfer` 日志

### 2026-04-23（协作复盘）

本次学到的内容：

- 当任务开始从“概念设计”走向“真实运行验证”时，必须先检查工程环境是否已经准备完成
- Day 4 的真正前置条件包括：明确当前工程目录、准备 PostgreSQL、创建 `transfer_events` 表、确认依赖与配置
- `day1-go-rpc` 适合作为 Day 1 练习样本，不应默认继续承载 Day 4 / Day 5 的完整工程逻辑

仍然不清楚的点：

- 下一次具体应该采用什么目录结构来承接阶段 1 的 indexer 工程
- PostgreSQL 本地环境准备应该采用什么最小路径

后续要补课的内容：

- 阶段 1 正式工程目录规划
- PostgreSQL 最小安装与建表实操
- 从环境准备到代码验证的最小执行清单

### 2026-04-25

本次学到的内容：

- 在 `ai-web3-risk-system/stage1-indexer-go` 中完成了阶段 1 正式工程的第一版可运行骨架
- 本地 PostgreSQL 环境已准备完成，并能够通过 `DATABASE_URL` 连接
- 完成了 `ProcessTransferLog -> ParseTransferEvent -> InsertTransferEvent` 的外层闭环串联
- 使用 `buildMockTransferLog` 构造出一条模拟 ERC20 `Transfer` 日志，并成功写入 PostgreSQL
- 通过 `go run .` 实际验证程序成功输出 `mock transfer log processed and inserted successfully`
- 理解了 mock 日志验证的价值：先验证解析与入库链路，再继续接入真实链上日志
- 已经把真实 Sepolia 日志查询接入到主流程，并成功把第一条真实 `Transfer` 日志写入 PostgreSQL
- 理解了“先查询历史日志拿到真实 `types.Log`，再交给 `ProcessTransferLog`”是从 mock 过渡到真实链上数据的自然路径

仍然不清楚的点：

- 历史日志抓取和 WebSocket 实时订阅，下一步应该先选哪个路径更适合作为长期 indexer 入口

后续要补课的内容：

- Day 5 查询 API 的最小设计与路由实现

### 2026-04-26

本次学到的内容：

- 进入 Day 5，先从 API 之前的数据库查询函数开始，而不是直接写 HTTP handler
- 完成了 `QueryTransferEventsByAddress(ctx, db, address)`，支持按 `from_address` 或 `to_address` 查询转账记录
- 学会了用 `QueryContext` 查询多行结果，并用 `rows.Scan` 映射回 `TransferEvent`
- 理解了 PostgreSQL `numeric` 读出后可以先用字符串承接，再通过 `big.Int.SetString(..., 10)` 转回 `*big.Int`
- 通过 `go run .` 验证已能查回真实链上日志入库后的转账记录
- 完成了 `TransfersHandler(db)`，实现了 `GET /transfers?address=...`
- 理解了 API 返回 JSON 时需要将 `*big.Int` 转为字符串，避免大整数表达不清
- 通过 HTTP 请求验证了地址查询接口可以返回真实转账记录
- 验证了缺少 `address` 参数时接口会返回 `400`

仍然不清楚的点：

- 当前所有逻辑还集中在 `main.go`，后续需要学习如何拆分文件和模块
- 后续 API 是否需要更多查询条件，例如合约地址、区块范围、分页

后续要补课的内容：

- Go 项目文件拆分方式
- API 查询参数扩展设计
- Week 1 项目复盘与 README 更新

### 2026-04-30

本次学到的内容：

- 将单条真实日志处理升级为小区块范围批量处理，实际在区块 `7900000` 到 `7900009` 查询到 `978` 条 Transfer 签名日志
- 理解了批量处理要统计 `total`、`succeeded`、`skipped`、`failed`，否则很难判断程序到底处理到了什么程度
- 理解了 `topics[0]` 相同只说明事件签名相同，不保证日志布局完全符合标准 ERC20 Transfer
- 观察到真实日志中存在 `data` 长度为 `0` 的 Transfer 签名日志，这类日志当前不适合按 ERC20 `value` 解析，应归类为 `skipped`
- 学会了用 Go 的哨兵错误表达“可预期、可跳过”的错误类别，例如 `ErrUnsupportedTransferLog`
- 学会了用 `fmt.Errorf("%w", err)` 保留错误类别，再用 `errors.Is(err, ErrUnsupportedTransferLog)` 在外层做分类判断

仍然不清楚的点：

- ERC721、ERC1155 等不同 token 标准的 Transfer/TransferSingle 日志布局还有待后续对比学习
- 当前区块范围仍然写在代码里，后续需要学习如何把它变成配置项

后续要补课的内容：

- Go 中配置读取的最小做法，例如从环境变量读取 `FROM_BLOCK` 和 `TO_BLOCK`
- 不同 token 标准的事件日志结构差异
- 批量 indexer 后续如何设计断点续扫
