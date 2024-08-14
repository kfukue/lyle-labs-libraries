package gethlyletrades

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

func GetAllGethTradeTaxTransfersByTradeID(dbConnPgx utils.PgxIface, gethTradeID *int) ([]GethTradeTaxTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT 
		geth_trade_transfers.geth_trade_id,
		geth_trade_transfers.geth_transfer_id,
		geth_trade_transfers.tax_id,
		geth_trade_transfers.uuid, 
		geth_trade_transfers.name, 
		geth_trade_transfers.alternate_name, 
		geth_trade_transfers.description,
		geth_trade_transfers.created_by, 
		geth_trade_transfers.created_at, 
		geth_trade_transfers.updated_by, 
		geth_trade_transfers.updated_at 
	FROM geth_trades 
	LEFT JOIN geth_trade_transfers
	ON geth_trade_transfers.geth_trade_id = geth_trades.id 
	WHERE 
	geth_trades.id = $1
	`, *gethTradeID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethTradeTaxTransfers, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTradeTaxTransfer])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTradeTaxTransfers, nil
}
func GetGethTradeTaxTransfer(dbConnPgx utils.PgxIface, gethTradeID, gethGethTransferID *int) (*GethTradeTaxTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `
	SELECT 
		geth_trade_id,
		geth_transfer_id,
		tax_id,
		uuid, 
		name, 
		alternate_name, 
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM geth_trade_transfers 
	WHERE 
	geth_trade_id = $1
	AND geth_transfer_id = $2
	`, *gethTradeID, *gethGethTransferID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	gethTradeTaxTransfer, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethTradeTaxTransfer])

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethTradeTaxTransfer, nil
}

func RemoveGethTradeTaxTransfer(dbConnPgx utils.PgxIface, gethTradeID, gethGethTransferID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethTradeTaxTransfer DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_trade_transfers WHERE geth_trade_id =$1 AND geth_transfer_id = $2`

	if _, err := dbConnPgx.Exec(ctx, sql, *gethTradeID, *gethGethTransferID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetGethTradeTaxTransferList(dbConnPgx utils.PgxIface, gethTradeIds []int, swapIds []int) ([]GethTradeTaxTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT 
		geth_trade_id,
		geth_transfer_id,
		tax_id,
		uuid, 
		name, 
		alternate_name, 
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM geth_trade_transfers`
	if len(gethTradeIds) > 0 || len(swapIds) > 0 {
		additionalQuery := ` WHERE`
		if len(gethTradeIds) > 0 {
			gethTradeStrIds := utils.SplitToString(gethTradeIds, ",")
			additionalQuery += fmt.Sprintf(`geth_trade_id IN (%s)`, gethTradeStrIds)
		}
		if len(swapIds) > 0 {
			if len(gethTradeIds) > 0 {
				additionalQuery += `AND `
			}
			swapStrIds := utils.SplitToString(swapIds, ",")
			additionalQuery += fmt.Sprintf(`geth_transfer_id IN (%s)`, swapStrIds)
		}
		sql += additionalQuery
	}
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethTradeTaxTransfers, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTradeTaxTransfer])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTradeTaxTransfers, nil
}

func UpdateGethTradeTaxTransfer(dbConnPgx utils.PgxIface, gethTradeTaxTransfer *GethTradeTaxTransfer) error {
	// if the gethTradeTaxTransfer id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (gethTradeTaxTransfer.GethTransferID == nil || *gethTradeTaxTransfer.GethTransferID == 0) || (gethTradeTaxTransfer.GethTradeID == nil || *gethTradeTaxTransfer.GethTradeID == 0) {
		return errors.New("gethTradeTaxTransfer has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethTradeTaxTransfer DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_trade_transfers SET 
		tax_id 				=$1,	
		name				=$2,  
		alternate_name		=$3, 
		description			=$4,
		updated_by			=$5, 
		updated_at			=current_timestamp at time zone 'UTC'
		WHERE 
			geth_trade_id=$6 AND
			geth_transfer_id=$7
		`
	if _, err := dbConnPgx.Exec(ctx, sql,
		gethTradeTaxTransfer.TaxID,          //1
		gethTradeTaxTransfer.Name,           //2
		gethTradeTaxTransfer.AlternateName,  //3
		gethTradeTaxTransfer.Description,    //4
		gethTradeTaxTransfer.UpdatedBy,      //5
		gethTradeTaxTransfer.GethTradeID,    //6
		gethTradeTaxTransfer.GethTransferID, //7
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertGethTradeTaxTransfer(dbConnPgx utils.PgxIface, gethTradeTaxTransfer *GethTradeTaxTransfer) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertGethTradeTaxTransfer DbConn.Begin   %s", err.Error())
		return -1, -1, err
	}
	var GethTransferID int
	var GethTradeID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO geth_trade_transfers  
	(
		geth_trade_id,
		geth_transfer_id,
		tax_id, 
		uuid,	 
		name, 
		alternate_name,  
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
		) VALUES ($1,
			$2,
			$3,
			uuid_generate_v4(),
			$4,
			$5,
			$6,
			$7,
			current_timestamp at time zone 'UTC',
			$7,
			current_timestamp at time zone 'UTC'
		)
		RETURNING geth_trade_id, geth_transfer_id`,
		gethTradeTaxTransfer.GethTradeID,    //1
		gethTradeTaxTransfer.GethTransferID, //2
		gethTradeTaxTransfer.TaxID,          //3
		gethTradeTaxTransfer.Name,           //4
		gethTradeTaxTransfer.AlternateName,  //5
		gethTradeTaxTransfer.Description,    //6
		gethTradeTaxTransfer.CreatedBy,      //7
	).Scan(&GethTradeID, &GethTransferID)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, -1, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		log.Println(err.Error())
		return -1, -1, err
	}
	return int(GethTradeID), int(GethTransferID), nil
}

func InsertGethTradeTaxTransfers(dbConnPgx utils.PgxIface, gethTradeTaxTransfers []GethTradeTaxTransfer) error {
	// need to supply uuid
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range gethTradeTaxTransfers {
		gethTradeTaxTransfer := gethTradeTaxTransfers[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(gethTradeTaxTransfer.UUID)
		row := []interface{}{
			gethTradeTaxTransfer.GethTradeID,    //1
			gethTradeTaxTransfer.GethTransferID, //2
			gethTradeTaxTransfer.TaxID,          //3
			uuidString,                          //4
			gethTradeTaxTransfer.Name,           //5
			gethTradeTaxTransfer.AlternateName,  //6
			gethTradeTaxTransfer.Description,    //7
			gethTradeTaxTransfer.CreatedBy,      //8
			&now,                                //9
			gethTradeTaxTransfer.CreatedBy,      //10
			&now,                                //11
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_trade_transfers"},
		[]string{
			"geth_trade_id",    //1
			"geth_transfer_id", //2
			"tax_id",           //3
			"uuid",             //4
			"name",             //5
			"alternate_name",   //6
			"description",      //7
			"created_by",       //8
			"created_at",       //9
			"updated_by",       //10
			"updated_at",       //11
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertGethTradeTaxTransfers: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// for refinedev
func GetGethTradeTaxTransferListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]GethTradeTaxTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	sql := `SELECT 
	geth_trade_id,
		geth_trade_id,
		geth_transfer_id,
		tax_id,
		uuid, 
		name, 
		alternate_name, 
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM geth_trade_transfers 
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

	gethTradeTaxTransfers, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTradeTaxTransfer])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTradeTaxTransfers, nil
}

func GetTotalGethTradeTaxTransferCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM geth_trade_transfers
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
