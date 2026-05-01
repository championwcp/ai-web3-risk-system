# Week 2 任务拆解

本周主题：

`单条日志处理 -> 小区块范围多日志处理`

## Week 2 目标

- 将 `FetchFirstRealTransferLog` 扩展为 `FetchTransferLogs`
- 支持一次查询一个小区块范围内的多条 ERC20 `Transfer` 日志
- 循环调用 `ProcessTransferLog`，将多条日志写入 PostgreSQL
- 输出批量处理统计信息
- 使用现有 `GET /transfers?address=...` API 验证多条结果

## 本周需要掌握的知识点

- `FilterLogs` 返回多条日志时如何处理
- 批量处理中的成功、失败、跳过统计
- 幂等写入和 `ON CONFLICT DO NOTHING` 的作用
- 为什么先做小范围历史日志，再考虑 WebSocket 实时监听

## 当前状态

Week 2 已完成，准备进入 Week 3。

## 已完成

- 将单条日志查询函数升级为 `FetchTransferLogs`
- 新增 `ProcessTransferLogsResult`
- 新增 `ProcessTransferLogs`
- 主流程已从处理单条 `types.Log` 改为处理 `[]types.Log`
- 已通过 `go test ./...` 编译检查
- 已通过 `go run .` 验证小区块范围批量查询和处理
- 实际查询区块 `7900000` 到 `7900009`，共获取 `978` 条 Transfer 签名日志
- 批量处理结果为 `total=978 succeeded=963 skipped=15 failed=0`
- 新增 `skipped` 统计，用于区分“不支持的 Transfer 日志布局”和“真正处理失败”
- 使用 `ErrUnsupportedTransferLog`、`fmt.Errorf("%w", ...)` 和 `errors.Is(...)` 替代字符串匹配错误
- 已将区块范围从硬编码改为环境变量 `FROM_BLOCK` / `TO_BLOCK`
- 已为 `GET /transfers?address=...` 增加 `limit` 参数，支持控制返回条数
- 已验证 `limit=10` 返回 `200`，`limit=abc`、`limit=0`、`limit=1000` 返回 `400`
- 已为 `GET /transfers?address=...` 增加可选 `contract` 参数，支持按 ERC20 合约地址过滤
- 已验证不传 `contract` 返回正常记录，传正确 `contract` 返回正常记录，传不存在的 `contract` 返回空数组 `[]`
- 已将 `LoadBlockRangeFromEnv` 从 `indexer.go` 拆分到 `config.go`，让配置读取和日志索引职责分开

## 当前还差

- Week 3 优先处理启动流程整理：把“抓日志入库”和“启动 API”拆成更清晰的运行模式
- 后续再考虑分页游标、断点续扫、更多 API 查询条件

## 完成标志

做到下面这些，就算 Week 2 第一阶段完成：

- 程序能查询到多条 `Transfer` 日志：已完成
- 批量处理统计输出正常：已完成
- PostgreSQL 中能看到多条真实转账记录：已完成
- API 能按地址返回转账记录：已完成

## Week 2 复盘

- 最重要的能力提升：对索引服务有了基本认知，已经从数据库搭建、拉取日志、解析、统计、入库到 API 查询打通了完整流程。
- 当前系统短板：只有一个查询接口，还不是一个完备系统；启动流程也还混合了“抓日志”和“启动 API”两类职责。
- Week 3 方向：优先做启动流程整理，把 `go run . index` 和 `go run . serve` 这类运行模式逐步拆清楚。

## 暂不做

- WebSocket 实时监听
- 断点续扫
- 大范围历史同步
- 多链支持
