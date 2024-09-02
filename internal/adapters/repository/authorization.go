package repository

import (
	"context"
	"fmt"
)

const queryGetUserCredentials = `select hash_password from user_data where email = $1`

func (r *Repository) GetUserCredentials(ctx context.Context, email string) (password string, err error) {
	err = r.db.GetContext(ctx, &password, queryGetUserCredentials, email)
	if err != nil {
		err = fmt.Errorf("r.db.ExecContext(...): %w", err)
		return "", err
	}
	return password, nil
}

const querySaveUserCredentials = `insert into user_data (email, hash_password) values ($1, $2)`

func (r *Repository) SaveUserCredentials(ctx context.Context, email string, hashPassword string) (err error) {

	_, err = r.db.ExecContext(ctx, querySaveUserCredentials, email, hashPassword)
	if err != nil {
		err = fmt.Errorf("r.db.ExecContext(...): %w", err)
		return err
	}
	return nil
}
