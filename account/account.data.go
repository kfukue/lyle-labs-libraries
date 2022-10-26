package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
)

func GetAccount(accountID int) (*Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return account, nil
}

func GetAccountByAddress(address string) (*Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return account, nil
}

func GetAccountByAlternateName(altenateName string) (*Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return account, nil
}

func RemoveAccount(accountID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM accounts WHERE id = $1`, accountID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetAccountList(ids []int) ([]Account, error) {
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
	results, err := database.DbConn.QueryContext(ctx, sql)
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

func UpdateAccount(account Account) error {
	// if the account id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if account.ID == nil || *account.ID == 0 {
		return errors.New("account has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE accounts SET 
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
		WHERE id=$11`,
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
		account.ID,
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertAccount(account Account) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var ID int
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO accounts  
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
