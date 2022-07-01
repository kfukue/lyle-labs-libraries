CREATE TABLE steps
(
  id SERIAL,
  pool_id  INT NOT NULL,
  parent_step_id INT NULL,
  name VARCHAR(255) NOT NULL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4() ,
  alternate_name VARCHAR(255) NULL,
  start_date timestamp NULL,
  end_date timestamp NULL,
  description TEXT NULL,
  action_type_id INT NULL,
  function_name VARCHAR(255) NOT NULL,
  step_order INT NOT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_pool_id FOREIGN KEY(pool_id) REFERENCES pools(id),
  CONSTRAINT fk_parent_step_id FOREIGN KEY(parent_step_id) REFERENCES steps(id),
  CONSTRAINT action_type_id FOREIGN KEY(action_type_id) REFERENCES structured_values(id)
);
