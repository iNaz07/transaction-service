package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

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
		log.Fatal().Err(err).Msg("get configuration error")
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
		log.Fatal().Err(err).Msg("cannot start the server")
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
		log.Fatal().Err(err).Msg("parse dsn error")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	db, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatal().Err(err).Msg("connect db error")
	}

	if err := db.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("db ping error")
	}

	//Move all to init.sql
	//TODO: optimize this, should save old data
	// if _, err = db.Exec(ctx, `DROP TABLE accounts;`); err != nil {
	// 	log.Fatal().Err(err).Msg("Drop table error")
	// }

	if _, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY UNIQUE,
		userid BIGINT NOT NULL,
        iin VARCHAR (255) NOT NULL,
		balance BIGINT NOT NULL,
        number VARCHAR (255) NOT NULL UNIQUE,
		registerDate TEXT NOT NULL,
		lasttransaction TEXT NOT NULL
	);`); err != nil {
		log.Fatal().Err(err).Msg("Error create table: accounts")
	}

	if _, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY UNIQUE,
        sender VARCHAR (255) NOT NULL,
		sender_number VARCHAR (255) NOT NULL,
        recipient_number VARCHAR (255) NOT NULL,
		recipient VARCHAR (255) NOT NULL,
		amount BIGINT NOT NULL,
		date VARCHAR (255) NOT NULL
	);
	`); err != nil {
		log.Fatal().Err(err).Msg("Error create table: transactions")
	}

	if _, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS deposits (
		id SERIAL PRIMARY KEY UNIQUE,
        iin VARCHAR (255) NOT NULL,
        number VARCHAR (255) NOT NULL,
		amount BIGINT NOT NULL,
		date VARCHAR (255) NOT NULL
	);
	`); err != nil {
		log.Fatal().Err(err).Msg("Error create table: deposits")
	}
	return db
}
