package asset

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

func GetAsset(dbConnPgx utils.PgxIface, assetID *int) (*Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	cusip,
	ticker,
	base_asset_id,
	quote_asset_id,
	description,
	asset_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	chain_id,
	category_id,
	sub_category_id,
	is_default_quote,
	ignore_market_data,
	decimals,
	contract_address,
	starting_block_number,
	import_geth,
	import_geth_initial
	FROM assets 
	WHERE id = $1`, *assetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	// from https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx
	defer row.Close()
	asset, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Asset])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &asset, nil
}

// GetAssetByTicker : get asset by ticker
func GetAssetByTicker(dbConnPgx utils.PgxIface, ticker string) (*Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	cusip,
	ticker,
	base_asset_id,
	quote_asset_id,
	description,
	asset_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	chain_id,
	category_id,
	sub_category_id,
	is_default_quote,
	ignore_market_data,
	decimals,
	contract_address,
	starting_block_number,
	import_geth,
	import_geth_initial
	FROM assets 
	WHERE ticker = $1`, ticker)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer row.Close()
	asset, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Asset])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &asset, nil
}

// GetAssetByContractAddress : get asset by contract address
func GetAssetByContractAddress(dbConnPgx utils.PgxIface, contractAddress string) (*Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	cusip,
	ticker,
	base_asset_id,
	quote_asset_id,
	description,
	asset_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	chain_id,
	category_id,
	sub_category_id,
	is_default_quote,
	ignore_market_data,
	decimals,
	contract_address,
	starting_block_number,
	import_geth,
	import_geth_initial
	FROM assets 
	WHERE contract_address = $1`, contractAddress)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer row.Close()
	asset, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Asset])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &asset, nil
}

// GetAssetByCusip : get asset by cusip
func GetAssetByCusip(dbConnPgx utils.PgxIface, cusip string) (*Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	cusip,
	ticker,
	base_asset_id,
	quote_asset_id,
	description,
	asset_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	chain_id,
	category_id,
	sub_category_id,
	is_default_quote,
	ignore_market_data,
	decimals,
	contract_address,
	starting_block_number,
	import_geth,
	import_geth_initial
	FROM assets 
	WHERE cusip = $1`, cusip)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer row.Close()
	asset, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Asset])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &asset, nil
}

// GetAssetByBaseAndQuoteID : get asset by base and quote id
func GetAssetByBaseAndQuoteID(dbConnPgx utils.PgxIface, baseAssetID *int, quoteAssetID *int) (*Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row, err := dbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	cusip,
	ticker,
	base_asset_id,
	quote_asset_id,
	description,
	asset_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	chain_id,
	category_id,
	sub_category_id,
	is_default_quote,
	ignore_market_data,
	decimals,
	contract_address,
	starting_block_number,
	import_geth,
	import_geth_initial
	FROM assets 
	WHERE base_asset_id = $1 
	AND quote_asset_id = $2`,
		*baseAssetID, *quoteAssetID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer row.Close()
	asset, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Asset])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return &asset, nil
}

func GetGethImportAssets(dbConnPgx utils.PgxIface) ([]Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	cusip,
	ticker,
	base_asset_id,
	quote_asset_id,
	description,
	asset_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	chain_id,
	category_id,
	sub_category_id,
	is_default_quote,
	ignore_market_data,
	decimals,
	contract_address,
	starting_block_number,
	import_geth,
	import_geth_initial
	FROM assets 
	WHERE import_geth = TRUE
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assets, err := pgx.CollectRows(results, pgx.RowToStructByName[Asset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return assets, nil
}

func RemoveAsset(dbConnPgx utils.PgxIface, assetID *int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in RemoveAsset DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `DELETE FROM assets WHERE id = $1`
	//defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql, *assetID); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func GetCurrentTradingAssets(dbConnPgx utils.PgxIface) ([]Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	cusip,
	ticker,
	base_asset_id,
	quote_asset_id,
	description,
	asset_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	chain_id,
	category_id,
	sub_category_id,
	is_default_quote,
	ignore_market_data,
	decimals,
	contract_address,
	starting_block_number,
	import_geth,
	import_geth_initial
	FROM public.get_current_assets`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assets, err := pgx.CollectRows(results, pgx.RowToStructByName[Asset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return assets, nil

}

func GetCryptoAssets(dbConnPgx utils.PgxIface) ([]Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	cusip,
	ticker,
	base_asset_id,
	quote_asset_id,
	description,
	asset_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	chain_id,
	category_id,
	sub_category_id,
	is_default_quote,
	ignore_market_data,
	decimals,
	contract_address,
	starting_block_number,
	import_geth,
	import_geth_initial
	FROM assets
	where asset_type_id = 1
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assets, err := pgx.CollectRows(results, pgx.RowToStructByName[Asset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return assets, nil
}

func GetAssetsByAssetTypeAndSource(dbConnPgx utils.PgxIface, assetTypeID *int, sourceID *int, excludeIgnoreMarketData bool) ([]AssetWithSources, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
		assets.id,
		assets.uuid, 
		assets.name, 
		assets.alternate_name, 
		assets.cusip,
		assets.ticker,
		assets.base_asset_id,
		assets.quote_asset_id,
		assets.description,
		assets.asset_type_id,
		assets.created_by, 
		assets.created_at, 
		assets.updated_by, 
		assets.updated_at,
		assets.chain_id,
		assets.category_id,
		assets.sub_category_id,
		assets.is_default_quote,
		assets.ignore_market_data,
		assets.decimals,
		assets.contract_address,
		assets.starting_block_number,
		assets.import_geth,
		assets.import_geth_initial,
		assetSources.source_id,
		assetSources.source_identifier
		FROM assets assets
		JOIN asset_sources assetSources ON assets.id = assetSources.asset_id
		WHERE assets.asset_type_id = $1
		AND assetSources.source_id = $2`
	if excludeIgnoreMarketData {
		sql += `AND ignore_market_data = FALSE
		`
	}
	results, err := dbConnPgx.Query(ctx, sql, *assetTypeID, *sourceID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assetsWithSources := make([]AssetWithSources, 0)
	for results.Next() {
		var asset AssetWithSources
		results.Scan(
			&asset.Asset.ID,
			&asset.Asset.UUID,
			&asset.Asset.Name,
			&asset.Asset.AlternateName,
			&asset.Asset.Cusip,
			&asset.Asset.Ticker,
			&asset.Asset.BaseAssetID,
			&asset.Asset.QuoteAssetID,
			&asset.Asset.Description,
			&asset.Asset.AssetTypeID,
			&asset.Asset.CreatedBy,
			&asset.Asset.CreatedAt,
			&asset.Asset.UpdatedBy,
			&asset.Asset.UpdatedAt,
			&asset.Asset.ChainID,
			&asset.Asset.CategoryID,
			&asset.Asset.SubCategoryID,
			&asset.Asset.IsDefaultQuote,
			&asset.Asset.IgnoreMarketData,
			&asset.Asset.Decimals,
			&asset.Asset.ContractAddress,
			&asset.StartingBlockNumber,
			&asset.ImportGeth,
			&asset.ImportGethInitial,
			&asset.SourceID,
			&asset.SourceIdentifier,
		)

		assetsWithSources = append(assetsWithSources, asset)
	}
	return assetsWithSources, nil
}

func GetCryptoAssetsBySourceId(dbConnPgx utils.PgxIface, sourceID *int, excludeIgnoreMarketData bool) ([]AssetWithSources, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
	assets.id,
	assets.uuid, 
	assets.name, 
	assets.alternate_name, 
	assets.cusip,
	assets.ticker,
	assets.base_asset_id,
	assets.quote_asset_id,
	assets.description,
	assets.asset_type_id,
	assets.created_by, 
	assets.created_at, 
	assets.updated_by, 
	assets.updated_at,
	assets.chain_id,
	assets.category_id,
	assets.sub_category_id,
	assets.is_default_quote,
	assets.ignore_market_data,
	assets.decimals,
	contract_address,
	starting_block_number,
	import_geth,
	import_geth_initial,
	assetSources.source_id,
	assetSources.source_identifier
	FROM assets assets
	JOIN asset_sources assetSources ON assets.id = assetSources.asset_id
	WHERE assetSources.source_id = $1
	AND assets.asset_type_id = 1`
	if excludeIgnoreMarketData {
		sql += `AND ignore_market_data = FALSE
		`
	}
	results, err := dbConnPgx.Query(ctx, sql, *sourceID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assetsWithSources := make([]AssetWithSources, 0)
	for results.Next() {
		var asset AssetWithSources
		results.Scan(
			&asset.Asset.ID,
			&asset.Asset.UUID,
			&asset.Asset.Name,
			&asset.Asset.AlternateName,
			&asset.Asset.Cusip,
			&asset.Asset.Ticker,
			&asset.Asset.BaseAssetID,
			&asset.Asset.QuoteAssetID,
			&asset.Asset.Description,
			&asset.Asset.AssetTypeID,
			&asset.Asset.CreatedBy,
			&asset.Asset.CreatedAt,
			&asset.Asset.UpdatedBy,
			&asset.Asset.UpdatedAt,
			&asset.Asset.ChainID,
			&asset.Asset.CategoryID,
			&asset.Asset.SubCategoryID,
			&asset.Asset.IsDefaultQuote,
			&asset.Asset.IgnoreMarketData,
			&asset.Asset.Decimals,
			&asset.Asset.ContractAddress,
			&asset.StartingBlockNumber,
			&asset.ImportGeth,
			&asset.ImportGethInitial,
			&asset.SourceID,
			&asset.SourceIdentifier,
		)

		assetsWithSources = append(assetsWithSources, asset)
	}
	return assetsWithSources, nil

}

func GetCryptoAssetsBySourceID(dbConnPgx utils.PgxIface, sourceID *int, excludeIgnoreMarketData bool) ([]Asset, error) {
	assetsWithSources, err := GetCryptoAssetsBySourceId(dbConnPgx, sourceID, excludeIgnoreMarketData)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	results := make([]Asset, 0)
	for _, assetsWithSource := range assetsWithSources {
		results = append(results, assetsWithSource.Asset)
	}
	return results, nil
}

func GetAssetWithSourceByAssetIdAndSourceID(dbConnPgx utils.PgxIface, assetID, sourceID *int, excludeIgnoreMarketData bool) (*AssetWithSources, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	query := `SELECT 
	assets.id,
	assets.uuid, 
	assets.name, 
	assets.alternate_name, 
	assets.cusip,
	assets.ticker,
	assets.base_asset_id,
	assets.quote_asset_id,
	assets.description,
	assets.asset_type_id,
	assets.created_by, 
	assets.created_at, 
	assets.updated_by, 
	assets.updated_at,
	assets.chain_id,
	assets.category_id,
	assets.sub_category_id,
	assets.is_default_quote,
	assets.ignore_market_data,
	assets.decimals,
	assets.contract_address,
	assets.starting_block_number,
	assets.import_geth,
	assets.import_geth_initial,
	assetSources.source_id as source_id,
	assetSources.source_identifier as source_identifier
	FROM assets assets
	JOIN asset_sources assetSources ON assets.id = assetSources.asset_id
	WHERE 
		assets.id = $1
	AND assetSources.source_id = $2`
	if excludeIgnoreMarketData {
		query += `AND ignore_market_data = FALSE
		`
	}
	row := dbConnPgx.QueryRow(ctx, query, *assetID, *sourceID)
	assetWithSources := &AssetWithSources{}
	err := row.Scan(
		&assetWithSources.Asset.ID,
		&assetWithSources.Asset.UUID,
		&assetWithSources.Asset.Name,
		&assetWithSources.Asset.AlternateName,
		&assetWithSources.Asset.Cusip,
		&assetWithSources.Asset.Ticker,
		&assetWithSources.Asset.BaseAssetID,
		&assetWithSources.Asset.QuoteAssetID,
		&assetWithSources.Asset.Description,
		&assetWithSources.Asset.AssetTypeID,
		&assetWithSources.Asset.CreatedBy,
		&assetWithSources.Asset.CreatedAt,
		&assetWithSources.Asset.UpdatedBy,
		&assetWithSources.Asset.UpdatedAt,
		&assetWithSources.Asset.ChainID,
		&assetWithSources.Asset.CategoryID,
		&assetWithSources.Asset.SubCategoryID,
		&assetWithSources.Asset.IsDefaultQuote,
		&assetWithSources.Asset.IgnoreMarketData,
		&assetWithSources.Asset.Decimals,
		&assetWithSources.Asset.ContractAddress,
		&assetWithSources.Asset.StartingBlockNumber,
		&assetWithSources.Asset.ImportGeth,
		&assetWithSources.Asset.ImportGethInitial,
		&assetWithSources.SourceID,
		&assetWithSources.SourceIdentifier,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return assetWithSources, nil
}

func GetAssetWithSourceByAssetIdsAndSourceID(dbConnPgx utils.PgxIface, assetIDs []int, sourceID *int, excludeIgnoreMarketData bool) ([]AssetWithSources, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	query := `SELECT 
	assets.id,
	assets.uuid, 
	assets.name, 
	assets.alternate_name, 
	assets.cusip,
	assets.ticker,
	assets.base_asset_id,
	assets.quote_asset_id,
	assets.description,
	assets.asset_type_id,
	assets.created_by, 
	assets.created_at, 
	assets.updated_by, 
	assets.updated_at,
	assets.chain_id,
	assets.category_id,
	assets.sub_category_id,
	assets.is_default_quote,
	assets.ignore_market_data,
	assets.decimals,
	assets.contract_address,
	assets.starting_block_number,
	assets.import_geth,
	assets.import_geth_initial,
	assetSources.source_id,
	assetSources.source_identifier
	FROM assets assets
	JOIN asset_sources assetSources ON assets.id = assetSources.asset_id
	WHERE 
	assets.id = ANY($1)
	AND assetSources.source_id = $2`
	if excludeIgnoreMarketData {
		query += `AND ignore_market_data = FALSE
		`
	}
	results, err := dbConnPgx.Query(ctx, query, pq.Array(assetIDs), *sourceID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assetsWithSources := make([]AssetWithSources, 0)
	for results.Next() {
		var asset AssetWithSources
		results.Scan(
			&asset.Asset.ID,
			&asset.Asset.UUID,
			&asset.Asset.Name,
			&asset.Asset.AlternateName,
			&asset.Asset.Cusip,
			&asset.Asset.Ticker,
			&asset.Asset.BaseAssetID,
			&asset.Asset.QuoteAssetID,
			&asset.Asset.Description,
			&asset.Asset.AssetTypeID,
			&asset.Asset.CreatedBy,
			&asset.Asset.CreatedAt,
			&asset.Asset.UpdatedBy,
			&asset.Asset.UpdatedAt,
			&asset.Asset.ChainID,
			&asset.Asset.CategoryID,
			&asset.Asset.SubCategoryID,
			&asset.Asset.IsDefaultQuote,
			&asset.Asset.IgnoreMarketData,
			&asset.Asset.Decimals,
			&asset.Asset.ContractAddress,
			&asset.Asset.StartingBlockNumber,
			&asset.Asset.ImportGeth,
			&asset.Asset.ImportGethInitial,
			&asset.SourceID,
			&asset.SourceIdentifier,
		)

		assetsWithSources = append(assetsWithSources, asset)
	}
	return assetsWithSources, nil
}

func GetAssetList(dbConnPgx utils.PgxIface, ids []int) ([]Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	sql := `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	cusip,
	ticker,
	base_asset_id,
	quote_asset_id,
	description,
	asset_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	chain_id,
	category_id,
	sub_category_id,
	is_default_quote,
	ignore_market_data,
	decimals,
	contract_address,
	starting_block_number,
	import_geth,
	import_geth_initial
	FROM assets`
	if len(ids) > 0 {
		strIds := utils.SplitToString(ids, ",")
		additionalQuery := fmt.Sprintf(` WHERE id IN (%s)`, strIds)
		sql += additionalQuery
	}
	results, err := dbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assets, err := pgx.CollectRows(results, pgx.RowToStructByName[Asset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return assets, nil
}

func GetAssetsByChainId(dbConnPgx utils.PgxIface, chainID *int) ([]Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	cusip,
	ticker,
	base_asset_id,
	quote_asset_id,
	description,
	asset_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	chain_id,
	category_id,
	sub_category_id,
	is_default_quote,
	ignore_market_data,
	decimals,
	contract_address,
	starting_block_number,
	import_geth,
	import_geth_initial
	FROM assets
	WHERE chain_id = $1`, *chainID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assets, err := pgx.CollectRows(results, pgx.RowToStructByName[Asset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return assets, nil
}

// for refinedev
func GetAssetListByPagination(dbConnPgx utils.PgxIface, _start, _end *int, _order, _sort string, _filters []string) ([]Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	sql := `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	cusip,
	ticker,
	base_asset_id,
	quote_asset_id,
	description,
	asset_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	chain_id,
	category_id,
	sub_category_id,
	is_default_quote,
	ignore_market_data,
	decimals,
	contract_address,
	starting_block_number,
	import_geth,
	import_geth_initial
	FROM assets
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
	assets, err := pgx.CollectRows(results, pgx.RowToStructByName[Asset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return assets, nil
}

func GetTotalAssetCount(dbConnPgx utils.PgxIface) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	row := dbConnPgx.QueryRow(ctx, `SELECT COUNT(*) FROM assets`)
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
func GetDefaultQuoteAssetListBySourceID(dbConnPgx utils.PgxIface, sourceID *int) ([]AssetWithSources, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	assets.id,
	assets.uuid, 
	assets.name, 
	assets.alternate_name, 
	assets.cusip,
	assets.ticker,
	assets.base_asset_id,
	assets.quote_asset_id,
	assets.description,
	assets.asset_type_id,
	assets.created_by, 
	assets.created_at, 
	assets.updated_by, 
	assets.updated_at,
	assets.chain_id,
	assets.category_id,
	assets.sub_category_id,
	assets.is_default_quote,
	assets.ignore_market_data,
	assets.decimals,
	assets.contract_address,
	assets.starting_block_number,
	assets.import_geth,
	assets.import_geth_initial,
	assetSources.source_id,
	assetSources.source_identifier
	FROM get_default_quotes assets
	JOIN asset_sources assetSources ON assets.id = assetSources.asset_id
	WHERE assetSources.source_id = $1
	`, *sourceID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assetWithSources := make([]AssetWithSources, 0)
	for results.Next() {
		var asset AssetWithSources
		results.Scan(
			&asset.Asset.ID,
			&asset.Asset.UUID,
			&asset.Asset.Name,
			&asset.Asset.AlternateName,
			&asset.Asset.Cusip,
			&asset.Asset.Ticker,
			&asset.Asset.BaseAssetID,
			&asset.Asset.QuoteAssetID,
			&asset.Asset.Description,
			&asset.Asset.AssetTypeID,
			&asset.Asset.CreatedBy,
			&asset.Asset.CreatedAt,
			&asset.Asset.UpdatedBy,
			&asset.Asset.UpdatedAt,
			&asset.Asset.ChainID,
			&asset.Asset.CategoryID,
			&asset.Asset.SubCategoryID,
			&asset.Asset.IsDefaultQuote,
			&asset.Asset.IgnoreMarketData,
			&asset.Asset.Decimals,
			&asset.Asset.ContractAddress,
			&asset.StartingBlockNumber,
			&asset.ImportGeth,
			&asset.ImportGethInitial,
			&asset.SourceID,
			&asset.SourceIdentifier,
		)

		assetWithSources = append(assetWithSources, asset)
	}
	return assetWithSources, nil

}
func GetDefaultQuoteAssetList(dbConnPgx utils.PgxIface) ([]Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := dbConnPgx.Query(ctx, `SELECT 
	id,
	uuid, 
	name, 
	alternate_name, 
	cusip,
	ticker,
	base_asset_id,
	quote_asset_id,
	description,
	asset_type_id,
	created_by, 
	created_at, 
	updated_by, 
	updated_at,
	chain_id,
	category_id,
	sub_category_id,
	is_default_quote,
	ignore_market_data,
	decimals,
	contract_address,
	starting_block_number,
	import_geth,
	import_geth_initial
	FROM get_default_quotes`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assets, err := pgx.CollectRows(results, pgx.RowToStructByName[Asset])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return assets, nil
}

func UpdateAsset(dbConnPgx utils.PgxIface, asset *Asset) error {
	// if the asset id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if asset.ID == nil || *asset.ID == 0 {
		return errors.New("asset has invalid ID")
	}
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in UpdateAsset DbConn.Begin   %s", err.Error())
		return err
	}
	sql := `UPDATE assets SET 
		name=$1,  
		alternate_name=$2, 
		cusip=$3,
		ticker=$4,
		base_asset_id=$5,
		quote_asset_id=$6,
		description=$7,
		asset_type_id=$8,
		updated_by=$9, 
		updated_at=current_timestamp at time zone 'UTC',
		chain_id = $10,
		category_id = $11,
		sub_category_id = $12,
		is_default_quote = $13,
		ignore_market_data = $14,
		decimals = $15,
		contract_address=$16,
		starting_block_number=$17,
		import_geth = $18,
		import_geth_initial = $19
		WHERE id=$20`
	//defer dbConnPgx.Close()
	if _, err := dbConnPgx.Exec(ctx, sql,
		asset.Name,                //1
		asset.AlternateName,       //2
		asset.Cusip,               //3
		asset.Ticker,              //4
		asset.BaseAssetID,         //5
		asset.QuoteAssetID,        //6
		asset.Description,         //7
		asset.AssetTypeID,         //8
		asset.UpdatedBy,           //9
		asset.ChainID,             //10
		asset.CategoryID,          //11
		asset.SubCategoryID,       //12
		asset.IsDefaultQuote,      //13
		asset.IgnoreMarketData,    //14
		asset.Decimals,            //15
		asset.ContractAddress,     //16
		asset.StartingBlockNumber, //17
		asset.ImportGeth,          //18
		asset.ImportGethInitial,   //19
		asset.ID,                  //20
	); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func InsertAsset(dbConnPgx utils.PgxIface, asset *Asset) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	tx, err := dbConnPgx.Begin(ctx)
	if err != nil {
		log.Printf("Error in InsertAsset DbConn.Begin   %s", err.Error())
		return -1, err
	}
	var insertID int
	err = dbConnPgx.QueryRow(ctx, `INSERT INTO assets  
	(
		name, 
		uuid,
		alternate_name, 
		cusip,
		ticker,
		base_asset_id,
		quote_asset_id,
		description,
		asset_type_id,
		created_by, 
		created_at, 
		updated_by, 
		updated_at,
		chain_id,
		category_id,
		sub_category_id,
		is_default_quote,
		ignore_market_data,
		decimals,
		contract_address,
		starting_block_number,
		import_geth,
		import_geth_initial
		) VALUES (
			$1,
			uuid_generate_v4(), 
			$2, 
			$3, 
			$4, 
			$5,
			$6,
			$7,
			$8,
			$9,
			current_timestamp at time zone 'UTC',
			$9,
			current_timestamp at time zone 'UTC',
			$10,
			$11,
			$12,
			$13,
			$14,
			$15,
			$16,
			$17,
			$18,
			$19
		)
		RETURNING id`,
		asset.Name,                //1
		asset.AlternateName,       //2
		asset.Cusip,               //3
		asset.Ticker,              //4
		asset.BaseAssetID,         //5
		asset.QuoteAssetID,        //6
		asset.Description,         //7
		asset.AssetTypeID,         //8
		asset.CreatedBy,           //9
		asset.ChainID,             //10
		asset.CategoryID,          //11
		asset.SubCategoryID,       //12
		asset.IsDefaultQuote,      //13
		asset.IgnoreMarketData,    //14
		asset.Decimals,            //15
		asset.ContractAddress,     //16
		asset.StartingBlockNumber, //17
		asset.ImportGeth,          //18
		asset.ImportGethInitial,   //19
	).Scan(&insertID)

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
	return int(insertID), nil
}

func InsertAssets(dbConnPgx utils.PgxIface, assets []Asset) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rows := [][]interface{}{}
	for i := range assets {
		asset := assets[i]
		uuidString := &pgtype.UUID{}
		uuidString.Set(asset.UUID)
		row := []interface{}{
			uuidString,                //1
			asset.Name,                //2
			asset.AlternateName,       //3
			asset.Cusip,               //4
			asset.Ticker,              //5
			asset.BaseAssetID,         //6
			asset.QuoteAssetID,        //7
			asset.Description,         //8
			asset.AssetTypeID,         //9
			asset.CreatedBy,           //10
			&now,                      //11
			asset.CreatedBy,           //12
			&now,                      //13
			asset.ChainID,             //14
			asset.CategoryID,          //15
			asset.SubCategoryID,       //16
			asset.IsDefaultQuote,      //17
			asset.IgnoreMarketData,    //18
			asset.Decimals,            //19
			asset.ContractAddress,     //20
			asset.StartingBlockNumber, //21
			asset.ImportGeth,          //22
			asset.ImportGethInitial,   //23

		}
		rows = append(rows, row)
	}
	copyCount, err := dbConnPgx.CopyFrom(
		ctx,
		pgx.Identifier{"assets"},
		[]string{
			"uuid",                  //1
			"name",                  //2
			"alternate_name",        //3
			"cusip",                 //4
			"ticker",                //5
			"base_asset_id",         //6
			"quote_asset_id",        //7
			"description",           //8
			"asset_type_id",         //9
			"created_by",            //10
			"created_at",            //11
			"updated_by",            //12
			"updated_at",            //13
			"chain_id",              //14
			"category_id",           //15
			"sub_category_id",       //16
			"is_default_quote",      //17
			"ignore_market_data",    //18
			"decimals",              //19
			"contract_address",      //20
			"starting_block_number", //21
			"import_geth",           //22
			"import_geth_initial",   //23
		},
		pgx.CopyFromRows(rows),
	)
	log.Println(fmt.Printf("InsertAssets: copy count: %d", copyCount))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
