CREATE TABLE step_assets
(
  id SERIAL,
  step_id INT NOT NULL,
  asset_id  INT NOT NULL,
  swap_asset_id INT NULL,
  target_pool_id INT NULL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4() ,
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  start_date timestamp NULL,
  end_date timestamp NULL,
  description TEXT NULL,
  action_parameter NUMERIC NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_asset_id FOREIGN KEY(asset_id) REFERENCES assets(id),
  CONSTRAINT fk_step_id FOREIGN KEY(step_id) REFERENCES steps(id),
  CONSTRAINT fk_swap_asset_id FOREIGN KEY(asset_id) REFERENCES assets(id),
  CONSTRAINT fk_pool_id FOREIGN KEY(target_pool_id) REFERENCES pools(id)
);
