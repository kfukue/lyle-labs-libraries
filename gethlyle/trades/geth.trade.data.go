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
	"github.com/lib/pq"
)

func GetGethTrade(dbConnPgx utils.PgxIface, gethTradeID *int) (*GethTrade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT
		id,
		uuid,
		name,
		alternate_name,
		address_str,
		address_id,
		trade_date,
		txn_hash,
		token0_amount,
		token0_amount_decimal_adj,
	  	token1_amount,
  		token1_amount_decimal_adj,
		is_buy,
		price,
		price_usd,
		lp_token1_price_usd,
		total_amount_usd,
		token0_asset_id,
		token1_asset_id,
		geth_process_job_id,
		status_id,
		trade_type_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		base_asset_id,
		oracle_price_usd,
  		oracle_price_asset_id
	FROM geth_trades
	WHERE id = $1
	`, *gethTradeID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	gethTrade, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethTrade])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethTrade, nil
}

func GetGethTradeByStartAndEndDates(dbConnPgx utils.PgxIface, startDate, endDate time.Time) ([]GethTrade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
		SELECT
			id,
			uuid,
			name,
			alternate_name,
			address_str,
			address_id,
			trade_date,
			txn_hash,
			token0_amount,
			token0_amount_decimal_adj,
			token1_amount,
			token1_amount_decimal_adj,
			is_buy,
			price,
			price_usd,
			lp_token1_price_usd,
			total_amount_usd,
			token0_asset_id,
			token1_asset_id,
			geth_process_job_id,
			status_id,
			trade_type_id,
			description,
			created_by,
			created_at,
			updated_by,
			updated_at,
			base_asset_id,
			oracle_price_usd,
			oracle_price_asset_id
		FROM geth_trades
		WHERE trade_date BETWEEN $1 AND $2
		`,
		startDate.Format(utils.LayoutPostgres), endDate.Format(utils.LayoutPostgres))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethTrades, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTrade])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTrades, nil
}

func GetGethTradeByFromAddress(dbConnPgx utils.PgxIface, addressStr string) ([]GethTrade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
		SELECT
			id,
			uuid,
			name,
			alternate_name,
			address_str,
			address_id,
			trade_date,
			txn_hash,
			token0_amount,
			token0_amount_decimal_adj,
			token1_amount,
			token1_amount_decimal_adj,
			is_buy,
			price,
			price_usd,
			lp_token1_price_usd,
			total_amount_usd,
			token0_asset_id,
			token1_asset_id,
			geth_process_job_id,
			status_id,
			trade_type_id,
			description,
			created_by,
			created_at,
			updated_by,
			updated_at,
			base_asset_id,
			oracle_price_usd,
			oracle_price_asset_id
		FROM geth_trades
		WHERE
		address_str = $1
		ORDER BY gethTrade_date asc`,
		addressStr,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethTrades, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTrade])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTrades, nil
}

func GetGethTradeByFromAddressId(dbConnPgx utils.PgxIface, addressID *int) ([]GethTrade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
		SELECT
			id,
			uuid,
			name,
			alternate_name,
			address_str,
			address_id,
			trade_date,
			txn_hash,
			token0_amount,
			token0_amount_decimal_adj,
			token1_amount,
			token1_amount_decimal_adj,
			is_buy,
			price,
			price_usd,
			lp_token1_price_usd,
			total_amount_usd,
			token0_asset_id,
			token1_asset_id,
			geth_process_job_id,
			status_id,
			trade_type_id,
			description,
			created_by,
			created_at,
			updated_by,
			updated_at,
			base_asset_id,
			oracle_price_usd,
			oracle_price_asset_id
		FROM geth_trades
		WHERE
		address_id = $1
		ORDER BY trade_date asc`,
		*addressID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethTrades, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTrade])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTrades, nil
}

func GetGethTradeByUUIDs(dbConnPgx utils.PgxIface, UUIDList []string) ([]GethTrade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
		SELECT 
			id,
			uuid,
			name,
			alternate_name,
			address_str,
			address_id,
			trade_date,
			txn_hash,
			token0_amount,
			token0_amount_decimal_adj,
			token1_amount,
			token1_amount_decimal_adj,
			is_buy,
			price,
			price_usd,
			lp_token1_price_usd,
			total_amount_usd,
			token0_asset_id,
			token1_asset_id,
			geth_process_job_id,
			status_id,
			trade_type_id,
			description,
			created_by,
			created_at,
			updated_by,
			updated_at,
			base_asset_id,
			oracle_price_usd,
			oracle_price_asset_id
		FROM geth_trades
		WHERE text(uuid) = ANY($1)
	`, pq.Array(UUIDList))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethTrades, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTrade])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTrades, nil
}

func GetNetTransfersByTxnHashAndAddressStrs(dbConnPgx utils.PgxIface, txnHash, addressStr string, baseAssetID *int) ([]NetTransferByAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	WITH to_address as (
		SELECT to_address as receiving_address, asset_id as asset_id, SUM(amount) as in_amount FROM geth_transfers 
		WHERE txn_hash = $1
		AND base_asset_id =$3
		GROUP BY to_address, asset_id
		),
		sender_address as (
		SELECT sender_address as sender_address,asset_id as asset_id, SUM(-amount) as out_amount FROM geth_transfers 
		WHERE txn_hash = $1
		AND base_asset_id =$3
		GROUP BY sender_address, asset_id
			)
	SELECT
		txn_hash,
		address,
		asset_id,
		net_amount,
		assets.*
		FROM (
		SELECT receiving_address as address, asset_id, in_amount as net_amount FROM to_address
		UNION 
		SELECT  sender_address as address, asset_id, out_amount as net_amount FROM sender_address
			) addresses 
		LEFT JOIN assets assets
				ON addresses.asset_id = assets.id
	WHERE address = $2
	`, txnHash, addressStr, *baseAssetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	netTransfersByAddress := make([]NetTransferByAddress, 0)
	for results.Next() {
		var netTransferByAddress NetTransferByAddress
		results.Scan(
			&netTransferByAddress.TxnHash,
			&netTransferByAddress.AddressStr,
			&netTransferByAddress.AssetID,
			&netTransferByAddress.NetAmount,
			&netTransferByAddress.Asset.ID,
			&netTransferByAddress.Asset.UUID,
			&netTransferByAddress.Asset.Name,
			&netTransferByAddress.Asset.AlternateName,
			&netTransferByAddress.Asset.Cusip,
			&netTransferByAddress.Asset.Ticker,
			&netTransferByAddress.Asset.BaseAssetID,
			&netTransferByAddress.Asset.QuoteAssetID,
			&netTransferByAddress.Asset.Description,
			&netTransferByAddress.Asset.AssetTypeID,
			&netTransferByAddress.Asset.CreatedBy,
			&netTransferByAddress.Asset.CreatedAt,
			&netTransferByAddress.Asset.UpdatedBy,
			&netTransferByAddress.Asset.UpdatedAt,
			&netTransferByAddress.Asset.ChainID,
			&netTransferByAddress.Asset.CategoryID,
			&netTransferByAddress.Asset.SubCategoryID,
			&netTransferByAddress.Asset.IsDefaultQuote,
			&netTransferByAddress.Asset.IgnoreMarketData,
			&netTransferByAddress.Asset.Decimals,
			&netTransferByAddress.Asset.ContractAddress,
			&netTransferByAddress.StartingBlockNumber,
			&netTransferByAddress.ImportGeth,
			&netTransferByAddress.ImportGethInitial,
		)
	}
	return netTransfersByAddress, nil
}

func GetFromNetTransfersByTxnHashesAndAddressStrs(dbConnPgx utils.PgxIface, txnHashes []string, baseAssetID *int) ([]NetTransferByAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	WITH to_address as (
		SELECT 
			txn_hash as txn_hash,
			to_address as receiving_address,
			asset_id as asset_id,
			SUM(amount) as in_amount 
		FROM geth_transfers 
		WHERE
			txn_hash = ANY($1)
			AND base_asset_id =$2
		GROUP BY 
			txn_hash,
			to_address,
			asset_id
		),
	sender_address as (
		SELECT
			txn_hash as txn_hash,
			sender_address as sender_address,
			asset_id as asset_id,
			SUM(-amount) as out_amount 
		FROM geth_transfers 
		WHERE 
			txn_hash = ANY($1)
			AND base_asset_id =$2
		GROUP BY
			txn_hash,
			sender_address,
			asset_id
	)
	SELECT
		addresses.txn_hash,
		addresses.address,
		addresses.asset_id,
		addresses.net_amount,
		assets.*
		FROM (
			SELECT txn_hash, receiving_address as address, asset_id, in_amount as net_amount FROM to_address
			UNION 
			SELECT txn_hash, sender_address as address, asset_id, out_amount as net_amount FROM sender_address
		) addresses 
		LEFT JOIN assets assets
			ON addresses.asset_id = assets.id
	`, pq.Array(txnHashes), *baseAssetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	netTransfersByAddress := make([]NetTransferByAddress, 0)
	for results.Next() {
		var netTransferByAddress NetTransferByAddress
		results.Scan(
			&netTransferByAddress.TxnHash,
			&netTransferByAddress.AddressStr,
			&netTransferByAddress.AssetID,
			&netTransferByAddress.NetAmount,
			&netTransferByAddress.Asset.ID,
			&netTransferByAddress.Asset.UUID,
			&netTransferByAddress.Asset.Name,
			&netTransferByAddress.Asset.AlternateName,
			&netTransferByAddress.Asset.Cusip,
			&netTransferByAddress.Asset.Ticker,
			&netTransferByAddress.Asset.BaseAssetID,
			&netTransferByAddress.Asset.QuoteAssetID,
			&netTransferByAddress.Asset.Description,
			&netTransferByAddress.Asset.AssetTypeID,
			&netTransferByAddress.Asset.CreatedBy,
			&netTransferByAddress.Asset.CreatedAt,
			&netTransferByAddress.Asset.UpdatedBy,
			&netTransferByAddress.Asset.UpdatedAt,
			&netTransferByAddress.Asset.ChainID,
			&netTransferByAddress.Asset.CategoryID,
			&netTransferByAddress.Asset.SubCategoryID,
			&netTransferByAddress.Asset.IsDefaultQuote,
			&netTransferByAddress.Asset.IgnoreMarketData,
			&netTransferByAddress.Asset.Decimals,
			&netTransferByAddress.Asset.ContractAddress,
			&netTransferByAddress.StartingBlockNumber,
			&netTransferByAddress.ImportGeth,
			&netTransferByAddress.ImportGethInitial,
		)
	}
	return netTransfersByAddress, nil
}

func GetStartAndEndBlockForNewTradesByBaseAssetID(dbConnPgx utils.PgxIface, baseAssetID *int) (*uint64, *uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	// TODO: Need to relook at this logic
	row := dbConnPgx.QueryRow(ctx, `SELECT
		MIN(geth_swaps.BlockNumber) as start_block_number,
		MAX(geth_swaps.BlockNumber) as end_block_number
	FROM geth_trade_swaps
	LEFT JOIN geth_swaps ON geth_trade_swaps.get_swap_id = geth_swaps.id
	WHERE geth_swaps.base_asset_id = $1
	AND geth_swaps.get_swap_id = NULL
	`, *baseAssetID)

	var startBlockNumber, endBlockNumber uint64
	err := row.Scan(
		&startBlockNumber,
		&endBlockNumber,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	return &startBlockNumber, &endBlockNumber, nil

}

func RemoveGethTrade(dbConnPgx utils.PgxIface, gethTradeID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethTrade DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_trades WHERE id = $1`

	if _, err := dbConnPgx.Exec(ctx, sql, *gethTradeID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func DeleteGethTradesByBaseAssetId(dbConnPgx utils.PgxIface, baseAssetID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveGethTrade DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM geth_trades WHERE base_asset_id = $1`

	if _, err := dbConnPgx.Exec(ctx, sql, *baseAssetID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetGethTradeList(dbConnPgx utils.PgxIface) ([]GethTrade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `
	SELECT
	id,
		uuid,
		name,
		alternate_name,
		address_str,
		address_id,
		trade_date,
		txn_hash,
		token0_amount,
		token0_amount_decimal_adj,
		token1_amount,
		token1_amount_decimal_adj,
		is_buy,
		price,
		price_usd,
		lp_token1_price_usd,
		total_amount_usd,
		token0_asset_id,
		token1_asset_id,
		geth_process_job_id,
		status_id,
		trade_type_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		base_asset_id,
		oracle_price_usd,
		oracle_price_asset_id
	FROM geth_trades `)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	gethTrades, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTrade])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTrades, nil
}

func UpdateGethTrade(dbConnPgx utils.PgxIface, gethTrade *GethTrade) error {
	// if the gethTrade id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if gethTrade.ID == nil {
		return errors.New("gethTrade has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateGethTrade DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE geth_trades SET 
		name=$1,
		alternate_name=$2,
		address_str=$3,
		address_id=$4,
		trade_date=$5,
		txn_hash=$6,
		token0_amount=$7,
		token0_amount_decimal_adj=$8,
		token1_amount=$9,
		token1_amount_decimal_adj=$10,
		is_buy=$11,
		price=$12,
		price_usd=$13,
		lp_token1_price_usd=$14,
		total_amount_usd=$15,
		token0_asset_id=$16,
		token1_asset_id=$17,
		geth_process_job_id=$18,
		status_id=$19,
		trade_type_id=$20,
		description=$21,
		updated_by=$22,
		updated_at=current_timestamp at time zone 'UTC',
		base_asset_id=$23,
		oracle_price_usd=$24,
		oracle_price_asset_id=$25
		WHERE id=$26`

	if _, err := dbConnPgx.Exec(ctx, sql,

		gethTrade.Name,                   //1
		gethTrade.AlternateName,          //2
		gethTrade.AddressStr,             //3
		gethTrade.AddressID,              //4
		gethTrade.TradeDate,              //5
		gethTrade.TxnHash,                //6
		gethTrade.Token0Amount,           //7
		gethTrade.Token0AmountDecimalAdj, //8
		gethTrade.Token1Amount,           //9
		gethTrade.Token1AmountDecimalAdj, //10
		gethTrade.IsBuy,                  //11
		gethTrade.Price,                  //12
		gethTrade.PriceUSD,               //13
		gethTrade.LPToken1PriceUSD,       //14
		gethTrade.TotalAmountUSD,         //15
		gethTrade.Token0AssetId,          //16
		gethTrade.Token1AssetId,          //17
		gethTrade.GethProcessJobID,       //18
		gethTrade.StatusID,               //19
		gethTrade.TradeTypeID,            //20
		gethTrade.Description,            //21
		gethTrade.UpdatedBy,              //22
		gethTrade.BaseAssetID,            //23
		gethTrade.OraclePriceUSD,         //24
		gethTrade.OraclePriceAssetID,     //25
		gethTrade.ID,                     //26
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertGethTrade(dbConnPgx utils.PgxIface, gethTrade *GethTrade) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertGethTrade DbConn.Begin   %s", err.Error())
		return -1, "", err
	}
	var gethTradeID int
	var gethTradeUUID string
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO geth_trades
	(
		uuid,
		name,
		alternate_name,
		address_str,
		address_id,
		trade_date,
		txn_hash,
		token0_amount,
		token0_amount_decimal_adj,		
		token1_amount,
		token1_amount_decimal_adj,
		is_buy,
		price,
		price_usd,
		lp_token1_price_usd,
		total_amount_usd,
		token0_asset_id,
		token1_asset_id,
		geth_process_job_id,
		status_id,
		trade_type_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		base_asset_id,
		oracle_price_usd,
		oracle_price_asset_id
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
		$15,
		$16,
		$17,
		$18,
		$19,
		$20,
		$21,
		$22,
		current_timestamp at time zone 'UTC',
		$22,
		current_timestamp at time zone 'UTC',
		$23,
		$24,
		$25
		)
		RETURNING id, uuid`,
		gethTrade.Name,                   //1
		gethTrade.AlternateName,          //2
		gethTrade.AddressStr,             //3
		gethTrade.AddressID,              //4
		gethTrade.TradeDate,              //5
		gethTrade.TxnHash,                //6
		gethTrade.Token0Amount,           //7
		gethTrade.Token0AmountDecimalAdj, //8
		gethTrade.Token1Amount,           //9
		gethTrade.Token1AmountDecimalAdj, //10
		gethTrade.IsBuy,                  //11
		gethTrade.Price,                  //12
		gethTrade.PriceUSD,               //13
		gethTrade.LPToken1PriceUSD,       //14
		gethTrade.TotalAmountUSD,         //15
		gethTrade.Token0AssetId,          //16
		gethTrade.Token1AssetId,          //17
		gethTrade.GethProcessJobID,       //18
		gethTrade.StatusID,               //19
		gethTrade.TradeTypeID,            //20
		gethTrade.Description,            //21
		gethTrade.CreatedBy,              //22
		gethTrade.BaseAssetID,            //23
		gethTrade.OraclePriceUSD,         //24
		gethTrade.OraclePriceAssetID,     //25
	).Scan(&gethTradeID, &gethTradeUUID)
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
	return int(gethTradeID), gethTradeUUID, nil
}

func InsertGethTrades(dbConnPgx utils.PgxIface, gethTrades []GethTrade) error {
	// need to supply uuid
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range gethTrades {
		gethTrade := gethTrades[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(gethTrade.UUID)
		row := []interface{}{
			uuidString,                       //1
			gethTrade.Name,                   //2
			gethTrade.AlternateName,          //3
			gethTrade.AddressStr,             //4
			gethTrade.AddressID,              //5
			gethTrade.TradeDate,              //6
			gethTrade.TxnHash,                //7
			gethTrade.Token0Amount,           //8
			gethTrade.Token0AmountDecimalAdj, //9
			gethTrade.Token1Amount,           //10
			gethTrade.Token1AmountDecimalAdj, //11
			gethTrade.IsBuy,                  //12
			gethTrade.Price,                  //13
			gethTrade.PriceUSD,               //14
			gethTrade.LPToken1PriceUSD,       //15
			gethTrade.TotalAmountUSD,         //16
			gethTrade.Token0AssetId,          //17
			gethTrade.Token1AssetId,          //18
			gethTrade.GethProcessJobID,       //19
			gethTrade.StatusID,               //20
			gethTrade.TradeTypeID,            //21
			gethTrade.Description,            //22
			gethTrade.CreatedBy,              //23
			&now,                             //24
			gethTrade.CreatedBy,              //25
			&now,                             //26
			gethTrade.BaseAssetID,            //27
			gethTrade.OraclePriceUSD,         //28
			gethTrade.OraclePriceAssetID,     //29
		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_trades"},
		[]string{
			"uuid",                      //1
			"name",                      //2
			"alternate_name",            //3
			"address_str",               //4
			"address_id",                //5
			"trade_date",                //6
			"txn_hash",                  //7
			"token0_amount",             //8
			"token0_amount_decimal_adj", //9
			"token1_amount",             //10
			"token1_amount_decimal_adj", //11
			"is_buy",                    //12
			"price",                     //13
			"price_usd",                 //14
			"lp_token1_price_usd",       //15
			"total_amount_usd",          //16
			"token0_asset_id",           //17
			"token1_asset_id",           //18
			"geth_process_job_id",       //19
			"status_id",                 //20
			"trade_type_id",             //21
			"description",               //22
			"created_by",                //23
			"created_at",                //24
			"updated_by",                //25
			"updated_at",                //26
			"base_asset_id",             //27
			"oracle_price_usd",          //28
			"oracle_price_asset_id",     //29
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertGethTrades: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetLatestGethTradeFromAssetIDAnDate(dbConnPgx utils.PgxIface, assetID *int, asOfDate time.Time, isBefore *bool) (*GethTrade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	selectSQL := `SELECT
		id,
		uuid,
		name,
		alternate_name,
		address_str,
		address_id,
		trade_date,
		txn_hash,
		token0_amount,
		token0_amount_decimal_adj,
		token1_amount,
		token1_amount_decimal_adj,
		is_buy,
		price,
		price_usd,
		lp_token1_price_usd,
		total_amount_usd,
		token0_asset_id,
		token1_asset_id,
		geth_process_job_id,
		status_id,
		trade_type_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		base_asset_id,
		oracle_price_usd,
		oracle_price_asset_id
	FROM geth_trades
	WHERE base_asset_id = $1 `
	tradeDateSQL := ``
	if *isBefore {
		tradeDateSQL = ` AND trade_date <= $2 ORDER BY trade_date desc`
	} else {
		tradeDateSQL = ` AND trade_date >= $2 ORDER BY trade_date asc`
	}
	sql := selectSQL + tradeDateSQL + ` LIMIT 1`
	row, err := dbConnPgx.Query(ctx, sql, *assetID, asOfDate.Format(utils.LayoutPostgres))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	gethTrade, err := pgx.CollectOneRow(row, pgx.RowToStructByName[GethTrade])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &gethTrade, nil
}

// for refinedev
func GetGethTradeListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]GethTrade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT
		id,
		uuid,
		name,
		alternate_name,
		address_str,
		address_id,
		trade_date,
		txn_hash,
		token0_amount,
		token0_amount_decimal_adj,
	  	token1_amount,
  		token1_amount_decimal_adj,
		is_buy,
		price,
		price_usd,
		lp_token1_price_usd,
		total_amount_usd,
		token0_asset_id,
		token1_asset_id,
		geth_process_job_id,
		status_id,
		trade_type_id,
		description,
		created_by,
		created_at,
		updated_by,
		updated_at,
		base_asset_id,
		oracle_price_usd,
  		oracle_price_asset_id
	FROM geth_trades 
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

	gethTrades, err := pgx.CollectRows(results, pgx.RowToStructByName[GethTrade])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTrades, nil
}

func GetTotalTradesCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT 
	COUNT(*)
	FROM geth_trades
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
