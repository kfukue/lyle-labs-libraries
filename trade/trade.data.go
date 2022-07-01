package trade

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
)

func GetTrade(tradeID int) (*Trade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	source_data,
	is_active,
	created_by,
	created_at,
	updated_by,
	updated_at
	FROM trades 
	WHERE id = $1
	`, tradeID)

	trade := &Trade{}
	err := row.Scan(
		&trade.ID,
		&trade.ParentTradeID,
		&trade.FromAccountID,
		&trade.ToAccountID,
		&trade.AssetID,
		&trade.SourceID,
		&trade.UUID,
		&trade.TransactionIDFromSource,
		&trade.OrderIDFromSource,
		&trade.TradeIDFromSource,
		&trade.Name,
		&trade.AlternateName,
		&trade.TradeTypeID,
		&trade.TradeDate,
		&trade.SettleDate,
		&trade.TransferDate,
		&trade.FromQuantity,
		&trade.ToQuantity,
		&trade.Price,
		&trade.TotalAmount,
		&trade.FeesAmount,
		&trade.FeesAssetID,
		&trade.RealizedReturnAmount,
		&trade.RealizedReturnAssetID,
		&trade.CostBasisAmount,
		&trade.CostBasisTradeID,
		&trade.Description,
		&trade.SourceData,
		&trade.IsActive,
		&trade.CreatedBy,
		&trade.CreatedAt,
		&trade.UpdatedBy,
		&trade.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return trade, nil
}
func GetTradeByTransactionIDFromSource(transactionIDFromSource string, fromAccountID *int, toAccountID *int) (*Trade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	source_data,
	is_active,
	created_by,
	created_at,
	updated_by,
	updated_at
	FROM trades 
	WHERE transaction_id = $1
	AND from_account_id = $2
	AND to_account_id = $3
	`, transactionIDFromSource, fromAccountID, toAccountID)

	trade := &Trade{}
	err := row.Scan(
		&trade.ID,
		&trade.ParentTradeID,
		&trade.FromAccountID,
		&trade.ToAccountID,
		&trade.AssetID,
		&trade.SourceID,
		&trade.UUID,
		&trade.TransactionIDFromSource,
		&trade.OrderIDFromSource,
		&trade.TradeIDFromSource,
		&trade.Name,
		&trade.AlternateName,
		&trade.TradeTypeID,
		&trade.TradeDate,
		&trade.SettleDate,
		&trade.TransferDate,
		&trade.FromQuantity,
		&trade.ToQuantity,
		&trade.Price,
		&trade.TotalAmount,
		&trade.FeesAmount,
		&trade.FeesAssetID,
		&trade.RealizedReturnAmount,
		&trade.RealizedReturnAssetID,
		&trade.CostBasisAmount,
		&trade.CostBasisTradeID,
		&trade.Description,
		&trade.SourceData,
		&trade.IsActive,
		&trade.CreatedBy,
		&trade.CreatedAt,
		&trade.UpdatedBy,
		&trade.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return trade, nil
}
func GetTradeByTradeID(tradeID int) (*Trade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	source_data,
	is_active,
	created_by,
	created_at,
	updated_by,
	updated_at
	FROM trades 
	WHERE id = $1	
	`, tradeID)

	trade := &Trade{}
	err := row.Scan(
		&trade.ID,
		&trade.ParentTradeID,
		&trade.FromAccountID,
		&trade.ToAccountID,
		&trade.AssetID,
		&trade.SourceID,
		&trade.UUID,
		&trade.TransactionIDFromSource,
		&trade.OrderIDFromSource,
		&trade.TradeIDFromSource,
		&trade.Name,
		&trade.AlternateName,
		&trade.TradeTypeID,
		&trade.TradeDate,
		&trade.SettleDate,
		&trade.TransferDate,
		&trade.FromQuantity,
		&trade.ToQuantity,
		&trade.Price,
		&trade.TotalAmount,
		&trade.FeesAmount,
		&trade.FeesAssetID,
		&trade.RealizedReturnAmount,
		&trade.RealizedReturnAssetID,
		&trade.CostBasisAmount,
		&trade.CostBasisTradeID,
		&trade.Description,
		&trade.SourceData,
		&trade.IsActive,
		&trade.CreatedBy,
		&trade.CreatedAt,
		&trade.UpdatedBy,
		&trade.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return trade, nil
}

func GetTradeByStartAndEndDates(startDate, endDate time.Time) ([]Trade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
		source_data,
		is_active,
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
	trades := make([]Trade, 0)
	for results.Next() {
		var trade Trade
		results.Scan(
			&trade.ID,
			&trade.ParentTradeID,
			&trade.FromAccountID,
			&trade.ToAccountID,
			&trade.AssetID,
			&trade.SourceID,
			&trade.UUID,
			&trade.TransactionIDFromSource,
			&trade.OrderIDFromSource,
			&trade.TradeIDFromSource,
			&trade.Name,
			&trade.AlternateName,
			&trade.TradeTypeID,
			&trade.TradeDate,
			&trade.SettleDate,
			&trade.TransferDate,
			&trade.FromQuantity,
			&trade.ToQuantity,
			&trade.Price,
			&trade.TotalAmount,
			&trade.FeesAmount,
			&trade.FeesAssetID,
			&trade.RealizedReturnAmount,
			&trade.RealizedReturnAssetID,
			&trade.CostBasisAmount,
			&trade.CostBasisTradeID,
			&trade.Description,
			&trade.SourceData,
			&trade.IsActive,
			&trade.CreatedBy,
			&trade.CreatedAt,
			&trade.UpdatedBy,
			&trade.UpdatedAt,
		)

		trades = append(trades, trade)
	}
	return trades, nil
}

func GetTradeByFromAccountStartAndEndDates(fromAccountID *int, startDate, endDate time.Time) ([]Trade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
		source_data,
		is_active,
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
		fromAccountID,
		startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	trades := make([]Trade, 0)
	for results.Next() {
		var trade Trade
		results.Scan(
			&trade.ID,
			&trade.ParentTradeID,
			&trade.FromAccountID,
			&trade.ToAccountID,
			&trade.AssetID,
			&trade.SourceID,
			&trade.UUID,
			&trade.TransactionIDFromSource,
			&trade.OrderIDFromSource,
			&trade.TradeIDFromSource,
			&trade.Name,
			&trade.AlternateName,
			&trade.TradeTypeID,
			&trade.TradeDate,
			&trade.SettleDate,
			&trade.TransferDate,
			&trade.FromQuantity,
			&trade.ToQuantity,
			&trade.Price,
			&trade.TotalAmount,
			&trade.FeesAmount,
			&trade.FeesAssetID,
			&trade.RealizedReturnAmount,
			&trade.RealizedReturnAssetID,
			&trade.CostBasisAmount,
			&trade.CostBasisTradeID,
			&trade.Description,
			&trade.SourceData,
			&trade.IsActive,
			&trade.CreatedBy,
			&trade.CreatedAt,
			&trade.UpdatedBy,
			&trade.UpdatedAt,
		)
		trades = append(trades, trade)
	}
	return trades, nil
}

func GetTradeByToAccountStartAndEndDates(toAccountID *int, startDate, endDate time.Time) ([]Trade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
		source_data,
		is_active,
		created_by,
		created_at,
		updated_by,
		updated_at
		FROM trades
		WHERE
		to_account_id = $1
		AND trade_date BETWEEN $2 AND $3
		AND is_active = true`,
		toAccountID,
		startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	trades := make([]Trade, 0)
	for results.Next() {
		var trade Trade
		results.Scan(
			&trade.ID,
			&trade.ParentTradeID,
			&trade.FromAccountID,
			&trade.ToAccountID,
			&trade.AssetID,
			&trade.SourceID,
			&trade.UUID,
			&trade.TransactionIDFromSource,
			&trade.OrderIDFromSource,
			&trade.TradeIDFromSource,
			&trade.Name,
			&trade.AlternateName,
			&trade.TradeTypeID,
			&trade.TradeDate,
			&trade.SettleDate,
			&trade.TransferDate,
			&trade.FromQuantity,
			&trade.ToQuantity,
			&trade.Price,
			&trade.TotalAmount,
			&trade.FeesAmount,
			&trade.FeesAssetID,
			&trade.RealizedReturnAmount,
			&trade.RealizedReturnAssetID,
			&trade.CostBasisAmount,
			&trade.CostBasisTradeID,
			&trade.Description,
			&trade.SourceData,
			&trade.IsActive,
			&trade.CreatedBy,
			&trade.CreatedAt,
			&trade.UpdatedBy,
			&trade.UpdatedAt,
		)
		trades = append(trades, trade)
	}
	return trades, nil
}

func RemoveTrade(tradeID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM trades WHERE id = $1`, tradeID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetTradeList() ([]Trade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
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
	source_data,
	is_active,
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
	trades := make([]Trade, 0)
	for results.Next() {
		var trade Trade
		results.Scan(
			&trade.ID,
			&trade.ParentTradeID,
			&trade.FromAccountID,
			&trade.ToAccountID,
			&trade.AssetID,
			&trade.SourceID,
			&trade.UUID,
			&trade.TransactionIDFromSource,
			&trade.OrderIDFromSource,
			&trade.TradeIDFromSource,
			&trade.Name,
			&trade.AlternateName,
			&trade.TradeTypeID,
			&trade.TradeDate,
			&trade.SettleDate,
			&trade.TransferDate,
			&trade.FromQuantity,
			&trade.ToQuantity,
			&trade.Price,
			&trade.TotalAmount,
			&trade.FeesAmount,
			&trade.FeesAssetID,
			&trade.RealizedReturnAmount,
			&trade.RealizedReturnAssetID,
			&trade.CostBasisAmount,
			&trade.CostBasisTradeID,
			&trade.Description,
			&trade.SourceData,
			&trade.IsActive,
			&trade.CreatedBy,
			&trade.CreatedAt,
			&trade.UpdatedBy,
			&trade.UpdatedAt,
		)

		trades = append(trades, trade)
	}
	return trades, nil
}

func UpdateTrade(trade Trade) error {
	// if the trade id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (trade.SourceID == nil || *trade.SourceID == 0) || (trade.AssetID == nil || *trade.AssetID == 0) {
		return errors.New("trade has invalid ID")
	}
	_, err := database.DbConn.ExecContext(ctx, `UPDATE trades SET 
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
		source_data=$26,
		is_active=$27,
		updated_by=$28, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE id=$29`,
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
		trade.SourceData,              //26
		trade.IsActive,                //27
		trade.UpdatedBy,               //28
		trade.ID,                      //29
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertTrade(trade Trade) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var tradeID int
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO trades  
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
		source_data,
		is_active,
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
		&trade.ParentTradeID,           //1
		&trade.FromAccountID,           //2
		&trade.ToAccountID,             //3
		&trade.AssetID,                 //4
		&trade.SourceID,                //5
		&trade.TransactionIDFromSource, //6
		&trade.OrderIDFromSource,       //7
		&trade.TradeIDFromSource,       //8
		&trade.Name,                    //9
		&trade.AlternateName,           //10
		&trade.TradeTypeID,             //11
		&trade.TradeDate,               //12
		&trade.SettleDate,              //13
		&trade.TransferDate,            //14
		&trade.FromQuantity,            //15
		&trade.ToQuantity,              //16
		&trade.Price,                   //17
		&trade.TotalAmount,             //18
		&trade.FeesAmount,              //19
		&trade.FeesAssetID,             //20
		&trade.RealizedReturnAmount,    //21
		&trade.RealizedReturnAssetID,   //22
		&trade.CostBasisAmount,         //23
		&trade.CostBasisTradeID,        //24
		&trade.Description,             //25
		&trade.SourceData,              //26
		&trade.IsActive,                //27
		trade.CreatedBy,                //28
	).Scan(&tradeID)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(tradeID), nil
}
func GetMinAndMaxTradeDates() (*MinMaxTradeDates, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
	MIN(trade_date),
	MAX(trade_date)
	FROM trades 
	`)

	minMaxTradeDates := &MinMaxTradeDates{}
	err := row.Scan(
		&minMaxTradeDates.MinTradeDate,
		&minMaxTradeDates.MaxTradeDate,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return minMaxTradeDates, nil
}
