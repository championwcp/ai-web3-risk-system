# Week 1 任务拆解

本周主题：

`链上数据 -> PostgreSQL -> 基础查询`

## Week 1 目标

- 连接 Sepolia RPC / WebSocket
- 获取最新区块
- 解析 ERC20 `Transfer` 事件
- 将事件写入 PostgreSQL
- 提供地址转账查询 API

## 本周需要掌握的知识点

- 区块、交易、日志的基本概念
- RPC / WebSocket 的作用
- ERC20 `Transfer` 事件结构
- PostgreSQL 基础建表与插入
- Go 中最小可运行程序结构

## Day 1

目标：

连接 Sepolia RPC，并获取最新区块高度。

当前状态：

已完成。

今天实际完成：

- 在 `ai-web3-risk-system/day1-go-rpc` 初始化了最小 Go 模块
- 使用 `go-ethereum/ethclient` 连接 Sepolia RPC
- 成功获取并打印最新区块高度
- 将 RPC 地址改为从环境变量 `SEPOLIA_RPC_URL` 读取
- 为 `client.BlockNumber(...)` 增加了超时控制
- 能用自己的话解释 `client` 与 `ctx` 的职责区别

你需要掌握的知识点：

- RPC 是什么
- 区块高度是什么
- Go 程序如何连接外部服务

你需要自己完成的部分：

- 申请或配置一个可用的 Sepolia RPC 地址
- 自己写出第一个最小 Go 程序

完成标志：

- 程序可以成功打印最新区块高度

复盘：

- 这一节的重点不是写很多代码，而是先打通“配置 -> 连接 -> 请求 -> 返回结果”的最小闭环
- 当前已经完成 Day 1，可以在明天直接进入 Day 2，开始理解区块、交易、日志三者关系

## Day 2

目标：

理解区块、交易、日志，以及 ERC20 `Transfer` 事件在日志中的位置。

明天的起步动作：

- 先不用写新代码
- 先用自己的话说明区块、交易、日志三者的关系
- 再进入 ERC20 `Transfer` 事件为什么会出现在日志里的解释

你需要掌握的知识点：

- 区块包含什么
- 一笔交易和日志的关系
- 为什么 `Transfer` 是事件日志

你需要自己完成的部分：

- 自己画一张简单数据流图
- 自己解释一次 `Transfer` 事件字段含义

完成标志：

- 能用自己的话说清楚区块、交易、日志、事件的关系

当前状态：

已完成核心理解部分。

今天实际完成：

- 用自己的话解释了区块、交易、日志三者关系
- 理解了为什么外部程序需要通过日志而不是只看区块来获取 ERC20 转账细节
- 明确了 `Transfer` 事件的核心字段是 `from`、`to`、`value`
- 理解了 `topics[0] / topics[1] / topics[2] / data` 在 `Transfer` 事件中的分工
- 理解了为什么 `from`、`to` 要从 `topics` 中还原地址，`value` 要从 `data` 中还原成大整数

复盘：

- Day 2 的重点是建立事件日志的心智模型，而不是一开始就写很多代码
- 当前已经可以进入 Day 3 和 Day 4 的工程设计部分

## Day 3

目标：

设计 PostgreSQL 最小表结构，用于保存转账事件。

你需要掌握的知识点：

- 表结构设计最小原则
- 哪些字段必须存
- 主键和索引的最小考虑

你需要自己完成的部分：

- 自己先写一版 `CREATE TABLE`
- 自己思考地址查询需要哪些索引

完成标志：

- 有一张可执行的最小建表 SQL

当前状态：

已完成第一版设计。

今天实际完成：

- 列出了 `transfer_events` 最小字段清单
- 设计了 `transaction_hash + log_index` 作为联合主键候选
- 为字段选择了第一版 PostgreSQL 类型
- 写出了第一版 `CREATE TABLE transfer_events (...)`
- 明确了第一版索引优先考虑 `from_address`、`to_address`、`contract_address`

复盘：

- Day 3 的重点不是一步到位设计最完整的表，而是先为“唯一保存”和“高频查询”建立最小骨架
- 当前已经具备进入 Day 4 的数据库入库设计基础

## Day 4

目标：

把一批 ERC20 `Transfer` 事件写入 PostgreSQL。

你需要掌握的知识点：

- 日志解析结果如何映射到数据库字段
- Go 里如何连接 PostgreSQL 并执行插入

你需要自己完成的部分：

- 自己完成插入逻辑的第一版
- 自己验证数据库里确实写进了数据

完成标志：

- 数据库中存在可查询的转账记录

当前状态：

已完成。

今天实际完成：

- 设计了 `TransferEvent` 结构体及字段类型
- 设计了 `InsertTransferEvent(ctx, db, event)` 的第一版函数签名与执行流程
- 明确了 `Value *big.Int` 在入库前需要转换为十进制字符串
- 写出了 `INSERT INTO transfer_events (...) VALUES (...)` 的第一版 SQL
- 设计并补全了 `ParseTransferEvent(rawLog)` 的第一版解析骨架
- 完成了 `from`、`to` 地址从 `topics` 中还原，以及 `value` 从 `data` 中还原为 `*big.Int` 的思路与第一版代码
- 明确了解析前要做的最小安全检查：`topics` 数量、事件签名、`data` 长度
- 新建了正式工程目录 `ai-web3-risk-system/stage1-indexer-go`
- 完成本地 PostgreSQL 环境准备、数据库连接和 `transfer_events` 表验证
- 写出了 `ProcessTransferLog(...)`，把“识别 -> 解析 -> 入库”串成一条主流程
- 用 `buildMockTransferLog` 构造模拟日志，并成功插入 PostgreSQL
- 通过 `go run .` 验证程序成功输出 `mock transfer log processed and inserted successfully`
- 把真实 Sepolia 历史日志查询接入正式工程
- 用真实 `types.Log` 调用 `ProcessTransferLog(...)`
- 通过 `go run .` 实际验证程序成功输出 `first real transfer log processed and inserted successfully`

当前还差：

- 在当前入库结果基础上进入 Day 5，开始提供按地址查询的最小 API

复盘：

- Day 4 已经从“概念设计 -> mock 闭环 -> 真实链上日志入库”完整走通
- 下一步不再是继续补入库，而是进入 Day 5 的查询 API

## Day 5

目标：

提供一个最小 API，支持按地址查询转账记录。

你需要掌握的知识点：

- API 的最小职责
- 查询参数和返回结构怎么设计
- 如何验证接口是否可用

你需要自己完成的部分：

- 自己写第一版查询接口
- 自己用 Postman 或 curl 做验证

完成标志：

- 可以通过接口查到某地址的转账记录

当前状态：

已完成最小 API 验证。

今天实际完成：

- 设计并实现了 `QueryTransferEventsByAddress(ctx, db, address)`
- 查询条件覆盖 `from_address = $1 OR to_address = $1`
- 查询结果按 `block_number DESC` 排序，并限制返回 20 条
- 将数据库中的 `numeric` 先扫描为字符串，再转换回 `*big.Int`
- 通过 `go run .` 验证可以按真实链上日志中的地址查回转账记录
- 设计并实现了 `TransferEventResponse`
- 设计并实现了 `TransfersHandler(db)`
- 实现了 `GET /transfers?address=...`
- 验证了缺少 `address` 参数时返回 `400`
- 验证了传入真实地址时 API 返回 JSON 转账记录

当前还差：

- Week 1 可以进入收尾复盘
- 后续可以考虑整理工程结构，把 indexer 和 API 逻辑拆分到更清晰的文件中

复盘：

- Day 5 不应该一上来就写 HTTP API，先把“地址 -> 数据库查询 -> 结果对象”打通是更稳的路径
- 当前已完成“地址 -> HTTP API -> 数据库查询 -> JSON 返回”的最小闭环

## 本周验收标准

完成以下内容即视为 Week 1 完成：

- 能连接 Sepolia
- 能拿到区块或日志数据
- 能识别并解析 ERC20 `Transfer`
- 能落库
- 能查询

## 本周最大风险

- RPC 配置不通，导致第一步卡住
- 事件解析概念不清，导致代码虽然能跑但不理解
- 一开始任务过大，导致迟迟没有闭环

## 卡住时应该如何提问

提问时请尽量带上：

- 报错信息
- 关键代码片段
- 当前预期
- 实际结果
