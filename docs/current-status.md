# 当前状态

本文件是新线程 / 新协作开始时的第一入口。

## 必读文件

新线程继续项目时，必须先读取：

- `docs/current-status.md`
- `docs/codex-prompt.md`
- `progress/week3.md`
- `progress/backlog.md`

如果需要了解历史，再读取：

- `progress/week2.md`
- `docs/weekly-progress.md`
- `docs/learning-log.md`

## 当前阶段

- 当前阶段：阶段 1，链上数据采集与入库
- 当前周次：Week 3
- 当前主题：整理启动流程
- 当前正式工程目录：`ai-web3-risk-system/stage1-indexer-go`

## 当前已完成

- 连接 Sepolia RPC
- 查询真实历史 Transfer 日志
- 解析标准 ERC20 Transfer 日志
- 写入 PostgreSQL
- 批量处理小区块范围内的 Transfer 签名日志
- 支持 `total / succeeded / skipped / failed` 批量统计
- 使用 `ErrUnsupportedTransferLog` 区分可跳过日志和真正失败
- 从 `FROM_BLOCK` / `TO_BLOCK` 环境变量读取扫描区块范围
- 提供 `GET /transfers?address=...` 查询 API
- API 支持 `limit` 参数
- API 支持可选 `contract` 参数
- 已同步到 GitHub 公开仓库

## 当前短板

- 程序启动流程仍然混合了两类职责：
  - 抓取日志并写入数据库
  - 启动 HTTP API 服务
- 当前还不是一个完整系统，只有最小 indexer 和查询 API。
- 还没有断点续扫、WebSocket、生产级任务调度。

## Week 3 目标

Week 3 只做一件主线任务：

> 将“抓日志入库”和“启动 API”拆成两个更清晰的运行模式。

目标运行方式：

```text
go run . index
```

只执行：

```text
连接 RPC -> 查询日志 -> 解析 -> 入库 -> 输出统计 -> 退出
```

```text
go run . serve
```

只执行：

```text
连接数据库 -> 启动 HTTP API
```

## 当前协作原则

- 教学优先，不要默认直接完整实现代码。
- 每次代码任务必须先说明：
  - 本次目标
  - 需要掌握的知识点
  - 应该改哪个文件
  - 用户需要自己完成哪一部分
  - 完成标志
- 如果只是文档整理、注释、格式化、低风险清理，Codex 可以直接执行。
- 进入代码实现前，必须确认当前工程目录和环境依赖。

## 下一步最小任务

不要直接改完整启动流程。

下一步先让用户理解并设计：

- `os.Args` 是什么
- 为什么可以用第一个命令行参数区分 `index` 和 `serve`
- `main.go` 当前混合了哪些职责
- 第一版应该如何拆出两个函数：
  - `runIndex(...)`
  - `runServe(...)`
