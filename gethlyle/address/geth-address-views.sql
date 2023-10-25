CREATE VIEW geth_addresses
 SELECT assets.id,
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
    assets.is_default_quote,
    assets.ignore_market_data,
    assets.decimals,
    assets.contract_address,
    assets.starting_block_number
   FROM assets
  WHERE (assets.id IN ( SELECT DISTINCT assets_1.id
           FROM trades
             LEFT JOIN assets assets_1 ON trades.asset_id = assets_1.id
          WHERE trades.is_active = true AND assets_1.base_asset_id IS NULL AND assets_1.quote_asset_id IS NULL)) OR (assets.id IN ( SELECT DISTINCT assets_1.base_asset_id
           FROM trades
             LEFT JOIN assets assets_1 ON trades.asset_id = assets_1.id
          WHERE trades.is_active = true AND assets_1.base_asset_id IS NOT NULL)) OR (assets.id IN ( SELECT DISTINCT assets_1.quote_asset_id
           FROM trades
             LEFT JOIN assets assets_1 ON trades.asset_id = assets_1.id
          WHERE trades.is_active = true AND assets_1.quote_asset_id IS NOT NULL));