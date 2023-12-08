ROLLBACK
BEGIN TRANSACTION;
DROP TABLE IF EXISTS geth_trade_swaps CASCADE;


CREATE TABLE geth_trade_swaps
(
  geth_trade_id                 INT NOT NULL,
  geth_swap_id                  INT NOT NULL,
  uuid                          uuid NOT NULL DEFAULT uuid_generate_v4(),
  name                          VARCHAR(255) NOT NULL,
  alternate_name                VARCHAR(255) NULL,
  description                   TEXT NULL,
  created_by                    VARCHAR(255) NOT NULL,
  created_at                    timestamp NOT NULL,
  updated_by                    VARCHAR(255) NOT NULL,
  updated_at                    timestamp NOT NULL,
  PRIMARY KEY(geth_trade_id, geth_swap_id),
  CONSTRAINT fk_geth_trade FOREIGN KEY(geth_trade_id) REFERENCES geth_trades(id),
  CONSTRAINT fk_geth_swap FOREIGN KEY(geth_swap_id) REFERENCES geth_swaps(id)
);


GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-user";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-user";

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-api";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-api";

-- create index
CREATE INDEX geth_trade_swaps_swap_id ON geth_trade_swaps(geth_swap_id);
CREATE INDEX geth_trade_swaps_trade_id ON geth_trade_swaps(geth_trade_id);