# 系统架构

## 当前架构目标

第一阶段只建立一个最小可运行闭环：

`链上数据 -> Indexer -> PostgreSQL -> API -> 查询结果`

## 当前架构图

```text
Sepolia RPC
        |
        v
stage1-indexer-go
  - FetchFirstRealTransferLog
  - ParseTransferEvent
  - ProcessTransferLog
        |
        v
   PostgreSQL
  transfer_events
        |
        v
API (Go net/http)
  GET /transfers?address=...
        |
        v
JSON Transfer Records
```

## 当前模块职责

### 1. Indexer

职责：

- 连接 Sepolia RPC
- 查询历史日志事件
- 解析 ERC20 `Transfer` 事件
- 将标准化后的转账数据写入 PostgreSQL

第一阶段暂不追求：

- 高并发优化
- 多链并行采集
- 重试框架完善
- 复杂任务调度
- WebSocket 实时监听

### 2. PostgreSQL

职责：

- 存储转账事件
- 支持按地址、区块范围、时间等维度查询

第一阶段先从单表或少量表开始，不做过度设计。

### 3. API

职责：

- 提供 `GET /transfers?address=...`
- 根据地址查询转账记录
- 以 JSON 返回转账结果

第一阶段暂不追求：

- 认证授权
- 复杂聚合查询
- 多版本接口设计
- 分页和高级过滤

## 第一阶段建议数据流

当前已跑通的数据流：

1. Go 程序连接 Sepolia RPC
2. 使用 `FilterLogs` 查询一小段历史区块中的 ERC20 `Transfer` 日志
3. 通过 `topics[0]` 判断事件类型
4. 从 `topics[1]`、`topics[2]` 解析 `from`、`to`
5. 从 `data` 解析 `value`
6. 组装为 `TransferEvent`
7. 写入 PostgreSQL 的 `transfer_events` 表
8. HTTP API 按地址查询并返回 JSON

Week 2 计划扩展的数据流：

1. 输入一个区块范围
2. 查询范围内多条 `Transfer` 日志
3. 循环调用 `ProcessTransferLog`
4. 批量写入 PostgreSQL
5. 通过 API 查询多条结果

## 第一阶段重点关注的学习点

- 区块、交易、日志的关系
- RPC / WebSocket 的区别
- 事件日志和 ABI 解码的基础概念
- 数据库存储与查询的最小设计
- Go 项目中最小可运行服务的组织方式

## 当前架构决策

- 先单体推进，不拆微服务
- 先单链推进，不做多链
- 先做事件采集，不做复杂风险引擎
- 先做基础查询，不做 AI 分析
- 先用历史日志查询跑通稳定闭环，再考虑 WebSocket 实时监听
- 当前代码保持一个 Go module，按文件职责拆分，不引入多包结构

后续每次架构变化，都需要同步更新本文件。
