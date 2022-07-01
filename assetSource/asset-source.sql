CREATE TABLE asset_sources
(
  source_id INT NOT NULL,
  asset_id INT NOT NULL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  source_identifier VARCHAR(255) NOT NULL,
  description TEXT NULL,
  source_data jsonb NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(source_id, asset_id),
  CONSTRAINT fk_source FOREIGN KEY(source_id) REFERENCES sources(id),
  CONSTRAINT fk_asset FOREIGN KEY(asset_id) REFERENCES assets(id)
);
