CREATE TABLE chains
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  address VARCHAR(255) NULL,
  chain_type_id INT NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  chain_id INT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_chain_type FOREIGN KEY(chain_type_id) REFERENCES structured_values(id)
);

