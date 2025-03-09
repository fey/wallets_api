package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/fey/wallets_api/docs"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type (
	OperationType string

	WalletOperationRequest struct {
		WalletId      string        `json:"walletId" validate:"required,uuid"`
		OperationType OperationType `json:"operationType" validate:"required,oneof=DEPOSIT WITHDRAW"`
		Amount        float64       `json:"amount" validate:"required,gte=0"`
	}
	Wallet struct {
		WalletId string  `json:"wallet_id"`
		Balance  float64 `json:"balance"`
	}

	ValidationError struct {
		Field string
		Tag   string
		Value string
	}
)

const (
	Deposit  OperationType = "DEPOSIT"
	Withdraw OperationType = "WITHDRAW"
)

var db *sql.DB

func connect() error {
	var err error
	err = godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("error to load config.env: %v", err)
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	// Формируем строку подключения
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(time.Minute * 5)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}
	return nil
}

// @title			Wallet API
// @version		1.0
// @description	API для управления кошельками
// @host			localhost:8080
// @BasePath		/api/v1
func main() {
	if err := connect(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := fiber.New()
	app.Get("/", RootHandler)
	app.Get("/api/v1/wallets/:uuid", GetWalletHandler)
	app.Post("/api/v1/wallets", WalletOperationHandler)

	app.Get("/swagger/*", swagger.HandlerDefault)
	log.Fatal(app.Listen(":8080"))
}

func RootHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "/",
	})
}

func GetWalletHandler(ctx *fiber.Ctx) error {
	uuidParam := ctx.Params("uuid", "")

	if _, err := uuid.Parse(uuidParam); err != nil {
		uuidError := ValidationError{
			Field: "uuid",
			Tag:   "invalid",
			Value: uuidParam,
		}
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": []ValidationError{uuidError},
		})
	}

	var wallet Wallet
	row := db.QueryRow("SELECT id, balance FROM wallets WHERE id = $1", uuidParam)

	err := row.Scan(&wallet.WalletId, &wallet.Balance)

	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	jsonBytes, _ := json.Marshal(&wallet)

	return ctx.SendString(string(jsonBytes))
}

// @Summary		Perform wallet operation
// @Description	Deposit or withdraw an amount from a wallet
// @Tags			wallets
// @Accept			json
// @Produce		json
// @Param			request	body		WalletOperationRequest	true	"Wallet operation request"
// @Success		200		{object}	Wallet					"Updated wallet details"
// @Failure		422		{object}	[]ValidationError		"Validation errors"
// @Failure		404		{object}	string					"Wallet not found"
// @Failure		500		{object}	string					"Internal server error"
// @Router			/api/v1/wallets [post]
func WalletOperationHandler(ctx *fiber.Ctx) error {
	var req WalletOperationRequest

	if err := ctx.BodyParser(&req); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	errs := validateWalletOperation(req)

	if errs != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": errs,
		})
	}

	var wallet Wallet
	row := db.QueryRow("SELECT id, balance FROM wallets WHERE id = $1", req.WalletId)

	err := row.Scan(&wallet.WalletId, &wallet.Balance)

	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if req.OperationType == Deposit {
		_, err := tx.Exec("UPDATE wallets SET balance = balance + $1  WHERE id = $2", req.Amount, req.WalletId)
		if err != nil {
			tx.Rollback()
			return err
		}
		wallet.Balance += req.Amount
	}

	if req.OperationType == Withdraw {
		_, err := tx.Exec("UPDATE wallets SET balance = balance - $1  WHERE id = $2", req.Amount, req.WalletId)
		if err != nil {
			tx.Rollback()
			return err
		}

		wallet.Balance -= req.Amount
	}

	_, err = db.Exec("INSERT INTO transactions(wallet_id, operation_type, amount) VALUES($1, $2, $3)", req.WalletId, req.OperationType, req.Amount)
	if err != nil {
		return err
	}
	tx.Commit()
	return ctx.JSON(wallet)
}

func validateWalletOperation(req WalletOperationRequest) []*ValidationError {
	validate := validator.New(validator.WithRequiredStructEnabled())

	var errors []*ValidationError

	err := validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el ValidationError
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			errors = append(errors, &el)
		}
		return errors
	}

	return nil
}
