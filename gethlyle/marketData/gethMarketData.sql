COMMIT
BEGIN TRANSACTION;
DROP TABLE IF EXISTS geth_market_data CASCADE;

CREATE TABLE geth_market_data
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4() ,
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  start_date timestamp NOT NULL,
  end_date timestamp NOT NULL,
  asset_id INT NOT NULL,
  open_usd NUMERIC NULL,
  close_usd NUMERIC NULL,
  high_usd NUMERIC NULL,
  low_usd NUMERIC NULL,
  price_usd NUMERIC NULL,
  volume_usd NUMERIC NULL,
  market_cap_usd NUMERIC NULL,
  ticker VARCHAR(255) NULL,
  description TEXT NULL,
  interval_id INT NOT NULL,
  market_data_type_id INT NOT NULL,
  source_id INT NOT NULL,
  total_supply NUMERIC NULL,
  max_supply NUMERIC NULL,
  circulating_supply NUMERIC NULL,
  sparkline_7d NUMERIC [] NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  geth_process_job_id INT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_asset FOREIGN KEY(asset_id) REFERENCES assets(id),
  CONSTRAINT fk_structured_value_interval FOREIGN KEY(interval_id) REFERENCES structured_values(id),
  CONSTRAINT fk_source FOREIGN KEY(source_id) REFERENCES sources(id),
  CONSTRAINT fk_structured_value_market_data_type FOREIGN KEY(market_data_type_id) REFERENCES structured_values(id),
  CONSTRAINT fk_geth_process_jobs FOREIGN KEY(geth_process_job_id) REFERENCES geth_process_jobs(id)
);

CREATE INDEX geth_market_data_asset_id ON geth_market_data(id);
CREATE INDEX geth_market_data_market_data_type_id ON geth_market_data(market_data_type_id);
CREATE INDEX geth_market_data_start_date ON geth_market_data(start_date);

CREATE INDEX geth_market_data_interval ON geth_market_data(interval_id);
CREATE INDEX geth_market_data_source ON geth_market_data(source_id);
CREATE INDEX geth_market_data_geth_process_job ON geth_market_data(geth_process_job_id);

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-user";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-user";

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-api";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-api";

