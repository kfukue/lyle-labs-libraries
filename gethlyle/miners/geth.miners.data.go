package gethlyleminers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
)

func GetGethMiner(dbConnPgx utils.PgxIface, gethMinerID *int) (*GethMiner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		uuid,
		name,
		alternate_name,
		chain_id,
		exchange_id,
		starting_block_number,
		created_txn_hash,
		last_block_number,
		contract_address,
		contract_address_id,
		developer_address,
		developer_address_id,
		mining_asset_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM geth_miners
	WHERE id = $1
	`, *gethMinerID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	// from https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx
	defer row.Close()
	gethMiner, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethMiner])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethMiner, nil
}

func RemoveGethMiner(dbConnPgx utils.PgxIface, gethMinerID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethMiner DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_miners WHERE id = $1`
	//defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *gethMinerID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetGethMinerList(dbConnPgx utils.PgxIface) ([]GethMiner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT
		id,
		uuid,
		name,
		alternate_name,
		chain_id,
		exchange_id,
		starting_block_number,
		created_txn_hash,
		last_block_number,
		contract_address,
		contract_address_id,
		developer_address,
		developer_address_id,
		mining_asset_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM geth_miners `)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethMiners, err := pgx.CollectRows(results, pgx.RowToStructByName[GethMiner])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethMiners, nil
}

func GetGethMinerListByMiningAssetId(dbConnPgx utils.PgxIface, miningAssetID *int) ([]GethMiner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT
		id,
		uuid,
		name,
		alternate_name,
		chain_id,
		exchange_id,
		starting_block_number,
		created_txn_hash,
		last_block_number,
		contract_address,
		contract_address_id,
		developer_address,
		developer_address_id,
		mining_asset_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM geth_miners
	WHERE
		mining_asset_id = $1  `, *miningAssetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethMiners, err := pgx.CollectRows(results, pgx.RowToStructByName[GethMiner])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethMiners, nil
}

func UpdateGethMiner(dbConnPgx utils.PgxIface, gethMiner *GethMiner) error {
	// if the gethMiner id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethMiner.ID == nil {
		return errors.New("gethMiner has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethMiner DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_miners SET 
		name=$1,
		alternate_name=$2,
		chain_id=$3,
		exchange_id=$4,
		starting_block_number=$5,
		created_txn_hash=$6,
		last_block_number=$7,
		contract_address=$8,
		contract_address_id=$9,
		developer_address=$10,
		developer_address_id=$11,
		mining_asset_id=$12,
		description=$13,
		updated_by=$14,
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$15`
	//defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		gethMiner.Name,                //1
		gethMiner.AlternateName,       //2
		gethMiner.ChainID,             //3
		gethMiner.ExchangeID,          //4
		gethMiner.StartingBlockNumber, //5
		gethMiner.CreatedTxnHash,      //6
		gethMiner.LastBlockNumber,     //7
		gethMiner.ContractAddress,     //8
		gethMiner.ContractAddressID,   //9
		gethMiner.DeveloperAddress,    //10
		gethMiner.DeveloperAddressID,  //11
		gethMiner.MiningAssetID,       //12
		gethMiner.Description,         //13
		gethMiner.UpdatedBy,           //14
		gethMiner.ID,                  //15
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertGethMiner(dbConnPgx utils.PgxIface, gethMiner *GethMiner) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertGethMiner DbConn.Begin   %s", err.Error())
		return -1, "", err
	}
	var gethMinerID int
	var gethMinerUUID string
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO geth_miners
	(
		uuid,
		name,
		alternate_name,
		chain_id,
		exchange_id,
		starting_block_number,
		created_txn_hash,
		last_block_number,
		contract_address,
		contract_address_id,
		developer_address,
		developer_address_id,
		mining_asset_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
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
		$10,
		$11,
		$12,
		$13,
		$14,
		current_timestamp at time zone 'UTC', 
		$14,
		current_timestamp at time zone 'UTC'
		)
		RETURNING id, uuid`,
		gethMiner.Name,                //1
		gethMiner.AlternateName,       //2
		gethMiner.ChainID,             //3
		gethMiner.ExchangeID,          //4
		gethMiner.StartingBlockNumber, //5
		gethMiner.CreatedTxnHash,      //6
		gethMiner.LastBlockNumber,     //7
		gethMiner.ContractAddress,     //8
		gethMiner.ContractAddressID,   //9
		gethMiner.DeveloperAddress,    //10
		gethMiner.DeveloperAddressID,  //11
		gethMiner.MiningAssetID,       //12
		gethMiner.Description,         //13
		gethMiner.CreatedBy,           //14
	).Scan(&gethMinerID, &gethMinerUUID)

	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, "", err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, "", err
	}
	return int(gethMinerID), gethMinerUUID, nil
}
func InsertGethMiners(dbConnPgx utils.PgxIface, gethMiners []GethMiner) error {
	// need to supply uuid
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range gethMiners {
		gethMiner := gethMiners[i]
		uuidString := pgtype.UUID{}
		uuidString.Set(gethMiner.UUID)
		row := []interface{}{
			uuidString,                    //1
			gethMiner.Name,                //2
			gethMiner.AlternateName,       //3
			gethMiner.ChainID,             //4
			gethMiner.ExchangeID,          //5
			gethMiner.StartingBlockNumber, //6
			gethMiner.CreatedTxnHash,      //7
			gethMiner.LastBlockNumber,     //8
			gethMiner.ContractAddress,     //9
			gethMiner.ContractAddressID,   //10
			gethMiner.DeveloperAddress,    //11
			gethMiner.DeveloperAddressID,  //12
			gethMiner.MiningAssetID,       //13
			gethMiner.Description,         //14
			gethMiner.CreatedBy,           //15
			now,                           //16
			gethMiner.CreatedBy,           //17
			now,                           //18

		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_miners"},
		[]string{
			"uuid",                  //1
			"name",                  //2
			"alternate_name",        //3
			"chain_id",              //4
			"exchange_id",           //5
			"starting_block_number", //6
			"created_txn_hash",      //7
			"last_block_number",     //8
			"contract_address",      //9
			"contract_address_id",   //10
			"developer_address",     //11
			"developer_address_id",  //12
			"mining_asset_id",       //13
			"description",           //14
			"created_by",            //15
			"created_at",            //16
			"updated_by",            //17
			"updated_at",            //18
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertGethMiners: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func UpdateGethMinerAddresses(dbConnPgx utils.PgxIface, gethMinerID *int) error {
	// update address ids from existing addresses in geth_addresses
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethMinerAddresses DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_miners SET 
		contract_address_id = ga.id from geth_addresses as ga
			WHERE LOWER(gm.contract_address) = LOWER(ga.address_str)
			AND gm.contract_address_id IS NULL
			AND gm.id = $1
			`
	//defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *gethMinerID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	sql2 := `UPDATE geth_miners SET 
			UPDATE geth_miners as gm SET
			developer_address_id = ga.id from geth_addresses as ga
			WHERE LOWER(gm.developer_address) = LOWER(ga.address_str)
			AND gm.developer_address_id IS NULL
			AND gm.id = $1
			`
	if _, err := dbConnPgx.Exec(ctx, sql2, *gethMinerID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

// for refinedev
func GetGethMinerListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]GethMiner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		id,
		uuid,
		name,
		alternate_name,
		chain_id,
		exchange_id,
		starting_block_number,
		created_txn_hash,
		last_block_number,
		contract_address,
		contract_address_id,
		developer_address,
		developer_address_id,
		mining_asset_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM geth_miners 
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
	gethMiners, err := pgx.CollectRows(results, pgx.RowToStructByName[GethMiner])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return gethMiners, nil
}

func GetTotalGethMinersCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM geth_miners
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
