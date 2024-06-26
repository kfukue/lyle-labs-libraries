CREATE TABLE chains
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  base_asset_id VARCHAR(255) NULL,
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  address VARCHAR(255) NULL,
  chain_type_id INT NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  rpc_url VARCHAR(255) NULL,
  chain_id INT NULL,
  block_explorer_url VARCHAR(255) NULL,
  COLUMN rpc_url_dev VARCHAR(255) NULL,
  rpc_url_prod VARCHAR(255) NULL,
  rpc_url_archive VARCHAR(255) NULL
  PRIMARY KEY(id),
  CONSTRAINT fk_chain_type FOREIGN KEY(chain_type_id) REFERENCES structured_values(id)
);

-- new columns 2022-08-04
ROLLBACK
START TRANSACTION;
ALTER TABLE chains
  ADD COLUMN base_asset_id VARCHAR(255) NULL,
  ADD COLUMN rpc_url VARCHAR(255) NULL,
  ADD COLUMN chain_id INT NULL,
  ADD COLUMN block_explorer_url VARCHAR(255) NULL,
  CONSTRAINT fk_base_asset FOREIGN KEY(base_asset_id) REFERENCES assets(id)
  COMMIT
-- end 2022-08-04

-- new columns 2023-09-11
ROLLBACK
START TRANSACTION;
ALTER TABLE chains
  ADD COLUMN rpc_url_dev VARCHAR(255) NULL,
  ADD COLUMN rpc_url_prod VARCHAR(255) NULL
  COMMIT
-- end 2023-09-11


-- new columns 2024-05-14
ROLLBACK
START TRANSACTION;
ALTER TABLE chains
  ADD COLUMN rpc_url_archive VARCHAR(255) NULL
  COMMIT
-- end 2024-05-14
