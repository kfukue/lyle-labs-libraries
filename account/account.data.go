package account

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	pgx "github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgx/v5/pgconn"
	"github.com/kfukue/lyle-labs-libraries/utils"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
	Prepare(context.Context, string, string) (*pgconn.StatementDescription, error)
	Close()
}

type TestAccount struct {
	ID   *int   `db:"id"`
	Name string `db:"name"`
}

func GetAccount(dbConnPgx utils.PgxIface, accountID int) (*Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := dbConnPgx.QueryRow(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	address,
	name_from_source,
	portfolio_id,
	source_id,
	account_type_id,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	chain_id
	FROM accounts 
	WHERE id = $1
	`, accountID)

	account := &Account{}
	err := row.Scan(
		&account.ID,
		&account.UUID,
		&account.Name,
		&account.AlternateName,
		&account.Address,
		&account.NameFromSource,
		&account.PortfolioID,
		&account.SourceID,
		&account.AccountTypeID,
		&account.Description,
		&account.CreatedBy,
		&account.CreatedAt,
		&account.UpdatedBy,
		&account.UpdatedAt,
		&account.ChainID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return account, nil
}

func GetAccountByAddress(dbConnPgx utils.PgxIface, address string) (*Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := dbConnPgx.QueryRow(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	address,
	name_from_source,
	portfolio_id,
	source_id,
	account_type_id,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	chain_id 
	FROM accounts 
	WHERE address = $1
	`, address)

	account := &Account{}
	err := row.Scan(
		&account.ID,
		&account.UUID,
		&account.Name,
		&account.AlternateName,
		&account.Address,
		&account.NameFromSource,
		&account.PortfolioID,
		&account.SourceID,
		&account.AccountTypeID,
		&account.Description,
		&account.CreatedBy,
		&account.CreatedAt,
		&account.UpdatedBy,
		&account.UpdatedAt,
		&account.ChainID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return account, nil
}

func GetAccountByAlternateName(dbConnPgx utils.PgxIface, altenateName string) (*Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := dbConnPgx.QueryRow(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	address,
	name_from_source,
	portfolio_id,
	source_id,
	account_type_id,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	chain_id 
	FROM accounts 
	WHERE alternate_name = $1
	`, altenateName)

	account := &Account{}
	err := row.Scan(
		&account.ID,
		&account.UUID,
		&account.Name,
		&account.AlternateName,
		&account.Address,
		&account.NameFromSource,
		&account.PortfolioID,
		&account.SourceID,
		&account.AccountTypeID,
		&account.Description,
		&account.CreatedBy,
		&account.CreatedAt,
		&account.UpdatedBy,
		&account.UpdatedAt,
		&account.ChainID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return account, nil
}

func RemoveAccount(dbConnPgx utils.PgxIface, accountID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveAccount DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM accounts WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, accountID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetAccountList(dbConnPgx utils.PgxIface, ids []int) ([]Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT
	id,
	uuid,
	name,
	alternate_name,
	address,
	name_from_source,
	portfolio_id,
	source_id,
	account_type_id,
	description,
	created_by,
	created_at,
	updated_by,
	updated_at,
	chain_id
	FROM accounts`
	if len(ids) > 0 {
		strIds := utils.SplitToString(ids, ",")
		additionalQuery := fmt.Sprintf(` WHERE id IN (%s)`, strIds)
		sql += additionalQuery
	}
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	accounts := make([]Account, 0)
	for results.Next() {
		var account Account
		results.Scan(
			&account.ID,
			&account.UUID,
			&account.Name,
			&account.AlternateName,
			&account.Address,
			&account.NameFromSource,
			&account.PortfolioID,
			&account.SourceID,
			&account.AccountTypeID,
			&account.Description,
			&account.CreatedBy,
			&account.CreatedAt,
			&account.UpdatedBy,
			&account.UpdatedAt,
			&account.ChainID,
		)

		accounts = append(accounts, account)
	}
	return accounts, nil
}

func UpdateAccount(dbConnPgx utils.PgxIface, account Account) error {
	// if the account id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if account.ID == nil || *account.ID == 0 {
		return errors.New("account has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateAccount DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE accounts SET 
		name=$1,
		alternate_name=$2, 
		address=$3,
		name_from_source=$4,
		portfolio_id=$5,
		source_id=$6,
		account_type_id=$7,
		description=$8,
		updated_by=$9, 
		updated_at=current_timestamp at time zone 'UTC',
		chain_id = $10
		WHERE id=$11`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		account.Name,
		account.AlternateName,
		account.Address,
		account.NameFromSource,
		account.PortfolioID,
		account.SourceID,
		account.AccountTypeID,
		account.Description,
		account.UpdatedBy,
		account.ChainID,
		account.ID); err != nil {
		tx.Rollback(ctx)
		return err
	}

	// )
	// _, err := dbConnPgx.Query(ctx, `UPDATE accounts SET
	// name=$1,
	// alternate_name=$2,
	// address=$3,
	// name_from_source=$4,
	// portfolio_id=$5,
	// source_id=$6,
	// account_type_id=$7,
	// description=$8,
	// 	updated_by=$9,
	// 	updated_at=current_timestamp at time zone 'UTC',
	// 	chain_id = $10
	// 	WHERE id=$11`,
	// 	account.Name,
	// 	account.AlternateName,
	// 	account.Address,
	// 	account.NameFromSource,
	// 	account.PortfolioID,
	// 	account.SourceID,
	// 	account.AccountTypeID,
	// 	account.Description,
	// 	account.UpdatedBy,
	// 	account.ChainID,
	// 	account.ID,
	// )
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return err
	// }
	return tx.Commit(ctx)
}

func InsertAccount(dbConnPgx utils.PgxIface, account Account) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var ID int
	err := dbConnPgx.QueryRow(ctx, `INSERT INTO accounts  
	(
		uuid, 
		name, 
		alternate_name, 
		address,
		name_from_source,
		portfolio_id,
		source_id,
		account_type_id,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at,
		chain_id 
		) VALUES (
			uuid_generate_v4(),
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			current_timestamp at time zone 'UTC',
			$9,
			current_timestamp at time zone 'UTC',
			$10
		)
		RETURNING id`,
		account.Name,           //1
		account.AlternateName,  //2
		account.Address,        //3
		account.NameFromSource, //4
		account.PortfolioID,    //5
		account.SourceID,       //6
		account.AccountTypeID,  //7
		account.Description,    //8
		account.CreatedBy,      //9
		account.ChainID,        //10
	).Scan(&ID)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(ID), nil
}
