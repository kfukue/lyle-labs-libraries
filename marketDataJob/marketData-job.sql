CREATE TABLE market_data_jobs
(
  market_data_id INT NOT NULL,
  job_id INT NOT NULL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  start_date timestamp NOT NULL,
  end_date timestamp NOT NULL,
  description TEXT NULL,
  status_id int NOT NULL,
  response_status TEXT,
  request_url TEXT NULL,
  request_body TEXT NULL,
  request_method TEXT NULL,
  response_data TEXT NULL,
  response_data_json jsonb NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(market_data_id, job_id),
  CONSTRAINT fk_market_data FOREIGN KEY(market_data_id) REFERENCES market_data(id),
  CONSTRAINT fk_jobs FOREIGN KEY(job_id) REFERENCES jobs(id),
  CONSTRAINT status_id FOREIGN KEY(status_id) REFERENCES structured_values(id)
);
