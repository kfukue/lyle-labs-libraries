SELECT 
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
	assets.chainlink_usd_address,
	assets.chainlink_usd_chain_id,
	assetSources.source_id,
	assetSources.source_identifier
	FROM assets assets
	JOIN asset_sources assetSources ON assets.id = assetSources.asset_id
	WHERE assetSources.source_id = 3
	AND assets.asset_type_id = 1
  AND assets.ignore_market_data = FALSE