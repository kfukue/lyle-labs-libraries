BEGIN TRANSACTION;
DROP TABLE IF EXISTS geth_trade_swaps CASCADE;
DROP TABLE IF EXISTS geth_trade_tax_transfers CASCADE;
DROP TABLE IF EXISTS geth_trades CASCADE;
DROP TABLE IF EXISTS geth_addresses CASCADE;
DROP TABLE IF EXISTS geth_process_jobs CASCADE;
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
  PRIMARY KEY(id),
  CONSTRAINT fk_chains FOREIGN KEY(chain_id) REFERENCES chains(id),
  CONSTRAINT fk_assets FOREIGN KEY(asset_id) REFERENCES assets(id),
  CONSTRAINT fk_token_address FOREIGN KEY(token_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_sender_address FOREIGN KEY(sender_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_to_address FOREIGN KEY(to_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_geth_process_jobs FOREIGN KEY(geth_process_job_id) REFERENCES geth_process_jobs(id),
  CONSTRAINT fk_statuses FOREIGN KEY(status_id) REFERENCES structured_values(id)
);

--create index
CREATE INDEX geth_transfers_to_address_id ON geth_transfers(to_address_id);
CREATE INDEX geth_transfers_sender_address_id ON geth_transfers(sender_address_id);
CREATE INDEX geth_transfers_to_address ON geth_transfers(to_address);
CREATE INDEX geth_transfers_sender_address ON geth_transfers(sender_address);
CREATE INDEX geth_transfers_asset_id ON geth_transfers(asset_id);

CREATE TABLE geth_swaps
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  chain_id INT NOT NULL,
  exchange_id INT NOT NULL,
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
  PRIMARY KEY(id),
  CONSTRAINT fk_chain FOREIGN KEY(chain_id) REFERENCES chains(id),
  CONSTRAINT fk_exchange FOREIGN KEY(exchange_id) REFERENCES exchanges(id),
  CONSTRAINT fk_trade_type FOREIGN KEY(trade_type_id) REFERENCES structured_values(id),
  CONSTRAINT fk_liquidity_pool FOREIGN KEY(liquidity_pool_id) REFERENCES liquidity_pools(id),
  CONSTRAINT fk_maker_address FOREIGN KEY(maker_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_geth_process_jobs FOREIGN KEY(geth_process_job_id) REFERENCES geth_process_jobs(id),
  CONSTRAINT fk_statuses FOREIGN KEY(status_id) REFERENCES structured_values(id)
);

--create index
CREATE INDEX geth_swaps_maker_address_id ON geth_swaps(maker_address_id);
CREATE INDEX geth_swaps_token0_asset_id ON geth_swaps(token0_asset_id);


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
  PRIMARY KEY(id),
  CONSTRAINT fk_address FOREIGN KEY(address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_token0_asset FOREIGN KEY(token0_asset_id) REFERENCES assets(id),
  CONSTRAINT fk_token1_asset FOREIGN KEY(token1_asset_id) REFERENCES assets(id),
  CONSTRAINT fk_geth_process_jobs FOREIGN KEY(geth_process_job_id) REFERENCES geth_process_jobs(id),
  CONSTRAINT fk_status FOREIGN KEY(status_id) REFERENCES structured_values(id),
  CONSTRAINT fk_trade_type FOREIGN KEY(trade_type_id) REFERENCES structured_values(id)
);


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



CREATE TABLE geth_trade_tax_transfers
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


GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-user";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-user";

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-api";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-api";