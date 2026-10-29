package main

import (
	"fmt"
	"log"

	"github.com/amirjbr/shared-expense/config"
	UserRepo "github.com/amirjbr/shared-expense/internal/account/adapter/repository"
	UserHandlers "github.com/amirjbr/shared-expense/internal/account/app/handler"
	UserRoutes "github.com/amirjbr/shared-expense/internal/account/app/routes"
	UserService "github.com/amirjbr/shared-expense/internal/account/core/service"
	GroupRepo "github.com/amirjbr/shared-expense/internal/groups/adapter/repository"
	GroupHandlers "github.com/amirjbr/shared-expense/internal/groups/app/handler"
	GroupRoutes "github.com/amirjbr/shared-expense/internal/groups/app/routes"
	GroupService "github.com/amirjbr/shared-expense/internal/groups/core/service"
	"github.com/amirjbr/shared-expense/internal/platform/database"
	"github.com/amirjbr/shared-expense/internal/platform/http"
	"github.com/amirjbr/shared-expense/pkg/conv"
	"github.com/amirjbr/shared-expense/pkg/logger"
	"github.com/amirjbr/shared-expense/pkg/migrator"
	"github.com/joho/godotenv"
)

// TODO separate technical erros with bussines errors
func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	logg := logger.New()

	db, err := database.InitDB(*conf)
	if err != nil {
		fmt.Println(err)
	}

	mig := migrator.NewMigrator(db, "postgres")
	mig.Up()

	// user module
	userRepo, err := UserRepo.NewUserRepo(db)
	if err != nil {
		fmt.Println(err)
	}
	userSvc, err := UserService.NewUserService(userRepo, conv.ToBytes(conf.JwtSecret), logg, conf.Auth.TokenExpiresMinute, conf.Auth.TokenRefreshMinute)
	if err != nil {
		fmt.Println(err)
	}
	userHandler := UserHandlers.NewUserHandler(userSvc)
	userRoutes := UserRoutes.NewUserRoutes(userHandler, logg)
	userLookUpSvc, err := UserService.NewUserLookUp(userRepo, logg)
	if err != nil {
		fmt.Println(err)
	}

	//group module

	groupRep, err := GroupRepo.NewGroupRepo(db)
	if err != nil {
		fmt.Println(err)
	}
	groupSvc, err := GroupService.NewGroupService(groupRep, userLookUpSvc, logg)
	groupHandler := GroupHandlers.NewGroupHandler(groupSvc)
	groupRoute := GroupRoutes.NewGroupRoutes(groupHandler, logg)

	server := http.NewServer("localhost:8080")
	server.Register("api/v1/share_expense", userRoutes, groupRoute)

	if err = server.Run(); err != nil {
		log.Fatal(err)
	}

}
