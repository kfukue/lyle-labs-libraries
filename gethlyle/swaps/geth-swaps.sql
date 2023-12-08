COMMIT
BEGIN TRANSACTION;
DROP TABLE IF EXISTS geth_swaps CASCADE;

CREATE TABLE geth_swaps
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  chain_id INT NOT NULL,
  exchange_id INT NULL,
  block_number NUMERIC NULL,
  index_number NUMERIC NULL,
  swap_date timestamp NULL,
  trade_type_id INT NULL,
  txn_hash VARCHAR(255) NOT NULL,
  maker_address VARCHAR(255) NOT NULL,
  maker_address_id INT NULL,
  is_buy BOOLEAN NULL,
  price         NUMERIC NULL,
  price_usd NUMERIC NULL,
  token1_price_usd NUMERIC NULL,
  total_amount_usd NUMERIC NULL,
	pair_address   VARCHAR(255) NULL,
  liquidity_pool_id INT NULL,
	token0_asset_id INT NULL,
	token1_asset_id INT NULL,
	token0_amount  NUMERIC NULL,
	token1_amount  NUMERIC NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  geth_process_job_id INT NULL,
  topics_str TEXT[] NULL,
  status_id INT NOT NULL,
  base_asset_id INT NOT NULL,
  oracle_price_usd NUMERIC NULL,
  oracle_price_asset_id INT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_chain FOREIGN KEY(chain_id) REFERENCES chains(id),
  CONSTRAINT fk_exchange FOREIGN KEY(exchange_id) REFERENCES exchanges(id),
  CONSTRAINT fk_trade_type FOREIGN KEY(trade_type_id) REFERENCES structured_values(id),
  CONSTRAINT fk_liquidity_pool FOREIGN KEY(liquidity_pool_id) REFERENCES liquidity_pools(id),
  CONSTRAINT fk_maker_address FOREIGN KEY(maker_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_geth_process_jobs FOREIGN KEY(geth_process_job_id) REFERENCES geth_process_jobs(id),
  CONSTRAINT fk_statuses FOREIGN KEY(status_id) REFERENCES structured_values(id),
  CONSTRAINT fk_base_asset FOREIGN KEY(base_asset_id) REFERENCES assets(id),
  CONSTRAINT fk_oracle_price_asset FOREIGN KEY(oracle_price_asset_id) REFERENCES assets(id)
);

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-user";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-user";

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-api";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-api";

--create index
CREATE INDEX geth_swaps_maker_address_id ON geth_swaps(maker_address_id);
CREATE INDEX geth_swaps_txn_hash ON geth_swaps(txn_hash);
CREATE INDEX geth_swaps_maker_address ON geth_swaps(maker_address);
CREATE INDEX geth_swaps_token0_asset_id ON geth_swaps(token0_asset_id);
CREATE INDEX geth_swaps_fk_chain ON geth_swaps(chain_id);
CREATE INDEX geth_swaps_fk_exchange ON geth_swaps(exchange_id);
CREATE INDEX geth_swaps_fk_trade_type ON geth_swaps(trade_type_id);
CREATE INDEX geth_swaps_fk_liquidity_pool ON geth_swaps(liquidity_pool_id);
CREATE INDEX geth_swaps_fk_maker_address ON geth_swaps(maker_address_id);
CREATE INDEX geth_swaps_fk_statuses ON geth_swaps(status_id);
CREATE INDEX geth_swaps_fk_base_asset ON geth_swaps(base_asset_id);
CREATE INDEX geth_swaps_fk_oracle_price_asset ON geth_swaps(oracle_price_asset_id);



-- new column 2023-12-08
ROLLBACK
START TRANSACTION;
ALTER TABLE geth_swaps
  ADD  oracle_price_usd NUMERIC NULL,
  ADD  oracle_price_asset_id INT NULL,
  ADD CONSTRAINT fk_oracle_price_asset FOREIGN KEY(oracle_price_asset_id) REFERENCES assets(id)
  COMMIT
-- end


-- update existing rows to oracle price and oracle usd (use token1_price and token1_asset_id)
ROLLBACK
START TRANSACTION;
UPDATE geth_swaps SET
oracle_price_usd = token1_price_usd
oracle_price_asset_id = token1_asset_id
