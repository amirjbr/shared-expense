package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/amirjbr/shared-expense/config"
	_ "github.com/lib/pq"
)

func InitDB(cfg config.Config) (*sql.DB, error) {
	opts := ConvertConfigToDBOptions(cfg)
	dbInstance, err := sql.Open("postgres", opts.PostgresDSN())
	if err != nil {
		log.Fatalf("err on opening database %v", err)
	}
	err = dbInstance.Ping()
	if err != nil {
		log.Fatalf("error with pinging to database %v", err)
	}
	fmt.Println("Successfully connected to database")
	return dbInstance, nil
}

type DBConnectionOptions struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func (o DBConnectionOptions) PostgresDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		o.Host, o.Port, o.Username, o.Password, o.Database)
}

func ConvertConfigToDBOptions(cfg config.Config) *DBConnectionOptions {
	return &DBConnectionOptions{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
		Database: cfg.DB.Database,
	}
}
