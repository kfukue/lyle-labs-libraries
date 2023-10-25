ROLLBACK
BEGIN TRANSACTION;
DROP TABLE IF EXISTS geth_addresses CASCADE;

CREATE TABLE geth_addresses
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  description TEXT NULL,
  address_str VARCHAR(255) UNIQUE NOT NULL,
  address_type_id INT NOT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_address_type FOREIGN KEY(address_type_id) REFERENCES structured_values(id)
);

