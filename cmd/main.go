package main

import (
	"fmt"
	"log"

	"github.com/amirjbr/shared-expense/config"
	"github.com/amirjbr/shared-expense/internal/account/adapter/repository"
	"github.com/amirjbr/shared-expense/internal/account/app/handler"
	"github.com/amirjbr/shared-expense/internal/account/app/routes"
	"github.com/amirjbr/shared-expense/internal/account/core/service"
	repository2 "github.com/amirjbr/shared-expense/internal/groups/adapter/repository"
	handler2 "github.com/amirjbr/shared-expense/internal/groups/app/handler"
	routes2 "github.com/amirjbr/shared-expense/internal/groups/app/routes"
	service2 "github.com/amirjbr/shared-expense/internal/groups/core/service"
	"github.com/amirjbr/shared-expense/internal/platform/database"
	"github.com/amirjbr/shared-expense/internal/platform/http"
	"github.com/amirjbr/shared-expense/pkg/conv"
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

	db, err := database.InitDB(*conf)
	if err != nil {
		fmt.Println(err)
	}

	mig := migrator.NewMigrator(db, "postgres")
	mig.Up()

	userRepo, err := repository.NewUserRepo(db)
	if err != nil {
		fmt.Println(err)
	}
	userSvc, err := service.NewUserService(userRepo, conv.ToBytes(conf.JwtSecret), loggger, conf.Auth.TokenExpiresMinute, conf.Auth.TokenRefreshMinute)
	if err != nil {
		fmt.Println(err)
	}
	userHandler := handler.NewUserHandler(userSvc)
	userRoutes := routes.NewUserRoutes(userHandler)

	groupRepo, err := repository2.NewGroupRepo(db)
	if err != nil {
		fmt.Println(err)
	}
	groupSvc, err := service2.NewGroupService(groupRepo, loggger)
	groupHandler := handler2.NewGroupHandler(groupSvc)
	groupRoutes := routes2.NewGroupRoutes(groupHandler)

	server := http.NewServer("localhost:8080")
	server.Register("api/v1/share_expense", userRoutes, groupRoutes)

	if err = server.Run(); err != nil {
		log.Fatal(err)
	}

}
