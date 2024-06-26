package chain

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/utils"
)

func GetChain(dbConnPgx utils.PgxIface, chainID *int) (*Chain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	base_asset_id,
	name, 
	alternate_name, 
	address,
	chain_type_id,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	COALESCE(rpc_url, '') ,
	chain_id,
	COALESCE(block_explorer_url, ''),
	COALESCE(rpc_url_dev, ''),
	COALESCE(rpc_url_prod, ''),
	COALESCE(rpc_url_archive, '')
	FROM chains 
	WHERE id = $1
	`, *chainID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	// from https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx
	defer row.Close()
	chain, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Chain])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &chain, nil
}

func GetChainByAddress(dbConnPgx utils.PgxIface, address string) (*Chain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	base_asset_id,
	name, 
	alternate_name, 
	address,
	chain_type_id,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	COALESCE(rpc_url, '') ,
	chain_id,
	COALESCE(block_explorer_url, ''),
	COALESCE(rpc_url_dev, ''),
	COALESCE(rpc_url_prod, ''),
	COALESCE(rpc_url_archive, '')
	FROM chains 
	WHERE address = $1
	`, address)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer row.Close()
	chain, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Chain])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &chain, nil
}

func GetChainByAlternateName(dbConnPgx utils.PgxIface, altenateName string) (*Chain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	base_asset_id,
	name, 
	alternate_name, 
	address,
	chain_type_id,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	COALESCE(rpc_url, '') ,
	chain_id,
	COALESCE(block_explorer_url, ''),
	COALESCE(rpc_url_dev, ''),
	COALESCE(rpc_url_prod, ''),
	COALESCE(rpc_url_archive, '')
	FROM chains 
	WHERE alternate_name = $1
	`, altenateName)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer row.Close()
	chain, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Chain])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &chain, nil
}

func RemoveChain(dbConnPgx utils.PgxIface, chainID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveChain DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM chains WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *chainID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetChainList(dbConnPgx utils.PgxIface, ids []int) ([]Chain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
	id,
	uuid, 
	base_asset_id,
	name, 
	alternate_name, 
	address,
	chain_type_id,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	COALESCE(rpc_url, ''),
	chain_id,
	COALESCE(block_explorer_url, ''),
	COALESCE(rpc_url_dev, ''),
	COALESCE(rpc_url_prod, ''),
	COALESCE(rpc_url_archive, '')
	FROM chains`
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
	chains, err := pgx.CollectRows(results, pgx.RowToStructByName[Chain])
	return chains, nil
}

func GetChainListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]Chain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	sql := `SELECT 
	id,
	uuid, 
	base_asset_id,
	name, 
	alternate_name, 
	address,
	chain_type_id,
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	COALESCE(rpc_url, ''),
	chain_id,
	COALESCE(block_explorer_url, ''),
	COALESCE(rpc_url_dev, ''),
	COALESCE(rpc_url_prod, ''),
	COALESCE(rpc_url_archive, '')
	FROM chains
	`
	if len(_filters) > 0 {
		sql += "WHERE "
		for i, filter := range _filters {
			sql += filter
			if i < len(_filters)-1 {
				sql += " AND "
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
	chains, err := pgx.CollectRows(results, pgx.RowToStructByName[Chain])
	return chains, nil
}

func GetTotalChainCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT COUNT(*) FROM chains`)
	totalCount := 0
	err := row.Scan(
		&totalCount,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &totalCount, nil
}

func UpdateChain(dbConnPgx utils.PgxIface, chain *Chain) error {
	// if the chain id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if chain.ID == nil || *chain.ID == 0 {
		return errors.New("chain has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateAsset DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE chains SET 
		name=$1,
		alternate_name=$2, 
		address=$3,
		chain_type_id=$4,
		description=$5,
		updated_by=$6, 
		updated_at=current_timestamp at time zone 'UTC',
		base_asset_id = $7,
		rpc_url = $8,
		chain_id = $9,
		block_explorer_url = $10,
		rpc_url_dev=$11,
		rpc_url_prod=$12,
		rpc_url_archive=$13
		WHERE id=$14`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		chain.Name,             //1
		chain.AlternateName,    //2
		chain.Address,          //3
		chain.ChainTypeID,      //4
		chain.Description,      //5
		chain.UpdatedBy,        //6
		chain.BaseAssetID,      //7
		chain.RpcURL,           //8
		chain.ChainID,          //9
		chain.BlockExplorerURL, //10
		chain.RpcURLDev,        //11
		chain.RpcURLProd,       //12
		chain.RpcURLArchive,    //13
		chain.ID,               //14
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertChain(dbConnPgx utils.PgxIface, chain *Chain) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertAsset DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var ID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO chains  
	(
		uuid, 
		name, 
		alternate_name, 
		address,
		chain_type_id,
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at,
		base_asset_id ,
		rpc_url,
		chain_id,
		block_explorer_url,
		rpc_url_dev,
		rpc_url_prod,
		rpc_url_archive
		) VALUES (
			uuid_generate_v4(),
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			current_timestamp at time zone 'UTC',
			$6,
			current_timestamp at time zone 'UTC',
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13
		)
		RETURNING id`,
		chain.Name,             //1
		chain.AlternateName,    //2
		chain.Address,          //3
		chain.ChainTypeID,      //4
		chain.Description,      //5
		chain.CreatedBy,        //6
		chain.BaseAssetID,      //7
		chain.RpcURL,           //8
		chain.ChainID,          //9
		chain.BlockExplorerURL, //10
		chain.RpcURLDev,        //11
		chain.RpcURLProd,       //12
		chain.RpcURLArchive,    //13
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
