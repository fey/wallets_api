package main_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB

func setup() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	// Выполняем SQL-запросы из файла init.sql
	initSQL, err := os.ReadFile("init.sql")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(string(initSQL))
	if err != nil {
		panic(err)
	}
}

func teardown() {
	// Здесь вы можете очистить базу данных после тестов
	db.Close()
}

func TestRootHandler(t *testing.T) {
	setup()
	defer teardown()

	app := fiber.New()
	app.Get("/", RootHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "/", response["message"])
}

func TestGetWalletHandler(t *testing.T) {
	setup()
	defer teardown()

	app := fiber.New()
	app.Get("/api/v1/wallets/:uuid", GetWalletHandler)

	// Предположим, что у вас есть кошелек с UUID "test-uuid"
	req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/550e8400-e29b-41d4-a716-446655440000", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var wallet Wallet
	json.NewDecoder(resp.Body).Decode(&wallet)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", wallet.WalletId)
}

func TestWalletOperationHandler(t *testing.T) {
	setup()
	defer teardown()

	app := fiber.New()
	app.Post("/api/v1/wallets", WalletOperationHandler)

	// Пример запроса на депозит
	reqBody := WalletOperationRequest{
		WalletId:      "550e8400-e29b-41d4-a716-446655440000",
		OperationType: Deposit,
		Amount:        100.0,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var wallet Wallet
	json.NewDecoder(resp.Body).Decode(&wallet)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", wallet.WalletId)
	assert.Equal(t, 200.0, wallet.Balance) // Предполагаем, что начальный баланс был 0
}

func TestWalletOperationHandlerValidationError(t *testing.T) {
	setup()
	defer teardown()

	app := fiber.New()
	app.Post("/api/v1/wallets", WalletOperationHandler)

	// Пример запроса с ошибкой валидации (недостаточно данных)
	reqBody := WalletOperationRequest{
		WalletId:      "", // Неверный UUID
		OperationType: Deposit,
		Amount:        -50.0, // Неверная сумма
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	var validationErrors map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&validationErrors)
	assert.NotEmpty(t, validationErrors["errors"])
}
