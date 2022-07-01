CREATE TABLE strategies
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4() ,
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  start_date timestamp NULL,
  end_date timestamp NULL,
  description TEXT NULL,
  strategy_type_id INT NOT NULL,
-- end 2022-03-06
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_strategy_type_id FOREIGN KEY(strategy_type_id) REFERENCES structured_values(id)
);