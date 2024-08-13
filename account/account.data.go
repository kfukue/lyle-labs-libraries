package account

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	pgx "github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
)

func GetAccount(dbConnPgx utils.PgxIface, accountID *int) (*Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
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
	`, *accountID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	// from https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx
	account, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Account])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &account, nil
}

func GetAccountByAddress(dbConnPgx utils.PgxIface, address string) (*Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
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
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	account, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Account])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &account, nil
}

func GetAccountByAlternateName(dbConnPgx utils.PgxIface, altenateName string) (*Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
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
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	account, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Account])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &account, nil
}

func RemoveAccount(dbConnPgx utils.PgxIface, accountID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveAccount DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM accounts WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *accountID); err != nil {
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
	accounts, err := pgx.CollectRows(results, pgx.RowToStructByName[Account])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return accounts, nil
}

func UpdateAccount(dbConnPgx utils.PgxIface, account *Account) error {
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
	return tx.Commit(ctx)
}

func InsertAccount(dbConnPgx utils.PgxIface, account *Account) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertAccount DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var ID int
	defer dbConnPgx.Close()

	err = dbConnPgx.QueryRow(ctx, `INSERT INTO accounts
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
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, err
	}
	return int(ID), nil
}

func InsertAccounts(dbConnPgx utils.PgxIface, accounts []Account) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i, _ := range accounts {
		account := accounts[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(account.UUID)
		row := []interface{}{
			uuidString,             //1
			account.Name,           //2
			account.AlternateName,  //3
			account.Address,        //4
			account.NameFromSource, //5
			account.PortfolioID,    //6
			account.SourceID,       //7
			account.AccountTypeID,  //8
			account.Description,    //9
			account.CreatedBy,      //10
			&now,                   //11
			account.CreatedBy,      //12
			&now,                   //13
			account.ChainID,        //14
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"accounts"},
		[]string{
			"uuid",             //1
			"name",             //2
			"alternate_name",   //3
			"address",          //4
			"name_from_source", //5
			"portfolio_id",     //6
			"source_id",        //7
			"account_type_id",  //8
			"description",      //9
			"created_by",       //10
			"created_at",       //11
			"updated_by",       //12
			"updated_at",       //13
			"chain_id",         //14
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertAccounts: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// for refinedev
func GetAccountListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]Account, error) {
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
	FROM accounts 
	`
	if len(_filters) > 0 {
		sql += "WHERE "
		for i, filter := range _filters {
			sql += filter
			if i < len(_filters)-1 {
				sql += " OR "
			}
		}
	}
	if _order != "" && _sort != "" {
		sql += fmt.Sprintf(" ORDER BY %s %s ", _sort, _order)
	}
	if (_start != nil && *_start > 0) && (_end != nil && *_end > 0) {
		pageSize := *_end - *_start
		sql += fmt.Sprintf(" OFFSET %d LIMIT %d ", *_start, pageSize)
	}

	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	accounts, err := pgx.CollectRows(results, pgx.RowToStructByName[Account])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()

	return accounts, nil
}

func GetTotalAccountsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM accounts
	`)
	totalCount := 0
	err := row.Scan(
		&totalCount,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &totalCount, nil
}
