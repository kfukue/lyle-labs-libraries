

COMMIT
BEGIN TRANSACTION;

DROP TABLE IF EXISTS geth_transaction_inputs CASCADE;
CREATE TABLE geth_transaction_inputs
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  function_name VARCHAR(255) NULL,
  method_id_str VARCHAR(255) NULL,
  num_of_parameters INT NOT NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id)
);
DROP TABLE IF EXISTS geth_transactions CASCADE;
CREATE TABLE geth_transactions
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  chain_id INT NOT NULL,
  exchange_id INT NULL,
  block_number NUMERIC NULL,
  index_number NUMERIC NULL,
  txn_date timestamp NULL,
  txn_hash VARCHAR(255) NOT NULL,
  from_address VARCHAR(255) NOT NULL,
  from_address_id INT NULL,
  to_address VARCHAR(255) NOT NULL,
  to_address_id INT NULL,
  interacted_contract_address VARCHAR(255) NULL,
  interacted_contract_address_id INT NULL,
  native_asset_id INT NOT NULL,
  geth_process_job_id INT NULL,
  value NUMERIC NULL,
  geth_transaction_input_id INT NULL,
  status_id INT NOT NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_chain FOREIGN KEY(chain_id) REFERENCES chains(id),
  CONSTRAINT fk_exchange FOREIGN KEY(exchange_id) REFERENCES exchanges(id),
  CONSTRAINT fk_from_address FOREIGN KEY(from_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_to_address FOREIGN KEY(to_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_interacted_contract_address FOREIGN KEY(interacted_contract_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_geth_process_jobs FOREIGN KEY(geth_process_job_id) REFERENCES geth_process_jobs(id),
  CONSTRAINT fk_geth_transaction_input FOREIGN KEY(geth_transaction_input_id) REFERENCES geth_transaction_inputs(id),
  CONSTRAINT fk_statuses FOREIGN KEY(status_id) REFERENCES structured_values(id)
);

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-user";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-user";

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-api";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-api";
