package handler

import (
	"net/http"

	"github.com/amirjbr/shared-expense/internal/ledger/app/dto"
	"github.com/amirjbr/shared-expense/internal/ledger/core/entity"
	"github.com/amirjbr/shared-expense/internal/ledger/core/service"
	"github.com/amirjbr/shared-expense/internal/ledger/core/service/command"
	"github.com/amirjbr/shared-expense/internal/platform/http/httpcontext"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LedgerHandler struct {
	service *service.LedgerService
}

func NewLedgerHandler(service *service.LedgerService) *LedgerHandler {
	return &LedgerHandler{
		service: service,
	}
}

func (r *LedgerHandler) CreateExpenseHandler(c *gin.Context) {
	var req dto.CreateExpenseRequest

	// -----------------------
	// Parse request body
	// -----------------------

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	// -----------------------
	// Parse group ID from URL
	// -----------------------

	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid group ID",
		})
		return
	}

	// -----------------------
	// Get authenticated user
	// -----------------------

	creatorID, ok := httpcontext.GetAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "authenticated user information is unavailable",
		})
		return
	}

	// -----------------------
	// Parse payers
	// -----------------------

	payers := make([]command.CreateExpensePayer, 0, len(req.Payers))

	for _, payer := range req.Payers {
		userID, err := uuid.Parse(payer.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid payer user ID",
			})
			return
		}

		payers = append(payers, command.CreateExpensePayer{
			UserID: userID,
			Amount: entity.Money(payer.Amount),
		})
	}

	// -----------------------
	// Parse participants
	// -----------------------

	participants := make([]uuid.UUID, 0, len(req.Participants))

	for _, participant := range req.Participants {
		userID, err := uuid.Parse(participant)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid participant user ID",
			})
			return
		}
		participants = append(participants, userID)
	}

	// -----------------------
	// Parse custom shares
	// -----------------------

	shares := make([]command.CreateExpenseShare, 0, len(req.Shares))

	for _, share := range req.Shares {
		userID, err := uuid.Parse(share.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid share user ID",
			})
			return
		}

		shares = append(shares, command.CreateExpenseShare{
			UserID: userID,
			Amount: entity.Money(share.Amount),
		})
	}

	// -----------------------
	// Build application command
	// -----------------------

	cmd := command.CreateExpenseCommand{
		GroupID:         groupID,
		CreatedByUserID: creatorID,
		Title:           req.Title,
		TotalAmount:     entity.Money(req.TotalAmount),
		SplitType:       entity.SplitType(req.SplitType),
		OccurredAt:      req.OccurredAt,
		Payers:          payers,
		Participants:    participants,
		Shares:          shares,
	}

	// -----------------------
	// Call application service
	// -----------------------

	expenseID, err := r.service.CreateExpense(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create expense",
		})
		return
	}

	// -----------------------
	// Response
	// -----------------------

	c.JSON(http.StatusCreated, gin.H{
		"id": expenseID.String(),
	})
}
