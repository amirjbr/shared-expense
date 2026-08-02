package main

import (
	"fmt"
	"log"

	"github.com/amirjbr/shared-expense/config"
	"github.com/amirjbr/shared-expense/internal/platform/database"
	"github.com/amirjbr/shared-expense/pkg/logger"
	"github.com/amirjbr/shared-expense/pkg/migrator"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	loggger := logger.NewLogger()

	fmt.Println(loggger)

	db, err := database.InitDB(*conf)
	if err != nil {
		fmt.Println(err)
	}

	mig := migrator.NewMigrator(db, "postgres")
	mig.Up()

}
