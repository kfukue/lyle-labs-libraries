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
	gethlyleswaps "github.com/kfukue/lyle-labs-libraries/gethlyle/swaps"
	"github.com/kfukue/lyle-labs-libraries/utils"
)

func GetAllGethTradeSwapsByTradeID(gethTradeID int) ([]GethTradeSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
	SELECT 
		geth_trade_swaps.geth_trade_id,
		geth_trade_swaps.geth_swap_id,
		geth_trade_swaps.uuid, 
		geth_trade_swaps.name, 
		geth_trade_swaps.alternate_name, 
		geth_trade_swaps.description,
		geth_trade_swaps.created_by, 
		geth_trade_swaps.created_at, 
		geth_trade_swaps.updated_by, 
		geth_trade_swaps.updated_at 
	FROM geth_trades 
	LEFT JOIN geth_trade_swaps
	ON geth_trade_swaps.geth_trade_id = geth_trades.id 
	WHERE 
	geth_trades.id = $1
	`, gethTradeID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTradeSwaps := make([]GethTradeSwap, 0)
	for results.Next() {
		var gethTradeSwap GethTradeSwap
		results.Scan(
			&gethTradeSwap.GethTradeID,
			&gethTradeSwap.GethSwapID,
			&gethTradeSwap.UUID,
			&gethTradeSwap.Name,
			&gethTradeSwap.AlternateName,
			&gethTradeSwap.Description,
			&gethTradeSwap.CreatedBy,
			&gethTradeSwap.CreatedAt,
			&gethTradeSwap.UpdatedBy,
			&gethTradeSwap.UpdatedAt,
		)

		gethTradeSwaps = append(gethTradeSwaps, gethTradeSwap)
	}
	return gethTradeSwaps, nil
}
func GetGethTradeSwap(gethGethSwapID int, gethTradeID int) (*GethTradeSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConnPgx.QueryRow(ctx, `
	SELECT 
		geth_trade_id,
		geth_swap_id,
		uuid, 
		name, 
		alternate_name, 
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM geth_trade_swaps 
	WHERE 
	geth_trade_id = $1
	AND geth_swap_id = $2
	`, gethTradeID, gethGethSwapID)

	gethTradeSwap := &GethTradeSwap{}
	err := row.Scan(
		&gethTradeSwap.GethTradeID,
		&gethTradeSwap.GethSwapID,
		&gethTradeSwap.UUID,
		&gethTradeSwap.Name,
		&gethTradeSwap.AlternateName,
		&gethTradeSwap.Description,
		&gethTradeSwap.CreatedBy,
		&gethTradeSwap.CreatedAt,
		&gethTradeSwap.UpdatedBy,
		&gethTradeSwap.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return gethTradeSwap, nil
}

func RemoveGethTradeSwap(gethTradeID int, gethGethSwapID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM geth_trade_swaps WHERE 
	geth_trade_id =$1 AND geth_swap_id = $2`, gethTradeID, gethGethSwapID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetGethTradeSwapList(gethTradeIds []int, swapIds []int) ([]GethTradeSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `
	SELECT 
		geth_trade_id,
		geth_swap_id,
		uuid, 
		name, 
		alternate_name, 
		description,
		created_by, 
		created_at, 
		updated_by, 
		updated_at 
	FROM geth_trade_swaps`
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
			additionalQuery += fmt.Sprintf(`geth_swap_id IN (%s)`, swapStrIds)
		}
		sql += additionalQuery
	}
	results, err := database.DbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethTradeSwaps := make([]GethTradeSwap, 0)
	for results.Next() {
		var gethTradeSwap GethTradeSwap
		results.Scan(
			&gethTradeSwap.GethTradeID,
			&gethTradeSwap.GethSwapID,
			&gethTradeSwap.UUID,
			&gethTradeSwap.Name,
			&gethTradeSwap.AlternateName,
			&gethTradeSwap.Description,
			&gethTradeSwap.CreatedBy,
			&gethTradeSwap.CreatedAt,
			&gethTradeSwap.UpdatedBy,
			&gethTradeSwap.UpdatedAt,
		)

		gethTradeSwaps = append(gethTradeSwaps, gethTradeSwap)
	}
	return gethTradeSwaps, nil
}

func UpdateGethTradeSwap(gethTradeSwap GethTradeSwap) error {
	// if the gethTradeSwap id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if (gethTradeSwap.GethSwapID == nil || *gethTradeSwap.GethSwapID == 0) || (gethTradeSwap.GethTradeID == nil || *gethTradeSwap.GethTradeID == 0) {
		return errors.New("gethTradeSwap has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE geth_trade_swaps SET 
		name=$1,  
		alternate_name=$2, 
		description=$3,
		updated_by=$4, 
		updated_at=current_timestamp at time zone 'UTC'
		WHERE 
			geth_trade_id=$5 AND
			geth_swap_id=$6 
		`,

		gethTradeSwap.Name,          //1
		gethTradeSwap.AlternateName, //2
		gethTradeSwap.Description,   //3
		gethTradeSwap.UpdatedBy,     //4
		gethTradeSwap.GethTradeID,   //5
		gethTradeSwap.GethSwapID,    //6
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertGethTradeSwap(gethTradeSwap GethTradeSwap) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var GethSwapID int
	var GethTradeID int
	err := database.DbConnPgx.QueryRow(ctx, `INSERT INTO geth_trade_swaps  
	(
		geth_trade_id,
		geth_swap_id,
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
			uuid_generate_v4(),
			$3,
			$4,
			$5,
			$6,
			current_timestamp at time zone 'UTC',
			$6,
			current_timestamp at time zone 'UTC'
		)
		RETURNING geth_trade_id, geth_swap_id`,
		gethTradeSwap.GethTradeID,   //1
		gethTradeSwap.GethSwapID,    //2
		gethTradeSwap.Name,          //3
		gethTradeSwap.AlternateName, //4
		gethTradeSwap.Description,   //5
		gethTradeSwap.CreatedBy,     //6
	).Scan(&GethTradeID, &GethSwapID)
	if err != nil {
		log.Println(err.Error())
		return 0, 0, err
	}
	return int(GethTradeID), int(GethSwapID), nil
}

func InsertGethTradeSwaps(gethTradeSwaps []*GethTradeSwap) error {
	// need to supply uuid
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range gethTradeSwaps {
		gethTradeSwap := gethTradeSwaps[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(gethTradeSwap.UUID)
		row := []interface{}{
			gethTradeSwap.GethTradeID,   //1
			gethTradeSwap.GethSwapID,    //2
			uuidString,                  //3
			gethTradeSwap.Name,          //4
			gethTradeSwap.AlternateName, //5
			gethTradeSwap.Description,   //6
			gethTradeSwap.CreatedBy,     //7
			&now,                        //8
			gethTradeSwap.CreatedBy,     //9
			&now,                        //10
		}
		rows = append(rows, row)
	}
	copyCount, err := database.DbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"geth_trade_swaps"},
		[]string{
			"geth_trade_id",  //1
			"geth_swap_id",   //2
			"uuid",           //3
			"name",           //4
			"alternate_name", //5
			"description",    //6
			"created_by",     //7
			"created_at",     //8
			"updated_by",     //9
			"updated_at",     //10
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

func GetMissingTradesFromSwapsByBaseAssetID(baseAssetID *int) ([]gethlyleswaps.GethSwap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
		SELECT
			gs.id,
			gs.uuid,
			gs.chain_id,
			gs.exchange_id,
			gs.block_number,
			gs.index_number,
			gs.swap_date,
			gs.trade_type_id,
			gs.txn_hash,
			gs.maker_address,
			gs.maker_address_id,
			gs.is_buy,
			gs.price,
			gs.price_usd,
			gs.token1_price_usd,
			gs.total_amount_usd,
			gs.pair_address,
			gs.liquidity_pool_id,
			gs.token0_asset_id,
			gs.token1_asset_id,
			gs.token0_amount,
			gs.token1_Amount,
			gs.description,
			gs.created_by,
			gs.created_at,
			gs.updated_by,
			gs.updated_at,
			gs.geth_process_job_id,
			gs.topics_str,
			gs.status_id,
			gs.base_asset_id
		FROM geth_swaps gs
		LEFT JOIN geth_trade_swaps gts
			ON gs.id = gts.geth_swap_id
		WHERE
		gts.geth_swap_id IS NULL
		AND 
		gs.base_asset_id = $1
		ORDER BY gs.block_number asc
		`,
		baseAssetID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	gethSwaps := make([]gethlyleswaps.GethSwap, 0)
	for results.Next() {
		var gethSwap gethlyleswaps.GethSwap
		results.Scan(
			&gethSwap.ID,
			&gethSwap.UUID,
			&gethSwap.ChainID,
			&gethSwap.ExchangeID,
			&gethSwap.BlockNumber,
			&gethSwap.IndexNumber,
			&gethSwap.SwapDate,
			&gethSwap.TradeTypeID,
			&gethSwap.TxnHash,
			&gethSwap.MakerAddress,
			&gethSwap.MakerAddressID,
			&gethSwap.IsBuy,
			&gethSwap.Price,
			&gethSwap.PriceUSD,
			&gethSwap.Token1PriceUSD,
			&gethSwap.TotalAmountUSD,
			&gethSwap.PairAddress,
			&gethSwap.LiquidityPoolID,
			&gethSwap.Token0AssetId,
			&gethSwap.Token1AssetId,
			&gethSwap.Token0Amount,
			&gethSwap.Token1Amount,
			&gethSwap.Description,
			&gethSwap.CreatedBy,
			&gethSwap.CreatedAt,
			&gethSwap.UpdatedBy,
			&gethSwap.UpdatedAt,
			&gethSwap.GethProcessJobID,
			&gethSwap.TopicsStr,
			&gethSwap.StatusID,
			&gethSwap.BaseAssetID,
		)

		gethSwaps = append(gethSwaps, gethSwap)
	}
	return gethSwaps, nil
}

func GetMissingTxnHashesFromSwapsByBaseAssetID(baseAssetID, maxBlockNumber *int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `
		SELECT
			DISTINCT gs.txn_hash
		FROM geth_swaps gs
		LEFT JOIN geth_trade_swaps gts
			ON gs.id = gts.geth_swap_id
		WHERE
			gts.geth_swap_id IS NULL
				AND 
			gs.base_asset_id = $1
				AND
			gs.status_id = $2
			AND
			gs.block_number > $3
		`,
		*baseAssetID, utils.SUCCESS_STRUCTURED_VALUE_ID, *maxBlockNumber,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	txnHashes := make([]string, 0)
	for results.Next() {
		var txnHash string
		results.Scan(
			&txnHash,
		)

		txnHashes = append(txnHashes, txnHash)
	}
	return txnHashes, nil
}

func GetMinMaxBlocksOfMissingSwapByBaseAssetID(baseAssetID *int) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var minBlock, maxBlock int
	err := database.DbConnPgx.QueryRow(ctx, `SELECT
			MIN(gs.block_number),
			MAX(gs.block_number)
		FROM geth_swaps gs
		LEFT JOIN geth_trade_swaps gts
			ON gs.id = gts.geth_swap_id
		WHERE
			gts.geth_swap_id IS NULL
				AND 
			gs.base_asset_id = $1
				AND
			gs.status_id = $2
		`,
		*baseAssetID, utils.SUCCESS_STRUCTURED_VALUE_ID,
	).Scan(&minBlock, &maxBlock)
	if err != nil {
		log.Println(err.Error())
		return -1, -1, err
	}
	return minBlock, maxBlock, nil
}

func GetFirstNonProcessedSwapBlockNumberForTrades(baseAssetID *int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var startingBlock int
	err := database.DbConnPgx.QueryRow(ctx, `
	WITH max_existing_block_swaps as (
		SELECT COALESCE(MAX(block_number),0) as min_block_number from  geth_swaps gs
		LEFT JOIN geth_trade_swaps gts
			ON gs.id = gts.geth_swap_id
		WHERE
			gts.geth_swap_id IS NOT NULL
				AND 
			gs.base_asset_id = $1
				AND
			gs.status_id = $2
		),
		min_missing_block_swaps as (
		SELECT COALESCE(MIN(block_number),0) as min_block_number from  geth_swaps gs
		LEFT JOIN geth_trade_swaps gts
			ON gs.id = gts.geth_swap_id
		WHERE
			gts.geth_swap_id IS NULL
				AND 
			gs.base_asset_id = $1
				AND
			gs.status_id = $2
			)
		SELECT MAX(min_block_number) 
		FROM (
			SELECT min_block_number FROM max_existing_block_swaps
			UNION 
			SELECT min_block_number FROM min_missing_block_swaps
				)   missing_and_existing_swap
		`, *baseAssetID, utils.SUCCESS_STRUCTURED_VALUE_ID,
	).Scan(&startingBlock)
	if err != nil {
		log.Println(err.Error())
		return -1, err
	}
	return startingBlock, nil
}
