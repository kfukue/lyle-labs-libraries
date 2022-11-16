CREATE TABLE portfolios
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4() ,
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  start_date timestamp NULL,
  end_date timestamp NULL,
  user_email VARCHAR(255) NOT NULL,
  description TEXT NULL,
  base_asset_id INT NOT NULL,
  portfolio_type_id INT NOT NULL,
  parent_id INT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_asset_base_asset_id FOREIGN KEY(base_asset_id) REFERENCES assets(id),
  CONSTRAINT fk_portfolio_parent_id FOREIGN KEY(parent_id) REFERENCES portfolios(id),
  CONSTRAINT fk_structured_value_portfolio_type FOREIGN KEY(portfolio_type_id) REFERENCES structured_values(id)
);
