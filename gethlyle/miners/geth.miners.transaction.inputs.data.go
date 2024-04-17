package gethlyleminers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
)

func GetAllGethMinerTransactionInputsByMinerID(minerID *int) ([]GethMinerTransactionInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
	SELECT 
		miner_id,
		transaction_input_id,
		uuid,
		name,
		alternate_name,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
	FROM geth_miners_transaction_inputs 
	WHERE 
	miner_id = $1
	`, minerID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethMinerTransactionInputs := make([]GethMinerTransactionInput, 0)
	for results.Next() {
		var gethMinerTransactionInput GethMinerTransactionInput
		results.Scan(
			&gethMinerTransactionInput.MinerID,
			&gethMinerTransactionInput.TransactionInputID,
			&gethMinerTransactionInput.UUID,
			&gethMinerTransactionInput.Name,
			&gethMinerTransactionInput.AlternateName,
			&gethMinerTransactionInput.Description,
			&gethMinerTransactionInput.CreatedBy,
			&gethMinerTransactionInput.CreatedAt,
			&gethMinerTransactionInput.UpdatedBy,
			&gethMinerTransactionInput.UpdatedAt,
		)

		gethMinerTransactionInputs = append(gethMinerTransactionInputs, gethMinerTransactionInput)
	}
	return gethMinerTransactionInputs, nil
}

func GetAllGethMinerTransactionInputsByTransactionInputID(transactionInputID *int) ([]GethMinerTransactionInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
	SELECT 
		miner_id,
		transaction_input_id,
		uuid,
		name,
		alternate_name,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
	FROM geth_miners_transaction_inputs 
	WHERE 
	transaction_input_id = $1
	`, transactionInputID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethMinerTransactionInputs := make([]GethMinerTransactionInput, 0)
	for results.Next() {
		var gethMinerTransactionInput GethMinerTransactionInput
		results.Scan(
			&gethMinerTransactionInput.MinerID,
			&gethMinerTransactionInput.TransactionInputID,
			&gethMinerTransactionInput.UUID,
			&gethMinerTransactionInput.Name,
			&gethMinerTransactionInput.AlternateName,
			&gethMinerTransactionInput.Description,
			&gethMinerTransactionInput.CreatedBy,
			&gethMinerTransactionInput.CreatedAt,
			&gethMinerTransactionInput.UpdatedBy,
			&gethMinerTransactionInput.UpdatedAt,
		)

		gethMinerTransactionInputs = append(gethMinerTransactionInputs, gethMinerTransactionInput)
	}
	return gethMinerTransactionInputs, nil
}

func GetMinerTransactionInput(minerID, transactionInputID *int) (*GethMinerTransactionInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `
	SELECT 
		miner_id,
		transaction_input_id,
		uuid,
		name,
		alternate_name,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
	FROM geth_miners_transaction_inputs 
	WHERE miner_id = $1
	AND transaction_input_id = $2
	`, minerID, transactionInputID)

	gethMinerTransactionInput := &GethMinerTransactionInput{}
	err := row.Scan(
		&gethMinerTransactionInput.MinerID,
		&gethMinerTransactionInput.TransactionInputID,
		&gethMinerTransactionInput.UUID,
		&gethMinerTransactionInput.Name,
		&gethMinerTransactionInput.AlternateName,
		&gethMinerTransactionInput.Description,
		&gethMinerTransactionInput.CreatedBy,
		&gethMinerTransactionInput.CreatedAt,
		&gethMinerTransactionInput.UpdatedBy,
		&gethMinerTransactionInput.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethMinerTransactionInput, nil
}

func RemoveMinerTransactionInput(minerID, transactionInputID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM geth_miners_transaction_inputs WHERE 
	miner_id = $1 AND transaction_input_id =$2`, minerID, transactionInputID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetMinerTransactionInputList(minerIDs, transactionInputIDs []int) ([]GethMinerTransactionInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT 
		miner_id,
		transaction_input_id,
		uuid,
		name,
		alternate_name,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
	FROM geth_miners_transaction_inputs `
	if len(minerIDs) > 0 || len(transactionInputIDs) > 0 {
		additionalQuery := ` WHERE`
		if len(minerIDs) > 0 {
			strIds := utils.SplitToString(minerIDs, ",")
			additionalQuery += fmt.Sprintf(`miner_id IN (%s)`, strIds)
		}
		if len(transactionInputIDs) > 0 {
			if len(minerIDs) > 0 {
				additionalQuery += `AND `
			}
			strIds := utils.SplitToString(transactionInputIDs, ",")
			additionalQuery += fmt.Sprintf(`transaction_input_id IN (%s)`, strIds)
		}
		sql += additionalQuery
	}
	results, err := database.DbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	minerTransactionInputs := make([]GethMinerTransactionInput, 0)
	for results.Next() {
		var minerTransactionInput GethMinerTransactionInput
		results.Scan(
			&minerTransactionInput.MinerID,
			&minerTransactionInput.TransactionInputID,
			&minerTransactionInput.UUID,
			&minerTransactionInput.Name,
			&minerTransactionInput.AlternateName,
			&minerTransactionInput.Description,
			&minerTransactionInput.CreatedBy,
			&minerTransactionInput.CreatedAt,
			&minerTransactionInput.UpdatedBy,
			&minerTransactionInput.UpdatedAt,
		)

		minerTransactionInputs = append(minerTransactionInputs, minerTransactionInput)
	}
	return minerTransactionInputs, nil
}

func UpdateMinerTransactionInput(minerTransactionInput GethMinerTransactionInput) error {
	// if the minerTransactionInput id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (minerTransactionInput.MinerID == nil || *minerTransactionInput.MinerID == 0) || (minerTransactionInput.TransactionInputID == nil || *minerTransactionInput.TransactionInputID == 0) {
		return errors.New("minerTransactionInput has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE geth_miners_transaction_inputs SET 
		name=$1,  
		alternate_name=$2, 
		description=$3,
		updated_by=$4, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE miner_id=$5 AND transaction_input_id=$6`,
		minerTransactionInput.Name,               //1
		minerTransactionInput.AlternateName,      //2
		minerTransactionInput.Description,        //3
		minerTransactionInput.UpdatedBy,          //4
		minerTransactionInput.MinerID,            //5
		minerTransactionInput.TransactionInputID, //6
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertMinerTransactionInput(minerTransactionInput GethMinerTransactionInput) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var MinerID int
	var TransactionInputID int
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO geth_miners_transaction_inputs  
	(
		miner_id,
		transaction_input_id,
		uuid,	 
		name, 
		alternate_name,  
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		) VALUES (
			$1,
			$2,
			uuid_generate_v4(),
			$3,
			$4,
			$5,
			$6,
			current_timestamp at time zone 'UTC',
			$6,
			current_timestamp at time zone 'UTC'
		)
		RETURNING miner_id, transaction_input_id`,
		minerTransactionInput.MinerID,            //1
		minerTransactionInput.TransactionInputID, //2
		minerTransactionInput.Name,               //3
		minerTransactionInput.AlternateName,      //4
		minerTransactionInput.Description,        //5
		minerTransactionInput.CreatedBy,          //6
	).Scan(&MinerID, &TransactionInputID)
	if err != nil {
		log.Println(err.Error())
		return 0, 0, err
	}
	return int(MinerID), int(TransactionInputID), nil
}

func InsertGethMinersTransactionInputs(gethMinersTransactionInputs []*GethMinerTransactionInput) error {
	// need to supply uuid
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range gethMinersTransactionInputs {
		gethTransactionInput := gethMinersTransactionInputs[i]
		uuidString := pgtype.UUID{}
		uuidString.Set(gethTransactionInput.UUID)
		row := []interface{}{
			gethTransactionInput.MinerID,            //1
			gethTransactionInput.TransactionInputID, //2
			uuidString,                              //3
			gethTransactionInput.Name,               //4
			gethTransactionInput.AlternateName,      //5
			gethTransactionInput.Description,        //6
			gethTransactionInput.CreatedBy,          //7
			now,                                     //8
			gethTransactionInput.CreatedBy,          //9
			now,                                     //10

		}
		rows = append(rows, row)
	}
	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_miners_transaction_inputs"},
		[]string{
			"miner_id",             //1
			"transaction_input_id", //2
			"uuid",                 //3
			"name",                 //4
			"alternate_name",       //5
			"description",          //6
			"created_by",           //7
			"created_at",           //8
			"updated_by",           //9
			"updated_at",           //10
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("geth_miners_transaction_inputs copy count: %d", copyCount))
	if err != nil {
		log.Fatal(err)
		// handle error that occurred while using *pgx.Conn
	}
	return nil
}

// for refinedev
func GetMinerTransactionInputListByPagination(_start, _end *int, _order, _sort string, _filters []string) ([]GethMinerTransactionInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	sql := `SELECT 
	miner_id,
	transaction_input_id,
	uuid, 
	name, 
	alternate_name, 
	description,
	created_by, 
	created_at, 
	updated_by, 
	updated_at 
	FROM geth_miners_transaction_inputs 
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
	minerTransactionInputs := make([]GethMinerTransactionInput, 0)
	for results.Next() {
		var minerTransactionInput GethMinerTransactionInput
		results.Scan(
			&minerTransactionInput.MinerID,
			&minerTransactionInput.TransactionInputID,
			&minerTransactionInput.UUID,
			&minerTransactionInput.Name,
			&minerTransactionInput.AlternateName,
			&minerTransactionInput.Description,
			&minerTransactionInput.CreatedBy,
			&minerTransactionInput.CreatedAt,
			&minerTransactionInput.UpdatedBy,
			&minerTransactionInput.UpdatedAt,
		)

		minerTransactionInputs = append(minerTransactionInputs, minerTransactionInput)
	}
	return minerTransactionInputs, nil
}

func GetTotalMinerTransactionInputCount() (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := database.DbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM geth_miners_transaction_inputs
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
