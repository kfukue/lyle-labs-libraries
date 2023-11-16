ROLLBACK
BEGIN TRANSACTION;
DROP TABLE IF EXISTS geth_trade_tax_transfers CASCADE;


CREATE TABLE geth_trade_tax_transfers
(
  geth_trade_id                 INT NOT NULL,
  geth_transfer_id              INT NOT NULL,
  tax_id                        INT NULL,
  uuid                          uuid NOT NULL DEFAULT uuid_generate_v4(),
  name                          VARCHAR(255) NOT NULL,
  alternate_name                VARCHAR(255) NULL,
  description                   TEXT NULL,
  created_by                    VARCHAR(255) NOT NULL,
  created_at                    timestamp NOT NULL,
  updated_by                    VARCHAR(255) NOT NULL,
  updated_at                    timestamp NOT NULL,
  PRIMARY KEY(geth_trade_id, geth_transfer_id),
  CONSTRAINT fk_geth_trade FOREIGN KEY(geth_trade_id) REFERENCES geth_trades(id),
  CONSTRAINT fk_geth_transfers FOREIGN KEY(geth_transfer_id) REFERENCES geth_transfers(id)
);


GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-user";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-user";

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-api";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-api";