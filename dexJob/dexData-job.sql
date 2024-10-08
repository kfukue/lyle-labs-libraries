CREATE TABLE dex_txn_jobs
(
  id INT NOT NULL,
  job_id INT NOT NULL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  start_date timestamp NOT NULL,
  end_date timestamp NOT NULL,
  description TEXT NULL,
  status_id int NOT NULL,
  chain_id INT NULL,
  exchange_id INT NULL,
  transaction_hashes VARCHAR[],
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_jobs FOREIGN KEY(job_id) REFERENCES jobs(id),
  CONSTRAINT status_id FOREIGN KEY(status_id) REFERENCES structured_values(id),
  CONSTRAINT fk_chains FOREIGN KEY(chain_id) REFERENCES chains(id),
  CONSTRAINT fk_exchanges FOREIGN KEY(exchange_id) REFERENCES exchanges(id)
);
