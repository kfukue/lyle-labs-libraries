

COMMIT
BEGIN TRANSACTION;
DROP TABLE IF EXISTS geth_transfers CASCADE;

CREATE TABLE geth_transfers
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  chain_id INT NOT NULL,
  token_address VARCHAR(255) NULL,
  token_address_id INT NULL,
  asset_id INT NULL,
  block_number NUMERIC NULL,
  index_number NUMERIC NULL,
  transfer_date timestamp NULL,
  txn_hash VARCHAR(70) NOT NULL,
  sender_address VARCHAR(255) NOT NULL,
  sender_address_id INT NULL,
  to_address VARCHAR(255) NOT NULL,
  to_address_id INT  NULL,
  amount NUMERIC NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  geth_process_job_id INT NULL,
  topics_str TEXT[] NULL,
  status_id INT NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_chains FOREIGN KEY(chain_id) REFERENCES chains(id),
  CONSTRAINT fk_assets FOREIGN KEY(asset_id) REFERENCES assets(id),
  CONSTRAINT fk_token_address FOREIGN KEY(token_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_sender_address FOREIGN KEY(sender_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_to_address FOREIGN KEY(to_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_geth_process_jobs FOREIGN KEY(geth_process_job_id) REFERENCES geth_process_jobs(id),
  CONSTRAINT fk_statuses FOREIGN KEY(status_id) REFERENCES structured_values(id)
);
GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-user";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-user";

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-api";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-api";

--create index
CREATE INDEX geth_transfers_to_address_id ON geth_transfers(to_address_id);
CREATE INDEX geth_transfers_sender_address_id ON geth_transfers(sender_address_id);
CREATE INDEX geth_transfers_to_address ON geth_transfers(to_address);
CREATE INDEX geth_transfers_sender_address ON geth_transfers(sender_address);
CREATE INDEX geth_transfers_asset_id ON geth_transfers(asset_id);

