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
  PRIMARY KEY(id),
  CONSTRAINT fk_asset FOREIGN KEY(asset_id) REFERENCES assets(id),
  CONSTRAINT fk_structured_value_interval FOREIGN KEY(interval_id) REFERENCES structured_values(id),
  CONSTRAINT fk_source FOREIGN KEY(source_id) REFERENCES sources(id),
  CONSTRAINT fk_structured_value_market_data_type FOREIGN KEY(market_data_type_id) REFERENCES structured_values(id)
);

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-user";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-user";

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-api";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-api";

