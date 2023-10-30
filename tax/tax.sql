DROP TABLE IF EXISTS taxes CASCADE;

CREATE TABLE taxes
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4() ,
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  start_date timestamp NULL,
  end_date timestamp NULL,
  start_block INT NULL,
  end_block INT NULL,
  tax_rate NUMERIC NOT NULL,
  tax_rate_type_id INT NULL,
  contract_address_str VARCHAR(255) UNIQUE NULL,
  contract_address_id INT NULL,
  tax_type_id INT NOT NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_tax_rate_type_id FOREIGN KEY(tax_rate_type_id) REFERENCES structured_values(id),
  CONSTRAINT fk_contract_address_id FOREIGN KEY(contract_address_id) REFERENCES geth_addresses(id),
  CONSTRAINT fk_tax_type_id FOREIGN KEY(tax_type_id) REFERENCES structured_values(id)
);


GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-user";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-user";

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-api";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-api";