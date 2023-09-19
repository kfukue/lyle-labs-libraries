package asset

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kfukue/lyle-labs-libraries/database"
	"github.com/kfukue/lyle-labs-libraries/utils"
	"github.com/lib/pq"
)

func GetAsset(assetID int) (*Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	import_geth
	FROM assets 
	WHERE id = $1`, assetID)

	asset := &Asset{}
	err := row.Scan(
		&asset.ID,
		&asset.UUID,
		&asset.Name,
		&asset.AlternateName,
		&asset.Cusip,
		&asset.Ticker,
		&asset.BaseAssetID,
		&asset.QuoteAssetID,
		&asset.Description,
		&asset.AssetTypeID,
		&asset.CreatedBy,
		&asset.CreatedAt,
		&asset.UpdatedBy,
		&asset.UpdatedAt,
		&asset.ChainID,
		&asset.CategoryID,
		&asset.SubCategoryID,
		&asset.IsDefaultQuote,
		&asset.IgnoreMarketData,
		&asset.Decimals,
		&asset.ContractAddress,
		&asset.StartingBlockNumber,
		&asset.ImportGeth,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return asset, nil
}

// GetAssetByTicker : get asset by ticker
func GetAssetByTicker(ticker string) (*Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	import_geth
	FROM assets 
	WHERE ticker = $1`, ticker)

	asset := &Asset{}
	err := row.Scan(
		&asset.ID,
		&asset.UUID,
		&asset.Name,
		&asset.AlternateName,
		&asset.Cusip,
		&asset.Ticker,
		&asset.BaseAssetID,
		&asset.QuoteAssetID,
		&asset.Description,
		&asset.AssetTypeID,
		&asset.CreatedBy,
		&asset.CreatedAt,
		&asset.UpdatedBy,
		&asset.UpdatedAt,
		&asset.ChainID,
		&asset.CategoryID,
		&asset.SubCategoryID,
		&asset.IsDefaultQuote,
		&asset.IgnoreMarketData,
		&asset.Decimals,
		&asset.ContractAddress,
		&asset.StartingBlockNumber,
		&asset.ImportGeth,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return asset, nil
}

// GetAssetByContractAddress : get asset by contract address
func GetAssetByContractAddress(contractAddress string) (*Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	import_geth
	FROM assets 
	WHERE contract_address = $1`, contractAddress)

	asset := &Asset{}
	err := row.Scan(
		&asset.ID,
		&asset.UUID,
		&asset.Name,
		&asset.AlternateName,
		&asset.Cusip,
		&asset.Ticker,
		&asset.BaseAssetID,
		&asset.QuoteAssetID,
		&asset.Description,
		&asset.AssetTypeID,
		&asset.CreatedBy,
		&asset.CreatedAt,
		&asset.UpdatedBy,
		&asset.UpdatedAt,
		&asset.ChainID,
		&asset.CategoryID,
		&asset.SubCategoryID,
		&asset.IsDefaultQuote,
		&asset.IgnoreMarketData,
		&asset.Decimals,
		&asset.ContractAddress,
		&asset.StartingBlockNumber,
		&asset.ImportGeth,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return asset, nil
}

// GetAssetByCusip : get asset by cusip
func GetAssetByCusip(cusip string) (*Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	import_geth
	FROM assets 
	WHERE cusip = $1`, cusip)

	asset := &Asset{}
	err := row.Scan(
		&asset.ID,
		&asset.UUID,
		&asset.Name,
		&asset.AlternateName,
		&asset.Cusip,
		&asset.Ticker,
		&asset.BaseAssetID,
		&asset.QuoteAssetID,
		&asset.Description,
		&asset.AssetTypeID,
		&asset.CreatedBy,
		&asset.CreatedAt,
		&asset.UpdatedBy,
		&asset.UpdatedAt,
		&asset.ChainID,
		&asset.CategoryID,
		&asset.SubCategoryID,
		&asset.IsDefaultQuote,
		&asset.IgnoreMarketData,
		&asset.Decimals,
		&asset.ContractAddress,
		&asset.StartingBlockNumber,
		&asset.ImportGeth,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return asset, nil
}

// GetAssetByBaseAndQuoteID : get asset by base and quote id
func GetAssetByBaseAndQuoteID(baseAssetID *int, quoteAssetID *int) (*Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
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
	import_geth
	FROM assets 
	WHERE base_asset_id = $1 AND 
	quote_asset_id = $2`, *baseAssetID, *quoteAssetID)

	asset := &Asset{}
	err := row.Scan(
		&asset.ID,
		&asset.UUID,
		&asset.Name,
		&asset.AlternateName,
		&asset.Cusip,
		&asset.Ticker,
		&asset.BaseAssetID,
		&asset.QuoteAssetID,
		&asset.Description,
		&asset.AssetTypeID,
		&asset.CreatedBy,
		&asset.CreatedAt,
		&asset.UpdatedBy,
		&asset.UpdatedAt,
		&asset.ChainID,
		&asset.CategoryID,
		&asset.SubCategoryID,
		&asset.IsDefaultQuote,
		&asset.IgnoreMarketData,
		&asset.Decimals,
		&asset.ContractAddress,
		&asset.StartingBlockNumber,
		&asset.ImportGeth,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return asset, nil
}

func GetTopTenAssets() ([]Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	import_geth
	FROM assets 
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assets := make([]Asset, 0)
	for results.Next() {
		var asset Asset
		results.Scan(
			&asset.ID,
			&asset.UUID,
			&asset.Name,
			&asset.AlternateName,
			&asset.Cusip,
			&asset.Ticker,
			&asset.BaseAssetID,
			&asset.QuoteAssetID,
			&asset.Description,
			&asset.AssetTypeID,
			&asset.CreatedBy,
			&asset.CreatedAt,
			&asset.UpdatedBy,
			&asset.UpdatedAt,
			&asset.ChainID,
			&asset.CategoryID,
			&asset.SubCategoryID,
			&asset.IsDefaultQuote,
			&asset.IgnoreMarketData,
			&asset.Decimals,
			&asset.ContractAddress,
			&asset.StartingBlockNumber,
			&asset.ImportGeth,
		)

		assets = append(assets, asset)
	}
	return assets, nil
}

func GetGethImportAssets() ([]Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	import_geth
	FROM assets 
	WHERE import_geth = TRUE
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assets := make([]Asset, 0)
	for results.Next() {
		var asset Asset
		results.Scan(
			&asset.ID,
			&asset.UUID,
			&asset.Name,
			&asset.AlternateName,
			&asset.Cusip,
			&asset.Ticker,
			&asset.BaseAssetID,
			&asset.QuoteAssetID,
			&asset.Description,
			&asset.AssetTypeID,
			&asset.CreatedBy,
			&asset.CreatedAt,
			&asset.UpdatedBy,
			&asset.UpdatedAt,
			&asset.ChainID,
			&asset.CategoryID,
			&asset.SubCategoryID,
			&asset.IsDefaultQuote,
			&asset.IgnoreMarketData,
			&asset.Decimals,
			&asset.ContractAddress,
			&asset.StartingBlockNumber,
			&asset.ImportGeth,
		)

		assets = append(assets, asset)
	}
	return assets, nil
}

func RemoveAsset(assetID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, err := database.DbConnPgx.Query(ctx, `DELETE FROM assets WHERE id = $1`, assetID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetCurrentTradingAssets() ([]Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	&asset.UpdatedAt,
	FROM public.get_current_assets`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assets := make([]Asset, 0)
	for results.Next() {
		var asset Asset
		results.Scan(
			&asset.ID,
			&asset.UUID,
			&asset.Name,
			&asset.AlternateName,
			&asset.Cusip,
			&asset.Ticker,
			&asset.BaseAssetID,
			&asset.QuoteAssetID,
			&asset.Description,
			&asset.AssetTypeID,
			&asset.CreatedBy,
			&asset.CreatedAt,
			&asset.UpdatedBy,
			&asset.UpdatedAt,
			&asset.ChainID,
		)

		assets = append(assets, asset)
	}
	return assets, nil

}

func GetCryptoAssets() ([]Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	import_geth
	FROM public.assets
	where asset_type_id = 1
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assets := make([]Asset, 0)
	for results.Next() {
		var asset Asset
		results.Scan(
			&asset.ID,
			&asset.UUID,
			&asset.Name,
			&asset.AlternateName,
			&asset.Cusip,
			&asset.Ticker,
			&asset.BaseAssetID,
			&asset.QuoteAssetID,
			&asset.Description,
			&asset.AssetTypeID,
			&asset.CreatedBy,
			&asset.CreatedAt,
			&asset.UpdatedBy,
			&asset.UpdatedAt,
			&asset.ChainID,
			&asset.CategoryID,
			&asset.SubCategoryID,
			&asset.IsDefaultQuote,
			&asset.IgnoreMarketData,
			&asset.Decimals,
			&asset.ContractAddress,
			&asset.StartingBlockNumber,
			&asset.ImportGeth,
		)

		assets = append(assets, asset)
	}
	return assets, nil
}

func GetAssetsByAssetTypeAndSource(assetTypeID *int, sourceID *int, excludeIgnoreMarketData bool) ([]AssetWithSources, error) {
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
	assetSources.source_id,
	assetSources.source_identifier,
	assets.decimals,
	contract_address,
	starting_block_number,
	import_geth
	FROM assets assets
	JOIN asset_sources assetSources ON assets.id = assetSources.asset_id
	WHERE assets.asset_type_id = $1
	AND assetSources.source_id = $2
		`
	if excludeIgnoreMarketData {
		sql += `AND ignore_market_data = FALSE
		`
	}
	results, err := database.DbConnPgx.Query(ctx, sql, *assetTypeID, *sourceID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assets := make([]AssetWithSources, 0)
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
			&asset.SourceID,
			&asset.SourceIdentifier,
			&asset.Asset.Decimals,
			&asset.Asset.ContractAddress,
			&asset.StartingBlockNumber,
			&asset.ImportGeth,
		)

		assets = append(assets, asset)
	}
	return assets, nil
}

func GetCryptoAssetsBySourceID(sourceID *int, excludeIgnoreMarketData bool) ([]Asset, error) {
	assetsWithSources, err := GetCryptoAssetsBySourceId(sourceID, excludeIgnoreMarketData)
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

func GetCryptoAssetsBySourceId(sourceID *int, excludeIgnoreMarketData bool) ([]AssetWithSources, error) {
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
	assetSources.source_id,
	assetSources.source_identifier,
	assets.decimals,
	contract_address,
	starting_block_number,
	import_geth
	FROM assets assets
	JOIN asset_sources assetSources ON assets.id = assetSources.asset_id
	WHERE assetSources.source_id = $1
	AND assets.asset_type_id = 1
	
		`
	if excludeIgnoreMarketData {
		sql += `AND ignore_market_data = FALSE
		`
	}
	results, err := database.DbConnPgx.Query(ctx, sql, *sourceID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assets := make([]AssetWithSources, 0)
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
			&asset.SourceID,
			&asset.SourceIdentifier,
			&asset.Asset.Decimals,
			&asset.Asset.ContractAddress,
			&asset.StartingBlockNumber,
			&asset.ImportGeth,
		)

		assets = append(assets, asset)
	}
	return assets, nil

}

func GetAssetWithSourceByAssetIdAndSourceID(assetID *int, sourceID *int, excludeIgnoreMarketData bool) (*AssetWithSources, error) {
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
	assetSources.source_id,
	assetSources.source_identifier,
	assets.decimals,
	contract_address,
	starting_block_number,
	import_geth
	FROM assets assets
	JOIN asset_sources assetSources ON assets.id = assetSources.asset_id
	WHERE 
		assets.id = $1
	AND assetSources.source_id = $2
	
		`
	if excludeIgnoreMarketData {
		query += `AND ignore_market_data = FALSE
		`
	}
	row := database.DbConn.QueryRowContext(ctx, query, &assetID, &sourceID)
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
		&assetWithSources.SourceID,
		&assetWithSources.SourceIdentifier,
		&assetWithSources.Asset.Decimals,
		&assetWithSources.Asset.ContractAddress,
		&assetWithSources.Asset.StartingBlockNumber,
		&assetWithSources.Asset.ImportGeth,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return assetWithSources, nil
}

func GetAssetWithSourceByAssetIdsAndSourceID(assetIDs []int, sourceID *int, excludeIgnoreMarketData bool) ([]AssetWithSources, error) {
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
	assetSources.source_id,
	assetSources.source_identifier,
	assets.decimals,
	contract_address,
	starting_block_number,
	import_geth
	FROM assets assets
	JOIN asset_sources assetSources ON assets.id = assetSources.asset_id
	WHERE 
	assets.id = ANY($1)
	AND assetSources.source_id = $2
	
		`
	if excludeIgnoreMarketData {
		query += `AND ignore_market_data = FALSE
		`
	}
	results, err := database.DbConnPgx.Query(ctx, query, pq.Array(assetIDs), &sourceID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assets := make([]AssetWithSources, 0)
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
			&asset.SourceID,
			&asset.SourceIdentifier,
			&asset.Asset.Decimals,
			&asset.Asset.ContractAddress,
			&asset.Asset.StartingBlockNumber,
			&asset.Asset.ImportGeth,
		)

		assets = append(assets, asset)
	}
	return assets, nil
}

func GetAssetList(ids []int) ([]Asset, error) {
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
	import_geth
	FROM assets`
	if len(ids) > 0 {
		strIds := utils.SplitToString(ids, ",")
		additionalQuery := fmt.Sprintf(` WHERE id IN (%s)`, strIds)
		sql += additionalQuery
	}
	results, err := database.DbConnPgx.Query(ctx, sql)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assets := make([]Asset, 0)
	for results.Next() {
		var asset Asset
		results.Scan(
			&asset.ID,
			&asset.UUID,
			&asset.Name,
			&asset.AlternateName,
			&asset.Cusip,
			&asset.Ticker,
			&asset.BaseAssetID,
			&asset.QuoteAssetID,
			&asset.Description,
			&asset.AssetTypeID,
			&asset.CreatedBy,
			&asset.CreatedAt,
			&asset.UpdatedBy,
			&asset.UpdatedAt,
			&asset.ChainID,
			&asset.CategoryID,
			&asset.SubCategoryID,
			&asset.IsDefaultQuote,
			&asset.IgnoreMarketData,
			&asset.Decimals,
			&asset.ContractAddress,
			&asset.StartingBlockNumber,
			&asset.ImportGeth,
		)

		assets = append(assets, asset)
	}
	return assets, nil
}
func GetDefaultQuoteAssetListBySourceID(sourceID *int) ([]AssetWithSources, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	assetSources.source_id,
	assetSources.source_identifier,
	assets.decimals,
	contract_address,
	starting_block_number,
	import_geth
	FROM get_default_quotes assets
	JOIN asset_sources assetSources ON assets.id = assetSources.asset_id
	WHERE assetSources.source_id = $1
	`, *sourceID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assets := make([]AssetWithSources, 0)
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
			&asset.SourceID,
			&asset.SourceIdentifier,
			&asset.Asset.Decimals,
			&asset.Asset.ContractAddress,
			&asset.StartingBlockNumber,
			&asset.ImportGeth,
		)

		assets = append(assets, asset)
	}
	return assets, nil

}
func GetDefaultQuoteAssetList() ([]Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	results, err := database.DbConnPgx.Query(ctx, `SELECT 
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
	import_geth
	FROM get_default_quotes`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	assets := make([]Asset, 0)
	for results.Next() {
		var asset Asset
		results.Scan(
			&asset.ID,
			&asset.UUID,
			&asset.Name,
			&asset.AlternateName,
			&asset.Cusip,
			&asset.Ticker,
			&asset.BaseAssetID,
			&asset.QuoteAssetID,
			&asset.Description,
			&asset.AssetTypeID,
			&asset.CreatedBy,
			&asset.CreatedAt,
			&asset.UpdatedBy,
			&asset.UpdatedAt,
			&asset.ChainID,
			&asset.CategoryID,
			&asset.SubCategoryID,
			&asset.IsDefaultQuote,
			&asset.IgnoreMarketData,
			&asset.Decimals,
			&asset.ContractAddress,
			&asset.StartingBlockNumber,
			&asset.ImportGeth,
		)

		assets = append(assets, asset)
	}
	return assets, nil
}

func UpdateAsset(asset Asset) error {
	// if the asset id is set, update, otherwise add
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	if asset.ID == nil || *asset.ID == 0 {
		return errors.New("asset has invalid ID")
	}
	_, err := database.DbConnPgx.Query(ctx, `UPDATE assets SET 
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
		import_geth = $18
		WHERE id=$19`,
		asset.Name,                // 1
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
		asset.ID)                  //19

	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func InsertAsset(asset Asset) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	var insertID int
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO assets  
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
		import_geth
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
			$18
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
	).Scan(&insertID)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
