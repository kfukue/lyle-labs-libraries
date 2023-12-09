COMMIT

DROP TABLE liquidity_pools;
DROP TABLE liquidity_pool_assets;
BEGIN TRANSACTION;
CREATE TABLE liquidity_pools
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4() ,
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  pair_address VARCHAR(255) NULL,
  chain_id INT NULL,
  exchange_id INT NULL,
  liquidity_pool_type_id INT NULL,
  token0_id INT NULL,
  token1_id INT NULL,
  url varchar(255)NULL,
  start_block INT NULL,
  latest_block_synced INT NULL,
  created_txn_hash VARCHAR(255) NULL,
  IsActive BOOLEAN NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  base_asset_id INT NULL,
  quote_asset_id INT NULL,
  quote_asset_chainlink_address_usd VARCHAR(255) NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_chain_id FOREIGN KEY(chain_id) REFERENCES chains(id),
  CONSTRAINT fk_exchange_id FOREIGN KEY(exchange_id) REFERENCES exchanges(id),
  CONSTRAINT fk_liquidity_pool_type_id FOREIGN KEY(liquidity_pool_type_id) REFERENCES structured_values(id),
  CONSTRAINT fk_token0_id FOREIGN KEY(token0_id) REFERENCES assets(id),
  CONSTRAINT fk_token1_id FOREIGN KEY(token1_id) REFERENCES assets(id),
  CONSTRAINT fk_base_asset_id FOREIGN KEY(base_asset_id) REFERENCES assets(id),
  CONSTRAINT fk_quote_asset_id FOREIGN KEY(quote_asset_id) REFERENCES assets(id)
);
-- new colun 2023-11-30
ROLLBACK
START TRANSACTION;
ALTER TABLE liquidity_pools
  ADD base_asset_id INT NULL;
ALTER TABLE liquidity_pools
  ADD CONSTRAINT fk_base_asset_id FOREIGN KEY(base_asset_id) REFERENCES assets(id);
COMMIT
-- end

-- new colun 2023-12-08
ROLLBACK
START TRANSACTION;
ALTER TABLE liquidity_pools
  ADD quote_asset_id INT NULL,
  ADD quote_asset_chainlink_address_usd VARCHAR(255) NULL,
  ADD CONSTRAINT fk_quote_asset_id FOREIGN KEY(quote_asset_id) REFERENCES assets(id);
COMMIT;
-- end

CREATE TABLE liquidity_pool_assets
(
  uuid uuid NOT NULL DEFAULT uuid_generate_v4() ,
  liquidity_pool_id INT NOT NULL,
  asset_id  INT NOT NULL,
  token_number INT NOT NULL,
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(liquidity_pool_id, asset_id,token_number),
  CONSTRAINT fk_liquidity_pool_id FOREIGN KEY(liquidity_pool_id) REFERENCES liquidity_pools(id),
  CONSTRAINT fk_asset_id FOREIGN KEY(asset_id) REFERENCES assets(id)
);
