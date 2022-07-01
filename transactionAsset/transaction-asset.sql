CREATE TABLE transaction_assets
(
  transaction_id INT NOT NULL,
  asset_id INT NOT NULL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  description TEXT NULL,
  quantity NUMERIC NULL,
  quantity_usd  NUMERIC NULL,
  market_data_id INT NULL,
  manual_exchange_rate_usd NUMERIC NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(transaction_id, asset_id),
  CONSTRAINT fk_transaction FOREIGN KEY(transaction_id) REFERENCES transactions(id),
  CONSTRAINT fk_asset FOREIGN KEY(asset_id) REFERENCES assets(id),
  CONSTRAINT fk_market_data FOREIGN KEY(market_data_id) REFERENCES market_data(id)
);
