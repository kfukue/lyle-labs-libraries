COMMIT
BEGIN TRANSACTION;
DROP TABLE IF EXISTS geth_miners CASCADE;

CREATE TABLE geth_miners
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  chain_id INT NOT NULL,
  exchange_id INT NULL,
  starting_block_number NUMERIC NULL,
  created_txn_hash VARCHAR(255) NOT NULL,
  last_block_number NUMERIC NULL,
  contract_address VARCHAR(255) NOT NULL,
  contract_address_id INT NULL,
  developer_address VARCHAR(255) NOT NULL,
  developer_address_id INT NULL,
  mining_asset_id INT NOT NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_chain FOREIGN KEY(chain_id) REFERENCES chains(id),
  CONSTRAINT fk_exchange FOREIGN KEY(exchange_id) REFERENCES exchanges(id),
  CONSTRAINT fk_contract_address FOREIGN KEY(contract_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_developer_address FOREIGN KEY(developer_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_mining_asset FOREIGN KEY(mining_asset_id) REFERENCES assets(id)
);

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-user";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-user";

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-api";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-api";

--create index
CREATE INDEX geth_miners_contract_address_id ON geth_miners(contract_address_id);
CREATE INDEX geth_miners_developer_address_id ON geth_miners(developer_address_id);
CREATE INDEX geth_miners_mining_asset ON geth_miners(mining_asset_id);
CREATE INDEX geth_miners_fk_chain ON geth_miners(chain_id);

DROP TABLE IF EXISTS geth_miners_transaction_inputs CASCADE;
CREATE TABLE geth_miners_transaction_inputs
(
  miner_id INT NOT NULL,
  transaction_input_id INT NOT NULL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(miner_id, transaction_input_id),
  CONSTRAINT fk_miner_id FOREIGN KEY(miner_id) REFERENCES geth_miners(id),
  CONSTRAINT fk_transaction_input_id FOREIGN KEY(transaction_input_id) REFERENCES geth_transaction_inputs(id)
);

CREATE INDEX geth_miners_transaction_inputs_miner_id ON geth_miners_transaction_inputs(miner_id);
CREATE INDEX geth_miners_transaction_inputs_transaction_input_id ON geth_miners_transaction_inputs(transaction_input_id);

DROP TABLE IF EXISTS geth_miners_transactions CASCADE;
CREATE TABLE geth_miners_transactions
(
  miner_id INT NOT NULL,
  transaction_id INT NOT NULL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(miner_id, transaction_id),
  CONSTRAINT fk_miner_id FOREIGN KEY(miner_id) REFERENCES geth_miners(id),
  CONSTRAINT fk_transaction_id FOREIGN KEY(transaction_id) REFERENCES geth_transaction_inputs(id)
);

CREATE INDEX geth_miners_transactions_miner_id ON geth_miners_transactions(miner_id);
CREATE INDEX geth_miners_transactions_transactiont_id ON geth_miners_transactions(transaction_id);