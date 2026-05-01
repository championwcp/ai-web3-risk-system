package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// 第一步：连接 Sepolia RPC，用来读取真实链上的历史日志。
	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
		log.Fatal("missing SEPOLIA_RPC_URL")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("failed to connect rpc:", err)
	}
	defer client.Close()

	// 第二步：连接 PostgreSQL，解析后的 Transfer 事件会写入这里，API 也从这里查询。
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("missing DATABASE_URL")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal("failed to ping db:", err)
	}

	// Week 2 启动流程：先抓一小段区块里的真实日志，再批量解析并入库。
	logs, err := FetchTransferLogs(ctx, client)
	if err != nil {
		log.Fatal("failed to fetch transfer logs:", err)
	}

	result := ProcessTransferLogs(ctx, db, logs)
	fmt.Printf(
		"transfer logs processed: total=%d succeeded=%d skipped=%d failed=%d\n",
		result.Total,
		result.Succeeded,
		result.Skipped,
		result.Failed,
	)

	if result.Failed > 0 {
		log.Fatal("some transfer logs failed to process")
	}

	// 启动 API 前做一次本地验证：用第一条日志里的 from 地址反查数据库。
	event, err := ParseTransferEvent(logs[0])
	if err != nil {
		log.Fatal("failed to parse real transfer log for query test:", err)
	}

	events, err := QueryTransferEventsByAddress(ctx, db, event.FromAddress, 20)
	if err != nil {
		log.Fatal("failed to query transfer events by address:", err)
	}

	fmt.Printf("queried %d transfer events for address %s\n", len(events), event.FromAddress)
	for i, item := range events {
		fmt.Printf(
			"[%d] tx=%s block=%d from=%s to=%s value=%s contract=%s logIndex=%d\n",
			i,
			item.TransactionHash,
			item.BlockNumber,
			item.FromAddress,
			item.ToAddress,
			item.Value.String(),
			item.ContractAddress,
			item.LogIndex,
		)
	}

	// Day 5 最小 API：按地址查询转账记录。
	http.HandleFunc("/transfers", TransfersHandler(db))

	addr := ":8080"
	fmt.Println("api server listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
