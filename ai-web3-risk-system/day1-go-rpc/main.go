package main

import (
    "context"
    "fmt"
    "log"
	"os"
	"time"

    "github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
    	log.Fatal("missing SEPOLIA_RPC_URL")
	}
	client, err := ethclient.Dial(rpcURL)
	if err != nil{
		log.Fatal("failed to connect rpc:",err)
	}
	ctx, cancel := context.WithTimeout(context.Background(),15 * time.Second)
	defer cancel()
	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		log.Fatal("failed to get latest block number:",err)

	}
	fmt.Println("Latest Block:",blockNumber)
}