package main

import (
	"context"
	"fmt"
	"log"
	"time"
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
	_, err = db.Exec(ctx, `
	DROP TABLE users;
	`)
	if err != nil {
		log.Fatalf("Drop table error: %v", err)
	}
	_, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
        iin VARCHAR (255) NOT NULL,
		balance BIGINT,
        number BIGINT,
		registerDate VARCHAR (255) NOT NULL
	);
	`)
	if err != nil {
		log.Fatalf("Create table error: %v", err)
	}
	// TODO: need this?
	// _, err = db.Exec(ctx,
	// 	`INSERT INTO users(username, password, iin, role) VALUES ($1, $2, $3, $4)`,
	// 	"admin", "pass", "940217200216", "admin")
	// if err != nil {
	// 	log.Fatalf("Add admin error: %v", err)
	// }
	return db
}
