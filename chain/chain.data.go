package chain

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/kfukue/lyle-labs-libraries/database"
)

func GetChain(chainID int) (*Chain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	rpc_url,
	chain_id,
	block_explorer_url
	FROM chains 
	WHERE id = $1
	`, chainID)

	chain := &Chain{}
	err := row.Scan(
		&chain.ID,
		&chain.UUID,
		&chain.BaseAssetID,
		&chain.Name,
		&chain.AlternateName,
		&chain.Address,
		&chain.ChainTypeID,
		&chain.Description,
		&chain.CreatedBy,
		&chain.CreatedAt,
		&chain.UpdatedBy,
		&chain.UpdatedAt,
		&chain.RpcURL,
		&chain.ChainID,
		&chain.BlockExplorerURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return chain, nil
}

func GetChainByAddress(address string) (*Chain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	rpc_url,
	chain_id,
	block_explorer_url
	FROM chains 
	WHERE address = $1
	`, address)

	chain := &Chain{}
	err := row.Scan(
		&chain.ID,
		&chain.UUID,
		&chain.BaseAssetID,
		&chain.Name,
		&chain.AlternateName,
		&chain.Address,
		&chain.ChainTypeID,
		&chain.Description,
		&chain.CreatedBy,
		&chain.CreatedAt,
		&chain.UpdatedBy,
		&chain.UpdatedAt,
		&chain.RpcURL,
		&chain.ChainID,
		&chain.BlockExplorerURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return chain, nil
}

func GetChainByAlternateName(altenateName string) (*Chain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	rpc_url,
	chain_id,
	block_explorer_url
	FROM chains 
	WHERE alternate_name = $1
	`, altenateName)

	chain := &Chain{}
	err := row.Scan(
		&chain.ID,
		&chain.UUID,
		&chain.BaseAssetID,
		&chain.Name,
		&chain.AlternateName,
		&chain.Address,
		&chain.ChainTypeID,
		&chain.Description,
		&chain.CreatedBy,
		&chain.CreatedAt,
		&chain.UpdatedBy,
		&chain.UpdatedAt,
		&chain.RpcURL,
		&chain.ChainID,
		&chain.BlockExplorerURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return chain, nil
}

func RemoveChain(chainID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM chains WHERE id = $1`, chainID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetChainList() ([]Chain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
	rpc_url,
	chain_id,
	block_explorer_url
	FROM chains`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	chains := make([]Chain, 0)
	for results.Next() {
		var chain Chain
		results.Scan(
			&chain.ID,
			&chain.UUID,
			&chain.BaseAssetID,
			&chain.Name,
			&chain.AlternateName,
			&chain.Address,
			&chain.ChainTypeID,
			&chain.Description,
			&chain.CreatedBy,
			&chain.CreatedAt,
			&chain.UpdatedBy,
			&chain.UpdatedAt,
			&chain.RpcURL,
			&chain.ChainID,
			&chain.BlockExplorerURL,
		)

		chains = append(chains, chain)
	}
	return chains, nil
}

func UpdateChain(chain Chain) error {
	// if the chain id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if chain.ID == nil || *chain.ID == 0 {
		return errors.New("chain has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE chains SET 
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
		WHERE id=$11`,
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
		chain.ID,               //11
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertChain(chain Chain) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var ID int
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO chains  
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
		block_explorer_url
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
			$10
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
	).Scan(&ID)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(ID), nil
}
