package gethlyletransactions

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

func GetGethTransactionInput(dbConnPgx utils.PgxIface, gethTransactionInputID *int) (*GethTransactionInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		name,
		alternate_name,
		function_name,
		method_id_str,
		num_of_parameters,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM geth_transaction_inputs
	WHERE id = $1
	`, *gethTransactionInputID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	gethTransactionInput, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethTransactionInput])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethTransactionInput, nil
}

func GetGethTransactionInputByFromToAddress(dbConnPgx utils.PgxIface, fromToAddressID *int) ([]GethTransactionInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT
		id,
		uuid,
		name ,
		alternate_name,
		function_name,
		method_id_str,
		num_of_parameters,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM geth_transaction_inputs
	WHERE
		(from_address_id =$1 OR 
			to_address_id = $1)
		`,
		*fromToAddressID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactionInputs, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransactionInput])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransactionInputs, nil
}

func RemoveGethTransactionInput(dbConnPgx utils.PgxIface, gethTransactionInputID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethTransactionInput DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_transaction_inputs WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *gethTransactionInputID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetGethTransactionInputList(dbConnPgx utils.PgxIface) ([]GethTransactionInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT
		id,
		uuid,
		chain_id,
		name,
		alternate_name,
		function_name,
		method_id_str,
		num_of_parameters,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM geth_transaction_inputs `)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactionInputs, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransactionInput])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransactionInputs, nil
}

func UpdateGethTransactionInput(dbConnPgx utils.PgxIface, gethTransactionInput *GethTransactionInput) error {
	// if the gethTransactionInput id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethTransactionInput.ID == nil {
		return errors.New("gethTransactionInput has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethTransactionInput DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_transaction_inputs SET 
		name=$1,
		alternate_name=$2,
		function_name=$3,
		method_id_str=$4,
		num_of_parameters=$5,
		description=$6,
		updated_by=$7,
		updated_at=current_timestamp at time zone 'UTC',
		WHERE id=$8`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		gethTransactionInput.Name,            //1
		gethTransactionInput.AlternateName,   //2
		gethTransactionInput.FunctionName,    //3
		gethTransactionInput.MethodIDStr,     //4
		gethTransactionInput.NumOfParameters, //5
		gethTransactionInput.Description,     //6
		gethTransactionInput.UpdatedBy,       //7
		gethTransactionInput.ID,              //8
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertGethTransactionInput(dbConnPgx utils.PgxIface, gethTransactionInput *GethTransactionInput) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertGethTransactionInput DbConn.Begin   %s", err.Error())
		return -1, "", err
	}
	var gethTransactionInputID int
	var gethTransactionInputUUID string
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO geth_transaction_inputs
	(
		uuid,
		name,
		alternate_name,
		function_name,
		method_id_str,
		num_of_parameters,
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
		current_timestamp at time zone 'UTC', 
		$7,
		current_timestamp at time zone 'UTC'
		)
		RETURNING id, uuid`,
		gethTransactionInput.Name,            //1
		gethTransactionInput.AlternateName,   //2
		gethTransactionInput.FunctionName,    //3
		gethTransactionInput.MethodIDStr,     //4
		gethTransactionInput.NumOfParameters, //5
		gethTransactionInput.Description,     //6
		gethTransactionInput.CreatedBy,       //7
	).Scan(&gethTransactionInputID, &gethTransactionInputUUID)
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
	return int(gethTransactionInputID), gethTransactionInputUUID, nil
}
func InsertGethTransactionInputs(dbConnPgx utils.PgxIface, gethTransactionInputs []GethTransactionInput) error {
	// need to supply uuid
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range gethTransactionInputs {
		gethTransactionInput := gethTransactionInputs[i]
		uuidString := pgtype.UUID{}
		uuidString.Set(gethTransactionInput.UUID)
		row := []interface{}{
			uuidString,                           //1
			gethTransactionInput.Name,            //2
			gethTransactionInput.AlternateName,   //3
			gethTransactionInput.FunctionName,    //4
			gethTransactionInput.MethodIDStr,     //5
			gethTransactionInput.NumOfParameters, //6
			gethTransactionInput.Description,     //7
			gethTransactionInput.CreatedBy,       //8
			now,                                  //9
			gethTransactionInput.CreatedBy,       //10
			now,                                  //11

		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_transaction_inputs"},
		[]string{
			"uuid",              //1
			"name",              //2
			"alternate_name",    //3
			"function_name",     //4
			"method_id_str",     //5
			"num_of_parameters", //6
			"description",       //7
			"created_by",        //8
			"created_at",        //9
			"updated_by",        //10
			"updated_at",        //11
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertGethTransactionInputs: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// for refinedev
func GetTransactionInputListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]GethTransactionInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		id,
		uuid,
		name,
		alternate_name,
		function_name,
		method_id_str,
		num_of_parameters,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM geth_transaction_inputs 
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
	gethTransactionInputs, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransactionInput])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransactionInputs, nil
}

func GetTotalTransactionInputsCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM geth_transaction_inputs
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

func GetGethTransactionInputByFromMinerID(dbConnPgx utils.PgxIface, minerID *int) ([]GethTransactionInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT
		gti.id,
		gti.uuid,
		gti.name ,
		gti.alternate_name,
		gti.function_name,
		gti.method_id_str,
		gti.num_of_parameters,
		gti.description,
		gti.created_by,
		gti.created_at,
		gti.updated_by,
		gti.updated_at
	FROM geth_transaction_inputs gti 
	JOIN geth_miners_transaction_inputs gmti
		ON gti.id = gmti.transaction_input_id
	WHERE
		gmti.miner_id = $1
		`,
		*minerID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactionInputs, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTransactionInput])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTransactionInputs, nil
}
