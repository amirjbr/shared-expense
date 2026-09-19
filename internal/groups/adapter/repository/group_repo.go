package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/amirjbr/shared-expense/internal/groups/core/entity"
	"github.com/google/uuid"
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

func (r *GroupRepo) CreateGroup(ctx context.Context, group entity.Group) (string, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	queryGroup := `INSERT INTO groups(id,name,owner_user_id,created_at,updated_at) 
				VALUES ($1 , $2 ,$3 ,$4,$5)`

	queryGroupMember := `INSERT INTO group_members(group_id,user_id,joined_at)
				VALUES ($1, $2, $3)`

	_, err = tx.ExecContext(ctx, queryGroup, group.ID, group.Name, group.OwnerID,
		group.CreatedAt, group.UpdatedAt,
	)
	if err != nil {
		defer func(tx *sql.Tx) {
			err := tx.Rollback()
			if err != nil {
				panic(err)
			}
		}(tx)
		return "", err
	}

	_, err = tx.ExecContext(ctx, queryGroupMember, group.ID, group.OwnerID, group.CreatedAt)
	if err != nil {
		defer func(tx *sql.Tx) {
			err := tx.Rollback()
			if err != nil {
				panic(err)
			}
		}(tx)
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}
	return group.ID.String(), nil
}
func (r *GroupRepo) GetGroupByID(ctx context.Context, groupID string) (*entity.Group, error) {
	query := `SELECT * FROM groups WHERE id = $1`
	group := &entity.Group{}
	err := r.DB.QueryRowContext(ctx, query, groupID).Scan(
		&group.ID, &group.Name, &group.OwnerID, &group.CreatedAt, &group.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// TODO handle this later
			return nil, nil
		}
		return nil, err
	}
	return group, nil

}

func (r *GroupRepo) IsGroupMember(ctx context.Context, groupID string, userID string) (bool, error) {
	query := `SELECT EXISTS (
    	SELECT 1
    	FROM group_members
    	WHERE group_id = $1
      	AND user_id = $2
	);`
	var exists bool

	err := r.DB.QueryRowContext(ctx, query, groupID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *GroupRepo) CreateGroupInvitation(ctx context.Context, groupInvitation entity.GroupInvitation) error {
	query := `INSERT INTO group_invitations(id,group_id,
                              invited_user_id,
                              invited_by_user_id,
                              status,created_at
                              )
				VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.DB.ExecContext(ctx, query, groupInvitation.ID, groupInvitation.GroupID,
		groupInvitation.InvitedUserID,
		groupInvitation.InvitedByUserID,
		groupInvitation.Status, groupInvitation.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *GroupRepo) IsHaveInvitation(ctx context.Context, invitationUserID string, groupID string) (bool, error) {
	query := `SELECT EXISTS (
				SELECT 1 
				From group_invitations WHERE invited_user_id =$1 
				AND group_id = $2
				AND status = "pending");`

	var exists bool

	err := r.DB.QueryRowContext(ctx, query, invitationUserID, groupID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *GroupRepo) AcceptInvitation(ctx context.Context, invitationID string, userID string) error {
	queryGroupInvitation := `UPDATE group_invitations SET 
                             status = 'accepted',
                             responded_at = NOW()
							WHERE id = $1 
							AND invited_user_id = $2 AND status = 'pending'
							RETURNING group_id;
								`
	queryAddingGroupMember := `INSERT INTO group_members(group_id,user_id)
								VALUES($1,$2)
								`

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var groupID uuid.UUID
	err = tx.QueryRowContext(ctx, queryGroupInvitation, invitationID, userID).Scan(&groupID)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("invitation not found or already responded")
	}

	_, err = tx.ExecContext(ctx, queryAddingGroupMember, groupID, userID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil

}

func (r *GroupRepo) RejectInvitation(ctx context.Context, invitationID string, userID string) error {
	queryGroupInvitation := `UPDATE group_invitations SET 
                             status = 'rejected',
                             responded_at = NOW()
							WHERE id = $1 
							AND invited_user_id = $2 AND status = 'pending'
							RETURNING group_id;
								`
	queryAddingGroupMember := `INSERT INTO group_members(group_id,user_id)
								VALUES($1,$2)
								`

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var groupID uuid.UUID
	err = tx.QueryRowContext(ctx, queryGroupInvitation, invitationID, userID).Scan(&groupID)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("invitation not found or already responded")
	}

	_, err = tx.ExecContext(ctx, queryAddingGroupMember, groupID, userID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil

}
