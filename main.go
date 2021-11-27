package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_userDeliveryHttp "github.com/bazeeko/investor-social-network/user/delivery/http"
	_userRepoMysql "github.com/bazeeko/investor-social-network/user/repository/mysql"
	_userUsecase "github.com/bazeeko/investor-social-network/user/usecase"

	_stockDeliveryHttp "github.com/bazeeko/investor-social-network/stock/delivery/http"
	_stockRepoMysql "github.com/bazeeko/investor-social-network/stock/repository/mysql"
	_stockUsecase "github.com/bazeeko/investor-social-network/stock/usecase"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func connectDB(config string) (*sql.DB, error) {
	conn, err := sql.Open("mysql", config)
	if err != nil {
		return nil, fmt.Errorf("connectDB: %w", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("connectDB: %w", err)
	}

	// _, err = conn.Exec(`CREATE DATABASE IF NOT EXISTS investordb`)
	// if err != nil {
	// 	return nil, fmt.Errorf("connectDB: %w", err)
	// }

	_, err = conn.Exec(`USE heroku_c64bdd5da1fe53c`)
	if err != nil {
		return nil, fmt.Errorf("connectDB: %w", err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS users (
		id BIGINT NOT NULL UNIQUE AUTO_INCREMENT,
		username VARCHAR(40) NOT NULL UNIQUE,
		password VARCHAR(40) NOT NULL,
		rating BIGINT NOT NULL,
		created_at INT NOT NULL,
		PRIMARY KEY (id)
	);`)

	if err != nil {
		return nil, fmt.Errorf("connectDB: %w", err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS threads (
		id BIGINT NOT NULL AUTO_INCREMENT UNIQUE,
		user_id BIGINT NOT NULL,
		topic TEXT NOT NULL,
		body TEXT,
		image_url TEXT,
		created_at INT NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`)

	if err != nil {
		return nil, fmt.Errorf("connectDB: %w", err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS category (
		id BIGINT NOT NULL AUTO_INCREMENT UNIQUE,
		user_id BIGINT NOT NULL,
		topic TEXT NOT NULL,
		body TEXT,
		image_url TEXT,
		created_at INT NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`)

	if err != nil {
		return nil, fmt.Errorf("connectDB: %w", err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS stocks (
		symbol VARCHAR(10) NOT NULL UNIQUE,
		name TEXT,
		info TEXT,
		image_url TEXT,
		PRIMARY KEY (symbol)
	);`)

	if err != nil {
		return nil, fmt.Errorf("connectDB: %w", err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS comments (
		id BIGINT NOT NULL AUTO_INCREMENT UNIQUE,
		user_id BIGINT NOT NULL,
		thread_id BIGINT NOT NULL,
		body TEXT,
		image_url TEXT,
		created_at INT NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (thread_id) REFERENCES threads(id)
	);`)

	if err != nil {
		return nil, fmt.Errorf("connectDB: %w", err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS favourite_stocks (
		user_id BIGINT NOT NULL,
		stock_symbol VARCHAR(10) NOT NULL,
		PRIMARY KEY (user_id, stock_symbol)
	);`)

	if err != nil {
		log.Println(err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS favourite_people (
		user_id BIGINT NOT NULL,
		favourite_user_id BIGINT NOT NULL,
		PRIMARY KEY (user_id, favourite_user_id)
	);`)

	if err != nil {
		log.Println(err)
	}

	_, err = conn.Exec(`INSERT stocks (symbol, name, info, image_url) VALUES (?, ?, ?, ?)`,
		"AAPL", "Apple Inc.", "Apple Bla bla bla bla bla", "https://example.com/")

	if err != nil {
		log.Println(err)
	}

	_, err = conn.Exec(`INSERT users (username, password, rating, created_at) VALUES (?, ?, ?, ?)`,
		"DWAdmin", "12345", 0, 0)

	if err != nil {
		log.Println(err)
	}

	for i := 0; i < 10; i++ {
		_, err = conn.Exec(`INSERT users (username, password, rating, created_at) VALUES (?, ?, ?, ?)`,
			fmt.Sprintf("user_%d", i), "12345", 0, 0)

		if err != nil {
			log.Println(err)
		}
	}

	return conn, nil
}

func main() {
	// tcp(127.0.0.1:3306)
	// dbURL := os.Getenv("CLEARDB_DATABASE_URL")

	// mysql://b37bfcbb24c371:17b6ee02@us-cdbr-east-04.cleardb.com/heroku_c64bdd5da1fe53c?reconnect=true

	config := "b37bfcbb24c371:17b6ee02@tcp(us-cdbr-east-04.cleardb.com)/heroku_c64bdd5da1fe53c"

	// config := "root:password@tcp(tcp(127.0.0.1:3306))/"
	db, err := connectDB(config)
	if err != nil {
		log.Fatalf("main: %s\n", err)
	}

	e := echo.New()

	userRepoMysql := _userRepoMysql.NewMysqlUserRepository(db)
	userUsecase := _userUsecase.NewUserUsecase(userRepoMysql)

	stockRepoMysql := _stockRepoMysql.NewMysqlStockRepository(db)
	stockUsecase := _stockUsecase.NewStockUsecase(stockRepoMysql)

	_stockDeliveryHttp.NewStockHandler(e, stockUsecase, userUsecase)
	_userDeliveryHttp.NewUserHandler(e, userUsecase, stockUsecase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatalln(e.Start(":" + port))
	// conn, err := sql.Open("mysql", config)
}
