CREATE TABLE transactions
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4() ,
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  start_date timestamp NOT NULL,
  end_date timestamp NOT NULL,
  description TEXT NULL,
  tx_hash VARCHAR(255) NULL,
  status_id int NOT NULL,
  from_account_id INT NULL,
  to_account_id INT NULL,
  chain_id INT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_status_id FOREIGN KEY(status_id) REFERENCES structured_values(id),
  CONSTRAINT fk_from_account_id FOREIGN KEY(from_account_id) REFERENCES accounts(id),
  CONSTRAINT fk_to_account_id FOREIGN KEY(to_account_id) REFERENCES accounts(id),
  CONSTRAINT fk_chain_id FOREIGN KEY(chain_id) REFERENCES chains(id)
);
