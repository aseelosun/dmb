package db

import (
	"database/sql"
	"doMassageBot/internal/config"
	"fmt"
	"log"
)

func ConnectingToDb(conf *config.Configuration) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s", conf.DbConfig.Host, conf.DbConfig.Port, conf.DbConfig.User, conf.DbConfig.Password, conf.DbConfig.DbName, conf.DbConfig.SearchPath)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	db.SetMaxIdleConns(conf.DbConfig.IdleConns)
	db.SetMaxOpenConns(conf.DbConfig.OpenConns)
	db.Exec(`set search_path='doMassageBot'`)
	fmt.Printf("Postgres Connected!\n")
	return db, err
}
