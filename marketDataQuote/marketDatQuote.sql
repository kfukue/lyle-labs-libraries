CREATE TABLE market_data_quotes
(
  market_data_id INT NOT NULL,
  base_asset_id  INT NOT NULL,
  quote_asset_id INT NOT NULL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4() ,
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  open NUMERIC NULL,
  close NUMERIC NULL,
  high_24h NUMERIC NULL,
  low_24h NUMERIC NULL,
  price NUMERIC NULL,
  volume NUMERIC NULL,
  market_cap NUMERIC NULL,
  ticker VARCHAR(255) NULL,
  description TEXT NULL,
  source_id INT NOT NULL,
  -- new columns 2022-03-06
  fully_diluted_valution NUMERIC NULL,
  ath NUMERIC NULL,
  ath_date date,
  atl NUMERIC NULL,
  atl_date date,
  rename to high_24h NUMERIC NULL,
  rename to low_24h NUMERIC NULL,
  price_change_1h NUMERIC NULL,
  price_change_24h NUMERIC NULL,
  price_change_7d NUMERIC NULL,
  price_change_30d  NUMERIC NULL,
  price_change_60d NUMERIC NULL,
  price_change_200d NUMERIC NULL,
  price_change_1y NUMERIC NULL,
  -- end 2022-03-06

  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(market_data_id, base_asset_id, quote_asset_id),
  CONSTRAINT fk_market_data FOREIGN KEY(market_data_id) REFERENCES market_data(id),
  CONSTRAINT fk_base_asset_id FOREIGN KEY(base_asset_id) REFERENCES assets(id),
  CONSTRAINT fk_quote_asset_id FOREIGN KEY(quote_asset_id) REFERENCES assets(id),
  CONSTRAINT fk_source FOREIGN KEY(source_id) REFERENCES sources(id)
);

-- new columns 2022-03-06
ROLLBACK
START TRANSACTION;
ALTER TABLE market_data_quotes
  ADD COLUMN fully_diluted_valution NUMERIC NULL,
  ADD COLUMN ath NUMERIC NULL,
  ADD COLUMN ath_date date,
  ADD COLUMN atl NUMERIC NULL,
  ADD COLUMN atl_date date,
  ADD COLUMN price_change_1h NUMERIC NULL,
  ADD COLUMN price_change_24h NUMERIC NULL,
  ADD COLUMN price_change_7d NUMERIC NULL,
  ADD COLUMN price_change_30d  NUMERIC NULL,
  ADD COLUMN price_change_60d NUMERIC NULL,
  ADD COLUMN price_change_200d NUMERIC NULL,
  ADD COLUMN price_change_1y NUMERIC NULL

  COMMIT
-- end 2022-03-06
