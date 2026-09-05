package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/amirjbr/shared-expense/internal/groups/core/entity"
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

	//TODO complete handler part and test the functionality
	//TODO add err handling in both services
}
