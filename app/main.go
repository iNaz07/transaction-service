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
	"transaction-service/domain"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("get configuration error: ", err)
	}
}

func main() {

	token := &domain.JwtToken{
		AccessSecret: viper.GetString(`token.secret`),
		AccessTtl:    viper.GetDuration(`token.ttl`) * time.Minute,
	}
	timeout := viper.GetDuration(`timeout`) * time.Second
	db := connectDB()
	defer db.Close()

	accRepo := _repo.NewAccountRepo(db)
	accUsecase := _usecase.NewAccountUsecase(accRepo, timeout)
	jwtUsecase := _usecase.NewJWTUseCase(token)

	e := echo.New()
	_handler.NewAccountHandler(e, accUsecase, jwtUsecase)

	err := e.Start(viper.GetString(`addr`))
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(`shutting down the server`, err)
	}
}

func connectDB() *pgxpool.Pool {

	username := viper.GetString(`postgres.user`)
	password := viper.GetString(`postgres.password`)
	hostname := viper.GetString(`postgres.host`)
	port := viper.GetInt(`postgres.port`)
	dbname := viper.GetString(`postgres.dbname`)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, hostname, port, dbname)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Open db error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	db, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatalf("Connect db error: %v", err)
	}

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Ping db error: %v", err)
	}

	//Move all to init.sql
	//TODO: optimize this, should save old data
	// if _, err = db.Exec(ctx, `DROP TABLE accounts;`); err != nil {
	// 	log.Fatalf("Drop table error: %v", err)
	// }

	if _, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY UNIQUE,
		userid BIGINT NOT NULL,
        iin VARCHAR (255) NOT NULL,
		balance BIGINT NOT NULL,
        number VARCHAR (255) NOT NULL UNIQUE,
		registerDate TEXT NOT NULL,
		lasttransaction TEXT
	);`); err != nil {
		log.Fatalf("Create accounts table error: %v", err)
	}

	if _, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY UNIQUE,
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
		id SERIAL PRIMARY KEY UNIQUE,
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
