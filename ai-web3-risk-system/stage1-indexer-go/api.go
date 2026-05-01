package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"strconv"
)

// TransfersHandler 提供当前最小可用 API：
// GET /transfers?address=0x...
func TransfersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "missing address", http.StatusBadRequest)
			return
		}

		// contract 是可选过滤条件；空字符串表示查询该地址相关的所有 token 转账。
		contractAddress := r.URL.Query().Get("contract")

		limit := 20

		limitStr := r.URL.Query().Get("limit")
		if limitStr != "" {
			parsedLimit, err := strconv.Atoi(limitStr)
			if err != nil {
				http.Error(w, "invalid limit", http.StatusBadRequest)
				return
			}

			if parsedLimit <= 0 || parsedLimit > 100 {
				http.Error(w, "limit must be between 1 and 100", http.StatusBadRequest)
				return
			}

			limit = parsedLimit
		}

		events, err := QueryTransferEventsByAddress(r.Context(), db, address, contractAddress, limit)
		if err != nil {
			http.Error(w, "failed to query transfer events", http.StatusInternalServerError)
			return
		}

		resp := make([]TransferEventResponse, 0, len(events))
		for _, event := range events {
			valueStr := ""
			if event.Value != nil {
				valueStr = event.Value.String()
			}

			resp = append(resp, TransferEventResponse{
				TransactionHash: event.TransactionHash,
				LogIndex:        event.LogIndex,
				BlockNumber:     event.BlockNumber,
				ContractAddress: event.ContractAddress,
				FromAddress:     event.FromAddress,
				ToAddress:       event.ToAddress,
				Value:           valueStr,
			})
		}

		// 先保持返回结构简单：直接返回转账记录数组。
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
