package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "sync"
  "github.com/swaggo/http-swagger"
  _ "github.com/fey/wallets_api/docs"
)

type Wallet struct {
  ID      string  `json:"walletId"`
  Balance float64 `json:"balance"`
}

var (
  wallets = make(map[string]*Wallet)
  mu sync.Mutex
)
// @title Wallet API
// @version 1.0
// @description API для управления кошельками
// @host localhost:8080
// @BasePath /api/v1
func main() {
  http.HandleFunc("/api/v1/wallet", handleWallet)
  http.HandleFunc("/api/v1/wallets/", handleGetWallet)

  // Swagger UI
  http.Handle("/swagger/", httpSwagger.WrapHandler)

  fmt.Println("Starting server on :8080...")
  if err := http.ListenAndServe(":8080", nil); err != nil {
      log.Fatal(err)
  }
}

// @Summary Создание или обновление кошелька
// @Description Создает новый кошелек или обновляет существующий
// @Accept json
// @Produce json
// @Param wallet body struct {
//     WalletID     string  `json:"walletId"`
//     OperationType string  `json:"operationType"`
//     Amount       float64 `json:"amount"`
// } true "Wallet"
// @Success 200 {object} Wallet
// @Failure 400 {string} string "Invalid request"
// @Router /wallet [post]
func handleWallet(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodPost {
      var req struct {
          WalletID     string  `json:"walletId"`
          OperationType string  `json:"operationType"`
          Amount       float64 `json:"amount"`
      }

      if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
          http.Error(w, "Invalid request", http.StatusBadRequest)
          return
      }

      mu.Lock()
      defer mu.Unlock()

      wallet, exists := wallets[req.WalletID]
      if !exists {
          wallet = &Wallet{ID: req.WalletID, Balance: 0}
          wallets[req.WalletID] = wallet
      }

      switch req.OperationType {
      case "DEPOSIT":
          wallet.Balance += req.Amount
      case "WITHDRAW":
          if wallet.Balance < req.Amount {
              http.Error(w, "Insufficient funds", http.StatusBadRequest)
              return
          }
          wallet.Balance -= req.Amount
      default:
          http.Error(w, "Invalid operation type", http.StatusBadRequest)
          return
      }

      w.WriteHeader(http.StatusOK)
      json.NewEncoder(w).Encode(wallet)
  } else {
      http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
  }
}

// @Summary Получение информации о кошельке
// @Description Получает информацию о кошельке по его ID
// @Produce json
// @Param walletId path string true "Wallet ID"
// @Success 200 {object} Wallet
// @Failure 404 {string} string "Wallet not found"
// @Router /wallets/{walletId} [get]
func handleGetWallet(w http.ResponseWriter, r *http.Request) {
  walletID := r.URL.Path[len("/api/v1/wallets/"):]

  mu.Lock()
  defer mu.Unlock()

  wallet, exists := wallets[walletID]
  if !exists {
      http.Error(w, "Wallet not found", http.StatusNotFound)
      return
  }

  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(wallet)
}
