package trade

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

func GetTrade(dbConnPgx utils.PgxIface, tradeID *int) (*Trade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		parent_trade_id,
		from_account_id,
		to_account_id,
		asset_id,
		source_id,
		uuid,
		transaction_id,
		order_id,
		trade_id,
		name,
		alternate_name,
		trade_type_id,
		trade_date,
		settle_date,
		transfer_date,
		from_quantity,
		to_quantity,
		price,
		total_amount,
		fees_amount,
		fees_asset_id,
		realized_return_amount,
		realized_return_asset_id,
		cost_basis_amount,
		cost_basis_trade_id,
		description,
		is_active,
		source_data,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM trades 
	WHERE id = $1
	`, *tradeID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	trade, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Trade])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &trade, nil
}
func GetTradeByTransactionIDFromSource(dbConnPgx utils.PgxIface, transactionIDFromSource string, fromAccountID, toAccountID *int) (*Trade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		parent_trade_id,
		from_account_id,
		to_account_id,
		asset_id,
		source_id,
		uuid,
		transaction_id,
		order_id,
		trade_id,
		name,
		alternate_name,
		trade_type_id,
		trade_date,
		settle_date,
		transfer_date,
		from_quantity,
		to_quantity,
		price,
		total_amount,
		fees_amount,
		fees_asset_id,
		realized_return_amount,
		realized_return_asset_id,
		cost_basis_amount,
		cost_basis_trade_id,
		description,
		is_active,
		source_data,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM trades 
	WHERE transaction_id = $1
	AND from_account_id = $2
	AND to_account_id = $3
	`, transactionIDFromSource, *fromAccountID, *toAccountID)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	trade, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Trade])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &trade, nil
}

func GetTradeByStartAndEndDates(dbConnPgx utils.PgxIface, startDate, endDate time.Time) ([]Trade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		parent_trade_id,
		from_account_id,
		to_account_id,
		asset_id,
		source_id,
		uuid,
		transaction_id,
		order_id,
		trade_id,
		name,
		alternate_name,
		trade_type_id,
		trade_date,
		settle_date,
		transfer_date,
		from_quantity,
		to_quantity,
		price,
		total_amount,
		fees_amount,
		fees_asset_id,
		realized_return_amount,
		realized_return_asset_id,
		cost_basis_amount,
		cost_basis_trade_id,
		description,
		is_active,
		source_data,
		created_by,
		created_at,
		updated_by,
		updated_at
		FROM trades
		WHERE trade_date BETWEEN $1 AND $2 
		AND is_active = true`,
		startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	trades, err := pgx.CollectRows(results, pgx.RowToStructByName[Trade])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return trades, nil
}

func GetTradeByFromAccountStartAndEndDates(dbConnPgx utils.PgxIface, fromAccountID *int, startDate, endDate time.Time) ([]Trade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		parent_trade_id,
		from_account_id,
		to_account_id,
		asset_id,
		source_id,
		uuid,
		transaction_id,
		order_id,
		trade_id,
		name,
		alternate_name,
		trade_type_id,
		trade_date,
		settle_date,
		transfer_date,
		from_quantity,
		to_quantity,
		price,
		total_amount,
		fees_amount,
		fees_asset_id,
		realized_return_amount,
		realized_return_asset_id,
		cost_basis_amount,
		cost_basis_trade_id,
		description,
		is_active,
		source_data,
		created_by,
		created_at,
		updated_by,
		updated_at
		FROM trades
		WHERE
		from_account_id = $1
		AND trade_date BETWEEN $2 AND $3
		AND is_active = true
		ORDER BY trade_date asc, asset_id`,
		*fromAccountID,
		startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	trades, err := pgx.CollectRows(results, pgx.RowToStructByName[Trade])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return trades, nil
}

func GetTradeByToAccountStartAndEndDates(dbConnPgx utils.PgxIface, toAccountID *int, startDate, endDate time.Time) ([]Trade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
		id,
		parent_trade_id,
		from_account_id,
		to_account_id,
		asset_id,
		source_id,
		uuid,
		transaction_id,
		order_id,
		trade_id,
		name,
		alternate_name,
		trade_type_id,
		trade_date,
		settle_date,
		transfer_date,
		from_quantity,
		to_quantity,
		price,
		total_amount,
		fees_amount,
		fees_asset_id,
		realized_return_amount,
		realized_return_asset_id,
		cost_basis_amount,
		cost_basis_trade_id,
		description,
		is_active,
		source_data,
		created_by,
		created_at,
		updated_by,
		updated_at
		FROM trades
		WHERE
		to_account_id = $1
		AND trade_date BETWEEN $2 AND $3
		AND is_active = true`,
		*toAccountID,
		startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	trades, err := pgx.CollectRows(results, pgx.RowToStructByName[Trade])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return trades, nil
}

func RemoveTrade(dbConnPgx utils.PgxIface, tradeID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveTrade DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM trades WHERE id = $1`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *tradeID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetTradeList(dbConnPgx utils.PgxIface) ([]Trade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	id,
	parent_trade_id,
	from_account_id,
	to_account_id,
	asset_id,
	source_id,
	uuid,
	transaction_id,
	order_id,
	trade_id,
	name,
	alternate_name,
	trade_type_id,
	trade_date,
	settle_date,
	transfer_date,
	from_quantity,
	to_quantity,
	price,
	total_amount,
	fees_amount,
	fees_asset_id,
	realized_return_amount,
	realized_return_asset_id,
	cost_basis_amount,
	cost_basis_trade_id,
	description,
	is_active,
	source_data,
	created_by,
	created_at,
	updated_by,
	updated_at
	FROM trades `)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	trades, err := pgx.CollectRows(results, pgx.RowToStructByName[Trade])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return trades, nil
}

func UpdateTrade(dbConnPgx utils.PgxIface, trade *Trade) error {
	// if the trade id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (trade.ID == nil || trade.SourceID == nil || *trade.SourceID == 0) || (trade.AssetID == nil || *trade.AssetID == 0) {
		return errors.New("trade has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateTrade DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE trades SET 
		parent_trade_id=$1,
		from_account_id=$2,
		to_account_id=$3,
		asset_id=$4,
		source_id=$5,
		transaction_id=$6,
		order_id=$7,
		trade_id=$8,
		name=$9,
		alternate_name=$10,
		trade_type_id=$11,
		trade_date=$12,
		settle_date=$13,
		transfer_date=$14,
		from_quantity=$15,
		to_quantity=$16,
		price=$17,
		total_amount=$18,
		fees_amount=$19,
		fees_asset_id=$20,
		realized_return_amount=$21,
		realized_return_asset_id=$22,
		cost_basis_amount=$23,
		cost_basis_trade_id=$24,
		description=$25,
		is_active=$26,
		source_data=$27,
		updated_by=$28, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$29`
	defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		trade.ParentTradeID,           //1
		trade.FromAccountID,           //2
		trade.ToAccountID,             //3
		trade.AssetID,                 //4
		trade.SourceID,                //5
		trade.TransactionIDFromSource, //6
		trade.OrderIDFromSource,       //7
		trade.TradeIDFromSource,       //8
		trade.Name,                    //9
		trade.AlternateName,           //10
		trade.TradeTypeID,             //11
		trade.TradeDate,               //12
		trade.SettleDate,              //13
		trade.TransferDate,            //14
		trade.FromQuantity,            //15
		trade.ToQuantity,              //16
		trade.Price,                   //17
		trade.TotalAmount,             //18
		trade.FeesAmount,              //19
		trade.FeesAssetID,             //20
		trade.RealizedReturnAmount,    //21
		trade.RealizedReturnAssetID,   //22
		trade.CostBasisAmount,         //23
		trade.CostBasisTradeID,        //24
		trade.Description,             //25
		trade.IsActive,                //26
		trade.SourceData,              //27
		trade.UpdatedBy,               //28
		trade.ID,                      //29
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertTrade(dbConnPgx utils.PgxIface, trade *Trade) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertTrade DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var tradeID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO trades  
	(
		parent_trade_id,
		from_account_id,
		to_account_id,
		asset_id,
		source_id,
		uuid,
		transaction_id,
		order_id,
		trade_id,
		name,
		alternate_name,
		trade_type_id,
		trade_date,
		settle_date,
		transfer_date,
		from_quantity,
		to_quantity,
		price,
		total_amount,
		fees_amount,
		fees_asset_id,
		realized_return_amount,
		realized_return_asset_id,
		cost_basis_amount,
		cost_basis_trade_id,
		description,
		is_active,
		source_data,
		created_by,
		created_at,
		updated_by,
		updated_at
		) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		uuid_generate_v4(),
		$6,
		$7,
		$8,
		$9,
		$10,
		$11,
		$12,
		$13,
		$14,
		$15,
		$16,
		$17,
		$18,
		$19,
		$20,
		$21,
		$22,
		$23,
		$24,
		$25,
		$26,
		$27,
		$28,
		current_timestamp at time zone 'UTC',
		$28,
		current_timestamp at time zone 'UTC'
		)
		RETURNING id`,
		trade.ParentTradeID,           //1
		trade.FromAccountID,           //2
		trade.ToAccountID,             //3
		trade.AssetID,                 //4
		trade.SourceID,                //5
		trade.TransactionIDFromSource, //6
		trade.OrderIDFromSource,       //7
		trade.TradeIDFromSource,       //8
		trade.Name,                    //9
		trade.AlternateName,           //10
		trade.TradeTypeID,             //11
		trade.TradeDate,               //12
		trade.SettleDate,              //13
		trade.TransferDate,            //14
		trade.FromQuantity,            //15
		trade.ToQuantity,              //16
		trade.Price,                   //17
		trade.TotalAmount,             //18
		trade.FeesAmount,              //19
		trade.FeesAssetID,             //20
		trade.RealizedReturnAmount,    //21
		trade.RealizedReturnAssetID,   //22
		trade.CostBasisAmount,         //23
		trade.CostBasisTradeID,        //24
		trade.Description,             //25
		trade.IsActive,                //26
		trade.SourceData,              //27
		trade.CreatedBy,               //28
	).Scan(&tradeID)

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
	return int(tradeID), nil
}

func InsertTrades(dbConnPgx utils.PgxIface, trades []Trade) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range trades {
		trade := trades[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(trade.UUID)
		row := []interface{}{
			trade.ParentTradeID,           //1
			trade.FromAccountID,           //2
			trade.ToAccountID,             //3
			trade.AssetID,                 //4
			trade.SourceID,                //5
			uuidString,                    //6
			trade.TransactionIDFromSource, //7
			trade.OrderIDFromSource,       //8
			trade.TradeIDFromSource,       //9
			trade.Name,                    //10
			trade.AlternateName,           //11
			trade.TradeTypeID,             //12
			trade.TradeDate,               //13
			trade.SettleDate,              //14
			trade.TransferDate,            //15
			trade.FromQuantity,            //16
			trade.ToQuantity,              //17
			trade.Price,                   //18
			trade.TotalAmount,             //19
			trade.FeesAmount,              //20
			trade.FeesAssetID,             //21
			trade.RealizedReturnAmount,    //22
			trade.RealizedReturnAssetID,   //23
			trade.CostBasisAmount,         //24
			trade.CostBasisTradeID,        //25
			trade.Description,             //26
			trade.IsActive,                //27
			trade.SourceData,              //28
			trade.CreatedBy,               //29
			&now,                          //30
			trade.CreatedBy,               //31
			&now,                          //32
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"trades"},
		[]string{
			"parent_trade_id",          //1
			"from_account_id",          //2
			"to_account_id",            //3
			"asset_id",                 //4
			"source_id",                //5
			"uuid",                     //6
			"transaction_id",           //7
			"order_id",                 //8
			"trade_id",                 //9
			"name",                     //10
			"alternate_name",           //11
			"trade_type_id",            //12
			"trade_date",               //13
			"settle_date",              //14
			"transfer_date",            //15
			"from_quantity",            //16
			"to_quantity",              //17
			"price",                    //18
			"total_amount",             //19
			"fees_amount",              //20
			"fees_asset_id",            //21
			"realized_return_amount",   //22
			"realized_return_asset_id", //23
			"cost_basis_amount",        //24
			"cost_basis_trade_id",      //25
			"description",              //26
			"is_active",                //27
			"source_data",              //28
			"created_by",               //29
			"created_at",               //30
			"updated_by",               //31
			"updated_at",               //32
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertTrades: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func GetMinAndMaxTradeDates(dbConnPgx utils.PgxIface) (*MinMaxTradeDates, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := dbConnPgx.QueryRow(ctx, `SELECT 
		MIN(trade_date) as min_date,
		MAX(trade_date) as max_date
	FROM trades 
	`)

	minMaxTradeDates := &MinMaxTradeDates{}
	err := row.Scan(
		&minMaxTradeDates.MinTradeDate,
		&minMaxTradeDates.MaxTradeDate,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return minMaxTradeDates, nil
}

func GetTradeListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]Trade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	sql := `SELECT 
		id,
		parent_trade_id,
		from_account_id,
		to_account_id,
		asset_id,
		source_id,
		uuid,
		transaction_id,
		order_id,
		trade_id,
		name,
		alternate_name,
		trade_type_id,
		trade_date,
		settle_date,
		transfer_date,
		from_quantity,
		to_quantity,
		price,
		total_amount,
		fees_amount,
		fees_asset_id,
		realized_return_amount,
		realized_return_asset_id,
		cost_basis_amount,
		cost_basis_trade_id,
		description,
		is_active,
		source_data,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM trades 
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
	trades, err := pgx.CollectRows(results, pgx.RowToStructByName[Trade])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return trades, nil
}

func GetTotalTradesCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM trades
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
