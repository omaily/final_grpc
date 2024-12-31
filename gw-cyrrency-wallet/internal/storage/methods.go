package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/pkg/model"
)

func (st *Instance) checkExistUser(ctx context.Context, acc *model.Account) bool {
	var counter int
	query := `SELECT count(*) FROM account WHERE name = $1 OR email = $2`
	st.pool.QueryRow(ctx, query, acc.Username, acc.Email).Scan(&counter)
	return counter > 0
}

func (st *Instance) CreateAccount(ctx context.Context, user *model.Account) (*[16]byte, error) {
	if st.checkExistUser(ctx, user) {
		return nil, fmt.Errorf("username or email already exists")
	}

	passwordCript, err := user.SetPassword()
	if err != nil {
		slog.Error("bcrypt library generation error", slog.String("err", err.Error()))
		return nil, err
	}

	query := `INSERT INTO account(uuid, 
	name, 
	pass, 
	email)
VALUES(gen_random_uuid(), 
	@name, 
	@pass, 
	@email)
RETURNING uuid;`

	args := pgx.NamedArgs{
		"name":  user.Username,
		"pass":  passwordCript,
		"email": user.Email,
	}
	batch := &pgx.Batch{}
	batch.Queue(query, args)
	results := st.pool.SendBatch(ctx, batch)
	defer results.Close()

	var userID pgtype.UUID
	err = results.QueryRow().Scan(&userID)
	if err != nil {
		slog.Error("unable to insert row: ", slog.String("err", err.Error()))
		return nil, err
	}

	return &userID.Bytes, nil
}

func (st *Instance) FindAccount(ctx context.Context, user *model.Account) (string, error) {
	var storedUser model.Account
	query := `SELECT name, pass FROM account WHERE name = $1`
	err := st.pool.QueryRow(ctx, query, user.Username).Scan(&storedUser.Username, &storedUser.Password)
	if err != nil {
		slog.Error("Error scanning row:", slog.String("err", err.Error()))
		return "", err
	}

	if storedUser.Username == "" {
		slog.Error("this user not found")
		return "", errors.New("this user not found")
	}

	return storedUser.Password, nil
}
