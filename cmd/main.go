package main

import (
	"fmt"
	"log"

	"github.com/amirjbr/shared-expense/config"
	"github.com/amirjbr/shared-expense/internal/account/adapter/repository"
	"github.com/amirjbr/shared-expense/internal/account/core/service"
	"github.com/amirjbr/shared-expense/internal/platform/database"
	"github.com/amirjbr/shared-expense/internal/platform/http"
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

	repo, err := repository.NewUserRepo(db)
	if err != nil {
		fmt.Println(err)
	}

	userSvc, err := service.NewUserService(repo, loggger)
	if err != nil {
		fmt.Println(err)
	}

	app := http.NewApp(*conf, loggger, userSvc)
	app.RunAndListen()

}
