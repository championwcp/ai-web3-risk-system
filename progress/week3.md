# Week 3 任务拆解

本周主题：

`整理启动流程：index / serve 分离`

## Week 3 目标

把当前混在 `main.go` 里的启动流程拆清楚。

当前程序启动后会同时做：

```text
连接 RPC -> 连接数据库 -> 抓日志入库 -> 本地查询验证 -> 启动 API
```

Week 3 目标是逐步整理成：

```text
go run . index
```

只负责：

```text
抓取链上日志 -> 解析 -> 入库 -> 输出统计 -> 退出
```

```text
go run . serve
```

只负责：

```text
启动 HTTP API 服务
```

## 本周需要掌握的知识点

- Go 程序如何读取命令行参数
- `os.Args` 的基本用法
- 为什么一个程序可以有不同运行模式
- `main.go` 如何做最小职责拆分
- 如何在不改变行为的前提下做小步重构

## 当前状态

未开始。

## 本周第一步

先不写完整代码。

第一步先理解并设计：

- 当前 `main.go` 做了哪些事情
- 哪些属于 `index` 模式
- 哪些属于 `serve` 模式
- 第一版函数边界应该如何划分

## 计划拆解

- Step 1：阅读 `main.go`，列出当前启动流程
- Step 2：理解 `os.Args`，设计命令格式
- Step 3：把抓日志入库逻辑抽成 `runIndex`
- Step 4：把 API 启动逻辑抽成 `runServe`
- Step 5：让 `main` 根据参数选择运行模式
- Step 6：验证 `go run . index`
- Step 7：验证 `go run . serve`
- Step 8：更新 README 和文档

## 暂不做

- 不做断点续扫
- 不做 WebSocket
- 不做后台任务调度
- 不做复杂 CLI 框架
- 不引入 Cobra 等第三方命令行库
