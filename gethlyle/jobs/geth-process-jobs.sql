
COMMIT
BEGIN TRANSACTION;
DROP TABLE IF EXISTS geth_process_jobs CASCADE;

CREATE TABLE geth_process_jobs
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  start_date timestamp NOT NULL,
  end_date timestamp  NULL,
  description TEXT NULL,
  status_id int NOT NULL,
  job_category_id INT NULL,
  import_type_id INT NULL,
  chain_id INT NULL,
  start_block_number INT  NULL,
  end_block_number INT  NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  asset_id INT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_statuses FOREIGN KEY(status_id) REFERENCES structured_values(id),
  CONSTRAINT fk_job_categories FOREIGN KEY(job_category_id) REFERENCES structured_values(id),
  CONSTRAINT fk_import_type FOREIGN KEY(import_type_id) REFERENCES structured_values(id),
  CONSTRAINT fk_chains FOREIGN KEY(chain_id) REFERENCES chains(id)
);

-- new column 2023-09-13
ROLLBACK
START TRANSACTION;
ALTER TABLE geth_process_jobs
  ADD  asset_id INT NULL
  COMMIT
-- end

-- create index
CREATE INDEX geth_process_jobs_asset_id ON geth_process_jobs(asset_id);
CREATE INDEX geth_process_jobs_start_date ON geth_process_jobs(start_date);
CREATE INDEX geth_process_jobs_end_date ON geth_process_jobs(end_date);
CREATE INDEX geth_process_jobs_status_id ON geth_process_jobs(status_id);
CREATE INDEX geth_process_jobs_import_type_id ON geth_process_jobs(import_type_id);