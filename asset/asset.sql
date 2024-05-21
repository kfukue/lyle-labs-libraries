CREATE TABLE assets
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4() ,
  name TEXT NOT NULL,
  alternate_name TEXT NULL,
  cusip VARCHAR(255) NULL,
  ticker VARCHAR(255) NULL,
  base_asset_id INT NULL,
  quote_asset_id INT NULL,
  description TEXT NULL,
  asset_type_id INT NOT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  category_id INT NULL,
  sub_category_id INT NULL,
  chain_id INT NULL,
  is_default_quote BOOLEAN NOT NULL,
  ignore_market_data BOOLEAN NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_base_asset FOREIGN KEY(base_asset_id) REFERENCES assets(id),
  CONSTRAINT fk_quote_asset FOREIGN KEY(quote_asset_id) REFERENCES assets(id),
  CONSTRAINT fk_structured_value_asset_type FOREIGN KEY(asset_type_id) REFERENCES structured_values(id)
  CONSTRAINT fk_chain FOREIGN KEY REFERENCES chains(id)
  CONSTRAINT fk_structured_value_category FOREIGN KEY(category_id) REFERENCES structured_values(id)
  CONSTRAINT fk_structured_value__sub_category FOREIGN KEY(sub_category_id) REFERENCES structured_values(id)
);

-- new columns 2022-07-02
ROLLBACK
START TRANSACTION;
ALTER TABLE assets
  ADD COLUMN decimals INT NULL,
  ADD COLUMN contract_address VARCHAR(255) NULL
  COMMIT
-- end 2022-07-02
-- new colun 2023-07-22
ROLLBACK
START TRANSACTION;
ALTER TABLE assets
  ADD  starting_block_number NUMERIC NULL
  COMMIT
-- end

-- new colun 2023-08-28
ROLLBACK
START TRANSACTION;
ALTER TABLE assets
  ADD  import_geth BOOLEAN NOT NULL DEFAULT FALSE;
  COMMIT
-- end
-- make name alt name to text 2023-11-14
ROLLBACK
START TRANSACTION;
ALTER TABLE assets
  ALTER COLUMN name TYPE text;
  ALTER COLUMN alternate_name TYPE text;
  
  COMMIT
-- end

-- Add Import Geth Initial As New Column
ROLLBACK
START TRANSACTION;
ALTER TABLE assets
  ADD  import_geth_initial BOOLEAN NOT NULL DEFAULT FALSE;
  COMMIT
-- end
-- add indexes

CREATE INDEX assets_id ON assets(id);
CREATE INDEX assets_address ON assets(contract_address);
CREATE INDEX assets_base_asset ON assets(base_asset_id);
CREATE INDEX assets_quote_asset ON assets(quote_asset_id);
CREATE INDEX assets_chains ON assets(chain_id);
CREATE INDEX assets_asset_type ON assets(asset_type_id);
