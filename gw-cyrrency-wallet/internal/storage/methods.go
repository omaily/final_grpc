package storage

import (
	"context"
	"encoding/hex"
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

func (st *Instance) CreateAccount(ctx context.Context, user *model.Account) ([16]byte, error) {
	if st.checkExistUser(ctx, user) {
		return [16]byte{}, fmt.Errorf("username or email already exists")
	}

	passwordCript, err := user.SetPassword()
	if err != nil {
		slog.Error("bcrypt library generation error", slog.String("err", err.Error()))
		return [16]byte{}, err
	}

	query := `INSERT INTO account(uuid, name, pass, email)
VALUES(gen_random_uuid(), @name, @pass, @email)
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
		return [16]byte{}, err
	}

	return userID.Bytes, nil
}

func (st *Instance) FindAccount(ctx context.Context, user *model.Account) ([16]byte, string, error) {
	var uuid pgtype.UUID
	var storedUser model.Account
	query := `SELECT uuid, name, pass FROM account WHERE name = $1`
	err := st.pool.QueryRow(ctx, query, user.Username).Scan(&uuid, &storedUser.Username, &storedUser.Password)
	if err != nil {
		slog.Error("Error scanning row:", slog.String("err", err.Error()))
		return [16]byte{}, "", err
	}

	if storedUser.Username == "" {
		slog.Error("this user not found")
		return [16]byte{}, "", errors.New("this user not found")
	}

	return uuid.Bytes, storedUser.Password, nil
}

func moveStringToUUID(strUUID string) pgtype.UUID {
	buf, _ := hex.DecodeString(strUUID)
	return pgtype.UUID{Bytes: [16]byte(buf), Valid: true}
}

func (st *Instance) CheckBalance(ctx context.Context, strUUID string) (map[string]float64, error) {
	uuid := moveStringToUUID(strUUID)
	var currency string
	var count float64

	query := `SELECT currency, count FROM wallet WHERE user_id = $1`
	rows, err := st.pool.Query(ctx, query, uuid) //.Scan(&currency, &count)
	if err != nil {
		return nil, err
	}

	balance := make(map[string]float64)
	pgx.ForEachRow(rows, []any{&currency, &count}, func() error {
		balance[currency] = count
		return nil
	})

	return balance, nil
}

func (st *Instance) PutMoney(ctx context.Context, strUUID string, deposit model.Transfer) (map[string]float64, error) {
	uuid := moveStringToUUID(strUUID)
	query := `INSERT INTO wallet (user_id, currency, count)
VALUES (@user_id, @currency, @count)
ON CONFLICT ON CONSTRAINT uniq_wallet_user_cur
DO UPDATE SET count = wallet.count + @count;`
	args := pgx.NamedArgs{
		"user_id":  uuid,
		"currency": deposit.Currency,
		"count":    deposit.Amount,
	}

	_, err := st.pool.Exec(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to insert row: %w", err)
	}

	return st.CheckBalance(ctx, strUUID)
}

func (st *Instance) TakeMoney(ctx context.Context, strUUID string, deposit model.Transfer) (map[string]float64, error) {
	uuid := moveStringToUUID(strUUID)
	query := `UPDATE wallet SET count = wallet.count - @count WHERE user_id = @user_id and currency = @currency`
	args := pgx.NamedArgs{
		"user_id":  uuid,
		"currency": deposit.Currency,
		"count":    deposit.Amount,
	}

	_, err := st.pool.Exec(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("db %w", err)
	}

	return st.CheckBalance(ctx, strUUID)
}

func (st *Instance) takeMoney(tx pgx.Tx, strUUID string, deposit model.Transfer) {

}

func (st *Instance) ChangeMoney(ctx context.Context, strUUID string, ex model.Exchange, odds float64) error {
	uuid := moveStringToUUID(strUUID)

	tx, err := st.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	// takeMoney
	query := `UPDATE wallet SET count = wallet.count - @count WHERE user_id = @user_id and currency = @currency`
	args := pgx.NamedArgs{
		"user_id":  uuid,
		"currency": ex.FromCurrency,
		"count":    ex.Amount,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db %w", err)
	}

	// putMoney
	query = `INSERT INTO wallet (user_id, currency, count)
	VALUES (@user_id, @currency, @count)
	ON CONFLICT ON CONSTRAINT uniq_wallet_user_cur
	DO UPDATE SET count = wallet.count + @count;`
	args = pgx.NamedArgs{
		"user_id":  uuid,
		"currency": ex.ToCurrency,
		"count":    ex.Amount * odds,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}
