CREATE TABLE accounts
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  address VARCHAR(255) NULL,
  name_from_source VARCHAR(255) NULL,
  portfolio_id INT NULL,
  source_id INT NULL,
  account_type_id INT NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  chain_id INT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_source FOREIGN KEY(source_id) REFERENCES sources(id),
  CONSTRAINT fk_portfolios FOREIGN KEY(portfolio_id) REFERENCES portfolios(id),
  CONSTRAINT fk_account_type FOREIGN KEY(account_type_id) REFERENCES structured_values(id),
  CONSTRAINT fk_chain FOREIGN KEY(chain_id) REFERENCES chains(id)
);

