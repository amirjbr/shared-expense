package main

import (
	"fmt"
	"github.com/amirjbr/shared-expense/config"
	"github.com/joho/godotenv"
	"log"
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
	fmt.Println(conf)
}
