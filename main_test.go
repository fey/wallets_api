package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleWallet(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    WalletRequest
		expectedStatus int
		expectedBalance float64
	}{
		{
			name: "Deposit money",
			requestBody: WalletRequest{
				WalletID:      "wallet1",
				OperationType: Deposit,
				Amount:        100.0,
			},
			expectedStatus: http.StatusOK,
			expectedBalance: 100.0,
		},
		{
			name: "Withdraw money",
			requestBody: WalletRequest{
				WalletID:      "wallet1",
				OperationType: Withdraw,
				Amount:        50.0,
			},
			expectedStatus: http.StatusOK,
			expectedBalance: 50.0,
		},
		{
			name: "Withdraw more than balance",
			requestBody: WalletRequest{
				WalletID:      "wallet1",
				OperationType: Withdraw,
				Amount:        100.0,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid operation type",
			requestBody: WalletRequest{
				WalletID:      "wallet1",
				OperationType: "INVALID",
				Amount:        50.0,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, err := http.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(HandleWallet)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, status)
			}

			if tt.expectedStatus == http.StatusOK {
				var wallet Wallet
				if err := json.NewDecoder(rr.Body).Decode(&wallet); err != nil {
					t.Fatal(err)
				}
				if wallet.Balance != tt.expectedBalance {
					t.Errorf("expected balance %v, got %v", tt.expectedBalance, wallet.Balance)
				}
			}
		})
	}
}

func TestHandleGetWallet(t *testing.T) {
	// Создаем кошелек для тестирования
	wallets["wallet1"] = &Wallet{ID: "wallet1", Balance: 100.0}

	tests := []struct {
		name           string
		walletID      string
		expectedStatus int
		expectedBalance float64
	}{
		{
			name: "Get existing wallet",
			walletID: "wallet1",
			expectedStatus: http.StatusOK,
			expectedBalance: 100.0,
		},
		{
			name: "Get non-existing wallet",
			walletID: "wallet2",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/v1/wallets/"+tt.walletID, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handleGetWallet)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, status)
			}

			if tt.expectedStatus == http.StatusOK {
				var wallet Wallet
				if err := json.NewDecoder(rr.Body).Decode(&wallet); err != nil {
					t.Fatal(err)
				}
				if wallet.Balance != tt.expectedBalance {
					t.Errorf("expected balance %v, got %v", tt.expectedBalance, wallet.Balance)
				}
			}
		})
	}
}
