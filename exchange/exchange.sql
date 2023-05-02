CREATE TABLE exchanges
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4() ,
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  exchange_type_id INT NULL,
  url varchar(255)NULL,
  start_date timestamp NULL,
  end_date timestamp NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_structured_value_exchange_type FOREIGN KEY(exchange_type_id) REFERENCES structured_values(id)
);

CREATE TABLE exchange_chains
(
  uuid uuid NOT NULL DEFAULT uuid_generate_v4() ,
  exchange_id INT NOT NULL,
  chain_id  INT NOT NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(exchange_id, chain_id),
  CONSTRAINT fk_exchanged_id FOREIGN KEY(exchange_id) REFERENCES exchanges(id),
  CONSTRAINT fk_chain_id FOREIGN KEY(chain_id) REFERENCES chains(id)
);
