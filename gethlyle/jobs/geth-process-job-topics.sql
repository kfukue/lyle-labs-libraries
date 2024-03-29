CREATE TABLE geth_process_job_topics
(
  id SERIAL,
  geth_process_job_id INT NULL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  description TEXT NULL,
  status_id int NOT NULL,
  topic_str TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_geth_process_jobs FOREIGN KEY(geth_process_job_id) REFERENCES geth_process_jobs(id),
  CONSTRAINT fk_statuses FOREIGN KEY(status_id) REFERENCES structured_values(id)
);
-- create index
CREATE INDEX geth_process_job_topics_geth_process_job_id ON geth_process_job_topics(geth_process_job_id);
CREATE INDEX geth_process_job_topics_topic_str ON geth_process_job_topics(topic_str);