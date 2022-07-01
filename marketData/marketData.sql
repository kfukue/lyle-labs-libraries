CREATE TABLE market_data
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
  -- new columns 2022-03-06
  total_supply NUMERIC NULL,
  max_supply NUMERIC NULL,
  circulating_supply NUMERIC NULL,
  sparkline_7d NUMERIC [] NULL,
-- end 2022-03-06
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_asset FOREIGN KEY(asset_id) REFERENCES assets(id),
  CONSTRAINT fk_structured_value_interval FOREIGN KEY(interval_id) REFERENCES structured_values(id),
  CONSTRAINT fk_source FOREIGN KEY(source_id) REFERENCES sources(id)
);

-- new columns 2022-03-06
ROLLBACK
START TRANSACTION;
ALTER TABLE market_data
  ADD COLUMN total_supply NUMERIC NULL,
  ADD COLUMN max_supply NUMERIC NULL,
  ADD COLUMN circulating_supply NUMERIC NULL,
  ADD COLUMN sparkline_7d NUMERIC [] NULL
  
  COMMIT
-- end 2022-03-06
