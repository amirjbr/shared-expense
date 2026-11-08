package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/amirjbr/shared-expense/internal/ledger/core/entity"
	"github.com/amirjbr/shared-expense/internal/ledger/core/port"
	"github.com/amirjbr/shared-expense/internal/ledger/core/service/command"
	"github.com/amirjbr/shared-expense/internal/ledger/core/utils"
	"github.com/google/uuid"
)

type LedgerService struct {
	repo        port.LedgerRepo
	logger      *slog.Logger
	groupFinder port.GroupFinder
}

func NewLedgerService(repo port.LedgerRepo, groupFinder port.GroupFinder, logger *slog.Logger) (*LedgerService, error) {
	if repo == nil {
		return nil, errors.New("repo is nil")
	}
	return &LedgerService{repo: repo,
		groupFinder: groupFinder, logger: logger}, nil
}

func (s *LedgerService) CreateExpense(ctx context.Context, cmd command.CreateExpenseCommand) (uuid.UUID, error) {

	// -----------------------
	// Basic validation
	// -----------------------

	if cmd.GroupID == uuid.Nil {
		return uuid.Nil, errors.New("group ID is empty")
	}

	if cmd.CreatedByUserID == uuid.Nil {
		return uuid.Nil, errors.New("creator user ID is empty")
	}

	if strings.TrimSpace(cmd.Title) == "" {
		return uuid.Nil, errors.New("expense title is empty")
	}

	if cmd.TotalAmount <= 0 {
		return uuid.Nil, errors.New(
			"total amount must be greater than zero",
		)
	}

	if len(cmd.Payers) == 0 {
		return uuid.Nil, errors.New(
			"expense must have at least one payer",
		)
	}

	if cmd.SplitType != entity.EqualSplitType &&
		cmd.SplitType != entity.CustomSplitType {

		return uuid.Nil, errors.New("invalid split type")
	}

	// -----------------------
	// Validate payers
	// -----------------------

	payers := make(
		[]entity.ExpensePayer,
		0,
		len(cmd.Payers),
	)

	payerIDs := make(
		map[uuid.UUID]struct{},
		len(cmd.Payers),
	)

	var totalPaid entity.Money

	for _, payer := range cmd.Payers {

		if payer.UserID == uuid.Nil {
			return uuid.Nil, errors.New(
				"payer user ID is empty",
			)
		}

		if payer.Amount <= 0 {
			return uuid.Nil, errors.New(
				"payer amount must be greater than zero",
			)
		}

		if _, exists := payerIDs[payer.UserID]; exists {
			return uuid.Nil, errors.New("duplicate payer")
		}

		payerIDs[payer.UserID] = struct{}{}

		totalPaid += payer.Amount

		payers = append(
			payers,
			entity.ExpensePayer{
				UserID: payer.UserID,
				Amount: payer.Amount,
			},
		)
	}

	if totalPaid != cmd.TotalAmount {
		return uuid.Nil, errors.New(
			"sum of payer amounts must equal total amount",
		)
	}

	// -----------------------
	// Users requiring membership validation
	// -----------------------

	memberIDs := make(
		[]uuid.UUID,
		0,
		1+len(cmd.Payers)+len(cmd.Participants)+len(cmd.Shares),
	)

	// Creator must be a group member.
	memberIDs = append(
		memberIDs,
		cmd.CreatedByUserID,
	)

	for _, payer := range cmd.Payers {
		memberIDs = append(
			memberIDs,
			payer.UserID,
		)
	}

	// -----------------------
	// Build expense shares
	// -----------------------

	var shares []entity.ExpenseShare

	switch cmd.SplitType {

	case entity.EqualSplitType:

		if len(cmd.Participants) == 0 {
			return uuid.Nil, errors.New(
				"equal split requires at least one participant",
			)
		}

		if len(cmd.Shares) != 0 {
			return uuid.Nil, errors.New(
				"shares should not be provided for equal split",
			)
		}

		/*
			Because Money is whole Toman and every share
			must be positive:

			    total = 2
			    participants = 3

			would produce:
			    1, 1, 0

			which is invalid.
		*/
		if cmd.TotalAmount < entity.Money(len(cmd.Participants)) {
			return uuid.Nil, errors.New(
				"total amount is too small for the number of participants",
			)
		}

		participantIDs := make(
			map[uuid.UUID]struct{},
			len(cmd.Participants),
		)

		count := entity.Money(len(cmd.Participants))

		baseShare := cmd.TotalAmount / count
		remainder := cmd.TotalAmount % count

		shares = make(
			[]entity.ExpenseShare,
			0,
			len(cmd.Participants),
		)

		for i, participantID := range cmd.Participants {

			if participantID == uuid.Nil {
				return uuid.Nil, errors.New(
					"participant user ID is empty",
				)
			}

			if _, exists := participantIDs[participantID]; exists {
				return uuid.Nil, errors.New(
					"duplicate participant",
				)
			}

			participantIDs[participantID] = struct{}{}

			amount := baseShare

			/*
				Deterministically distribute the remainder
				one Toman at a time in participant order.
			*/
			if entity.Money(i) < remainder {
				amount++
			}

			shares = append(
				shares,
				entity.ExpenseShare{
					UserID: participantID,
					Amount: amount,
				},
			)

			memberIDs = append(
				memberIDs,
				participantID,
			)
		}

	case entity.CustomSplitType:

		if len(cmd.Shares) == 0 {
			return uuid.Nil, errors.New(
				"custom split requires at least one share",
			)
		}

		if len(cmd.Participants) != 0 {
			return uuid.Nil, errors.New(
				"participants should not be provided for custom split",
			)
		}

		shareIDs := make(
			map[uuid.UUID]struct{},
			len(cmd.Shares),
		)

		shares = make(
			[]entity.ExpenseShare,
			0,
			len(cmd.Shares),
		)

		var totalShares entity.Money

		for _, share := range cmd.Shares {

			if share.UserID == uuid.Nil {
				return uuid.Nil, errors.New(
					"share user ID is empty",
				)
			}

			if share.Amount <= 0 {
				return uuid.Nil, errors.New(
					"share amount must be greater than zero",
				)
			}

			if _, exists := shareIDs[share.UserID]; exists {
				return uuid.Nil, errors.New(
					"duplicate share participant",
				)
			}

			shareIDs[share.UserID] = struct{}{}

			totalShares += share.Amount

			shares = append(
				shares,
				entity.ExpenseShare{
					UserID: share.UserID,
					Amount: share.Amount,
				},
			)

			memberIDs = append(
				memberIDs,
				share.UserID,
			)
		}

		if totalShares != cmd.TotalAmount {
			return uuid.Nil, errors.New(
				"sum of shares must equal total amount",
			)
		}
	}

	// -----------------------
	// Remove duplicates before membership lookup
	// -----------------------

	memberIDs = utils.UniqueUUIDs(memberIDs)

	// -----------------------
	// Validate group membership
	// -----------------------

	areMembers, err := s.groupFinder.AreGroupMembers(ctx, cmd.GroupID, memberIDs)
	if err != nil {
		return uuid.Nil, fmt.Errorf(
			"check group membership: %w",
			err,
		)
	}

	if !areMembers {
		return uuid.Nil, errors.New(
			"creator, payer or participant is not a member of the group",
		)
	}

	// -----------------------
	// Time
	// -----------------------

	now := time.Now()

	occurredAt := cmd.OccurredAt

	if occurredAt.IsZero() {
		occurredAt = now
	} else {
		occurredAt = occurredAt.UTC()
	}

	// -----------------------
	// Build Expense aggregate
	// -----------------------

	expense := entity.Expense{
		ID:              uuid.New(),
		GroupID:         cmd.GroupID,
		Title:           strings.TrimSpace(cmd.Title),
		CreatedByUserID: cmd.CreatedByUserID,
		TotalAmount:     cmd.TotalAmount,
		SplitType:       cmd.SplitType,
		OccurredAt:      occurredAt,
		Payers:          payers,
		Shares:          shares,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// -----------------------
	// Persist aggregate
	// -----------------------

	if err := s.repo.CreateExpense(ctx, expense); err != nil {
		return uuid.Nil, fmt.Errorf(
			"create expense: %w",
			err,
		)
	}

	return expense.ID, nil
}
