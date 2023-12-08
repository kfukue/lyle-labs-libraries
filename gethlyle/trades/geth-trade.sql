ROLLBACK
BEGIN TRANSACTION;
DROP TABLE IF EXISTS geth_trades CASCADE;

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


GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-user";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-user";

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-api";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-api";



-- new column 2023-12-08
ROLLBACK
START TRANSACTION;
ALTER TABLE geth_trades
  ADD  oracle_price_usd NUMERIC NULL,
  ADD  oracle_price_asset_id INT NULL,
  ADD CONSTRAINT fk_oracle_price_asset FOREIGN KEY(oracle_price_asset_id) REFERENCES assets(id)
  COMMIT
-- end