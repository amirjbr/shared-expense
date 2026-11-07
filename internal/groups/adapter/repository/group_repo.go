package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/amirjbr/shared-expense/internal/groups/core/entity"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type GroupRepo struct {
	DB *sql.DB
}

func NewGroupRepo(db *sql.DB) (*GroupRepo, error) {
	if db == nil {
		return nil, errors.New("db connection is nil")
	}

	return &GroupRepo{
		DB: db,
	}, nil
}

// ---------------------------------------------------------
// Create Group
// ---------------------------------------------------------

func (r *GroupRepo) CreateGroup(ctx context.Context, group entity.Group) error {

	const queryGroup = `
		INSERT INTO groups (
			id,
			name,
			owner_user_id,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	const queryGroupMember = `
		INSERT INTO group_members (
			group_id,
			user_id,
			joined_at
		)
		VALUES ($1, $2, $3)
	`

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin create group transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(
		ctx,
		queryGroup,
		group.ID,
		group.Name,
		group.OwnerID,
		group.CreatedAt,
		group.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert group: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		queryGroupMember,
		group.ID,
		group.OwnerID,
		group.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert group owner as member: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit create group transaction: %w", err)
	}

	return nil
}

// ---------------------------------------------------------
// Get Group By ID
// ---------------------------------------------------------

func (r *GroupRepo) GetGroupByID(ctx context.Context, groupID uuid.UUID) (entity.Group, error) {

	const query = `
		SELECT
			id,
			name,
			owner_user_id,
			created_at,
			updated_at
		FROM groups
		WHERE id = $1
	`

	var group entity.Group

	err := r.DB.QueryRowContext(ctx, query, groupID).Scan(
		&group.ID,
		&group.Name,
		&group.OwnerID,
		&group.CreatedAt,
		&group.UpdatedAt,
	)
	if err != nil {
		return entity.Group{}, fmt.Errorf(
			"get group by id: %w",
			err,
		)
	}

	return group, nil
}

// ---------------------------------------------------------
// Is Group Member
// ---------------------------------------------------------

func (r *GroupRepo) IsGroupMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) (bool, error) {

	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM group_members
			WHERE group_id = $1
			  AND user_id = $2
		)
	`

	var exists bool

	err := r.DB.QueryRowContext(
		ctx,
		query,
		groupID,
		userID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf(
			"check group membership: %w",
			err,
		)
	}

	return exists, nil
}

// ---------------------------------------------------------
// Create Group Invitation
// ---------------------------------------------------------

func (r *GroupRepo) CreateGroupInvitation(ctx context.Context, groupInvitation entity.GroupInvitation) error {

	const query = `
		INSERT INTO group_invitations (
			id,
			group_id,
			invited_user_id,
			invited_by_user_id,
			status,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		groupInvitation.ID,
		groupInvitation.GroupID,
		groupInvitation.InvitedUserID,
		groupInvitation.InvitedByUserID,
		groupInvitation.Status,
		groupInvitation.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf(
			"create group invitation: %w",
			err,
		)
	}

	return nil
}

// ---------------------------------------------------------
// Has Pending Invitation
// ---------------------------------------------------------

func (r *GroupRepo) HasPendingInvitation(ctx context.Context, groupID uuid.UUID, invitedUserID uuid.UUID) (bool, error) {

	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM group_invitations
			WHERE group_id = $1
			  AND invited_user_id = $2
			  AND status = 'pending'
		)
	`

	var exists bool

	err := r.DB.QueryRowContext(
		ctx,
		query,
		groupID,
		invitedUserID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf(
			"check pending invitation: %w",
			err,
		)
	}

	return exists, nil
}

// ---------------------------------------------------------
// Accept Invitation
// ---------------------------------------------------------

func (r *GroupRepo) AcceptInvitation(ctx context.Context, invitationID uuid.UUID, userID uuid.UUID) error {

	const queryInvitation = `
		UPDATE group_invitations
		SET
			status = 'accepted',
			responded_at = NOW()
		WHERE id = $1
		  AND invited_user_id = $2
		  AND status = 'pending'
		RETURNING group_id
	`

	const queryAddMember = `
		INSERT INTO group_members (
			group_id,
			user_id
		)
		VALUES ($1, $2)
	`

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf(
			"begin accept invitation transaction: %w",
			err,
		)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var groupID uuid.UUID

	err = tx.QueryRowContext(
		ctx,
		queryInvitation,
		invitationID,
		userID,
	).Scan(&groupID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New(
				"invitation not found or already responded",
			)
		}

		return fmt.Errorf(
			"accept invitation: %w",
			err,
		)
	}

	_, err = tx.ExecContext(
		ctx,
		queryAddMember,
		groupID,
		userID,
	)
	if err != nil {
		return fmt.Errorf(
			"add accepted user to group: %w",
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(
			"commit accept invitation transaction: %w",
			err,
		)
	}

	return nil
}

// ---------------------------------------------------------
// Reject Invitation
// ---------------------------------------------------------

func (r *GroupRepo) RejectInvitation(ctx context.Context, invitationID uuid.UUID, userID uuid.UUID) error {

	const query = `
		UPDATE group_invitations
		SET
			status = 'rejected',
			responded_at = NOW()
		WHERE id = $1
		  AND invited_user_id = $2
		  AND status = 'pending'
	`

	result, err := r.DB.ExecContext(
		ctx,
		query,
		invitationID,
		userID,
	)
	if err != nil {
		return fmt.Errorf(
			"reject invitation: %w",
			err,
		)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"get rejected invitation rows affected: %w",
			err,
		)
	}

	if rowsAffected == 0 {
		return errors.New(
			"invitation not found or already responded",
		)
	}

	return nil
}

// ---------------------------------------------------------
// Get Pending Invitations
// ---------------------------------------------------------

func (r *GroupRepo) GetPendingInvitations(ctx context.Context, userID uuid.UUID) ([]entity.GroupInvitation, error) {

	const query = `
		SELECT
			id,
			group_id,
			invited_user_id,
			invited_by_user_id,
			status,
			created_at
		FROM group_invitations
		WHERE invited_user_id = $1
		  AND status = 'pending'
	`

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get pending invitations: %w",
			err,
		)
	}

	defer rows.Close()

	invitations := make(
		[]entity.GroupInvitation,
		0,
	)

	for rows.Next() {
		var invitation entity.GroupInvitation

		err := rows.Scan(
			&invitation.ID,
			&invitation.GroupID,
			&invitation.InvitedUserID,
			&invitation.InvitedByUserID,
			&invitation.Status,
			&invitation.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan group invitation: %w",
				err,
			)
		}

		invitations = append(
			invitations,
			invitation,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate pending invitations: %w",
			err,
		)
	}

	return invitations, nil
}

// ---------------------------------------------------------
// Get Groups Of User
// ---------------------------------------------------------

func (r *GroupRepo) GetAllGroupsOfUser(ctx context.Context, userID uuid.UUID) ([]entity.Group, error) {

	const query = `
		SELECT
			g.id,
			g.name,
			g.owner_user_id,
			g.created_at,
			g.updated_at
		FROM groups g
		INNER JOIN group_members gm
			ON gm.group_id = g.id
		WHERE gm.user_id = $1
	`

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get groups of user: %w",
			err,
		)
	}

	defer rows.Close()

	groups := make([]entity.Group, 0)

	for rows.Next() {
		var group entity.Group

		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.OwnerID,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan group: %w",
				err,
			)
		}

		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate user groups: %w",
			err,
		)
	}

	return groups, nil
}

// ---------------------------------------------------------
// Get Group Members
// ---------------------------------------------------------

func (r *GroupRepo) GetMembersOfGroup(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error) {

	const query = `
		SELECT user_id
		FROM group_members
		WHERE group_id = $1
	`

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		groupID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get group members: %w",
			err,
		)
	}

	defer rows.Close()

	userIDs := make([]uuid.UUID, 0)

	for rows.Next() {
		var userID uuid.UUID

		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf(
				"scan group member: %w",
				err,
			)
		}

		userIDs = append(
			userIDs,
			userID,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate group members: %w",
			err,
		)
	}

	return userIDs, nil
}

// ---------------------------------------------------------
// Check Multiple Group Members
// ---------------------------------------------------------

func (r *GroupRepo) AreGroupMembers(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) (bool, error) {

	if len(userIDs) == 0 {
		return true, nil
	}

	const query = `
		SELECT NOT EXISTS (
			SELECT 1
			FROM unnest($2::uuid[]) AS requested(user_id)
			WHERE NOT EXISTS (
				SELECT 1
				FROM group_members gm
				WHERE gm.group_id = $1
				  AND gm.user_id = requested.user_id
			)
		)
	`

	// lib/pq handles []string more reliably for Postgres arrays,
	// so conversion stays inside the Postgres adapter.
	idStrings := make([]string, len(userIDs))

	for i, id := range userIDs {
		idStrings[i] = id.String()
	}

	var areMembers bool

	err := r.DB.QueryRowContext(
		ctx,
		query,
		groupID,
		pq.Array(idStrings),
	).Scan(&areMembers)

	if err != nil {
		return false, fmt.Errorf(
			"check group members: %w",
			err,
		)
	}

	return areMembers, nil
}
