package main

import (
	"fmt"

	"github.com/amirjbr/shared-expense/internal/ledger/app/dto"
)

func main() {
	userInput := 1000000
	participantLength := 3

	baghimande := userInput % participantLength

	var createResponseReq dto.CreateExpenseRequest
	createResponseReq.TotalAmount = userInput

	if baghimande != 0 {
		for i := 1; i <= int(baghimande); i++ {
			for j := 1; j <= i; j++ {
				createExpenseReq.Shares[j].Amount += 1
			}
		}
		// Calculate each share
	}
}
