package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
	"strings"
)

func RegisterHandlers(mux *http.ServeMux, srv *Service) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	mux.HandleFunc("POST /api/v1/wallet", srv.Wallet)
	mux.HandleFunc("GET /api/v1/wallets/{wallet_id}", srv.GetBalance)
}

func (srv Service) Wallet(w http.ResponseWriter, r *http.Request) {
	var body WalletBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	op := strings.ToUpper(body.OperationType)
	if op != "DEPOSIT" && op != "WITHDRAW" {
		http.Error(w, "operation_type must be DEPOSIT or WITHDRAW", http.StatusBadRequest)
		return
	}
	if body.WalletId == "" {
		http.Error(w, "wallet_id is required", http.StatusBadRequest)
		return
	}
	if body.Amount.LessThanOrEqual(decimal.Zero) {
		http.Error(w, "amount must be > 0", http.StatusBadRequest)
		return
	}

	newAmount, err := srv.updateBalance(r.Context(), body.WalletId, op, body.Amount)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			http.Error(w, "wallet not found", http.StatusNotFound)
		case errors.Is(err, ErrInsufficient):
			http.Error(w, "insufficient funds", http.StatusConflict) // 409, не 5xx
		default:
			log.Println("wallet update error:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"amount": newAmount.String()})
}

func (srv Service) GetBalance(w http.ResponseWriter, r *http.Request) {
	walletId := r.PathValue("wallet_id")
	if walletId == "" {
		http.Error(w, "wallet_id is required", http.StatusBadRequest)
		return
	}

	const sql = `SELECT amount::text FROM wallet WHERE uuid = $1`
	var amountStr string
	if err := srv.pool.QueryRow(r.Context(), sql, walletId).Scan(&amountStr); err != nil {
		http.Error(w, "wallet not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"amount": amountStr})
}

func (srv Service) updateBalance(ctx context.Context, walletID, op string, amt decimal.Decimal) (decimal.Decimal, error) {
	const sql = `
        UPDATE wallet
        SET amount = CASE
                       WHEN $2 = 'DEPOSIT'  THEN amount + $3
                       WHEN $2 = 'WITHDRAW' THEN amount - $3
                     END
        WHERE uuid = $1
          AND ($2 = 'DEPOSIT' OR ($2 = 'WITHDRAW' AND amount >= $3))
        RETURNING amount::text;
    `
	var newAmountStr string
	if err := srv.pool.QueryRow(ctx, sql, walletID, op, amt.String()).Scan(&newAmountStr); err != nil {
		var exists bool
		_ = srv.pool.QueryRow(ctx, `SELECT true FROM wallet WHERE uuid=$1`, walletID).Scan(&exists)
		if !exists {
			return decimal.Zero, ErrNotFound
		}
		return decimal.Zero, ErrInsufficient
	}
	return decimal.NewFromString(newAmountStr)
}
