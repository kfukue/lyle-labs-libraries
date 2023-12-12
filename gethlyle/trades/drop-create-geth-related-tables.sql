BEGIN TRANSACTION;
DROP TABLE IF EXISTS geth_trade_swaps CASCADE;
DROP TABLE IF EXISTS geth_trade_transfers CASCADE;
DROP TABLE IF EXISTS geth_trades CASCADE;
DROP TABLE IF EXISTS geth_addresses CASCADE;
DROP TABLE IF EXISTS geth_process_jobs CASCADE;
DROP TABLE IF EXISTS geth_process_job_topics;
DROP TABLE IF EXISTS geth_swaps CASCADE;
DROP TABLE IF EXISTS geth_transfers CASCADE;

CREATE TABLE geth_addresses
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  description TEXT NULL,
  address_str VARCHAR(255) UNIQUE NOT NULL,
  address_type_id INT NOT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_address_type FOREIGN KEY(address_type_id) REFERENCES structured_values(id)
);

-- create index
CREATE INDEX geth_addresses_address_str ON geth_addresses(address_str);

CREATE TABLE geth_process_jobs
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  start_date timestamp NOT NULL,
  end_date timestamp  NULL,
  description TEXT NULL,
  status_id int NOT NULL,
  job_category_id INT NULL,
  import_type_id INT NULL,
  chain_id INT NULL,
  start_block_number INT  NULL,
  end_block_number INT  NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  asset_id INT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_statuses FOREIGN KEY(status_id) REFERENCES structured_values(id),
  CONSTRAINT fk_job_categories FOREIGN KEY(job_category_id) REFERENCES structured_values(id),
  CONSTRAINT fk_import_type FOREIGN KEY(import_type_id) REFERENCES structured_values(id),
  CONSTRAINT fk_chains FOREIGN KEY(chain_id) REFERENCES chains(id)
);

-- create index
CREATE INDEX geth_process_jobs_asset_id ON geth_process_jobs(asset_id);
CREATE INDEX geth_process_jobs_start_date ON geth_process_jobs(start_date);
CREATE INDEX geth_process_jobs_end_date ON geth_process_jobs(end_date);
CREATE INDEX geth_process_jobs_status_id ON geth_process_jobs(status_id);
CREATE INDEX geth_process_jobs_import_type_id ON geth_process_jobs(import_type_id);

CREATE TABLE geth_process_job_topics
(
  id SERIAL,
  geth_process_job_id INT NULL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  description TEXT NULL,
  status_id int NOT NULL,
  topic_str TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_geth_process_jobs FOREIGN KEY(geth_process_job_id) REFERENCES geth_process_jobs(id),
  CONSTRAINT fk_statuses FOREIGN KEY(status_id) REFERENCES structured_values(id)
);

-- create index
CREATE INDEX geth_process_job_topics_geth_process_job_id ON geth_process_job_topics(geth_process_job_id);
CREATE INDEX geth_process_job_topics_topic_str ON geth_process_job_topics(topic_str);


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


CREATE TABLE geth_transfers
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  chain_id INT NOT NULL,
  token_address VARCHAR(255) NULL,
  token_address_id INT NULL,
  asset_id INT NULL,
  block_number NUMERIC NULL,
  index_number NUMERIC NULL,
  transfer_date timestamp NULL,
  txn_hash VARCHAR(70) NOT NULL,
  sender_address VARCHAR(255) NOT NULL,
  sender_address_id INT NULL,
  to_address VARCHAR(255) NOT NULL,
  to_address_id INT  NULL,
  amount NUMERIC NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  geth_process_job_id INT NULL,
  topics_str TEXT[] NULL,
  status_id INT NOT NULL,
  base_asset_id INT NOT NULL,
  transfer_type_id INT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_chains FOREIGN KEY(chain_id) REFERENCES chains(id),
  CONSTRAINT fk_assets FOREIGN KEY(asset_id) REFERENCES assets(id),
  CONSTRAINT fk_token_address FOREIGN KEY(token_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_sender_address FOREIGN KEY(sender_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_to_address FOREIGN KEY(to_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_geth_process_jobs FOREIGN KEY(geth_process_job_id) REFERENCES geth_process_jobs(id),
  CONSTRAINT fk_statuses FOREIGN KEY(status_id) REFERENCES structured_values(id),
  CONSTRAINT fk_base_asset FOREIGN KEY(base_asset_id) REFERENCES assets(id),
  CONSTRAINT fk_transfer_types FOREIGN KEY(transfer_type_id) REFERENCES structured_values(id)
);
--create index
CREATE INDEX geth_transfers_to_address_id ON geth_transfers(to_address_id);
CREATE INDEX geth_transfers_sender_address_id ON geth_transfers(sender_address_id);
CREATE INDEX geth_transfers_to_address ON geth_transfers(to_address);
CREATE INDEX geth_transfers_sender_address ON geth_transfers(sender_address);
CREATE INDEX geth_transfers_asset_id ON geth_transfers(asset_id);
CREATE INDEX geth_transfers_base_asset_id ON geth_transfers(base_asset_id);

CREATE INDEX  geth_transfers_txn_hash ON geth_transfers(txn_hash);
CREATE INDEX  geth_transfers_fk_geth_process_jobs ON geth_transfers(geth_process_job_id);
CREATE INDEX  geth_transfers_fk_transfer_types ON geth_transfers(transfer_type_id);


CREATE TABLE geth_trades
(
  id                            SERIAL,
  uuid                          uuid NOT NULL DEFAULT uuid_generate_v4(),
  name                          VARCHAR(255) NOT NULL,
  alternate_name                VARCHAR(255) NULL,
  address_str                   VARCHAR(255) NOT NULL,
  address_id                    INT NOT NULL,
  trade_date                    timestamp NULL,
  txn_hash                      VARCHAR(255) NULL,
  is_buy                        BOOLEAN NULL,
  token0_amount                 NUMERIC NULL,
  token0_amount_decimal_adj     NUMERIC NULL,
  token1_amount                 NUMERIC NULL,
  token1_amount_decimal_adj     NUMERIC NULL,
  price                         NUMERIC NULL,
  price_usd                     NUMERIC NULL,
  lp_token1_price_usd           NUMERIC NULL,
  total_amount_usd              NUMERIC NULL,
  token0_asset_id               INT NOT NULL,
  token1_asset_id               INT NULL,
  geth_process_job_id           INT NULL,
  status_id                     INT NOT NULL,
  trade_type_id                 INT NULL,
  description                   TEXT NULL,
  created_by                    VARCHAR(255) NOT NULL,
  created_at                    timestamp NOT NULL,
  updated_by                    VARCHAR(255) NOT NULL,
  updated_at                    timestamp NOT NULL,
  base_asset_id                 INT NOT NULL,
  oracle_price_usd              NUMERIC NULL,
  oracle_price_asset_id         INT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_address FOREIGN KEY(address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_token0_asset FOREIGN KEY(token0_asset_id) REFERENCES assets(id),
  CONSTRAINT fk_token1_asset FOREIGN KEY(token1_asset_id) REFERENCES assets(id),
  CONSTRAINT fk_geth_process_jobs FOREIGN KEY(geth_process_job_id) REFERENCES geth_process_jobs(id),
  CONSTRAINT fk_status FOREIGN KEY(status_id) REFERENCES structured_values(id),
  CONSTRAINT fk_trade_type FOREIGN KEY(trade_type_id) REFERENCES structured_values(id),
  CONSTRAINT fk_base_asset FOREIGN KEY(base_asset_id) REFERENCES assets(id),
  CONSTRAINT fk_oracle_price_asset FOREIGN KEY(oracle_price_asset_id) REFERENCES assets(id)
);

--create index
CREATE INDEX geth_trades_address_str ON geth_trades(address_str);
CREATE INDEX geth_trades_address_id ON geth_trades(address_id);
CREATE INDEX geth_trades_txn_hash ON geth_trades(txn_hash);
CREATE INDEX geth_trades_fk_geth_process_jobs ON geth_trades(geth_process_job_id);
CREATE INDEX geth_trades_fk_statuses ON geth_trades(status_id);
CREATE INDEX geth_trades_fk_base_asset ON geth_trades(base_asset_id);
CREATE INDEX geth_trades_fk_oracle_price_asset ON geth_trades(oracle_price_asset_id);

CREATE TABLE geth_trade_swaps
(
  geth_trade_id                 INT NOT NULL,
  geth_swap_id                  INT NOT NULL,
  uuid                          uuid NOT NULL DEFAULT uuid_generate_v4(),
  name                          VARCHAR(255) NOT NULL,
  alternate_name                VARCHAR(255) NULL,
  description                   TEXT NULL,
  created_by                    VARCHAR(255) NOT NULL,
  created_at                    timestamp NOT NULL,
  updated_by                    VARCHAR(255) NOT NULL,
  updated_at                    timestamp NOT NULL,
  PRIMARY KEY(geth_trade_id, geth_swap_id),
  CONSTRAINT fk_geth_trade FOREIGN KEY(geth_trade_id) REFERENCES geth_trades(id),
  CONSTRAINT fk_geth_swap FOREIGN KEY(geth_swap_id) REFERENCES geth_swaps(id)
);
-- create index
CREATE INDEX geth_trade_swaps_swap_id ON geth_trade_swaps(geth_swap_id);
CREATE INDEX geth_trade_swaps_trade_id ON geth_trade_swaps(geth_trade_id);


CREATE TABLE geth_trade_transfers
(
  geth_trade_id                 INT NOT NULL,
  geth_transfer_id              INT NOT NULL,
  tax_id                        INT NULL,
  uuid                          uuid NOT NULL DEFAULT uuid_generate_v4(),
  name                          VARCHAR(255) NOT NULL,
  alternate_name                VARCHAR(255) NULL,
  description                   TEXT NULL,
  created_by                    VARCHAR(255) NOT NULL,
  created_at                    timestamp NOT NULL,
  updated_by                    VARCHAR(255) NOT NULL,
  updated_at                    timestamp NOT NULL,
  PRIMARY KEY(geth_trade_id, geth_transfer_id),
  CONSTRAINT fk_geth_trade FOREIGN KEY(geth_trade_id) REFERENCES geth_trades(id),
  CONSTRAINT fk_geth_transfers FOREIGN KEY(geth_transfer_id) REFERENCES geth_transfers(id)
);

-- create index
CREATE INDEX geth_trade_transfers_transfer_id ON geth_trade_transfers(geth_transfer_id);
CREATE INDEX geth_trade_transfers_trade_id ON geth_trade_transfers(geth_trade_id);



GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-user";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-user";

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-api";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-api";
