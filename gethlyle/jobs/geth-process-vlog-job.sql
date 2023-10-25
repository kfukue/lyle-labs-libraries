CREATE TABLE geth_process_vlog_jobs
(
  id SERIAL,
  geth_process_job_id INT NULL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  start_date timestamp NOT NULL,
  end_date timestamp NOT NULL,
  description TEXT NULL,
  status_id int NOT NULL,
  job_category_id INT NULL,
  asset_id INT NULL,
  chain_id INT NULL,
  txn_hash VARCHAR (255) NULL,
  address_id INT NULL,
  block_number NUMERIC NULL,
  index_number NUMERIC NULL,
  topics_str TEXT[] NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_geth_process_jobs FOREIGN KEY(geth_process_job_id) REFERENCES geth_process_jobs(id),
  CONSTRAINT fk_statuses FOREIGN KEY(status_id) REFERENCES structured_values(id),
  CONSTRAINT fk_job_categories FOREIGN KEY(job_category_id) REFERENCES structured_values(id),
  CONSTRAINT fk_assets FOREIGN KEY(asset_id) REFERENCES assets(id),
  CONSTRAINT fk_chains FOREIGN KEY(chain_id) REFERENCES chains(id),
  CONSTRAINT fk_addresses FOREIGN KEY(address_id) REFERENCES geth_addresses(id)
);

ALTER TABLE geth_process_vlog_jobs ALTER COLUMN topics_str TYPE TEXT[];

