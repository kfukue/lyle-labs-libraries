package gethlyletrades

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

func GetAllGethTradeTaxTransfersByTradeID(gethTradeID int) ([]GethTradeTaxTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
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
	`, gethTradeID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTradeTaxTransfers := make([]GethTradeTaxTransfer, 0)
	for results.Next() {
		var gethTradeTaxTransfer GethTradeTaxTransfer
		results.Scan(
			&gethTradeTaxTransfer.GethTradeID,
			&gethTradeTaxTransfer.GethTransferID,
			&gethTradeTaxTransfer.TaxID,
			&gethTradeTaxTransfer.UUID,
			&gethTradeTaxTransfer.Name,
			&gethTradeTaxTransfer.AlternateName,
			&gethTradeTaxTransfer.Description,
			&gethTradeTaxTransfer.CreatedBy,
			&gethTradeTaxTransfer.CreatedAt,
			&gethTradeTaxTransfer.UpdatedBy,
			&gethTradeTaxTransfer.UpdatedAt,
		)

		gethTradeTaxTransfers = append(gethTradeTaxTransfers, gethTradeTaxTransfer)
	}
	return gethTradeTaxTransfers, nil
}
func GetGethTradeTaxTransfer(gethTradeID int, gethGethTransferID int) (*GethTradeTaxTransfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `
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
	`, gethTradeID, gethGethTransferID)

	gethTradeTaxTransfer := &GethTradeTaxTransfer{}
	err := row.Scan(
		&gethTradeTaxTransfer.GethTradeID,
		&gethTradeTaxTransfer.GethTransferID,
		&gethTradeTaxTransfer.TaxID,
		&gethTradeTaxTransfer.UUID,
		&gethTradeTaxTransfer.Name,
		&gethTradeTaxTransfer.AlternateName,
		&gethTradeTaxTransfer.Description,
		&gethTradeTaxTransfer.CreatedBy,
		&gethTradeTaxTransfer.CreatedAt,
		&gethTradeTaxTransfer.UpdatedBy,
		&gethTradeTaxTransfer.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTradeTaxTransfer, nil
}

func RemoveGethTradeTaxTransfer(gethTradeID int, gethGethTransferID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM geth_trade_transfers WHERE 
	geth_trade_id =$1 AND geth_transfer_id = $2`, gethTradeID, gethGethTransferID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetGethTradeTaxTransferList(gethTradeIds []int, swapIds []int) ([]GethTradeTaxTransfer, error) {
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
	results, err := database.DbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTradeTaxTransfers := make([]GethTradeTaxTransfer, 0)
	for results.Next() {
		var gethTradeTaxTransfer GethTradeTaxTransfer
		results.Scan(
			&gethTradeTaxTransfer.GethTradeID,
			&gethTradeTaxTransfer.GethTransferID,
			&gethTradeTaxTransfer.TaxID,
			&gethTradeTaxTransfer.UUID,
			&gethTradeTaxTransfer.Name,
			&gethTradeTaxTransfer.AlternateName,
			&gethTradeTaxTransfer.Description,
			&gethTradeTaxTransfer.CreatedBy,
			&gethTradeTaxTransfer.CreatedAt,
			&gethTradeTaxTransfer.UpdatedBy,
			&gethTradeTaxTransfer.UpdatedAt,
		)

		gethTradeTaxTransfers = append(gethTradeTaxTransfers, gethTradeTaxTransfer)
	}
	return gethTradeTaxTransfers, nil
}

func UpdateGethTradeTaxTransfer(gethTradeTaxTransfer GethTradeTaxTransfer) error {
	// if the gethTradeTaxTransfer id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (gethTradeTaxTransfer.GethTransferID == nil || *gethTradeTaxTransfer.GethTransferID == 0) || (gethTradeTaxTransfer.GethTradeID == nil || *gethTradeTaxTransfer.GethTradeID == 0) {
		return errors.New("gethTradeTaxTransfer has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE geth_trade_transfers SET 
		tax_id 				=$1,	
		name				=$2,  
		alternate_name		=$3, 
		description			=$4,
		updated_by			=$5, 
		updated_at			=current_timestamp at time zone 'UTC'
		WHERE 
			geth_trade_id=$6 AND
			geth_transfer_id=$7
		`,
		gethTradeTaxTransfer.TaxID,          //1
		gethTradeTaxTransfer.Name,           //2
		gethTradeTaxTransfer.AlternateName,  //3
		gethTradeTaxTransfer.Description,    //4
		gethTradeTaxTransfer.UpdatedBy,      //5
		gethTradeTaxTransfer.GethTradeID,    //6
		gethTradeTaxTransfer.GethTransferID, //7
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertGethTradeTaxTransfer(gethTradeTaxTransfer GethTradeTaxTransfer) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var GethTransferID int
	var GethTradeID int
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO geth_trade_transfers  
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
		log.Println(err.Error())
		return 0, 0, err
	}
	return int(GethTradeID), int(GethTransferID), nil
}

func InsertGethTradeTaxTransfers(gethTradeTaxTransfers []*GethTradeTaxTransfer) error {
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
	copyCount, err := database.DbConnPgx.CopyFrom(
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
	log.Println(fmt.Printf("copy count: %d", copyCount))
	if err != nil {
		log.Fatal(err)
		// handle error that occurred while using *pgx.Conn
	}
	return nil
}
