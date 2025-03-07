package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/fey/wallets_api/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/lib/pq"
)

type OperationType string

const (
	Deposit  OperationType = "DEPOSIT"
	Withdraw OperationType = "WITHDRAW"
)

type WalletOperationRequest struct {
	WalletId    string        `json:"walletId"`
	OperationType OperationType `json:"operationType"`
	Amount        float64       `json:"amount"`
}

type Wallet struct {
	WalletId      string  `json:"walletId"`
	Balance float64 `json:"balance"`
}

// @title			Wallet API
// @version		1.0
// @description	API для управления кошельками
// @host			localhost:8080
// @BasePath		/api/v1
func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to db!")

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		var greeting string
		err := db.QueryRow("SELECT 'Hello, World!'").Scan(&greeting)
		if err != nil {
			return err
		}
		return c.SendString(greeting)
	})
	app.Get("/api/v1/wallets/:uuid", GetWalletHandler)
	app.Post("api/v1/wallets", PostWalletHandler)
	app.Get("/swagger/*", swagger.HandlerDefault)
	log.Fatal(app.Listen(":8080"))
}

func GetWalletHandler(c *fiber.Ctx) error {
	uuid := c.Params("uuid", "")

	wallet := Wallet{
		WalletId: uuid,
		Balance: 500.0,
	}

	jsonBytes, _ := json.Marshal(&wallet)

	return c.SendString(string(jsonBytes))

	// if uuid == "" {
	// 	return c.SendStatus(fiber.StatusNotFound)
	// }

	// return c.SendString(uuid)
}

func PostWalletHandler(c *fiber.Ctx) error {
	var request WalletOperationRequest

	if err := c.BodyParser(&request); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	fmt.Println(request.WalletId)

	wallet := Wallet{
		WalletId: request.WalletId,
		Balance: 500,
	}

	c.Status(fiber.StatusCreated)

	return c.JSON(wallet)
}
