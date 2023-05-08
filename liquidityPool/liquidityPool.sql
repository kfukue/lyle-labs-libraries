CREATE TABLE liquidity_pools
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4() ,
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  pair_address VARCHAR(255) NULL,
  chain_id INT NULL,
  exchange_id INT NULL,
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
  PRIMARY KEY(id),
  CONSTRAINT fk_exchanged_id FOREIGN KEY(exchange_id) REFERENCES exchanges(id),
  CONSTRAINT fk_token0_id FOREIGN KEY(token0_id) REFERENCES assets(id),
  CONSTRAINT fk_token1_id FOREIGN KEY(token1_id) REFERENCES assets(id)

);

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
