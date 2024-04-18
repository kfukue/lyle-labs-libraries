package gethlyletransactions

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/lib/pq"
)

func GetGethTransactionInput(gethTransactionInputID int) (*GethTransactionInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `SELECT
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
	`, gethTransactionInputID)

	gethTransactionInput := GethTransactionInput{}
	err := row.Scan(
		&gethTransactionInput.ID,
		&gethTransactionInput.UUID,
		&gethTransactionInput.Name,
		&gethTransactionInput.AlternateName,
		&gethTransactionInput.FunctionName,
		&gethTransactionInput.MethodIDStr,
		&gethTransactionInput.NumOfParameters,
		&gethTransactionInput.Description,
		&gethTransactionInput.CreatedBy,
		&gethTransactionInput.CreatedAt,
		&gethTransactionInput.UpdatedBy,
		&gethTransactionInput.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethTransactionInput, nil
}

func GetGethTransactionInputByFromToAddress(fromToAddressID *int) ([]GethTransactionInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
		AND (from_address_id =$1 OR 
			to_address_id = $1)
		`,
		fromToAddressID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactionInputs := make([]GethTransactionInput, 0)
	for results.Next() {
		var gethTransactionInput GethTransactionInput
		results.Scan(
			&gethTransactionInput.ID,
			&gethTransactionInput.UUID,
			&gethTransactionInput.Name,
			&gethTransactionInput.AlternateName,
			&gethTransactionInput.FunctionName,
			&gethTransactionInput.MethodIDStr,
			&gethTransactionInput.NumOfParameters,
			&gethTransactionInput.Description,
			&gethTransactionInput.CreatedBy,
			&gethTransactionInput.CreatedAt,
			&gethTransactionInput.UpdatedBy,
			&gethTransactionInput.UpdatedAt,
		)
		gethTransactionInputs = append(gethTransactionInputs, gethTransactionInput)
	}
	return gethTransactionInputs, nil
}

func GetGethTransactionInputByFromAddressAndBeforeBlockNumber(fromAddressID, blockNumber *int) ([]GethTransactionInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
	WHERE
		(from_address_id =$1 OR 
			to_address_id = $1)
		AND block_number <= $2
		`,
		fromAddressID, blockNumber,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactionInputs := make([]GethTransactionInput, 0)
	for results.Next() {
		var gethTransactionInput GethTransactionInput
		results.Scan(
			&gethTransactionInput.ID,
			&gethTransactionInput.UUID,
			&gethTransactionInput.Name,
			&gethTransactionInput.AlternateName,
			&gethTransactionInput.FunctionName,
			&gethTransactionInput.MethodIDStr,
			&gethTransactionInput.NumOfParameters,
			&gethTransactionInput.Description,
			&gethTransactionInput.CreatedBy,
			&gethTransactionInput.CreatedAt,
			&gethTransactionInput.UpdatedBy,
			&gethTransactionInput.UpdatedAt,
		)
		gethTransactionInputs = append(gethTransactionInputs, gethTransactionInput)
	}
	return gethTransactionInputs, nil
}

func GetGethTransactionInputsByTxnHash(txnHash string) ([]GethTransactionInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
		WHERE
		txn_hash = $1
		`,
		txnHash,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactionInputs := make([]GethTransactionInput, 0)
	for results.Next() {
		var gethTransactionInput GethTransactionInput
		results.Scan(
			&gethTransactionInput.ID,
			&gethTransactionInput.UUID,
			&gethTransactionInput.Name,
			&gethTransactionInput.AlternateName,
			&gethTransactionInput.FunctionName,
			&gethTransactionInput.MethodIDStr,
			&gethTransactionInput.NumOfParameters,
			&gethTransactionInput.Description,
			&gethTransactionInput.CreatedBy,
			&gethTransactionInput.CreatedAt,
			&gethTransactionInput.UpdatedBy,
			&gethTransactionInput.UpdatedAt,
		)
		gethTransactionInputs = append(gethTransactionInputs, gethTransactionInput)
	}
	return gethTransactionInputs, nil
}

func GetGethTransactionInputsByTxnHashes(txnHashes []string) ([]GethTransactionInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
		FROM geth_transaction_inputs
		WHERE
		txn_hash = ANY($1)
		`,
		pq.Array(txnHashes),
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactionInputs := make([]GethTransactionInput, 0)
	for results.Next() {
		var gethTransactionInput GethTransactionInput
		results.Scan(
			&gethTransactionInput.ID,
			&gethTransactionInput.UUID,
			&gethTransactionInput.Name,
			&gethTransactionInput.AlternateName,
			&gethTransactionInput.FunctionName,
			&gethTransactionInput.MethodIDStr,
			&gethTransactionInput.NumOfParameters,
			&gethTransactionInput.Description,
			&gethTransactionInput.CreatedBy,
			&gethTransactionInput.CreatedAt,
			&gethTransactionInput.UpdatedBy,
			&gethTransactionInput.UpdatedAt,
		)
		gethTransactionInputs = append(gethTransactionInputs, gethTransactionInput)
	}
	return gethTransactionInputs, nil
}

func RemoveGethTransactionInput(gethTransactionInputID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `DELETE FROM geth_transaction_inputs WHERE id = $1`, gethTransactionInputID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func RemoveGethTransactionInputsFromChainIDAndStartBlockNumber(chainID, startBlockNumber *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `DELETE FROM geth_transaction_inputs WHERE chain_id = $1 AND block_number >=  $2`, chainID, startBlockNumber)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func RemoveGethTransactionInputsFromChainID(chainID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `DELETE FROM geth_transaction_inputs WHERE chain_id = $1`, *chainID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetGethTransactionInputList() ([]GethTransactionInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
	gethTransactionInputs := make([]GethTransactionInput, 0)
	for results.Next() {
		var gethTransactionInput GethTransactionInput
		results.Scan(
			&gethTransactionInput.ID,
			&gethTransactionInput.UUID,
			&gethTransactionInput.Name,
			&gethTransactionInput.AlternateName,
			&gethTransactionInput.FunctionName,
			&gethTransactionInput.MethodIDStr,
			&gethTransactionInput.NumOfParameters,
			&gethTransactionInput.Description,
			&gethTransactionInput.CreatedBy,
			&gethTransactionInput.CreatedAt,
			&gethTransactionInput.UpdatedBy,
			&gethTransactionInput.UpdatedAt,
		)

		gethTransactionInputs = append(gethTransactionInputs, gethTransactionInput)
	}
	return gethTransactionInputs, nil
}

func UpdateGethTransactionInput(gethTransactionInput GethTransactionInput) error {
	// if the gethTransactionInput id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethTransactionInput.ID == nil {
		return errors.New("gethTransactionInput has invalid ID")
	}
	_, err := database.DbConnPgx.Exec(ctx, `UPDATE geth_transaction_inputs SET
		name=$1,
		alternate_name=$2,
		function_name=$3,
		method_id_str=$4,
		num_of_parameters=$5,
		description=$6,
		updated_by=$7,
		updated_at=current_timestamp at time zone 'UTC',
		WHERE id=$8`,
		gethTransactionInput.Name,            //1
		gethTransactionInput.AlternateName,   //2
		gethTransactionInput.FunctionName,    //3
		gethTransactionInput.MethodIDStr,     //4
		gethTransactionInput.NumOfParameters, //5
		gethTransactionInput.Description,     //6
		gethTransactionInput.UpdatedBy,       //7
		gethTransactionInput.ID,              //8
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertGethTransactionInput(gethTransactionInput GethTransactionInput) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var gethTransactionInputID int
	var gethTransactionInputUUID string
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO geth_transaction_inputs
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
		log.Println(err.Error())
		return 0, "", err
	}
	return int(gethTransactionInputID), gethTransactionInputUUID, nil
}
func InsertGethTransactionInputs(gethTransactionInputs []*GethTransactionInput) error {
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
	copyCount, err := database.DbConnPgx.CopyFrom(
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
	log.Println(fmt.Printf("geth_transaction_inputs copy count: %d", copyCount))
	if err != nil {
		log.Fatal(err)
		// handle error that occurred while using *pgx.Conn
	}
	return nil
}

func UpdateGethTransactionInputAddresses() error {
	// update address ids from existing addresses in geth_addresses
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Exec(ctx, `
		UPDATE geth_transaction_inputs as gt SET
			from_address_id = ga.id from geth_addresses as ga
			WHERE LOWER(gt.from_address) = LOWER(ga.address_str)
			AND gt.from_address_id IS NULL
			`,
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	_, err = database.DbConnPgx.Exec(ctx, `
			UPDATE geth_transaction_inputs as gt SET
			to_address_id = ga.id from geth_addresses as ga
			WHERE LOWER(gt.to_address) = LOWER(ga.address_str)
			AND gt.to_address_id IS NULL
			`,
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetNullAddressStrsFromTransactionInputs() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
		WITH sender_table as(
			SELECT DISTINCT LOWER(gt.from_address) as address  
			FROM geth_transaction_inputs gt
				LEFT JOIN geth_addresses as ga
				ON LOWER(gt.from_address) = LOWER(ga.address_str)
			WHERE gt.from_address_id IS NULL
			AND gt.base_asset_id = $1
			AND ga.id IS NULL
		),
		to_table as(
			SELECT DISTINCT LOWER(gt.to_address) as address 
			FROM geth_transaction_inputs gt
				LEFT JOIN geth_addresses as ga
				ON LOWER(gt.to_address) = LOWER(ga.address_str)
			WHERE gt.to_address_id IS NULL
			AND gt.base_asset_id = $1
			AND ga.id IS NULL
		)
		SELECT * FROM sender_table
		UNION
		SELECT * FROM to_table
		`,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethNullAddressStrs := make([]string, 0)
	for results.Next() {
		var gethNullAddressStr string
		results.Scan(
			&gethNullAddressStr,
		)
		gethNullAddressStrs = append(gethNullAddressStrs, gethNullAddressStr)
	}
	return gethNullAddressStrs, nil
}

// for refinedev
func GetTransactionInputListByPagination(_start, _end *int, _order, _sort string, _filters []string) ([]GethTransactionInput, error) {
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

	results, err := database.DbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTransactionInputs := make([]GethTransactionInput, 0)
	for results.Next() {
		var gethTransactionInput GethTransactionInput
		results.Scan(
			&gethTransactionInput.ID,
			&gethTransactionInput.UUID,
			&gethTransactionInput.Name,
			&gethTransactionInput.AlternateName,
			&gethTransactionInput.FunctionName,
			&gethTransactionInput.MethodIDStr,
			&gethTransactionInput.NumOfParameters,
			&gethTransactionInput.Description,
			&gethTransactionInput.CreatedBy,
			&gethTransactionInput.CreatedAt,
			&gethTransactionInput.UpdatedBy,
			&gethTransactionInput.UpdatedAt,
		)

		gethTransactionInputs = append(gethTransactionInputs, gethTransactionInput)
	}
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()

	return gethTransactionInputs, nil
}

func GetTotalTransactionInputsCount() (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM geth_transaction_inputs
	`)
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
