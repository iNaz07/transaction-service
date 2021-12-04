package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	_handler "transaction-service/accounts/delivery/http"
	_repo "transaction-service/accounts/repository/postgres"
	_usecase "transaction-service/accounts/usecase"

	"github.com/labstack/echo/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

//TODO: get this params from environment
const (
	username = "postgres"
	password = "password"
	hostname = "localhost" //change as in docker-compose
	port     = 5432
	dbname   = "transaction" //change as in init.sql
)

func main() {
	db := connectDB()
	defer db.Close()

	accRepo := _repo.NewAccountRepo(db)
	accUsecase := _usecase.NewAccountUsecase(accRepo)
	e := echo.New()
	_handler.NewAccountHandler(e, accUsecase)

	err := e.Start("127.0.0.1:8181")
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(`shutting down the server`, err)
	}
}

func connectDB() *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, hostname, port, dbname)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Open db error: %v", err)
	}
	// should read bout them
	// config.MaxConns = 25
	// config.MaxConnLifetime = 5 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	db, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatalf("Connect db error: %v", err)
	}

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Ping db error: %v", err)
	}
	//TODO: optimize this, should save old data
	if _, err = db.Exec(ctx, `DROP TABLE accounts;`); err != nil {
		log.Fatalf("Drop table error: %v", err)
	}
	if _, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
        iin VARCHAR (255) NOT NULL,
		balance BIGINT,
        number VARCHAR (255) NOT NULL,
		registerDate VARCHAR (255) NOT NULL
	);`); err != nil {
		log.Fatalf("Create accounts table error: %v", err)
	}

	if _, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
        sender VARCHAR (255) NOT NULL,
		sender_number VARCHAR (255) NOT NULL,
        recipient_number VARCHAR (255) NOT NULL,
		recipient VARCHAR (255) NOT NULL,
		amount BIGINT,
		date VARCHAR (255) NOT NULL
	);
	`); err != nil {
		log.Fatalf("Create transactions table error: %v", err)
	}

	if _, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS deposits (
		id SERIAL PRIMARY KEY,
        iin VARCHAR (255) NOT NULL,
        number VARCHAR (255) NOT NULL,
		amount BIGINT,
		date VARCHAR (255) NOT NULL
	);
	`); err != nil {
		log.Fatalf("Create deposits table error: %v", err)
	}
	return db
}
