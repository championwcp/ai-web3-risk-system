# AI-Driven Cross-Chain Risk & Intelligence System

一个以学习驱动为核心的 AI x Web3 项目。

这个项目的目标不是一开始就做复杂系统，而是通过一个真实可运行、可逐步扩展、可开源展示的项目，系统学习以下能力：

- Web3 基础
- Go 基础
- Rust 基础
- Solidity 基础
- AI 应用开发基础
- 后端系统设计基础

## 项目定位

这是一个长期推进的学习型项目，强调：

- 边学边做
- 先建立最小可运行闭环
- 只学习当前任务真正依赖的知识
- 每一步都保留文档、进度和决策记录

当前阶段不追求：

- 一开始就做复杂 Agent 编排
- 一开始就同时覆盖多链、跨链、风控、RAG、微调
- 一开始就写大量抽象层
- 一开始就追求工程完美或生产级性能

## 项目主题

项目名称：

`AI-Driven Cross-Chain Risk & Intelligence System`

中文描述：

`AI 驱动的跨链风险分析与资产智能系统`

一句话说明：

基于链上数据索引、风险计算和 LLM 分析能力，构建一个面向 DeFi / 跨链场景的学习型系统项目。

## 当前开发策略

采用分阶段推进：

1. 阶段 1：链上数据采集与入库
2. 阶段 2：查询 API 与地址分析
3. 阶段 3：风险规则与基础风控
4. 阶段 4：接入 LLM 做解释型分析
5. 阶段 5：增加 Agent / Tool Calling
6. 阶段 6：完善 README、架构图、开源展示

## 当前阶段

当前聚焦：`阶段 1：链上数据采集与入库`

第一阶段目标：

- 连接 Sepolia RPC / WebSocket
- 监听新区块
- 解析 ERC20 `Transfer` 事件
- 将事件写入 PostgreSQL
- 提供最基础的地址转账查询 API

第一阶段暂时不做：

- AI Agent
- RAG
- 多链聚合
- 跨链风控
- 深入 Solidity 合约开发
- 高级性能优化

## 第一阶段建议技术栈

- Indexer：Go
- API：Go
- 数据库：PostgreSQL
- 网络：Sepolia
- AI：第一阶段暂不接入

说明：

第一阶段优先选择 Go，是为了降低同时学习多门语言的负担，先完成最小闭环。后续可以把 Indexer 重写为 Rust，作为阶段增强项。

## 新线程继续项目

如果在新的 Codex / AI 线程中继续本项目，请先让助手读取：

```text
docs/current-status.md
docs/codex-prompt.md
progress/week3.md
progress/backlog.md
```

然后再继续当前任务。

## 当前进度

Week 1 和 Week 2 已完成，当前准备进入 Week 3。

已跑通的最小闭环：

```text
Sepolia 真实日志
-> 解析 ERC20 Transfer
-> 写入 PostgreSQL
-> 按地址查询 API
-> 返回 JSON
```

当前已完成：

- 连接 Sepolia RPC
- 查询真实历史 `Transfer` 日志
- 解析 `topics[0] / topics[1] / topics[2] / data`
- 将 `Transfer` 事件写入 PostgreSQL
- 提供 `GET /transfers?address=...` 最小查询 API
- 支持小区块范围批量日志处理
- 支持 `FROM_BLOCK` / `TO_BLOCK` 配置扫描区块范围
- API 支持 `limit` 参数
- API 支持可选 `contract` 参数

当前暂未完成：

- WebSocket 实时监听新区块
- 断点续扫
- 启动流程拆分
- 更完整的分页和地址分析能力
- 工程测试补充

## 建议目录结构

```text
aiWeb3learning/
├── README.md
├── docs/
│   ├── project-overview.md
│   ├── architecture.md
│   ├── learning-log.md
│   ├── roadmap.md
│   ├── weekly-progress.md
│   └── decisions.md
├── progress/
│   ├── week1.md
│   └── backlog.md
├── indexer/
├── api/
├── scripts/
└── examples/
```

当前阶段 1 正式工程位于：

```text
ai-web3-risk-system/stage1-indexer-go/
```

当前工程文件职责：

```text
stage1-indexer-go/
├── main.go        # 程序启动、环境变量、数据库连接、API 启动
├── models.go      # TransferEvent 和 JSON response 数据结构
├── parser.go      # ERC20 Transfer 日志解析
├── indexer.go     # 真实日志查询与日志处理编排
├── config.go      # 配置读取，例如 FROM_BLOCK / TO_BLOCK
├── repository.go  # PostgreSQL 写入和查询
└── api.go         # HTTP handler
```

Day 1 的独立练习样本保留在：

```text
ai-web3-risk-system/day1-go-rpc/
```

## 本地运行

### 前置条件

- Go
- PostgreSQL
- Sepolia RPC URL

### 环境变量

在 `ai-web3-risk-system/stage1-indexer-go/.env` 中准备：

```text
SEPOLIA_RPC_URL=你的 Sepolia RPC 地址
DATABASE_URL=postgres://用户名:密码@localhost:5432/数据库名?sslmode=disable
```

### 数据表

当前最小表结构：

```sql
CREATE TABLE transfer_events (
    transaction_hash varchar(66),
    block_number bigint,
    contract_address varchar(42),
    from_address varchar(42),
    to_address varchar(42),
    value numeric,
    log_index integer,
    PRIMARY KEY (transaction_hash, log_index)
);
```

### 启动

进入阶段 1 工程：

```powershell
cd D:\aiproject\aiWeb3learning\ai-web3-risk-system\stage1-indexer-go
go run .
```

当前程序会先做一次学习用的启动流程：

1. 连接 Sepolia
2. 查询一条真实历史 `Transfer` 日志
3. 解析并写入 PostgreSQL
4. 按地址做一次本地查询验证
5. 启动 HTTP API 服务

默认服务地址：

```text
http://localhost:8080
```

### 查询 API

```text
GET /transfers?address=0x...
```

PowerShell 示例：

```powershell
Invoke-RestMethod "http://localhost:8080/transfers?address=0xEe12063a08584b501d7F49Aa3751841EBA07e716"
```

返回示例：

```json
[
  {
    "transaction_hash": "0x...",
    "log_index": 0,
    "block_number": 7900000,
    "contract_address": "0x...",
    "from_address": "0x...",
    "to_address": "0x...",
    "value": "5005000"
  }
]
```

## 当前仓库文档

- [当前状态](D:\aiproject\aiWeb3learning\docs\current-status.md)
- [项目概览](D:\aiproject\aiWeb3learning\docs\project-overview.md)
- [系统架构](D:\aiproject\aiWeb3learning\docs\architecture.md)
- [学习记录](D:\aiproject\aiWeb3learning\docs\learning-log.md)
- [路线图](D:\aiproject\aiWeb3learning\docs\roadmap.md)
- [周进度](D:\aiproject\aiWeb3learning\docs\weekly-progress.md)
- [技术决策](D:\aiproject\aiWeb3learning\docs\decisions.md)
- [Week 1 任务拆解](D:\aiproject\aiWeb3learning\progress\week1.md)
- [Week 2 任务拆解](D:\aiproject\aiWeb3learning\progress\week2.md)
- [Week 3 任务拆解](D:\aiproject\aiWeb3learning\progress\week3.md)
- [Backlog](D:\aiproject\aiWeb3learning\progress\backlog.md)

## 如何推进这个项目

每次推进时，优先回答下面 5 个问题：

1. 当前阶段目标是否变化？
2. 本周任务是否完成？
3. 新学到的知识点是什么？
4. 当前最大的阻塞点是什么？
5. 下一步最小任务是什么？

## 当前下一步

进入 Week 3：

- 整理启动流程
- 将“抓日志入库”和“启动 API”拆成两个运行模式
- 目标形式：`go run . index` 和 `go run . serve`
