# 每周进度

用于记录每周完成了什么、卡住了什么、下一周准备做什么。

## 记录模板

### Week X

本周完成：

- 

本周卡点：

- 

学到的内容：

- 

下周计划：

- 

## 当前记录

### Week 1

本周完成：

- 初始化项目文档结构
- 明确第一阶段范围和边界
- 明确 Week 1 聚焦“链上数据 -> 数据库 -> 查询”
- 完成 Day 1 最小闭环：连接 Sepolia RPC 并成功打印最新区块高度
- 在 `ai-web3-risk-system/day1-go-rpc` 中完成第一个 Go 练习程序
- 完成从硬编码 RPC 地址到环境变量读取的改造
- 为最新区块请求增加了超时控制
- 完成 Day 2 的核心概念理解：区块、交易、日志、ERC20 `Transfer` 事件之间的关系
- 完成 Day 3 的第一版数据库设计：`transfer_events` 表结构、主键候选、索引候选
- 完成 Day 4 的大部分设计工作：`TransferEvent` 结构体、`InsertTransferEvent(...)`、`ParseTransferEvent(...)` 的第一版设计与解析骨架
- 完成阶段 1 正式工程目录 `stage1-indexer-go`
- 完成本地 PostgreSQL 环境准备和连接验证
- 使用 `buildMockTransferLog` 跑通了“识别 -> 解析 -> 入库”本地闭环
- 已把真实 Sepolia Transfer 日志接入 `ProcessTransferLog`，并成功完成真实日志入库验证
- 完成 Day 5 的第一步：按地址查询数据库的函数 `QueryTransferEventsByAddress`
- 验证可以按真实 `from_address` 查回已入库的转账记录
- 完成最小 HTTP API：`GET /transfers?address=...`
- 验证 API 能返回 JSON 转账记录，并能在缺少 `address` 时返回 `400`

本周卡点：

- 今天暴露出协作问题：在进入 Day 4 收尾时，没有先明确当前工程目录、数据库环境、建表状态等实操前置条件

学到的内容：

- 学习型项目需要把阶段目标、进度、学习点和决策过程一起沉淀
- 第一阶段不应该同时引入 AI、跨链、多语言重写等复杂目标
- Go 可以用 `go-ethereum` 的 `ethclient` 很快建立最小 RPC 读取程序
- RPC 地址不应直接写在代码里，更适合通过环境变量管理
- 单次 RPC 请求适合加上 `context.WithTimeout`，避免网络调用无限等待
- 事件日志处理可以分成“先识别事件类型，再解码参数，再组装对象”
- `indexed` 参数与非 `indexed` 参数会分别落在 `topics` 和 `data`
- `uint256` 这类链上数值在 Go 中更适合先用 `*big.Int` 承接，再在入库前转换
- 当任务从设计走向运行验证时，必须先明确工程位置和环境准备状态，再安排后续代码执行
- mock 日志是连接“函数设计”和“真实链上日志”之间的一个很有价值的中间验证步骤

下周计划：

- Week 1 已完成最小闭环，下一步进入 Week 2
- Week 2 目标：从“处理一条真实日志”升级为“处理一个小区块范围内的多条 Transfer 日志”
- 将当前 `FetchFirstRealTransferLog` 扩展为返回多条日志的查询函数
- 循环调用 `ProcessTransferLog`，验证多条日志可写入 PostgreSQL
- 保持 API 不变，验证同一个地址可以查到多条记录

### Week 2

本周完成：

- 将单条日志查询函数升级为 `FetchTransferLogs`
- 新增 `ProcessTransferLogsResult`，用于记录批量处理总数、成功数、跳过数和失败数
- 新增 `ProcessTransferLogs`，循环处理一个区块范围内的多条 `Transfer` 日志
- 主流程已从单条日志处理升级为批量日志处理
- 通过 `go test ./...` 完成编译检查
- 通过 `go run .` 完成真实运行验证，实际查询到 `978` 条 Transfer 签名日志
- 批量处理结果为 `total=978 succeeded=963 skipped=15 failed=0`
- 保持现有 API 服务可启动，并继续支持 `GET /transfers?address=...`
- 将不支持的 Transfer 日志布局归类为 `skipped`，避免误认为程序失败
- 将日志查询区块范围改为从 `FROM_BLOCK` / `TO_BLOCK` 环境变量读取
- 为 `GET /transfers?address=...` 增加 `limit` 参数，默认 `20`，合法范围 `1` 到 `100`
- 验证 `limit=10` 返回 `200`，`limit=abc`、`limit=0`、`limit=1000` 返回 `400`
- 为 `GET /transfers?address=...` 增加可选 `contract` 参数，用于按 ERC20 合约地址过滤
- 验证不传 `contract`、传正确 `contract`、传不存在 `contract` 三种情况均符合预期
- 将区块范围配置读取逻辑拆到 `config.go`，让 `indexer.go` 更专注于日志查询和处理

本周卡点：

- 真实链上日志里，`topics[0]` 相同并不一定代表一定是标准 ERC20 Transfer 布局
- 部分日志虽然事件签名也是 `Transfer(address,address,uint256)`，但 `data` 长度为 `0`，更像 ERC721 等不同 indexed 布局，需要跳过而不是按失败处理
- SQL 占位符 `$1`、`$2` 与 `QueryContext(ctx, query, ...)` 中 query 后面的参数列表对应，需要注意和 Go 函数完整参数位置区分

学到的内容：

- 从单条处理升级到批量处理时，第一步不是引入复杂调度，而是把返回值从一条日志改成日志切片
- 批量处理需要最小统计信息，否则很难判断到底处理了多少、失败了多少
- `skipped` 和 `failed` 的含义不同：`skipped` 表示当前阶段不支持但可预期的数据形态，`failed` 表示真正异常
- Go 里可以用哨兵错误表达一类可识别错误，再用 `%w` 包装上下文，用 `errors.Is` 在外层判断错误类别
- API handler 负责解析和校验用户输入，repository 只接收已经整理好的查询参数
- `http.Error` 返回的文本会成为 HTTP 响应体，所以 `invalid limit` 等错误信息来自 handler 中显式写出的错误文本
- `contract` 和 `address` 含义不同：`address` 是用户钱包地址，`contract` 是 token 合约地址
- SQL 中混用 `OR` 和 `AND` 时要用括号明确优先级，例如 `(from_address = $1 OR to_address = $1) AND contract_address = $2`
- Go 同一个 `package main` 下的多个 `.go` 文件可以共享函数，拆文件主要是为了组织职责，不等于拆 package

下周计划：

- Week 3 优先做启动流程整理，把“抓日志入库”和“启动 API”拆成两个更清晰的运行模式
- 暂时不急着继续堆 API 参数，也不急着上 WebSocket 或 AI
- 继续保持教学优先：每次写入代码前先明确写在哪个文件、为什么写、运行前需要什么环境依赖

Week 2 复盘：

- 最重要的能力提升是对索引服务形成了基本认知，已经打通“数据库搭建 -> 拉取日志 -> 解析 -> 统计 -> 入库 -> API 查询”的完整流程
- 当前系统最明显的短板是接口和运行模式都还比较简单，还不是一个完备系统
- Week 3 选择启动流程整理，因为先拆清运行职责，后面做断点续扫和更稳定的 indexer 会更自然
