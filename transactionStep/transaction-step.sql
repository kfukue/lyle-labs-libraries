CREATE TABLE transaction_steps
(
  transaction_id INT NOT NULL,
  step_id INT NOT NULL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(transaction_id, step_id),
  CONSTRAINT fk_transaction FOREIGN KEY(transaction_id) REFERENCES transactions(id),
  CONSTRAINT fk_step FOREIGN KEY(step_id) REFERENCES steps(id)
);
