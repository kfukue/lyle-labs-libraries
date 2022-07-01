CREATE TABLE source_jobs
(
  source_id INT NOT NULL,
  job_id INT NOT NULL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(source_id, job_id),
  CONSTRAINT fk_source FOREIGN KEY(source_id) REFERENCES sources(id),
  CONSTRAINT fk_jobs FOREIGN KEY(job_id) REFERENCES jobs(id)
);
